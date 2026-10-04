use crate::query::QueryError;
use account_facts::AccountFact;
use shared_types::Network;

pub struct ParseContext<'a> {
    pub network: &'a Network,
    pub account: &'a str,
    pub raw_log_id: &'a str,
    pub fills: &'a [u8],
    pub meta: &'a [u8],
    pub spot_meta: &'a [u8],
}
pub struct Parsed {
    pub source_records: usize,
    pub duplicate_records: usize,
    pub spot_records: usize,
    pub unsupported_records: usize,
    pub invalid_records: usize,
    pub trades: Vec<AccountFact>,
    pub warnings: Vec<String>,
}
pub trait ProtocolParser: Send + Sync {
    fn parse(&self, context: ParseContext<'_>) -> Result<Parsed, QueryError>;
}
