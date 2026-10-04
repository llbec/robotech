use serde::Deserialize;
use service_runtime::config::{LoggingConfig, ServerConfig};
use shared_types::Network;
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
    pub hyperliquid: HyperliquidConfig,
    pub query: QueryConfig,
    pub evidence: EvidenceConfig,
    pub internal: InternalConfig,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct HyperliquidConfig {
    pub network: Network,
    pub connect_timeout_seconds: u64,
    pub request_timeout_seconds: u64,
}
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct QueryConfig {
    pub timeout_seconds: u64,
    pub max_concurrency: usize,
    pub max_response_bytes: usize,
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
        let text =
            std::fs::read_to_string(path).map_err(|_| "cannot read trade-log configuration")?;
        Self::parse(&text, env)
    }
    pub fn parse(text: &str, env: &BTreeMap<String, String>) -> Result<Self, String> {
        let mut config: Self = toml::from_str(text)
            .map_err(|e: toml::de::Error| format!("invalid configuration: {}", e.message()))?;
        service_runtime::config::apply(&mut config.server, &mut config.logging, env)?;
        if config.config_version != 1 {
            return Err("config_version must be 1".into());
        }
        service_runtime::config::validate(&config.server, &config.logging)?;
        if !(1..=60).contains(&config.hyperliquid.connect_timeout_seconds)
            || !(1..=60).contains(&config.hyperliquid.request_timeout_seconds)
        {
            return Err("hyperliquid timeouts must be 1–60 seconds".into());
        }
        if config.hyperliquid.connect_timeout_seconds > config.hyperliquid.request_timeout_seconds {
            return Err("connect timeout must not exceed request timeout".into());
        }
        if !(1..=120).contains(&config.query.timeout_seconds) {
            return Err("query.timeout_seconds must be 1–120".into());
        }
        if !(1..=64).contains(&config.query.max_concurrency) {
            return Err("query.max_concurrency must be 1–64".into());
        }
        if !(1024..=64 * 1024 * 1024).contains(&config.query.max_response_bytes) {
            return Err("query.max_response_bytes must be 1024–67108864".into());
        }
        if config.evidence.directory.as_os_str().is_empty()
            || config.internal.credential_file.as_os_str().is_empty()
        {
            return Err("evidence.directory and internal.credential_file are required".into());
        }
        Ok(config)
    }
}
