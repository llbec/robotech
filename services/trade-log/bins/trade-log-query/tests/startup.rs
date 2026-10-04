use std::{collections::BTreeMap, process::Command};
use trade_log_query::config::Config;
const CONFIG: &str = include_str!("../../../../../config/trade-log.toml");
#[test]
fn service_configuration_and_cli_failures() {
    Config::parse(CONFIG, &BTreeMap::new()).unwrap();
    for (from, to) in [
        ("network = \"mainnet\"", "network = \"invalid\""),
        ("max_concurrency = 4", "max_concurrency = 0"),
        ("timeout_seconds = 30", "timeout_seconds = 0"),
        ("max_response_bytes = 16777216", "max_response_bytes = 0"),
        (
            "directory = \"/var/lib/robotech/trade-log\"",
            "directory = \"\"",
        ),
    ] {
        assert!(Config::parse(&CONFIG.replace(from, to), &BTreeMap::new()).is_err());
    }
    let command = || {
        let mut c = Command::new(env!("CARGO_BIN_EXE_trade-log-query"));
        for key in service_runtime::config::ENV_KEYS {
            c.env_remove(key);
        }
        c
    };
    for option in ["--help", "--version"] {
        assert!(
            command()
                .args([option, "--config", "/missing"])
                .output()
                .unwrap()
                .status
                .success()
        );
    }
    let missing = command().args(["--config", "/missing"]).output().unwrap();
    assert_eq!(missing.status.code(), Some(2));
    let directory = std::env::temp_dir().join(format!("robotech-config-{}", uuid::Uuid::new_v4()));
    std::fs::create_dir(&directory).unwrap();
    let config_path = directory.join("config.toml");
    std::fs::write(
        &config_path,
        CONFIG.replace(
            "/run/secrets/trade-log-token",
            directory.join("missing-token").to_str().unwrap(),
        ),
    )
    .unwrap();
    let missing = command()
        .arg("--config")
        .arg(&config_path)
        .output()
        .unwrap();
    assert_eq!(missing.status.code(), Some(2));
    assert!(!String::from_utf8_lossy(&missing.stderr).contains("server_started"));
    std::fs::remove_dir_all(directory).unwrap();
}
