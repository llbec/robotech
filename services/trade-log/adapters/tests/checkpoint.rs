mod common;
use common::*;
use trade_log::checkpoint::CollectionStatusReader;
use trade_log::stored_query::StoredQuery;
const BASE: i64 = 1791097200000;
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn quiet_success_advances_watermark_and_resume_ignores_new_config_start() {
    let db = database().await;
    let c = collection_config();
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db
        .collection_work(&l, &c, BASE + 5000)
        .await
        .unwrap()
        .unwrap();
    let end = w.range.end_ms;
    assert!(
        runtime(&db, c.clone(), serde_json::json!([]))
            .round(&l, &mut w)
            .await
            .unwrap()
    );
    let status = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert_eq!(
        status.scanned_through.as_deref(),
        Some(trade_log::collection::time(end).as_str())
    );
    assert!(status.last_success_at.is_some());
    assert!(status.last_trade_at.is_none());
    assert_eq!(count(&db, "account_facts_current").await, 0);
    db.release_collection(&l).await.unwrap();
    let mut changed = c.clone();
    changed.start_time = "2026-10-04T07:00:02Z".into();
    let l2 = db
        .acquire_collection(&changed, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    no_schedule(&db).await;
    let w2 = db
        .collection_work(&l2, &changed, BASE + 9000)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(w2.range.start_ms, BASE + 1000);
    assert!(w2.range.end_ms > end);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn old_lease_cannot_commit_and_query_restart_does_not_interrupt_collector() {
    let db = database().await;
    let c = collection_config();
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db
        .collection_work(&l, &c, BASE + 5000)
        .await
        .unwrap()
        .unwrap();
    assert!(
        db.acquire_collection(&c, uuid::Uuid::new_v4())
            .await
            .unwrap()
            .is_none()
    );
    db.interrupt().await.unwrap();
    assert_eq!(db.job_result(&w.query_id).await.unwrap().0, "RUNNING");
    sqlx::query(
        "UPDATE trade_log.collection_checkpoints SET lease_expires_at=now()-interval '1 second'",
    )
    .execute(&db.pool)
    .await
    .unwrap();
    let newer = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    assert!(newer.epoch > l.epoch);
    assert_eq!(
        runtime(&db, c.clone(), fills(1))
            .round(&l, &mut w)
            .await
            .unwrap_err()
            .code,
        "LEASE_LOST"
    );
    assert_eq!(count(&db, "account_facts_current").await, 0);
    let mut restored = db
        .collection_work(&newer, &c, BASE + 5000)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(restored.query_id, w.query_id);
    assert!(
        runtime(&db, c, fills(1))
            .round(&newer, &mut restored)
            .await
            .unwrap()
    );
    assert_eq!(count(&db, "account_facts_current").await, 1);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn conflict_rolls_back_new_facts_and_watermark_together() {
    let db = database().await;
    collect(&db, "query_original", fills(1), 100).await.unwrap();
    let c = collection_config();
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db
        .collection_work(&l, &c, BASE + 5000)
        .await
        .unwrap()
        .unwrap();
    let mut changed = fills(2);
    changed[0]["fee"] = serde_json::json!("123");
    let e = runtime(&db, c, changed)
        .round(&l, &mut w)
        .await
        .unwrap_err();
    assert_eq!(e.code, "VERSION_CONFLICT");
    db.collection_failure(&l, &e, 1, false).await.unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 1);
    let status = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert!(status.scanned_through.is_none());
    assert!(status.last_success_at.is_none());
    assert_eq!(status.status, "FAILED");
    assert_eq!(status.consecutive_failures, 1);
    assert_eq!(db.job_result(&w.query_id).await.unwrap().0, "FAILED");
    assert!(count(&db, "raw_logs").await > 3);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn splitting_budget_resume_and_multi_page_reparse_preserve_all_records() {
    let db = database().await;
    let mut c = collection_config();
    c.safety_delay_seconds = 0;
    c.max_requests_per_round = 3;
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db
        .collection_work(&l, &c, BASE + 2001)
        .await
        .unwrap()
        .unwrap();
    assert!(
        !runtime(&db, c.clone(), fills(2001))
            .round(&l, &mut w)
            .await
            .unwrap()
    );
    assert_eq!(w.remaining.len(), 2);
    assert_eq!(count(&db, "account_facts_current").await, 0);
    assert!(
        db.collection_status(ACCOUNT).await.unwrap().items[0]
            .scanned_through
            .is_none()
    );
    db.release_collection(&l).await.unwrap();
    let l2 = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut resumed = db
        .collection_work(&l2, &c, BASE + 2001)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(resumed.query_id, w.query_id);
    assert_eq!(resumed.remaining.len(), 2);
    assert!(
        runtime(&db, c, fills(2001))
            .round(&l2, &mut resumed)
            .await
            .unwrap()
    );
    assert_eq!(count(&db, "account_facts_current").await, 2001);
    assert_eq!(resumed.pages.len(), 2);
    let report = db
        .reparse(&w.query_id, &hyperliquid::parser::HyperliquidParser)
        .await
        .unwrap();
    assert_eq!(report.comparison, "SAME");
    assert_eq!(report.source_records, 2001);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn saturated_millisecond_never_advances_checkpoint() {
    let db = database().await;
    let mut c = collection_config();
    c.safety_delay_seconds = 0;
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db.collection_work(&l, &c, BASE + 1).await.unwrap().unwrap();
    let mut data = fills(2000);
    for f in data.as_array_mut().unwrap() {
        f["time"] = serde_json::json!(BASE);
    }
    let e = runtime(&db, c, data).round(&l, &mut w).await.unwrap_err();
    assert_eq!(e.code, "SOURCE_WINDOW_SATURATED");
    db.collection_failure(&l, &e, 1, false).await.unwrap();
    let status = &db.collection_status(ACCOUNT).await.unwrap().items[0];
    assert_eq!(status.status, "FAILED");
    assert!(status.scanned_through.is_none());
    assert_eq!(count(&db, "account_facts_current").await, 0);
    assert_eq!(count(&db, "raw_logs").await, 3);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn overlapping_round_and_manual_collection_share_fact_identity() {
    let db = database().await;
    let mut c = collection_config();
    c.safety_delay_seconds = 0;
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db
        .collection_work(&l, &c, BASE + 500)
        .await
        .unwrap()
        .unwrap();
    let r = runtime(&db, c.clone(), fills(2));
    assert!(r.round(&l, &mut w).await.unwrap());
    no_schedule(&db).await;
    let mut w2 = db
        .collection_work(&l, &c, BASE + 1000)
        .await
        .unwrap()
        .unwrap();
    assert!(r.round(&l, &mut w2).await.unwrap());
    let manual = collect(&db, "query_manual_after", fills(2), 100)
        .await
        .unwrap();
    assert_eq!(manual.persistence.unwrap().inserted_records, 0);
    assert_eq!(count(&db, "account_facts_current").await, 2);
    assert_eq!(count(&db, "account_fact_versions").await, 2);
    assert_eq!(count(&db, "fact_observations").await, 6);
    let status = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert_eq!(
        status.last_success_query_id.as_deref(),
        Some(w2.query_id.as_str())
    );
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn checkpoint_failure_after_fact_writes_rolls_back_the_entire_commit() {
    use trade_log_adapters::postgres::checkpoint::CollectionCommit;
    let db = database().await;
    let mut c = collection_config();
    c.safety_delay_seconds = 0;
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut w = db.collection_work(&l, &c, BASE + 1).await.unwrap().unwrap();
    assert!(
        runtime(&db, c.clone(), fills(1))
            .round(&l, &mut w)
            .await
            .unwrap()
    );
    no_schedule(&db).await;
    let mut pending = db.collection_work(&l, &c, BASE + 3).await.unwrap().unwrap();
    let mut budget = c.clone();
    budget.max_requests_per_round = 3;
    // Reading all pages consumes the budget, but commit is attempted as part of round.
    // Trigger a database failure in the checkpoint update, after fact INSERTs.
    sqlx::raw_sql("CREATE FUNCTION trade_log.reject_checkpoint() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.status='WAITING' THEN RAISE EXCEPTION 'test rollback'; END IF; RETURN NEW; END $$; CREATE TRIGGER reject_checkpoint BEFORE UPDATE ON trade_log.collection_checkpoints FOR EACH ROW EXECUTE FUNCTION trade_log.reject_checkpoint();").execute(&db.pool).await.unwrap();
    assert!(
        runtime(&db, budget, fills(3))
            .round(&l, &mut pending)
            .await
            .is_err()
    );
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert_eq!(count(&db, "account_fact_versions").await, 1);
    assert_eq!(db.job_result(&pending.query_id).await.unwrap().0, "RUNNING");
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0]
            .scanned_through
            .as_deref(),
        Some(trade_log::collection::time(BASE + 1).as_str())
    );
    sqlx::query("DROP TRIGGER reject_checkpoint ON trade_log.collection_checkpoints")
        .execute(&db.pool)
        .await
        .unwrap();
    let ids: Vec<_> = pending.pages.iter().map(|p| p.raw_id.clone()).collect();
    let (result, observations) = db
        .collection_result(
            &pending.query_id,
            ACCOUNT,
            &ids,
            pending.meta.as_deref().unwrap(),
            pending.spot_meta.as_deref().unwrap(),
            &hyperliquid::parser::HyperliquidParser,
        )
        .await
        .unwrap();
    db.commit_collection(
        &result,
        &CollectionCommit {
            lease: &l,
            work: &pending,
            observations: &observations,
            interval_seconds: 1,
        },
    )
    .await
    .unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 3);
    // A lost reply can repeat a completed transaction without re-inserting facts.
    db.commit_collection(
        &result,
        &CollectionCommit {
            lease: &l,
            work: &pending,
            observations: &observations,
            interval_seconds: 1,
        },
    )
    .await
    .unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 3);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn incremental_migration_preserves_legacy_facts_and_collector_role_is_limited() {
    use sqlx::{Connection, PgConnection, postgres::PgConnectOptions};
    use std::str::FromStr;
    let options =
        PgConnectOptions::from_str(&std::env::var("ROBOTECH_TEST_DATABASE_URL").unwrap()).unwrap();
    let mut admin = PgConnection::connect_with(&options).await.unwrap();
    let name = format!("robotech_upgrade_{}", uuid::Uuid::new_v4().simple());
    sqlx::query(&format!("CREATE DATABASE {name}"))
        .execute(&mut admin)
        .await
        .unwrap();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .connect_with(options.clone().database(&name))
        .await
        .unwrap();
    let first = trade_log_adapters::postgres::migration::MIGRATOR
        .iter()
        .next()
        .unwrap()
        .clone();
    let initial = sqlx::migrate::Migrator {
        migrations: std::borrow::Cow::Owned(vec![first]),
        ..sqlx::migrate::Migrator::DEFAULT
    };
    initial.run(&pool).await.unwrap();
    let db = trade_log_adapters::postgres::Postgres {
        pool,
        network: shared_types::Network::Mainnet,
        mirror: None,
    };
    let saved = legacy_collect(&db, "query_v02_upgrade", fills(2))
        .await
        .unwrap();
    let ids: Vec<_> = saved.trades.iter().map(|f| f.fact_id.clone()).collect();
    assert!(db.check_schema().await.is_err());
    sqlx::raw_sql("DO $$ BEGIN IF NOT EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_collector') THEN CREATE ROLE trade_log_collector LOGIN; END IF; END $$;").execute(&mut admin).await.unwrap();
    db.migrate().await.unwrap();
    db.check_schema().await.unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 2);
    let reread = db.stored(&request(100)).await.unwrap();
    assert_eq!(
        reread
            .trades
            .iter()
            .map(|f| f.fact_id.clone())
            .collect::<Vec<_>>(),
        ids
    );
    assert_eq!(
        db.reparse("query_v02_upgrade", &hyperliquid::parser::HyperliquidParser)
            .await
            .unwrap()
            .comparison,
        "SAME"
    );
    let pool = sqlx::postgres::PgPoolOptions::new()
        .connect_with(options.database(&name).username("trade_log_collector"))
        .await
        .unwrap();
    let collector = trade_log_adapters::postgres::Postgres {
        pool,
        network: db.network.clone(),
        mirror: None,
    };
    collector.check_schema().await.unwrap();
    let c = collection_config();
    let l = collector
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut work = collector
        .collection_work(&l, &c, BASE + 5000)
        .await
        .unwrap()
        .unwrap();
    assert!(
        runtime(&collector, c, fills(2))
            .round(&l, &mut work)
            .await
            .unwrap()
    );
    assert!(
        sqlx::query("CREATE TABLE trade_log.denied(id INT)")
            .execute(&collector.pool)
            .await
            .is_err()
    );
    assert!(
        sqlx::query("DELETE FROM trade_log.collection_checkpoints")
            .execute(&collector.pool)
            .await
            .is_err()
    );
    collector.pool.close().await;
    db.pool.close().await;
}
