use axum::{Router, routing::get};
use query_api::{
    config::Config,
    http::router,
    lifecycle::{self, RunError},
    state::AppState,
};
use serde_json::Value;
use std::{
    collections::BTreeMap,
    process::Command,
    sync::{
        Arc,
        atomic::{AtomicBool, Ordering},
    },
    time::Duration,
};
use tokio::{
    net::TcpListener,
    sync::{Notify, oneshot},
};

const VALID: &str = include_str!("../fixtures/config/valid.toml");

#[test]
fn configuration_contract() {
    let empty = BTreeMap::new();
    let config = Config::parse(VALID, &empty).unwrap();
    assert_eq!(config.server.port, 8080);
    for (from, to) in [
        ("config_version = 1", "config_version = 2"),
        ("port = 8080", "port = 0"),
        ("port = 8080", "port = 65536"),
        ("port = 8080", ""),
        ("level = \"info\"", "level = \"verbose\""),
        ("format = \"json\"", "format = \"xml\""),
        (
            "shutdown_timeout_seconds = 10",
            "shutdown_timeout_seconds = 0",
        ),
        (
            "shutdown_timeout_seconds = 10",
            "shutdown_timeout_seconds = 61",
        ),
        ("host = \"127.0.0.1\"", "host = \"localhost\""),
        ("port = 8080", "port = 8080\nunknown = true"),
    ] {
        assert!(
            Config::parse(&VALID.replace(from, to), &empty).is_err(),
            "{to}"
        );
    }
    assert!(Config::parse("not toml", &empty).is_err());
    let overrides = BTreeMap::from([
        ("ROBOTECH_SERVER_PORT".into(), "8081".into()),
        ("ROBOTECH_SERVER_HOST".into(), "::1".into()),
        ("ROBOTECH_LOG_LEVEL".into(), "debug".into()),
        ("ROBOTECH_LOG_FORMAT".into(), "text".into()),
        ("ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS".into(), "20".into()),
    ]);
    let config = Config::parse(VALID, &overrides).unwrap();
    assert_eq!(config.server.port, 8081);
    assert_eq!(config.server.host.to_string(), "::1");
    assert_eq!(config.logging.level, "debug");
    assert_eq!(config.logging.format, "text");
    assert_eq!(config.server.shutdown_timeout_seconds, 20);
    for key in query_api::config::ENV_KEYS {
        for value in ["", "invalid"] {
            assert!(Config::parse(VALID, &BTreeMap::from([(key.into(), value.into())])).is_err());
        }
    }
}

async fn start(
    app: Router,
    timeout: Duration,
) -> (
    String,
    std::net::SocketAddr,
    oneshot::Sender<()>,
    tokio::task::JoinHandle<Result<(), RunError>>,
) {
    let listener = TcpListener::bind("127.0.0.1:0").await.unwrap();
    let address = listener.local_addr().unwrap();
    let (stop, receive) = oneshot::channel();
    let task = tokio::spawn(lifecycle::serve(
        listener,
        app,
        async {
            let _ = receive.await;
        },
        timeout,
    ));
    (format!("http://{address}"), address, stop, task)
}

#[tokio::test]
async fn http_contract_and_trace_identity() {
    let (base, address, stop, task) =
        start(router::build(AppState::new()), Duration::from_secs(1)).await;
    let client = reqwest::Client::new();
    let mut startup = None;
    let mut traces = std::collections::HashSet::new();
    for route in ["health", "health", "version"] {
        let response = client
            .get(format!("{base}/api/v1/{route}"))
            .header("x-trace-id", "client-controlled")
            .send()
            .await
            .unwrap();
        assert_eq!(response.status(), 200);
        assert_eq!(response.headers()["content-type"], "application/json");
        let trace = response.headers()["x-trace-id"]
            .to_str()
            .unwrap()
            .to_owned();
        let body: Value = response.json().await.unwrap();
        assert_eq!(body["meta"]["trace_id"], trace);
        assert_eq!(body["meta"]["schema_version"], 1);
        assert_ne!(trace, "client-controlled");
        assert!(traces.insert(trace));
        assert_eq!(body["data"]["service"], "query-api");
        if route == "health" {
            assert_eq!(body["data"]["status"], "ok");
            let value = body["data"]["started_at"].as_str().unwrap().to_owned();
            chrono::DateTime::parse_from_rfc3339(&value).unwrap();
            assert!(value.ends_with('Z'));
            if let Some(previous) = startup.as_ref() {
                assert_eq!(previous, &value);
            }
            startup = Some(value);
        } else {
            assert_eq!(body["data"]["version"], query_api::VERSION);
        }
    }
    for route in ["health", "version"] {
        let response = client
            .post(format!("{base}/api/v1/{route}"))
            .send()
            .await
            .unwrap();
        assert_eq!(response.status(), 405);
        let allow = response.headers()["allow"].to_str().unwrap();
        assert!(allow.contains("GET") && allow.contains("HEAD"));
        let trace = response.headers()["x-trace-id"]
            .to_str()
            .unwrap()
            .to_owned();
        let body: Value = response.json().await.unwrap();
        assert_eq!(body["code"], "METHOD_NOT_ALLOWED");
        assert_eq!(body["trace_id"], trace);
        let response = client
            .head(format!("{base}/api/v1/{route}"))
            .send()
            .await
            .unwrap();
        assert_eq!(response.status(), 200);
        assert_eq!(response.headers()["content-type"], "application/json");
        assert!(response.bytes().await.unwrap().is_empty());
    }
    let response = client
        .get(format!("{base}/unknown?private=value"))
        .send()
        .await
        .unwrap();
    assert_eq!(response.status(), 404);
    let trace = response.headers()["x-trace-id"]
        .to_str()
        .unwrap()
        .to_owned();
    let body: Value = response.json().await.unwrap();
    assert_eq!(body["code"], "RESOURCE_NOT_FOUND");
    assert_eq!(body["trace_id"], trace);
    stop.send(()).unwrap();
    task.await.unwrap().unwrap();
    TcpListener::bind(address).await.unwrap();
}

struct DropFlag(Arc<AtomicBool>);
impl Drop for DropFlag {
    fn drop(&mut self) {
        self.0.store(true, Ordering::SeqCst);
    }
}

#[tokio::test]
async fn draining_completes_existing_request() {
    let entered = Arc::new(Notify::new());
    let release = Arc::new(Notify::new());
    let app = Router::new().route(
        "/slow",
        get({
            let entered = entered.clone();
            let release = release.clone();
            move || {
                let entered = entered.clone();
                let release = release.clone();
                async move {
                    entered.notify_one();
                    release.notified().await;
                    "finished"
                }
            }
        }),
    );
    let (base, address, stop, task) = start(app, Duration::from_secs(1)).await;
    let request = tokio::spawn(async move {
        reqwest::get(format!("{base}/slow"))
            .await
            .unwrap()
            .text()
            .await
            .unwrap()
    });
    entered.notified().await;
    stop.send(()).unwrap();
    // Leave the handler pending while graceful draining begins.
    tokio::time::sleep(Duration::from_millis(30)).await;
    assert!(!task.is_finished());
    release.notify_one();
    assert_eq!(request.await.unwrap(), "finished");
    task.await.unwrap().unwrap();
    TcpListener::bind(address).await.unwrap();
}

#[tokio::test]
async fn shutdown_timeout_cancels_inflight_work() {
    let entered = Arc::new(Notify::new());
    let dropped = Arc::new(AtomicBool::new(false));
    let app = Router::new().route(
        "/hang",
        get({
            let entered = entered.clone();
            let dropped = dropped.clone();
            move || {
                let entered = entered.clone();
                let dropped = dropped.clone();
                async move {
                    let _guard = DropFlag(dropped);
                    entered.notify_one();
                    std::future::pending::<()>().await;
                    "unreachable"
                }
            }
        }),
    );
    let (base, address, stop, task) = start(app, Duration::from_millis(50)).await;
    let request = tokio::spawn(async move { reqwest::get(format!("{base}/hang")).await });
    entered.notified().await;
    stop.send(()).unwrap();
    assert!(matches!(
        tokio::time::timeout(Duration::from_secs(1), task)
            .await
            .unwrap()
            .unwrap(),
        Err(RunError::ShutdownTimeout)
    ));
    tokio::time::timeout(Duration::from_secs(1), async {
        while !dropped.load(Ordering::SeqCst) {
            tokio::task::yield_now().await;
        }
    })
    .await
    .unwrap();
    TcpListener::bind(address).await.unwrap();
    let _ = request.await;
}

fn command() -> Command {
    let mut cmd = Command::new(env!("CARGO_BIN_EXE_query-api"));
    for key in query_api::config::ENV_KEYS {
        cmd.env_remove(key);
    }
    cmd
}

#[test]
fn cli_help_version_and_configuration_errors() {
    let help = command()
        .args(["--help", "--config", "/missing"])
        .output()
        .unwrap();
    assert!(help.status.success());
    assert!(String::from_utf8_lossy(&help.stdout).contains("--config"));
    let version = command()
        .args(["--version", "--config", "/missing"])
        .output()
        .unwrap();
    assert!(version.status.success());
    assert_eq!(
        String::from_utf8_lossy(&version.stdout).trim(),
        format!("query-api {}", query_api::VERSION)
    );
    let missing = command()
        .args(["--config", "/missing/robotech.toml"])
        .output()
        .unwrap();
    assert_eq!(missing.status.code(), Some(2));
    assert!(!String::from_utf8_lossy(&missing.stderr).contains("server_started"));
    assert_eq!(
        command().arg("--unknown").output().unwrap().status.code(),
        Some(2)
    );
    let path = std::path::Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("../../tests/fixtures/config/valid.toml");
    let invalid = command()
        .arg("--config")
        .arg(path)
        .env("ROBOTECH_SERVER_PORT", "")
        .output()
        .unwrap();
    assert_eq!(invalid.status.code(), Some(2));
}

#[tokio::test]
async fn occupied_port_is_startup_failure() {
    let listener = TcpListener::bind("127.0.0.1:0").await.unwrap();
    let path = std::path::Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("../../tests/fixtures/config/valid.toml");
    let result = command()
        .arg("--config")
        .arg(path)
        .env(
            "ROBOTECH_SERVER_PORT",
            listener.local_addr().unwrap().port().to_string(),
        )
        .output()
        .unwrap();
    assert_eq!(result.status.code(), Some(3));
    let stderr = String::from_utf8_lossy(&result.stderr);
    assert!(stderr.contains("bind_failed"));
    assert!(!stderr.contains("server_started"));
}

#[cfg(unix)]
struct ChildGuard(Option<std::process::Child>);
#[cfg(unix)]
impl Drop for ChildGuard {
    fn drop(&mut self) {
        if let Some(child) = self.0.as_mut() {
            let _ = child.kill();
            let _ = child.wait();
        }
    }
}

#[cfg(unix)]
#[tokio::test]
async fn os_signals_and_request_logs() {
    for signal in ["-TERM", "-INT"] {
        let listener = TcpListener::bind("127.0.0.1:0").await.unwrap();
        let address = listener.local_addr().unwrap();
        drop(listener);
        let path = std::path::Path::new(env!("CARGO_MANIFEST_DIR"))
            .join("../../tests/fixtures/config/valid.toml");
        let mut child = ChildGuard(Some(
            command()
                .arg("--config")
                .arg(path)
                .env("ROBOTECH_SERVER_PORT", address.port().to_string())
                .stdout(std::process::Stdio::null())
                .stderr(std::process::Stdio::piped())
                .spawn()
                .unwrap(),
        ));
        let client = reqwest::Client::builder()
            .timeout(Duration::from_secs(1))
            .build()
            .unwrap();
        let base = format!("http://{address}");
        let health = tokio::time::timeout(Duration::from_secs(5), async {
            loop {
                if let Ok(response) = client.get(format!("{base}/api/v1/health")).send().await {
                    break response;
                }
                assert!(child.0.as_mut().unwrap().try_wait().unwrap().is_none());
                tokio::time::sleep(Duration::from_millis(20)).await;
            }
        })
        .await
        .unwrap();
        let health_trace = health.headers()["x-trace-id"].to_str().unwrap().to_owned();
        let _ = health.bytes().await.unwrap();
        let response = client
            .get(format!("{base}/private-path?secret=do-not-log"))
            .header("Authorization", "Bearer do-not-log")
            .send()
            .await
            .unwrap();
        let error_trace = response.headers()["x-trace-id"]
            .to_str()
            .unwrap()
            .to_owned();
        let _ = response.bytes().await.unwrap();
        let pid = child.0.as_ref().unwrap().id().to_string();
        assert!(
            Command::new("/bin/kill")
                .args([signal, &pid])
                .status()
                .unwrap()
                .success()
        );
        tokio::time::timeout(Duration::from_secs(3), async {
            while child.0.as_mut().unwrap().try_wait().unwrap().is_none() {
                tokio::time::sleep(Duration::from_millis(10)).await;
            }
        })
        .await
        .unwrap();
        let output = child.0.take().unwrap().wait_with_output().unwrap();
        assert!(output.status.success());
        let logs = String::from_utf8(output.stderr).unwrap();
        assert!(!logs.contains("do-not-log") && !logs.contains("private-path"));
        let events: Vec<Value> = logs
            .lines()
            .map(|line| serde_json::from_str(line).unwrap())
            .collect();
        for message in ["server_started", "shutdown_started", "shutdown_completed"] {
            assert!(
                events
                    .iter()
                    .any(|event| event["fields"]["message"] == message)
            );
        }
        for (trace, status, level) in [(health_trace, 200, "INFO"), (error_trace, 404, "WARN")] {
            let event = events
                .iter()
                .find(|event| event["fields"]["trace_id"] == trace)
                .unwrap();
            assert_eq!(event["fields"]["status"], status);
            assert_eq!(event["level"], level);
            assert_eq!(event["fields"]["service"], "query-api");
            assert!(event["fields"]["duration_ms"].is_number());
        }
        TcpListener::bind(address).await.unwrap();
    }
}

#[tokio::test]
async fn trade_queries_validate_before_contacting_disabled_dependency() {
    let (base, _, stop, task) = start(router::build(AppState::new()), Duration::from_secs(1)).await;
    let client = reqwest::Client::new();
    let address = "0x0000000000000000000000000000000000000001";
    for query in [
        "account=invalid".to_owned(),
        format!("account={address}&limit=0"),
        format!("account={address}&start=1"),
        format!("account={address}&limit=no"),
        format!("account={address}&source=unknown"),
        format!("account={address}&source=live&cursor=invalid"),
        format!("account={address}&source=stored&start_time=invalid"),
        format!("account={address}&source=stored&cursor="),
    ] {
        let response = client
            .get(format!("{base}/api/v1/trade-events?{query}"))
            .send()
            .await
            .unwrap();
        assert_eq!(response.status(), 400);
        let body: Value = response.json().await.unwrap();
        assert_eq!(body["code"], "VALIDATION_ERROR");
    }
    assert_eq!(
        client
            .get(format!("{base}/api/v1/watch-accounts"))
            .send()
            .await
            .unwrap()
            .status(),
        503
    );
    let response = client
        .get(format!("{base}/api/v1/trade-events?account={address}"))
        .send()
        .await
        .unwrap();
    assert_eq!(response.status(), 503);
    assert_eq!(
        response.json::<Value>().await.unwrap()["code"],
        "DEPENDENCY_UNAVAILABLE"
    );
    stop.send(()).unwrap();
    task.await.unwrap().unwrap();
}
