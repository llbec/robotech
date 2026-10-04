use account_facts::{AccountFact, TradeFact, fact_id};
use serde_json::{Value, json};
use shared_types::{Decimal, decimal, timestamp};
use std::collections::{BTreeMap, BTreeSet};
use trade_log::{
    parsing::{ParseContext, Parsed, ProtocolParser},
    query::QueryError,
};

pub struct HyperliquidParser;
fn object(bytes: &[u8]) -> Result<Value, QueryError> {
    serde_json::from_slice(bytes).map_err(|_| QueryError::incomplete("Invalid source JSON"))
}
fn text<'a>(v: &'a Value, key: &str) -> Result<&'a str, &'static str> {
    v.get(key)
        .and_then(Value::as_str)
        .ok_or("missing string field")
}
fn amount(v: &Value, key: &str) -> Result<Decimal, &'static str> {
    decimal(text(v, key)?)
}
fn uint(v: &Value, key: &str) -> Result<u64, &'static str> {
    v.get(key)
        .and_then(Value::as_u64)
        .ok_or("missing integer identity")
}

fn metadata(bytes: &[u8]) -> Result<Vec<Value>, QueryError> {
    object(bytes)
        .map_err(|_| QueryError::unavailable("Source market metadata invalid"))?
        .get("universe")
        .and_then(Value::as_array)
        .cloned()
        .ok_or_else(|| QueryError::unavailable("Source market metadata invalid"))
}

impl ProtocolParser for HyperliquidParser {
    fn parse(&self, c: ParseContext<'_>) -> Result<Parsed, QueryError> {
        let meta = metadata(c.meta)?;
        let spots = metadata(c.spot_meta)?;
        let mut perps = BTreeSet::new();
        for market in &meta {
            perps.insert(
                text(market, "name")
                    .map_err(|_| QueryError::unavailable("Source market metadata invalid"))?
                    .to_owned(),
            );
        }
        let mut spot_names = BTreeSet::new();
        for market in &spots {
            spot_names.insert(
                text(market, "name")
                    .map_err(|_| QueryError::unavailable("Source spot metadata invalid"))?
                    .to_owned(),
            );
            spot_names.insert(format!(
                "@{}",
                uint(market, "index")
                    .map_err(|_| QueryError::unavailable("Source spot metadata invalid"))?
            ));
        }
        if !perps.is_disjoint(&spot_names) {
            return Err(QueryError::unavailable(
                "Conflicting source market metadata",
            ));
        }
        let fills = object(c.fills)?;
        let fills = fills
            .as_array()
            .ok_or_else(|| QueryError::incomplete("Source fills must be an array"))?;
        let mut result = Parsed {
            source_records: fills.len(),
            duplicate_records: 0,
            spot_records: 0,
            unsupported_records: 0,
            invalid_records: 0,
            trades: Vec::new(),
            warnings: Vec::new(),
        };
        let mut seen: BTreeMap<(String, u64), (Value, Option<usize>)> = BTreeMap::new();
        for (index, fill) in fills.iter().enumerate() {
            let (Ok(coin), Ok(tid)) = (text(fill, "coin"), uint(fill, "tid")) else {
                result.invalid_records += 1;
                result
                    .warnings
                    .push(format!("INVALID_RECORD:{index}:identity"));
                continue;
            };
            let key = (coin.to_owned(), tid);
            if let Some((previous, trade_index)) = seen.get(&key) {
                if previous != fill {
                    return Err(QueryError::conflict());
                }
                result.duplicate_records += 1;
                if let Some(i) = trade_index {
                    result.trades[*i].payload.extension["source_indices"]
                        .as_array_mut()
                        .expect("indices array")
                        .push(index.into());
                }
                continue;
            }
            seen.insert(key.clone(), (fill.clone(), None));
            if spot_names.contains(coin) {
                result.spot_records += 1;
                continue;
            }
            if !perps.contains(coin) {
                result.unsupported_records += 1;
                result
                    .warnings
                    .push(format!("UNSUPPORTED_OR_UNKNOWN_MARKET:{index}"));
                continue;
            }
            match parse_fill(&c, fill, index, coin, tid) {
                Ok(fact) => {
                    seen.get_mut(&key).expect("inserted identity").1 = Some(result.trades.len());
                    result.trades.push(fact);
                }
                Err(reason) => {
                    result.invalid_records += 1;
                    result
                        .warnings
                        .push(format!("INVALID_RECORD:{index}:{reason}"));
                }
            }
        }
        Ok(result)
    }
}

fn parse_fill(
    c: &ParseContext<'_>,
    fill: &Value,
    index: usize,
    coin: &str,
    tid: u64,
) -> Result<AccountFact, &'static str> {
    let price = amount(fill, "px")?;
    let quantity = amount(fill, "sz")?;
    let fee = amount(fill, "fee")?;
    if price <= Decimal::ZERO || quantity <= Decimal::ZERO {
        return Err("nonpositive price or quantity");
    }
    let price = price.normalize();
    let quantity = quantity.normalize();
    let product_scale = price
        .scale()
        .checked_add(quantity.scale())
        .ok_or("unsupported product scale")?;
    if product_scale > 28 {
        return Err("unsupported product scale");
    }
    let coefficient = price
        .mantissa()
        .checked_mul(quantity.mantissa())
        .ok_or("notional overflow")?;
    let notional = Decimal::try_from_i128_with_scale(coefficient, product_scale)
        .map_err(|_| "notional precision loss")?;
    let side = match text(fill, "side")? {
        "B" => "BUY",
        "A" => "SELL",
        _ => return Err("unknown side"),
    };
    let occurred_ms = fill
        .get("time")
        .and_then(Value::as_i64)
        .ok_or("missing time")?;
    let occurred_at = timestamp(occurred_ms)?;
    let start = amount(fill, "startPosition")?;
    let delta = if side == "BUY" { quantity } else { -quantity };
    let end = start.checked_add(delta).ok_or("position overflow")?;
    let effect = if start == Decimal::ZERO {
        "OPEN"
    } else if end == Decimal::ZERO {
        "CLOSE"
    } else if start.is_sign_negative() != end.is_sign_negative() {
        "REVERSE"
    } else if end.abs() > start.abs() {
        "INCREASE"
    } else {
        "DECREASE"
    };
    let dir = text(fill, "dir")?;
    let forced = fill.get("liquidation").is_some_and(|v| !v.is_null())
        || dir.to_ascii_lowercase().contains("liquidat");
    let recognized = [
        "Open Long",
        "Open Short",
        "Close Long",
        "Close Short",
        "Long > Short",
        "Short > Long",
    ]
    .contains(&dir);
    let trigger = if forced {
        "LIQUIDATION"
    } else if recognized {
        "USER"
    } else {
        "PROTOCOL"
    };
    let copy_eligible = !forced && recognized;
    let fee_asset = fill
        .get("feeToken")
        .and_then(Value::as_str)
        .map(str::trim)
        .filter(|s| !s.is_empty())
        .map(str::to_owned);
    let reported = fill
        .get("closedPnl")
        .map(|_| amount(fill, "closedPnl"))
        .transpose()?
        .map(|v| v.to_string());
    let chain_id = format!("hyperliquid:{}", c.network.name());
    Ok(AccountFact {
        fact_id: fact_id(c.network, c.account, coin, tid),
        fact_type: "TRADE".into(),
        revision: 1,
        change_type: "UPSERT".into(),
        confirmation_status: "OBSERVED".into(),
        account_key: format!("{chain_id}:hyperliquid:{}", c.account),
        chain_id,
        protocol: "hyperliquid".into(),
        account: c.account.into(),
        ordering_key: format!("{occurred_ms:020}:{tid:020}"),
        sub_index: 0,
        source: "OFFICIAL_API".into(),
        source_ref: tid.to_string(),
        raw_log_id: c.raw_log_id.into(),
        occurred_at,
        payload: TradeFact {
            market: format!("hyperliquid:{coin}-USDC"),
            instrument_type: "PERPETUAL".into(),
            base_asset: coin.into(),
            quote_asset: "USDC".into(),
            action: "TRADE".into(),
            trigger_type: trigger.into(),
            copy_eligible,
            side: side.into(),
            position_effect: effect.into(),
            order_id: fill
                .get("oid")
                .and_then(Value::as_u64)
                .map(|v| v.to_string()),
            operation_id: None,
            price: price.to_string(),
            quantity: quantity.to_string(),
            notional: notional.to_string(),
            fee: fee.to_string(),
            fee_asset,
            reported_realized_pnl: reported,
            reported_pnl_asset: None,
            reported_pnl_includes_fee: None,
            reported_pnl_includes_funding: None,
            transaction_hash: fill.get("hash").and_then(Value::as_str).map(str::to_owned),
            extension: json!({"source_fill":fill,"source_indices":[index],"time_ms":occurred_ms,"source_dir":dir,"trigger_unconfirmed":!forced && !recognized}),
        },
    })
}
