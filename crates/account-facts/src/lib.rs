use serde::{Deserialize, Serialize};
use serde_json::Value;
use sha2::{Digest, Sha256};
use shared_types::Network;

#[derive(Clone, Debug, Serialize, Deserialize, PartialEq)]
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
    pub raw_log_id: String,
    pub occurred_at: String,
    pub payload: TradeFact,
}
#[derive(Clone, Debug, Serialize, Deserialize, PartialEq)]
pub struct TradeFact {
    pub market: String,
    pub instrument_type: String,
    pub base_asset: String,
    pub quote_asset: String,
    pub action: String,
    pub trigger_type: String,
    pub copy_eligible: bool,
    pub side: String,
    pub position_effect: String,
    pub order_id: Option<String>,
    pub operation_id: Option<String>,
    pub price: String,
    pub quantity: String,
    pub notional: String,
    pub fee: String,
    pub fee_asset: Option<String>,
    pub reported_realized_pnl: Option<String>,
    pub reported_pnl_asset: Option<String>,
    pub reported_pnl_includes_fee: Option<bool>,
    pub reported_pnl_includes_funding: Option<bool>,
    pub transaction_hash: Option<String>,
    pub extension: Value,
}
// v1 identity uses JSON array framing, fixed field order and UTF-8, then SHA-256.
// No query IDs, response order or transport-specific fields participate.
pub fn fact_id(network: &Network, account: &str, market: &str, tid: u64) -> String {
    let canonical = serde_json::to_vec(&[
        "hl-fill-v1",
        network.name(),
        "hyperliquid",
        account,
        market,
        &tid.to_string(),
    ])
    .expect("string serialization");
    format!("hl_fill_v1_{:x}", Sha256::digest(canonical))
}

/// Stable candidate event. Retry headers may change; serialized event bytes do not.
#[derive(Clone, Debug, serde::Serialize, serde::Deserialize)]
pub struct AccountFactEnvelope {
    pub schema_version: u32,
    pub event_type: String,
    pub event_id: String,
    pub partition_key: String,
    pub occurred_at: String,
    pub received_at: String,
    pub stored_at: String,
    pub published_at: Option<String>,
    pub expires_at: String,
    pub observation: serde_json::Value,
    pub fact: AccountFact,
}
