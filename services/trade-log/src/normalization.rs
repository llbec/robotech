use crate::{
    domain::{HyperliquidFill, TradeFact, TradePayload},
    market_metadata::MetadataStore,
};
use anyhow::{bail, Context, Result};
use chrono::{DateTime, Utc};
use rust_decimal::Decimal;
use sha2::{Digest, Sha256};
use std::str::FromStr;
use uuid::Uuid;

pub async fn normalize(
    account: &str,
    fill: &HyperliquidFill,
    raw_log_id: Uuid,
    observed_at: DateTime<Utc>,
    metadata: &MetadataStore,
) -> Result<TradeFact> {
    let (instrument_type, market) = metadata.resolve(&fill.coin).await?;
    let occurred_at = DateTime::<Utc>::from_timestamp_millis(fill.time)
        .context("fill time is outside supported range")?;
    let price = Decimal::from_str(&fill.px).context("invalid fill price")?;
    let quantity = Decimal::from_str(&fill.sz).context("invalid fill size")?;
    if price.is_sign_negative() || quantity <= Decimal::ZERO {
        bail!("price and size must be positive")
    }
    let side = match fill.side.as_str() {
        "B" => "BUY",
        "A" => "SELL",
        value => bail!("unknown fill side: {value}"),
    };
    let forced = is_forced(&fill.dir);
    let trigger_type = if fill.dir.to_ascii_lowercase().contains("liquidat") {
        "LIQUIDATION"
    } else if forced {
        "PROTOCOL"
    } else {
        "USER"
    };
    let position_effect = if instrument_type == "SPOT" {
        "NONE".into()
    } else {
        derivative_effect(&fill.start_position, &fill.sz, &fill.side)?
    };
    let source_ref = format!("{}:{}", account, fill.tid);
    let fact_id = stable_id("hl-fill", &source_ref);
    Ok(TradeFact {
        event_id: stable_id("account-fact-v1", &fact_id),
        fact_id,
        revision: 1,
        fact_type: "TRADE".into(),
        schema_version: 1,
        chain_id: "hyperliquid:mainnet".into(),
        protocol: "hyperliquid".into(),
        account: account.into(),
        account_key: format!("hyperliquid:mainnet:hyperliquid:{account}"),
        occurred_at,
        observed_at,
        ordering_key: format!("{:013}:{:020}", fill.time, fill.tid),
        source: "HYPERLIQUID_OFFICIAL".into(),
        source_ref,
        raw_log_id,
        payload: TradePayload {
            instrument_type,
            market: market.market,
            base_asset: market.base_asset,
            quote_asset: market.quote_asset,
            side: side.into(),
            position_effect,
            trigger_type: trigger_type.into(),
            copy_eligible: trigger_type == "USER",
            price: price.normalize().to_string(),
            quantity: quantity.normalize().to_string(),
            notional: (price * quantity).normalize().to_string(),
            fee: fill.fee.clone(),
            fee_asset: fill
                .fee_token
                .as_ref()
                .map(|value| value.trim().to_string()),
            reported_realized_pnl: fill.closed_pnl.clone(),
            order_id: fill.oid.to_string(),
            transaction_hash: fill.hash.clone(),
        },
    })
}

pub fn source_event_id(account: &str, fill: &HyperliquidFill) -> String {
    format!("hl-fill:{account}:{}", fill.tid)
}

fn stable_id(namespace: &str, value: &str) -> String {
    let mut hash = Sha256::new();
    hash.update(namespace.as_bytes());
    hash.update(b":");
    hash.update(value.as_bytes());
    hex::encode(hash.finalize())
}

fn is_forced(dir: &str) -> bool {
    let value = dir.to_ascii_lowercase();
    value.contains("liquidat") || value.contains("deleverag") || value.contains("settlement")
}

fn derivative_effect(start: &str, size: &str, side: &str) -> Result<String> {
    let start = Decimal::from_str(start).context("invalid start position")?;
    let size = Decimal::from_str(size).context("invalid fill size")?;
    let delta = match side {
        "B" => size,
        "A" => -size,
        _ => bail!("unknown side"),
    };
    let end = start + delta;
    let effect = if start.is_zero() && !end.is_zero() {
        "OPEN"
    } else if end.is_zero() && !start.is_zero() {
        "CLOSE"
    } else if start.is_sign_positive() != end.is_sign_positive() {
        "REVERSE"
    } else if end.abs() > start.abs() {
        "INCREASE"
    } else if end.abs() < start.abs() {
        "DECREASE"
    } else {
        "UNKNOWN"
    };
    Ok(effect.into())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn derives_position_effects_from_position_delta() {
        assert_eq!(derivative_effect("0", "2", "B").unwrap(), "OPEN");
        assert_eq!(derivative_effect("2", "1", "B").unwrap(), "INCREASE");
        assert_eq!(derivative_effect("2", "1", "A").unwrap(), "DECREASE");
        assert_eq!(derivative_effect("2", "2", "A").unwrap(), "CLOSE");
        assert_eq!(derivative_effect("2", "3", "A").unwrap(), "REVERSE");
    }

    #[test]
    fn source_identity_is_independent_of_transport() {
        let fill: HyperliquidFill = serde_json::from_value(serde_json::json!({
            "coin":"BTC","px":"100","sz":"1","side":"B","time":1,
            "startPosition":"0","dir":"Open Long","closedPnl":"0","hash":"0x1",
            "oid":2,"crossed":true,"fee":"0.1","feeToken":"USDC","tid":3
        }))
        .unwrap();
        assert_eq!(source_event_id("0xabc", &fill), "hl-fill:0xabc:3");
    }
}
