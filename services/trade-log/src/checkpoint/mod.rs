use crate::query::QueryError;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Serialize, Deserialize, PartialEq, Eq)]
pub struct ScanRange {
    pub start_ms: i64,
    pub end_ms: i64,
}
impl ScanRange {
    pub fn split(&self) -> Result<(Self, Self), QueryError> {
        if self.end_ms - self.start_ms <= 1 {
            return Err(crate::collection::error(
                "SOURCE_WINDOW_SATURATED",
                "A millisecond exceeds the source record limit",
                false,
            ));
        }
        let middle = self.start_ms + (self.end_ms - self.start_ms) / 2;
        Ok((
            Self {
                start_ms: self.start_ms,
                end_ms: middle,
            },
            Self {
                start_ms: middle,
                end_ms: self.end_ms,
            },
        ))
    }
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct AcceptedPage {
    pub range: ScanRange,
    pub raw_id: String,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct PendingWork {
    pub query_id: String,
    pub range: ScanRange,
    pub remaining: Vec<ScanRange>,
    pub pages: Vec<AcceptedPage>,
    pub meta: Option<String>,
    pub spot_meta: Option<String>,
}
#[derive(Clone, Debug)]
pub struct Lease {
    pub key: String,
    pub owner: String,
    pub epoch: i64,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct CollectionStatus {
    pub account: String,
    pub account_key: String,
    pub network: String,
    pub status: String,
    pub coverage: String,
    pub initial_start_time: String,
    pub scanned_through: Option<String>,
    pub last_trade_at: Option<String>,
    pub last_attempt_at: Option<String>,
    pub last_success_at: Option<String>,
    pub consecutive_failures: i32,
    pub next_run_at: Option<String>,
    pub pending_range: Option<serde_json::Value>,
    pub last_query_id: Option<String>,
    pub last_success_query_id: Option<String>,
    pub last_error: Option<serde_json::Value>,
    pub heartbeat_at: Option<String>,
    pub lease_expires_at: Option<String>,
    pub warnings: Vec<String>,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct CollectionList {
    pub items: Vec<CollectionStatus>,
}
#[async_trait]
pub trait CollectionStatusReader: Send + Sync {
    async fn collection_status(&self, account: &str) -> Result<CollectionList, QueryError>;
}

/// The adapter must save all observations, facts and this checkpoint atomically.
pub struct CollectionCommit<'a> {
    pub lease: &'a Lease,
    pub work: &'a PendingWork,
    pub observations: &'a [account_facts::AccountFact],
    pub interval_seconds: u64,
}
