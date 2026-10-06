use clap::{Parser, Subcommand};
use std::{path::PathBuf, sync::Arc};
use tokio_util::sync::CancellationToken;
use trade_log_adapters::{
    postgres::Postgres,
    publisher_http::{self, PublisherState},
    publisher_runtime::PublisherRuntime,
    webhook::Webhook,
};
#[derive(Parser)]
#[command(name="trade-parser-publisher",version=env!("CARGO_PKG_VERSION"))]
struct Args {
    #[arg(
        long,
        global = true,
        default_value = "/etc/robotech/trade-parser-publisher.toml"
    )]
    config: PathBuf,
    #[command(subcommand)]
    command: Option<Command>,
}
#[derive(Subcommand)]
enum Command {
    Retry {
        #[arg(long)]
        event_id: String,
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
    let (account, network) = c.account().map_err(|e| (2, e))?;
    let key = format!("hyperliquid:{}:hyperliquid:{}", network.name(), account);
    let token = read_token(&c.internal.credential_file)?;
    let url = std::fs::read_to_string(&c.database.url_file)
        .map_err(|_| (2, "Cannot read database credential".into()))?;
    if !url.trim().starts_with("postgresql://") && !url.trim().starts_with("postgres://") {
        return Err((2, "Invalid database credential".into()));
    }
    let webhook = if c.publishing.enabled {
        Some(Arc::new(
            Webhook::new(
                &c.publishing.webhook_url,
                &read_token(&c.publishing.credential_file)?,
                c.publishing.request_timeout_seconds,
                c.publishing.max_response_bytes,
            )
            .map_err(|e| (2, e))?,
        )
            as Arc<dyn trade_log::publishing::PublicationTarget>)
    } else {
        None
    };
    service_runtime::logging::init(&c.logging)
        .map_err(|_| (3, "Logging initialization failed".into()))?;
    let store = Postgres::connect(
        url.trim(),
        network,
        c.database.max_connections,
        c.database.connect_timeout_seconds,
        c.database.statement_timeout_seconds,
    )
    .await
    .map_err(|e| (3, e.message))?;
    store.check_schema().await.map_err(|e| (3, e.message))?;
    if let Some(Command::Retry { event_id }) = args.command {
        store
            .retry_publication(&key, &event_id)
            .await
            .map_err(|e| (3, e.message))?;
        println!(
            "{}",
            serde_json::json!({"event_id":event_id,"status":"RETRY_WAIT"})
        );
        store.pool.close().await;
        return Ok(());
    }
    store
        .configure_publishing(
            &key,
            &c.publishing.webhook_url,
            c.publishing.enabled,
            &c.policy(),
        )
        .await
        .map_err(|e| (3, format!("{}: {}", e.code, e.message)))?;
    let listener = tokio::net::TcpListener::bind(c.server.address())
        .await
        .map_err(|_| (3, "Cannot bind publisher address".into()))?;
    let signal = service_runtime::lifecycle::shutdown_signal()
        .map_err(|_| (3, "Cannot initialize shutdown signal".into()))?;
    let cancel = CancellationToken::new();
    let runtime = PublisherRuntime {
        store: store.clone(),
        key: key.clone(),
        webhook,
        lease_seconds: c.publishing.lease_seconds,
        poll_interval_ms: c.publishing.poll_interval_ms,
        retry_base_seconds: c.publishing.retry_base_seconds,
        retry_max_seconds: c.publishing.retry_max_seconds,
    };
    let mut worker = tokio::spawn(runtime.run(cancel.clone()));
    let signal_cancel = cancel.clone();
    tracing::info!(service="trade-parser-publisher",version=env!("CARGO_PKG_VERSION"),address=%c.server.address(),"server_started");
    let result = service_runtime::lifecycle::serve(
        "trade-parser-publisher",
        listener,
        publisher_http::router(PublisherState {
            reader: Arc::new(store.clone()),
            account_key: key,
            credential: Arc::new(token),
        }),
        async move {
            signal.await;
            signal_cancel.cancel();
        },
        c.server.shutdown_timeout(),
    )
    .await
    .map_err(|_| (4, "Publisher HTTP shutdown failed".into()));
    cancel.cancel();
    let drained = match tokio::time::timeout(c.server.shutdown_timeout(), &mut worker).await {
        Ok(r) => r.map_err(|_| (4, "Publisher worker failed".into())),
        Err(_) => {
            worker.abort();
            let _ = worker.await;
            Err((4, "Publisher shutdown timed out".into()))
        }
    };
    store.pool.close().await;
    result?;
    drained
}
fn read_token(path: &std::path::Path) -> Result<String, (u8, String)> {
    let s = std::fs::read_to_string(path).map_err(|_| (2, "Cannot read credential file".into()))?;
    let s = s.trim();
    if s.len() < 32 || !s.bytes().all(|b| b.is_ascii_alphanumeric()) {
        return Err((2, "Invalid credential".into()));
    }
    Ok(s.into())
}
