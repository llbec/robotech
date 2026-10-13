use account_facts::AccountFact;
use serde_json::{Value, json};
use trade_log::{parsing::ParseContext, query::QueryError};

pub fn market<'a>(meta: &'a Value, coin: &str) -> Result<&'a Value, QueryError> {
    let markets = meta["universe"]
        .as_array()
        .ok_or_else(|| QueryError::incomplete("Missing spot markets"))?;
    let mut matching = markets.iter().filter(|m| {
        m["name"].as_str() == Some(coin)
            || m["index"].as_u64().is_some_and(|i| coin == format!("@{i}"))
    });
    let m = matching
        .next()
        .ok_or_else(|| QueryError::incomplete("Unknown spot market"))?;
    if matching.next().is_some() {
        return Err(QueryError::incomplete("Ambiguous spot market"));
    }
    Ok(m)
}
fn token(meta: &Value, index: u64) -> Result<&Value, QueryError> {
    let tokens = meta["tokens"]
        .as_array()
        .ok_or_else(|| QueryError::incomplete("Missing spot token metadata"))?;
    let matches: Vec<_> = tokens
        .iter()
        .filter(|t| t["index"].as_u64() == Some(index))
        .collect();
    if matches.len() != 1
        || matches[0]["name"].as_str().is_none_or(str::is_empty)
        || matches[0]["tokenId"].as_str().is_none_or(|s| {
            s.len() != 34 || !s.starts_with("0x") || !s[2..].bytes().all(|b| b.is_ascii_hexdigit())
        })
    {
        return Err(QueryError::incomplete("Invalid spot token identity"));
    }
    Ok(matches[0])
}
pub fn parse(
    c: &ParseContext<'_>,
    fill: &Value,
    index: usize,
    coin: &str,
    tid: u64,
    meta: &Value,
) -> Result<AccountFact, QueryError> {
    let m = market(meta, coin)?;
    let pair = m["index"]
        .as_u64()
        .ok_or_else(|| QueryError::incomplete("Invalid spot pair index"))?;
    let ids = m["tokens"]
        .as_array()
        .filter(|v| v.len() == 2)
        .ok_or_else(|| QueryError::incomplete("Invalid spot pair tokens"))?;
    let base = token(
        meta,
        ids[0]
            .as_u64()
            .ok_or_else(|| QueryError::incomplete("Invalid base token"))?,
    )?;
    let quote = token(
        meta,
        ids[1]
            .as_u64()
            .ok_or_else(|| QueryError::incomplete("Invalid quote token"))?,
    )?;
    if base["index"] == quote["index"] || base["tokenId"] == quote["tokenId"] {
        return Err(QueryError::incomplete("Spot pair has identical assets"));
    }
    // Reuse exact-decimal and envelope validation; perpetual position inference is discarded.
    let mut normalized = fill.clone();
    normalized["startPosition"] = json!("0");
    normalized["dir"] = json!("Open Long");
    normalized
        .as_object_mut()
        .ok_or_else(|| QueryError::incomplete("Invalid fill"))?
        .remove("closedPnl");
    let identity = format!("spot:{pair}");
    let mut f = super::parser::parse_fill(c, &normalized, index, &identity, tid)
        .map_err(QueryError::incomplete)?;
    let p = &mut f.payload;
    p.market = format!("hyperliquid:spot:{pair}");
    p.instrument_type = "SPOT".into();
    p.base_asset = base["name"].as_str().expect("validated").into();
    p.quote_asset = quote["name"].as_str().expect("validated").into();
    p.action = p.side.clone();
    p.position_effect = "NOT_APPLICABLE".into();
    p.reported_realized_pnl = None;
    let fee_name = fill["feeToken"].as_str();
    let fees: Vec<_> = meta["tokens"]
        .as_array()
        .expect("validated tokens")
        .iter()
        .filter(|t| {
            t["name"].as_str() == fee_name && fee_name.is_some() && t["tokenId"].is_string()
        })
        .collect();
    p.fee_asset = if fees.len() == 1
        && fees[0]["index"]
            .as_u64()
            .is_some_and(|index| token(meta, index).is_ok())
    {
        fee_name.map(str::to_owned)
    } else {
        None
    };
    let forced = fill.get("liquidation").is_some_and(|v| !v.is_null());
    let dir = fill["dir"].as_str().unwrap_or("");
    let user = (p.side == "BUY" && dir == "Buy") || (p.side == "SELL" && dir == "Sell");
    p.trigger_type = if forced {
        "LIQUIDATION"
    } else if user {
        "USER"
    } else {
        "UNKNOWN"
    }
    .into();
    p.copy_eligible = !forced && user && p.fee_asset.is_some();
    p.extension = json!({"source_fill":fill,"source_indices":[index],"time_ms":fill["time"],"source_dir":dir,"pair_index":pair,"pair_name":m["name"],"base_token":base,"quote_token":quote,"fee_token":if p.fee_asset.is_some(){fees.first().copied()}else{None},"parser_version":"hyperliquid-spot-v1","fee_asset_unconfirmed":p.fee_asset.is_none()});
    Ok(f)
}
