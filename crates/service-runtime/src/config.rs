use serde::Deserialize;
use std::{
    collections::BTreeMap,
    net::{IpAddr, SocketAddr},
    time::Duration,
};

pub const ENV_KEYS: [&str; 5] = [
    "ROBOTECH_SERVER_HOST",
    "ROBOTECH_SERVER_PORT",
    "ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS",
    "ROBOTECH_LOG_LEVEL",
    "ROBOTECH_LOG_FORMAT",
];

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

pub fn validate(server: &ServerConfig, logging: &LoggingConfig) -> Result<(), String> {
    if server.port == 0 {
        return Err("server.port must be 1–65535".into());
    }
    if !(1..=60).contains(&server.shutdown_timeout_seconds) {
        return Err("server.shutdown_timeout_seconds must be 1–60".into());
    }
    if !["trace", "debug", "info", "warn", "error"].contains(&logging.level.as_str()) {
        return Err("logging.level must be trace, debug, info, warn or error".into());
    }
    if !["json", "text"].contains(&logging.format.as_str()) {
        return Err("logging.format must be json or text".into());
    }
    Ok(())
}
pub fn apply(
    server: &mut ServerConfig,
    logging: &mut LoggingConfig,
    overrides: &BTreeMap<String, String>,
) -> Result<(), String> {
    for (key, value) in overrides {
        let invalid = || format!("invalid environment override {key}");
        match key.as_str() {
            "ROBOTECH_SERVER_HOST" => server.host = value.parse().map_err(|_| invalid())?,
            "ROBOTECH_SERVER_PORT" => server.port = value.parse().map_err(|_| invalid())?,
            "ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS" => {
                server.shutdown_timeout_seconds = value.parse().map_err(|_| invalid())?
            }
            "ROBOTECH_LOG_LEVEL" => logging.level = value.clone(),
            "ROBOTECH_LOG_FORMAT" => logging.format = value.clone(),
            _ => {}
        }
    }
    Ok(())
}
pub fn environment() -> Result<BTreeMap<String, String>, String> {
    let mut values = BTreeMap::new();
    for key in ENV_KEYS {
        match std::env::var(key) {
            Ok(value) => {
                values.insert(key.to_owned(), value);
            }
            Err(std::env::VarError::NotPresent) => {}
            Err(_) => return Err(format!("environment override {key} is not UTF-8")),
        }
    }
    Ok(values)
}
impl ServerConfig {
    pub fn address(&self) -> SocketAddr {
        SocketAddr::new(self.host, self.port)
    }
    pub fn shutdown_timeout(&self) -> Duration {
        Duration::from_secs(self.shutdown_timeout_seconds)
    }
}
