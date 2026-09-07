use account_facts::AccountFactEnvelope;
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use serde_json::Value;
use uuid::Uuid;

use crate::{normalization::FailedTrade, raw_log::RawTrade};

pub struct PageWrite<'a> {
    pub page: u32,
    pub response: &'a Value,
    pub trades: &'a [RawTrade],
    pub facts: &'a [AccountFactEnvelope],
    pub failures: &'a [FailedTrade],
}

#[async_trait]
pub trait TradeRepository: Send + Sync {
    async fn start_run(
        &self,
        address: &str,
        from: DateTime<Utc>,
        to: DateTime<Utc>,
    ) -> Result<Uuid, String>;
    async fn persist_page(
        &self,
        run_id: Uuid,
        write: PageWrite<'_>,
    ) -> Result<Vec<AccountFactEnvelope>, String>;
    async fn complete_run(
        &self,
        run_id: Uuid,
        pages: u32,
        raw: u64,
        normalized: u64,
        failed: u64,
    ) -> Result<(), String>;
    async fn fail_run(&self, run_id: Uuid, error: &str) -> Result<(), String>;
}
