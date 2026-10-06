pub mod facts;
pub mod migration;
pub mod raw_log;
pub mod stored_query;
use crate::file_evidence::FileEvidence;
use shared_types::Network;
use sqlx::{PgPool, postgres::PgPoolOptions};
use std::{sync::Arc, time::Duration};
use trade_log::query::QueryError;
#[derive(Clone)]
pub struct Postgres {
    pub pool: PgPool,
    pub network: Network,
    pub mirror: Option<Arc<FileEvidence>>,
}
fn db_error(_: sqlx::Error) -> QueryError {
    let mut error = QueryError::unavailable("Trade database unavailable or operation failed");
    error.retryable = true;
    error
}
impl Postgres {
    pub async fn connect(
        url: &str,
        network: Network,
        max: u32,
        connect_seconds: u64,
        statement_seconds: u64,
    ) -> Result<Self, QueryError> {
        let pool=PgPoolOptions::new().max_connections(max).acquire_timeout(Duration::from_secs(connect_seconds))
            .after_connect(move |connection,_|Box::pin(async move {
                sqlx::query("SELECT set_config('statement_timeout',$1,false), set_config('TimeZone','UTC',false)")
                    .bind(format!("{}ms",statement_seconds*1000)).execute(connection).await?;
                Ok(())
            })).connect(url).await.map_err(db_error)?;
        Ok(Self {
            pool,
            network,
            mirror: None,
        })
    }
    pub async fn interrupt(&self) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.collection_jobs SET status='INTERRUPTED',updated_at=now() WHERE status='RUNNING' AND job_origin='MANUAL'").execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
}

pub mod checkpoint;
pub mod realtime;
pub mod replay;

pub mod publishing;
