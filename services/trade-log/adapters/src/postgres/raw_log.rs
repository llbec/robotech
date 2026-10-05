use super::*;
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use protocol_api::QueryKind;
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use trade_log::{
    acquisition::SourceResponse, persistence::Persistence, query::QueryResult,
    raw_log::RawEvidenceStore, validation::QueryRequest,
};
use uuid::Uuid;
#[async_trait]
impl RawEvidenceStore for Postgres {
    async fn begin(&self, id: &str, request: &QueryRequest, trace: &str) -> Result<(), QueryError> {
        let account = shared_types::account(&request.account).map_err(QueryError::validation)?;
        sqlx::query("INSERT INTO trade_log.collection_jobs(id,query_id,chain_id,account,network,status,request,trace_id) VALUES($1,$2,$3,$4,$5,'RUNNING',$6,$7)")
            .bind(Uuid::new_v4()).bind(id).bind(format!("hyperliquid:{}",self.network.name())).bind(account).bind(self.network.name()).bind(json!(request)).bind(trace)
            .execute(&self.pool).await.map_err(db_error)?;
        if let Some(m) = &self.mirror
            && let Err(e) = m.begin(id, request, trace).await
        {
            tracing::warn!(query_id=id,code=%e.code,"evidence_mirror_failed");
        }
        Ok(())
    }
    async fn save_response(
        &self,
        id: &str,
        kind: QueryKind,
        attempt: usize,
        response: &SourceResponse,
        request: &QueryRequest,
    ) -> Result<String, QueryError> {
        let raw_id = format!("raw_{id}_{}_{}", kind.api_name(), attempt);
        let received: DateTime<Utc> = DateTime::parse_from_rfc3339(&response.received_at)
            .map_err(|_| QueryError::storage())?
            .with_timezone(&Utc);
        let mut body = json!({"type":kind.api_name()});
        if kind == QueryKind::UserFills {
            body["user"] = request.account.clone().into();
            body["aggregateByTime"] = false.into();
        }
        let hash = format!("{:x}", Sha256::digest(&response.body));
        let inserted=sqlx::query("INSERT INTO trade_log.raw_logs(id,collection_job_id,chain_id,source_event_id,kind,attempt,chain_position,ordering_key,payload,body,http_status,request,sha256,observed_at) SELECT $1,id,chain_id,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13 FROM trade_log.collection_jobs WHERE query_id=$2 ON CONFLICT DO NOTHING")
            .bind(Uuid::new_v4()).bind(id).bind(&raw_id).bind(kind.api_name()).bind(attempt as i32)
            .bind(json!({"query_id":id,"kind":kind.api_name(),"attempt":attempt,"retry_after_seconds":response.retry_after_seconds}))
            .bind(format!("{}:{raw_id}",response.received_at)).bind(serde_json::from_slice::<Value>(&response.body).ok())
            .bind(&response.body).bind(i32::from(response.status)).bind(json!({"endpoint":self.network.endpoint(),"method":"POST","body":body})).bind(&hash).bind(received)
            .execute(&self.pool).await.map_err(db_error)?;
        if inserted.rows_affected() == 0 {
            let previous:Option<(String,i32,Value,DateTime<Utc>)>=sqlx::query_as("SELECT sha256,http_status,request,observed_at FROM trade_log.raw_logs WHERE source_event_id=$1 AND collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$2)").bind(&raw_id).bind(id).fetch_optional(&self.pool).await.map_err(db_error)?;
            let expected_request =
                json!({"endpoint":self.network.endpoint(),"method":"POST","body":body});
            if previous != Some((hash, i32::from(response.status), expected_request, received)) {
                return Err(QueryError::conflict());
            }
        }
        if let Some(m) = &self.mirror
            && let Err(e) = m.save_response(id, kind, attempt, response, request).await
        {
            tracing::warn!(query_id=id,code=%e.code,"evidence_mirror_failed");
        }
        Ok(raw_id)
    }
    async fn finish(
        &self,
        id: &str,
        result: Result<&QueryResult, &QueryError>,
    ) -> Result<(), QueryError> {
        let saved = match result {
            Ok(r) => self.commit_facts(id, r).await,
            Err(e) => self.fail(id, e).await,
        };
        if let Err(e) = &saved {
            let _ = self.fail(id, e).await;
        }
        if let Some(m) = &self.mirror {
            let outcome = if let Err(e) = &saved { Err(e) } else { result };
            if let Err(e) = m.finish(id, outcome).await {
                tracing::warn!(query_id=id,code=%e.code,"evidence_mirror_failed");
            }
        }
        saved
    }
    async fn persistence(&self, id: &str) -> Result<Option<Persistence>, QueryError> {
        let value: Value = sqlx::query_scalar(
            "SELECT result FROM trade_log.collection_jobs WHERE query_id=$1 AND status='COMPLETED'",
        )
        .bind(id)
        .fetch_one(&self.pool)
        .await
        .map_err(db_error)?;
        serde_json::from_value(value["persistence"].clone())
            .map(Some)
            .map_err(|_| QueryError::storage())
    }
}
impl Postgres {
    pub async fn fail(&self, id: &str, error: &QueryError) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.collection_jobs SET status='FAILED',error=$2,updated_at=now() WHERE query_id=$1 AND status='RUNNING'").bind(id).bind(json!(error)).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
}
