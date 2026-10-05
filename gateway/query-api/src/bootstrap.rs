use crate::{config, http, lifecycle, logging, state::AppState};
use clap::Parser;
use std::path::PathBuf;
use tokio::net::TcpListener;

#[derive(Parser)]
#[command(name = "query-api", version = crate::VERSION, about = "Robotech query HTTP gateway")]
struct Args {
    #[arg(long, default_value = "config/query-api.toml")]
    config: PathBuf,
}

#[derive(Debug, thiserror::Error)]
#[error("{message}")]
pub struct StartupError {
    message: String,
    code: u8,
}
impl StartupError {
    pub fn exit_code(&self) -> u8 {
        self.code
    }
    fn new(code: u8, message: impl ToString) -> Self {
        Self {
            message: message.to_string(),
            code,
        }
    }
}

pub async fn run() -> Result<(), StartupError> {
    let args = match Args::try_parse() {
        Ok(args) => args,
        Err(error)
            if matches!(
                error.kind(),
                clap::error::ErrorKind::DisplayHelp | clap::error::ErrorKind::DisplayVersion
            ) =>
        {
            error.print().map_err(|e| StartupError::new(3, e))?;
            return Ok(());
        }
        Err(error) => return Err(StartupError::new(2, error)),
    };
    let overrides = config::environment().map_err(|e| StartupError::new(2, e))?;
    let config =
        config::Config::load(&args.config, &overrides).map_err(|e| StartupError::new(2, e))?;
    logging::init(&config.logging)
        .map_err(|e| StartupError::new(3, format!("logging initialization failed: {e}")))?;
    let mut state = AppState::new();
    state.trade_log = config
        .trade_log
        .as_ref()
        .map(crate::clients::trade_log::TradeLogClient::from_config)
        .transpose()
        .map_err(|e| StartupError::new(2, e))?
        .flatten();
    state.collector = config
        .collector
        .as_ref()
        .map(crate::clients::collector::CollectorClient::from_config)
        .transpose()
        .map_err(|e| StartupError::new(2, e))?
        .flatten();
    let signal = lifecycle::shutdown_signal()
        .map_err(|e| StartupError::new(3, format!("signal initialization failed: {e}")))?;
    let address = config.server.address();
    let listener = TcpListener::bind(address).await.map_err(|e| {
        tracing::error!(service = crate::SERVICE, %address, error = %e, "bind_failed");
        StartupError::new(3, format!("cannot bind {address}: {e}"))
    })?;

    tracing::info!(service = crate::SERVICE, version = crate::VERSION, %address,
        config_path = %args.config.display(), log_level = %config.logging.level,
        log_format = %config.logging.format, shutdown_timeout_seconds = config.server.shutdown_timeout_seconds,
        "server_started");
    lifecycle::serve(
        listener,
        http::router::build(state),
        signal,
        config.server.shutdown_timeout(),
    )
    .await
    .map_err(|e| {
        tracing::error!(service = crate::SERVICE, error = %e, "service_failed");
        StartupError::new(4, e)
    })
}
