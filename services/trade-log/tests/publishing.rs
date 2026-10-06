use account_facts::AccountFact;
use chrono::{Duration, Utc};
use serde_json::{Value, json};
use trade_log::publishing::*;
fn fact() -> AccountFact {
    serde_json::from_value(json!({"fact_id":"fact","fact_type":"TRADE","revision":1,"change_type":"UPSERT","confirmation_status":"CONFIRMED","chain_id":"hyperliquid:mainnet","protocol":"hyperliquid","account":"account","account_key":"key","ordering_key":"1","sub_index":0,"source":"hyperliquid","source_ref":"1","raw_log_id":"raw","occurred_at":Utc::now().to_rfc3339(),"payload":{"market":"ETH","instrument_type":"PERPETUAL","base_asset":"ETH","quote_asset":"USDC","action":"TRADE","trigger_type":"USER","copy_eligible":true,"side":"BUY","position_effect":"OPEN","order_id":null,"operation_id":"1","price":"1","quantity":"1","notional":"1","fee":"0","fee_asset":null,"reported_realized_pnl":null,"reported_pnl_asset":null,"reported_pnl_includes_fee":null,"reported_pnl_includes_funding":null,"transaction_hash":null,"extension":{}}})).unwrap()
}
fn context() -> CandidateContext {
    let now = Utc::now();
    CandidateContext {
        enabled: true,
        account_key: "key".into(),
        activated_at: now - Duration::seconds(60),
        now,
        received_at: now,
        subscribed_at: Some(now - Duration::seconds(1)),
        snapshot_sequence: Some(1),
        sequence: 2,
        message_mode: "LIVE_UPDATE".into(),
        flag_absent: false,
        metadata_stale: false,
    }
}
fn policy() -> CandidatePolicy {
    CandidatePolicy {
        max_event_age_seconds: 30,
        signal_ttl_seconds: 60,
        clock_skew_tolerance_seconds: 5,
        version: "candidate-v1".into(),
    }
}
#[test]
fn fixed_modes() {
    let cases: Value = serde_json::from_str(include_str!(
        "../../../tests/fixtures/v0.5/candidate-cases.json"
    ))
    .unwrap();
    for case in cases.as_array().unwrap() {
        let mut c = context();
        c.message_mode = case["mode"].as_str().unwrap().into();
        c.flag_absent = case["flag_absent"].as_bool().unwrap();
        c.snapshot_sequence = case["snapshot"].as_i64();
        let result = decide(&fact(), &c, &policy());
        if let Some(reason) = case["reason"].as_str() {
            assert_eq!(result, Ok(reason));
        } else {
            assert_eq!(result, Err(case["suppression"].as_str().unwrap()));
        }
    }
}
#[test]
fn all_gates() {
    let p = policy();
    let mut c = context();
    let f = fact();
    c.enabled = false;
    assert_eq!(decide(&f, &c, &p), Err("DISABLED"));
    c = context();
    c.subscribed_at = None;
    assert_eq!(decide(&f, &c, &p), Err("SESSION_UNCONFIRMED"));
    c = context();
    c.snapshot_sequence = None;
    c.message_mode = "UNKNOWN".into();
    c.flag_absent = true;
    assert_eq!(decide(&f, &c, &p), Err("MODE_UNCONFIRMED"));
    c = context();
    c.metadata_stale = true;
    assert_eq!(decide(&f, &c, &p), Err("METADATA_STALE"));
    c = context();
    c.activated_at = Utc::now() + Duration::seconds(1);
    assert_eq!(decide(&f, &c, &p), Err("BEFORE_ACTIVATION"));
    c = context();
    let mut old = f.clone();
    old.occurred_at = (c.now - Duration::seconds(31)).to_rfc3339();
    assert_eq!(decide(&old, &c, &p), Err("STALE_EVENT"));
    c.activated_at = c.now - Duration::seconds(3600);
    c.now += Duration::seconds(100);
    assert_eq!(decide(&f, &c, &p), Err("STALE_EVENT"));
    c = context();
    let mut future = f.clone();
    future.occurred_at = (c.now + Duration::seconds(6)).to_rfc3339();
    assert_eq!(decide(&future, &c, &p), Err("FUTURE_EVENT"));
    let mut forced = f.clone();
    forced.payload.trigger_type = "LIQUIDATION".into();
    assert_eq!(decide(&forced, &c, &p), Err("NOT_COPY_ELIGIBLE"));
    forced = f.clone();
    forced.payload.copy_eligible = false;
    assert_eq!(decide(&forced, &c, &p), Err("NOT_COPY_ELIGIBLE"));
}
