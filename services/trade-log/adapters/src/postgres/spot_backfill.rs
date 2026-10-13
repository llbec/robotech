use super::*;
use serde_json::{Value, json};
use sqlx::Row;
use uuid::Uuid;
impl Postgres {
    pub async fn submit_spot_backfill(
        &self,
        account: &str,
        start: i64,
        end: i64,
    ) -> Result<Value, QueryError> {
        if start < 0 || end <= start || end > chrono::Utc::now().timestamp_millis() {
            return Err(QueryError::validation("Invalid backfill time range"));
        }
        let account = shared_types::account(account).map_err(QueryError::validation)?;
        let key = self.checkpoint_key(&account);
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        sqlx::query("SELECT pg_advisory_xact_lock(hashtextextended($1,6))")
            .bind(&key)
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let registered: bool = sqlx::query_scalar(
            "SELECT EXISTS(SELECT FROM trade_log.watch_accounts WHERE account_key=$1)",
        )
        .bind(&key)
        .fetch_one(&mut *tx)
        .await
        .map_err(db_error)?;
        if !registered {
            return Err(trade_log::collection::error(
                "WATCH_ACCOUNT_NOT_FOUND",
                "Register the account before backfill",
                false,
            ));
        }
        let same:Option<Value>=sqlx::query_scalar("SELECT to_jsonb(b) FROM trade_log.spot_backfill_jobs b WHERE account_key=$1 AND start_ms=$2 AND end_ms=$3").bind(&key).bind(start).bind(end).fetch_optional(&mut *tx).await.map_err(db_error)?;
        if let Some(mut v) = same {
            if v["status"] == "BLOCKED" {
                v=sqlx::query_scalar("UPDATE trade_log.spot_backfill_jobs SET status='PENDING',retry_at=NULL,updated_at=clock_timestamp() WHERE id=$1::text::uuid RETURNING to_jsonb(spot_backfill_jobs)").bind(v["id"].as_str()).fetch_one(&mut *tx).await.map_err(db_error)?;
            }
            tx.commit().await.map_err(db_error)?;
            return Ok(v);
        }
        let busy:bool=sqlx::query_scalar("SELECT EXISTS(SELECT FROM trade_log.spot_backfill_jobs WHERE account_key=$1 AND status<>'COMPLETED')").bind(&key).fetch_one(&mut *tx).await.map_err(db_error)?;
        if busy {
            return Err(QueryError::conflict());
        }
        let id = Uuid::new_v4();
        let v:Value=sqlx::query_scalar("INSERT INTO trade_log.spot_backfill_jobs(id,account_key,start_ms,end_ms,scanned_ms) VALUES($1,$2,$3,$4,$3) RETURNING to_jsonb(spot_backfill_jobs)").bind(id).bind(&key).bind(start).bind(end).fetch_one(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        Ok(v)
    }
    pub async fn spot_backfill_status(&self, id: Uuid) -> Result<Value, QueryError> {
        sqlx::query_scalar("SELECT to_jsonb(b) FROM trade_log.spot_backfill_jobs b WHERE id=$1 AND account_key LIKE $2").bind(id).bind(format!("hyperliquid:{}:%",self.network.name())).fetch_optional(&self.pool).await.map_err(db_error)?.ok_or_else(||QueryError::incomplete("Unknown backfill job"))
    }
    pub async fn claim_spot_backfill(&self, owner: Uuid) -> Result<Option<Value>, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let row=sqlx::query("SELECT b.id,w.generation FROM trade_log.spot_backfill_jobs b JOIN trade_log.watch_accounts w USING(account_key) WHERE w.network=$1 AND w.desired_state='ACTIVE' AND b.status IN ('PENDING','RUNNING','RETRY_WAIT') AND (b.lease_until IS NULL OR b.lease_until<clock_timestamp()) AND (b.retry_at IS NULL OR b.retry_at<=clock_timestamp()) ORDER BY b.updated_at,b.id LIMIT 1 FOR UPDATE OF b SKIP LOCKED").bind(self.network.name()).fetch_optional(&mut *tx).await.map_err(db_error)?;
        let Some(row) = row else { return Ok(None) };
        let id: Uuid = row.get("id");
        let mut v:Value=sqlx::query_scalar("UPDATE trade_log.spot_backfill_jobs SET status='RUNNING',lease_owner=$2,lease_until=clock_timestamp()+interval '120 seconds',updated_at=clock_timestamp() WHERE id=$1 RETURNING to_jsonb(spot_backfill_jobs)").bind(id).bind(owner).fetch_one(&mut *tx).await.map_err(db_error)?;
        v["generation"] = json!(row.get::<i64, _>("generation"));
        tx.commit().await.map_err(db_error)?;
        Ok(Some(v))
    }
    pub async fn finish_spot_batch(
        &self,
        job: &Value,
        scanned: i64,
        cursor: Option<Uuid>,
        evidence_done: bool,
    ) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.spot_backfill_jobs SET scanned_ms=$3,evidence_cursor=COALESCE($4,evidence_cursor),evidence_complete=$5,status=CASE WHEN $5 AND $3>=end_ms THEN 'COMPLETED' ELSE 'PENDING' END,lease_owner=NULL,lease_until=NULL,last_error=NULL,attempts=0,updated_at=clock_timestamp() WHERE id=$1::text::uuid AND lease_owner=$2::text::uuid")
            .bind(job["id"].as_str()).bind(job["lease_owner"].as_str()).bind(scanned).bind(cursor).bind(evidence_done).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
    pub async fn fail_spot_batch(&self, job: &Value, error: &QueryError) {
        let saturated = error.code == "SOURCE_WINDOW_SATURATED";
        let window = job["window_ms"].as_i64().unwrap_or(1);
        let blocked = (!error.retryable && !saturated) || (saturated && window <= 1);
        let _=sqlx::query("UPDATE trade_log.spot_backfill_jobs SET status=$3,last_error=$4,attempts=attempts+1,retry_at=GREATEST(COALESCE(retry_at,clock_timestamp()),clock_timestamp()+interval '5 seconds'*LEAST(60,power(2,LEAST(attempts,6))::int)),window_ms=$5,lease_owner=NULL,lease_until=NULL,updated_at=clock_timestamp() WHERE id=$1::text::uuid AND lease_owner=$2::text::uuid")
            .bind(job["id"].as_str()).bind(job["lease_owner"].as_str()).bind(if blocked {"BLOCKED"} else {"RETRY_WAIT"}).bind(&error.code).bind(if saturated {(window/2).max(1)} else {window}).execute(&self.pool).await;
    }
}
