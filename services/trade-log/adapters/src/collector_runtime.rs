use crate::postgres::{Postgres, checkpoint::CollectionCommit};
use std::{sync::Arc, time::Duration};
use tokio_util::sync::CancellationToken;
use trade_log::{
    acquisition::SourceReader,
    checkpoint::{AcceptedPage, Lease, PendingWork},
    collection::{CollectionConfig, error},
    parsing::ProtocolParser,
    persistence::CollectionFactStore,
    query::QueryError,
};

pub struct CollectorRuntime {
    pub store: Postgres,
    pub config: CollectionConfig,
    pub source: Arc<dyn SourceReader>,
    pub parser: Arc<dyn ProtocolParser>,
}
impl CollectorRuntime {
    pub async fn round(&self, l: &Lease, work: &mut PendingWork) -> Result<bool, QueryError> {
        for _ in 0..self.config.max_requests_per_round {
            let (kind, range) = if work.meta.is_none() {
                ("meta", None)
            } else if work.spot_meta.is_none() {
                ("spotMeta", None)
            } else if let Some(r) = work.remaining.last() {
                ("userFillsByTime", Some(r.clone()))
            } else {
                break;
            };
            let response = match &range {
                Some(r) => {
                    self.source
                        .fetch_range(&self.config.account, r.start_ms, r.end_ms)
                        .await?
                }
                None => {
                    self.source
                        .fetch(
                            if kind == "meta" {
                                protocol_api::QueryKind::Meta
                            } else {
                                protocol_api::QueryKind::SpotMeta
                            },
                            &self.config.account,
                        )
                        .await?
                }
            };
            let raw = self
                .store
                .save_collection_response(l, work, kind, range.as_ref(), &response)
                .await?;
            if !(200..300).contains(&response.status) {
                if response.status == 429 || response.status >= 500 {
                    let mut e = if response.status == 429 {
                        QueryError::limited()
                    } else {
                        QueryError::unavailable("Source returned an unsuccessful response")
                    };
                    e.retryable = true;
                    return Err(e);
                }
                return Err(error(
                    "INCOMPLETE_DATA",
                    "Source rejected the collection request",
                    false,
                ));
            }
            if let Some(range) = range {
                let fills: Vec<serde_json::Value> = serde_json::from_slice(&response.body)
                    .map_err(|_| QueryError::incomplete("Invalid fills response"))?;
                // Reject unexpected range data before allowing a checkpoint to advance.
                for fill in &fills {
                    let at = fill["time"]
                        .as_i64()
                        .ok_or_else(|| QueryError::incomplete("Missing fill timestamp"))?;
                    if at < range.start_ms || at >= range.end_ms {
                        return Err(QueryError::incomplete(
                            "Source fill outside requested range",
                        ));
                    }
                }
                work.remaining.pop();
                if fills.len() >= 2000 {
                    let (left, right) = range.split()?;
                    work.remaining.push(right);
                    work.remaining.push(left);
                } else {
                    work.pages.push(AcceptedPage { range, raw_id: raw });
                }
            } else if kind == "meta" {
                work.meta = Some(raw);
            } else {
                work.spot_meta = Some(raw);
            }
            self.store.save_work(l, work).await?;
        }
        if !work.remaining.is_empty() || work.meta.is_none() || work.spot_meta.is_none() {
            return Ok(false);
        }
        let ids: Vec<String> = work.pages.iter().map(|p| p.raw_id.clone()).collect();
        let (result, observations) = self
            .store
            .collection_result(
                &work.query_id,
                &self.config.account,
                &ids,
                work.meta.as_deref().expect("meta"),
                work.spot_meta.as_deref().expect("spot meta"),
                self.parser.as_ref(),
            )
            .await?;
        self.store
            .persist_collection(
                &result,
                &CollectionCommit {
                    lease: l,
                    work,
                    observations: &observations,
                    interval_seconds: self.config.interval_seconds,
                },
            )
            .await?;
        tracing::info!(service="trade-collector",account_key=%l.key,query_id=%work.query_id,lease_epoch=l.epoch,start_ms=work.range.start_ms,end_ms=work.range.end_ms,facts=result.trades.len(),"collection_completed");
        Ok(true)
    }
    pub async fn run(&self, cancel: CancellationToken) {
        self.supervise(cancel, None).await;
    }
    pub async fn run_with_realtime(
        &self,
        cancel: CancellationToken,
        config: trade_log::realtime::WebsocketConfig,
        source: Arc<dyn trade_log::realtime::RealtimeSource>,
    ) {
        let stream = config.enabled.then_some((config, source));
        self.supervise(cancel, stream).await;
    }
    async fn supervise(
        &self,
        cancel: CancellationToken,
        stream: Option<(
            trade_log::realtime::WebsocketConfig,
            Arc<dyn trade_log::realtime::RealtimeSource>,
        )>,
    ) {
        let instance = uuid::Uuid::new_v4();
        let reconnect_schedule = Arc::new(tokio::sync::Mutex::new(
            trade_log::realtime::ReconnectSchedule::default(),
        ));
        let clock_origin = tokio::time::Instant::now();
        loop {
            if cancel.is_cancelled() {
                break;
            }
            let acquired = tokio::select! { _=cancel.cancelled()=>break, result=self.store.acquire_collection(&self.config,instance)=>result };
            if let Ok(Some(lease)) = acquired {
                let ready = tokio::select! { _=cancel.cancelled()=>false,result=self.store.prepare_stream(&lease,&self.config,stream.is_some())=>result.is_ok() };
                if ready {
                    let workers = cancel.child_token();
                    let http = self.run_http(&lease, workers.clone());
                    let ws = async {
                        if let Some((config, source)) = &stream {
                            crate::websocket_runtime::WebsocketRuntime {
                                store: self.store.clone(),
                                collection: self.config.clone(),
                                config: config.clone(),
                                source: source.clone(),
                                http: self.source.clone(),
                                parser: self.parser.clone(),
                                schedule: reconnect_schedule.clone(),
                                clock_origin,
                            }
                            .run(&lease, workers.clone())
                            .await;
                        } else {
                            workers.cancelled().await;
                        }
                    };
                    tokio::pin!(http, ws);
                    let mut heartbeat = tokio::time::interval(Duration::from_secs(
                        (self.config.lease_seconds / 3).max(1),
                    ));
                    heartbeat.tick().await;
                    let mut http_done = false;
                    let mut ws_done = false;
                    loop {
                        tokio::select! {
                            biased;
                            _=cancel.cancelled()=>break,
                            _=&mut http=>{http_done=true;break;},
                            _=&mut ws=>{ws_done=true;break;},
                            _=heartbeat.tick()=>{
                                let renewed=tokio::select! { _=cancel.cancelled()=>false,result=self.store.renew_collection(&lease,self.config.lease_seconds)=>result.is_ok() };
                                if !renewed { break; }
                            }
                        }
                    }
                    workers.cancel();
                    if !http_done {
                        http.await;
                    }
                    if !ws_done {
                        ws.await;
                    }
                }
                let _ = tokio::time::timeout(
                    Duration::from_secs(2),
                    self.store.release_collection(&lease),
                )
                .await;
            }
            tokio::select! { _=cancel.cancelled()=>break, _=tokio::time::sleep(Duration::from_secs(1))=>{} }
        }
    }
    async fn run_http(&self, l: &Lease, cancel: CancellationToken) {
        loop {
            let available = tokio::select! { _=cancel.cancelled()=>break,result=self.store.collection_work(l,&self.config,chrono::Utc::now().timestamp_millis())=>result };
            match available {
                Ok(Some(mut work)) => {
                    let outcome = tokio::select! {
                        _=cancel.cancelled()=>break,
                        value=tokio::time::timeout(Duration::from_secs(self.config.round_timeout_seconds),self.round(l,&mut work))=>value.unwrap_or_else(|_|Err(error("ROUND_TIMEOUT","Collection round timed out",true)))
                    };
                    let (e, delay, budget) = match outcome {
                        Ok(true) => {
                            continue;
                        }
                        Ok(false) => (
                            error(
                                "ROUND_BUDGET_EXHAUSTED",
                                "Collection request budget exhausted",
                                true,
                            ),
                            self.config.interval_seconds,
                            true,
                        ),
                        Err(mut e) => {
                            if e.code == "LEASE_LOST" {
                                break;
                            }
                            if e.code == "DEPENDENCY_UNAVAILABLE" {
                                e.retryable = true;
                            }
                            let failures=sqlx::query_scalar::<_,i32>("SELECT consecutive_failures FROM trade_log.collection_checkpoints WHERE partition_key=$1").bind(&l.key).fetch_one(&self.store.pool).await.unwrap_or(0);
                            let delay = self
                                .config
                                .retry_base_seconds
                                .saturating_mul(1_u64 << failures.min(20))
                                .min(self.config.retry_max_seconds);
                            let source_delay = self
                                .store
                                .collection_retry_after(&work.query_id)
                                .await
                                .unwrap_or(0);
                            (e, delay.max(source_delay), false)
                        }
                    };
                    tracing::warn!(service="trade-collector",account_key=%l.key,query_id=%work.query_id,code=%e.code,"collection_failed");
                    let recorded = tokio::select! { _=cancel.cancelled()=>break,result=self.store.collection_failure(l,&e,delay,budget)=>result };
                    if recorded.is_err() {
                        break;
                    }
                }
                Ok(None) => {}
                Err(e) => {
                    if e.code == "LEASE_LOST" {
                        break;
                    }
                    let recorded = tokio::select! { _=cancel.cancelled()=>break,result=self.store.collection_failure(l,&e,self.config.retry_base_seconds,false)=>result };
                    if recorded.is_err() {
                        break;
                    }
                }
            }
            tokio::select! { _=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_secs(1))=>{} }
        }
    }
}
