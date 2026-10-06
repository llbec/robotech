pub mod policy;
pub use policy::{CandidateContext, CandidatePolicy, decide};

/// Delivery port: acquisition and candidate policy do not depend on a transport client.
#[async_trait::async_trait]
pub trait PublicationTarget: Send + Sync {
    async fn deliver(
        &self,
        event_id: &str,
        attempt: i32,
        body: &[u8],
        retry_base: u64,
        retry_cap: u64,
    ) -> DeliveryOutcome;
}
pub struct DeliveryOutcome {
    pub status: &'static str,
    pub http: Option<u16>,
    pub error: Option<&'static str>,
    pub delay: u64,
}

#[derive(Clone, Debug, serde::Serialize, serde::Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum PublishingState {
    Disabled,
    Starting,
    Running,
    Degraded,
    Blocked,
    Stopped,
}
#[derive(Clone, Debug, serde::Serialize, serde::Deserialize)]
pub struct PublishingStatus {
    pub enabled: bool,
    pub status: PublishingState,
    pub account: String,
    pub account_key: String,
    pub target_id: String,
    pub activation_epoch: i64,
    pub activated_at: chrono::DateTime<chrono::Utc>,
    pub heartbeat_at: Option<chrono::DateTime<chrono::Utc>>,
    pub policy: CandidatePolicy,
    pub outbox: std::collections::BTreeMap<String, i64>,
    pub oldest_pending_at: Option<chrono::DateTime<chrono::Utc>>,
    pub last_published_at: Option<chrono::DateTime<chrono::Utc>>,
    pub last_delivered_at: Option<chrono::DateTime<chrono::Utc>>,
    pub last_error: Option<String>,
    pub expired_pending: i64,
    pub suppression_counts: std::collections::BTreeMap<String, i64>,
}
