mod common;
use common::*;
use hyperliquid::parser::HyperliquidParser;
use serde_json::{Value, json};
use trade_log::{
    checkpoint::Lease,
    collection::CollectionConfig,
    publishing::CandidatePolicy,
    realtime::{StreamEvent, StreamMessage},
};
use trade_log_adapters::{
    postgres::{Postgres, realtime::StreamCommit},
    webhook::{Webhook, classify},
};
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

fn policy() -> CandidatePolicy {
    CandidatePolicy {
        max_event_age_seconds: 30,
        signal_ttl_seconds: 60,
        clock_skew_tolerance_seconds: 5,
        version: "candidate-v1".into(),
    }
}
fn fresh(tid: u64) -> Value {
    let mut f = fills(1);
    f[0]["time"] = json!(chrono::Utc::now().timestamp_millis());
    f[0]["tid"] = json!(tid);
    f
}
async fn send(
    db: &Postgres,
    l: &Lease,
    session: Uuid,
    seq: i64,
    f: Value,
    mode: Option<bool>,
) -> String {
    let q = db
        .archive_stream(
            l,
            session,
            seq,
            &message(f, mode),
            META.as_bytes(),
            SPOT.as_bytes(),
        )
        .await
        .unwrap();
    commit(db, l, session, seq, &q).await.unwrap();
    q
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL"]
async fn candidate_atomicity_http_race_snapshot_and_fencing() {
    let db = database().await;
    let (_, l, session) = start(&db).await;
    db.configure_publishing(&l.key, "https://example.invalid/events", true, &policy())
        .await
        .unwrap();
    let activated: chrono::DateTime<chrono::Utc> =
        sqlx::query_scalar("SELECT activated_at FROM trade_log.publishing_control")
            .fetch_one(&db.pool)
            .await
            .unwrap();
    db.configure_publishing(&l.key, "https://example.invalid/events", true, &policy())
        .await
        .unwrap();
    assert_eq!(
        activated,
        sqlx::query_scalar::<_, chrono::DateTime<chrono::Utc>>(
            "SELECT activated_at FROM trade_log.publishing_control"
        )
        .fetch_one(&db.pool)
        .await
        .unwrap()
    );
    assert_eq!(
        db.configure_publishing(&l.key, "https://other.invalid/events", true, &policy())
            .await
            .unwrap_err()
            .code,
        "TARGET_CONFIGURATION_CHANGED"
    );
    send(&db, &l, session, 1, fresh(1), None).await;
    assert_eq!(count(&db, "outbox_events").await, 0);
    send(&db, &l, session, 2, fresh(2), Some(true)).await;
    assert_eq!(count(&db, "outbox_events").await, 0);
    let duplicate = fresh(3);
    let original = collect(&db, "query_http_first", duplicate.clone(), 100)
        .await
        .unwrap();
    send(&db, &l, session, 3, duplicate.clone(), None).await;
    send(&db, &l, session, 4, duplicate, None).await;
    assert_eq!(count(&db, "outbox_events").await, 1);
    let payload: Value = sqlx::query_scalar("SELECT payload FROM trade_log.outbox_events")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    assert_eq!(payload["fact"], json!(original.trades[0]));
    assert_eq!(
        payload["observation"]["realtime_reason"],
        "POST_SNAPSHOT_UNFLAGGED"
    );
    assert_ne!(
        payload["observation"]["raw_log_id"],
        payload["fact"]["raw_log_id"]
    );
    let malformed = message(fresh(4), None);
    let mut value: Value = serde_json::from_slice(&malformed.body).unwrap();
    value["data"]["isSnapshot"] = Value::Null;
    let malformed = StreamMessage {
        body: serde_json::to_vec(&value).unwrap(),
        ..malformed
    };
    let q = db
        .archive_stream(&l, session, 5, &malformed, META.as_bytes(), SPOT.as_bytes())
        .await
        .unwrap();
    commit(&db, &l, session, 5, &q).await.unwrap();
    assert_eq!(count(&db, "outbox_events").await, 1);
    // Force an outbox failure; no fact, decision or message commit may survive.
    sqlx::raw_sql("CREATE FUNCTION fail_outbox() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'test rollback'; END $$; CREATE TRIGGER fail_outbox BEFORE INSERT ON trade_log.outbox_events FOR EACH ROW EXECUTE FUNCTION fail_outbox();").execute(&db.pool).await.unwrap();
    let q = archive(&db, &l, session, 6, fresh(6)).await;
    let before = count(&db, "account_fact_versions").await;
    assert!(commit(&db, &l, session, 6, &q).await.is_err());
    assert_eq!(count(&db, "account_fact_versions").await, before);
    sqlx::query("DROP TRIGGER fail_outbox ON trade_log.outbox_events")
        .execute(&db.pool)
        .await
        .unwrap();
    commit(&db, &l, session, 6, &q).await.unwrap();
    let first = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    sqlx::query("UPDATE trade_log.outbox_events SET lease_expires_at=now()-interval '1 second' WHERE event_id=$1").bind(&first.event_id).execute(&db.pool).await.unwrap();
    let second = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(first.event_id, second.event_id);
    assert_eq!(first.body, second.body);
    assert_eq!(second.attempt, 2);
    assert_eq!(
        db.finish_publication(&first, "DELIVERED", Some(200), None, 0)
            .await
            .unwrap_err()
            .code,
        "LEASE_LOST"
    );
    db.finish_publication(
        &second,
        "BLOCKED",
        Some(401),
        Some("HTTP_PERMANENT_REJECTION"),
        0,
    )
    .await
    .unwrap();
    db.retry_publication(&l.key, &second.event_id)
        .await
        .unwrap();
    let third = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(third.body, first.body);
    assert_eq!(third.attempt, 3);
    db.finish_publication(&third, "DELIVERED", Some(200), None, 0)
        .await
        .unwrap();
    db.configure_publishing(&l.key, "", false, &policy())
        .await
        .unwrap();
    assert!(
        db.claim_publication(&l.key, Uuid::new_v4(), 30)
            .await
            .unwrap()
            .is_none()
    );
    db.configure_publishing(&l.key, "https://example.invalid/events", true, &policy())
        .await
        .unwrap();
    let status = db.publishing_status(&l.key).await.unwrap();
    assert_eq!(status["activation_epoch"], 2);
}

#[test]
fn responses_and_retry_after() {
    let now = chrono::Utc::now();
    for s in [200, 201, 204, 299] {
        assert_eq!(classify(Some(s), None, 1, 5, 300, now).status, "DELIVERED");
    }
    for s in [408, 425, 429, 500, 503, 599] {
        assert_eq!(classify(Some(s), None, 1, 5, 300, now).status, "RETRY_WAIT");
    }
    for s in [301, 302, 400, 401, 403, 404, 409] {
        assert_eq!(classify(Some(s), None, 1, 5, 300, now).status, "BLOCKED");
    }
    assert_eq!(classify(None, None, 1, 5, 300, now).status, "RETRY_WAIT");
    assert_eq!(classify(Some(429), Some("120"), 1, 5, 300, now).delay, 120);
    assert_eq!(
        classify(Some(503), Some("86401"), 1, 5, 300, now).status,
        "BLOCKED"
    );
    assert_eq!(
        classify(
            Some(429),
            Some("99999999999999999999999999999"),
            1,
            5,
            300,
            now
        )
        .status,
        "BLOCKED"
    );
    assert!(classify(Some(500), None, 100, 5, 300, now).delay <= 300);
    let date = (now + chrono::Duration::seconds(120)).to_rfc2822();
    assert!(classify(Some(503), Some(&date), 1, 5, 300, now).delay >= 119);
}

struct ReceiverProcess {
    child: std::process::Child,
    directory: std::path::PathBuf,
}
impl Drop for ReceiverProcess {
    fn drop(&mut self) {
        let _ = self.child.kill();
        let _ = self.child.wait();
        let _ = std::fs::remove_dir_all(&self.directory);
    }
}
async fn receiver() -> (ReceiverProcess, String, String) {
    let directory =
        std::env::temp_dir().join(format!("robotech_receiver_{}", Uuid::new_v4().simple()));
    std::fs::create_dir(&directory).unwrap();
    let token = format!("{}{}", Uuid::new_v4().simple(), Uuid::new_v4().simple());
    std::fs::write(directory.join("token"), &token).unwrap();
    let socket = std::net::TcpListener::bind("127.0.0.1:0").unwrap();
    let port = socket.local_addr().unwrap().port();
    drop(socket);
    let script = std::env::var("ROBOTECH_RECEIVER_SCRIPT").unwrap_or_else(|_| {
        format!(
            "{}/../../../scripts/webhook-receiver.py",
            env!("CARGO_MANIFEST_DIR")
        )
    });
    let child = std::process::Command::new("python3")
        .arg(script)
        .args([
            "--host",
            "127.0.0.1",
            "--port",
            &port.to_string(),
            "--database",
        ])
        .arg(directory.join("receipts.sqlite"))
        .arg("--credential-file")
        .arg(directory.join("token"))
        .spawn()
        .unwrap();
    let process = ReceiverProcess { child, directory };
    let url = format!("http://127.0.0.1:{port}");
    let client = reqwest::Client::new();
    for _ in 0..100 {
        if client.get(format!("{url}/health")).send().await.is_ok() {
            return (process, url, token);
        }
        tokio::time::sleep(std::time::Duration::from_millis(20)).await;
    }
    panic!("receiver failed to start")
}
async fn control(url: &str, value: Value) {
    assert_eq!(
        reqwest::Client::new()
            .post(format!("{url}/control"))
            .json(&value)
            .send()
            .await
            .unwrap()
            .status(),
        200
    );
}
async fn receipts(url: &str) -> Value {
    reqwest::Client::new()
        .get(format!("{url}/receipts"))
        .send()
        .await
        .unwrap()
        .json()
        .await
        .unwrap()
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL and Python3; deterministic server acceptance suite"]
async fn server_fixture_suite() {
    let db = database().await;
    let (_receiver, url, token) = receiver().await;
    let mut config = collection_config();
    config.lease_seconds = 120;
    let l = db
        .acquire_collection(&config, Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    db.prepare_stream(&l, &config, true).await.unwrap();
    let session = Uuid::new_v4();
    db.begin_stream(&l, &config, session).await.unwrap();
    db.stream_connected(&l, session).await.unwrap();
    db.stream_subscribed(&l, session).await.unwrap();
    db.configure_publishing(&l.key, &format!("{url}/events"), true, &policy())
        .await
        .unwrap();
    send(&db, &l, session, 1, fresh(10), Some(true)).await;
    let duplicate = fresh(11);
    let http_first = collect(&db, "query_fixture_http_first", duplicate.clone(), 100)
        .await
        .unwrap();
    send(&db, &l, session, 2, duplicate, None).await;
    let envelope: Value = sqlx::query_scalar("SELECT payload FROM trade_log.outbox_events")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    assert_eq!(envelope["fact"], json!(http_first.trades[0]));
    assert_ne!(
        envelope["fact"]["raw_log_id"],
        envelope["observation"]["raw_log_id"]
    );
    let webhook = Webhook::new(&format!("{url}/events"), &token, 2, 65536).unwrap();
    control(&url, json!({"response_status":500})).await;
    let first = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    let d = webhook.send(&first, 1, 4).await;
    assert_eq!(d.status, "RETRY_WAIT");
    db.finish_publication(&first, d.status, d.http, d.error, d.delay)
        .await
        .unwrap();
    // Collector can commit fresh facts while receiver is failing; no external send owns ingestion locks.
    let q = archive(&db, &l, session, 3, fresh(12)).await;
    sqlx::raw_sql("CREATE FUNCTION fail_fixture_outbox() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'test atomic rollback'; END $$; CREATE TRIGGER fail_fixture_outbox BEFORE INSERT ON trade_log.outbox_events FOR EACH ROW EXECUTE FUNCTION fail_fixture_outbox();").execute(&db.pool).await.unwrap();
    let facts_before = count(&db, "account_fact_versions").await;
    let decisions_before = count(&db, "publication_decisions").await;
    assert!(commit(&db, &l, session, 3, &q).await.is_err());
    assert_eq!(count(&db, "account_fact_versions").await, facts_before);
    assert_eq!(count(&db, "publication_decisions").await, decisions_before);
    sqlx::query("DROP TRIGGER fail_fixture_outbox ON trade_log.outbox_events")
        .execute(&db.pool)
        .await
        .unwrap();
    commit(&db, &l, session, 3, &q).await.unwrap();
    assert_eq!(count(&db, "outbox_events").await, 2);
    control(
        &url,
        json!({"response_status":200,"drop_response_after_store_once":true}),
    )
    .await;
    sqlx::query("UPDATE trade_log.outbox_events SET next_retry_at=now() WHERE event_id=$1")
        .bind(&first.event_id)
        .execute(&db.pool)
        .await
        .unwrap();
    let second = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(first.body, second.body);
    let d = webhook.send(&second, 1, 4).await;
    assert_eq!(d.status, "RETRY_WAIT");
    db.finish_publication(&second, d.status, d.http, d.error, 0)
        .await
        .unwrap();
    let third = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    let d = webhook.send(&third, 1, 4).await;
    assert_eq!(d.status, "DELIVERED");
    db.finish_publication(&third, d.status, d.http, d.error, 0)
        .await
        .unwrap();
    assert_eq!(first.body, third.body);
    let r = receipts(&url).await;
    assert_eq!(r["receipts"].as_array().unwrap().len(), 1);
    assert_eq!(r["attempts"].as_array().unwrap().len(), 3);
    assert_eq!(r["attempts"][2]["duplicate"], 1);
    control(&url, json!({"response_status":401})).await;
    let blocked = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    let d = webhook.send(&blocked, 1, 4).await;
    assert_eq!(d.status, "BLOCKED");
    db.finish_publication(&blocked, d.status, d.http, d.error, 0)
        .await
        .unwrap();
    assert!(
        db.claim_publication(&l.key, Uuid::new_v4(), 30)
            .await
            .unwrap()
            .is_none()
    );
    control(&url, json!({"response_status":429,"retry_after":"3"})).await;
    db.retry_publication(&l.key, &blocked.event_id)
        .await
        .unwrap();
    let retry = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    let d = webhook.send(&retry, 1, 4).await;
    assert_eq!(d.status, "RETRY_WAIT");
    assert_eq!(d.delay, 3);
    db.finish_publication(&retry, d.status, d.http, d.error, d.delay)
        .await
        .unwrap();
    assert!(
        db.claim_publication(&l.key, Uuid::new_v4(), 30)
            .await
            .unwrap()
            .is_none()
    );
    // Simulate process interruption after claim, then another owner reclaims unchanged bytes.
    sqlx::query("UPDATE trade_log.outbox_events SET next_retry_at=now() WHERE event_id=$1")
        .bind(&retry.event_id)
        .execute(&db.pool)
        .await
        .unwrap();
    let interrupted = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    sqlx::query("UPDATE trade_log.outbox_events SET lease_expires_at=now()-interval '1 second' WHERE event_id=$1").bind(&retry.event_id).execute(&db.pool).await.unwrap();
    control(&url, json!({"response_status":200,"retry_after":null})).await;
    let restored = db
        .claim_publication(&l.key, Uuid::new_v4(), 30)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(restored.body, interrupted.body);
    assert!(
        db.finish_publication(&interrupted, "DELIVERED", Some(200), None, 0)
            .await
            .is_err()
    );
    let d = webhook.send(&restored, 1, 4).await;
    db.finish_publication(&restored, d.status, d.http, d.error, 0)
        .await
        .unwrap();
    assert_eq!(
        db.publishing_status(&l.key).await.unwrap()["outbox"]["DELIVERED"],
        2
    );
    // Session boundary, pre-activation, stale, forced and metadata gates also exercised in the transaction.
    let other = Uuid::new_v4();
    db.begin_stream(&l, &config, other).await.unwrap();
    db.stream_connected(&l, other).await.unwrap();
    db.stream_subscribed(&l, other).await.unwrap();
    send(&db, &l, other, 1, fresh(20), None).await;
    assert_eq!(count(&db, "outbox_events").await, 2);
    let mut old = fresh(21);
    old[0]["time"] = json!(chrono::Utc::now().timestamp_millis() - 120000);
    send(&db, &l, other, 2, old, Some(false)).await;
    assert_eq!(count(&db, "outbox_events").await, 2);
    sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || '{\"metadata_stale\":true}'::jsonb").execute(&db.pool).await.unwrap();
    send(&db, &l, other, 3, fresh(22), Some(false)).await;
    assert_eq!(count(&db, "outbox_events").await, 2);
    sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || '{\"metadata_stale\":false}'::jsonb").execute(&db.pool).await.unwrap();
    send(&db, &l, other, 4, fresh(23), Some(true)).await;
    let mut forced = fresh(24);
    forced[0]["liquidation"] = json!({"method":"market"});
    send(&db, &l, other, 5, forced, Some(false)).await;
    let mut future = fresh(25);
    future[0]["time"] = json!(chrono::Utc::now().timestamp_millis() + 100000);
    send(&db, &l, other, 6, future, Some(false)).await;
    // Shift only this isolated test's cutover to distinguish stale from pre-activation.
    sqlx::query("UPDATE trade_log.publishing_control SET activated_at=now()-interval '1 minute'")
        .execute(&db.pool)
        .await
        .unwrap();
    let mut stale = fresh(26);
    stale[0]["time"] = json!(chrono::Utc::now().timestamp_millis() - 31000);
    send(&db, &l, other, 7, stale, Some(false)).await;
    let mut invalid = message(fresh(27), None);
    let mut body: Value = serde_json::from_slice(&invalid.body).unwrap();
    body["data"]["isSnapshot"] = json!("false");
    invalid.body = serde_json::to_vec(&body).unwrap();
    let q = db
        .archive_stream(&l, other, 8, &invalid, META.as_bytes(), SPOT.as_bytes())
        .await
        .unwrap();
    commit(&db, &l, other, 8, &q).await.unwrap();
    let mut premature = message(fresh(28), Some(false));
    premature.received_at = (chrono::Utc::now() - chrono::Duration::seconds(120)).to_rfc3339();
    let q = db
        .archive_stream(&l, other, 9, &premature, META.as_bytes(), SPOT.as_bytes())
        .await
        .unwrap();
    commit(&db, &l, other, 9, &q).await.unwrap();
    db.configure_publishing(&l.key, "", false, &policy())
        .await
        .unwrap();
    send(&db, &l, other, 10, fresh(29), Some(false)).await;
    assert_eq!(count(&db, "outbox_events").await, 2);
    let reasons: Vec<String> = sqlx::query_scalar(
        "SELECT DISTINCT reason FROM trade_log.publication_decisions WHERE result='SUPPRESSED'",
    )
    .fetch_all(&db.pool)
    .await
    .unwrap();
    for reason in [
        "SNAPSHOT",
        "MODE_UNCONFIRMED",
        "BEFORE_ACTIVATION",
        "STALE_EVENT",
        "FUTURE_EVENT",
        "NOT_COPY_ELIGIBLE",
        "METADATA_STALE",
        "SESSION_UNCONFIRMED",
        "DISABLED",
    ] {
        assert!(reasons.iter().any(|s| s == reason), "missing {reason}");
    }
    let events = receipts(&url).await;
    assert_eq!(events["receipts"].as_array().unwrap().len(), 2);
    executable_restart(&db, &l, other, &url, &token, &_receiver.directory).await;
    println!(
        "[PASS] fixed candidate/suppression gates, durable outbox, 500 recovery, 429 Retry-After, 401 explicit retry, lost acknowledgement, fenced reclaim, immutable bytes and receiver deduplication"
    );
}

struct ChildProcess(std::process::Child);
impl ChildProcess {
    fn stop(&mut self) {
        assert!(
            std::process::Command::new("kill")
                .args(["-TERM", &self.0.id().to_string()])
                .status()
                .unwrap()
                .success()
        );
        assert!(self.0.wait().unwrap().success());
    }
}
impl Drop for ChildProcess {
    fn drop(&mut self) {
        let _ = self.0.kill();
        let _ = self.0.wait();
    }
}
async fn wait_outbox(db: &Postgres, status: &str, expected: i64) {
    for _ in 0..200 {
        let count: i64 =
            sqlx::query_scalar("SELECT count(*) FROM trade_log.outbox_events WHERE status=$1")
                .bind(status)
                .fetch_one(&db.pool)
                .await
                .unwrap();
        if count == expected {
            return;
        }
        tokio::time::sleep(std::time::Duration::from_millis(50)).await;
    }
    panic!("outbox did not reach {status}={expected}");
}
async fn executable_restart(
    db: &Postgres,
    l: &Lease,
    session: Uuid,
    url: &str,
    token: &str,
    directory: &std::path::Path,
) {
    let binary = std::env::var("ROBOTECH_PUBLISHER_BINARY").unwrap_or_else(|_| {
        format!(
            "{}/../../../target/debug/trade-parser-publisher",
            env!("CARGO_MANIFEST_DIR")
        )
    });
    let mut database_url =
        reqwest::Url::parse(&std::env::var("ROBOTECH_TEST_DATABASE_URL").unwrap()).unwrap();
    database_url.set_path(db.pool.connect_options().get_database().unwrap());
    let publisher_role: bool = sqlx::query_scalar(
        "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname='trade_log_publisher')",
    )
    .fetch_one(&db.pool)
    .await
    .unwrap();
    if publisher_role {
        let privileges:bool=sqlx::query_scalar("SELECT NOT has_table_privilege('trade_log_publisher','trade_log.account_fact_versions','UPDATE') AND NOT has_table_privilege('trade_log_publisher','trade_log.collection_checkpoints','UPDATE') AND has_table_privilege('trade_log_publisher','trade_log.outbox_events','UPDATE')").fetch_one(&db.pool).await.unwrap();
        assert!(
            privileges,
            "publisher role must not mutate facts or checkpoints"
        );
        database_url
            .query_pairs_mut()
            .append_pair("options", "-c role=trade_log_publisher");
    }

    std::fs::write(directory.join("publisher-db-url"), database_url.as_str()).unwrap();
    std::fs::write(
        directory.join("collector.toml"),
        format!("[hyperliquid]\nnetwork=\"mainnet\"\n[collection]\naccount=\"{ACCOUNT}\"\n"),
    )
    .unwrap();
    let socket = std::net::TcpListener::bind("127.0.0.1:0").unwrap();
    let port = socket.local_addr().unwrap().port();
    drop(socket);
    let config = include_str!("../../../../config/trade-parser-publisher.toml")
        .replace("host = \"0.0.0.0\"", "host = \"127.0.0.1\"")
        .replace("port = 8083", &format!("port = {port}"))
        .replace("enabled = false", "enabled = true")
        .replace(
            "webhook_url = \"\"",
            &format!("webhook_url = \"{url}/events\""),
        )
        .replace("allow_plain_http = false", "allow_plain_http = true")
        .replace(
            "/run/secrets/trade-log-token",
            directory.join("token").to_str().unwrap(),
        )
        .replace(
            "/run/secrets/webhook-token",
            directory.join("token").to_str().unwrap(),
        )
        .replace(
            "/run/secrets/trade-log-publisher-database-url",
            directory.join("publisher-db-url").to_str().unwrap(),
        )
        .replace(
            "/etc/robotech/trade-collector.toml",
            directory.join("collector.toml").to_str().unwrap(),
        );
    std::fs::write(directory.join("publisher.toml"), config).unwrap();
    let spawn = || {
        ChildProcess(
            std::process::Command::new(&binary)
                .arg("--config")
                .arg(directory.join("publisher.toml"))
                .spawn()
                .unwrap(),
        )
    };
    let mut process = spawn();
    let client = reqwest::Client::new();
    let status_url = format!("http://127.0.0.1:{port}/internal/v1/publishing-status");
    let baseline = wait_publisher(&client, &status_url, token).await;
    assert_eq!(baseline["data"]["status"], "RUNNING");
    send(db, l, session, 11, fresh(30), Some(false)).await;
    wait_outbox(db, "DELIVERED", 3).await;
    control(url, json!({"response_status":500})).await;
    send(db, l, session, 12, fresh(31), Some(false)).await;
    wait_outbox(db, "RETRY_WAIT", 1).await;
    let queued: Value =
        sqlx::query_scalar("SELECT payload FROM trade_log.outbox_events WHERE status='RETRY_WAIT'")
            .fetch_one(&db.pool)
            .await
            .unwrap();
    process.stop();
    let mut restarted = spawn();
    let after = wait_publisher(&client, &status_url, token).await;
    assert_eq!(
        after["data"]["activated_at"],
        baseline["data"]["activated_at"]
    );
    assert_eq!(
        after["data"]["activation_epoch"],
        baseline["data"]["activation_epoch"]
    );
    control(url, json!({"response_status":200})).await;
    wait_outbox(db, "DELIVERED", 4).await;
    let delivered: Value =
        sqlx::query_scalar("SELECT payload FROM trade_log.outbox_events WHERE event_id=$1")
            .bind(queued["event_id"].as_str().unwrap())
            .fetch_one(&db.pool)
            .await
            .unwrap();
    assert_eq!(queued, delivered);
    control(url, json!({"response_status":401})).await;
    send(db, l, session, 13, fresh(32), Some(false)).await;
    wait_outbox(db, "BLOCKED", 1).await;
    let event: String =
        sqlx::query_scalar("SELECT event_id FROM trade_log.outbox_events WHERE status='BLOCKED'")
            .fetch_one(&db.pool)
            .await
            .unwrap();
    control(url, json!({"response_status":200})).await;
    assert!(
        std::process::Command::new(&binary)
            .arg("--config")
            .arg(directory.join("publisher.toml"))
            .args(["retry", "--event-id", &event])
            .status()
            .unwrap()
            .success()
    );
    wait_outbox(db, "DELIVERED", 5).await;
    restarted.stop();
    println!(
        "[PASS] actual publisher executable: durable delivery, 500 recovery after process restart, preserved activation/body/TTL, 401 and retry CLI"
    );
}
async fn wait_publisher(client: &reqwest::Client, url: &str, token: &str) -> Value {
    for _ in 0..100 {
        if let Ok(response) = client.get(url).bearer_auth(token).send().await
            && response.status() == 200
        {
            return response.json().await.unwrap();
        }
        tokio::time::sleep(std::time::Duration::from_millis(50)).await;
    }
    panic!("publisher startup timed out");
}
