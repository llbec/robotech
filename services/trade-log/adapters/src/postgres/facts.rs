use super::*;
use account_facts::AccountFact;
use chrono::{DateTime, Utc};
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use sqlx::Row;
use trade_log::{persistence::Persistence, query::QueryResult};
use uuid::Uuid;
pub fn content_hash(fact: &AccountFact) -> Result<String, QueryError> {
    let mut value = json!(fact);
    value
        .as_object_mut()
        .ok_or_else(QueryError::storage)?
        .remove("raw_log_id");
    value["payload"]["extension"]
        .as_object_mut()
        .ok_or_else(QueryError::storage)?
        .remove("source_indices");
    Ok(format!(
        "{:x}",
        Sha256::digest(serde_json::to_vec(&value).map_err(|_| QueryError::storage())?)
    ))
}
impl Postgres {
    pub async fn commit_facts(&self, id: &str, result: &QueryResult) -> Result<(), QueryError> {
        self.commit_facts_at(id, result, None).await
    }
    pub async fn commit_facts_at(
        &self,
        id: &str,
        result: &QueryResult,
        finished: Option<DateTime<Utc>>,
    ) -> Result<(), QueryError> {
        let created = if finished.is_some() {
            Some(
                DateTime::parse_from_rfc3339(&result.queried_at)
                    .map_err(|_| QueryError::storage())?
                    .with_timezone(&Utc),
            )
        } else {
            None
        };
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let mut seq: i64 = sqlx::query_scalar(
            "SELECT committed_seq FROM trade_log.ingestion_state WHERE id=1 FOR UPDATE",
        )
        .fetch_one(&mut *tx)
        .await
        .map_err(db_error)?;
        let status: String = sqlx::query_scalar(
            "SELECT status FROM trade_log.collection_jobs WHERE query_id=$1 FOR UPDATE",
        )
        .bind(id)
        .fetch_one(&mut *tx)
        .await
        .map_err(db_error)?;
        if status == "COMPLETED" {
            return Ok(());
        }
        if status != "RUNNING" {
            return Err(QueryError::conflict());
        }
        let mut inserted = 0;
        for fact in &result.trades {
            let hash = content_hash(fact)?;
            let previous:Option<String>=sqlx::query_scalar("SELECT content_hash FROM trade_log.account_fact_versions WHERE fact_id=$1 AND revision=1").bind(&fact.fact_id).fetch_optional(&mut *tx).await.map_err(db_error)?;
            let raw:Uuid=sqlx::query_scalar("SELECT r.id FROM trade_log.raw_logs r JOIN trade_log.collection_jobs j ON j.id=r.collection_job_id WHERE r.source_event_id=$1 AND j.query_id=$2").bind(&fact.raw_log_id).bind(id).fetch_one(&mut *tx).await.map_err(db_error)?;
            if let Some(previous) = previous {
                if previous != hash {
                    return Err(QueryError::conflict());
                }
            } else {
                let at = DateTime::parse_from_rfc3339(&fact.occurred_at)
                    .map_err(|_| QueryError::storage())?
                    .with_timezone(&Utc);
                sqlx::query("INSERT INTO trade_log.account_fact_versions(fact_id,revision,event_id,account_key,fact_type,ordering_key,sub_index,change_type,confirmation_status,occurred_at,payload,raw_log_id,parser_version,content_hash) VALUES($1,1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,'hyperliquid-v1',$12)")
                    .bind(&fact.fact_id).bind(format!("{}:1",fact.fact_id)).bind(&fact.account_key).bind(&fact.fact_type).bind(&fact.ordering_key).bind(fact.sub_index as i32).bind(&fact.change_type).bind(&fact.confirmation_status).bind(at).bind(json!(fact)).bind(raw).bind(&hash).execute(&mut *tx).await.map_err(db_error)?;
                seq = seq.checked_add(1).ok_or_else(QueryError::storage)?;
                sqlx::query("INSERT INTO trade_log.account_facts_current(fact_id,current_revision,account_key,fact_type,ordering_key,sub_index,occurred_at,source_tid,ingest_seq) VALUES($1,1,$2,$3,$4,$5,$6,$7::text::numeric,$8)")
                    .bind(&fact.fact_id).bind(&fact.account_key).bind(&fact.fact_type).bind(&fact.ordering_key).bind(fact.sub_index as i32).bind(at).bind(&fact.source_ref).bind(seq).execute(&mut *tx).await.map_err(db_error)?;
                inserted += 1;
            }
            for index in fact.payload.extension["source_indices"]
                .as_array()
                .ok_or_else(QueryError::storage)?
            {
                let index = index
                    .as_u64()
                    .and_then(|v| i32::try_from(v).ok())
                    .ok_or_else(QueryError::storage)?;
                sqlx::query("INSERT INTO trade_log.fact_observations VALUES($1,1,$2,$3) ON CONFLICT DO NOTHING").bind(&fact.fact_id).bind(raw).bind(index).execute(&mut *tx).await.map_err(db_error)?;
            }
        }
        sqlx::query("UPDATE trade_log.ingestion_state SET committed_seq=$1 WHERE id=1")
            .bind(seq)
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let mut full = result.clone();
        full.display_truncated = false;
        full.counts.returned_records = full.trades.len();
        full.persistence = Some(Persistence {
            status: "COMMITTED".into(),
            inserted_records: inserted,
            existing_records: full.trades.len() - inserted,
        });
        sqlx::query("UPDATE trade_log.raw_logs SET parse_status='PARSED' WHERE collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$1) AND http_status BETWEEN 200 AND 299").bind(id).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("UPDATE trade_log.collection_jobs SET status='COMPLETED',result=$2,error=NULL,updated_at=COALESCE($3,now()),created_at=COALESCE($4,created_at) WHERE query_id=$1").bind(id).bind(json!(full)).bind(finished).bind(created).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn job_result(&self, id: &str) -> Result<(String, String, Value), QueryError> {
        let row = sqlx::query(
            "SELECT status,network,result FROM trade_log.collection_jobs WHERE query_id=$1",
        )
        .bind(id)
        .fetch_one(&self.pool)
        .await
        .map_err(db_error)?;
        Ok((
            row.get("status"),
            row.get("network"),
            row.try_get::<Option<Value>, _>("result")
                .map_err(db_error)?
                .unwrap_or(Value::Null),
        ))
    }
}

#[async_trait::async_trait]
impl trade_log::persistence::FactStore for Postgres {
    async fn persist(&self, id: &str, complete: &QueryResult) -> Result<Persistence, QueryError> {
        self.commit_facts(id, complete).await?;
        trade_log::raw_log::RawEvidenceStore::persistence(self, id)
            .await?
            .ok_or_else(QueryError::storage)
    }
}
