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
        let instance = uuid::Uuid::new_v4();
        let mut lease: Option<Lease> = None;
        loop {
            if cancel.is_cancelled() {
                break;
            }
            if lease.is_none() {
                match self.store.acquire_collection(&self.config, instance).await {
                    Ok(value) => lease = value,
                    Err(e) => {
                        tracing::warn!(service="trade-collector",code=%e.code,"collection_lease_unavailable")
                    }
                }
            }
            if let Some(l) = lease.clone() {
                if self
                    .store
                    .renew_collection(&l, self.config.lease_seconds)
                    .await
                    .is_err()
                {
                    lease = None;
                } else {
                    match self
                        .store
                        .collection_work(&l, &self.config, chrono::Utc::now().timestamp_millis())
                        .await
                    {
                        Ok(Some(mut work)) => {
                            let outcome = {
                                let deadline = tokio::time::sleep(Duration::from_secs(
                                    self.config.round_timeout_seconds,
                                ));
                                tokio::pin!(deadline);
                                let operation = self.round(&l, &mut work);
                                tokio::pin!(operation);
                                let mut heartbeat = tokio::time::interval(Duration::from_secs(
                                    (self.config.lease_seconds / 3).max(1),
                                ));
                                heartbeat.tick().await;
                                loop {
                                    tokio::select! {
                                     biased;
                                     _=cancel.cancelled()=>{break None;}
                                     _=&mut deadline=>{break Some(Err(error("ROUND_TIMEOUT","Collection round timed out",true)));}
                                     result=&mut operation=>{break Some(result);}
                                     _=heartbeat.tick()=>{if self.store.renew_collection(&l,self.config.lease_seconds).await.is_err(){lease=None;break None;}}
                                    }
                                }
                            };
                            if let Some(result) = outcome {
                                match result {
                                    Ok(true) => {}
                                    Ok(false) => {
                                        let e = error(
                                            "ROUND_BUDGET_EXHAUSTED",
                                            "Collection request budget exhausted",
                                            true,
                                        );
                                        let _ = self
                                            .store
                                            .collection_failure(
                                                &l,
                                                &e,
                                                self.config.interval_seconds,
                                                true,
                                            )
                                            .await;
                                    }
                                    Err(mut e) => {
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
                                        let delay = delay
                                            .max(source_delay)
                                            .min(self.config.retry_max_seconds);
                                        let jitter = u64::from(instance.as_bytes()[0]) % 3;
                                        tracing::warn!(service="trade-collector",account_key=%l.key,query_id=%work.query_id,code=%e.code,"collection_failed");
                                        if self
                                            .store
                                            .collection_failure(
                                                &l,
                                                &e,
                                                delay
                                                    .saturating_add(jitter)
                                                    .min(self.config.retry_max_seconds),
                                                false,
                                            )
                                            .await
                                            .is_err()
                                        {
                                            lease = None;
                                        }
                                    }
                                }
                            }
                        }
                        Ok(None) => {}
                        Err(e) => {
                            let _ = self
                                .store
                                .collection_failure(&l, &e, self.config.retry_base_seconds, false)
                                .await;
                            if e.code == "LEASE_LOST" {
                                lease = None;
                            }
                        }
                    }
                }
            }
            tokio::select! {_=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_secs(1))=>{}}
        }
        if let Some(l) = lease {
            let _ = tokio::time::timeout(Duration::from_secs(2), self.store.release_collection(&l))
                .await;
        }
    }
}
