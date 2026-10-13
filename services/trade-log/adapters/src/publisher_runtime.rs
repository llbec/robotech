use crate::postgres::Postgres;
use std::sync::Arc;
use std::time::Duration;
use tokio_util::sync::CancellationToken;
use trade_log::publishing::PublicationTarget;
pub struct PublisherRuntime {
    pub store: Postgres,
    pub webhook: Option<Arc<dyn PublicationTarget>>,
    pub lease_seconds: u64,
    pub poll_interval_ms: u64,
    pub retry_base_seconds: u64,
    pub retry_max_seconds: u64,
}
impl PublisherRuntime {
    pub async fn run_multi(
        self,
        cancel: CancellationToken,
        url: String,
        enabled: bool,
        policy: trade_log::publishing::CandidatePolicy,
        max_inflight: usize,
    ) {
        self.run_multi_with_spot(cancel, url, enabled, policy, max_inflight, false)
            .await
    }
    pub async fn run_multi_with_spot(
        self,
        cancel: CancellationToken,
        url: String,
        enabled: bool,
        policy: trade_log::publishing::CandidatePolicy,
        max_inflight: usize,
        spot_enabled: bool,
    ) {
        let semaphore = Arc::new(tokio::sync::Semaphore::new(max_inflight));
        let mut workers = DeliveryWorkers(std::collections::BTreeMap::new());
        let mut rotate = 0_usize;
        loop {
            if cancel.is_cancelled() {
                break;
            }
            let finished: Vec<String> = workers
                .0
                .iter()
                .filter(|(_, h)| h.is_finished())
                .map(|(k, _)| k.clone())
                .collect();
            for key in finished {
                if let Some(handle) = workers.0.remove(&key)
                    && handle.await.is_err()
                {
                    let _ = self
                        .store
                        .publishing_heartbeat(&key, Some("WORKER_PANIC"))
                        .await;
                }
            }
            match self.store.watch_rows().await {
                Ok(mut rows) => {
                    if !rows.is_empty() {
                        let n = rows.len();
                        rows.rotate_left(rotate % n);
                        rotate = rotate.wrapping_add(1);
                    }
                    for row in rows {
                        let Some(key) = row["account_key"].as_str() else {
                            continue;
                        };
                        if let Err(e) = self
                            .store
                            .configure_publishing(key, &url, enabled, &policy)
                            .await
                        {
                            tracing::warn!(account_key=key,code=%e.code,"publication_control_failed");
                            let _ = self.store.publishing_heartbeat(key, Some(&e.code)).await;
                            continue;
                        }
                        if let Err(e) = self
                            .store
                            .configure_spot_publishing(key, enabled && spot_enabled, &policy)
                            .await
                        {
                            let _ = self.store.publishing_heartbeat(key, Some(&e.code)).await;
                            continue;
                        }
                        let _ = self.store.touch_publishing_heartbeat(key).await;
                        if workers.0.contains_key(key) {
                            continue;
                        }
                        let Ok(permit) = semaphore.clone().try_acquire_owned() else {
                            continue;
                        };
                        let store = self.store.clone();
                        let key = key.to_owned();
                        let webhook = self.webhook.clone();
                        let seconds = self.lease_seconds;
                        let base = self.retry_base_seconds;
                        let max = self.retry_max_seconds;
                        workers.0.insert(key.clone(),tokio::spawn(async move {
       let _permit=permit;
       let result=async {
        if let Some(target)=webhook && let Some(claim)=store.claim_publication(&key,uuid::Uuid::new_v4(),seconds).await?{
         let delivered=target.deliver(&claim.event_id,claim.attempt,&claim.body,base,max).await;
         store.finish_publication(&claim,delivered.status,delivered.http,delivered.error,delivered.delay).await?;
         tracing::info!(account_key=%key,event_id=%claim.event_id,status=delivered.status,"publication_completed");
        }Ok::<(),trade_log::query::QueryError>(())
       }.await;
       let _=store.publishing_heartbeat(&key,result.as_ref().err().map(|e|e.code.as_str())).await;
      }));
                    }
                }
                Err(e) => tracing::warn!(code=%e.code,"publication_registry_unavailable"),
            }
            tokio::select! {_=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_millis(self.poll_interval_ms))=>{}}
        }
        // In-flight sends retain their permit and claim until confirmation; outer shutdown is bounded.
        while let Some(key) = workers.0.keys().next().cloned() {
            if let Some(handle) = workers.0.get_mut(&key) {
                let _ = handle.await;
            }
            workers.0.remove(&key);
        }
    }
}

struct DeliveryWorkers(std::collections::BTreeMap<String, tokio::task::JoinHandle<()>>);
impl Drop for DeliveryWorkers {
    fn drop(&mut self) {
        for handle in self.0.values() {
            handle.abort()
        }
    }
}
