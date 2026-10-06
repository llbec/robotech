use serde::Deserialize;
use service_runtime::config::{LoggingConfig, ServerConfig};
use std::{
    collections::BTreeMap,
    path::{Path, PathBuf},
};
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Config {
    pub config_version: u32,
    pub server: ServerConfig,
    pub logging: LoggingConfig,
    pub internal: Internal,
    pub database: Database,
    pub publishing: Publishing,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Internal {
    pub credential_file: PathBuf,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Database {
    pub url_file: PathBuf,
    pub max_connections: u32,
    pub connect_timeout_seconds: u64,
    pub statement_timeout_seconds: u64,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Publishing {
    pub enabled: bool,
    pub collection_config_path: PathBuf,
    pub webhook_url: String,
    pub credential_file: PathBuf,
    pub allow_plain_http: bool,
    pub request_timeout_seconds: u64,
    pub max_response_bytes: usize,
    pub lease_seconds: u64,
    pub poll_interval_ms: u64,
    pub retry_base_seconds: u64,
    pub retry_max_seconds: u64,
    pub max_event_age_seconds: i64,
    pub signal_ttl_seconds: i64,
    pub clock_skew_tolerance_seconds: i64,
}
impl Config {
    pub fn load(path: &Path, env: &BTreeMap<String, String>) -> Result<Self, String> {
        Self::parse(
            &std::fs::read_to_string(path).map_err(|_| "Cannot read publisher configuration")?,
            env,
        )
    }
    pub fn parse(text: &str, env: &BTreeMap<String, String>) -> Result<Self, String> {
        let mut c: Self =
            toml::from_str(text).map_err(|_: toml::de::Error| "Invalid publisher configuration")?;
        service_runtime::config::apply(&mut c.server, &mut c.logging, env)?;
        service_runtime::config::validate(&c.server, &c.logging)?;
        let p = &c.publishing;
        if c.config_version != 1
            || !(1..=64).contains(&c.database.max_connections)
            || !(1..=60).contains(&c.database.connect_timeout_seconds)
            || !(1..=120).contains(&c.database.statement_timeout_seconds)
            || c.database.url_file.as_os_str().is_empty()
            || c.internal.credential_file.as_os_str().is_empty()
            || p.collection_config_path.as_os_str().is_empty()
        {
            return Err("Invalid publisher database or credential configuration".into());
        }
        if !(1..=120).contains(&p.request_timeout_seconds)
            || p.lease_seconds <= p.request_timeout_seconds + c.database.statement_timeout_seconds
            || p.lease_seconds > 7200
            || c.server.shutdown_timeout_seconds
                < p.request_timeout_seconds + c.database.statement_timeout_seconds
            || !(50..=60000).contains(&p.poll_interval_ms)
            || !(1024..=1048576).contains(&p.max_response_bytes)
            || p.retry_base_seconds == 0
            || p.retry_base_seconds > p.retry_max_seconds
            || p.retry_max_seconds > 86400
            || !(1..=3600).contains(&p.max_event_age_seconds)
            || p.signal_ttl_seconds < p.max_event_age_seconds
            || p.signal_ttl_seconds > 86400
            || !(0..=60).contains(&p.clock_skew_tolerance_seconds)
        {
            return Err("Invalid publishing timing or limits".into());
        }
        if p.enabled || !p.webhook_url.is_empty() {
            let u = reqwest::Url::parse(&p.webhook_url).map_err(|_| "Invalid webhook URL")?;
            if !(u.scheme() == "https" || (p.allow_plain_http && u.scheme() == "http"))
                || u.host_str().is_none()
                || !u.username().is_empty()
                || u.password().is_some()
                || u.query().is_some()
                || u.fragment().is_some()
                || p.credential_file.as_os_str().is_empty()
            {
                return Err("Invalid webhook URL or credential path".into());
            }
        }
        Ok(c)
    }
    pub fn account(&self) -> Result<(String, shared_types::Network), String> {
        let text = std::fs::read_to_string(&self.publishing.collection_config_path)
            .map_err(|_| "Cannot read collector configuration")?;
        let value: toml::Value =
            toml::from_str(&text).map_err(|_| "Invalid collector configuration")?;
        let account = shared_types::account(
            value
                .get("collection")
                .and_then(|v| v.get("account"))
                .and_then(|v| v.as_str())
                .ok_or("Collector account required")?,
        )?;
        let network: shared_types::Network = value
            .get("hyperliquid")
            .and_then(|v| v.get("network"))
            .ok_or("Collector network required")?
            .clone()
            .try_into()
            .map_err(|_| "Invalid collector network")?;
        Ok((account, network))
    }
    pub fn policy(&self) -> trade_log::publishing::CandidatePolicy {
        let p = &self.publishing;
        trade_log::publishing::CandidatePolicy {
            max_event_age_seconds: p.max_event_age_seconds,
            signal_ttl_seconds: p.signal_ttl_seconds,
            clock_skew_tolerance_seconds: p.clock_skew_tolerance_seconds,
            version: if (
                p.max_event_age_seconds,
                p.signal_ttl_seconds,
                p.clock_skew_tolerance_seconds,
            ) == (30, 60, 5)
            {
                "candidate-v1".into()
            } else {
                format!(
                    "candidate-v1-age{}-ttl{}-skew{}",
                    p.max_event_age_seconds, p.signal_ttl_seconds, p.clock_skew_tolerance_seconds
                )
            },
        }
    }
}
