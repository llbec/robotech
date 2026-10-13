use hyperliquid::parser::HyperliquidParser;
use serde_json::{Value, json};
use shared_types::Network;
use trade_log::parsing::{ParseContext, ProtocolParser};
fn parse(
    fills: Value,
    meta: Value,
) -> Result<trade_log::parsing::Parsed, trade_log::query::QueryError> {
    HyperliquidParser.parse(ParseContext {
        network: &Network::Mainnet,
        account: "0x0000000000000000000000000000000000000001",
        raw_log_id: "raw_spot",
        fills: &serde_json::to_vec(&fills).unwrap(),
        meta: include_bytes!("../../../../tests/fixtures/hyperliquid/meta.json"),
        spot_meta: &serde_json::to_vec(&meta).unwrap(),
    })
}
fn fixtures() -> (Value, Value) {
    (
        serde_json::from_str(include_str!(
            "../../../../tests/fixtures/hyperliquid/spot/fills.json"
        ))
        .unwrap(),
        serde_json::from_str(include_str!(
            "../../../../tests/fixtures/hyperliquid/spot/meta.json"
        ))
        .unwrap(),
    )
}
#[test]
fn aliases_deduplicate_and_base_fee_does_not_change_gross_quantity() {
    let (mut f, m) = fixtures();
    let mut alias = f[0].clone();
    alias["coin"] = json!("@0");
    f.as_array_mut().unwrap().push(alias);
    let p = parse(f, m).unwrap();
    assert_eq!(p.trades.len(), 1);
    assert_eq!(p.duplicate_records, 1);
    assert_eq!(p.spot_records, 1);
    let t = &p.trades[0];
    assert_eq!(t.payload.instrument_type, "SPOT");
    assert_eq!(t.payload.quantity, "12.5");
    assert_eq!(t.payload.notional, "1.543209750");
    assert_eq!(t.payload.fee_asset.as_deref(), Some("PURR"));
    assert_eq!(t.payload.action, "BUY");
    assert_eq!(t.payload.position_effect, "NOT_APPLICABLE");
    assert!(t.payload.copy_eligible);
    assert_eq!(t.payload.reported_realized_pnl, None);
}
#[test]
fn sell_rebate_and_non_usdc_quote_are_precise() {
    let (mut f, mut m) = fixtures();
    m["tokens"][0]["name"] = json!("USDH");
    f[0]["side"] = json!("A");
    f[0]["dir"] = json!("Sell");
    f[0]["feeToken"] = json!("USDH");
    f[0]["fee"] = json!("-0.00001");
    let p = parse(f, m).unwrap();
    let t = &p.trades[0];
    assert_eq!(t.payload.action, "SELL");
    assert_eq!(t.payload.quote_asset, "USDH");
    assert_eq!(t.payload.fee, "-0.00001");
    assert!(t.payload.copy_eligible);
}
#[test]
fn missing_or_ambiguous_fee_and_forced_messages_cannot_copy() {
    let (f, m) = fixtures();
    for variant in 0..4 {
        let mut f = f.clone();
        let mut m = m.clone();
        match variant {
            0 => {
                f[0].as_object_mut().unwrap().remove("feeToken");
            }
            1 => {
                let mut t = m["tokens"][1].clone();
                t["index"] = json!(3);
                t["tokenId"] = json!("0x00000000000000000000000000000009");
                m["tokens"].as_array_mut().unwrap().push(t);
            }
            2 => f[0]["liquidation"] = json!({"method":"liquidation"}),
            _ => {
                f[0]["feeToken"] = json!("TOKEN");
                m["tokens"][2]["tokenId"] = json!("");
            }
        };
        let p = parse(f, m).unwrap();
        assert!(!p.trades[0].payload.copy_eligible);
    }
}
#[test]
fn bad_metadata_and_conflicting_aliases_do_not_silently_skip() {
    let (f, mut m) = fixtures();
    m["universe"][0]["tokens"] = json!([999, 0]);
    assert!(parse(f.clone(), m).is_err());
    let (mut f, m) = fixtures();
    let mut second = f[0].clone();
    second["coin"] = json!("@0");
    second["sz"] = json!("20");
    f.as_array_mut().unwrap().push(second);
    assert_eq!(parse(f, m).err().unwrap().code, "VERSION_CONFLICT");
}
