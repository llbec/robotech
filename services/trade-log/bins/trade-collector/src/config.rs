use serde::Deserialize;
use service_runtime::config::{LoggingConfig, ServerConfig};
use std::{
    collections::BTreeMap,
    path::{Path, PathBuf},
};
use trade_log::collection::{CollectionConfig, start_ms};
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Config {
    pub config_version: u32,
    pub server: ServerConfig,
    pub logging: LoggingConfig,
    pub hyperliquid: SourceConfig,
    pub database: DatabaseConfig,
    pub evidence: EvidenceConfig,
    pub internal: InternalConfig,
    pub collection: CollectionConfig,
    #[serde(default)]
    pub websocket: trade_log::realtime::WebsocketConfig,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SourceConfig {
    pub network: shared_types::Network,
    pub connect_timeout_seconds: u64,
    pub request_timeout_seconds: u64,
    pub max_response_bytes: usize,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct DatabaseConfig {
    pub url_file: PathBuf,
    pub max_connections: u32,
    pub connect_timeout_seconds: u64,
    pub statement_timeout_seconds: u64,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct EvidenceConfig {
    pub directory: PathBuf,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct InternalConfig {
    pub credential_file: PathBuf,
}
impl Config {
    pub fn load(path: &Path, env: &BTreeMap<String, String>) -> Result<Self, String> {
        let s = std::fs::read_to_string(path).map_err(|_| "cannot read collector configuration")?;
        Self::parse(&s, env)
    }
    pub fn parse(s: &str, env: &BTreeMap<String, String>) -> Result<Self, String> {
        let mut c: Self = toml::from_str(s)
            .map_err(|e: toml::de::Error| format!("invalid configuration: {}", e.message()))?;
        service_runtime::config::apply(&mut c.server, &mut c.logging, env)?;
        service_runtime::config::validate(&c.server, &c.logging)?;
        if c.config_version != 1 {
            return Err("config_version must be 1".into());
        }
        c.collection.account = shared_types::account(&c.collection.account)?;
        start_ms(&c.collection.start_time).map_err(|e| e.message)?;
        let t = &c.collection;
        if !(1..=86400).contains(&t.interval_seconds)
            || !(1..=86400).contains(&t.overlap_seconds)
            || !(1..=86400).contains(&t.max_window_seconds)
            || t.overlap_seconds >= t.max_window_seconds
            || t.safety_delay_seconds > 3600
            || !(1..=1000).contains(&t.max_requests_per_round)
            || !(1..=3600).contains(&t.round_timeout_seconds)
            || !(1..=7200).contains(&t.lease_seconds)
            || t.lease_seconds <= t.round_timeout_seconds
            || !(1..=3600).contains(&t.retry_base_seconds)
            || !(1..=86400).contains(&t.retry_max_seconds)
            || t.retry_base_seconds > t.retry_max_seconds
        {
            return Err("Invalid collection timing or request limits".into());
        }
        if !(1..=60).contains(&c.hyperliquid.connect_timeout_seconds)
            || !(1..=60).contains(&c.hyperliquid.request_timeout_seconds)
            || c.hyperliquid.connect_timeout_seconds > c.hyperliquid.request_timeout_seconds
            || !(1024..=67108864).contains(&c.hyperliquid.max_response_bytes)
        {
            return Err("Invalid source limits".into());
        }
        if c.database.url_file.as_os_str().is_empty()
            || !(1..=64).contains(&c.database.max_connections)
            || !(1..=60).contains(&c.database.connect_timeout_seconds)
            || !(1..=120).contains(&c.database.statement_timeout_seconds)
        {
            return Err("Invalid database configuration".into());
        }
        if c.internal.credential_file.as_os_str().is_empty()
            || c.evidence.directory.as_os_str().is_empty()
        {
            return Err("Credential and evidence paths are required".into());
        }
        c.websocket.validate(c.collection.lease_seconds)?;
        Ok(c)
    }
}
