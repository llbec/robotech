use crate::query::QueryError;
use async_trait::async_trait;
use protocol_api::QueryKind;
#[derive(Clone)]
pub struct SourceResponse {
    pub status: u16,
    pub body: Vec<u8>,
    pub received_at: String,
    pub retry_after_seconds: Option<u64>,
}
#[async_trait]
pub trait SourceReader: Send + Sync {
    async fn wait_retry(&self, seconds: u64);
    async fn fetch(&self, kind: QueryKind, account: &str) -> Result<SourceResponse, QueryError>;
}
