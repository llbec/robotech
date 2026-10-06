use crate::postgres::Postgres;
use std::sync::Arc;
use std::time::Duration;
use tokio_util::sync::CancellationToken;
use trade_log::publishing::PublicationTarget;
pub struct PublisherRuntime {
    pub store: Postgres,
    pub key: String,
    pub webhook: Option<Arc<dyn PublicationTarget>>,
    pub lease_seconds: u64,
    pub poll_interval_ms: u64,
    pub retry_base_seconds: u64,
    pub retry_max_seconds: u64,
}
impl PublisherRuntime {
    pub async fn run(self, cancel: CancellationToken) {
        let owner = uuid::Uuid::new_v4();
        let heartbeat_store = self.store.clone();
        let heartbeat_key = self.key.clone();
        let heartbeat_cancel = cancel.clone();
        let heartbeat = tokio::spawn(async move {
            loop {
                tokio::select! {_=heartbeat_cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_secs(2))=>{let _=heartbeat_store.touch_publishing_heartbeat(&heartbeat_key).await;}}
            }
        });
        let mut had_error = false;
        loop {
            if cancel.is_cancelled() {
                break;
            }
            let result=async {
                if let Some(webhook)=&self.webhook && let Some(claim)=self.store.claim_publication(&self.key,owner,self.lease_seconds).await? {
                    let delivery=webhook.deliver(&claim.event_id,claim.attempt,&claim.body,self.retry_base_seconds,self.retry_max_seconds).await;
                    self.store.finish_publication(&claim,delivery.status,delivery.http,delivery.error,delivery.delay).await?;
                    tracing::info!(service="trade-parser-publisher",event_id=%claim.event_id,attempt=claim.attempt,status=delivery.status,http_status=delivery.http,"publication_completed");
                }
                Ok::<(),trade_log::query::QueryError>(())
            }.await;
            if let Err(e) = result {
                had_error = true;
                tracing::warn!(service="trade-parser-publisher",code=%e.code,"publication_failed");
                let _ = self
                    .store
                    .publishing_heartbeat(&self.key, Some(&e.code))
                    .await;
            } else if had_error
                && self
                    .store
                    .publishing_heartbeat(&self.key, None)
                    .await
                    .is_ok()
            {
                had_error = false;
            }
            tokio::select! {_=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_millis(self.poll_interval_ms))=>{}}
        }
        cancel.cancel();
        let _ = heartbeat.await;
    }
}
