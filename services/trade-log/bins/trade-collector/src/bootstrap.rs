use clap::Parser;
use std::{path::PathBuf, sync::Arc, time::Duration};
use tokio_util::sync::CancellationToken;
use trade_log_adapters::{
    collector_http::{self, CollectorState},
    collector_runtime::CollectorRuntime,
    postgres::Postgres,
};
#[derive(Parser)]
#[command(name="trade-collector",version=env!("CARGO_PKG_VERSION"))]
struct Args {
    #[arg(long, default_value = "/etc/robotech/trade-collector.toml")]
    config: PathBuf,
}
pub async fn run() -> Result<(), (u8, String)> {
    let args = match Args::try_parse() {
        Ok(a) => a,
        Err(e)
            if matches!(
                e.kind(),
                clap::error::ErrorKind::DisplayHelp | clap::error::ErrorKind::DisplayVersion
            ) =>
        {
            e.print()
                .map_err(|_| (3, "Cannot write CLI output".into()))?;
            return Ok(());
        }
        Err(e) => return Err((2, e.to_string())),
    };
    let env = service_runtime::config::environment().map_err(|e| (2, e))?;
    let c = crate::config::Config::load(&args.config, &env).map_err(|e| (2, e))?;
    let token = std::fs::read_to_string(&c.internal.credential_file)
        .map_err(|_| (2, "Cannot read internal credential file".into()))?;
    let token = token.trim().to_owned();
    if token.len() < 32 || !token.bytes().all(|b| b.is_ascii_alphanumeric()) {
        return Err((2, "Invalid service credential".into()));
    }
    let url = std::fs::read_to_string(&c.database.url_file)
        .map_err(|_| (2, "Cannot read database credential file".into()))?;
    if !url.trim().starts_with("postgres://") && !url.trim().starts_with("postgresql://") {
        return Err((2, "Invalid PostgreSQL credential".into()));
    }
    service_runtime::logging::init(&c.logging).map_err(|e| (3, e.to_string()))?;
    let mut store = Postgres::connect(
        url.trim(),
        c.hyperliquid.network.clone(),
        c.database.max_connections,
        c.database.connect_timeout_seconds,
        c.database.statement_timeout_seconds,
    )
    .await
    .map_err(|e| (3, e.message))?;
    store.check_schema().await.map_err(|e| (3, e.message))?;
    let mirror = Arc::new(trade_log_adapters::file_evidence::FileEvidence {
        directory: c.evidence.directory.clone(),
        network: c.hyperliquid.network.clone(),
    });
    if let Err(e) = mirror.check().await {
        tracing::warn!(code=%e.code,"evidence_mirror_unavailable");
    }
    store.mirror = Some(mirror);
    // A persisted checkpoint takes precedence over the configuration start time.
    let exists: bool = sqlx_exists(&store, &c.collection.account).await?;
    if !exists
        && trade_log::collection::start_ms(&c.collection.start_time).map_err(|e| (2, e.message))?
            > chrono::Utc::now().timestamp_millis()
                - (c.collection.safety_delay_seconds * 1000) as i64
    {
        return Err((2, "Initial collection start is in the future".into()));
    }
    let source = hyperliquid::source::HttpSource::new(
        c.hyperliquid.network.endpoint(),
        Duration::from_secs(c.hyperliquid.connect_timeout_seconds),
        Duration::from_secs(c.hyperliquid.request_timeout_seconds),
        c.hyperliquid.max_response_bytes,
    )
    .map_err(|e| (3, e.message))?;
    let address = c.server.address();
    let signal = service_runtime::lifecycle::shutdown_signal().map_err(|e| (3, e.to_string()))?;
    let listener = tokio::net::TcpListener::bind(address)
        .await
        .map_err(|e| (3, format!("Cannot bind {address}: {e}")))?;
    let state = CollectorState {
        reader: Arc::new(store.clone()),
        account: c.collection.account.clone(),
        credential: Arc::new(token),
    };
    let runtime = CollectorRuntime {
        store: store.clone(),
        config: c.collection,
        source: Arc::new(source),
        parser: Arc::new(hyperliquid::parser::HyperliquidParser),
    };
    let cancel = CancellationToken::new();
    let worker_cancel = cancel.clone();
    let worker = tokio::spawn(async move {
        runtime.run(worker_cancel).await;
    });
    tracing::info!(service="trade-collector",version=env!("CARGO_PKG_VERSION"),%address,config_path=%args.config.display(),"server_started");
    let signal_cancel = cancel.clone();
    let result = service_runtime::lifecycle::serve(
        "trade-collector",
        listener,
        collector_http::router(state),
        async move {
            signal.await;
            signal_cancel.cancel();
        },
        c.server.shutdown_timeout(),
    )
    .await
    .map_err(|e| (4, e.to_string()));
    cancel.cancel();
    let mut worker = worker;
    let worker_result = match tokio::time::timeout(c.server.shutdown_timeout(), &mut worker).await {
        Ok(result) => result.map_err(|_| (4, "Collector worker failed".into())),
        Err(_) => {
            worker.abort();
            let _ = worker.await;
            Err((4, "Collector shutdown timed out".into()))
        }
    };
    store.pool.close().await;
    result?;
    worker_result
}
async fn sqlx_exists(store: &Postgres, account: &str) -> Result<bool, (u8, String)> {
    store
        .checkpoint_exists(account)
        .await
        .map_err(|e| (3, e.message))
}
