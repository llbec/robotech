use chrono::{DateTime, NaiveDateTime, Utc};
use serde::{de::Error, Deserialize, Deserializer, Serialize};
use serde_json::Value;

#[derive(Debug, Clone)]
pub struct NansenPage {
    pub page: u32,
    pub is_last_page: bool,
    pub response: Value,
    pub trades: Vec<RawTrade>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RawTrade {
    pub user: String,
    #[serde(deserialize_with = "deserialize_nansen_time")]
    pub timestamp: DateTime<Utc>,
    pub block_number: Option<u64>,
    pub transaction_hash: String,
    pub oid: Option<Value>,
    pub token_symbol: String,
    pub side: String,
    pub action: String,
    pub price: Value,
    pub size: Value,
    pub value_usd: Value,
    pub fee_usd: Value,
    pub fee_token_symbol: String,
    pub closed_pnl: Option<Value>,
    pub start_position: Option<Value>,
    pub crossed: Option<bool>,
    #[serde(flatten)]
    pub extra: serde_json::Map<String, Value>,
}

fn deserialize_nansen_time<'de, D>(deserializer: D) -> Result<DateTime<Utc>, D::Error>
where
    D: Deserializer<'de>,
{
    let value = String::deserialize(deserializer)?;
    if let Ok(timestamp) = DateTime::parse_from_rfc3339(&value) {
        return Ok(timestamp.with_timezone(&Utc));
    }
    for format in ["%Y-%m-%dT%H:%M:%S%.f", "%Y-%m-%d %H:%M:%S%.f"] {
        if let Ok(timestamp) = NaiveDateTime::parse_from_str(&value, format) {
            return Ok(timestamp.and_utc());
        }
    }
    Err(D::Error::custom("invalid Nansen UTC timestamp"))
}
