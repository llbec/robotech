mod common;
use common::*;
use hyperliquid::parser::HyperliquidParser;
use serde_json::{Value, json};
use trade_log::stored_query::StoredQuery;
use trade_log::{
    checkpoint::{CollectionStatusReader, Lease},
    collection::CollectionConfig,
    realtime::{StreamEvent, StreamMessage},
};
use trade_log_adapters::postgres::{Postgres, realtime::StreamCommit};
use uuid::Uuid;

async fn start(db: &Postgres) -> (CollectionConfig, Lease, Uuid) {
    let c = collection_config();
    let l = db
        .acquire_collection(&c, Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    db.prepare_stream(&l, &c, true).await.unwrap();
    let session = Uuid::new_v4();
    db.begin_stream(&l, &c, session).await.unwrap();
    db.stream_connected(&l, session).await.unwrap();
    db.stream_subscribed(&l, session).await.unwrap();
    (c, l, session)
}
fn message(fills: Value, mode: Option<bool>) -> StreamMessage {
    let mut v = json!({"channel":"userFills","data":{"user":ACCOUNT,"fills":fills}});
    if let Some(mode) = mode {
        v["data"]["isSnapshot"] = json!(mode);
    }
    let StreamEvent::Data(m) =
        hyperliquid::websocket::decode(&serde_json::to_vec(&v).unwrap(), ACCOUNT).unwrap()
    else {
        panic!()
    };
    m
}
async fn archive(db: &Postgres, l: &Lease, session: Uuid, seq: i64, fills: Value) -> String {
    db.archive_stream(
        l,
        session,
        seq,
        &message(fills, Some(false)),
        META.as_bytes(),
        SPOT.as_bytes(),
    )
    .await
    .unwrap()
}
async fn commit(
    db: &Postgres,
    l: &Lease,
    session: Uuid,
    seq: i64,
    id: &str,
) -> Result<(), trade_log::query::QueryError> {
    let result = db.stream_result(id, &HyperliquidParser).await?;
    db.commit_stream(
        &result,
        &StreamCommit {
            lease: l,
            session_id: session,
            sequence: seq,
        },
    )
    .await
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn http_ws_duplicates_preserve_legacy_hash_and_both_observations() {
    let db = database().await;
    let old = collect(&db, "query_before_ws", fills(1), 100)
        .await
        .unwrap();
    sqlx::query("UPDATE trade_log.account_fact_versions SET semantic_hash_version=NULL,semantic_content_hash=NULL").execute(&db.pool).await.unwrap();
    let hash: String =
        sqlx::query_scalar("SELECT content_hash FROM trade_log.account_fact_versions")
            .fetch_one(&db.pool)
            .await
            .unwrap();
    let (_, l, session) = start(&db).await;
    let mut duplicate = fills(1);
    duplicate[0]["px"] = json!("2500.2500");
    duplicate[0]["unusedOptionalField"] = json!(true);
    let id = archive(&db, &l, session, 1, duplicate.clone()).await;
    // Exact repeated archive is idempotent, but another envelope is a new observation.
    assert_eq!(archive(&db, &l, session, 1, duplicate.clone()).await, id);
    commit(&db, &l, session, 1, &id).await.unwrap();
    let next = archive(&db, &l, session, 2, duplicate).await;
    commit(&db, &l, session, 2, &next).await.unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert_eq!(count(&db, "account_fact_versions").await, 1);
    assert_eq!(count(&db, "fact_observations").await, 3);
    assert_eq!(
        db.stored(&request(100)).await.unwrap().trades[0],
        old.trades[0]
    );
    assert_eq!(
        sqlx::query_scalar::<_, String>("SELECT content_hash FROM trade_log.account_fact_versions")
            .fetch_one(&db.pool)
            .await
            .unwrap(),
        hash
    );
    assert!(
        sqlx::query_scalar::<_, String>(
            "SELECT semantic_content_hash FROM trade_log.account_fact_versions"
        )
        .fetch_one(&db.pool)
        .await
        .is_ok()
    );
    for query in ["query_before_ws", &id, &next] {
        assert_eq!(
            db.reparse(query, &HyperliquidParser)
                .await
                .unwrap()
                .comparison,
            "SAME"
        );
    }
    let status = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert!(status.scanned_through.is_none());
    assert!(status.websocket["last_committed_at"].is_string());
    assert_eq!(status.monitoring_status, "RECOVERING");
    let bytes: Vec<u8> = sqlx::query_scalar(
        "SELECT body FROM trade_log.raw_logs WHERE transport='WEBSOCKET' LIMIT 1",
    )
    .fetch_one(&db.pool)
    .await
    .unwrap();
    assert_eq!(
        serde_json::from_slice::<Value>(&bytes).unwrap()["channel"],
        "userFills"
    );
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn genuine_ws_conflict_rolls_back_new_facts_observations_and_positions() {
    let db = database().await;
    collect(&db, "query_initial", fills(1), 100).await.unwrap();
    let (_, l, session) = start(&db).await;
    let mut batch = fills(2);
    batch[0]["fee"] = json!("99");
    let id = archive(&db, &l, session, 1, batch).await;
    let e = commit(&db, &l, session, 1, &id).await.unwrap_err();
    assert_eq!(e.code, "VERSION_CONFLICT");
    db.fail_stream(&l, &id, &e).await.unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert_eq!(count(&db, "fact_observations").await, 1);
    let state = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert!(state.websocket["last_committed_at"].is_null());
    assert!(state.scanned_through.is_none());
    assert!(db.stream_fatal_error(&l).await.unwrap().is_some());
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn stale_epoch_cannot_archive_commit_or_change_ws_state() {
    let db = database().await;
    let (c, old, session) = start(&db).await;
    let id = archive(&db, &old, session, 1, fills(1)).await;
    sqlx::query(
        "UPDATE trade_log.collection_checkpoints SET lease_expires_at=now()-interval '1 second'",
    )
    .execute(&db.pool)
    .await
    .unwrap();
    let new = db
        .acquire_collection(&c, Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    assert!(new.epoch > old.epoch);
    assert_eq!(
        commit(&db, &old, session, 1, &id).await.unwrap_err().code,
        "LEASE_LOST"
    );
    assert!(
        db.stream_state(&old, json!({"status":"LIVE"}))
            .await
            .is_err()
    );
    assert!(
        db.archive_stream(
            &old,
            session,
            2,
            &message(fills(1), Some(false)),
            META.as_bytes(),
            SPOT.as_bytes()
        )
        .await
        .is_err()
    );
    assert_eq!(count(&db, "account_facts_current").await, 0);
    db.prepare_stream(&new, &c, true).await.unwrap();
    db.resume_stream_job(&new, &id).await.unwrap();
    commit(&db, &new, session, 1, &id).await.unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 1);
    let stopped: String = sqlx::query_scalar(
        "SELECT status FROM trade_log.collection_stream_sessions WHERE session_id=$1",
    )
    .bind(session)
    .fetch_one(&db.pool)
    .await
    .unwrap();
    assert_eq!(stopped, "INTERRUPTED");
    assert_eq!(
        db.reparse(&id, &HyperliquidParser)
            .await
            .unwrap()
            .comparison,
        "SAME"
    );
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn newer_ws_trade_never_skips_http_gap_and_gap_completion_is_atomic() {
    let db = database().await;
    let now = chrono::Utc::now().timestamp_millis();
    let mut c = collection_config();
    c.start_time = trade_log::collection::time(now - 5000);
    c.safety_delay_seconds = 0;
    let l = db
        .acquire_collection(&c, Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    db.prepare_stream(&l, &c, true).await.unwrap();
    let session = Uuid::new_v4();
    db.begin_stream(&l, &c, session).await.unwrap();
    db.stream_connected(&l, session).await.unwrap();
    db.stream_subscribed(&l, session).await.unwrap();
    let mut recent = fills(1);
    recent[0]["time"] = json!(now + 500);
    let id = archive(&db, &l, session, 1, recent.clone()).await;
    commit(&db, &l, session, 1, &id).await.unwrap();
    assert!(
        db.collection_status(ACCOUNT).await.unwrap().items[0]
            .scanned_through
            .is_none()
    );
    let mut history = fills(2);
    history[0] = recent[0].clone();
    history[1]["time"] = json!(now - 1000);
    let runner = runtime(&db, c.clone(), history);
    let mut work = db
        .collection_work(&l, &c, now + 2000)
        .await
        .unwrap()
        .unwrap();
    sqlx::raw_sql("CREATE FUNCTION trade_log.reject_gap() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'injected gap failure'; END $$; CREATE TRIGGER reject_gap BEFORE UPDATE ON trade_log.collection_gaps FOR EACH ROW EXECUTE FUNCTION trade_log.reject_gap();").execute(&db.pool).await.unwrap();
    assert!(runner.round(&l, &mut work).await.is_err());
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert!(
        db.collection_status(ACCOUNT).await.unwrap().items[0]
            .scanned_through
            .is_none()
    );
    sqlx::raw_sql("DROP TRIGGER reject_gap ON trade_log.collection_gaps")
        .execute(&db.pool)
        .await
        .unwrap();
    assert!(runner.round(&l, &mut work).await.unwrap());
    let status = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert_eq!(status.recovery["status"], "HTTP_SCANNED");
    assert!(status.scanned_through.is_some());
    assert_eq!(count(&db, "account_facts_current").await, 2);
    assert_eq!(status.recovery["open_gap_count"], 0);
    db.stream_pong(&l, session).await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].monitoring_status,
        "LIVE"
    );
    db.end_stream(
        &l,
        &c,
        session,
        &trade_log::realtime::stream_error("DISCONNECT", "test disconnect"),
        false,
    )
    .await
    .unwrap();
    db.prepare_stream(&l, &c, false).await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].monitoring_status,
        "HTTP_ONLY"
    );
    no_schedule(&db).await;
    let mut work = db
        .collection_work(&l, &c, now + 3000)
        .await
        .unwrap()
        .unwrap();
    runner.round(&l, &mut work).await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].recovery["status"],
        "DISABLED"
    );
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn empty_unknown_mode_malformed_and_metadata_refresh_are_auditable() {
    let db = database().await;
    let (_, l, session) = start(&db).await;
    let unknown = message(json!([]), None);
    let empty = db
        .archive_stream(&l, session, 1, &unknown, META.as_bytes(), SPOT.as_bytes())
        .await
        .unwrap();
    commit(&db, &l, session, 1, &empty).await.unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 0);
    assert_eq!(
        db.reparse(&empty, &HyperliquidParser)
            .await
            .unwrap()
            .comparison,
        "SAME"
    );
    let mut fill = fills(1);
    fill[0]["coin"] = json!("NEW");
    let id = archive(&db, &l, session, 2, fill).await;
    assert_eq!(
        db.stream_result(&id, &HyperliquidParser)
            .await
            .unwrap_err()
            .code,
        "INCOMPLETE_DATA"
    );
    let mut meta: Value = serde_json::from_str(META).unwrap();
    meta["universe"]
        .as_array_mut()
        .unwrap()
        .push(json!({"name":"NEW"}));
    db.update_stream_metadata(
        &l,
        &id,
        &serde_json::to_vec(&meta).unwrap(),
        SPOT.as_bytes(),
    )
    .await
    .unwrap();
    commit(&db, &l, session, 2, &id).await.unwrap();
    assert_eq!(
        db.reparse(&id, &HyperliquidParser)
            .await
            .unwrap()
            .comparison,
        "SAME"
    );
    let request: Value =
        sqlx::query_scalar("SELECT request FROM trade_log.collection_jobs WHERE query_id=$1")
            .bind(&id)
            .fetch_one(&db.pool)
            .await
            .unwrap();
    assert_eq!(request["meta_snapshots"].as_array().unwrap().len(), 2);
    assert_eq!(request["meta_snapshots"][0]["meta"], META);
    let invalid = StreamMessage {
        body: include_bytes!("../../../../tests/fixtures/v0.4/malformed.json").to_vec(),
        received_at: shared_types::now(),
        mode: "UNKNOWN".into(),
        fills: b"[]".to_vec(),
    };
    let bad = db
        .archive_stream(&l, session, 3, &invalid, META.as_bytes(), SPOT.as_bytes())
        .await
        .unwrap();
    assert!(db.stream_result(&bad, &HyperliquidParser).await.is_err());
    sqlx::query("UPDATE trade_log.raw_logs SET body='tampered'::bytea WHERE collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$1)").bind(&id).execute(&db.pool).await.unwrap();
    assert!(db.reparse(&id, &HyperliquidParser).await.is_err());
}

use async_trait::async_trait;
use std::{
    sync::{
        Arc,
        atomic::{AtomicUsize, Ordering},
    },
    time::Duration,
};
use tokio::{
    sync::{Mutex, mpsc},
    time::Instant,
};
use tokio_util::sync::CancellationToken;
use trade_log::{
    query::QueryError,
    realtime::{RealtimeConnection, RealtimeSource, WebsocketConfig, stream_error},
};
struct MockWs {
    receiver: Arc<Mutex<mpsc::Receiver<Result<StreamEvent, QueryError>>>>,
    connections: Arc<AtomicUsize>,
    ack: bool,
    pong: bool,
}
struct MockConnection {
    receiver: Arc<Mutex<mpsc::Receiver<Result<StreamEvent, QueryError>>>>,
    local: std::collections::VecDeque<StreamEvent>,
    pong: bool,
}
#[async_trait]
impl RealtimeSource for MockWs {
    async fn connect(&self, _: &str, _: usize) -> Result<Box<dyn RealtimeConnection>, QueryError> {
        self.connections.fetch_add(1, Ordering::SeqCst);
        let mut local = std::collections::VecDeque::new();
        if self.ack {
            local.push_back(StreamEvent::Subscribed);
        }
        Ok(Box::new(MockConnection {
            receiver: self.receiver.clone(),
            local,
            pong: self.pong,
        }))
    }
}
#[async_trait]
impl RealtimeConnection for MockConnection {
    async fn next(&mut self) -> Result<StreamEvent, QueryError> {
        if let Some(event) = self.local.pop_front() {
            return Ok(event);
        }
        self.receiver
            .lock()
            .await
            .recv()
            .await
            .unwrap_or_else(|| Err(stream_error("DISCONNECT", "fixture closed")))
    }
    async fn ping(&mut self) -> Result<(), QueryError> {
        if self.pong {
            self.local.push_back(StreamEvent::Pong);
        }
        Ok(())
    }
    async fn close(&mut self) {}
}
fn mock(ack: bool, pong: bool) -> (Arc<MockWs>, mpsc::Sender<Result<StreamEvent, QueryError>>) {
    let (sender, receiver) = mpsc::channel(256);
    (
        Arc::new(MockWs {
            receiver: Arc::new(Mutex::new(receiver)),
            connections: Arc::new(AtomicUsize::new(0)),
            ack,
            pong,
        }),
        sender,
    )
}
async fn wait_for(
    db: &Postgres,
    condition: impl Fn(&trade_log::checkpoint::CollectionStatus) -> bool,
) {
    tokio::time::timeout(Duration::from_secs(12), async {
        loop {
            if let Ok(list) = db.collection_status(ACCOUNT).await
                && condition(&list.items[0])
            {
                break;
            }
            tokio::time::sleep(Duration::from_millis(20)).await;
        }
    })
    .await
    .unwrap();
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn automatic_dual_channel_reconnects_compensates_and_heartbeats_a_quiet_account() {
    let db = database().await;
    let now = chrono::Utc::now().timestamp_millis();
    let mut c = collection_config();
    c.start_time = trade_log::collection::time(now - 5000);
    c.safety_delay_seconds = 0;
    c.interval_seconds = 1;
    let mut history = fills(2);
    history[0]["time"] = json!(now - 1000);
    history[1]["time"] = json!(now - 500);
    let rt = runtime(&db, c, history.clone());
    let (ws, sender) = mock(true, true);
    let cancel = CancellationToken::new();
    let worker_cancel = cancel.clone();
    let source = ws.clone();
    let cfg = WebsocketConfig {
        enabled: true,
        ping_interval_seconds: 1,
        pong_timeout_seconds: 1,
        reconnect_jitter_percent: 0,
        ..Default::default()
    };
    let worker = tokio::spawn(async move {
        rt.run_with_realtime(worker_cancel, cfg, source).await;
    });
    wait_for(&db, |s| s.websocket["status"] == "LIVE").await;
    let before = Instant::now();
    sender
        .send(Ok(StreamEvent::Data(message(history.clone(), Some(false)))))
        .await
        .unwrap();
    wait_for(&db, |s| s.websocket["last_committed_at"].is_string()).await;
    assert!(before.elapsed() < Duration::from_secs(5));
    wait_for(&db, |s| s.monitoring_status == "LIVE").await;
    let original_session =
        db.collection_status(ACCOUNT).await.unwrap().items[0].websocket["session_id"].clone();
    let pong_before =
        db.collection_status(ACCOUNT).await.unwrap().items[0].websocket["last_pong_at"].clone();
    wait_for(&db, |s| s.websocket["last_pong_at"] != pong_before).await;
    assert_eq!(ws.connections.load(Ordering::SeqCst), 1);
    sender
        .send(Err(stream_error("DISCONNECT", "fixture disconnect")))
        .await
        .unwrap();
    wait_for(&db, |s| s.websocket["status"] == "RECONNECT_WAIT").await;
    assert_eq!(ws.connections.load(Ordering::SeqCst), 1);
    wait_for(&db, |s| {
        s.websocket["session_id"] != original_session && s.websocket["status"] == "LIVE"
    })
    .await;
    sender
        .send(Ok(StreamEvent::Data(message(history, Some(true)))))
        .await
        .unwrap();
    wait_for(&db, |s| {
        s.monitoring_status == "LIVE" && s.websocket["reconnect_count"] == 1
    })
    .await;
    assert_eq!(count(&db, "account_facts_current").await, 2);
    assert!(count(&db, "fact_observations").await >= 4);
    cancel.cancel();
    tokio::time::timeout(Duration::from_secs(5), worker)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].status,
        "STOPPED"
    );
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn missing_ack_and_missing_pong_record_recovery_instead_of_claiming_live() {
    for (ack, pong, code) in [
        (false, true, "SUBSCRIBE_TIMEOUT"),
        (true, false, "HEARTBEAT_TIMEOUT"),
    ] {
        let db = database().await;
        let now = chrono::Utc::now().timestamp_millis();
        let mut c = collection_config();
        c.start_time = trade_log::collection::time(now - 5000);
        c.safety_delay_seconds = 0;
        c.interval_seconds = 1;
        let rt = runtime(&db, c, json!([]));
        let (ws, _sender) = mock(ack, pong);
        let cancel = CancellationToken::new();
        let wc = cancel.clone();
        let cfg = WebsocketConfig {
            enabled: true,
            subscribe_timeout_seconds: 1,
            ping_interval_seconds: 1,
            pong_timeout_seconds: 1,
            ..Default::default()
        };
        let worker = tokio::spawn(async move {
            rt.run_with_realtime(wc, cfg, ws).await;
        });
        wait_for(&db, |s| {
            s.websocket["last_error"]["code"] == code && s.websocket["status"] == "RECONNECT_WAIT"
        })
        .await;
        assert_ne!(
            db.collection_status(ACCOUNT).await.unwrap().items[0].monitoring_status,
            "LIVE"
        );
        let jobs: i64 = sqlx::query_scalar(
            "SELECT count(*) FROM trade_log.collection_jobs WHERE transport='WEBSOCKET'",
        )
        .fetch_one(&db.pool)
        .await
        .unwrap();
        assert_eq!(jobs, 0);
        cancel.cancel();
        tokio::time::timeout(Duration::from_secs(5), worker)
            .await
            .unwrap()
            .unwrap();
    }
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn bounded_queue_overflow_closes_socket_and_http_recovers_the_missing_fills() {
    let db = database().await;
    let now = chrono::Utc::now().timestamp_millis();
    let mut c = collection_config();
    c.start_time = trade_log::collection::time(now - 5000);
    c.safety_delay_seconds = 0;
    c.interval_seconds = 1;
    let mut history = fills(5);
    for (i, f) in history.as_array_mut().unwrap().iter_mut().enumerate() {
        f["time"] = json!(now + 1000 + i as i64);
    }
    let rt = runtime(&db, c, history.clone());
    let (ws, sender) = mock(true, true);
    let cancel = CancellationToken::new();
    let wc = cancel.clone();
    let cfg = WebsocketConfig {
        enabled: true,
        max_pending_messages: 1,
        ping_interval_seconds: 1,
        pong_timeout_seconds: 1,
        ..Default::default()
    };
    let worker = tokio::spawn(async move {
        rt.run_with_realtime(wc, cfg, ws).await;
    });
    wait_for(&db, |s| s.websocket["status"] == "LIVE").await;
    let mut held = db.pool.begin().await.unwrap();
    sqlx::query("SELECT committed_seq FROM trade_log.ingestion_state WHERE id=1 FOR UPDATE")
        .execute(&mut *held)
        .await
        .unwrap();

    for f in history.as_array().unwrap() {
        sender
            .send(Ok(StreamEvent::Data(message(json!([f]), Some(false)))))
            .await
            .unwrap();
    }
    wait_for(&db, |s| {
        s.websocket["last_error"]["code"] == "QUEUE_OVERFLOW"
    })
    .await;
    assert_eq!(count(&db, "account_facts_current").await, 0);
    held.rollback().await.unwrap();
    wait_for(&db, |s| s.monitoring_status == "LIVE").await;
    assert_eq!(count(&db, "account_facts_current").await, 5);
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].recovery["open_gap_count"],
        0
    );
    cancel.cancel();
    tokio::time::timeout(Duration::from_secs(5), worker)
        .await
        .unwrap()
        .unwrap();
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn successful_http_scan_clears_legacy_database_block_but_not_data_conflict() {
    let db = database().await;
    let (c, l, session) = start(&db).await;
    let error = trade_log::query::QueryError::unavailable("Legacy database timeout");
    db.end_stream(&l, &c, session, &error, false).await.unwrap();
    db.stream_subscribed(&l, session).await.unwrap();
    let end: i64 = sqlx::query_scalar("SELECT end_ms FROM trade_log.collection_gaps LIMIT 1")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    let mut tx = db.pool.begin().await.unwrap();
    trade_log_adapters::postgres::realtime::http_recovered(&mut tx, &l.key, end - 1)
        .await
        .unwrap();
    tx.commit().await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].recovery["status"],
        "BLOCKED"
    );
    let mut tx = db.pool.begin().await.unwrap();
    trade_log_adapters::postgres::realtime::http_recovered(&mut tx, &l.key, end)
        .await
        .unwrap();
    tx.commit().await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].recovery["open_gap_count"],
        0
    );
    sqlx::query("UPDATE trade_log.collection_gaps SET status='BLOCKED',last_error='{\"code\":\"VERSION_CONFLICT\"}'::jsonb")
        .execute(&db.pool).await.unwrap();
    let mut tx = db.pool.begin().await.unwrap();
    trade_log_adapters::postgres::realtime::http_recovered(&mut tx, &l.key, end + 1)
        .await
        .unwrap();
    tx.commit().await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].recovery["status"],
        "BLOCKED"
    );
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn renewal_waiting_on_http_lock_does_not_stall_http_commit() {
    let db = database().await;
    sqlx::query("CREATE FUNCTION delay_checkpoint_commit() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.position IS DISTINCT FROM OLD.position THEN PERFORM pg_sleep(3); END IF; RETURN NEW; END $$")
        .execute(&db.pool).await.unwrap();
    sqlx::query("CREATE TRIGGER delay_commit BEFORE UPDATE ON trade_log.collection_checkpoints FOR EACH ROW EXECUTE FUNCTION delay_checkpoint_commit()")
        .execute(&db.pool).await.unwrap();
    let mut c = collection_config();
    c.start_time = trade_log::collection::time(chrono::Utc::now().timestamp_millis() - 5000);
    c.lease_seconds = 6;
    c.round_timeout_seconds = 20;
    c.safety_delay_seconds = 0;
    let rt = runtime(&db, c, json!([]));
    let cancel = CancellationToken::new();
    let wc = cancel.clone();
    let worker = tokio::spawn(async move { rt.run(wc).await });
    let observed = tokio::time::timeout(Duration::from_secs(7), async {
        loop {
            let completed: i64 = sqlx::query_scalar("SELECT count(*) FROM trade_log.collection_jobs WHERE transport='HTTP' AND status='COMPLETED'")
                .fetch_one(&db.pool).await.unwrap();
            if completed > 0 { break; }
            tokio::time::sleep(Duration::from_millis(50)).await;
        }
    }).await;
    cancel.cancel();
    tokio::time::timeout(Duration::from_secs(15), worker)
        .await
        .unwrap()
        .unwrap();
    observed.expect("HTTP must commit while renewal waits for its checkpoint lock");
    let epoch: i64 = sqlx::query_scalar("SELECT lease_epoch FROM trade_log.collection_checkpoints")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    assert_eq!(epoch, 1);
}
