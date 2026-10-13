use clap::{Parser, Subcommand};
use std::{path::PathBuf, sync::Arc, time::Duration};
use tokio_util::sync::CancellationToken;
use trade_log_adapters::{
    collector_http::{self, CollectorState},
    collector_supervisor::CollectorSupervisor,
    postgres::Postgres,
};
#[derive(Parser)]
#[command(name="trade-collector",version=env!("CARGO_PKG_VERSION"))]
struct Args {
    #[command(subcommand)]
    command: Option<Command>,
    #[arg(long, default_value = "/etc/robotech/trade-collector.toml")]
    config: PathBuf,
}
#[derive(Subcommand)]
enum Command {
    BackfillSpot {
        #[arg(long)]
        account: String,
        #[arg(long)]
        start_time: String,
        #[arg(long)]
        end_time: String,
    },
    BackfillSpotStatus {
        #[arg(long)]
        job_id: uuid::Uuid,
    },
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
    store
        .monitoring_configure(&c.monitoring, &c.source_budget)
        .await
        .map_err(|e| (3, e.message))?;
    store
        .import_watch(&c.collection.account, &c.collection.start_time)
        .await
        .map_err(|e| (2, e.message))?;
    store.stream_memory = Some(Arc::new(tokio::sync::Semaphore::new(
        c.monitoring.max_total_pending_bytes,
    )));
    let source = hyperliquid::source::HttpSource::new(
        c.hyperliquid.network.endpoint(),
        Duration::from_secs(c.hyperliquid.connect_timeout_seconds),
        Duration::from_secs(c.hyperliquid.request_timeout_seconds),
        c.hyperliquid.max_response_bytes,
    )
    .map_err(|e| (3, e.message))?;
    if let Some(command) = args.command {
        let value = match command {
            Command::BackfillSpot {
                account,
                start_time,
                end_time,
            } => {
                let start =
                    trade_log::collection::start_ms(&start_time).map_err(|e| (2, e.message))?;
                let end = trade_log::collection::start_ms(&end_time).map_err(|e| (2, e.message))?;
                store
                    .submit_spot_backfill(&account, start, end)
                    .await
                    .map_err(|e| (3, e.message))?
            }
            Command::BackfillSpotStatus { job_id } => store
                .spot_backfill_status(job_id)
                .await
                .map_err(|e| (3, e.message))?,
        };
        println!("{value}");
        store.pool.close().await;
        return Ok(());
    }
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
    let realtime_source = Arc::new(hyperliquid::websocket::WebsocketSource::new(
        &c.hyperliquid.network,
    ));
    let websocket_config = c.websocket;
    let runtime = CollectorSupervisor {
        store: store.clone(),
        defaults: c.collection,
        monitoring: c.monitoring.clone(),
        websocket: websocket_config,
        realtime: Arc::new(trade_log_adapters::source_budget::BudgetRealtime {
            store: store.clone(),
            inner: realtime_source,
        }),
        source: Arc::new(trade_log_adapters::source_budget::BudgetSource {
            store: store.clone(),
            inner: Arc::new(source),
            concurrent: Arc::new(tokio::sync::Semaphore::new(
                c.monitoring.max_concurrent_http_requests,
            )),
        }),
        parser: Arc::new(hyperliquid::parser::HyperliquidParser),
    };
    let cancel = CancellationToken::new();
    let backfill_runtime = trade_log_adapters::spot_backfill_runtime::SpotBackfillRuntime {
        store: store.clone(),
        source: runtime.source.clone(),
        parser: runtime.parser.clone(),
    };
    let mut backfill_worker = tokio::spawn(backfill_runtime.run(cancel.clone()));
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
    if tokio::time::timeout(c.server.shutdown_timeout(), &mut backfill_worker)
        .await
        .is_err()
    {
        backfill_worker.abort();
        let _ = backfill_worker.await;
    }
    store.pool.close().await;
    result?;
    worker_result
}
