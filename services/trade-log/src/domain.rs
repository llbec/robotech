use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_json::{Map, Value};
use uuid::Uuid;

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct HyperliquidFill {
    pub coin: String,
    pub px: String,
    pub sz: String,
    pub side: String,
    pub time: i64,
    pub start_position: String,
    pub dir: String,
    #[serde(default)]
    pub closed_pnl: Option<String>,
    pub hash: String,
    pub oid: u64,
    pub crossed: bool,
    #[serde(default)]
    pub fee: Option<String>,
    #[serde(default)]
    pub fee_token: Option<String>,
    #[serde(default)]
    pub builder_fee: Option<String>,
    pub tid: u64,
    #[serde(flatten)]
    pub extra: Map<String, Value>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct WsUserFills {
    #[serde(default, rename = "isSnapshot")]
    pub is_snapshot: bool,
    pub user: String,
    pub fills: Vec<HyperliquidFill>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct WsEnvelope {
    pub channel: String,
    pub data: Value,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct TradePayload {
    pub instrument_type: String,
    pub market: String,
    pub base_asset: String,
    pub quote_asset: String,
    pub side: String,
    pub position_effect: String,
    pub trigger_type: String,
    pub copy_eligible: bool,
    pub price: String,
    pub quantity: String,
    pub notional: String,
    pub fee: Option<String>,
    pub fee_asset: Option<String>,
    pub reported_realized_pnl: Option<String>,
    pub order_id: String,
    pub transaction_hash: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct TradeFact {
    pub event_id: String,
    pub fact_id: String,
    pub revision: i32,
    pub fact_type: String,
    pub schema_version: i32,
    pub chain_id: String,
    pub protocol: String,
    pub account: String,
    pub account_key: String,
    pub occurred_at: DateTime<Utc>,
    pub observed_at: DateTime<Utc>,
    pub ordering_key: String,
    pub source: String,
    pub source_ref: String,
    pub raw_log_id: Uuid,
    pub payload: TradePayload,
}

#[derive(Clone, Debug, Serialize, sqlx::FromRow)]
pub struct MonitoredAddress {
    pub address: String,
    pub status: String,
    pub requested_start: DateTime<Utc>,
    pub coverage_start: Option<DateTime<Utc>>,
    pub history_complete: bool,
    pub last_event_time: Option<DateTime<Utc>>,
    pub last_received_at: Option<DateTime<Utc>>,
    pub last_error: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

pub fn validate_address(value: &str) -> Option<String> {
    let normalized = value.trim().to_ascii_lowercase();
    if normalized.len() != 42 || !normalized.starts_with("0x") {
        return None;
    }
    normalized[2..]
        .chars()
        .all(|c| c.is_ascii_hexdigit())
        .then_some(normalized)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn validates_and_normalizes_address() {
        assert_eq!(
            validate_address("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
            Some("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa".into())
        );
        assert!(validate_address("0x123").is_none());
    }
}
