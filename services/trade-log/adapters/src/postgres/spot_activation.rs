use super::*;
use serde_json::json;
use sqlx::Row;
use trade_log::publishing::CandidatePolicy;
impl Postgres {
    pub async fn configure_spot_publishing(
        &self,
        key: &str,
        enabled: bool,
        policy: &CandidatePolicy,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        sqlx::query("SELECT pg_advisory_xact_lock(hashtextextended($1,6))")
            .bind(key)
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let watch = sqlx::query(
            "SELECT generation,desired_state FROM trade_log.watch_accounts WHERE account_key=$1",
        )
        .bind(key)
        .fetch_optional(&mut *tx)
        .await
        .map_err(db_error)?;
        let generation = watch
            .as_ref()
            .map(|r| r.get::<i64, _>("generation"))
            .unwrap_or(1);
        let active = watch
            .as_ref()
            .is_none_or(|r| r.get::<String, _>("desired_state") == "ACTIVE");
        sqlx::query("INSERT INTO trade_log.instrument_activation(account_key,instrument_type,enabled,activation_epoch,monitoring_generation,policy) VALUES($1,'SPOT',$2,CASE WHEN $2 THEN 1 ELSE 0 END,$3,$4) ON CONFLICT(account_key,instrument_type) DO UPDATE SET enabled=EXCLUDED.enabled,activation_epoch=instrument_activation.activation_epoch+CASE WHEN EXCLUDED.enabled AND (NOT instrument_activation.enabled OR instrument_activation.policy<>EXCLUDED.policy OR ($5 AND instrument_activation.monitoring_generation<>EXCLUDED.monitoring_generation)) THEN 1 ELSE 0 END,activated_at=CASE WHEN EXCLUDED.enabled AND (NOT instrument_activation.enabled OR instrument_activation.policy<>EXCLUDED.policy OR ($5 AND instrument_activation.monitoring_generation<>EXCLUDED.monitoring_generation)) THEN clock_timestamp() ELSE instrument_activation.activated_at END,monitoring_generation=CASE WHEN $5 THEN EXCLUDED.monitoring_generation ELSE instrument_activation.monitoring_generation END,policy=EXCLUDED.policy")
            .bind(key).bind(enabled).bind(generation).bind(json!(policy)).bind(active).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
}
