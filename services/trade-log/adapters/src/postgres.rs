use account_facts::AccountFactEnvelope;
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use sqlx::{PgPool, Postgres, Transaction};
use trade_log::{normalization::raw_hash, PageWrite, TradeRepository};
use uuid::Uuid;

pub struct PostgresTradeRepository {
    pool: PgPool,
}
impl PostgresTradeRepository {
    pub async fn connect(database_url: &str) -> Result<Self, String> {
        let pool = PgPool::connect(database_url)
            .await
            .map_err(|e| format!("DATABASE_CONNECT_ERROR: {e}"))?;
        Ok(Self { pool })
    }
    pub async fn migrate(&self) -> Result<(), String> {
        sqlx::migrate!("../migrations")
            .run(&self.pool)
            .await
            .map_err(|e| format!("DATABASE_MIGRATION_ERROR: {e}"))
    }
}

#[async_trait]
impl TradeRepository for PostgresTradeRepository {
    async fn start_run(
        &self,
        address: &str,
        from: DateTime<Utc>,
        to: DateTime<Utc>,
    ) -> Result<Uuid, String> {
        let id = Uuid::new_v4();
        sqlx::query("INSERT INTO nansen_import_runs (id,address,from_time,to_time,status) VALUES ($1,$2,$3,$4,'RUNNING')")
            .bind(id).bind(address).bind(from).bind(to).execute(&self.pool).await.map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
        Ok(id)
    }

    async fn persist_page(
        &self,
        run_id: Uuid,
        write: PageWrite<'_>,
    ) -> Result<Vec<AccountFactEnvelope>, String> {
        let mut tx: Transaction<'_, Postgres> = self
            .pool
            .begin()
            .await
            .map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
        sqlx::query("INSERT INTO nansen_raw_pages (run_id,page_number,response) VALUES ($1,$2,$3) ON CONFLICT (run_id,page_number) DO NOTHING")
            .bind(run_id).bind(write.page as i32).bind(write.response).execute(&mut *tx).await.map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
        let mut saved = Vec::new();
        for (index, raw) in write.trades.iter().enumerate() {
            let hash = raw_hash(raw);
            let payload = serde_json::to_value(raw).map_err(|e| e.to_string())?;
            let raw_id: Uuid = sqlx::query_scalar("INSERT INTO raw_logs (id,run_id,page_number,record_index,source_id,content_hash,payload,transaction_ref,ordering_key,occurred_at,parse_status) VALUES ($1,$2,$3,$4,'NANSEN',$5,$6,$7,$8,$9,'PENDING') ON CONFLICT (content_hash) DO UPDATE SET content_hash=EXCLUDED.content_hash RETURNING id")
                .bind(Uuid::new_v4()).bind(run_id).bind(write.page as i32).bind(index as i32).bind(&hash).bind(payload).bind(&raw.transaction_hash)
                .bind(format!("{:020}:{}", raw.block_number.unwrap_or_default(), raw.transaction_hash)).bind(raw.timestamp).fetch_one(&mut *tx).await.map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
            if let Some(fact) = write.facts.iter().find(|f| f.fact.fact_id == hash) {
                let mut fact = fact.clone();
                fact.fact.raw_log_id = Some(raw_id.to_string());
                let payload =
                    serde_json::to_value(&fact.fact.payload).map_err(|e| e.to_string())?;
                let inserted = sqlx::query("INSERT INTO account_fact_versions (fact_id,revision,event_id,fact_type,account_key,ordering_key,sub_index,change_type,confirmation_status,occurred_at,payload,raw_log_id,schema_version) VALUES ($1,1,$2,'TRADE',$3,$4,$5,'UPSERT','OBSERVED',$6,$7,$8,1) ON CONFLICT (fact_id,revision) DO NOTHING")
                    .bind(&fact.fact.fact_id).bind(&fact.event_id).bind(&fact.fact.account_key).bind(&fact.fact.ordering_key).bind(fact.fact.sub_index as i32).bind(fact.occurred_at).bind(payload).bind(raw_id).execute(&mut *tx).await.map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
                sqlx::query(
                    "UPDATE raw_logs SET parse_status='PARSED', parse_error=NULL WHERE id=$1",
                )
                .bind(raw_id)
                .execute(&mut *tx)
                .await
                .map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
                let _was_inserted = inserted.rows_affected() > 0;
                saved.push(fact);
            } else if let Some(failure) = write.failures.iter().find(|f| f.raw_hash == hash) {
                sqlx::query(
                    "UPDATE raw_logs SET parse_status='FAILED', parse_error=$2 WHERE id=$1",
                )
                .bind(raw_id)
                .bind(&failure.reason)
                .execute(&mut *tx)
                .await
                .map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
            }
        }
        tx.commit()
            .await
            .map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
        Ok(saved)
    }

    async fn complete_run(
        &self,
        run_id: Uuid,
        pages: u32,
        raw: u64,
        normalized: u64,
        failed: u64,
    ) -> Result<(), String> {
        sqlx::query("UPDATE nansen_import_runs SET status='COMPLETED',pages_fetched=$2,raw_count=$3,normalized_count=$4,failed_count=$5,finished_at=now() WHERE id=$1")
            .bind(run_id).bind(pages as i32).bind(raw as i64).bind(normalized as i64).bind(failed as i64).execute(&self.pool).await.map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
        Ok(())
    }
    async fn fail_run(&self, run_id: Uuid, error: &str) -> Result<(), String> {
        let safe = if error.len() > 1000 {
            &error[..1000]
        } else {
            error
        };
        sqlx::query("UPDATE nansen_import_runs SET status='FAILED',last_error=$2,finished_at=now() WHERE id=$1").bind(run_id).bind(safe).execute(&self.pool).await.map_err(|e| format!("DATABASE_WRITE_ERROR: {e}"))?;
        Ok(())
    }
}
