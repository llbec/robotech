use crate::query::QueryError;
use serde::{Deserialize, Serialize};
#[derive(Clone, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct MonitoringConfig {
    pub max_active_accounts: i32,
    pub max_registered_accounts: i32,
    pub reconcile_interval_seconds: u64,
    pub worker_stale_after_seconds: i32,
    pub max_concurrent_http_requests: usize,
    pub max_total_pending_bytes: usize,
}
impl Default for MonitoringConfig {
    fn default() -> Self {
        Self {
            max_active_accounts: 5,
            max_registered_accounts: 100,
            reconcile_interval_seconds: 2,
            worker_stale_after_seconds: 30,
            max_concurrent_http_requests: 2,
            max_total_pending_bytes: 67108864,
        }
    }
}
impl MonitoringConfig {
    pub fn validate(&self) -> Result<(), String> {
        if !(1..=10).contains(&self.max_active_accounts)
            || !(self.max_active_accounts..=10000).contains(&self.max_registered_accounts)
            || !(1..=30).contains(&self.reconcile_interval_seconds)
            || !(10..=180).contains(&self.worker_stale_after_seconds)
            || !(1..=10).contains(&self.max_concurrent_http_requests)
            || !(16777216..=1073741824).contains(&self.max_total_pending_bytes)
        {
            return Err("Invalid monitoring limits".into());
        }
        Ok(())
    }
}
#[derive(Clone, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct BudgetConfig {
    pub http_weight_per_minute: i32,
    pub websocket_connections_per_minute: i32,
}
impl Default for BudgetConfig {
    fn default() -> Self {
        Self {
            http_weight_per_minute: 600,
            websocket_connections_per_minute: 15,
        }
    }
}
impl BudgetConfig {
    pub fn validate(&self) -> Result<(), String> {
        if !(120..=1200).contains(&self.http_weight_per_minute)
            || !(1..=30).contains(&self.websocket_connections_per_minute)
        {
            return Err("Invalid shared source budget".into());
        }
        Ok(())
    }
}
#[derive(Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct AddAccount {
    pub account: String,
    pub start_time: String,
}
impl AddAccount {
    pub fn validate(&mut self) -> Result<(), QueryError> {
        self.account = shared_types::account(&self.account).map_err(invalid)?;
        let ms = crate::collection::start_ms(&self.start_time).map_err(|e| invalid(e.message))?;
        if ms > chrono::Utc::now().timestamp_millis() {
            return Err(invalid("start_time must not be in the future"));
        }
        Ok(())
    }
}
pub fn invalid(message: impl Into<String>) -> QueryError {
    crate::collection::error("INVALID_ARGUMENT", &message.into(), false)
}
#[derive(Clone, serde::Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ListParams {
    pub state: Option<String>,
    pub limit: Option<usize>,
    pub cursor: Option<String>,
}
impl Default for ListParams {
    fn default() -> Self {
        Self {
            state: None,
            limit: Some(20),
            cursor: None,
        }
    }
}
