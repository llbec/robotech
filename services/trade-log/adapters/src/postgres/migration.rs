use super::*;
pub static MIGRATOR: sqlx::migrate::Migrator = sqlx::migrate!("../migrations");
impl Postgres {
    pub async fn migrate(&self) -> Result<(), QueryError> {
        MIGRATOR
            .run(&self.pool)
            .await
            .map_err(|_| QueryError::unavailable("Trade database migration failed"))?;
        sqlx::raw_sql("DO $$ BEGIN IF EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_app') THEN GRANT USAGE ON SCHEMA trade_log TO trade_log_app; GRANT SELECT,INSERT,UPDATE ON ALL TABLES IN SCHEMA trade_log TO trade_log_app; GRANT SELECT ON TABLE public._sqlx_migrations TO trade_log_app; END IF; END $$").execute(&self.pool).await.map_err(db_error)?;
        sqlx::raw_sql("DO $$ BEGIN IF EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_collector') THEN GRANT USAGE ON SCHEMA trade_log TO trade_log_collector; GRANT SELECT,INSERT,UPDATE ON ALL TABLES IN SCHEMA trade_log TO trade_log_collector; GRANT SELECT ON TABLE public._sqlx_migrations TO trade_log_collector; END IF; END $$").execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
    pub async fn check_schema(&self) -> Result<(), QueryError> {
        let rows: Vec<(i64, Vec<u8>, bool)> = sqlx::query_as(
            "SELECT version,checksum,success FROM public._sqlx_migrations ORDER BY version",
        )
        .fetch_all(&self.pool)
        .await
        .map_err(db_error)?;
        let expected: Vec<_> = MIGRATOR
            .iter()
            .filter(|m| !m.migration_type.is_down_migration())
            .collect();
        if rows.len() != expected.len()
            || rows
                .iter()
                .zip(expected)
                .any(|((v, c, s), m)| *v != m.version || c.as_slice() != m.checksum.as_ref() || !*s)
        {
            return Err(QueryError::unavailable(
                "Trade database schema migration required",
            ));
        }
        Ok(())
    }
}
