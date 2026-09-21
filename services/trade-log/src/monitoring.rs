use crate::{
    acquisition::{next_text, HyperliquidClient, UserSocket},
    config::Config,
    domain::{HyperliquidFill, WsEnvelope, WsUserFills},
    market_metadata::MetadataStore,
    Repository,
};
use anyhow::{Context, Result};
use chrono::{DateTime, TimeDelta, Utc};
use std::{collections::HashMap, sync::Arc};
use tokio::{sync::Mutex, task::JoinHandle};
use tracing::{error, info, warn};

#[derive(Clone)]
pub struct MonitorManager {
    inner: Arc<Inner>,
}

struct Inner {
    repo: Repository,
    client: HyperliquidClient,
    metadata: Arc<MetadataStore>,
    config: Config,
    tasks: Mutex<HashMap<String, JoinHandle<()>>>,
}

impl MonitorManager {
    pub fn new(
        repo: Repository,
        client: HyperliquidClient,
        metadata: Arc<MetadataStore>,
        config: Config,
    ) -> Self {
        Self {
            inner: Arc::new(Inner {
                repo,
                client,
                metadata,
                config,
                tasks: Mutex::new(HashMap::new()),
            }),
        }
    }

    pub async fn refresh_metadata(&self) -> Result<()> {
        let (perp, spot) = self.inner.client.metadata().await?;
        self.inner
            .metadata
            .replace_and_persist(self.inner.repo.pool(), perp, spot)
            .await
    }

    pub async fn load_persisted_metadata(&self) -> Result<bool> {
        self.inner
            .metadata
            .load_latest(self.inner.repo.pool())
            .await
    }

    pub async fn restore(&self) -> Result<()> {
        for monitor in self.inner.repo.list_monitors().await? {
            self.start(monitor.address).await;
        }
        Ok(())
    }

    pub async fn start(&self, address: String) {
        let mut tasks = self.inner.tasks.lock().await;
        if tasks.contains_key(&address) {
            return;
        }
        let inner = self.inner.clone();
        let task_address = address.clone();
        let handle = tokio::spawn(async move {
            inner.clone().run_address(task_address.clone()).await;
            inner.tasks.lock().await.remove(&task_address);
        });
        tasks.insert(address, handle);
    }

    pub async fn running_count(&self) -> usize {
        self.inner.tasks.lock().await.len()
    }
}

impl Inner {
    async fn run_address(self: Arc<Self>, address: String) {
        loop {
            let initial = match self.repo.get_monitor(&address).await {
                Ok(Some(monitor)) => monitor.coverage_start.is_none(),
                Ok(None) => return,
                Err(error) => {
                    error!(address, error = %error, "failed to read monitor state");
                    tokio::time::sleep(std::time::Duration::from_secs(2)).await;
                    continue;
                }
            };
            if let Err(error) = self.run_session(&address, initial).await {
                error!(address, error = %error, "monitor session failed");
                let message = format!("{error:#}");
                if let Err(update_error) = self
                    .repo
                    .update_monitor_status(&address, "DEGRADED", Some(&message))
                    .await
                {
                    error!(address, error = %update_error, "failed to persist degraded status");
                }
            }
            tokio::time::sleep(std::time::Duration::from_secs(2)).await;
        }
    }

    async fn run_session(&self, address: &str, initial: bool) -> Result<()> {
        if initial {
            self.repo
                .update_monitor_status(address, "BACKFILLING", None)
                .await?;
        }
        let connected_at = Utc::now();
        let mut socket = self.client.connect_user_fills(address).await?;
        self.consume_until_snapshot(address, &mut socket, !initial)
            .await?;

        let monitor = self
            .repo
            .get_monitor(address)
            .await?
            .context("monitor disappeared")?;
        let start = if initial {
            monitor.requested_start
        } else {
            monitor.last_event_time.unwrap_or(monitor.requested_start)
                - TimeDelta::from_std(self.config.overlap_window)?
        };
        let range = self
            .reconcile(address, start, connected_at, !initial)
            .await?;
        if initial {
            let coverage_start = if range.history_complete {
                start
            } else {
                range
                    .first_event
                    .map(|value| value.max(start))
                    .unwrap_or(start)
            };
            self.repo
                .complete_backfill(address, coverage_start, range.history_complete)
                .await?;
        } else {
            self.repo
                .update_monitor_status(address, "LIVE", None)
                .await?;
        }
        info!(address, "monitor is live");

        let mut interval = tokio::time::interval(self.config.reconcile_interval);
        interval.tick().await;
        loop {
            tokio::select! {
                message = next_text(&mut socket) => {
                    let Some(message) = message? else { anyhow::bail!("websocket closed") };
                    self.consume_ws_message(address, &message, true).await?;
                }
                _ = interval.tick() => {
                    let end = Utc::now();
                    let last = self.repo.last_event_time(address).await?.unwrap_or(start);
                    let from = last - TimeDelta::from_std(self.config.overlap_window)?;
                    self.reconcile(address, from, end, true).await?;
                }
            }
        }
    }

    async fn consume_until_snapshot(
        &self,
        address: &str,
        socket: &mut UserSocket,
        publish: bool,
    ) -> Result<()> {
        loop {
            let message = next_text(socket)
                .await?
                .context("websocket closed before snapshot")?;
            if self.consume_ws_message(address, &message, publish).await? {
                return Ok(());
            }
        }
    }

    async fn consume_ws_message(
        &self,
        address: &str,
        message: &str,
        publish: bool,
    ) -> Result<bool> {
        let envelope: WsEnvelope = match serde_json::from_str(message) {
            Ok(value) => value,
            Err(error) => {
                warn!(address, error = %error, "ignored non-envelope websocket message");
                return Ok(false);
            }
        };
        if envelope.channel != "userFills" {
            return Ok(false);
        }
        let data: WsUserFills =
            serde_json::from_value(envelope.data).context("invalid userFills websocket payload")?;
        for fill in &data.fills {
            if let Err(error) = self
                .repo
                .ingest_fill(
                    address,
                    fill,
                    "WEBSOCKET",
                    data.is_snapshot,
                    publish,
                    &self.metadata,
                )
                .await
            {
                warn!(address, tid = fill.tid, error = %error, "failed to ingest websocket fill");
            }
        }
        Ok(data.is_snapshot)
    }

    async fn reconcile(
        &self,
        address: &str,
        start: DateTime<Utc>,
        end: DateTime<Utc>,
        publish: bool,
    ) -> Result<ReconcileResult> {
        let kind = if self
            .repo
            .get_monitor(address)
            .await?
            .is_some_and(|m| m.status == "BACKFILLING")
        {
            "BACKFILL"
        } else {
            "RECONCILE"
        };
        let run_id = self
            .repo
            .start_collection_run(address, kind, start, end)
            .await?;
        let range = match self
            .client
            .fills_by_time(address, start.timestamp_millis(), end.timestamp_millis())
            .await
        {
            Ok(range) => range,
            Err(error) => {
                self.repo
                    .finish_collection_run(run_id, "FAILED", 0, 0, Some(&error.to_string()))
                    .await?;
                return Err(error);
            }
        };
        let first_event = range
            .fills
            .first()
            .and_then(|fill| DateTime::<Utc>::from_timestamp_millis(fill.time));
        let mut inserted = 0_i64;
        for fill in &range.fills {
            match self
                .repo
                .ingest_fill(address, fill, "HTTP", false, publish, &self.metadata)
                .await
            {
                Ok(result) => inserted += i64::from(result.inserted),
                Err(error) => {
                    warn!(address, tid = fill.tid, error = %error, "failed to ingest HTTP fill");
                    self.repo
                        .finish_collection_run(
                            run_id,
                            "FAILED",
                            range.fills.len() as i64,
                            inserted,
                            Some(&error.to_string()),
                        )
                        .await?;
                    return Err(error);
                }
            }
        }
        self.repo
            .finish_collection_run(
                run_id,
                "COMPLETED",
                range.fills.len() as i64,
                inserted,
                None,
            )
            .await?;
        Ok(ReconcileResult {
            first_event,
            history_complete: range.history_complete,
        })
    }
}

struct ReconcileResult {
    first_event: Option<DateTime<Utc>>,
    history_complete: bool,
}

#[allow(dead_code)]
fn _assert_fill_send_sync(_: HyperliquidFill) {}
