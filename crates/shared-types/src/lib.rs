use chrono::{SecondsFormat, Utc};
pub use rust_decimal::Decimal;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "lowercase")]
pub enum Network {
    Mainnet,
    Testnet,
}
impl Network {
    pub fn name(&self) -> &'static str {
        match self {
            Self::Mainnet => "mainnet",
            Self::Testnet => "testnet",
        }
    }
    pub fn endpoint(&self) -> &'static str {
        match self {
            Self::Mainnet => "https://api.hyperliquid.xyz/info",
            Self::Testnet => "https://api.hyperliquid-testnet.xyz/info",
        }
    }
}
pub fn account(value: &str) -> Result<String, &'static str> {
    if value.len() != 42
        || !value.starts_with("0x")
        || !value[2..].bytes().all(|b| b.is_ascii_hexdigit())
    {
        return Err("account must be 0x followed by 40 hexadecimal digits");
    }
    Ok(value.to_ascii_lowercase())
}
pub fn now() -> String {
    Utc::now().to_rfc3339_opts(SecondsFormat::Millis, true)
}
pub fn timestamp(ms: i64) -> Result<String, &'static str> {
    if ms < 0 {
        return Err("negative timestamp");
    }
    chrono::DateTime::from_timestamp_millis(ms)
        .map(|t| t.to_rfc3339_opts(SecondsFormat::Millis, true))
        .ok_or("timestamp out of range")
}
pub fn decimal(value: &str) -> Result<Decimal, &'static str> {
    // Exact parsing rejects rounding beyond the supported 28-digit scale.
    Decimal::from_str_exact(value).map_err(|_| "invalid or unsupported decimal precision")
}
