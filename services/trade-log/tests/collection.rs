use trade_log::{
    checkpoint::ScanRange,
    collection::{CollectionConfig, scan_range, start_ms},
};
fn config() -> CollectionConfig {
    serde_json::from_str(include_str!(
        "../../../tests/fixtures/v0.3/collection-cases.json"
    ))
    .unwrap()
}
#[test]
fn quiet_account_uses_scan_boundary_and_overlap_without_rewinding() {
    let c = config();
    let first = scan_range(1000, None, 10000, &c).unwrap().unwrap();
    assert_eq!(
        first,
        ScanRange {
            start_ms: 1000,
            end_ms: 8000
        }
    );
    let next = scan_range(1000, Some(8000), 15000, &c).unwrap().unwrap();
    assert_eq!(
        next,
        ScanRange {
            start_ms: 6000,
            end_ms: 13000
        }
    );
    assert!(scan_range(1000, Some(8000), 9000, &c).unwrap().is_none());
    assert_eq!(
        scan_range(1000, Some(8000), 7000, &c).unwrap_err().code,
        "CLOCK_BEHIND_WATERMARK"
    );
}
#[test]
fn millisecond_split_is_disjoint_and_saturation_is_visible() {
    let r = ScanRange {
        start_ms: 10,
        end_ms: 15,
    };
    let (a, b) = r.split().unwrap();
    assert_eq!(a.start_ms, 10);
    assert_eq!(a.end_ms, b.start_ms);
    assert_eq!(b.end_ms, 15);
    assert_eq!(
        ScanRange {
            start_ms: 10,
            end_ms: 11
        }
        .split()
        .unwrap_err()
        .code,
        "SOURCE_WINDOW_SATURATED"
    );
}
#[test]
fn checkpoint_start_has_the_same_precision_rules_as_stored_queries() {
    assert_eq!(
        start_ms("2026-10-04T15:00:00.001+08:00").unwrap(),
        1791097200001
    );
    for s in [
        "2026-10-04T07:00:00.0001Z",
        "2026-10-04T07:00:00",
        "1969-12-31T23:59:59Z",
    ] {
        assert!(start_ms(s).is_err());
    }
}
