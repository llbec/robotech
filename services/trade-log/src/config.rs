use anyhow::{bail, Context, Result};
use std::{env, net::SocketAddr, str::FromStr, time::Duration};

#[derive(Clone, Debug)]
pub struct Config {
    pub database_url: String,
    pub listen_addr: SocketAddr,
    pub hyperliquid_http_url: String,
    pub hyperliquid_ws_url: String,
    pub default_history_lookback: Duration,
    pub overlap_window: Duration,
    pub reconcile_interval: Duration,
    pub max_monitored_addresses: i64,
    pub webhook_url: String,
    pub webhook_secret: String,
}

impl Config {
    pub fn from_env() -> Result<Self> {
        let database_url = required("DATABASE_URL")?;
        let webhook_url = required("EVENT_WEBHOOK_URL")?;
        let webhook_secret = required("EVENT_WEBHOOK_SECRET")?;
        if webhook_secret.len() < 16 {
            bail!("EVENT_WEBHOOK_SECRET must contain at least 16 characters");
        }
        let max_monitored_addresses = parse_or("MAX_MONITORED_ADDRESSES", 10_i64)?;
        if !(1..=10).contains(&max_monitored_addresses) {
            bail!("MAX_MONITORED_ADDRESSES must be between 1 and 10");
        }
        Ok(Self {
            database_url,
            listen_addr: env::var("HTTP_LISTEN_ADDR")
                .unwrap_or_else(|_| "0.0.0.0:8080".into())
                .parse()
                .context("HTTP_LISTEN_ADDR is invalid")?,
            hyperliquid_http_url: env::var("HYPERLIQUID_HTTP_URL")
                .unwrap_or_else(|_| "https://api.hyperliquid.xyz/info".into()),
            hyperliquid_ws_url: env::var("HYPERLIQUID_WS_URL")
                .unwrap_or_else(|_| "wss://api.hyperliquid.xyz/ws".into()),
            default_history_lookback: Duration::from_secs(parse_or(
                "DEFAULT_HISTORY_LOOKBACK_SECONDS",
                7 * 24 * 60 * 60,
            )?),
            overlap_window: Duration::from_secs(parse_or("HTTP_OVERLAP_SECONDS", 120_u64)?),
            reconcile_interval: Duration::from_secs(parse_or(
                "HTTP_RECONCILE_INTERVAL_SECONDS",
                30_u64,
            )?),
            max_monitored_addresses,
            webhook_url,
            webhook_secret,
        })
    }
}

fn required(name: &str) -> Result<String> {
    let value = env::var(name).with_context(|| format!("{name} is required"))?;
    if value.trim().is_empty() {
        bail!("{name} cannot be empty");
    }
    Ok(value)
}

fn parse_or<T>(name: &str, default: T) -> Result<T>
where
    T: FromStr,
    T::Err: std::error::Error + Send + Sync + 'static,
{
    match env::var(name) {
        Ok(value) => value.parse().with_context(|| format!("{name} is invalid")),
        Err(_) => Ok(default),
    }
}
