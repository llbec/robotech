use std::collections::BTreeMap;
use trade_collector::config::Config;
fn valid() -> String {
    include_str!("../../../../../config/trade-collector.toml")
        .replace(
            "<实际账户地址>",
            "0x0000000000000000000000000000000000000001",
        )
        .replace("<RFC3339起点>", "2026-10-04T07:00:00Z")
}
#[test]
fn configuration_and_process_failure_contract() {
    let text = valid();
    assert!(Config::parse(&text, &BTreeMap::new()).is_ok());
    for (a, b) in [
        ("lease_seconds = 180", "lease_seconds = 120"),
        ("interval_seconds = 30", "interval_seconds = 0"),
        ("overlap_seconds = 60", "overlap_seconds = 3600"),
        ("max_requests_per_round = 20", "max_requests_per_round = 0"),
        ("retry_max_seconds = 300", "retry_max_seconds = 1"),
        ("2026-10-04T07:00:00Z", "2026-10-04T07:00:00.0001Z"),
        ("max_response_bytes = 16777216", "max_response_bytes = 1"),
    ] {
        assert!(
            Config::parse(&text.replace(a, b), &BTreeMap::new()).is_err(),
            "{b}"
        );
    }
    assert!(
        Config::parse(
            include_str!("../../../../../config/trade-collector.toml"),
            &BTreeMap::new()
        )
        .is_err()
    );
    let out = std::process::Command::new(env!("CARGO_BIN_EXE_trade-collector"))
        .args(["--config", "/nonexistent/collector.toml"])
        .output()
        .unwrap();
    assert_eq!(out.status.code(), Some(2));
    let out = std::process::Command::new(env!("CARGO_BIN_EXE_trade-collector"))
        .arg("--version")
        .output()
        .unwrap();
    assert!(out.status.success());
    assert!(
        String::from_utf8(out.stdout)
            .unwrap()
            .contains(env!("CARGO_PKG_VERSION"))
    );
}

#[test]
fn optional_websocket_config_keeps_legacy_http_mode_and_validates_limits() {
    let legacy = valid().split("[websocket]").next().unwrap().to_owned();
    assert!(
        !Config::parse(&legacy, &BTreeMap::new())
            .unwrap()
            .websocket
            .enabled
    );
    let text = valid();
    for (a, b) in [
        ("reconnect_base_seconds = 5", "reconnect_base_seconds = 0"),
        ("reconnect_max_seconds = 60", "reconnect_max_seconds = 61"),
        (
            "reconnect_reset_after_seconds = 60",
            "reconnect_reset_after_seconds = 59",
        ),
        (
            "reconnect_jitter_percent = 20",
            "reconnect_jitter_percent = 21",
        ),
        ("pong_timeout_seconds = 10", "pong_timeout_seconds = 40"),
        ("max_pending_messages = 256", "max_pending_messages = 0"),
        (
            "commit_timeout_seconds = 10",
            "commit_timeout_seconds = 180",
        ),
    ] {
        assert!(
            Config::parse(&text.replace(a, b), &BTreeMap::new()).is_err(),
            "{b}"
        );
    }
}
