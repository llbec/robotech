use axum::{Json, Router, extract::State, routing::post};
use hyperliquid::{parser::HyperliquidParser, source::HttpSource};
use serde_json::{Value, json};
use shared_types::Network;
use std::{collections::BTreeMap, sync::Arc, time::Duration};
use tokio::{
    net::TcpListener,
    sync::{Semaphore, oneshot},
};
use trade_log::{query::QueryService, validation::QueryRequest};
use trade_log_adapters::{
    file_evidence::FileEvidence,
    internal_http::{self, InternalState},
};
const ADDRESS: &str = "0x0000000000000000000000000000000000000001";
const TOKEN: &str = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";
const FILLS: &str = include_str!("../../../../tests/fixtures/hyperliquid/fills.json");
const META: &str = include_str!("../../../../tests/fixtures/hyperliquid/meta.json");
const SPOT: &str = include_str!("../../../../tests/fixtures/hyperliquid/spot-meta.json");
async fn start(router: Router) -> (String, oneshot::Sender<()>, tokio::task::JoinHandle<()>) {
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
async fn mock(State(fills): State<Arc<Value>>, Json(body): Json<Value>) -> Json<Value> {
    match body["type"].as_str().unwrap() {
        "userFills" => {
            assert_eq!(body["aggregateByTime"], false);
            Json((*fills).clone())
        }
        "meta" => Json(serde_json::from_str(META).unwrap()),
        "spotMeta" => Json(serde_json::from_str(SPOT).unwrap()),
        _ => panic!("unknown request"),
    }
}
struct Directory(std::path::PathBuf);
impl Drop for Directory {
    fn drop(&mut self) {
        let _ = std::fs::remove_dir_all(&self.0);
    }
}

#[tokio::test]
async fn gateway_to_source_roundtrip_with_persistent_evidence() {
    let dir =
        Directory(std::env::temp_dir().join(format!("robotech-test-{}", uuid::Uuid::new_v4())));
    let fills: Value = serde_json::from_str(FILLS).unwrap();
    let (source_url, source_stop, source_task) = start(
        Router::new()
            .route("/info", post(mock))
            .with_state(Arc::new(fills)),
    )
    .await;
    let evidence = Arc::new(FileEvidence {
        directory: dir.0.join("evidence"),
        network: Network::Mainnet,
    });
    evidence.check().await.unwrap();
    let service = Arc::new(QueryService {
        network: Network::Mainnet,
        source: Arc::new(
            HttpSource::new(
                &format!("{source_url}/info"),
                Duration::from_secs(1),
                Duration::from_secs(1),
                100000,
            )
            .unwrap(),
        ),
        parser: Arc::new(HyperliquidParser),
        evidence: evidence.clone(),
    });
    let state = InternalState {
        stored: None,
        service,
        credential: Arc::new(TOKEN.into()),
        permits: Arc::new(Semaphore::new(1)),
        timeout: Duration::from_secs(3),
    };
    let (internal_url, internal_stop, internal_task) = start(internal_http::router(state)).await;
    let token_path = dir.0.join("token");
    std::fs::write(&token_path, TOKEN).unwrap();
    let config = query_api::config::TradeLogConfig {
        enabled: true,
        base_url: Some(internal_url.clone()),
        credential_file: Some(token_path),
        request_timeout_seconds: Some(5),
    };
    let mut state = query_api::state::AppState::new();
    state.trade_log = query_api::clients::trade_log::TradeLogClient::from_config(&config).unwrap();
    let (gateway_url, gateway_stop, gateway_task) =
        start(query_api::http::router::build(state)).await;
    let client = reqwest::Client::new();
    let result = client
        .get(format!(
            "{gateway_url}/api/v1/trade-events?account={ADDRESS}&limit=1"
        ))
        .send()
        .await
        .unwrap();
    assert_eq!(result.status(), 200);
    let trace = result.headers()["x-trace-id"].to_str().unwrap().to_owned();
    let body: Value = result.json().await.unwrap();
    assert_eq!(body["meta"]["trace_id"], trace);
    assert_eq!(body["data"]["counts"]["perpetual_records"], 1);
    assert_eq!(body["data"]["trades"][0]["payload"]["notional"], "250.025");
    let id = body["data"]["query_id"].as_str().unwrap();
    let directory = dir.0.join("evidence").join(id);
    let manifest: Value =
        serde_json::from_slice(&std::fs::read(directory.join("manifest.json")).unwrap()).unwrap();
    assert_eq!(manifest["trace_id"], trace);
    assert_eq!(manifest["status"], "COMPLETED");
    assert_eq!(manifest["raw_log_metadata"].as_array().unwrap().len(), 3);
    let raw_id = body["data"]["trades"][0]["raw_log_id"].as_str().unwrap();
    assert_eq!(
        std::fs::read(directory.join("responses").join(format!("{raw_id}.body"))).unwrap(),
        serde_json::to_vec(&serde_json::from_str::<Value>(FILLS).unwrap()).unwrap()
    );
    let unauth = client
        .post(format!("{internal_url}/internal/v1/trade-queries"))
        .json(&json!({"account":ADDRESS}))
        .send()
        .await
        .unwrap();
    assert_eq!(unauth.status(), 401);
    let invalid = client
        .post(format!("{internal_url}/internal/v1/trade-queries"))
        .bearer_auth(TOKEN)
        .json(&json!({"account":ADDRESS,"unknown":true}))
        .send()
        .await
        .unwrap();
    assert_eq!(invalid.status(), 400);
    // Replaying evidence yields the same business fact; transport query IDs differ.
    use trade_log::parsing::{ParseContext, ProtocolParser};
    let replay = HyperliquidParser
        .parse(ParseContext {
            network: &Network::Mainnet,
            account: ADDRESS,
            raw_log_id: raw_id,
            fills: FILLS.as_bytes(),
            meta: META.as_bytes(),
            spot_meta: SPOT.as_bytes(),
        })
        .unwrap();
    assert_eq!(
        serde_json::to_value(&replay.trades[0]).unwrap(),
        body["data"]["trades"][0]
    );
    for stop in [gateway_stop, internal_stop, source_stop] {
        stop.send(()).unwrap();
    }
    for task in [gateway_task, internal_task, source_task] {
        task.await.unwrap();
    }
}

#[tokio::test]
async fn filtering_counts_and_display_limits() {
    for (fills, coverage, empty) in [
        (json!([]), "SOURCE_WINDOW_UNVERIFIED", true),
        (
            json!([{"coin":"@1","tid":1}]),
            "SOURCE_WINDOW_UNVERIFIED",
            true,
        ),
        (json!([{"coin":"unknown","tid":1}]), "LIMITED", true),
    ] {
        let dir =
            Directory(std::env::temp_dir().join(format!("robotech-test-{}", uuid::Uuid::new_v4())));
        let (url, stop, task) = start(
            Router::new()
                .route("/info", post(mock))
                .with_state(Arc::new(fills)),
        )
        .await;
        let evidence = Arc::new(FileEvidence {
            directory: dir.0.clone(),
            network: Network::Mainnet,
        });
        evidence.check().await.unwrap();
        let service = QueryService {
            network: Network::Mainnet,
            source: Arc::new(
                HttpSource::new(
                    &format!("{url}/info"),
                    Duration::from_secs(1),
                    Duration::from_secs(1),
                    100000,
                )
                .unwrap(),
            ),
            parser: Arc::new(HyperliquidParser),
            evidence,
        };
        let result = service
            .execute(
                "query_test",
                &QueryRequest {
                    account: ADDRESS.into(),
                    limit: 1,
                },
                "trace_test",
            )
            .await
            .unwrap();
        assert_eq!(result.coverage, coverage);
        assert_eq!(result.trades.is_empty(), empty);
        stop.send(()).unwrap();
        task.await.unwrap();
    }
    // Validate that the production deployment configuration is loadable; old gateway config remains valid.
    let path =
        std::path::Path::new(env!("CARGO_MANIFEST_DIR")).join("../../../config/query-api.toml");
    query_api::config::Config::load(&path, &BTreeMap::new()).unwrap();
}

#[derive(Clone)]
struct FailureSource {
    status: u16,
    body: String,
    calls: Arc<std::sync::atomic::AtomicUsize>,
    delay: Duration,
}
async fn failure_source(
    State(state): State<FailureSource>,
    Json(body): Json<Value>,
) -> axum::response::Response {
    use axum::response::IntoResponse;
    use std::sync::atomic::Ordering;
    if body["type"] == "userFills" {
        state.calls.fetch_add(1, Ordering::SeqCst);
        tokio::time::sleep(state.delay).await;
        (
            axum::http::StatusCode::from_u16(state.status).unwrap(),
            [("retry-after", "1")],
            state.body,
        )
            .into_response()
    } else {
        Json(
            serde_json::from_str::<Value>(if body["type"] == "meta" { META } else { SPOT })
                .unwrap(),
        )
        .into_response()
    }
}

#[tokio::test]
async fn source_failures_preserve_responses_and_bound_retries() {
    use std::sync::atomic::{AtomicUsize, Ordering};
    for (status, body, expected, attempts) in [
        (429, "limited", "RATE_LIMITED", 2),
        (500, "failed", "DEPENDENCY_UNAVAILABLE", 2),
        (400, "bad request", "DEPENDENCY_UNAVAILABLE", 1),
        (200, "not json", "INCOMPLETE_DATA", 1),
    ] {
        let dir =
            Directory(std::env::temp_dir().join(format!("robotech-test-{}", uuid::Uuid::new_v4())));
        let calls = Arc::new(AtomicUsize::new(0));
        let (url, stop, task) = start(
            Router::new()
                .route("/info", post(failure_source))
                .with_state(FailureSource {
                    status,
                    body: body.into(),
                    calls: calls.clone(),
                    delay: Duration::ZERO,
                }),
        )
        .await;
        let evidence = Arc::new(FileEvidence {
            directory: dir.0.clone(),
            network: Network::Mainnet,
        });
        evidence.check().await.unwrap();
        let service = QueryService {
            network: Network::Mainnet,
            source: Arc::new(
                HttpSource::new(
                    &format!("{url}/info"),
                    Duration::from_secs(1),
                    Duration::from_secs(1),
                    100000,
                )
                .unwrap(),
            ),
            parser: Arc::new(HyperliquidParser),
            evidence,
        };
        let result = service
            .execute(
                "query_error",
                &QueryRequest {
                    account: ADDRESS.into(),
                    limit: 1,
                },
                "trace_error",
            )
            .await;
        assert_eq!(result.unwrap_err().code, expected);
        assert_eq!(calls.load(Ordering::SeqCst), attempts);
        let manifest: Value = serde_json::from_slice(
            &std::fs::read(dir.0.join("query_error/manifest.json")).unwrap(),
        )
        .unwrap();
        assert_eq!(manifest["status"], "FAILED");
        assert_eq!(manifest["error"]["code"], expected);
        assert!(
            dir.0
                .join("query_error/responses/raw_query_error_userFills_1.body")
                .is_file()
        );
        stop.send(()).unwrap();
        task.await.unwrap();
    }
}

#[tokio::test]
async fn timeout_capacity_and_response_size_limits() {
    use std::sync::atomic::{AtomicUsize, Ordering};
    let dir =
        Directory(std::env::temp_dir().join(format!("robotech-test-{}", uuid::Uuid::new_v4())));
    let calls = Arc::new(AtomicUsize::new(0));
    let (url, source_stop, source_task) = start(
        Router::new()
            .route("/info", post(failure_source))
            .with_state(FailureSource {
                status: 200,
                body: "x".repeat(10000),
                calls: calls.clone(),
                delay: Duration::from_millis(100),
            }),
    )
    .await;
    let evidence = Arc::new(FileEvidence {
        directory: dir.0.clone(),
        network: Network::Mainnet,
    });
    evidence.check().await.unwrap();
    let service = Arc::new(QueryService {
        network: Network::Mainnet,
        source: Arc::new(
            HttpSource::new(
                &format!("{url}/info"),
                Duration::from_secs(1),
                Duration::from_secs(1),
                1000,
            )
            .unwrap(),
        ),
        parser: Arc::new(HyperliquidParser),
        evidence,
    });
    let permits = Arc::new(Semaphore::new(1));
    let state = InternalState {
        stored: None,
        service: service.clone(),
        credential: Arc::new(TOKEN.into()),
        permits: permits.clone(),
        timeout: Duration::from_millis(20),
    };
    let (internal_url, stop, task) = start(internal_http::router(state)).await;
    let client = reqwest::Client::new();
    let permit = permits.acquire().await.unwrap();
    let result = client
        .post(format!("{internal_url}/internal/v1/trade-queries"))
        .bearer_auth(TOKEN)
        .json(&json!({"account":ADDRESS}))
        .send()
        .await
        .unwrap();
    assert_eq!(result.status(), 429);
    assert_eq!(calls.load(Ordering::SeqCst), 0);
    drop(permit);
    let result = client
        .post(format!("{internal_url}/internal/v1/trade-queries"))
        .bearer_auth(TOKEN)
        .json(&json!({"account":ADDRESS}))
        .send()
        .await
        .unwrap();
    assert_eq!(result.status(), 503);
    assert_eq!(permits.available_permits(), 1);
    let request = QueryRequest {
        account: ADDRESS.into(),
        limit: 1,
    };
    let result = service.execute("query_size", &request, "trace_test").await;
    assert_eq!(result.unwrap_err().code, "DEPENDENCY_UNAVAILABLE");
    stop.send(()).unwrap();
    task.await.unwrap();
    source_stop.send(()).unwrap();
    source_task.await.unwrap();
}

#[tokio::test]
async fn full_evidence_survives_display_truncation_and_all_invalid_fails() {
    let base: Value = serde_json::from_str(FILLS).unwrap();
    for invalid in [false, true] {
        let mut fills = vec![];
        for i in 0..2000 {
            let mut f = base[0].clone();
            f["tid"] = (i + 1).into();
            if invalid {
                f["px"] = "bad".into();
            }
            fills.push(f);
        }
        let dir =
            Directory(std::env::temp_dir().join(format!("robotech-test-{}", uuid::Uuid::new_v4())));
        let (url, stop, task) = start(
            Router::new()
                .route("/info", post(mock))
                .with_state(Arc::new(json!(fills))),
        )
        .await;
        let evidence = Arc::new(FileEvidence {
            directory: dir.0.clone(),
            network: Network::Mainnet,
        });
        evidence.check().await.unwrap();
        let service = QueryService {
            network: Network::Mainnet,
            source: Arc::new(
                HttpSource::new(
                    &format!("{url}/info"),
                    Duration::from_secs(1),
                    Duration::from_secs(1),
                    2000000,
                )
                .unwrap(),
            ),
            parser: Arc::new(HyperliquidParser),
            evidence,
        };
        let result = service
            .execute(
                "query_limit",
                &QueryRequest {
                    account: ADDRESS.into(),
                    limit: 1,
                },
                "trace_test",
            )
            .await;
        if invalid {
            assert_eq!(result.unwrap_err().code, "INCOMPLETE_DATA");
        } else {
            let result = result.unwrap();
            assert_eq!(result.coverage, "LIMITED");
            assert!(result.display_truncated);
            assert_eq!(result.trades.len(), 1);
            assert_eq!(result.counts.perpetual_records, 2000);
            let full: Value = serde_json::from_slice(
                &std::fs::read(dir.0.join("query_limit/result.json")).unwrap(),
            )
            .unwrap();
            assert_eq!(full["trades"].as_array().unwrap().len(), 2000);
        }
        stop.send(()).unwrap();
        task.await.unwrap();
    }
}

#[tokio::test]
#[ignore = "explicit live Hyperliquid acceptance; requires internet access"]
async fn live_public_account_acceptance() {
    let address = std::env::var("ROBOTECH_ACCEPTANCE_ACCOUNT")
        .unwrap_or_else(|_| "0x010461c14e146ac35fe42271bdc1134ee31c703a".into());
    let dir = Directory(
        std::env::temp_dir().join(format!("robotech-live-token-{}", uuid::Uuid::new_v4())),
    );
    std::fs::create_dir_all(&dir.0).unwrap();
    let evidence_path =
        std::path::Path::new(env!("CARGO_MANIFEST_DIR")).join("../../../var/acceptance/v0.1");
    let evidence = Arc::new(FileEvidence {
        directory: evidence_path.clone(),
        network: Network::Mainnet,
    });
    evidence.check().await.unwrap();
    let service = Arc::new(QueryService {
        network: Network::Mainnet,
        source: Arc::new(
            HttpSource::new(
                Network::Mainnet.endpoint(),
                Duration::from_secs(5),
                Duration::from_secs(10),
                16777216,
            )
            .unwrap(),
        ),
        parser: Arc::new(HyperliquidParser),
        evidence,
    });
    let state = InternalState {
        stored: None,
        service,
        credential: Arc::new(TOKEN.into()),
        permits: Arc::new(Semaphore::new(1)),
        timeout: Duration::from_secs(30),
    };
    let (internal_url, internal_stop, internal_task) = start(internal_http::router(state)).await;
    let token = dir.0.join("token");
    std::fs::write(&token, TOKEN).unwrap();
    let mut state = query_api::state::AppState::new();
    state.trade_log = query_api::clients::trade_log::TradeLogClient::from_config(
        &query_api::config::TradeLogConfig {
            enabled: true,
            base_url: Some(internal_url),
            credential_file: Some(token),
            request_timeout_seconds: Some(35),
        },
    )
    .unwrap();
    let (gateway_url, gateway_stop, gateway_task) =
        start(query_api::http::router::build(state)).await;
    let client = reqwest::Client::builder()
        .timeout(Duration::from_secs(40))
        .build()
        .unwrap();
    let response = client
        .get(format!("{gateway_url}/api/v1/trade-events"))
        .query(&[("account", address.as_str()), ("limit", "100")])
        .send()
        .await
        .unwrap();
    let status = response.status();
    let value: Value = response.json().await.unwrap();
    gateway_stop.send(()).unwrap();
    internal_stop.send(()).unwrap();
    gateway_task.await.unwrap();
    internal_task.await.unwrap();
    assert_eq!(status, 200, "{value}");
    let id = value["data"]["query_id"].as_str().unwrap();
    let directory = evidence_path.join(id);
    let records = value["data"]["trades"].as_array().unwrap();
    assert!(
        !records.is_empty(),
        "account had no usable recent perpetual records"
    );
    let mut checked = 0;
    for fact in records {
        let raw_id = fact["raw_log_id"].as_str().unwrap();
        let source: Value = serde_json::from_slice(
            &std::fs::read(directory.join("responses").join(format!("{raw_id}.body"))).unwrap(),
        )
        .unwrap();
        let index = fact["payload"]["extension"]["source_indices"][0]
            .as_u64()
            .unwrap() as usize;
        let fill = &source[index];
        assert_eq!(
            fact["source_ref"],
            fill["tid"].as_u64().unwrap().to_string()
        );
        for (source_key, standard_key) in [("px", "price"), ("sz", "quantity"), ("fee", "fee")] {
            assert_eq!(
                shared_types::decimal(fill[source_key].as_str().unwrap()).unwrap(),
                shared_types::decimal(fact["payload"][standard_key].as_str().unwrap()).unwrap()
            );
        }
        assert_eq!(fact["payload"]["extension"]["time_ms"], fill["time"]);
        assert_eq!(
            fact["payload"]["fee_asset"].as_str(),
            fill["feeToken"].as_str().map(str::trim)
        );
        assert_eq!(
            fact["payload"]["side"],
            if fill["side"] == "B" { "BUY" } else { "SELL" }
        );
        assert_eq!(fact["payload"]["base_asset"], fill["coin"]);
        checked += 1;
    }
    println!(
        "live acceptance: account={address} source_records={} checked={checked} coverage={} evidence={}",
        value["data"]["counts"]["source_records"],
        value["data"]["coverage"],
        directory.display()
    );
}

#[tokio::test]
async fn unwritable_evidence_is_not_reported_as_empty_success() {
    let dir =
        Directory(std::env::temp_dir().join(format!("robotech-test-{}", uuid::Uuid::new_v4())));
    std::fs::create_dir(&dir.0).unwrap();
    let file = dir.0.join("not-a-directory");
    std::fs::write(&file, "occupied").unwrap();
    let evidence = Arc::new(FileEvidence {
        directory: file,
        network: Network::Mainnet,
    });
    assert_eq!(
        evidence.check().await.unwrap_err().code,
        "INTERNAL_INVARIANT_VIOLATION"
    );
    let service = QueryService {
        network: Network::Mainnet,
        source: Arc::new(
            HttpSource::new(
                "http://127.0.0.1:1",
                Duration::from_secs(1),
                Duration::from_secs(1),
                1000,
            )
            .unwrap(),
        ),
        parser: Arc::new(HyperliquidParser),
        evidence,
    };
    let result = service
        .execute(
            "query_storage",
            &QueryRequest {
                account: ADDRESS.into(),
                limit: 1,
            },
            "trace_test",
        )
        .await;
    assert_eq!(result.unwrap_err().status(), 500);
}
