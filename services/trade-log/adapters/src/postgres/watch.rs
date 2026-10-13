use super::*;
use base64::{Engine, engine::general_purpose::URL_SAFE_NO_PAD};
use chrono::{DateTime, Utc};
use serde_json::{Value, json};
use sqlx::{Postgres as Pg, Row, Transaction};
use trade_log::{
    checkpoint::{CollectionStatusReader, Lease},
    watch::{AddAccount, BudgetConfig, ListParams, MonitoringConfig, invalid},
};
use uuid::Uuid;
fn not_found() -> QueryError {
    trade_log::collection::error(
        "WATCH_ACCOUNT_NOT_FOUND",
        "Monitoring account not found",
        false,
    )
}
fn conflict(code: &str, message: &str) -> QueryError {
    trade_log::collection::error(code, message, false)
}
impl Postgres {
    pub async fn monitoring_configure(
        &self,
        c: &MonitoringConfig,
        b: &BudgetConfig,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        sqlx::query("SELECT id FROM trade_log.monitoring_settings WHERE id=1 FOR UPDATE")
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let (active,total):(i64,i64)=sqlx::query_as("SELECT count(*) FILTER(WHERE desired_state='ACTIVE'),count(*) FROM trade_log.watch_accounts").fetch_one(&mut *tx).await.map_err(db_error)?;
        if active > i64::from(c.max_active_accounts) || total > i64::from(c.max_registered_accounts)
        {
            return Err(conflict(
                "WATCH_ACCOUNT_LIMIT_REACHED",
                "Configured limits are below existing account counts",
            ));
        }
        sqlx::query("UPDATE trade_log.monitoring_settings SET max_active=$1,max_registered=$2,http_budget=$3,ws_budget=$4,stale_seconds=$5,max_http_requests=$6 WHERE id=1").bind(c.max_active_accounts).bind(c.max_registered_accounts).bind(b.http_weight_per_minute).bind(b.websocket_connections_per_minute).bind(c.worker_stale_after_seconds).bind(c.max_concurrent_http_requests as i32).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn import_watch(&self, account: &str, start: &str) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let imported: bool = sqlx::query_scalar(
            "SELECT imported FROM trade_log.monitoring_settings WHERE id=1 FOR UPDATE",
        )
        .fetch_one(&mut *tx)
        .await
        .map_err(db_error)?;
        if !imported && !account.is_empty() {
            let mut a = AddAccount {
                account: account.into(),
                start_time: start.into(),
            };
            a.validate()?;
            sqlx::query("INSERT INTO trade_log.watch_accounts(id,account,account_key,network,start_time) VALUES($1,$2,$3,$4,$5) ON CONFLICT(account_key) DO NOTHING").bind(Uuid::new_v4()).bind(&a.account).bind(self.checkpoint_key(&a.account)).bind(self.network.name()).bind(DateTime::parse_from_rfc3339(&a.start_time).map_err(|_|invalid("Invalid start_time"))?).execute(&mut *tx).await.map_err(db_error)?;
        }
        sqlx::query("UPDATE trade_log.monitoring_settings SET imported=true WHERE id=1")
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn add_watch(&self, mut a: AddAccount) -> Result<Value, QueryError> {
        a.validate()?;
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let settings =
            sqlx::query("SELECT * FROM trade_log.monitoring_settings WHERE id=1 FOR UPDATE")
                .fetch_one(&mut *tx)
                .await
                .map_err(db_error)?;
        let key = self.checkpoint_key(&a.account);
        if let Some(id) = sqlx::query_scalar::<_, Uuid>(
            "SELECT id FROM trade_log.watch_accounts WHERE account_key=$1",
        )
        .bind(&key)
        .fetch_optional(&mut *tx)
        .await
        .map_err(db_error)?
        {
            return Err(conflict(
                "WATCH_ACCOUNT_EXISTS",
                &format!("Monitoring account already exists; existing_id={id}"),
            ));
        }
        let (active,total):(i64,i64)=sqlx::query_as("SELECT count(*) FILTER(WHERE desired_state='ACTIVE'),count(*) FROM trade_log.watch_accounts").fetch_one(&mut *tx).await.map_err(db_error)?;
        if active >= i64::from(settings.get::<i32, _>("max_active"))
            || total >= i64::from(settings.get::<i32, _>("max_registered"))
        {
            return Err(conflict(
                "WATCH_ACCOUNT_LIMIT_REACHED",
                "Monitoring capacity reached",
            ));
        }
        let id = Uuid::new_v4();
        sqlx::query("INSERT INTO trade_log.watch_accounts(id,account,account_key,network,start_time) VALUES($1,$2,$3,$4,$5)").bind(id).bind(a.account).bind(key).bind(self.network.name()).bind(DateTime::parse_from_rfc3339(&a.start_time).map_err(|_|invalid("Invalid start_time"))?).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        self.watch_detail(id).await
    }
    pub async fn set_watch(&self, id: Uuid, active: bool) -> Result<Value, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let max: i32 = sqlx::query_scalar(
            "SELECT max_active FROM trade_log.monitoring_settings WHERE id=1 FOR UPDATE",
        )
        .fetch_one(&mut *tx)
        .await
        .map_err(db_error)?;
        let lock_key: Option<String> =
            sqlx::query_scalar("SELECT account_key FROM trade_log.watch_accounts WHERE id=$1")
                .bind(id)
                .fetch_optional(&mut *tx)
                .await
                .map_err(db_error)?;
        let lock_key = lock_key.ok_or_else(not_found)?;
        sqlx::query("SELECT pg_advisory_xact_lock(hashtextextended($1,6))")
            .bind(&lock_key)
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let row = sqlx::query(
            "SELECT account_key,desired_state FROM trade_log.watch_accounts WHERE id=$1 FOR UPDATE",
        )
        .bind(id)
        .fetch_optional(&mut *tx)
        .await
        .map_err(db_error)?
        .ok_or_else(not_found)?;
        let key: String = row.get("account_key");
        let state = if active { "ACTIVE" } else { "PAUSED" };
        if row.get::<String, _>("desired_state") != state {
            if active {
                let n: i64 = sqlx::query_scalar(
                    "SELECT count(*) FROM trade_log.watch_accounts WHERE desired_state='ACTIVE'",
                )
                .fetch_one(&mut *tx)
                .await
                .map_err(db_error)?;
                if n >= i64::from(max) {
                    return Err(conflict(
                        "WATCH_ACCOUNT_LIMIT_REACHED",
                        "Active monitoring capacity reached",
                    ));
                }
            }
            sqlx::query("UPDATE trade_log.watch_accounts SET desired_state=$2,generation=generation+1,updated_at=clock_timestamp(),paused_at=CASE WHEN $2='PAUSED' THEN clock_timestamp() ELSE paused_at END,resumed_at=CASE WHEN $2='ACTIVE' THEN clock_timestamp() ELSE resumed_at END,last_error=NULL,last_error_at=NULL WHERE id=$1").bind(id).bind(state).execute(&mut *tx).await.map_err(db_error)?;
            // Fence old workers immediately, without deleting pending work or advancing a cursor.
            sqlx::query("UPDATE trade_log.collection_checkpoints SET lease_epoch=lease_epoch+1,lease_owner=NULL,lease_expires_at=NULL,status='STOPPED' WHERE partition_key=$1 AND source_id='official_http'").bind(&key).execute(&mut *tx).await.map_err(db_error)?;
        }
        tx.commit().await.map_err(db_error)?;
        self.watch_detail(id).await
    }
    pub async fn watch_rows(&self) -> Result<Vec<Value>, QueryError> {
        sqlx::query_scalar("SELECT to_jsonb(w) FROM trade_log.watch_accounts w WHERE network=$1 ORDER BY created_at,id").bind(self.network.name()).fetch_all(&self.pool).await.map_err(db_error)
    }
    pub async fn watch_detail(&self, id: Uuid) -> Result<Value, QueryError> {
        let mut value: Value = sqlx::query_scalar(
            "SELECT to_jsonb(w) FROM trade_log.watch_accounts w WHERE id=$1 AND network=$2",
        )
        .bind(id)
        .bind(self.network.name())
        .fetch_optional(&self.pool)
        .await
        .map_err(db_error)?
        .ok_or_else(not_found)?;
        let account = value["account"].as_str().ok_or_else(QueryError::storage)?;
        let collection = self.collection_status(account).await;
        let state = if value["desired_state"] == "PAUSED" {
            "PAUSED"
        } else if value["last_error"].is_string() {
            "DEGRADED"
        } else if let Ok(c) = &collection {
            let c = &c.items[0];
            let stale: i32 = sqlx::query_scalar(
                "SELECT stale_seconds FROM trade_log.monitoring_settings WHERE id=1",
            )
            .fetch_one(&self.pool)
            .await
            .map_err(db_error)?;
            let recent = c
                .heartbeat_at
                .as_ref()
                .and_then(|v| DateTime::parse_from_rfc3339(v).ok())
                .is_some_and(|v| {
                    (Utc::now() - v.with_timezone(&Utc)).num_seconds() <= i64::from(stale)
                });
            if !recent
                || c.status == "FAILED"
                || c.status == "RETRY_WAIT"
                || c.monitoring_status == "DEGRADED"
                || c.monitoring_status == "FAILED"
            {
                "DEGRADED"
            } else if c.status == "STOPPED" || c.monitoring_status == "STARTING" {
                "STARTING"
            } else {
                "RUNNING"
            }
        } else {
            "STARTING"
        };
        value["observed_state"] = json!(state);
        value["observed_at"] = json!(shared_types::now());
        value["collection"] = match collection {
            Ok(c) => json!(c.items.into_iter().next()),
            Err(e)
                if e.code == "DEPENDENCY_UNAVAILABLE"
                    && e.message == "Collection checkpoint not initialized" =>
            {
                Value::Null
            }
            Err(e) => return Err(e),
        };
        let key = value["account_key"]
            .as_str()
            .ok_or_else(QueryError::storage)?;
        let exists: bool = sqlx::query_scalar(
            "SELECT EXISTS(SELECT FROM trade_log.publishing_control WHERE account_key=$1)",
        )
        .bind(key)
        .fetch_one(&self.pool)
        .await
        .map_err(db_error)?;
        let key = key.to_owned();
        value["publishing"] = if exists {
            self.publishing_status(&key).await?
        } else {
            Value::Null
        };
        let backfill: Option<serde_json::Value> = sqlx::query_scalar("SELECT to_jsonb(b) FROM trade_log.spot_backfill_jobs b WHERE account_key=$1 ORDER BY created_at DESC LIMIT 1").bind(key).fetch_optional(&self.pool).await.map_err(db_error)?;
        value["spot_backfill"] = serde_json::json!(backfill);
        Ok(value)
    }
    pub async fn watch_list(&self, p: ListParams) -> Result<Value, QueryError> {
        let limit = p.limit.unwrap_or(20);
        if !(1..=100).contains(&limit)
            || p.state
                .as_ref()
                .is_some_and(|s| s != "ACTIVE" && s != "PAUSED")
        {
            return Err(invalid("Invalid list parameters"));
        }
        let cursor: Option<(DateTime<Utc>, Uuid)> = p
            .cursor
            .map(|s| {
                URL_SAFE_NO_PAD
                    .decode(s)
                    .ok()
                    .and_then(|v| serde_json::from_slice(&v).ok())
                    .ok_or_else(|| invalid("Invalid cursor"))
            })
            .transpose()?;
        let (at, id) = cursor
            .map(|(t, i)| (Some(t), Some(i)))
            .unwrap_or((None, None));
        let rows:Vec<(Uuid,DateTime<Utc>)>=sqlx::query_as("SELECT id,created_at FROM trade_log.watch_accounts WHERE network=$1 AND ($2::text IS NULL OR desired_state=$2) AND ($3::timestamptz IS NULL OR (created_at,id)>($3,$4)) ORDER BY created_at,id LIMIT $5").bind(self.network.name()).bind(&p.state).bind(at).bind(id).bind(limit as i64+1).fetch_all(&self.pool).await.map_err(db_error)?;
        let next = if rows.len() > limit {
            let (id, at) = rows[limit - 1];
            Some(
                URL_SAFE_NO_PAD
                    .encode(serde_json::to_vec(&(at, id)).map_err(|_| QueryError::storage())?),
            )
        } else {
            None
        };
        let mut items = vec![];
        for (id, _) in rows.into_iter().take(limit) {
            items.push(self.watch_detail(id).await?)
        }
        let (total,active):(i64,i64)=sqlx::query_as("SELECT count(*),count(*) FILTER(WHERE desired_state='ACTIVE') FROM trade_log.watch_accounts WHERE network=$1").bind(self.network.name()).fetch_one(&self.pool).await.map_err(db_error)?;
        let max: i32 =
            sqlx::query_scalar("SELECT max_active FROM trade_log.monitoring_settings WHERE id=1")
                .fetch_one(&self.pool)
                .await
                .map_err(db_error)?;
        Ok(
            json!({"items":items,"next_cursor":next,"total":total,"active_count":active,"max_active_accounts":max}),
        )
    }
    pub async fn lock_watch(
        &self,
        tx: &mut Transaction<'_, Pg>,
        l: &Lease,
    ) -> Result<(), QueryError> {
        sqlx::query("SELECT pg_advisory_xact_lock(hashtextextended($1,6))")
            .bind(&l.key)
            .execute(&mut **tx)
            .await
            .map_err(db_error)?;
        let row=sqlx::query("SELECT generation,desired_state FROM trade_log.watch_accounts WHERE account_key=$1 FOR UPDATE").bind(&l.key).fetch_optional(&mut **tx).await.map_err(db_error)?;
        if let Some(row) = row
            && (row.get::<String, _>("desired_state") != "ACTIVE"
                || l.generation != Some(row.get("generation")))
        {
            return Err(trade_log::collection::error(
                "LEASE_LOST",
                "Monitoring generation lost",
                true,
            ));
        }
        Ok(())
    }
    pub async fn watch_worker_error(
        &self,
        key: &str,
        generation: i64,
        error: Option<&str>,
    ) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.watch_accounts SET last_error=$3,last_error_at=CASE WHEN $3::text IS NULL THEN NULL ELSE clock_timestamp() END WHERE account_key=$1 AND generation=$2").bind(key).bind(generation).bind(error).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
}

impl Postgres {
    pub async fn fence_worker(&self, key: &str, owner: uuid::Uuid) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.collection_checkpoints SET lease_owner=NULL,lease_expires_at=NULL,lease_epoch=lease_epoch+1,status='STOPPED' WHERE partition_key=$1 AND source_id='official_http' AND lease_owner=$2").bind(key).bind(owner).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
}
