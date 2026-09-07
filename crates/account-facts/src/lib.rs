use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct AccountFactEnvelope {
    pub event_id: String,
    pub event_type: String,
    pub schema_version: u32,
    pub occurred_at: DateTime<Utc>,
    pub observed_at: DateTime<Utc>,
    pub published_at: DateTime<Utc>,
    pub producer: String,
    pub trace_id: String,
    pub fact: AccountFact,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct AccountFact {
    pub fact_id: String,
    pub fact_type: String,
    pub revision: u32,
    pub change_type: String,
    pub confirmation_status: String,
    pub chain_id: String,
    pub protocol: String,
    pub account: String,
    pub account_key: String,
    pub ordering_key: String,
    pub sub_index: u32,
    pub source: String,
    pub source_ref: String,
    pub raw_log_id: Option<String>,
    pub payload: AccountFactPayload,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(untagged)]
pub enum AccountFactPayload {
    Trade(TradeFact),
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct TradeFact {
    pub market: String,
    pub instrument_type: String,
    pub base_asset: String,
    pub quote_asset: String,
    pub action: String,
    pub trigger_type: String,
    pub side: String,
    pub position_effect: String,
    pub order_id: Option<String>,
    pub operation_id: String,
    pub price: String,
    pub quantity: String,
    pub notional: String,
    pub fee: String,
    pub fee_asset: String,
    pub reported_realized_pnl: Option<String>,
    pub reported_pnl_asset: Option<String>,
    pub reported_pnl_includes_fee: Option<bool>,
    pub reported_pnl_includes_funding: Option<bool>,
    pub transaction_hash: String,
    pub extension: Value,
}
