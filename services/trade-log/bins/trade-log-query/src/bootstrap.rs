use clap::{Parser, Subcommand};
use hyperliquid::{parser::HyperliquidParser, source::HttpSource};
use std::{path::PathBuf, sync::Arc, time::Duration};
use tokio::{net::TcpListener, sync::Semaphore};
use trade_log::query::QueryService;
use trade_log_adapters::{
    file_evidence::FileEvidence,
    internal_http::{self, InternalState},
    postgres::Postgres,
};
#[derive(Parser)]
#[command(name="trade-log-query",version=env!("CARGO_PKG_VERSION"),about="Robotech internal trade query service")]
struct Args {
    #[arg(long, global = true, default_value = "/etc/robotech/trade-log.toml")]
    config: PathBuf,
    #[command(subcommand)]
    command: Option<Command>,
}
#[derive(Subcommand)]
enum Command {
    Migrate,
    Reparse {
        #[arg(long)]
        query_id: String,
    },
    ImportEvidence {
        #[arg(long)]
        directory: PathBuf,
    },
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
    let token = if args.command.is_none() {
        let token = std::fs::read_to_string(&config.internal.credential_file)
            .map_err(|_| error(2, "cannot read internal credential file"))?;
        let token = token.trim().to_owned();
        if token.len() < 32 || !token.bytes().all(|b| b.is_ascii_alphanumeric()) {
            return Err(error(
                2,
                "service credential must have at least 32 ASCII alphanumeric characters",
            ));
        }
        Some(token)
    } else {
        None
    };
    let file = if matches!(args.command, Some(Command::Migrate)) {
        config
            .database
            .migration_url_file
            .as_ref()
            .ok_or_else(|| error(2, "database.migration_url_file required"))?
    } else {
        &config.database.url_file
    };
    let url = std::fs::read_to_string(file)
        .map_err(|_| error(2, "cannot read database credential file"))?;
    if !url.trim().starts_with("postgres://") && !url.trim().starts_with("postgresql://") {
        return Err(error(2, "invalid PostgreSQL connection credential"));
    }
    let mut store = Postgres::connect(
        url.trim(),
        config.hyperliquid.network.clone(),
        config.database.max_connections,
        config.database.connect_timeout_seconds,
        config.database.statement_timeout_seconds,
    )
    .await
    .map_err(|e| error(3, e))?;
    if matches!(args.command, Some(Command::Migrate)) {
        store.migrate().await.map_err(|e| error(3, e))?;
        println!("migration_completed");
        return Ok(());
    }
    store.check_schema().await.map_err(|e| error(3, e))?;
    match args.command {
        Some(Command::Reparse { query_id }) => {
            if !query_id.starts_with("query_")
                || !query_id
                    .bytes()
                    .all(|b| b.is_ascii_alphanumeric() || b == b'_')
            {
                return Err(error(2, "invalid query_id"));
            }
            let report = match store.reparse(&query_id, &HyperliquidParser).await {
                Ok(report) => report,
                Err(e) if e.code == "INCOMPLETE_DATA" => {
                    println!(
                        "{}",
                        serde_json::json!({"query_id":query_id,"comparison":"INCOMPLETE","code":e.code})
                    );
                    return Err(error(4, e));
                }
                Err(e) => return Err(error(3, e)),
            };
            println!(
                "{}",
                serde_json::to_string_pretty(&report).map_err(|e| error(3, e))?
            );
            return if report.comparison == "SAME" {
                Ok(())
            } else {
                Err(error(4, "reparse_difference"))
            };
        }
        Some(Command::ImportEvidence { directory }) => {
            let mut entries = tokio::fs::read_dir(directory)
                .await
                .map_err(|_| error(2, "cannot read evidence directory"))?;
            let (mut imported, mut existing, mut failed) = (0, 0, 0);
            let mut storage_failed = false;
            let mut paths = Vec::new();
            while let Some(entry) = entries
                .next_entry()
                .await
                .map_err(|_| error(3, "cannot read evidence entry"))?
            {
                if entry
                    .file_type()
                    .await
                    .map_err(|_| error(3, "cannot inspect evidence entry"))?
                    .is_dir()
                {
                    paths.push(entry.path());
                }
            }
            paths.sort();
            for path in paths {
                match store.import_directory(&path, &HyperliquidParser).await {
                    Ok(true) => {
                        imported += 1;
                        println!(
                            "{}",
                            serde_json::json!({"directory":path.file_name().map(|n|n.to_string_lossy()),"status":"IMPORTED"})
                        );
                    }
                    Ok(false) => {
                        existing += 1;
                        println!(
                            "{}",
                            serde_json::json!({"directory":path.file_name().map(|n|n.to_string_lossy()),"status":"EXISTING"})
                        );
                    }
                    Err(e) => {
                        failed += 1;
                        storage_failed |= e.code == "DEPENDENCY_UNAVAILABLE";
                        println!(
                            "{}",
                            serde_json::json!({"directory":path.file_name().map(|n|n.to_string_lossy()),"status":"FAILED","code":e.code})
                        );
                    }
                }
            }
            println!(
                "{}",
                serde_json::json!({"imported":imported,"existing":existing,"failed":failed})
            );
            return if failed == 0 {
                Ok(())
            } else {
                Err(error(
                    if storage_failed { 3 } else { 4 },
                    "some evidence directories failed",
                ))
            };
        }
        Some(Command::Migrate) => unreachable!(),
        None => {}
    }
    store.interrupt().await.map_err(|e| error(3, e))?;
    let mirror = Arc::new(FileEvidence {
        directory: config.evidence.directory,
        network: config.hyperliquid.network.clone(),
    });
    if let Err(e) = mirror.check().await {
        tracing::warn!(code=%e.code,"evidence_mirror_unavailable");
    }
    store.mirror = Some(mirror);
    let source = HttpSource::new(
        config.hyperliquid.network.endpoint(),
        Duration::from_secs(config.hyperliquid.connect_timeout_seconds),
        Duration::from_secs(config.hyperliquid.request_timeout_seconds),
        config.query.max_response_bytes,
    )
    .map_err(|e| error(3, e))?;
    let store = Arc::new(store);
    let service = Arc::new(QueryService {
        network: config.hyperliquid.network,
        source: Arc::new(source),
        parser: Arc::new(HyperliquidParser),
        evidence: store.clone(),
    });
    let state = InternalState {
        service,
        stored: Some(store.clone()),
        credential: Arc::new(token.expect("service token checked")),
        permits: Arc::new(Semaphore::new(config.query.max_concurrency)),
        timeout: Duration::from_secs(config.query.timeout_seconds),
    };
    let signal = service_runtime::lifecycle::shutdown_signal().map_err(|e| error(3, e))?;
    let address = config.server.address();
    let listener = TcpListener::bind(address)
        .await
        .map_err(|e| error(3, format!("cannot bind {address}: {e}")))?;
    tracing::info!(service="trade-log-query",version=env!("CARGO_PKG_VERSION"),%address,config_path=%args.config.display(),"server_started");
    let result = service_runtime::lifecycle::serve(
        "trade-log-query",
        listener,
        internal_http::router(state),
        signal,
        config.server.shutdown_timeout(),
    )
    .await
    .map_err(|e| error(4, e));
    store.pool.close().await;
    result
}
