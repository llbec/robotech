use trade_log::stored_query::{StoredRequest, time};
#[test]
fn exact_time_zone_and_request_validation() {
    assert_eq!(
        time("2026-10-05T08:00:00.001+08:00").unwrap(),
        "2026-10-05T00:00:00.001Z"
    );
    for value in [
        "2026-10-05",
        "2026-10-05T00:00:00",
        "2026-10-05T00:00:00.000001Z",
        "1969-12-31T23:59:59Z",
        "2026-10-05T00:00:00.0000Z",
    ] {
        assert!(time(value).is_err());
    }
    let mut r = StoredRequest {
        instrument_type: Default::default(),
        account: "0x000000000000000000000000000000000000000A".into(),
        limit: 100,
        start_time: Some("2026-10-05T00:00:00Z".into()),
        end_time: Some("2026-10-05T00:00:00Z".into()),
        cursor: None,
    };
    assert!(r.validated().is_err());
    r.end_time = None;
    assert!(r.validated().unwrap().account.ends_with('a'));
    r.limit = 0;
    assert!(r.validated().is_err());
}

#[test]
fn legacy_cursor_defaults_to_perpetual_and_cannot_cross_instrument_type() {
    use base64::{Engine, engine::general_purpose::URL_SAFE_NO_PAD};
    let request: trade_log::stored_query::StoredRequest = serde_json::from_value(
        serde_json::json!({"account":"0x0000000000000000000000000000000000000001","limit":1}),
    )
    .unwrap();
    let old = serde_json::json!({"version":1,"account":request.account,"network":"mainnet","range":{"start_time":null,"end_time":null},"snapshot":1,"time":"2026-10-01T00:00:00.000Z","tid":"1","fact_id":format!("hl_fill_v1_{}","a".repeat(64))});
    let mut request = request;
    request.cursor = Some(URL_SAFE_NO_PAD.encode(serde_json::to_vec(&old).unwrap()));
    assert!(request.decode_cursor("mainnet").is_ok());
    request.instrument_type = trade_log::instrument::InstrumentType::Spot;
    assert!(request.decode_cursor("mainnet").is_err());
}
