use chrono::{DateTime, Utc};
use thiserror::Error;

#[derive(Debug, Error, PartialEq)]
pub enum ValidationError {
    #[error("address must be a 42-character 0x-prefixed hexadecimal value")]
    InvalidAddress,
    #[error("from must be earlier than to")]
    InvalidTimeRange,
}

pub fn validate_address(address: &str) -> Result<String, ValidationError> {
    let valid = address.len() == 42
        && address.starts_with("0x")
        && address[2..].chars().all(|c| c.is_ascii_hexdigit());
    valid
        .then(|| address.to_ascii_lowercase())
        .ok_or(ValidationError::InvalidAddress)
}

pub fn validate_range(from: DateTime<Utc>, to: DateTime<Utc>) -> Result<(), ValidationError> {
    (from < to)
        .then_some(())
        .ok_or(ValidationError::InvalidTimeRange)
}
