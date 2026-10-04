use serde::Deserialize;
use std::{
    collections::BTreeMap,
    net::{IpAddr, SocketAddr},
    path::Path,
    time::Duration,
};

pub const ENV_KEYS: [&str; 5] = [
    "ROBOTECH_SERVER_HOST",
    "ROBOTECH_SERVER_PORT",
    "ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS",
    "ROBOTECH_LOG_LEVEL",
    "ROBOTECH_LOG_FORMAT",
];

#[derive(Debug, thiserror::Error)]
#[error("{0}")]
pub struct ConfigError(String);

#[derive(Clone, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Config {
    pub config_version: u32,
    pub server: ServerConfig,
    pub logging: LoggingConfig,
}

#[derive(Clone, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ServerConfig {
    pub host: IpAddr,
    pub port: u16,
    pub shutdown_timeout_seconds: u64,
}

#[derive(Clone, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct LoggingConfig {
    pub level: String,
    pub format: String,
}

impl Config {
    pub fn load(path: &Path, overrides: &BTreeMap<String, String>) -> Result<Self, ConfigError> {
        let text = std::fs::read_to_string(path)
            .map_err(|e| ConfigError(format!("cannot read config {}: {e}", path.display())))?;
        Self::parse(&text, overrides)
    }

    pub fn parse(text: &str, overrides: &BTreeMap<String, String>) -> Result<Self, ConfigError> {
        // Do not include TOML source snippets in errors: future files may contain secrets.
        let mut config: Self = toml::from_str(text).map_err(|e: toml::de::Error| {
            ConfigError(format!("invalid configuration: {}", e.message()))
        })?;
        for (key, value) in overrides {
            let invalid = || ConfigError(format!("invalid environment override {key}"));
            match key.as_str() {
                "ROBOTECH_SERVER_HOST" => {
                    config.server.host = value.parse().map_err(|_| invalid())?
                }
                "ROBOTECH_SERVER_PORT" => {
                    config.server.port = value.parse().map_err(|_| invalid())?
                }
                "ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS" => {
                    config.server.shutdown_timeout_seconds = value.parse().map_err(|_| invalid())?
                }
                "ROBOTECH_LOG_LEVEL" => config.logging.level = value.clone(),
                "ROBOTECH_LOG_FORMAT" => config.logging.format = value.clone(),
                _ => {} // Only declared process environment keys participate.
            }
        }
        config.validate()?;
        Ok(config)
    }

    pub fn validate(&self) -> Result<(), ConfigError> {
        if self.config_version != 1 {
            return Err(ConfigError("config_version must be 1".into()));
        }
        if self.server.port == 0 {
            return Err(ConfigError("server.port must be 1–65535".into()));
        }
        if !(1..=60).contains(&self.server.shutdown_timeout_seconds) {
            return Err(ConfigError(
                "server.shutdown_timeout_seconds must be 1–60".into(),
            ));
        }
        if !["trace", "debug", "info", "warn", "error"].contains(&self.logging.level.as_str()) {
            return Err(ConfigError(
                "logging.level must be trace, debug, info, warn or error".into(),
            ));
        }
        if !["json", "text"].contains(&self.logging.format.as_str()) {
            return Err(ConfigError("logging.format must be json or text".into()));
        }
        Ok(())
    }
}

impl ServerConfig {
    pub fn address(&self) -> SocketAddr {
        SocketAddr::new(self.host, self.port)
    }
    pub fn shutdown_timeout(&self) -> Duration {
        Duration::from_secs(self.shutdown_timeout_seconds)
    }
}

pub fn environment() -> Result<BTreeMap<String, String>, ConfigError> {
    let mut values = BTreeMap::new();
    for key in ENV_KEYS {
        match std::env::var(key) {
            Ok(value) => {
                values.insert(key.to_owned(), value);
            }
            Err(std::env::VarError::NotPresent) => {}
            Err(_) => {
                return Err(ConfigError(format!(
                    "environment override {key} is not UTF-8"
                )));
            }
        }
    }
    Ok(values)
}
