use crate::query::QueryError;
use serde::{Deserialize, Serialize};
#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct QueryRequest {
    #[serde(default)]
    pub instrument_type: crate::instrument::InstrumentType,
    pub account: String,
    #[serde(default = "default_limit")]
    pub limit: usize,
}
fn default_limit() -> usize {
    100
}
impl QueryRequest {
    pub fn validated(&self) -> Result<Self, QueryError> {
        if !(1..=2000).contains(&self.limit) {
            return Err(QueryError::validation("limit must be 1–2000"));
        }
        Ok(Self {
            instrument_type: self.instrument_type,
            account: shared_types::account(&self.account).map_err(QueryError::validation)?,
            limit: self.limit,
        })
    }
}
