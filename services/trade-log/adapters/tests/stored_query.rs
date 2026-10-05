mod common;
use common::*;
use trade_log::stored_query::StoredQuery;
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn range_boundaries_empty_and_large_tid_sort() {
    let db = database().await;
    let mut source = fills(3);
    source[1]["time"] = source[0]["time"].clone();
    collect(&db, "query_range", source, 100).await.unwrap();
    let all = db.stored(&request(100)).await.unwrap();
    assert_eq!(all.trades[1].source_ref, "9007199254740994");
    let mut r = request(100);
    r.start_time = Some(all.trades[2].occurred_at.clone());
    r.end_time = Some(all.trades[0].occurred_at.clone());
    let data = db.stored(&r).await.unwrap();
    assert_eq!(data.matched_records, 2);
    r.start_time = None;
    r.end_time = Some(all.trades[2].occurred_at.clone());
    assert_eq!(db.stored(&r).await.unwrap().matched_records, 0);
    r.account = "0x0000000000000000000000000000000000000002".into();
    assert_eq!(db.stored(&r).await.unwrap().matched_records, 0);
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn cursor_snapshot_excludes_later_older_inserts_and_rejects_conditions() {
    let db = database().await;
    collect(&db, "query_page", fills(3), 100).await.unwrap();
    let first = db.stored(&request(1)).await.unwrap();
    assert!(first.has_more);
    let mut later = fills(1);
    later[0]["tid"] = serde_json::json!(999);
    later[0]["time"] = serde_json::json!(1791097100000_i64);
    collect(&db, "query_later", later, 100).await.unwrap();
    let mut r = request(1);
    r.cursor = first.next_cursor;
    let second = db.stored(&r).await.unwrap();
    assert_eq!(second.snapshot_seq, first.snapshot_seq);
    assert_eq!(second.matched_records, 3);
    assert_ne!(first.trades[0].fact_id, second.trades[0].fact_id);
    r.cursor = second.next_cursor;
    let third = db.stored(&r).await.unwrap();
    assert!(!third.has_more);
    assert!(third.next_cursor.is_none());
    r.account = "0x0000000000000000000000000000000000000002".into();
    assert_eq!(db.stored(&r).await.unwrap_err().code, "VALIDATION_ERROR");
    r = request(1);
    r.cursor = Some("broken".into());
    assert_eq!(db.stored(&r).await.unwrap_err().code, "VALIDATION_ERROR");
    assert_eq!(db.stored(&request(100)).await.unwrap().matched_records, 4);
    db.pool.close().await;
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn gateway_internal_postgres_roundtrip_and_error_mapping() {
    use hyperliquid::parser::HyperliquidParser;
    use std::{sync::Arc, time::Duration};
    use tokio::{
        net::TcpListener,
        sync::{Semaphore, oneshot},
    };
    use trade_log::query::QueryService;
    use trade_log_adapters::internal_http::{self, InternalState};
    async fn serve(
        router: axum::Router,
    ) -> (String, oneshot::Sender<()>, tokio::task::JoinHandle<()>) {
        let listener = TcpListener::bind("127.0.0.1:0").await.unwrap();
        let address = listener.local_addr().unwrap();
        let (tx, rx) = oneshot::channel();
        let task = tokio::spawn(async move {
            axum::serve(listener, router)
                .with_graceful_shutdown(async {
                    let _ = rx.await;
                })
                .await
                .unwrap();
        });
        (format!("http://{address}"), tx, task)
    }
    let db = database().await;
    let service = Arc::new(QueryService {
        network: db.network.clone(),
        source: Arc::new(Source(fills(3))),
        parser: Arc::new(HyperliquidParser),
        evidence: Arc::new(db.clone()),
    });
    let token = "0123456789abcdef0123456789abcdef";
    let (internal, stop_i, task_i) = serve(internal_http::router(InternalState {
        service,
        stored: Some(Arc::new(db.clone())),
        credential: Arc::new(token.into()),
        permits: Arc::new(Semaphore::new(4)),
        timeout: Duration::from_secs(3),
    }))
    .await;
    let token_path =
        std::env::temp_dir().join(format!("robotech-http-token-{}", uuid::Uuid::new_v4()));
    std::fs::write(&token_path, token).unwrap();
    let mut state = query_api::state::AppState::new();
    state.trade_log = query_api::clients::trade_log::TradeLogClient::from_config(
        &query_api::config::TradeLogConfig {
            enabled: true,
            base_url: Some(internal.clone()),
            credential_file: Some(token_path.clone()),
            request_timeout_seconds: Some(5),
        },
    )
    .unwrap();
    let (gateway, stop_g, task_g) = serve(query_api::http::router::build(state)).await;
    let client = reqwest::Client::new();
    let r = client
        .get(format!(
            "{gateway}/api/v1/trade-events?account={ACCOUNT}&limit=1"
        ))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 200);
    let trace = r.headers()["x-trace-id"].to_str().unwrap().to_owned();
    let result: serde_json::Value = r.json().await.unwrap();
    assert_eq!(result["meta"]["trace_id"], trace);
    assert_eq!(result["data"]["persistence"]["inserted_records"], 3);
    let first: serde_json::Value = client
        .get(format!(
            "{gateway}/api/v1/trade-events?account={ACCOUNT}&source=stored&limit=1"
        ))
        .send()
        .await
        .unwrap()
        .json()
        .await
        .unwrap();
    assert_eq!(first["data"]["coverage"], "STORED_RECORDS_ONLY");
    assert_eq!(first["data"]["matched_records"], 3);
    assert_eq!(count(&db, "raw_logs").await, 3);
    for query in [
        format!("account={ACCOUNT}&source=no"),
        format!("account={ACCOUNT}&start_time=2026-10-05T00:00:00Z"),
        format!("account={ACCOUNT}&source=stored&cursor=broken"),
        format!("account={ACCOUNT}&source=stored&start_time=bad"),
    ] {
        let r = client
            .get(format!("{gateway}/api/v1/trade-events?{query}"))
            .send()
            .await
            .unwrap();
        assert_eq!(r.status(), 400);
    }
    assert_eq!(
        client
            .get(format!("{internal}/internal/v1/health"))
            .send()
            .await
            .unwrap()
            .status(),
        401
    );
    db.pool.close().await;
    let r = client
        .get(format!(
            "{gateway}/api/v1/trade-events?account={ACCOUNT}&source=stored"
        ))
        .send()
        .await
        .unwrap();
    assert_eq!(r.status(), 503);
    assert_eq!(
        client
            .get(format!("{gateway}/api/v1/health"))
            .send()
            .await
            .unwrap()
            .status(),
        200
    );
    stop_g.send(()).unwrap();
    stop_i.send(()).unwrap();
    task_g.await.unwrap();
    task_i.await.unwrap();
    std::fs::remove_file(token_path).unwrap();
}
