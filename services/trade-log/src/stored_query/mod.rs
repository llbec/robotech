use crate::query::{ObservedRange, QueryError};
use account_facts::AccountFact;
use async_trait::async_trait;
use base64::{Engine, engine::general_purpose::URL_SAFE_NO_PAD};
use chrono::{DateTime, Datelike, SecondsFormat, Utc};
use serde::{Deserialize, Serialize};
fn default_limit() -> usize {
    100
}
#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct StoredRequest {
    pub account: String,
    #[serde(default = "default_limit")]
    pub limit: usize,
    pub start_time: Option<String>,
    pub end_time: Option<String>,
    pub cursor: Option<String>,
}
#[derive(Clone, Debug, Serialize, Deserialize, PartialEq)]
pub struct RequestRange {
    pub start_time: Option<String>,
    pub end_time: Option<String>,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct StoredResult {
    pub account: String,
    pub network: String,
    pub query_scope: String,
    pub coverage: String,
    pub request_range: RequestRange,
    pub snapshot_seq: String,
    pub matched_records: usize,
    pub returned_records: usize,
    pub observed_range: Option<ObservedRange>,
    pub trades: Vec<AccountFact>,
    pub has_more: bool,
    pub next_cursor: Option<String>,
    pub warnings: Vec<String>,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Cursor {
    pub version: u32,
    pub account: String,
    pub network: String,
    pub range: RequestRange,
    pub snapshot: i64,
    pub time: String,
    pub tid: String,
    pub fact_id: String,
}
pub fn time(value: &str) -> Result<String, QueryError> {
    let dt = DateTime::parse_from_rfc3339(value)
        .map_err(|_| QueryError::validation("Invalid RFC3339 time"))?;
    let fractional_digits = value.split_once('.').map_or(0, |(_, tail)| {
        tail.bytes().take_while(u8::is_ascii_digit).count()
    });
    if dt.timestamp_millis() < 0
        || dt.timestamp_subsec_nanos() % 1_000_000 != 0
        || fractional_digits > 3
        || dt.year() > 9999
    {
        return Err(QueryError::validation(
            "Time must be nonnegative with millisecond precision",
        ));
    }
    Ok(dt
        .with_timezone(&Utc)
        .to_rfc3339_opts(SecondsFormat::Millis, true))
}
impl StoredRequest {
    pub fn validated(&self) -> Result<Self, QueryError> {
        let account = shared_types::account(&self.account).map_err(QueryError::validation)?;
        if !(1..=2000).contains(&self.limit) {
            return Err(QueryError::validation("limit must be 1–2000"));
        }
        let start_time = self.start_time.as_deref().map(time).transpose()?;
        let end_time = self.end_time.as_deref().map(time).transpose()?;
        if let (Some(start), Some(end)) = (&start_time, &end_time)
            && start >= end
        {
            return Err(QueryError::validation("start_time must precede end_time"));
        }
        if self.cursor.as_ref().is_some_and(|c| c.len() > 8192) {
            return Err(QueryError::validation("Cursor too large"));
        }
        Ok(Self {
            account,
            limit: self.limit,
            start_time,
            end_time,
            cursor: self.cursor.clone(),
        })
    }
    pub fn range(&self) -> RequestRange {
        RequestRange {
            start_time: self.start_time.clone(),
            end_time: self.end_time.clone(),
        }
    }
    pub fn decode_cursor(&self, network: &str) -> Result<Option<Cursor>, QueryError> {
        let Some(encoded) = &self.cursor else {
            return Ok(None);
        };
        let invalid = || QueryError::validation("Invalid cursor or query conditions");
        let bytes = URL_SAFE_NO_PAD.decode(encoded).map_err(|_| invalid())?;
        let c: Cursor = serde_json::from_slice(&bytes).map_err(|_| invalid())?;
        if c.version != 1
            || c.account != self.account
            || c.network != network
            || c.range != self.range()
            || c.snapshot < 0
            || time(&c.time)? != c.time
            || c.tid.parse::<u64>().is_err()
            || c.fact_id.len() != 75
            || !c.fact_id.starts_with("hl_fill_v1_")
            || !c.fact_id[11..].bytes().all(|b| b.is_ascii_hexdigit())
        {
            return Err(invalid());
        }
        Ok(Some(c))
    }
}
impl Cursor {
    pub fn encode(&self) -> Result<String, QueryError> {
        Ok(URL_SAFE_NO_PAD.encode(serde_json::to_vec(self).map_err(|_| QueryError::storage())?))
    }
}
#[async_trait]
pub trait StoredQuery: Send + Sync {
    async fn stored(&self, request: &StoredRequest) -> Result<StoredResult, QueryError>;
}
