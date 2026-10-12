mod common;
use common::*;
use serde_json::{Value, json};
use std::{
    sync::{
        Arc,
        atomic::{AtomicBool, Ordering},
    },
    time::Duration,
};
use tokio_util::sync::CancellationToken;
use trade_log::{
    acquisition::{SourceReader, SourceResponse},
    checkpoint::CollectionStatusReader,
    query::QueryError,
    watch::{AddAccount, BudgetConfig, ListParams, MonitoringConfig},
};
use trade_log_adapters::{collector_supervisor::CollectorSupervisor, postgres::Postgres};
use uuid::Uuid;
const B: &str = "0x0000000000000000000000000000000000000002";
fn add(account: &str) -> AddAccount {
    AddAccount {
        account: account.into(),
        start_time: trade_log::collection::time(chrono::Utc::now().timestamp_millis() - 5000),
    }
}
fn id(v: &Value) -> Uuid {
    Uuid::parse_str(v["id"].as_str().unwrap()).unwrap()
}
async fn start(db: &Postgres, account: &str) -> trade_log::checkpoint::Lease {
    let mut c = collection_config();
    c.account = account.into();
    c.start_time = add(account).start_time;
    db.acquire_collection(&c, Uuid::new_v4())
        .await
        .unwrap()
        .unwrap()
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL"]
async fn server_multi_account_suite() {
    let db = database().await;
    let m = MonitoringConfig {
        max_active_accounts: 2,
        ..Default::default()
    };
    db.monitoring_configure(&m, &BudgetConfig::default())
        .await
        .unwrap();
    let a = db.add_watch(add(ACCOUNT)).await.unwrap();
    let b = db.add_watch(add(B)).await.unwrap();
    assert_eq!(
        db.add_watch(add(ACCOUNT)).await.unwrap_err().code,
        "WATCH_ACCOUNT_EXISTS"
    );
    assert_eq!(
        db.add_watch(add("0x0000000000000000000000000000000000000003"))
            .await
            .unwrap_err()
            .code,
        "WATCH_ACCOUNT_LIMIT_REACHED"
    );
    let page = db
        .watch_list(ListParams {
            limit: Some(1),
            ..Default::default()
        })
        .await
        .unwrap();
    assert_eq!(page["items"].as_array().unwrap().len(), 1);
    let next = db
        .watch_list(ListParams {
            limit: Some(1),
            cursor: Some(page["next_cursor"].as_str().unwrap().into()),
            state: None,
        })
        .await
        .unwrap();
    assert_ne!(page["items"][0]["id"], next["items"][0]["id"]);
    let lease = start(&db, ACCOUNT).await;
    let paused = db.set_watch(id(&a), false).await.unwrap();
    assert_eq!(paused["generation"], 2);
    assert_eq!(db.set_watch(id(&a), false).await.unwrap()["generation"], 2);
    assert_eq!(
        db.renew_collection(&lease, 180).await.unwrap_err().code,
        "LEASE_LOST"
    );
    let mut c = collection_config();
    c.start_time = add(ACCOUNT).start_time;
    assert!(
        db.acquire_collection(&c, Uuid::new_v4())
            .await
            .unwrap()
            .is_none()
    );
    let (r1, r2) = tokio::join!(
        db.add_watch(add("0x0000000000000000000000000000000000000003")),
        db.add_watch(add("0x0000000000000000000000000000000000000004"))
    );
    assert_ne!(r1.is_ok(), r2.is_ok());
    let third = r1.or(r2).unwrap();
    db.set_watch(id(&third), false).await.unwrap();
    db.import_watch(ACCOUNT, &add(ACCOUNT).start_time)
        .await
        .unwrap();
    assert_eq!(
        db.watch_detail(id(&a)).await.unwrap()["desired_state"],
        "PAUSED"
    );
    let restored = db.set_watch(id(&a), true).await.unwrap();
    assert_eq!(restored["generation"], 3);
    assert_eq!(db.set_watch(id(&a), true).await.unwrap()["generation"], 3);
    let nextlease = start(&db, ACCOUNT).await;
    assert_eq!(nextlease.generation, Some(3));
    // A stale task cannot save a response or advance pending work after a pause/resume.
    let mut tx = db.pool.begin().await.unwrap();
    assert_eq!(
        db.lock_collection(&mut tx, &lease).await.unwrap_err().code,
        "LEASE_LOST"
    );
    drop(tx);
    let policy = trade_log::publishing::CandidatePolicy {
        max_event_age_seconds: 30,
        signal_ttl_seconds: 60,
        clock_skew_tolerance_seconds: 5,
        version: "candidate-v1".into(),
    };
    db.configure_publishing(
        &nextlease.key,
        "https://example.invalid/events",
        true,
        &policy,
    )
    .await
    .unwrap();
    let activated = db.publishing_status(&nextlease.key).await.unwrap();
    db.configure_publishing(
        &nextlease.key,
        "https://example.invalid/events",
        true,
        &policy,
    )
    .await
    .unwrap();
    assert_eq!(
        activated["activated_at"],
        db.publishing_status(&nextlease.key).await.unwrap()["activated_at"]
    );
    db.set_watch(id(&a), false).await.unwrap();
    assert_eq!(
        db.publishing_status(&nextlease.key).await.unwrap()["enabled"],
        false
    );
    db.set_watch(id(&a), true).await.unwrap();
    db.configure_publishing(
        &nextlease.key,
        "https://example.invalid/events",
        true,
        &policy,
    )
    .await
    .unwrap();
    assert!(
        db.publishing_status(&nextlease.key).await.unwrap()["activation_epoch"]
            .as_i64()
            .unwrap()
            > activated["activation_epoch"].as_i64().unwrap()
    );
    db.release_collection(&nextlease).await.unwrap();
    // Shared HTTP weight and in-flight slots across independently connected adapters.
    let reservation = db.reserve_source("HTTP", 120).await.unwrap().unwrap();
    let second = db.reserve_source("HTTP", 120).await.unwrap().unwrap();
    assert!(db.reserve_source("HTTP", 20).await.unwrap().is_none());
    db.settle_source(reservation, 20).await.unwrap();
    assert!(db.reserve_source("HTTP", 120).await.unwrap().is_some());
    db.release_source(second).await.unwrap();
    sqlx::query("DELETE FROM trade_log.source_budget_reservations")
        .execute(&db.pool)
        .await
        .unwrap();
    for _ in 0..15 {
        assert!(db.reserve_source("WS", 1).await.unwrap().is_some())
    }
    assert!(db.reserve_source("WS", 1).await.unwrap().is_none());
    sqlx::query("UPDATE trade_log.source_budget_reservations SET reserved_at=clock_timestamp()-interval '61 seconds'").execute(&db.pool).await.unwrap();
    assert!(db.reserve_source("WS", 1).await.unwrap().is_some());
    println!(
        "[PASS] capacity race, normalization, pagination, persistent pause, stale lease fencing, resume activation, shared HTTP/WS budgets"
    );
    db.set_watch(id(&b), true).await.unwrap();
    publishing_isolation(&db, id(&a), id(&b)).await;
    supervisor_isolation(&db, id(&a), id(&b)).await;
    roles(&db).await;
    management_http(&db).await;
}
struct FaultSource {
    panic_once: AtomicBool,
}
#[async_trait::async_trait]
impl SourceReader for FaultSource {
    async fn wait_retry(&self, _: u64) {}
    async fn fetch(
        &self,
        k: protocol_api::QueryKind,
        a: &str,
    ) -> Result<SourceResponse, QueryError> {
        Source(json!([])).fetch(k, a).await
    }
    async fn fetch_range(
        &self,
        a: &str,
        start: i64,
        end: i64,
    ) -> Result<SourceResponse, QueryError> {
        if a == ACCOUNT {
            if self.panic_once.swap(false, Ordering::SeqCst) {
                std::panic::resume_unwind(Box::new("injected account A panic"));
            }
            return Err(QueryError::unavailable("Injected A failure").with_retry());
        }
        RangeSource(json!([])).fetch_range(a, start, end).await
    }
}
async fn wait_success(db: &Postgres, account: &str, after: Option<String>) {
    tokio::time::timeout(Duration::from_secs(15), async {
        loop {
            if let Ok(c) = db.collection_status(account).await
                && c.items[0].last_success_at.is_some()
                && c.items[0].last_success_at != after
            {
                break;
            }
            tokio::time::sleep(Duration::from_millis(50)).await
        }
    })
    .await
    .unwrap();
}
async fn supervisor_isolation(db: &Postgres, a: Uuid, b: Uuid) {
    let mut defaults = collection_config();
    defaults.interval_seconds = 1;
    defaults.safety_delay_seconds = 0;
    let ws = trade_log::realtime::WebsocketConfig {
        enabled: false,
        ..Default::default()
    };
    let supervisor = CollectorSupervisor {
        store: db.clone(),
        defaults,
        monitoring: MonitoringConfig {
            reconcile_interval_seconds: 1,
            max_active_accounts: 2,
            ..Default::default()
        },
        websocket: ws,
        source: Arc::new(FaultSource {
            panic_once: AtomicBool::new(true),
        }),
        realtime: Arc::new(hyperliquid::websocket::WebsocketSource::new(&db.network)),
        parser: Arc::new(hyperliquid::parser::HyperliquidParser),
    };
    let cancel = CancellationToken::new();
    let task = tokio::spawn(supervisor.run(cancel.clone()));
    wait_success(db, B, None).await;
    tokio::time::timeout(Duration::from_secs(15), async {
        loop {
            if let Ok(status) = db.collection_status(ACCOUNT).await
                && status.items[0].consecutive_failures > 0
            {
                break;
            }
            tokio::time::sleep(Duration::from_millis(50)).await;
        }
    })
    .await
    .expect("account A worker must restart after panic and attempt collection again");
    assert!(!task.is_finished());
    db.set_watch(a, false).await.unwrap();
    let baseline = db.collection_status(B).await.unwrap().items[0]
        .last_success_at
        .clone();
    wait_success(db, B, baseline).await;
    assert_eq!(db.watch_detail(a).await.unwrap()["desired_state"], "PAUSED");
    assert_eq!(db.watch_detail(b).await.unwrap()["desired_state"], "ACTIVE");
    cancel.cancel();
    tokio::time::timeout(Duration::from_secs(8), task)
        .await
        .unwrap()
        .unwrap();
    println!(
        "[PASS] account A panic/failure and pause do not stall account B; supervisor shutdown drains"
    );
}
async fn serve(router: axum::Router) -> (String, tokio::task::JoinHandle<()>) {
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let address = listener.local_addr().unwrap();
    let handle = tokio::spawn(async move { axum::serve(listener, router).await.unwrap() });
    (format!("http://{address}"), handle)
}
async fn management_http(db: &Postgres) {
    let token = "0123456789abcdef0123456789abcdef";
    let (internal, ih) = serve(trade_log_adapters::collector_http::router(
        trade_log_adapters::collector_http::CollectorState {
            reader: Arc::new(db.clone()),
            account: ACCOUNT.into(),
            credential: Arc::new(token.into()),
        },
    ))
    .await;
    let path = std::env::temp_dir().join(format!("robotech-v06-token-{}", Uuid::new_v4()));
    std::fs::write(&path, token).unwrap();
    let mut state = query_api::state::AppState::new();
    state.management_credential = Some(Arc::new(token.into()));
    state.collector = query_api::clients::collector::CollectorClient::from_config(
        &query_api::config::TradeLogConfig {
            enabled: true,
            base_url: Some(internal),
            credential_file: Some(path.clone()),
            request_timeout_seconds: Some(10),
        },
    )
    .unwrap();
    let (publisher_url, ph) = serve(trade_log_adapters::publisher_http::router(
        trade_log_adapters::publisher_http::PublisherState {
            reader: Arc::new(db.clone()),
            account_key: String::new(),
            credential: Arc::new(token.into()),
        },
    ))
    .await;
    state.publisher = query_api::clients::publisher::PublisherClient::from_config(
        &query_api::config::TradeLogConfig {
            enabled: true,
            base_url: Some(publisher_url),
            credential_file: Some(path.clone()),
            request_timeout_seconds: Some(10),
        },
    )
    .unwrap();
    let (url, gh) = serve(query_api::http::router::build(state)).await;
    let client = reqwest::Client::new();
    let r = client
        .post(format!("{url}/api/v1/monitoring/accounts/add"))
        .json(&json!({}))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 401);
    let r = client
        .post(format!("{url}/api/v1/monitoring/accounts/add"))
        .bearer_auth(token)
        .json(&json!({"account":"bad","start_time":"bad"}))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 400);
    let r = client
        .post(format!("{url}/api/v1/monitoring/accounts/add"))
        .bearer_auth(token)
        .json(&json!({"account":ACCOUNT,"start_time":add(ACCOUNT).start_time}))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 409);
    let trace = r.headers()["x-trace-id"].to_str().unwrap().to_owned();
    let data: Value = r.json().await.unwrap();
    assert_eq!(data["trace_id"], trace);
    assert_eq!(data["code"], "WATCH_ACCOUNT_EXISTS");
    let r = client
        .get(format!("{url}/api/v1/monitoring/accounts/list"))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 200);
    let data: Value = r.json().await.unwrap();
    assert_eq!(data["data"]["total"], 3);
    assert_eq!(
        client
            .get(format!(
                "{url}/api/v1/monitoring/accounts/detail?id=not-uuid"
            ))
            .send()
            .await
            .unwrap()
            .status(),
        400
    );
    assert_eq!(
        client
            .get(format!(
                "{url}/api/v1/monitoring/accounts/detail?id={}",
                Uuid::new_v4()
            ))
            .send()
            .await
            .unwrap()
            .status(),
        404
    );
    let row = data["data"]["items"]
        .as_array()
        .unwrap()
        .iter()
        .find(|r| !r["publishing"].is_null())
        .unwrap();
    let id = row["id"].as_str().unwrap();
    let r = client
        .get(format!(
            "{url}/api/v1/monitoring/accounts/publishing-status"
        ))
        .query(&[("id", id)])
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 200);
    let publication: Value = r.json().await.unwrap();
    assert_eq!(publication["data"]["account_key"], row["account_key"]);
    let r = client
        .get(format!("{url}/api/v1/publishing/status"))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 200);
    let summary: Value = r.json().await.unwrap();
    assert!(summary["data"]["items"].is_array());
    assert!(summary["data"]["outbox"].is_object());
    for action in ["pause", "resume"] {
        let r = client
            .post(format!("{url}/api/v1/monitoring/accounts/{action}"))
            .bearer_auth(token)
            .json(&json!({"id": id}))
            .send()
            .await
            .unwrap();
        assert_eq!(r.status(), 202);
        let body: Value = r.json().await.unwrap();
        assert_eq!(body["data"]["id"], id);
        assert_eq!(
            body["data"]["desired_state"],
            if action == "pause" {
                "PAUSED"
            } else {
                "ACTIVE"
            }
        );
    }
    let r = client
        .get(format!(
            "{url}/api/v1/monitoring/accounts/collection-status"
        ))
        .query(&[("id", id)])
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 200);
    let body: Value = r.json().await.unwrap();
    assert_eq!(body["data"]["id"], id);
    assert!(body["data"].get("collection").is_some());
    ih.abort();
    ph.abort();
    gh.abort();
    let _ = std::fs::remove_file(path);
    println!(
        "[PASS] gateway management authentication, 400/404/409 and trace propagation, read-only list"
    );
}
fn policy() -> trade_log::publishing::CandidatePolicy {
    trade_log::publishing::CandidatePolicy {
        max_event_age_seconds: 30,
        signal_ttl_seconds: 60,
        clock_skew_tolerance_seconds: 5,
        version: "candidate-v1".into(),
    }
}
async fn session(db: &Postgres, account: &str) -> (trade_log::checkpoint::Lease, Uuid) {
    let l = start(db, account).await;
    let mut c = collection_config();
    c.account = account.into();
    db.prepare_stream(&l, &c, true).await.unwrap();
    let s = Uuid::new_v4();
    db.begin_stream(&l, &c, s).await.unwrap();
    db.stream_connected(&l, s).await.unwrap();
    db.stream_subscribed(&l, s).await.unwrap();
    (l, s)
}
async fn archive_fill(
    db: &Postgres,
    l: &trade_log::checkpoint::Lease,
    s: Uuid,
    seq: i64,
    account: &str,
    tid: u64,
    snapshot: bool,
) -> String {
    let mut f = fills(1);
    f[0]["time"] = json!(chrono::Utc::now().timestamp_millis());
    f[0]["tid"] = json!(tid);
    let raw = serde_json::to_vec(
        &json!({"channel":"userFills","data":{"user":account,"fills":f,"isSnapshot":snapshot}}),
    )
    .unwrap();
    let trade_log::realtime::StreamEvent::Data(m) =
        hyperliquid::websocket::decode(&raw, account).unwrap()
    else {
        panic!()
    };
    db.archive_stream(l, s, seq, &m, META.as_bytes(), SPOT.as_bytes())
        .await
        .unwrap()
}
async fn commit_fill(
    db: &Postgres,
    l: &trade_log::checkpoint::Lease,
    s: Uuid,
    seq: i64,
    q: &str,
) -> Result<(), QueryError> {
    let result = db
        .stream_result(q, &hyperliquid::parser::HyperliquidParser)
        .await?;
    db.commit_stream(
        &result,
        &trade_log_adapters::postgres::realtime::StreamCommit {
            lease: l,
            session_id: s,
            sequence: seq,
        },
    )
    .await
}
struct Target;
#[async_trait::async_trait]
impl trade_log::publishing::PublicationTarget for Target {
    async fn deliver(
        &self,
        _: &str,
        _: i32,
        body: &[u8],
        _: u64,
        _: u64,
    ) -> trade_log::publishing::DeliveryOutcome {
        let v: Value = serde_json::from_slice(body).unwrap();
        let a = v["fact"]["account_key"]
            .as_str()
            .unwrap()
            .ends_with(ACCOUNT);
        if a {
            tokio::time::sleep(Duration::from_millis(100)).await;
        }
        trade_log::publishing::DeliveryOutcome {
            status: if a { "RETRY_WAIT" } else { "DELIVERED" },
            http: Some(if a { 500 } else { 200 }),
            error: if a { Some("TARGET_UNAVAILABLE") } else { None },
            delay: 1,
        }
    }
}
async fn publishing_isolation(db: &Postgres, a: Uuid, b: Uuid) {
    db.set_watch(a, true).await.unwrap();
    db.set_watch(b, true).await.unwrap();
    let (la, sa) = session(db, ACCOUNT).await;
    let (lb, sb) = session(db, B).await;
    for l in [&la, &lb] {
        db.configure_publishing(&l.key, "https://example.invalid/events", true, &policy())
            .await
            .unwrap();
    }
    for (l, s, account) in [(&la, sa, ACCOUNT), (&lb, sb, B)] {
        let q = archive_fill(db, l, s, 1, account, 600, true).await;
        commit_fill(db, l, s, 1, &q).await.unwrap();
        let q = archive_fill(db, l, s, 2, account, 601, false).await;
        commit_fill(db, l, s, 2, &q).await.unwrap();
    }
    assert_eq!(count(db, "outbox_events").await, 2);
    let q = archive_fill(db, &la, sa, 3, ACCOUNT, 602, false).await;
    db.set_watch(a, false).await.unwrap();
    assert_eq!(
        commit_fill(db, &la, sa, 3, &q).await.unwrap_err().code,
        "LEASE_LOST"
    );
    // Old pending WS data may be recovered as a fact, but must not become a fresh signal.
    db.set_watch(a, true).await.unwrap();
    let new_a = start(db, ACCOUNT).await;
    db.configure_publishing(&la.key, "https://example.invalid/events", true, &policy())
        .await
        .unwrap();
    db.resume_stream_job(&new_a, &q).await.unwrap();
    commit_fill(db, &new_a, sa, 3, &q).await.unwrap();
    assert_eq!(count(db, "outbox_events").await, 2);
    db.set_watch(a, false).await.unwrap();
    let cancel = CancellationToken::new();
    let runtime = trade_log_adapters::publisher_runtime::PublisherRuntime {
        store: db.clone(),
        webhook: Some(Arc::new(Target)),
        lease_seconds: 30,
        poll_interval_ms: 50,
        retry_base_seconds: 1,
        retry_max_seconds: 5,
    };
    let task = tokio::spawn(runtime.run_multi(
        cancel.clone(),
        "https://example.invalid/events".into(),
        true,
        policy(),
        1,
    ));
    tokio::time::timeout(Duration::from_secs(10),async{loop{
  let good:i64=sqlx::query_scalar("SELECT count(*) FROM trade_log.outbox_events WHERE partition_key=$1 AND status='DELIVERED'").bind(&lb.key).fetch_one(&db.pool).await.unwrap();
  let bad:i64=sqlx::query_scalar("SELECT COALESCE(max(attempts),0)::bigint FROM trade_log.outbox_events WHERE partition_key=$1").bind(&la.key).fetch_one(&db.pool).await.unwrap();
  if good==1&&bad>=1{break}tokio::time::sleep(Duration::from_millis(50)).await;
 }}).await.unwrap();
    cancel.cancel();
    task.await.unwrap();
    assert_eq!(db.watch_detail(a).await.unwrap()["desired_state"], "PAUSED");
    db.release_collection(&lb).await.unwrap();
    db.set_watch(a, true).await.unwrap();
    println!(
        "[PASS] two-account fact identity; snapshot suppression; stale archived WS cannot publish after resume; paused backlog drains; failing A does not starve B with one send slot"
    );
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL administrator"]
async fn upgrade_preserves_v05_checkpoint_control_and_fact_identity() {
    use sqlx::{Connection, PgConnection, postgres::PgConnectOptions};
    use std::str::FromStr;
    let options =
        PgConnectOptions::from_str(&std::env::var("ROBOTECH_TEST_DATABASE_URL").unwrap()).unwrap();
    let mut admin = PgConnection::connect_with(&options).await.unwrap();
    let name = format!("robotech_v06_upgrade_{}", Uuid::new_v4().simple());
    sqlx::query(&format!("CREATE DATABASE {name}"))
        .execute(&mut admin)
        .await
        .unwrap();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .connect_with(options.clone().database(&name))
        .await
        .unwrap();
    let initial = sqlx::migrate::Migrator {
        migrations: std::borrow::Cow::Owned(
            trade_log_adapters::postgres::migration::MIGRATOR
                .iter()
                .take(4)
                .cloned()
                .collect(),
        ),
        ..sqlx::migrate::Migrator::DEFAULT
    };
    initial.run(&pool).await.unwrap();
    let db = Postgres {
        pool,
        network: shared_types::Network::Mainnet,
        mirror: None,
        stream_memory: None,
    };
    let saved = legacy_collect(&db, "query_v05_migration", fills(2))
        .await
        .unwrap();
    let key = db.checkpoint_key(ACCOUNT);
    sqlx::query("INSERT INTO trade_log.collection_checkpoints(chain_id,partition_key,position) VALUES('hyperliquid:mainnet',$1,$2)").bind(&key).bind(json!({"initial_start_ms":1791097200000_i64,"scanned_through_ms":1791097201000_i64,"last_trade_at_ms":1791097200001_i64})).execute(&db.pool).await.unwrap();
    sqlx::query("INSERT INTO trade_log.publishing_control(account_key,target_id,target_url,enabled,activation_epoch,policy) VALUES($1,'legacy','https://example.invalid/events',true,7,$2)").bind(&key).bind(json!(policy())).execute(&db.pool).await.unwrap();
    let before: chrono::DateTime<chrono::Utc> =
        sqlx::query_scalar("SELECT activated_at FROM trade_log.publishing_control")
            .fetch_one(&db.pool)
            .await
            .unwrap();
    db.migrate().await.unwrap();
    db.migrate().await.unwrap();
    db.import_watch(ACCOUNT, "2026-01-01T00:00:00Z")
        .await
        .unwrap();
    let rows = db.watch_rows().await.unwrap();
    assert_eq!(rows.len(), 1);
    assert_eq!(rows[0]["desired_state"], "ACTIVE");
    assert_eq!(
        count(&db, "account_facts_current").await,
        saved.trades.len() as i64
    );
    db.configure_publishing(&key, "https://example.invalid/events", true, &policy())
        .await
        .unwrap();
    let status = db.publishing_status(&key).await.unwrap();
    assert_eq!(status["activation_epoch"], 7);
    assert_eq!(status["activated_at"], json!(before));
    let position: Value =
        sqlx::query_scalar("SELECT position FROM trade_log.collection_checkpoints")
            .fetch_one(&db.pool)
            .await
            .unwrap();
    assert_eq!(position["scanned_through_ms"], 1791097201000_i64);
    db.pool.close().await;
    sqlx::query(&format!("DROP DATABASE {name} WITH (FORCE)"))
        .execute(&mut admin)
        .await
        .unwrap();
    println!(
        "[PASS] v0.5 schema upgrade retains facts, HTTP watermark, original activation epoch and timestamp; rerun is idempotent"
    );
}

async fn roles(db: &Postgres) {
    sqlx::raw_sql("DO $$ BEGIN IF NOT EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_publisher') THEN CREATE ROLE trade_log_publisher LOGIN; END IF; IF NOT EXISTS(SELECT FROM pg_roles WHERE rolname='trade_log_app') THEN CREATE ROLE trade_log_app LOGIN; END IF; END $$").execute(&db.pool).await.unwrap();
    db.migrate().await.unwrap();
    let options = (*db.pool.connect_options()).clone();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .after_connect(|c, _| {
            Box::pin(async move {
                sqlx::query("SET ROLE trade_log_publisher")
                    .execute(c)
                    .await?;
                Ok(())
            })
        })
        .connect_with(options)
        .await
        .unwrap();
    let publisher = Postgres {
        pool,
        network: db.network.clone(),
        mirror: None,
        stream_memory: None,
    };
    let key = db.checkpoint_key(B);
    publisher
        .configure_publishing(&key, "https://example.invalid/events", true, &policy())
        .await
        .unwrap();
    assert!(
        sqlx::query("UPDATE trade_log.watch_accounts SET desired_state='PAUSED'")
            .execute(&publisher.pool)
            .await
            .is_err()
    );
    assert!(
        sqlx::query("UPDATE trade_log.collection_checkpoints SET status='STOPPED'")
            .execute(&publisher.pool)
            .await
            .is_err()
    );
    publisher.pool.close().await;
    let pool = sqlx::postgres::PgPoolOptions::new()
        .after_connect(|c, _| {
            Box::pin(async move {
                sqlx::query("SET ROLE trade_log_app").execute(c).await?;
                Ok(())
            })
        })
        .connect_with((*db.pool.connect_options()).clone())
        .await
        .unwrap();
    let app = Postgres {
        pool,
        network: db.network.clone(),
        mirror: None,
        stream_memory: None,
    };
    assert!(
        sqlx::query("UPDATE trade_log.watch_accounts SET desired_state='PAUSED'")
            .execute(&app.pool)
            .await
            .is_err()
    );
    let _ = app.reserve_source("HTTP", 20).await.unwrap();
    app.pool.close().await;
    println!(
        "[PASS] publisher reads registry and configures control without registry/checkpoint writes; query role reserves shared budget without management writes"
    );
}
