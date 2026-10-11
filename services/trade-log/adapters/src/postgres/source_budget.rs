use super::*;
use sqlx::Row;
use uuid::Uuid;
impl Postgres {
    pub async fn reserve_source(
        &self,
        kind: &str,
        weight: i32,
    ) -> Result<Option<Uuid>, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        sqlx::query("SELECT pg_advisory_xact_lock(600006)")
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let settings=sqlx::query("SELECT http_budget,ws_budget,max_http_requests FROM trade_log.monitoring_settings WHERE id=1").fetch_one(&mut *tx).await.map_err(db_error)?;
        sqlx::query("DELETE FROM trade_log.source_budget_reservations WHERE reserved_at<=clock_timestamp()-interval '60 seconds' AND (active_until IS NULL OR active_until<=clock_timestamp())").execute(&mut *tx).await.map_err(db_error)?;
        let used:i64=sqlx::query_scalar("SELECT COALESCE(sum(weight),0)::bigint FROM trade_log.source_budget_reservations WHERE network=$1 AND kind=$2 AND reserved_at>clock_timestamp()-interval '60 seconds'").bind(self.network.name()).bind(kind).fetch_one(&mut *tx).await.map_err(db_error)?;
        let max: i32 = settings.get(if kind == "HTTP" {
            "http_budget"
        } else {
            "ws_budget"
        });
        if weight <= 0 || weight > max {
            return Err(QueryError::limited());
        }
        if used + i64::from(weight) > i64::from(max) {
            return Ok(None);
        }
        if kind == "HTTP" {
            let inflight:i64=sqlx::query_scalar("SELECT count(*) FROM trade_log.source_budget_reservations WHERE network=$1 AND kind='HTTP' AND active_until>clock_timestamp()").bind(self.network.name()).fetch_one(&mut *tx).await.map_err(db_error)?;
            if inflight >= i64::from(settings.get::<i32, _>("max_http_requests")) {
                return Ok(None);
            }
        }
        let id = Uuid::new_v4();
        sqlx::query("INSERT INTO trade_log.source_budget_reservations(id,network,kind,weight,active_until) VALUES($1,$2,$3,$4,CASE WHEN $3='HTTP' THEN clock_timestamp()+interval '120 seconds' ELSE NULL END)").bind(id).bind(self.network.name()).bind(kind).bind(weight).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        Ok(Some(id))
    }
    pub async fn release_source(&self, id: Uuid) -> Result<(), QueryError> {
        sqlx::query(
            "UPDATE trade_log.source_budget_reservations SET active_until=NULL WHERE id=$1",
        )
        .bind(id)
        .execute(&self.pool)
        .await
        .map_err(db_error)?;
        Ok(())
    }
    pub async fn settle_source(&self, id: Uuid, weight: i32) -> Result<(), QueryError> {
        // Only reduce a conservative reservation; an uncertain result keeps its full cost.
        sqlx::query("UPDATE trade_log.source_budget_reservations SET active_until=NULL,weight=LEAST(weight,$2) WHERE id=$1 AND $2>0").bind(id).bind(weight).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
}
