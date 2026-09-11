use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum CoverageStatus {
    Complete,
    Partial,
    Unavailable,
    Unverified,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CoverageResult {
    pub category: String,
    pub endpoint: Option<String>,
    pub status: CoverageStatus,
    pub record_count: u64,
    pub range_start: Option<DateTime<Utc>>,
    pub range_end: Option<DateTime<Utc>>,
    pub missing_fields: Vec<String>,
    pub evidence: Value,
}

pub fn unavailable(category: &str, evidence: &str) -> CoverageResult {
    CoverageResult {
        category: category.into(),
        endpoint: None,
        status: CoverageStatus::Unavailable,
        record_count: 0,
        range_start: None,
        range_end: None,
        missing_fields: Vec::new(),
        evidence: serde_json::json!({ "reason": evidence }),
    }
}
