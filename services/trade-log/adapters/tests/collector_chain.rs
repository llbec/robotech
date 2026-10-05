mod common;
use async_trait::async_trait;
use common::*;
use hyperliquid::parser::HyperliquidParser;
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
};
struct Flaky {
    fail: Arc<AtomicBool>,
    body: serde_json::Value,
}
#[async_trait]
impl SourceReader for Flaky {
    async fn wait_retry(&self, _: u64) {}
    async fn fetch(
        &self,
        k: protocol_api::QueryKind,
        a: &str,
    ) -> Result<SourceResponse, QueryError> {
        Source(self.body.clone()).fetch(k, a).await
    }
    async fn fetch_range(
        &self,
        a: &str,
        start: i64,
        end: i64,
    ) -> Result<SourceResponse, QueryError> {
        if self.fail.load(Ordering::SeqCst) {
            return Ok(SourceResponse {
                status: 429,
                body: b"limited".to_vec(),
                received_at: shared_types::now(),
                retry_after_seconds: Some(3),
            });
        }
        RangeSource(self.body.clone())
            .fetch_range(a, start, end)
            .await
    }
}
async fn wait_status(db: &trade_log_adapters::postgres::Postgres, state: &str) {
    tokio::time::timeout(Duration::from_secs(12), async {
        loop {
            if let Ok(s) = db.collection_status(ACCOUNT).await
                && s.items[0].status == state
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
async fn automatic_loop_reports_failure_recovers_and_stops_with_a_persistent_checkpoint() {
    let db = database().await;
    let mut c = collection_config();
    c.start_time = trade_log::collection::time(chrono::Utc::now().timestamp_millis() - 5000);
    c.safety_delay_seconds = 0;
    let failure = Arc::new(AtomicBool::new(true));
    let cancel = CancellationToken::new();
    let runtime = trade_log_adapters::collector_runtime::CollectorRuntime {
        store: db.clone(),
        config: c.clone(),
        source: Arc::new(Flaky {
            fail: failure.clone(),
            body: serde_json::json!([]),
        }),
        parser: Arc::new(HyperliquidParser),
    };
    let run_cancel = cancel.clone();
    let worker = tokio::spawn(async move {
        runtime.run(run_cancel).await;
    });
    wait_status(&db, "RETRY_WAIT").await;
    let s = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert_eq!(s.consecutive_failures, 1);
    assert_eq!(s.last_error.unwrap()["code"], "RATE_LIMITED");
    assert!(s.scanned_through.is_none());
    assert!(s.last_success_at.is_none());
    let scheduled =
        chrono::DateTime::parse_from_rfc3339(s.next_run_at.as_deref().unwrap()).unwrap();
    let error_at =
        chrono::DateTime::parse_from_rfc3339(s.last_attempt_at.as_deref().unwrap()).unwrap();
    assert!(scheduled.signed_duration_since(error_at).num_milliseconds() >= 3000);
    failure.store(false, Ordering::SeqCst);
    wait_status(&db, "WAITING").await;
    let first = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert_eq!(first.consecutive_failures, 0);
    assert!(first.last_error.is_none());
    assert!(first.scanned_through.is_some());
    tokio::time::timeout(Duration::from_secs(6), async {
        loop {
            let s = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
            if s.last_success_query_id != first.last_success_query_id {
                break;
            }
            tokio::time::sleep(Duration::from_millis(30)).await;
        }
    })
    .await
    .unwrap();
    cancel.cancel();
    worker.await.unwrap();
    assert_eq!(
        db.collection_status(ACCOUNT).await.unwrap().items[0].status,
        "STOPPED"
    );
    assert!(
        db.collection_status(ACCOUNT).await.unwrap().items[0]
            .scanned_through
            .is_some()
    );
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn gateway_status_authentication_and_database_fault_mapping() {
    async fn serve(
        r: axum::Router,
    ) -> (
        String,
        tokio::sync::oneshot::Sender<()>,
        tokio::task::JoinHandle<()>,
    ) {
        let l = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
        let address = l.local_addr().unwrap();
        let (tx, rx) = tokio::sync::oneshot::channel();
        let h = tokio::spawn(async move {
            axum::serve(l, r)
                .with_graceful_shutdown(async {
                    let _ = rx.await;
                })
                .await
                .unwrap();
        });
        (format!("http://{address}"), tx, h)
    }
    let db = database().await;
    let c = collection_config();
    db.acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap();
    let token = "0123456789abcdef0123456789abcdef";
    let (internal, si, hi) = serve(trade_log_adapters::collector_http::router(
        trade_log_adapters::collector_http::CollectorState {
            reader: Arc::new(db.clone()),
            account: ACCOUNT.into(),
            credential: Arc::new(token.into()),
        },
    ))
    .await;
    let path =
        std::env::temp_dir().join(format!("robotech-collector-token-{}", uuid::Uuid::new_v4()));
    std::fs::write(&path, token).unwrap();
    let mut state = query_api::state::AppState::new();
    state.collector = query_api::clients::collector::CollectorClient::from_config(
        &query_api::config::TradeLogConfig {
            enabled: true,
            base_url: Some(internal.clone()),
            credential_file: Some(path.clone()),
            request_timeout_seconds: Some(5),
        },
    )
    .unwrap();
    let (gateway, sg, hg) = serve(query_api::http::router::build(state)).await;
    let client = reqwest::Client::new();
    assert_eq!(
        client
            .get(format!("{internal}/internal/v1/collection-status"))
            .send()
            .await
            .unwrap()
            .status(),
        401
    );
    let r = client
        .get(format!("{gateway}/api/v1/watch-accounts"))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 200);
    let trace = r.headers()["x-trace-id"].to_str().unwrap().to_string();
    let body: serde_json::Value = r.json().await.unwrap();
    assert_eq!(body["meta"]["trace_id"], trace);
    assert_eq!(body["data"]["items"][0]["account"], ACCOUNT);
    assert_eq!(
        body["data"]["items"][0]["coverage"],
        "SOURCE_HISTORY_NOT_VERIFIED"
    );
    assert_eq!(
        client
            .post(format!("{gateway}/api/v1/watch-accounts"))
            .send()
            .await
            .unwrap()
            .status(),
        405
    );
    db.pool.close().await;
    assert_eq!(
        client
            .get(format!("{gateway}/api/v1/watch-accounts"))
            .send()
            .await
            .unwrap()
            .status(),
        503
    );
    assert_eq!(
        client
            .get(format!("{gateway}/api/v1/health"))
            .send()
            .await
            .unwrap()
            .status(),
        200
    );
    sg.send(()).unwrap();
    si.send(()).unwrap();
    hg.await.unwrap();
    hi.await.unwrap();
    std::fs::remove_file(path).unwrap();
}
#[tokio::test]
async fn official_time_request_preserves_inclusive_millisecond_boundaries() {
    use axum::{Json, Router, routing::post};
    async fn inspect(Json(body): Json<serde_json::Value>) -> Json<serde_json::Value> {
        assert_eq!(body["type"], "userFillsByTime");
        assert_eq!(body["user"], ACCOUNT);
        assert_eq!(body["aggregateByTime"], false);
        assert_eq!(body["startTime"], 10);
        assert_eq!(body["endTime"], 14);
        Json(serde_json::json!([]))
    }
    let l = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let address = l.local_addr().unwrap();
    let (tx, rx) = tokio::sync::oneshot::channel();
    let server = tokio::spawn(async move {
        axum::serve(l, Router::new().route("/info", post(inspect)))
            .with_graceful_shutdown(async {
                let _ = rx.await;
            })
            .await
            .unwrap();
    });
    let source = hyperliquid::source::HttpSource::new(
        &format!("http://{address}/info"),
        Duration::from_secs(1),
        Duration::from_secs(1),
        1024,
    )
    .unwrap();
    let response = source.fetch_range(ACCOUNT, 10, 15).await.unwrap();
    assert_eq!(response.status, 200);
    assert_eq!(response.body, b"[]");
    assert!(source.fetch_range(ACCOUNT, 10, 10).await.is_err());
    tx.send(()).unwrap();
    server.await.unwrap();
}
struct Blocked(Arc<AtomicBool>);
#[async_trait]
impl SourceReader for Blocked {
    async fn wait_retry(&self, _: u64) {}
    async fn fetch(
        &self,
        k: protocol_api::QueryKind,
        a: &str,
    ) -> Result<SourceResponse, QueryError> {
        Source(serde_json::json!([])).fetch(k, a).await
    }
    async fn fetch_range(&self, _: &str, _: i64, _: i64) -> Result<SourceResponse, QueryError> {
        self.0.store(true, Ordering::SeqCst);
        std::future::pending().await
    }
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn cancellation_preserves_pending_work_and_a_new_instance_completes_it() {
    let db = database().await;
    let mut c = collection_config();
    c.start_time = trade_log::collection::time(chrono::Utc::now().timestamp_millis() - 5000);
    c.safety_delay_seconds = 0;
    let started = Arc::new(AtomicBool::new(false));
    let cancel = CancellationToken::new();
    let runtime = trade_log_adapters::collector_runtime::CollectorRuntime {
        store: db.clone(),
        config: c.clone(),
        source: Arc::new(Blocked(started.clone())),
        parser: Arc::new(HyperliquidParser),
    };
    let child = cancel.clone();
    let worker = tokio::spawn(async move {
        runtime.run(child).await;
    });
    tokio::time::timeout(Duration::from_secs(3), async {
        while !started.load(Ordering::SeqCst) {
            tokio::time::sleep(Duration::from_millis(10)).await;
        }
    })
    .await
    .unwrap();
    cancel.cancel();
    worker.await.unwrap();
    let old = db.collection_status(ACCOUNT).await.unwrap().items.remove(0);
    assert_eq!(old.status, "STOPPED");
    assert!(old.pending_range.is_some());
    assert!(old.scanned_through.is_none());
    assert_eq!(count(&db, "raw_logs").await, 2);
    let l = db
        .acquire_collection(&c, uuid::Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    let mut work = db
        .collection_work(&l, &c, chrono::Utc::now().timestamp_millis())
        .await
        .unwrap()
        .unwrap();
    assert_eq!(Some(work.query_id.clone()), old.last_query_id);
    assert!(
        common::runtime(&db, c, serde_json::json!([]))
            .round(&l, &mut work)
            .await
            .unwrap()
    );
    assert_eq!(count(&db, "raw_logs").await, 3);
    db.pool.close().await;
}
