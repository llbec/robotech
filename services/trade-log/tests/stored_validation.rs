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
