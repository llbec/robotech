use crate::{collector_runtime::CollectorRuntime, postgres::Postgres};
use std::{collections::BTreeMap, sync::Arc, time::Duration};
use tokio_util::sync::CancellationToken;
use trade_log::{
    acquisition::SourceReader,
    collection::CollectionConfig,
    parsing::ProtocolParser,
    realtime::{RealtimeSource, WebsocketConfig},
    watch::MonitoringConfig,
};
struct Worker {
    owner: uuid::Uuid,
    generation: i64,
    cancel: CancellationToken,
    handle: tokio::task::JoinHandle<()>,
}
impl Drop for Worker {
    fn drop(&mut self) {
        self.cancel.cancel();
        self.handle.abort();
    }
}
pub struct CollectorSupervisor {
    pub store: Postgres,
    pub defaults: CollectionConfig,
    pub monitoring: MonitoringConfig,
    pub websocket: WebsocketConfig,
    pub source: Arc<dyn SourceReader>,
    pub realtime: Arc<dyn RealtimeSource>,
    pub parser: Arc<dyn ProtocolParser>,
}
impl CollectorSupervisor {
    pub async fn run(self, cancel: CancellationToken) {
        let mut workers: BTreeMap<String, Worker> = BTreeMap::new();
        loop {
            if cancel.is_cancelled() {
                break;
            }
            match self.store.watch_rows().await {
                Ok(rows) => {
                    let active: BTreeMap<String, serde_json::Value> = rows
                        .into_iter()
                        .filter(|v| v["desired_state"] == "ACTIVE")
                        .filter_map(|v| Some((v["account_key"].as_str()?.to_owned(), v)))
                        .collect();
                    let remove: Vec<String> = workers
                        .iter()
                        .filter(|(k, w)| {
                            w.handle.is_finished()
                                || active
                                    .get(*k)
                                    .is_none_or(|v| v["generation"].as_i64() != Some(w.generation))
                        })
                        .map(|(k, _)| k.clone())
                        .collect();
                    for k in remove {
                        if let Some(mut w) = workers.remove(&k) {
                            let finished = w.handle.is_finished();
                            w.cancel.cancel();
                            match tokio::time::timeout(Duration::from_secs(2), &mut w.handle).await
                            {
                                Ok(Err(_)) => {
                                    let _ = self
                                        .store
                                        .watch_worker_error(&k, w.generation, Some("WORKER_PANIC"))
                                        .await;
                                }
                                Ok(Ok(())) if finished => {
                                    let _ = self
                                        .store
                                        .watch_worker_error(
                                            &k,
                                            w.generation,
                                            Some("WORKER_STOPPED"),
                                        )
                                        .await;
                                }
                                _ => {}
                            }
                            w.handle.abort();
                            let _ = self.store.fence_worker(&k, w.owner).await;
                        }
                    }
                    for (key, row) in active {
                        if workers.contains_key(&key) {
                            continue;
                        }
                        let mut config = self.defaults.clone();
                        config.account = row["account"].as_str().unwrap_or_default().into();
                        config.start_time = row["start_time"].as_str().unwrap_or_default().into();
                        let generation = row["generation"].as_i64().unwrap_or(0);
                        let runtime = CollectorRuntime {
                            store: self.store.clone(),
                            config,
                            source: self.source.clone(),
                            parser: self.parser.clone(),
                        };
                        let child = cancel.child_token();
                        let child_run = child.clone();
                        let ws = self.websocket.clone();
                        let source = self.realtime.clone();
                        let owner = uuid::Uuid::new_v4();
                        let handle = tokio::spawn(async move {
                            runtime.run_owned(child_run, ws, source, owner).await
                        });
                        workers.insert(
                            key,
                            Worker {
                                owner,
                                generation,
                                cancel: child,
                                handle,
                            },
                        );
                    }
                }
                Err(e) => {
                    tracing::warn!(code=%e.code,"monitoring_registry_unavailable");
                }
            }
            tokio::select! {_=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_secs(self.monitoring.reconcile_interval_seconds))=>{}}
        }
        for w in workers.values() {
            w.cancel.cancel()
        }
        for (_, mut w) in workers {
            let _ = tokio::time::timeout(Duration::from_secs(3), &mut w.handle).await;
        }
    }
}
