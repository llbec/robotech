use serde::{Deserialize, Serialize};
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Persistence {
    pub status: String,
    pub inserted_records: usize,
    pub existing_records: usize,
}

/// Used by future collectors after archiving their source responses.
#[async_trait::async_trait]
pub trait FactStore: Send + Sync {
    async fn persist(
        &self,
        query_id: &str,
        complete: &crate::query::QueryResult,
    ) -> Result<Persistence, crate::query::QueryError>;
}

#[async_trait::async_trait]
pub trait CollectionFactStore: Send + Sync {
    async fn persist_collection(
        &self,
        complete: &crate::query::QueryResult,
        commit: &crate::checkpoint::CollectionCommit<'_>,
    ) -> Result<Persistence, crate::query::QueryError>;
}
