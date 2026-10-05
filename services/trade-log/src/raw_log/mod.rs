use crate::{
    acquisition::SourceResponse,
    query::{QueryError, QueryResult},
    validation::QueryRequest,
};
use async_trait::async_trait;
use protocol_api::QueryKind;
#[async_trait]
pub trait RawEvidenceStore: Send + Sync {
    async fn persistence(
        &self,
        _id: &str,
    ) -> Result<Option<crate::persistence::Persistence>, QueryError> {
        Ok(None)
    }
    async fn begin(&self, id: &str, request: &QueryRequest, trace: &str) -> Result<(), QueryError>;
    async fn save_response(
        &self,
        id: &str,
        kind: QueryKind,
        attempt: usize,
        response: &SourceResponse,
        request: &QueryRequest,
    ) -> Result<String, QueryError>;
    async fn finish(
        &self,
        id: &str,
        result: Result<&QueryResult, &QueryError>,
    ) -> Result<(), QueryError>;
}
