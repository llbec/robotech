use super::*;
use async_trait::async_trait;
use chrono::{DateTime, SecondsFormat, Utc};
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use sqlx::{Postgres as Pg, Row, Transaction};
use trade_log::{
    acquisition::SourceResponse,
    checkpoint::*,
    collection::{CollectionConfig, time},
    parsing::{ParseContext, ProtocolParser},
    query::QueryResult,
};
use uuid::Uuid;

pub use trade_log::checkpoint::CollectionCommit;
impl Postgres {
    pub async fn collection_retry_after(&self, id: &str) -> Result<u64, QueryError> {
        let value:Option<serde_json::Value>=sqlx::query_scalar("SELECT chain_position FROM trade_log.raw_logs WHERE collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$1) ORDER BY page_no DESC LIMIT 1").bind(id).fetch_optional(&self.pool).await.map_err(db_error)?;
        Ok(value
            .and_then(|v| v["retry_after_seconds"].as_u64())
            .unwrap_or(0))
    }
    pub async fn checkpoint_exists(&self, account: &str) -> Result<bool, QueryError> {
        sqlx::query_scalar(
            "SELECT EXISTS(SELECT FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http')",
        )
        .bind(self.checkpoint_key(account))
        .fetch_one(&self.pool)
        .await
        .map_err(db_error)
    }
    pub fn checkpoint_key(&self, account: &str) -> String {
        format!(
            "hyperliquid:{}:hyperliquid:{}",
            self.network.name(),
            account
        )
    }
    pub async fn acquire_collection(
        &self,
        c: &CollectionConfig,
        owner: Uuid,
    ) -> Result<Option<Lease>, QueryError> {
        let key = self.checkpoint_key(&c.account);
        let initial = trade_log::collection::start_ms(&c.start_time)?;
        sqlx::query("INSERT INTO trade_log.collection_checkpoints(chain_id,partition_key,position) VALUES($1,$2,$3) ON CONFLICT DO NOTHING")
   .bind(format!("hyperliquid:{}",self.network.name())).bind(&key).bind(json!({"initial_start_ms":initial,"scanned_through_ms":null,"last_trade_at_ms":null})).execute(&self.pool).await.map_err(db_error)?;
        let epoch:Option<i64>=sqlx::query_scalar("UPDATE trade_log.collection_checkpoints SET lease_epoch=CASE WHEN lease_owner=$2 AND lease_expires_at>now() THEN lease_epoch ELSE lease_epoch+1 END,status=CASE WHEN lease_owner=$2 AND lease_expires_at>now() THEN status ELSE 'STARTING' END,lease_owner=$2,lease_expires_at=now()+make_interval(secs=>$3),heartbeat_at=now(),updated_at=now() WHERE partition_key=$1 AND source_id='official_http' AND (lease_owner=$2 OR lease_expires_at IS NULL OR lease_expires_at<=now()) RETURNING lease_epoch")
   .bind(&key).bind(owner).bind(c.lease_seconds as f64).fetch_optional(&self.pool).await.map_err(db_error)?;
        Ok(epoch.map(|epoch| Lease {
            key,
            owner: owner.to_string(),
            epoch,
        }))
    }
    pub async fn renew_collection(&self, l: &Lease, seconds: u64) -> Result<(), QueryError> {
        let changed=sqlx::query("UPDATE trade_log.collection_checkpoints SET heartbeat_at=now(),lease_expires_at=now()+make_interval(secs=>$4),updated_at=now() WHERE partition_key=$1 AND source_id='official_http' AND lease_owner=$2 AND lease_epoch=$3 AND lease_expires_at>now()")
   .bind(&l.key).bind(owner(l)?).bind(l.epoch).bind(seconds as f64).execute(&self.pool).await.map_err(db_error)?;
        if changed.rows_affected() != 1 {
            return Err(lost());
        }
        Ok(())
    }
    pub async fn lock_collection(
        &self,
        tx: &mut Transaction<'_, Pg>,
        l: &Lease,
    ) -> Result<Value, QueryError> {
        let position:Option<Value>=sqlx::query_scalar("SELECT position FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http' AND lease_owner=$2 AND lease_epoch=$3 AND lease_expires_at>now() FOR UPDATE")
   .bind(&l.key).bind(owner(l)?).bind(l.epoch).fetch_optional(&mut **tx).await.map_err(db_error)?;
        position.ok_or_else(lost)
    }
    pub async fn collection_work(
        &self,
        l: &Lease,
        c: &CollectionConfig,
        now_ms: i64,
    ) -> Result<Option<PendingWork>, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let pos = self.lock_collection(&mut tx, l).await?;
        let row=sqlx::query("SELECT pending_work,status,next_run_at FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).fetch_one(&mut *tx).await.map_err(db_error)?;
        if row.get::<String, _>("status") == "FAILED" {
            return Ok(None);
        }
        let next: Option<DateTime<Utc>> = row.get("next_run_at");
        if next.is_some_and(|t| t.timestamp_millis() > now_ms) {
            return Ok(None);
        }
        let pending: Option<Value> = row.get("pending_work");
        let work = if let Some(value) = pending {
            let work: PendingWork =
                serde_json::from_value(value).map_err(|_| QueryError::storage())?;
            sqlx::query("UPDATE trade_log.collection_jobs SET status='RUNNING',error=NULL,lease_epoch=$2,updated_at=now() WHERE query_id=$1 AND job_origin='COLLECTOR' AND status<>'COMPLETED'").bind(&work.query_id).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
            work
        } else {
            let range = trade_log::collection::scan_range(
                pos["initial_start_ms"]
                    .as_i64()
                    .ok_or_else(QueryError::storage)?,
                pos["scanned_through_ms"].as_i64(),
                now_ms,
                c,
            )?;
            let Some(range) = range else {
                return Ok(None);
            };
            let id = format!("query_{}", Uuid::new_v4().simple());
            sqlx::query("INSERT INTO trade_log.collection_jobs(id,query_id,chain_id,account,network,status,request,trace_id,job_origin,checkpoint_key,lease_epoch) VALUES($1,$2,$3,$4,$5,'RUNNING',$6,$7,'COLLECTOR',$8,$9)")
    .bind(Uuid::new_v4()).bind(&id).bind(format!("hyperliquid:{}",self.network.name())).bind(&c.account).bind(self.network.name()).bind(json!({"account":c.account,"start_ms":range.start_ms,"end_ms":range.end_ms})).bind(format!("trace_{}",Uuid::new_v4().simple())).bind(&l.key).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
            PendingWork {
                query_id: id,
                remaining: vec![range.clone()],
                range,
                pages: vec![],
                meta: None,
                spot_meta: None,
            }
        };
        sqlx::query("UPDATE trade_log.collection_checkpoints SET status='RUNNING',pending_work=$2,last_query_id=$3,last_attempt_at=now(),next_run_at=NULL,updated_at=now() WHERE partition_key=$1 AND source_id='official_http'")
   .bind(&l.key).bind(json!(work)).bind(&work.query_id).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        Ok(Some(work))
    }
    pub async fn save_work(&self, l: &Lease, work: &PendingWork) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("UPDATE trade_log.collection_checkpoints SET pending_work=$2,updated_at=now() WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(json!(work)).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn collection_failure(
        &self,
        l: &Lease,
        e: &QueryError,
        delay: u64,
        budget: bool,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        let state = if e.retryable || budget {
            "RETRY_WAIT"
        } else {
            "FAILED"
        };
        let diagnostic =
            json!({"code":e.code,"message":e.message,"occurred_at":shared_types::now()});
        sqlx::query("UPDATE trade_log.collection_checkpoints SET status=$2,pending_work=CASE WHEN $5 THEN jsonb_set(pending_work,'{budget_exhausted}','true'::jsonb) ELSE pending_work END,last_error=CASE WHEN $5 THEN last_error ELSE $3 END,consecutive_failures=consecutive_failures+CASE WHEN $5 THEN 0 ELSE 1 END,next_run_at=CASE WHEN $2='FAILED' THEN NULL ELSE now()+make_interval(secs=>$4) END,updated_at=now() WHERE partition_key=$1 AND source_id='official_http'")
   .bind(&l.key).bind(state).bind(diagnostic).bind(delay as f64).bind(budget).execute(&mut *tx).await.map_err(db_error)?;
        if !budget {
            sqlx::query("UPDATE trade_log.collection_jobs SET status='FAILED',error=$2,updated_at=now() WHERE query_id=(SELECT last_query_id FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http') AND status='RUNNING'").bind(&l.key).bind(json!(e)).execute(&mut *tx).await.map_err(db_error)?;
        }
        tx.commit().await.map_err(db_error)
    }
    pub async fn release_collection(&self, l: &Lease) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.collection_checkpoints SET status='STOPPED',lease_owner=NULL,lease_expires_at=NULL,updated_at=now() WHERE partition_key=$1 AND source_id='official_http' AND lease_owner=$2 AND lease_epoch=$3")
   .bind(&l.key).bind(owner(l)?).bind(l.epoch).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
    pub async fn save_collection_response(
        &self,
        l: &Lease,
        work: &PendingWork,
        kind: &str,
        range: Option<&ScanRange>,
        response: &SourceResponse,
    ) -> Result<String, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        let page:i32=sqlx::query_scalar("SELECT COALESCE(max(page_no),0)+1 FROM trade_log.raw_logs WHERE collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$1)").bind(&work.query_id).fetch_one(&mut *tx).await.map_err(db_error)?;
        let raw = format!("raw_{}_{}_p{}_1", work.query_id, kind, page);
        let mut body = json!({"type":kind});
        if let Some(r) = range {
            body["user"] = json!(self.account_for_key(&l.key)?);
            body["startTime"] = json!(r.start_ms);
            body["endTime"] = json!(r.end_ms - 1);
            body["aggregateByTime"] = json!(false);
        }
        let received = DateTime::parse_from_rfc3339(&response.received_at)
            .map_err(|_| QueryError::storage())?
            .with_timezone(&Utc);
        sqlx::query("INSERT INTO trade_log.raw_logs(id,collection_job_id,chain_id,source_event_id,kind,page_no,attempt,chain_position,ordering_key,payload,body,http_status,request,sha256,observed_at) SELECT $1,id,chain_id,$3,$4,$5,1,$6,$7,$8,$9,$10,$11,$12,$13 FROM trade_log.collection_jobs WHERE query_id=$2 AND job_origin='COLLECTOR' AND lease_epoch=$14")
   .bind(Uuid::new_v4()).bind(&work.query_id).bind(&raw).bind(kind).bind(page).bind(json!({"query_id":work.query_id,"page_no":page,"attempt":1,"retry_after_seconds":response.retry_after_seconds}))
   .bind(format!("{}:{raw}",response.received_at)).bind(serde_json::from_slice::<Value>(&response.body).ok()).bind(&response.body).bind(i32::from(response.status)).bind(json!({"endpoint":self.network.endpoint(),"method":"POST","body":body})).bind(format!("{:x}",Sha256::digest(&response.body))).bind(received).bind(l.epoch)
   .execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        if let Some(m) = &self.mirror
            && let Err(e) = m
                .collector_response(
                    &work.query_id,
                    &raw,
                    &json!({"endpoint":self.network.endpoint(),"method":"POST","body":body}),
                    response,
                )
                .await
        {
            tracing::warn!(query_id=%work.query_id,code=%e.code,"evidence_mirror_failed");
        }
        Ok(raw)
    }
    fn account_for_key<'a>(&self, key: &'a str) -> Result<&'a str, QueryError> {
        key.rsplit(':').next().ok_or_else(QueryError::storage)
    }
    pub async fn collection_body(&self, id: &str, raw: &str) -> Result<Vec<u8>, QueryError> {
        let row=sqlx::query("SELECT body,sha256 FROM trade_log.raw_logs WHERE source_event_id=$1 AND collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$2) AND http_status BETWEEN 200 AND 299").bind(raw).bind(id).fetch_one(&self.pool).await.map_err(db_error)?;
        let body: Vec<u8> = row.get("body");
        if format!("{:x}", Sha256::digest(&body)) != row.get::<String, _>("sha256") {
            return Err(QueryError::incomplete("Stored response checksum mismatch"));
        }
        Ok(body)
    }
    pub async fn collection_result(
        &self,
        id: &str,
        account: &str,
        ids: &[String],
        meta_id: &str,
        spot_id: &str,
        parser: &dyn ProtocolParser,
    ) -> Result<(QueryResult, Vec<account_facts::AccountFact>), QueryError> {
        let meta = self.collection_body(id, meta_id).await?;
        let spot = self.collection_body(id, spot_id).await?;
        let mut parts = vec![];
        for raw in ids {
            let body = self.collection_body(id, raw).await?;
            parts.push(parser.parse(ParseContext {
                network: &self.network,
                account,
                raw_log_id: raw,
                fills: &body,
                meta: &meta,
                spot_meta: &spot,
            })?);
        }
        let (parsed, observations) = trade_log::collection::merge(parts)?;
        let mut result = trade_log::normalization::result(
            id,
            account,
            &self.network,
            &shared_types::now(),
            usize::MAX,
            parsed,
        )?;
        result.query_scope = "AUTOMATIC_TIME_RANGE".into();
        result.coverage = "SOURCE_HISTORY_NOT_VERIFIED".into();
        result.warnings.retain(|w| w != "SOURCE_RECORD_LIMIT");
        result.warnings.push("SOURCE_HISTORY_NOT_VERIFIED".into());
        Ok((result, observations))
    }
}
fn owner(l: &Lease) -> Result<Uuid, QueryError> {
    Uuid::parse_str(&l.owner).map_err(|_| QueryError::storage())
}
fn lost() -> QueryError {
    trade_log::collection::error("LEASE_LOST", "Collection lease lost", true)
}
fn timestamp(value: Option<DateTime<Utc>>) -> Option<String> {
    value.map(|v| v.to_rfc3339_opts(SecondsFormat::Millis, true))
}
#[async_trait]
impl CollectionStatusReader for Postgres {
    async fn collection_status(&self, account: &str) -> Result<CollectionList, QueryError> {
        let key = self.checkpoint_key(account);
        let row=sqlx::query("SELECT *,lease_expires_at<=now() AS expired FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http'").bind(&key).fetch_optional(&self.pool).await.map_err(db_error)?;
        let Some(row) = row else {
            return Err(QueryError::unavailable(
                "Collection checkpoint not initialized",
            ));
        };
        let pos: Value = row.get("position");
        let mut state: String = row.get("status");
        let mut warnings = vec![
            "SOURCE_HISTORY_NOT_VERIFIED".into(),
            "LATE_DATA_OUTSIDE_OVERLAP_NOT_VERIFIED".into(),
        ];
        if row.get::<Option<bool>, _>("expired") == Some(true) && state != "STOPPED" {
            state = "FAILED".into();
            warnings.push("COLLECTOR_HEARTBEAT_EXPIRED".into());
        }
        let pending: Option<Value> = row.get("pending_work");
        if state == "RETRY_WAIT"
            && pending
                .as_ref()
                .is_some_and(|p| p["budget_exhausted"] == true)
        {
            warnings.push("ROUND_BUDGET_EXHAUSTED".into());
        }
        let status=CollectionStatus{account:account.into(),account_key:key,network:self.network.name().into(),status:state,coverage:"SOURCE_HISTORY_NOT_VERIFIED".into(),initial_start_time:time(pos["initial_start_ms"].as_i64().ok_or_else(QueryError::storage)?),scanned_through:pos["scanned_through_ms"].as_i64().map(time),last_trade_at:pos["last_trade_at_ms"].as_i64().map(time),last_attempt_at:timestamp(row.get("last_attempt_at")),last_success_at:timestamp(row.get("last_success_at")),consecutive_failures:row.get("consecutive_failures"),next_run_at:timestamp(row.get("next_run_at")),pending_range:pending.map(|p|json!({"start_time":time(p["range"]["start_ms"].as_i64().unwrap_or(0)),"end_time":time(p["range"]["end_ms"].as_i64().unwrap_or(0))})),last_query_id:row.get("last_query_id"),last_success_query_id:row.get("last_success_query_id"),last_error:row.get("last_error"),heartbeat_at:timestamp(row.get("heartbeat_at")),lease_expires_at:timestamp(row.get("lease_expires_at")),warnings};
        Ok(CollectionList {
            items: vec![status],
        })
    }
}
