use crate::postgres::{Postgres, realtime::StreamCommit};
use std::{sync::Arc, time::Duration};
use tokio::{
    sync::{OwnedSemaphorePermit, Semaphore, mpsc},
    time::Instant,
};
use tokio_util::sync::CancellationToken;
use trade_log::{
    acquisition::SourceReader,
    checkpoint::Lease,
    collection::CollectionConfig,
    parsing::ProtocolParser,
    query::QueryError,
    realtime::{
        RealtimeSource, ReconnectSchedule, StreamEvent, StreamMessage, WebsocketConfig,
        stream_error,
    },
};
use uuid::Uuid;

pub struct WebsocketRuntime {
    pub store: Postgres,
    pub collection: CollectionConfig,
    pub config: WebsocketConfig,
    pub source: Arc<dyn RealtimeSource>,
    pub http: Arc<dyn SourceReader>,
    pub parser: Arc<dyn ProtocolParser>,
    pub schedule: Arc<tokio::sync::Mutex<ReconnectSchedule>>,
    pub clock_origin: Instant,
}
struct Queued {
    sequence: i64,
    message: StreamMessage,
    rejected: Option<QueryError>,
    _bytes: OwnedSemaphorePermit,
}
struct Metadata {
    meta: Vec<u8>,
    spot: Vec<u8>,
}
impl WebsocketRuntime {
    async fn metadata(&self) -> Result<Metadata, QueryError> {
        let meta = self
            .http
            .fetch(protocol_api::QueryKind::Meta, &self.collection.account)
            .await?;
        let spot = self
            .http
            .fetch(protocol_api::QueryKind::SpotMeta, &self.collection.account)
            .await?;
        if !(200..300).contains(&meta.status) || !(200..300).contains(&spot.status) {
            return Err(stream_error(
                "METADATA_UNAVAILABLE",
                "Market metadata unavailable",
            ));
        }
        for body in [&meta.body, &spot.body] {
            let value: serde_json::Value = serde_json::from_slice(body)
                .map_err(|_| stream_error("METADATA_UNAVAILABLE", "Invalid market metadata"))?;
            if !value["universe"].is_array() {
                return Err(stream_error(
                    "METADATA_UNAVAILABLE",
                    "Invalid market metadata",
                ));
            }
        }
        Ok(Metadata {
            meta: meta.body,
            spot: spot.body,
        })
    }
    async fn bounded<T>(
        &self,
        future: impl std::future::Future<Output = Result<T, QueryError>>,
    ) -> Result<T, QueryError> {
        tokio::time::timeout(
            Duration::from_secs(self.config.commit_timeout_seconds),
            future,
        )
        .await
        .unwrap_or_else(|_| Err(stream_error("COMMIT_TIMEOUT", "Stream operation timed out")))
    }
    async fn recover(&self, l: &Lease, metadata: &Metadata) -> Result<(), QueryError> {
        if let Some(e) = self.store.stream_fatal_error(l).await? {
            return Err(e);
        }
        loop {
            let pending = self.store.pending_stream_jobs(l).await?;
            if pending.is_empty() {
                return Ok(());
            }
            for (id, session, seq) in pending {
                self.store.resume_stream_job(l, &id).await?;
                let result = match self.store.stream_result(&id, self.parser.as_ref()).await {
                    Err(e) if e.code == "INCOMPLETE_DATA" => {
                        self.store
                            .update_stream_metadata(l, &id, &metadata.meta, &metadata.spot)
                            .await?;
                        self.store.stream_result(&id, self.parser.as_ref()).await?
                    }
                    result => result?,
                };
                self.bounded(self.store.commit_stream(
                    &result,
                    &StreamCommit {
                        lease: l,
                        session_id: session,
                        sequence: seq,
                    },
                ))
                .await?;
            }
        }
    }
    pub async fn run(&self, l: &Lease, cancel: CancellationToken) {
        let origin = self.clock_origin;
        let mut schedule = self.schedule.lock().await;
        let mut wait = schedule.pending_wait_ms(origin.elapsed().as_millis() as u64);
        let mut fatal = false;
        loop {
            if cancel.is_cancelled() {
                break;
            }
            if fatal {
                cancel.cancelled().await;
                break;
            }
            wait = wait.max(schedule.rate_wait_ms(origin.elapsed().as_millis() as u64));
            if wait > 0 {
                let patch = serde_json::json!({"status":"RECONNECT_WAIT","next_retry_at":trade_log::collection::time(chrono::Utc::now().timestamp_millis().saturating_add(wait as i64))});
                let _ = self.bounded(self.store.stream_state(l, patch)).await;
                tokio::select! {_=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_millis(wait))=>{}}
            }
            // Verify archive access and process durable pending messages before a new socket.
            let prepared = tokio::select! {_=cancel.cancelled()=>break,result=async {self.bounded(self.store.stream_ready(l)).await?;let metadata=self.metadata().await?;self.recover(l,&metadata).await.map_err(|mut e|{if e.code=="INCOMPLETE_DATA" {e.retryable=true;} e})?;Ok::<_,QueryError>(metadata)}=>result};
            let metadata = match prepared {
                Ok(v) => v,
                Err(e) => {
                    tracing::warn!(code=%e.code,"websocket_prepare_failed");
                    if e.code == "DEPENDENCY_UNAVAILABLE"
                        || e.code == "LEASE_LOST"
                        || e.code == "COMMIT_TIMEOUT"
                    {
                        return;
                    }
                    let _=self.store.stream_state(l,serde_json::json!({"status":if e.retryable{"RECONNECT_WAIT"}else{"FAILED"},"last_error":{"code":e.code,"message":e.message,"occurred_at":shared_types::now()}})).await;
                    fatal = !e.retryable;
                    wait = schedule.wait_ms(
                        origin.elapsed().as_millis() as u64,
                        &self.config,
                        random(),
                    );
                    continue;
                }
            };
            let session = Uuid::new_v4();
            if self
                .bounded(self.store.begin_stream(l, &self.collection, session))
                .await
                .is_err()
            {
                return;
            }
            schedule.attempted(origin.elapsed().as_millis() as u64);
            let connected = tokio::select! {_=cancel.cancelled()=>Err(stream_error("STOPPED","Collector stopped")),result=tokio::time::timeout(Duration::from_secs(self.config.connect_timeout_seconds),self.source.connect(&self.collection.account,self.config.max_message_bytes))=>result.unwrap_or_else(|_|Err(stream_error("CONNECT_TIMEOUT","Websocket connect timed out")))};
            let result = match connected {
                Ok(mut socket) => {
                    let outcome=async{
                        self.bounded(self.store.stream_connected(l,session)).await?;
                        let (sender,receiver)=mpsc::channel(self.config.max_pending_messages);
                        let bytes=Arc::new(Semaphore::new(self.config.max_pending_bytes));
                        let reader=self.receive(l,session,socket.as_mut(),(sender,bytes.clone()),&mut schedule,&cancel);
                        let consumer=self.consume(l,session,receiver,bytes,metadata,&cancel);
                        tokio::pin!(reader,consumer);
                        // Both futures are owned here: dropping either cancels all in-flight work.
                        tokio::select!{result=&mut reader=>result,result=&mut consumer=>result,_=cancel.cancelled()=>Err(stream_error("STOPPED","Collector stopped"))}
                    }.await;
                    let _ = tokio::time::timeout(Duration::from_secs(1), socket.close()).await;
                    outcome
                }
                Err(e) => Err(e),
            };
            let e = result
                .err()
                .unwrap_or_else(|| stream_error("DISCONNECT", "Websocket worker stopped"));
            tracing::warn!(service="trade-collector",session_id=%session,lease_epoch=l.epoch,code=%e.code,"websocket_stopped");
            let recorded = self
                .bounded(self.store.end_stream(
                    l,
                    &self.collection,
                    session,
                    &e,
                    cancel.is_cancelled(),
                ))
                .await;
            wait = schedule.wait_ms(origin.elapsed().as_millis() as u64, &self.config, random());
            if recorded.is_err() || e.code == "LEASE_LOST" || e.code == "DEPENDENCY_UNAVAILABLE" {
                return;
            }
            fatal = !e.retryable;
        }
    }
    async fn receive(
        &self,
        l: &Lease,
        id: Uuid,
        socket: &mut dyn trade_log::realtime::RealtimeConnection,
        queue: (mpsc::Sender<Queued>, Arc<Semaphore>),
        schedule: &mut ReconnectSchedule,
        cancel: &CancellationToken,
    ) -> Result<(), QueryError> {
        let (sender, bytes) = queue;
        let deadline = Instant::now() + Duration::from_secs(self.config.subscribe_timeout_seconds);
        let mut subscribed: Option<Instant> = None;
        let mut sequence = 0_i64;
        let mut ping =
            tokio::time::interval(Duration::from_secs(self.config.ping_interval_seconds));
        ping.set_missed_tick_behavior(tokio::time::MissedTickBehavior::Delay);
        ping.tick().await;
        let mut pong_due: Option<Instant> = None;
        let mut heartbeat_ok = false;
        loop {
            let timeout = pong_due
                .or_else(|| subscribed.is_none().then_some(deadline))
                .unwrap_or_else(|| Instant::now() + Duration::from_secs(86400));
            tokio::select! {
                biased;
                _=cancel.cancelled()=>return Err(stream_error("STOPPED","Collector stopped")),
                _=tokio::time::sleep_until(timeout)=>return Err(stream_error(if subscribed.is_none(){"SUBSCRIBE_TIMEOUT"}else{"HEARTBEAT_TIMEOUT"},"Websocket confirmation or pong timed out")),
                _=ping.tick(),if subscribed.is_some() && pong_due.is_none()=>{
                    tokio::time::timeout(Duration::from_secs(self.config.pong_timeout_seconds),socket.ping()).await.map_err(|_|stream_error("WS_IO_ERROR","Websocket write timed out"))??;
                    pong_due=Some(Instant::now()+Duration::from_secs(self.config.pong_timeout_seconds));
                },
                event=socket.next()=>match event? {
                    StreamEvent::Subscribed=>{
                        if subscribed.is_none(){self.bounded(self.store.stream_subscribed(l,id)).await?;subscribed=Some(Instant::now());tracing::info!(session_id=%id,"websocket_subscribed");}
                    },
                    StreamEvent::Pong=>{
                        if subscribed.is_some() && pong_due.is_some(){pong_due=None;heartbeat_ok=true;self.bounded(self.store.stream_pong(l,id)).await?;}
                    },
                    StreamEvent::Ignore=>{},
                    data=>{
                        let (message,rejected)=match data {StreamEvent::Data(m)=>(m,None),StreamEvent::Rejected(m,e)=>(m,Some(e)),_=>unreachable!()};
                        if subscribed.is_none(){return Err(stream_error("SUBSCRIPTION_REJECTED","Data before subscription confirmation"));}
                        sequence=sequence.checked_add(1).ok_or_else(QueryError::storage)?;
                        // Both envelope and extracted fills occupy memory until the consumer finishes.
                        let length=message.body.len().saturating_add(message.fills.len());
                        let permit=bytes.clone().try_acquire_many_owned(u32::try_from(length).map_err(|_|stream_error("MESSAGE_TOO_LARGE","Message too large"))?).map_err(|_|stream_error("QUEUE_OVERFLOW","Websocket byte queue is full"))?;
                        sender.try_send(Queued{sequence,message,rejected,_bytes:permit}).map_err(|_|stream_error("QUEUE_OVERFLOW","Websocket message queue is full"))?;
                        self.bounded(self.store.stream_state(l,serde_json::json!({"pending_messages":self.config.max_pending_messages-sender.capacity(),"pending_bytes":self.config.max_pending_bytes-bytes.available_permits()}))).await?;
                    }
                }
            }
            if let Some(at) = subscribed {
                schedule.stable(
                    at.elapsed().as_secs(),
                    heartbeat_ok && pong_due.is_none(),
                    &self.config,
                );
            }
        }
    }
    async fn consume(
        &self,
        l: &Lease,
        id: Uuid,
        mut receiver: mpsc::Receiver<Queued>,
        bytes: Arc<Semaphore>,
        mut metadata: Metadata,
        cancel: &CancellationToken,
    ) -> Result<(), QueryError> {
        let mut refresh =
            tokio::time::interval(Duration::from_secs(self.config.metadata_refresh_seconds));
        refresh.set_missed_tick_behavior(tokio::time::MissedTickBehavior::Delay);
        refresh.tick().await;
        loop {
            let item = tokio::select! {
                _=cancel.cancelled()=>return Err(stream_error("STOPPED","Collector stopped")),
                _=refresh.tick()=>{
                    match self.metadata().await {
                        Ok(m)=>{metadata=m;self.bounded(self.store.stream_state(l,serde_json::json!({"metadata_stale":false}))).await?;},
                        Err(e)=>{tracing::warn!(code=%e.code,"metadata_stale");self.bounded(self.store.stream_state(l,serde_json::json!({"metadata_stale":true}))).await?;}
                    }
                    continue;
                },
                item=receiver.recv()=>item.ok_or_else(||stream_error("DISCONNECT","Stream queue closed"))?
            };
            let query = self
                .bounded(self.store.archive_stream(
                    l,
                    id,
                    item.sequence,
                    &item.message,
                    &metadata.meta,
                    &metadata.spot,
                ))
                .await?;
            let processed=async{
                if let Some(e)=&item.rejected{return Err(e.clone());}
                let parsed=self.store.stream_result(&query,self.parser.as_ref()).await;
                let result=match parsed {
                    Err(e) if e.code=="INCOMPLETE_DATA"=>{
                        // Unknown markets keep their archived evidence until metadata can be refreshed.
                        metadata=self.metadata().await?;
                        self.store.update_stream_metadata(l,&query,&metadata.meta,&metadata.spot).await?;
                        self.store.stream_result(&query,self.parser.as_ref()).await.or(Err(e))?
                    },
                    other=>other?
                };
                self.bounded(self.store.commit_stream(&result,&StreamCommit{lease:l,session_id:id,sequence:item.sequence})).await?;
                tracing::info!(service="trade-collector",session_id=%id,message_sequence=item.sequence,query_id=%query,message_mode=%item.message.mode,lease_epoch=l.epoch,transport="WEBSOCKET",received_at=%item.message.received_at,facts=result.trades.len(),"stream_message_committed");
                Ok(())
            }.await;
            if let Err(mut e) = processed {
                if e.code == "DEPENDENCY_UNAVAILABLE" || e.code == "INCOMPLETE_DATA" {
                    e.retryable = true;
                }
                let _ = self.store.fail_stream(l, &query, &e).await;
                return Err(e);
            }
            drop(item);
            self.bounded(self.store.stream_state(l,serde_json::json!({"pending_messages":receiver.len(),"pending_bytes":self.config.max_pending_bytes-bytes.available_permits()}))).await?;
        }
    }
}
fn random() -> u64 {
    u64::from_le_bytes(
        Uuid::new_v4().as_bytes()[..8]
            .try_into()
            .expect("uuid bytes"),
    )
}
