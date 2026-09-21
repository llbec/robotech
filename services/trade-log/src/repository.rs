use crate::{
    domain::{HyperliquidFill, MonitoredAddress, TradeFact},
    market_metadata::MetadataStore,
    normalization,
};
use anyhow::{bail, Context, Result};
use chrono::{DateTime, Utc};
use serde::Serialize;
use serde_json::Value;
use sqlx::{postgres::PgPoolOptions, types::Json, PgPool};
use uuid::Uuid;

#[derive(Clone)]
pub struct Repository {
    pool: PgPool,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct RawTradeView {
    pub id: Uuid,
    pub source_event_id: String,
    pub transport: String,
    pub event_time: DateTime<Utc>,
    pub observed_at: DateTime<Utc>,
    pub is_snapshot: bool,
    pub parse_status: String,
    pub parse_error: Option<String>,
    pub payload: Json<Value>,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct TradeFactView {
    pub fact_id: String,
    pub event_id: String,
    pub revision: i32,
    pub account_key: String,
    pub account: String,
    pub occurred_at: DateTime<Utc>,
    pub observed_at: DateTime<Utc>,
    pub ordering_key: String,
    pub source: String,
    pub source_ref: String,
    pub raw_log_id: Uuid,
    pub schema_version: i32,
    pub payload: Json<Value>,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct DeliveryView {
    pub event_id: String,
    pub fact_id: String,
    pub status: String,
    pub attempts: i32,
    pub available_at: DateTime<Utc>,
    pub last_error: Option<String>,
    pub delivered_at: Option<DateTime<Utc>>,
}

pub struct IngestResult {
    pub inserted: bool,
    pub fact: TradeFact,
}

impl Repository {
    pub async fn connect(database_url: &str) -> Result<Self> {
        let pool = PgPoolOptions::new()
            .max_connections(10)
            .connect(database_url)
            .await
            .context("connect PostgreSQL")?;
        Ok(Self { pool })
    }

    pub fn pool(&self) -> &PgPool {
        &self.pool
    }

    pub async fn migrate(&self) -> Result<()> {
        sqlx::migrate!("./migrations").run(&self.pool).await?;
        Ok(())
    }

    pub async fn ready(&self) -> bool {
        sqlx::query_scalar::<_, i32>("SELECT 1")
            .fetch_one(&self.pool)
            .await
            .is_ok()
    }

    pub async fn create_monitor(
        &self,
        address: &str,
        start: DateTime<Utc>,
        maximum: i64,
    ) -> Result<(MonitoredAddress, bool)> {
        let mut tx = self.pool.begin().await?;
        sqlx::query("SELECT pg_advisory_xact_lock(753000)")
            .execute(&mut *tx)
            .await?;
        if let Some(existing) = sqlx::query_as::<_, MonitoredAddress>(
            "SELECT * FROM monitored_addresses WHERE address=$1",
        )
        .bind(address)
        .fetch_optional(&mut *tx)
        .await?
        {
            tx.commit().await?;
            return Ok((existing, false));
        }
        let count: i64 = sqlx::query_scalar("SELECT count(*) FROM monitored_addresses")
            .fetch_one(&mut *tx)
            .await?;
        if count >= maximum {
            bail!("MONITOR_LIMIT_REACHED");
        }
        let monitor = sqlx::query_as::<_, MonitoredAddress>(
            "INSERT INTO monitored_addresses(address,status,requested_start) VALUES($1,'PENDING',$2) RETURNING *",
        )
        .bind(address)
        .bind(start)
        .fetch_one(&mut *tx)
        .await?;
        tx.commit().await?;
        Ok((monitor, true))
    }

    pub async fn list_monitors(&self) -> Result<Vec<MonitoredAddress>> {
        Ok(sqlx::query_as::<_, MonitoredAddress>(
            "SELECT * FROM monitored_addresses ORDER BY created_at, address",
        )
        .fetch_all(&self.pool)
        .await?)
    }

    pub async fn get_monitor(&self, address: &str) -> Result<Option<MonitoredAddress>> {
        Ok(sqlx::query_as::<_, MonitoredAddress>(
            "SELECT * FROM monitored_addresses WHERE address=$1",
        )
        .bind(address)
        .fetch_optional(&self.pool)
        .await?)
    }

    pub async fn update_monitor_status(
        &self,
        address: &str,
        status: &str,
        error: Option<&str>,
    ) -> Result<()> {
        sqlx::query("UPDATE monitored_addresses SET status=$2,last_error=$3,updated_at=now() WHERE address=$1")
            .bind(address)
            .bind(status)
            .bind(error.map(redact_error))
            .execute(&self.pool)
            .await?;
        Ok(())
    }

    pub async fn complete_backfill(
        &self,
        address: &str,
        coverage_start: DateTime<Utc>,
        history_complete: bool,
    ) -> Result<()> {
        sqlx::query("UPDATE monitored_addresses SET coverage_start=$2,history_complete=$3,status='LIVE',last_error=NULL,updated_at=now() WHERE address=$1")
            .bind(address)
            .bind(coverage_start)
            .bind(history_complete)
            .execute(&self.pool)
            .await?;
        Ok(())
    }

    pub async fn last_event_time(&self, address: &str) -> Result<Option<DateTime<Utc>>> {
        Ok(
            sqlx::query_scalar("SELECT last_event_time FROM monitored_addresses WHERE address=$1")
                .bind(address)
                .fetch_optional(&self.pool)
                .await?
                .flatten(),
        )
    }

    pub async fn ingest_fill(
        &self,
        address: &str,
        fill: &HyperliquidFill,
        transport: &str,
        is_snapshot: bool,
        publish: bool,
        metadata: &MetadataStore,
    ) -> Result<IngestResult> {
        let observed_at = Utc::now();
        let event_time =
            DateTime::<Utc>::from_timestamp_millis(fill.time).context("invalid fill timestamp")?;
        let source_event_id = normalization::source_event_id(address, fill);
        let raw_payload = serde_json::to_value(fill)?;
        let mut tx = self.pool.begin().await?;
        let raw_id = Uuid::new_v4();
        let inserted_raw: Option<Uuid> = sqlx::query_scalar(
            "INSERT INTO raw_trade_logs(id,address,source,transport,source_event_id,event_time,observed_at,is_snapshot,payload,parse_status) VALUES($1,$2,'HYPERLIQUID_OFFICIAL',$3,$4,$5,$6,$7,$8,'PARSED') ON CONFLICT(source,source_event_id) DO NOTHING RETURNING id",
        )
        .bind(raw_id)
        .bind(address)
        .bind(transport)
        .bind(&source_event_id)
        .bind(event_time)
        .bind(observed_at)
        .bind(is_snapshot)
        .bind(raw_payload)
        .fetch_optional(&mut *tx)
        .await?;
        let actual_raw_id = match inserted_raw {
            Some(id) => id,
            None => sqlx::query_scalar(
                "SELECT id FROM raw_trade_logs WHERE source='HYPERLIQUID_OFFICIAL' AND source_event_id=$1",
            )
            .bind(&source_event_id)
            .fetch_one(&mut *tx)
            .await?,
        };
        let fact =
            normalization::normalize(address, fill, actual_raw_id, observed_at, metadata).await?;
        let payload = serde_json::to_value(&fact.payload)?;
        let inserted_fact: Option<String> = sqlx::query_scalar(
            "INSERT INTO trade_fact_versions(fact_id,revision,event_id,account_key,account,occurred_at,observed_at,ordering_key,source,source_ref,raw_log_id,schema_version,payload) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT(fact_id,revision) DO NOTHING RETURNING fact_id",
        )
        .bind(&fact.fact_id)
        .bind(fact.revision)
        .bind(&fact.event_id)
        .bind(&fact.account_key)
        .bind(&fact.account)
        .bind(fact.occurred_at)
        .bind(fact.observed_at)
        .bind(&fact.ordering_key)
        .bind(&fact.source)
        .bind(&fact.source_ref)
        .bind(fact.raw_log_id)
        .bind(fact.schema_version)
        .bind(payload)
        .fetch_optional(&mut *tx)
        .await?;
        let inserted = inserted_fact.is_some();
        if inserted {
            sqlx::query("INSERT INTO trade_facts_current(fact_id,current_revision,account,occurred_at,ordering_key) VALUES($1,1,$2,$3,$4) ON CONFLICT(fact_id) DO NOTHING")
                .bind(&fact.fact_id)
                .bind(address)
                .bind(fact.occurred_at)
                .bind(&fact.ordering_key)
                .execute(&mut *tx)
                .await?;
            if publish && fact.payload.copy_eligible {
                let envelope = serde_json::to_value(&fact)?;
                sqlx::query("INSERT INTO outbox_events(id,event_id,fact_id,payload) VALUES($1,$2,$3,$4) ON CONFLICT(event_id) DO NOTHING")
                    .bind(Uuid::new_v4())
                    .bind(&fact.event_id)
                    .bind(&fact.fact_id)
                    .bind(envelope)
                    .execute(&mut *tx)
                    .await?;
            }
        }
        sqlx::query("UPDATE monitored_addresses SET last_event_time=GREATEST(COALESCE(last_event_time,$2),$2),last_received_at=$3,updated_at=now() WHERE address=$1")
            .bind(address)
            .bind(event_time)
            .bind(observed_at)
            .execute(&mut *tx)
            .await?;
        tx.commit().await?;
        Ok(IngestResult { inserted, fact })
    }

    pub async fn raw_trades(
        &self,
        address: &str,
        after: Option<(DateTime<Utc>, String)>,
        limit: i64,
    ) -> Result<Vec<RawTradeView>> {
        let (after_time, after_id) = after
            .map(|v| (Some(v.0), Some(v.1)))
            .unwrap_or((None, None));
        Ok(sqlx::query_as::<_, RawTradeView>("SELECT id,source_event_id,transport,event_time,observed_at,is_snapshot,parse_status,parse_error,payload FROM raw_trade_logs WHERE address=$1 AND ($2::timestamptz IS NULL OR (event_time,source_event_id)>($2,$3)) ORDER BY event_time,source_event_id LIMIT $4")
            .bind(address).bind(after_time).bind(after_id).bind(limit).fetch_all(&self.pool).await?)
    }

    pub async fn trade_facts(
        &self,
        address: &str,
        after: Option<(DateTime<Utc>, String)>,
        limit: i64,
    ) -> Result<Vec<TradeFactView>> {
        let (after_time, after_id) = after
            .map(|v| (Some(v.0), Some(v.1)))
            .unwrap_or((None, None));
        Ok(sqlx::query_as::<_, TradeFactView>("SELECT v.fact_id,v.event_id,v.revision,v.account_key,v.account,v.occurred_at,v.observed_at,v.ordering_key,v.source,v.source_ref,v.raw_log_id,v.schema_version,v.payload FROM trade_fact_versions v JOIN trade_facts_current c ON c.fact_id=v.fact_id AND c.current_revision=v.revision WHERE v.account=$1 AND ($2::timestamptz IS NULL OR (v.occurred_at,v.ordering_key)>($2,$3)) ORDER BY v.occurred_at,v.ordering_key LIMIT $4")
            .bind(address).bind(after_time).bind(after_id).bind(limit).fetch_all(&self.pool).await?)
    }

    pub async fn deliveries(&self, limit: i64) -> Result<Vec<DeliveryView>> {
        Ok(sqlx::query_as::<_, DeliveryView>("SELECT event_id,fact_id,status,attempts,available_at,last_error,delivered_at FROM outbox_events ORDER BY created_at DESC LIMIT $1")
            .bind(limit).fetch_all(&self.pool).await?)
    }

    pub async fn start_collection_run(
        &self,
        address: &str,
        kind: &str,
        start: DateTime<Utc>,
        end: DateTime<Utc>,
    ) -> Result<Uuid> {
        let id = Uuid::new_v4();
        sqlx::query("INSERT INTO collection_runs(id,address,kind,range_start,range_end,status) VALUES($1,$2,$3,$4,$5,'RUNNING')")
            .bind(id).bind(address).bind(kind).bind(start).bind(end)
            .execute(&self.pool).await?;
        Ok(id)
    }

    pub async fn finish_collection_run(
        &self,
        id: Uuid,
        status: &str,
        records_seen: i64,
        records_inserted: i64,
        error: Option<&str>,
    ) -> Result<()> {
        sqlx::query("UPDATE collection_runs SET status=$2,records_seen=$3,records_inserted=$4,error=$5,finished_at=now() WHERE id=$1")
            .bind(id).bind(status).bind(records_seen).bind(records_inserted)
            .bind(error.map(redact_error)).execute(&self.pool).await?;
        Ok(())
    }
}

fn redact_error(value: &str) -> String {
    value.chars().take(1000).collect()
}
