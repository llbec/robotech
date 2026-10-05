use serde::Deserialize;
pub use service_runtime::config::{ENV_KEYS, LoggingConfig, ServerConfig};
use std::{
    collections::BTreeMap,
    path::{Path, PathBuf},
};
#[derive(Debug, thiserror::Error)]
#[error("{0}")]
pub struct ConfigError(String);
#[derive(Clone, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Config {
    pub config_version: u32,
    pub server: ServerConfig,
    pub logging: LoggingConfig,
    #[serde(default)]
    pub trade_log: Option<TradeLogConfig>,
    #[serde(default)]
    pub collector: Option<TradeLogConfig>,
}
#[derive(Clone, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct TradeLogConfig {
    #[serde(default)]
    pub enabled: bool,
    pub base_url: Option<String>,
    pub credential_file: Option<PathBuf>,
    pub request_timeout_seconds: Option<u64>,
}
impl Config {
    pub fn load(path: &Path, overrides: &BTreeMap<String, String>) -> Result<Self, ConfigError> {
        let text = std::fs::read_to_string(path)
            .map_err(|e| ConfigError(format!("cannot read config {}: {e}", path.display())))?;
        Self::parse(&text, overrides)
    }
    pub fn parse(text: &str, overrides: &BTreeMap<String, String>) -> Result<Self, ConfigError> {
        let mut config: Self = toml::from_str(text).map_err(|e: toml::de::Error| {
            ConfigError(format!("invalid configuration: {}", e.message()))
        })?;
        service_runtime::config::apply(&mut config.server, &mut config.logging, overrides)
            .map_err(ConfigError)?;
        config.validate()?;
        Ok(config)
    }
    pub fn validate(&self) -> Result<(), ConfigError> {
        if self.config_version != 1 {
            return Err(ConfigError("config_version must be 1".into()));
        }
        service_runtime::config::validate(&self.server, &self.logging).map_err(ConfigError)?;
        for t in [self.trade_log.as_ref(), self.collector.as_ref()]
            .into_iter()
            .flatten()
            .filter(|t| t.enabled)
        {
            let url = t
                .base_url
                .as_deref()
                .and_then(|s| reqwest::Url::parse(s).ok())
                .ok_or_else(|| ConfigError("trade_log.base_url must be an HTTP base URL".into()))?;
            if !["http", "https"].contains(&url.scheme())
                || url.host_str().is_none()
                || !url.username().is_empty()
                || url.password().is_some()
                || url.query().is_some()
                || url.fragment().is_some()
                || url.path() != "/"
            {
                return Err(ConfigError(
                    "trade_log.base_url must contain only scheme, host and port".into(),
                ));
            }
            if t.credential_file
                .as_ref()
                .is_none_or(|v| v.as_os_str().is_empty())
            {
                return Err(ConfigError("trade_log.credential_file is required".into()));
            }
            if !matches!(t.request_timeout_seconds, Some(1..=120)) {
                return Err(ConfigError(
                    "trade_log.request_timeout_seconds must be 1–120".into(),
                ));
            }
        }
        Ok(())
    }
}
pub fn environment() -> Result<BTreeMap<String, String>, ConfigError> {
    service_runtime::config::environment().map_err(ConfigError)
}
