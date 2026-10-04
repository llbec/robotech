use hyperliquid::parser::HyperliquidParser;
use serde_json::{Value, json};
use shared_types::Network;
use trade_log::{
    parsing::{ParseContext, Parsed, ProtocolParser},
    query::QueryError,
};
const META: &[u8] = include_bytes!("../../../../tests/fixtures/hyperliquid/meta.json");
const SPOT: &[u8] = include_bytes!("../../../../tests/fixtures/hyperliquid/spot-meta.json");
fn fill() -> Value {
    serde_json::from_slice::<Value>(include_bytes!(
        "../../../../tests/fixtures/hyperliquid/fills.json"
    ))
    .unwrap()[0]
        .clone()
}
fn parse(fills: Vec<Value>) -> Result<Parsed, QueryError> {
    let bytes = serde_json::to_vec(&fills).unwrap();
    HyperliquidParser.parse(ParseContext {
        network: &Network::Mainnet,
        account: "0x0000000000000000000000000000000000000001",
        raw_log_id: "raw_test",
        fills: &bytes,
        meta: META,
        spot_meta: SPOT,
    })
}
#[test]
fn mixed_records_dedup_and_precise_values() {
    let f = fill();
    let mut spot = f.clone();
    spot["coin"] = "@1".into();
    let mut unsupported = f.clone();
    unsupported["coin"] = "xyz:ETH".into();
    let mut invalid = f.clone();
    invalid["tid"] = 2.into();
    invalid["px"] = "invalid".into();
    let result = parse(vec![f.clone(), spot, unsupported, invalid, f]).unwrap();
    assert_eq!(
        (
            result.source_records,
            result.duplicate_records,
            result.spot_records,
            result.unsupported_records,
            result.invalid_records,
            result.trades.len()
        ),
        (5, 1, 1, 1, 1, 1)
    );
    let fact = &result.trades[0];
    assert_eq!(fact.source_ref, "9007199254740993");
    assert_eq!(fact.payload.order_id.as_deref(), Some("9007199254740995"));
    assert_eq!(fact.payload.notional, "250.025");
    assert_eq!(fact.payload.fee, "-0.0125");
    assert_eq!(fact.payload.extension["source_indices"], json!([0, 4]));
    assert!(fact.payload.copy_eligible);
}
#[test]
fn position_effect_and_forced_trade_semantics() {
    for (start, side, size, effect) in [
        ("-1", "B", "1", "CLOSE"),
        ("1", "A", "1", "CLOSE"),
        ("1", "A", "2", "REVERSE"),
        ("1", "B", "1", "INCREASE"),
        ("-2", "B", "1", "DECREASE"),
    ] {
        let mut f = fill();
        f["startPosition"] = start.into();
        f["side"] = side.into();
        f["sz"] = size.into();
        assert_eq!(
            parse(vec![f]).unwrap().trades[0].payload.position_effect,
            effect
        );
    }
    let mut f = fill();
    f["liquidation"] = json!({"liquidator":"0xexample"});
    let result = parse(vec![f]).unwrap();
    assert!(!result.trades[0].payload.copy_eligible);
    assert_eq!(result.trades[0].payload.trigger_type, "LIQUIDATION");
    let mut f = fill();
    f["dir"] = "Settlement".into();
    assert!(!parse(vec![f]).unwrap().trades[0].payload.copy_eligible);
}
#[test]
fn conflict_missing_identity_and_precision_are_explicit() {
    let f = fill();
    let mut other = f.clone();
    other["sz"] = "2".into();
    assert_eq!(
        parse(vec![f, other]).err().unwrap().code,
        "VERSION_CONFLICT"
    );
    let mut f = fill();
    f.as_object_mut().unwrap().remove("tid");
    assert_eq!(parse(vec![f]).unwrap().invalid_records, 1);
    for key in ["px", "sz"] {
        let mut f = fill();
        f[key] = "0".into();
        assert_eq!(parse(vec![f]).unwrap().invalid_records, 1);
    }
    let mut f = fill();
    f["px"] = "0.00000000000000000000000000001".into();
    assert_eq!(parse(vec![f]).unwrap().invalid_records, 1);
    let mut f = fill();
    f["px"] = "79228162514264337593543950335".into();
    f["sz"] = "2".into();
    assert_eq!(parse(vec![f]).unwrap().invalid_records, 1);
    let mut f = fill();
    f.as_object_mut().unwrap().remove("feeToken");
    assert!(
        parse(vec![f]).unwrap().trades[0]
            .payload
            .fee_asset
            .is_none()
    );
}
#[test]
fn empty_and_corrupt_metadata() {
    assert_eq!(parse(vec![]).unwrap().source_records, 0);
    assert!(
        HyperliquidParser
            .parse(ParseContext {
                network: &Network::Mainnet,
                account: "a",
                raw_log_id: "r",
                fills: b"[]",
                meta: b"{}",
                spot_meta: SPOT
            })
            .is_err()
    );
}
