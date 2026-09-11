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
    Trade(Box<TradeFact>),
    Transfer(TransferFact),
    Fee(MonetaryFact),
    Funding(MonetaryFact),
    Reward(MonetaryFact),
    LiquidationFee(MonetaryFact),
    AccountSnapshot(Box<AccountSnapshotFact>),
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct TransferFact {
    pub asset: String,
    pub amount: String,
    pub direction: String,
    pub from_account: Option<String>,
    pub to_account: Option<String>,
    pub transfer_type: String,
    pub transaction_hash: Option<String>,
    pub extension: Value,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MonetaryFact {
    pub asset: String,
    pub amount: String,
    pub direction: String,
    pub related_market: Option<String>,
    pub transaction_hash: Option<String>,
    pub extension: Value,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct AccountSnapshotFact {
    pub snapshot_at: DateTime<Utc>,
    pub valuation_currency: String,
    pub account_value: Option<String>,
    pub available_balance: Option<String>,
    pub margin_used: Option<String>,
    pub unrealized_pnl: Option<String>,
    pub balances: Vec<Value>,
    pub positions: Vec<Value>,
    pub extension: Value,
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
