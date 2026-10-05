use crate::query::QueryError;
use account_facts::AccountFact;
use rust_decimal::Decimal;
use serde_json::Value;
use sha2::{Digest, Sha256};

pub const HASH_VERSION: &str = "hl-trade-content-v2";

/// Compare business meaning independently of transport envelopes and decimal spelling.
pub fn semantic_hash(fact: &AccountFact) -> Result<String, QueryError> {
    let mut value = serde_json::to_value(fact).map_err(|_| QueryError::storage())?;
    value
        .as_object_mut()
        .ok_or_else(QueryError::storage)?
        .remove("raw_log_id");
    // All business fields are explicit in TradeFact; this extension holds source evidence.
    value["payload"]
        .as_object_mut()
        .ok_or_else(QueryError::storage)?
        .remove("extension");
    for field in [
        "price",
        "quantity",
        "notional",
        "fee",
        "reported_realized_pnl",
    ] {
        if let Some(text) = value["payload"][field].as_str() {
            let decimal = Decimal::from_str_exact(text)
                .map_err(|_| QueryError::incomplete("Invalid fact decimal"))?;
            value["payload"][field] = Value::String(decimal.normalize().to_string());
        }
    }
    Ok(format!(
        "{:x}",
        Sha256::digest(
            serde_json::to_vec(&(HASH_VERSION, value)).map_err(|_| QueryError::storage())?
        )
    ))
}
