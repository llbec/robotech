use std::collections::BTreeMap;
use trade_parser_publisher::config::Config;
fn valid() -> String {
    include_str!("../../../../../config/trade-parser-publisher.toml").into()
}
#[test]
fn validates_config_and_cli() {
    let text = valid();
    assert!(Config::parse(&text, &BTreeMap::new()).is_ok());
    for (old, new) in [
        ("lease_seconds = 30", "lease_seconds = 10"),
        ("signal_ttl_seconds = 60", "signal_ttl_seconds = 1"),
        ("poll_interval_ms = 500", "poll_interval_ms = 0"),
        ("retry_base_seconds = 5", "retry_base_seconds = 301"),
        (
            "shutdown_timeout_seconds = 25",
            "shutdown_timeout_seconds = 10",
        ),
    ] {
        assert!(Config::parse(&text.replace(old, new), &BTreeMap::new()).is_err());
    }
    let enabled = text.replace("enabled = false", "enabled = true");
    assert!(Config::parse(&enabled, &BTreeMap::new()).is_err());
    for url in [
        "http://example.invalid/events",
        "https://u:p@example.invalid/events",
        "https://example.invalid/events?token=secret",
        "https://example.invalid/#secret",
    ] {
        assert!(
            Config::parse(
                &enabled.replace("webhook_url = \"\"", &format!("webhook_url = \"{url}\"")),
                &BTreeMap::new()
            )
            .is_err()
        );
    }
    assert!(
        Config::parse(
            &enabled.replace(
                "webhook_url = \"\"",
                "webhook_url = \"https://example.invalid/events\""
            ),
            &BTreeMap::new()
        )
        .is_ok()
    );
    let binary = env!("CARGO_BIN_EXE_trade-parser-publisher");
    assert_eq!(
        std::process::Command::new(binary)
            .args(["--config", "/missing/publisher.toml"])
            .output()
            .unwrap()
            .status
            .code(),
        Some(2)
    );
    assert!(
        std::process::Command::new(binary)
            .arg("--version")
            .output()
            .unwrap()
            .status
            .success()
    );
}
