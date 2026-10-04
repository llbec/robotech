use crate::{
    acquisition::SourceReader,
    parsing::{ParseContext, ProtocolParser},
    raw_log::RawEvidenceStore,
    validation::QueryRequest,
};
use account_facts::AccountFact;
use protocol_api::QueryKind;
use serde::{Deserialize, Serialize};
use shared_types::Network;
use std::sync::Arc;

#[derive(Clone, Debug, Serialize, Deserialize, thiserror::Error)]
#[error("{message}")]
pub struct QueryError {
    pub code: String,
    pub message: String,
    #[serde(skip)]
    pub retryable: bool,
}
impl QueryError {
    pub fn validation(message: impl Into<String>) -> Self {
        Self {
            retryable: false,
            code: "VALIDATION_ERROR".into(),
            message: message.into(),
        }
    }
    pub fn unavailable(message: impl Into<String>) -> Self {
        Self {
            retryable: false,
            code: "DEPENDENCY_UNAVAILABLE".into(),
            message: message.into(),
        }
    }
    pub fn incomplete(message: impl Into<String>) -> Self {
        Self {
            retryable: false,
            code: "INCOMPLETE_DATA".into(),
            message: message.into(),
        }
    }
    pub fn storage() -> Self {
        Self {
            retryable: false,
            code: "INTERNAL_INVARIANT_VIOLATION".into(),
            message: "Evidence storage failed".into(),
        }
    }
    pub fn limited() -> Self {
        Self {
            retryable: false,
            code: "RATE_LIMITED".into(),
            message: "Query capacity or source rate limit exceeded".into(),
        }
    }
    pub fn conflict() -> Self {
        Self {
            retryable: false,
            code: "VERSION_CONFLICT".into(),
            message: "Conflicting source fill identity".into(),
        }
    }
    pub fn with_retry(mut self) -> Self {
        self.retryable = true;
        self
    }
    pub fn status(&self) -> u16 {
        match self.code.as_str() {
            "VALIDATION_ERROR" => 400,
            "INCOMPLETE_DATA" => 422,
            "RATE_LIMITED" => 429,
            "VERSION_CONFLICT" => 409,
            "INTERNAL_INVARIANT_VIOLATION" => 500,
            _ => 503,
        }
    }
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Counts {
    pub source_records: usize,
    pub duplicate_records: usize,
    pub spot_records: usize,
    pub unsupported_records: usize,
    pub invalid_records: usize,
    pub perpetual_records: usize,
    pub returned_records: usize,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct ObservedRange {
    pub first_at: String,
    pub last_at: String,
}
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct QueryResult {
    pub query_id: String,
    pub account: String,
    pub network: String,
    pub queried_at: String,
    pub query_scope: String,
    pub coverage: String,
    pub counts: Counts,
    pub display_truncated: bool,
    pub observed_range: Option<ObservedRange>,
    pub trades: Vec<AccountFact>,
    pub warnings: Vec<String>,
    pub evidence_ref: String,
}
impl QueryResult {
    pub fn for_display(mut self) -> Self {
        self.trades.truncate(self.counts.returned_records);
        self
    }
}

pub struct QueryService {
    pub network: Network,
    pub source: Arc<dyn SourceReader>,
    pub parser: Arc<dyn ProtocolParser>,
    pub evidence: Arc<dyn RawEvidenceStore>,
}
impl QueryService {
    // The runtime adapter owns timeout/concurrency; the domain remains independent of Tokio.
    pub async fn execute(
        &self,
        id: &str,
        request: &QueryRequest,
        trace: &str,
    ) -> Result<QueryResult, QueryError> {
        let normalized = request.validated()?;
        self.evidence.begin(id, request, trace).await?;
        let request = normalized;
        let result = self.collect(id, &request).await;
        self.evidence.finish(id, result.as_ref()).await?;
        result.map(QueryResult::for_display)
    }
    async fn fetch(
        &self,
        id: &str,
        request: &QueryRequest,
        kind: QueryKind,
    ) -> Result<(Vec<u8>, String), QueryError> {
        for attempt in 1..=2 {
            let response = match self.source.fetch(kind, &request.account).await {
                Ok(response) => response,
                Err(error) if attempt == 1 && error.retryable => {
                    self.source.wait_retry(1).await;
                    continue;
                }
                Err(error) => return Err(error),
            };
            let raw_id = self
                .evidence
                .save_response(id, kind, attempt, &response, request)
                .await?;
            if response.status == 429 || response.status >= 500 {
                if attempt == 1 {
                    self.source
                        .wait_retry(response.retry_after_seconds.unwrap_or(1).max(1))
                        .await;
                    continue;
                }
                return Err(if response.status == 429 {
                    QueryError::limited()
                } else {
                    QueryError::unavailable("Source returned an unsuccessful response")
                });
            }
            if !(200..300).contains(&response.status) {
                return Err(QueryError::unavailable(
                    "Source returned an unsuccessful response",
                ));
            }
            return Ok((response.body, raw_id));
        }
        unreachable!("bounded retries return a result")
    }

    async fn collect(&self, id: &str, request: &QueryRequest) -> Result<QueryResult, QueryError> {
        let queried_at = shared_types::now();
        let (fills, raw_id) = self.fetch(id, request, QueryKind::UserFills).await?;
        let (meta, _) = self.fetch(id, request, QueryKind::Meta).await?;
        let (spot_meta, _) = self.fetch(id, request, QueryKind::SpotMeta).await?;
        let parsed = self.parser.parse(ParseContext {
            network: &self.network,
            account: &request.account,
            raw_log_id: &raw_id,
            fills: &fills,
            meta: &meta,
            spot_meta: &spot_meta,
        })?;
        crate::normalization::result(
            id,
            &request.account,
            &self.network,
            &queried_at,
            request.limit,
            parsed,
        )
    }
}
