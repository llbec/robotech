use crate::postgres::Postgres;
use async_trait::async_trait;
use protocol_api::QueryKind;
use std::{sync::Arc, time::Duration};
use tokio::sync::Semaphore;
use trade_log::{
    acquisition::{SourceReader, SourceResponse},
    query::QueryError,
    realtime::{RealtimeConnection, RealtimeSource},
};

struct Reservation {
    store: Postgres,
    id: uuid::Uuid,
}
impl Drop for Reservation {
    fn drop(&mut self) {
        let store = self.store.clone();
        let id = self.id;
        tokio::spawn(async move {
            let _ = store.release_source(id).await;
        });
    }
}
pub struct BudgetSource {
    pub store: Postgres,
    pub inner: Arc<dyn SourceReader>,
    pub concurrent: Arc<Semaphore>,
}
impl BudgetSource {
    async fn reserve(&self, account: &str, weight: i32) -> Result<uuid::Uuid, QueryError> {
        let mut reported = false;
        loop {
            if let Some(id) = self.store.reserve_source("HTTP", weight).await? {
                return Ok(id);
            }
            if !reported {
                tracing::info!(
                    message = "source_budget_wait",
                    transport = "HTTP",
                    account,
                    weight
                );
                reported = true;
            }
            tokio::time::sleep(Duration::from_millis(250)).await
        }
    }
    async fn settle(
        &self,
        id: uuid::Uuid,
        r: Result<SourceResponse, QueryError>,
        fills: bool,
    ) -> Result<SourceResponse, QueryError> {
        if let Ok(response) = &r
            && (200..300).contains(&response.status)
        {
            let weight = if fills {
                serde_json::from_slice::<Vec<serde_json::Value>>(&response.body)
                    .ok()
                    .map(|v| {
                        trade_log::acquisition::budget::response_weight(
                            QueryKind::UserFills,
                            v.len(),
                        )
                    })
            } else {
                Some(20)
            };
            if let Some(w) = weight {
                self.store.settle_source(id, w).await?
            }
        }
        r
    }
}
#[async_trait]
impl SourceReader for BudgetSource {
    async fn wait_retry(&self, seconds: u64) {
        self.inner.wait_retry(seconds).await
    }
    async fn fetch(&self, kind: QueryKind, account: &str) -> Result<SourceResponse, QueryError> {
        let fills = matches!(kind, QueryKind::UserFills);
        let _permit = self
            .concurrent
            .acquire()
            .await
            .map_err(|_| QueryError::unavailable("Source scheduler stopped"))?;
        let id = self
            .reserve(
                account,
                trade_log::acquisition::budget::reservation_weight(kind),
            )
            .await?;
        let _reservation = Reservation {
            store: self.store.clone(),
            id,
        };
        self.settle(id, self.inner.fetch(kind, account).await, fills)
            .await
    }
    async fn fetch_range(
        &self,
        account: &str,
        start: i64,
        end: i64,
    ) -> Result<SourceResponse, QueryError> {
        let _permit = self
            .concurrent
            .acquire()
            .await
            .map_err(|_| QueryError::unavailable("Source scheduler stopped"))?;
        let id = self
            .reserve(
                account,
                trade_log::acquisition::budget::reservation_weight(QueryKind::UserFillsByTime),
            )
            .await?;
        let _reservation = Reservation {
            store: self.store.clone(),
            id,
        };
        self.settle(id, self.inner.fetch_range(account, start, end).await, true)
            .await
    }
}
pub struct BudgetRealtime {
    pub store: Postgres,
    pub inner: Arc<dyn RealtimeSource>,
}
#[async_trait]
impl RealtimeSource for BudgetRealtime {
    async fn connect(
        &self,
        account: &str,
        max_bytes: usize,
    ) -> Result<Box<dyn RealtimeConnection>, QueryError> {
        let mut reported = false;
        loop {
            if self.store.reserve_source("WS", 1).await?.is_some() {
                break;
            }
            if !reported {
                tracing::info!(message = "source_budget_wait", transport = "WS", account);
                reported = true;
            }
            tokio::time::sleep(Duration::from_millis(250)).await
        }
        self.inner.connect(account, max_bytes).await
    }
}
