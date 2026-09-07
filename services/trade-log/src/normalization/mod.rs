use account_facts::{AccountFact, AccountFactEnvelope, AccountFactPayload, TradeFact};
use chrono::Utc;
use rust_decimal::Decimal;
use serde_json::{json, Value};
use sha2::{Digest, Sha256};
use std::str::FromStr;
use thiserror::Error;

use crate::raw_log::RawTrade;

#[derive(Debug, Clone)]
pub struct FailedTrade {
    pub raw_hash: String,
    pub reason: String,
}

#[derive(Debug, Error, PartialEq)]
pub enum NormalizationError {
    #[error("missing or invalid decimal field: {0}")]
    InvalidDecimal(&'static str),
    #[error("cannot determine trade side from side={side}, action={action}")]
    UnknownSide { side: String, action: String },
    #[error("transaction_hash is required")]
    MissingTransactionHash,
}

pub fn canonical_json(value: &Value) -> Value {
    match value {
        Value::Object(map) => Value::Object(
            map.iter()
                .map(|(k, v)| (k.clone(), canonical_json(v)))
                .collect(),
        ),
        Value::Array(values) => Value::Array(values.iter().map(canonical_json).collect()),
        other => other.clone(),
    }
}

pub fn raw_hash(raw: &RawTrade) -> String {
    let value = canonical_json(&serde_json::to_value(raw).expect("RawTrade is serializable"));
    hex::encode(Sha256::digest(
        serde_json::to_vec(&value).expect("JSON serialization cannot fail"),
    ))
}

fn decimal(value: &Value, field: &'static str) -> Result<String, NormalizationError> {
    let input = match value {
        Value::String(s) => s.clone(),
        Value::Number(n) => n.to_string(),
        _ => return Err(NormalizationError::InvalidDecimal(field)),
    };
    Decimal::from_str(&input)
        .map(|v| v.normalize().to_string())
        .map_err(|_| NormalizationError::InvalidDecimal(field))
}

fn trade_side(raw: &RawTrade) -> Result<&'static str, NormalizationError> {
    let side = raw.side.trim().to_ascii_lowercase();
    let action = raw.action.trim().to_ascii_lowercase();
    match (side.as_str(), action.as_str()) {
        ("buy" | "b", _) => Ok("BUY"),
        ("sell" | "s", _) => Ok("SELL"),
        ("long", action) if action.contains("open") || action.contains("add") => Ok("BUY"),
        ("short", action) if action.contains("open") || action.contains("add") => Ok("SELL"),
        ("long", action) if action.contains("close") || action.contains("reduce") => Ok("SELL"),
        ("short", action) if action.contains("close") || action.contains("reduce") => Ok("BUY"),
        _ => Err(NormalizationError::UnknownSide {
            side: raw.side.clone(),
            action: raw.action.clone(),
        }),
    }
}

fn position_effect(raw: &RawTrade) -> &'static str {
    let action = raw.action.to_ascii_lowercase();
    if action.contains("reverse") {
        "REVERSE"
    } else if action.contains("close") {
        "CLOSE"
    } else if action.contains("open") {
        "OPEN"
    } else if action.contains("add") || action.contains("increase") {
        "INCREASE"
    } else if action.contains("reduce") || action.contains("decrease") {
        "DECREASE"
    } else {
        "UNKNOWN"
    }
}

pub fn normalize_trade(
    raw: &RawTrade,
    observed_at: chrono::DateTime<Utc>,
    trace_id: &str,
) -> Result<AccountFactEnvelope, NormalizationError> {
    if raw.transaction_hash.trim().is_empty() {
        return Err(NormalizationError::MissingTransactionHash);
    }
    let hash = raw_hash(raw);
    let account = raw.user.to_ascii_lowercase();
    let order_id = raw.oid.as_ref().map(|v| match v {
        Value::String(s) => s.clone(),
        other => other.to_string(),
    });
    let operation_id = format!(
        "{}:{}",
        raw.transaction_hash,
        order_id.as_deref().unwrap_or("unknown")
    );
    let ordering_key = format!(
        "{:020}:{}",
        raw.block_number.unwrap_or_default(),
        operation_id
    );
    let now = Utc::now();
    let fact = AccountFact {
        fact_id: hash.clone(),
        fact_type: "TRADE".into(),
        revision: 1,
        change_type: "UPSERT".into(),
        confirmation_status: "OBSERVED".into(),
        chain_id: "hyperliquid".into(),
        protocol: "hyperliquid".into(),
        account: account.clone(),
        account_key: format!("hyperliquid:hyperliquid:{account}"),
        ordering_key,
        sub_index: 0,
        source: "NANSEN".into(),
        source_ref: operation_id.clone(),
        raw_log_id: None,
        payload: AccountFactPayload::Trade(TradeFact {
            market: format!("hyperliquid:{}-USDC", raw.token_symbol),
            instrument_type: "PERPETUAL".into(),
            base_asset: raw.token_symbol.clone(),
            quote_asset: "USDC".into(),
            action: "TRADE".into(),
            trigger_type: "USER".into(),
            side: trade_side(raw)?.into(),
            position_effect: position_effect(raw).into(),
            order_id,
            operation_id,
            price: decimal(&raw.price, "price")?,
            quantity: decimal(&raw.size, "size")?,
            notional: decimal(&raw.value_usd, "value_usd")?,
            fee: decimal(&raw.fee_usd, "fee_usd")?,
            fee_asset: raw.fee_token_symbol.clone(),
            reported_realized_pnl: raw
                .closed_pnl
                .as_ref()
                .map(|v| decimal(v, "closed_pnl"))
                .transpose()?,
            reported_pnl_asset: Some("USDC".into()),
            reported_pnl_includes_fee: None,
            reported_pnl_includes_funding: None,
            transaction_hash: raw.transaction_hash.clone(),
            extension: json!({}),
        }),
    };
    Ok(AccountFactEnvelope {
        event_id: format!("evt_{hash}"),
        event_type: "account.fact.v1".into(),
        schema_version: 1,
        occurred_at: raw.timestamp,
        observed_at,
        published_at: now,
        producer: "trade-log-service".into(),
        trace_id: trace_id.into(),
        fact,
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::TimeZone;
    use serde_json::json;

    fn fixture() -> RawTrade {
        serde_json::from_value(json!({"user":"0xABCDEFabcdefABCDEFabcdefABCDEFabcdefABCD","timestamp":"2026-09-01T01:02:03Z","block_number":42,"transaction_hash":"0xtx","oid":7,"token_symbol":"ETH","side":"Long","action":"Open Long","price":"3245.1200","size":"1.25","value_usd":"4056.40","fee_usd":"2.03","fee_token_symbol":"USDC","closed_pnl":"0","start_position":"0","crossed":false})).unwrap()
    }

    #[test]
    fn maps_nansen_trade_deterministically() {
        let observed = Utc.with_ymd_and_hms(2026, 9, 1, 1, 3, 0).unwrap();
        let first = normalize_trade(&fixture(), observed, "trace").unwrap();
        let second = normalize_trade(&fixture(), observed, "trace").unwrap();
        assert_eq!(first.fact.fact_id, second.fact.fact_id);
        assert_eq!(
            first.fact.account,
            "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
        );
        let AccountFactPayload::Trade(trade) = first.fact.payload;
        assert_eq!(trade.side, "BUY");
        assert_eq!(trade.position_effect, "OPEN");
        assert_eq!(trade.price, "3245.12");
    }

    #[test]
    fn rejects_unknown_side() {
        let mut raw = fixture();
        raw.side = "mystery".into();
        assert!(matches!(
            normalize_trade(&raw, Utc::now(), "trace"),
            Err(NormalizationError::UnknownSide { .. })
        ));
    }

    #[test]
    fn closing_long_is_a_sell() {
        let mut raw = fixture();
        raw.action = "Close".into();
        let fact = normalize_trade(&raw, Utc::now(), "trace").unwrap();
        let AccountFactPayload::Trade(trade) = fact.fact.payload;
        assert_eq!(trade.side, "SELL");
        assert_eq!(trade.position_effect, "CLOSE");
    }

    #[test]
    fn envelope_json_contains_direct_trade_payload() {
        let fact = normalize_trade(&fixture(), Utc::now(), "trace").unwrap();
        let value = serde_json::to_value(fact).unwrap();
        assert_eq!(value["fact"]["payload"]["market"], "hyperliquid:ETH-USDC");
        assert!(value["fact"]["payload"].get("type").is_none());
    }
}
