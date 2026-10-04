use clap::Parser;
use hyperliquid::{parser::HyperliquidParser, source::HttpSource};
use std::{path::PathBuf, sync::Arc, time::Duration};
use tokio::{net::TcpListener, sync::Semaphore};
use trade_log::query::QueryService;
use trade_log_adapters::{
    file_evidence::FileEvidence,
    internal_http::{self, InternalState},
};

#[derive(Parser)]
#[command(name="trade-log-query",version=env!("CARGO_PKG_VERSION"),about="Robotech internal trade query service")]
struct Args {
    #[arg(long, default_value = "config/trade-log.toml")]
    config: PathBuf,
}
#[derive(Debug, thiserror::Error)]
#[error("{message}")]
pub struct StartupError {
    pub code: u8,
    pub message: String,
}
fn error(code: u8, message: impl ToString) -> StartupError {
    StartupError {
        code,
        message: message.to_string(),
    }
}
pub async fn run() -> Result<(), StartupError> {
    let args = match Args::try_parse() {
        Ok(args) => args,
        Err(e)
            if matches!(
                e.kind(),
                clap::error::ErrorKind::DisplayHelp | clap::error::ErrorKind::DisplayVersion
            ) =>
        {
            e.print().map_err(|e| error(3, e))?;
            return Ok(());
        }
        Err(e) => return Err(error(2, e)),
    };
    let env = service_runtime::config::environment().map_err(|e| error(2, e))?;
    let config = crate::config::Config::load(&args.config, &env).map_err(|e| error(2, e))?;
    service_runtime::logging::init(&config.logging).map_err(|e| error(3, e))?;
    let token = std::fs::read_to_string(&config.internal.credential_file)
        .map_err(|_| error(2, "cannot read internal credential file"))?;
    let token = token.trim();
    if token.len() < 32 || !token.bytes().all(|b| b.is_ascii_alphanumeric()) {
        return Err(error(
            2,
            "service credential must have at least 32 ASCII alphanumeric characters",
        ));
    }
    let evidence = Arc::new(FileEvidence {
        directory: config.evidence.directory,
        network: config.hyperliquid.network.clone(),
    });
    evidence
        .check()
        .await
        .map_err(|_| error(3, "evidence directory is not writable"))?;
    let source = HttpSource::new(
        config.hyperliquid.network.endpoint(),
        Duration::from_secs(config.hyperliquid.connect_timeout_seconds),
        Duration::from_secs(config.hyperliquid.request_timeout_seconds),
        config.query.max_response_bytes,
    )
    .map_err(|e| error(3, e))?;
    let service = Arc::new(QueryService {
        network: config.hyperliquid.network,
        source: Arc::new(source),
        parser: Arc::new(HyperliquidParser),
        evidence,
    });
    let state = InternalState {
        service,
        credential: Arc::new(token.into()),
        permits: Arc::new(Semaphore::new(config.query.max_concurrency)),
        timeout: Duration::from_secs(config.query.timeout_seconds),
    };
    let signal = service_runtime::lifecycle::shutdown_signal().map_err(|e| error(3, e))?;
    let address = config.server.address();
    let listener = TcpListener::bind(address)
        .await
        .map_err(|e| error(3, format!("cannot bind {address}: {e}")))?;
    tracing::info!(service="trade-log-query",version=env!("CARGO_PKG_VERSION"),%address,config_path=%args.config.display(),"server_started");
    service_runtime::lifecycle::serve(
        "trade-log-query",
        listener,
        internal_http::router(state),
        signal,
        config.server.shutdown_timeout(),
    )
    .await
    .map_err(|e| error(4, e))
}
