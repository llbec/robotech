mod common;
use common::*;
use hyperliquid::parser::HyperliquidParser;
use serde_json::{Value, json};
use std::sync::Arc;
use trade_log::{
    acquisition::{SourceReader, SourceResponse},
    instrument::InstrumentType,
    parsing::{ParseContext, Parsed, ProtocolParser},
    query::{QueryError, QueryService},
    realtime::{StreamEvent, StreamMessage},
    stored_query::StoredQuery,
    validation::QueryRequest,
};
use trade_log_adapters::{
    postgres::{Postgres, realtime::StreamCommit},
    spot_backfill_runtime::SpotBackfillRuntime,
};
use uuid::Uuid;
fn spot(tid: u64, time: i64) -> Value {
    let mut f: Value = serde_json::from_str(include_str!(
        "../../../../tests/fixtures/hyperliquid/spot/fills.json"
    ))
    .unwrap();
    f[0]["tid"] = json!(tid);
    f[0]["time"] = json!(time);
    f
}
struct Legacy;
impl ProtocolParser for Legacy {
    fn parse(&self, c: ParseContext<'_>) -> Result<Parsed, QueryError> {
        let mut p = HyperliquidParser.parse(c)?;
        p.trades
            .retain(|f| f.payload.instrument_type == "PERPETUAL");
        Ok(p)
    }
}
async fn register(db: &Postgres) -> Uuid {
    db.import_watch(ACCOUNT, &collection_config().start_time)
        .await
        .unwrap();
    let rows = db.watch_rows().await.unwrap();
    Uuid::parse_str(rows[0]["id"].as_str().unwrap()).unwrap()
}
#[tokio::test]
#[ignore = "Requires isolated PostgreSQL"]
async fn server_spot_suite() {
    let db = database().await;
    let watch_id = register(&db).await;
    let t = 1791097200000_i64;
    let mut mixed = fills(1);
    mixed.as_array_mut().unwrap().push(spot(22, t)[0].clone());
    let q = "query_spot_legacy";
    QueryService {
        network: db.network.clone(),
        source: Arc::new(Source(mixed.clone())),
        parser: Arc::new(Legacy),
        evidence: Arc::new(db.clone()),
    }
    .execute(
        q,
        &QueryRequest {
            account: ACCOUNT.into(),
            limit: 10,
            instrument_type: InstrumentType::Perpetual,
        },
        "trace_spot",
    )
    .await
    .unwrap();
    assert_eq!(count(&db, "account_facts_current").await, 1);
    let legacy: Value =
        sqlx::query_scalar("SELECT result FROM trade_log.collection_jobs WHERE query_id=$1")
            .bind(q)
            .fetch_one(&db.pool)
            .await
            .unwrap();
    let job = db
        .submit_spot_backfill(ACCOUNT, t - 10, t + 10)
        .await
        .unwrap();
    let same = db
        .submit_spot_backfill(ACCOUNT, t - 10, t + 10)
        .await
        .unwrap();
    assert_eq!(job["id"], same["id"]);
    let runtime = SpotBackfillRuntime {
        store: db.clone(),
        source: Arc::new(RangeSource(json!([]))),
        parser: Arc::new(HyperliquidParser),
    };
    for _ in 0..5 {
        if let Some(j) = db.claim_spot_backfill(Uuid::new_v4()).await.unwrap() {
            runtime.batch(&j).await.unwrap();
        }
    }
    let done = db
        .spot_backfill_status(Uuid::parse_str(job["id"].as_str().unwrap()).unwrap())
        .await
        .unwrap();
    assert_eq!(done["status"], "COMPLETED");
    assert_eq!(done["scanned_ms"], t + 10);
    let preserved: Value =
        sqlx::query_scalar("SELECT result FROM trade_log.collection_jobs WHERE query_id=$1")
            .bind(q)
            .fetch_one(&db.pool)
            .await
            .unwrap();
    assert_eq!(preserved, legacy);
    let derived:String=sqlx::query_scalar("SELECT query_id FROM trade_log.collection_jobs WHERE request->>'origin_query_id'=$1 LIMIT 1").bind(q).fetch_one(&db.pool).await.unwrap();
    assert_eq!(
        db.reparse(&derived, &HyperliquidParser)
            .await
            .unwrap()
            .comparison,
        "SAME"
    );
    assert_eq!(count(&db, "account_facts_current").await, 2);
    assert_eq!(count(&db, "outbox_events").await, 0);
    assert_eq!(count(&db, "collection_checkpoints").await, 0);
    let r = collect(&db, "query_spot_mixed", mixed, 1).await.unwrap();
    assert_eq!(r.trades.len(), 1);
    assert_eq!(r.trades[0].payload.instrument_type, "PERPETUAL");
    assert_eq!(r.counts.counts_version, 2);
    assert_eq!(r.counts.spot_records, 1);
    assert_eq!(r.counts.perpetual_records, 1);
    assert_eq!(count(&db, "account_facts_current").await, 2);
    let typed = QueryService {
        network: db.network.clone(),
        source: Arc::new(Source({
            let mut m = fills(1);
            m.as_array_mut().unwrap().push(spot(22, t)[0].clone());
            m
        })),
        parser: Arc::new(HyperliquidParser),
        evidence: Arc::new(db.clone()),
    }
    .execute(
        "query_spot_typed",
        &QueryRequest {
            account: ACCOUNT.into(),
            limit: 1,
            instrument_type: InstrumentType::Spot,
        },
        "trace_typed",
    )
    .await
    .unwrap();
    assert_eq!(typed.trades.len(), 1);
    assert_eq!(typed.trades[0].payload.instrument_type, "SPOT");
    assert_eq!(typed.counts.matched_records, 1);
    assert_eq!(typed.persistence.unwrap().existing_records, 2);
    let mut req = request(1);
    req.instrument_type = InstrumentType::Spot;
    let s = db.stored(&req).await.unwrap();
    assert_eq!(s.matched_records, 1);
    assert_eq!(s.trades[0].payload.instrument_type, "SPOT");
    req.instrument_type = InstrumentType::All;
    let all = db.stored(&req).await.unwrap();
    assert!(all.has_more);
    let cursor = all.next_cursor.unwrap();
    req.cursor = Some(cursor);
    let second = db.stored(&req).await.unwrap();
    assert_ne!(second.trades[0].fact_id, all.trades[0].fact_id);
    req.instrument_type = InstrumentType::Spot;
    assert!(db.stored(&req).await.is_err());
    let another = db
        .submit_spot_backfill(ACCOUNT, t - 20, t + 20)
        .await
        .unwrap();
    let claimed = db
        .claim_spot_backfill(Uuid::new_v4())
        .await
        .unwrap()
        .unwrap();
    db.set_watch(watch_id, false).await.unwrap();
    assert!(
        db.claim_spot_backfill(Uuid::new_v4())
            .await
            .unwrap()
            .is_none()
    );
    assert!(runtime.batch(&claimed).await.is_err());
    assert_eq!(count(&db, "account_facts_current").await, 2);
    assert_ne!(another["status"], "COMPLETED");
    db.fail_spot_batch(&claimed, &QueryError::incomplete("fixed test error"))
        .await;
    let retried = db
        .submit_spot_backfill(ACCOUNT, t - 20, t + 20)
        .await
        .unwrap();
    assert_eq!(retried["id"], another["id"]);
    assert_eq!(retried["status"], "PENDING");
    db.set_watch(watch_id, true).await.unwrap();
    println!(
        "[PASS] spot archive backfill, source scan, idempotence, immutable old receipt, no watermark rewind or historic candidates; filtered pagination and paused generation fencing"
    );
    publication(&db).await;
    db.pool.close().await;
}
struct RangeSource(Value);
#[async_trait::async_trait]
impl SourceReader for RangeSource {
    async fn wait_retry(&self, _: u64) {}
    async fn fetch(
        &self,
        k: protocol_api::QueryKind,
        a: &str,
    ) -> Result<SourceResponse, QueryError> {
        Source(json!([])).fetch(k, a).await
    }
    async fn fetch_range(&self, _: &str, _: i64, _: i64) -> Result<SourceResponse, QueryError> {
        Ok(SourceResponse {
            status: 200,
            body: serde_json::to_vec(&self.0).unwrap(),
            received_at: shared_types::now(),
            retry_after_seconds: None,
        })
    }
}
fn message(fills: Value, snapshot: bool) -> StreamMessage {
    let StreamEvent::Data(m)=hyperliquid::websocket::decode(&serde_json::to_vec(&json!({"channel":"userFills","data":{"user":ACCOUNT,"fills":fills,"isSnapshot":snapshot}})).unwrap(),ACCOUNT).unwrap() else{panic!()};
    m
}
async fn publication(db: &Postgres) {
    let received = Arc::new(tokio::sync::Mutex::new(Vec::<Value>::new()));
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let url = format!("http://{}/events", listener.local_addr().unwrap());
    let app = axum::Router::new()
        .route(
            "/events",
            axum::routing::post(
                |axum::extract::State(received): axum::extract::State<
                    Arc<tokio::sync::Mutex<Vec<Value>>>,
                >,
                 axum::Json(body): axum::Json<Value>| async move {
                    received.lock().await.push(body);
                    axum::http::StatusCode::OK
                },
            ),
        )
        .with_state(received.clone());
    let receiver = tokio::spawn(async move { axum::serve(listener, app).await.unwrap() });
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
    let p = trade_log::publishing::CandidatePolicy {
        max_event_age_seconds: 30,
        signal_ttl_seconds: 60,
        clock_skew_tolerance_seconds: 5,
        version: "candidate-v1".into(),
    };
    db.configure_publishing(&l.key, &url, true, &p)
        .await
        .unwrap();
    let original = db.publishing_status(&l.key).await.unwrap();
    db.configure_spot_publishing(&l.key, false, &p)
        .await
        .unwrap();
    for (seq, snapshot, enabled) in [
        (1, true, false),
        (2, false, false),
        (3, true, true),
        (4, false, true),
    ] {
        if enabled {
            db.configure_spot_publishing(&l.key, true, &p)
                .await
                .unwrap();
        }
        tokio::time::sleep(std::time::Duration::from_millis(5)).await;
        let fills = spot(100 + seq as u64, chrono::Utc::now().timestamp_millis());
        let q = db
            .archive_stream(
                &l,
                session,
                seq,
                &message(fills, snapshot),
                META.as_bytes(),
                SPOT.as_bytes(),
            )
            .await
            .unwrap();
        let result = db.stream_result(&q, &HyperliquidParser).await.unwrap();
        db.commit_stream(
            &result,
            &StreamCommit {
                lease: &l,
                session_id: session,
                sequence: seq,
            },
        )
        .await
        .unwrap();
        assert_eq!(
            count(db, "outbox_events").await,
            if seq == 4 { 1 } else { 0 }
        );
    }
    let after = db.publishing_status(&l.key).await.unwrap();
    assert_eq!(original["activation_epoch"], after["activation_epoch"]);
    assert_eq!(original["activated_at"], after["activated_at"]);
    assert_eq!(after["instruments"]["SPOT"]["enabled"], true);
    let body: Value = sqlx::query_scalar("SELECT payload FROM trade_log.outbox_events")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    assert_eq!(body["fact"]["payload"]["instrument_type"], "SPOT");
    assert_eq!(body["fact"]["payload"]["fee_asset"], "PURR");
    let cancel = tokio_util::sync::CancellationToken::new();
    let runtime = trade_log_adapters::publisher_runtime::PublisherRuntime {
        store: db.clone(),
        webhook: Some(Arc::new(
            trade_log_adapters::webhook::Webhook::new(&url, &"1".repeat(64), 5, 65536).unwrap(),
        )),
        lease_seconds: 30,
        poll_interval_ms: 50,
        retry_base_seconds: 1,
        retry_max_seconds: 5,
    };
    let worker = tokio::spawn(runtime.run_multi_with_spot(cancel.clone(), url, true, p, 1, true));
    tokio::time::timeout(std::time::Duration::from_secs(5), async {
        loop {
            let n: i64 = sqlx::query_scalar(
                "SELECT count(*) FROM trade_log.outbox_events WHERE status='DELIVERED'",
            )
            .fetch_one(&db.pool)
            .await
            .unwrap();
            if n == 1 {
                break;
            }
            tokio::time::sleep(std::time::Duration::from_millis(20)).await;
        }
    })
    .await
    .unwrap();
    cancel.cancel();
    worker.await.unwrap();
    assert_eq!(received.lock().await.len(), 1);
    let wire: Vec<u8> = sqlx::query_scalar("SELECT wire_body FROM trade_log.outbox_events")
        .fetch_one(&db.pool)
        .await
        .unwrap();
    let mut published: Value = serde_json::from_slice(&wire).unwrap();
    assert_eq!(received.lock().await[0], published);
    assert!(published["published_at"].is_string());
    published["published_at"] = Value::Null;
    assert_eq!(published, body);
    let policy = trade_log::publishing::CandidatePolicy {
        max_event_age_seconds: 30,
        signal_ttl_seconds: 60,
        clock_skew_tolerance_seconds: 5,
        version: "candidate-v1".into(),
    };
    db.configure_spot_publishing(&l.key, false, &policy)
        .await
        .unwrap();
    db.configure_spot_publishing(&l.key, true, &policy)
        .await
        .unwrap();
    let old_time = chrono::DateTime::parse_from_rfc3339(body["occurred_at"].as_str().unwrap())
        .unwrap()
        .timestamp_millis();
    let q = db
        .archive_stream(
            &l,
            session,
            5,
            &message(spot(105, old_time), false),
            META.as_bytes(),
            SPOT.as_bytes(),
        )
        .await
        .unwrap();
    let result = db.stream_result(&q, &HyperliquidParser).await.unwrap();
    db.commit_stream(
        &result,
        &StreamCommit {
            lease: &l,
            session_id: session,
            sequence: 5,
        },
    )
    .await
    .unwrap();
    assert_eq!(count(db, "outbox_events").await, 1);
    let decision:(i64,String)=sqlx::query_as("SELECT activation_epoch,reason FROM trade_log.publication_decisions WHERE message_sequence=5").fetch_one(&db.pool).await.unwrap();
    assert_eq!(decision, (2, "BEFORE_ACTIVATION".into()));
    assert_eq!(
        db.publishing_status(&l.key).await.unwrap()["activation_epoch"],
        original["activation_epoch"]
    );
    receiver.abort();
    println!(
        "[PASS] spot disabled by default; snapshots suppressed; fresh spot delivered to HTTP receiver with immutable body; perpetual activation unchanged"
    );
}
