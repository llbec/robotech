mod common;
use common::*;
use trade_log::stored_query::StoredQuery;
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn complete_storage_repeat_and_concurrent_ingestion() {
    let db = database().await;
    let first = collect(&db, "query_first", fills(4), 1).await.unwrap();
    assert_eq!(first.trades.len(), 1);
    assert_eq!(first.persistence.unwrap().inserted_records, 4);
    assert_eq!(count(&db, "account_facts_current").await, 4);
    let (a, b) = tokio::join!(
        collect(&db, "query_a", fills(4), 1),
        collect(&db, "query_b", fills(4), 1)
    );
    for r in [a.unwrap(), b.unwrap()] {
        let p = r.persistence.unwrap();
        assert_eq!(p.inserted_records, 0);
        assert_eq!(p.existing_records, 4);
    }
    assert_eq!(count(&db, "account_fact_versions").await, 4);
    assert_eq!(count(&db, "fact_observations").await, 12);
    assert_eq!(count(&db, "raw_logs").await, 9);
    assert_eq!(db.stored(&request(2000)).await.unwrap().trades.len(), 4);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn conflict_rolls_back_entire_batch_preserves_raw() {
    let db = database().await;
    collect(&db, "query_initial", fills(1), 100).await.unwrap();
    let mut source = fills(2);
    source[0]["fee"] = serde_json::json!("1");
    assert_eq!(
        collect(&db, "query_conflict", source, 100)
            .await
            .unwrap_err()
            .code,
        "VERSION_CONFLICT"
    );
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert_eq!(count(&db, "fact_observations").await, 1);
    assert_eq!(count(&db, "raw_logs").await, 6);
    let (status, _, _) = db.job_result("query_conflict").await.unwrap();
    assert_eq!(status, "FAILED");
    assert_eq!(
        db.stored(&request(100)).await.unwrap().trades[0]
            .payload
            .fee,
        "-0.0125"
    );
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn interrupted_task_and_reconnect_do_not_lose_committed_facts() {
    use trade_log::{raw_log::RawEvidenceStore, validation::QueryRequest};
    let db = database().await;
    let fact = collect(&db, "query_committed", fills(1), 100)
        .await
        .unwrap()
        .trades[0]
        .clone();
    db.begin(
        "query_running",
        &QueryRequest {
            account: ACCOUNT.into(),
            limit: 100,
        },
        "trace_test",
    )
    .await
    .unwrap();
    db.interrupt().await.unwrap();
    assert_eq!(
        db.job_result("query_running").await.unwrap().0,
        "INTERRUPTED"
    );
    assert_eq!(db.stored(&request(100)).await.unwrap().trades[0], fact);
    // Kill an idle pooled connection; a later read must acquire a fresh connection.
    sqlx::query("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname=current_database() AND pid<>pg_backend_pid() AND state='idle'").execute(&db.pool).await.unwrap();
    assert_eq!(db.stored(&request(100)).await.unwrap().trades[0], fact);
    db.pool.close().await;
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn database_error_is_not_empty_success_and_file_mirror_is_optional() {
    use std::sync::Arc;
    use trade_log_adapters::file_evidence::FileEvidence;
    let mut db = database().await;
    let path = std::env::temp_dir().join(format!("robotech-mirror-{}", uuid::Uuid::new_v4()));
    std::fs::write(&path, b"not a directory").unwrap();
    db.mirror = Some(Arc::new(FileEvidence {
        directory: path.clone(),
        network: shared_types::Network::Mainnet,
    }));
    let result = collect(&db, "query_mirror", fills(1), 100).await.unwrap();
    assert_eq!(result.persistence.unwrap().status, "COMMITTED");
    assert_eq!(db.stored(&request(100)).await.unwrap().matched_records, 1);
    std::fs::remove_file(path).unwrap();
    db.pool.close().await;
    assert_eq!(
        db.stored(&request(100)).await.unwrap_err().code,
        "DEPENDENCY_UNAVAILABLE"
    );
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn migration_role_and_application_permissions() {
    use sqlx::Row;
    let db = database().await;
    let suffix = uuid::Uuid::new_v4().simple().to_string();
    let migrator = format!("migrator_{suffix}");
    let app = format!("app_{suffix}");
    sqlx::raw_sql(&format!("CREATE ROLE {migrator} LOGIN; CREATE ROLE {app} LOGIN; GRANT CONNECT,CREATE ON DATABASE \"{}\" TO {migrator}; GRANT USAGE,CREATE ON SCHEMA public TO {migrator}; GRANT USAGE ON SCHEMA trade_log TO {app}; GRANT SELECT,INSERT,UPDATE ON ALL TABLES IN SCHEMA trade_log TO {app}; GRANT SELECT ON public._sqlx_migrations TO {app};",db.pool.connect_options().get_database().unwrap())).execute(&db.pool).await.unwrap();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .connect_with((*db.pool.connect_options()).clone().username(&app))
        .await
        .unwrap();
    let app_store = trade_log_adapters::postgres::Postgres {
        pool,
        network: db.network.clone(),
        mirror: None,
    };
    app_store.check_schema().await.unwrap();
    collect(&app_store, "query_role", fills(1), 100)
        .await
        .unwrap();
    assert_eq!(
        app_store
            .stored(&request(100))
            .await
            .unwrap()
            .matched_records,
        1
    );
    assert!(
        sqlx::query("CREATE TABLE trade_log.forbidden(id int)")
            .execute(&app_store.pool)
            .await
            .is_err()
    );
    assert!(
        sqlx::query("DELETE FROM trade_log.account_fact_versions")
            .execute(&app_store.pool)
            .await
            .is_err()
    );
    let row = sqlx::query("SELECT count(*) AS n FROM trade_log.account_facts_current")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    assert_eq!(row.get::<i64, _>("n"), 1);
    app_store.pool.close().await;
    db.pool.close().await;
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn malformed_source_bytes_are_archived_before_failure() {
    use protocol_api::QueryKind;
    use trade_log::{
        acquisition::SourceResponse, query::QueryError, raw_log::RawEvidenceStore,
        validation::QueryRequest,
    };
    let db = database().await;
    let request = QueryRequest {
        account: ACCOUNT.into(),
        limit: 100,
    };
    db.begin("query_invalid_source", &request, "trace_test")
        .await
        .unwrap();
    let raw = db
        .save_response(
            "query_invalid_source",
            QueryKind::UserFills,
            1,
            &SourceResponse {
                status: 200,
                body: b"not json".to_vec(),
                received_at: shared_types::now(),
                retry_after_seconds: None,
            },
            &request,
        )
        .await
        .unwrap();
    db.finish(
        "query_invalid_source",
        Err(&QueryError::incomplete("Invalid source JSON")),
    )
    .await
    .unwrap();
    let (body, payload): (Vec<u8>, Option<serde_json::Value>) =
        sqlx::query_as("SELECT body,payload FROM trade_log.raw_logs WHERE source_event_id=$1")
            .bind(raw)
            .fetch_one(&db.pool)
            .await
            .unwrap();
    assert_eq!(body, b"not json");
    assert!(payload.is_none());
    assert_eq!(
        db.job_result("query_invalid_source").await.unwrap().0,
        "FAILED"
    );
    assert_eq!(count(&db, "account_facts_current").await, 0);
    db.pool.close().await;
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn default_migration_role_can_initialize_and_app_role_cannot_migrate() {
    use sqlx::{Connection, PgConnection, postgres::PgConnectOptions};
    use std::str::FromStr;
    let options =
        PgConnectOptions::from_str(&std::env::var("ROBOTECH_TEST_DATABASE_URL").unwrap()).unwrap();
    let mut admin = PgConnection::connect_with(&options).await.unwrap();
    let name = format!("robotech_role_test_{}", uuid::Uuid::new_v4().simple());
    sqlx::raw_sql("DO $$ BEGIN IF NOT EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_migrator') THEN CREATE ROLE trade_log_migrator LOGIN; END IF; IF NOT EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_app') THEN CREATE ROLE trade_log_app LOGIN; END IF; END $$;").execute(&mut admin).await.unwrap();
    sqlx::query(&format!("CREATE DATABASE {name}"))
        .execute(&mut admin)
        .await
        .unwrap();
    sqlx::query(&format!(
        "GRANT CONNECT,CREATE ON DATABASE {name} TO trade_log_migrator"
    ))
    .execute(&mut admin)
    .await
    .unwrap();
    let mut owner = PgConnection::connect_with(&options.clone().database(&name))
        .await
        .unwrap();
    sqlx::query("GRANT CREATE,USAGE ON SCHEMA public TO trade_log_migrator")
        .execute(&mut owner)
        .await
        .unwrap();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .connect_with(
            options
                .clone()
                .database(&name)
                .username("trade_log_migrator"),
        )
        .await
        .unwrap();
    let migrator = trade_log_adapters::postgres::Postgres {
        pool,
        network: shared_types::Network::Mainnet,
        mirror: None,
    };
    migrator.migrate().await.unwrap();
    migrator.migrate().await.unwrap();
    migrator.check_schema().await.unwrap();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .connect_with(options.database(&name).username("trade_log_app"))
        .await
        .unwrap();
    let app = trade_log_adapters::postgres::Postgres {
        pool,
        network: shared_types::Network::Mainnet,
        mirror: None,
    };
    app.check_schema().await.unwrap();
    collect(&app, "query_app_role", fills(1), 100)
        .await
        .unwrap();
    assert_eq!(app.stored(&request(100)).await.unwrap().matched_records, 1);
    assert!(app.migrate().await.is_err());
    app.pool.close().await;
    migrator.pool.close().await;
}
