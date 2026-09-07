use anyhow::Context;
use chrono::{DateTime, Utc};
use clap::Parser;
use std::{env, io::Write, process::ExitCode};
use tracing::error;
use tracing_subscriber::EnvFilter;
use trade_log::{ImportJob, ImportRequest};
use trade_log_adapters::{nansen::NansenClient, postgres::PostgresTradeRepository};

#[derive(Parser)]
#[command(version, about = "Import Hyperliquid perpetual trades from Nansen")]
struct Args {
    #[arg(long)]
    address: String,
    #[arg(long, value_parser = parse_time)]
    from: DateTime<Utc>,
    #[arg(long, value_parser = parse_time)]
    to: DateTime<Utc>,
}
fn parse_time(value: &str) -> Result<DateTime<Utc>, String> {
    value
        .parse()
        .map_err(|e| format!("invalid UTC ISO 8601 timestamp: {e}"))
}

#[tokio::main]
async fn main() -> ExitCode {
    tracing_subscriber::fmt()
        .with_env_filter(
            EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info")),
        )
        .with_writer(std::io::stderr)
        .init();
    let args = match Args::try_parse() {
        Ok(args) => args,
        Err(error) => {
            let _ = error.print();
            return ExitCode::from(2);
        }
    };
    match execute(args).await {
        Ok(()) => ExitCode::SUCCESS,
        Err(RunError::Configuration(error)) => {
            error!(error = %error, "configuration failed");
            ExitCode::from(2)
        }
        Err(RunError::Runtime(error)) => {
            error!(error = %format!("{error:#}"), "import failed");
            ExitCode::from(1)
        }
    }
}

enum RunError {
    Configuration(anyhow::Error),
    Runtime(anyhow::Error),
}

async fn execute(args: Args) -> Result<(), RunError> {
    trade_log::validation::validate_address(&args.address)
        .map_err(|e| RunError::Configuration(e.into()))?;
    trade_log::validation::validate_range(args.from, args.to)
        .map_err(|e| RunError::Configuration(e.into()))?;
    let api_key = env::var("NANSEN_API_KEY")
        .context("NANSEN_API_KEY is required")
        .map_err(RunError::Configuration)?;
    let database_url = env::var("DATABASE_URL")
        .context("DATABASE_URL is required")
        .map_err(RunError::Configuration)?;
    let source = NansenClient::new(api_key)
        .map_err(anyhow::Error::msg)
        .map_err(RunError::Configuration)?;
    let repository = PostgresTradeRepository::connect(&database_url)
        .await
        .map_err(anyhow::Error::msg)
        .map_err(RunError::Runtime)?;
    repository
        .migrate()
        .await
        .map_err(anyhow::Error::msg)
        .map_err(RunError::Runtime)?;
    let summary = ImportJob::new(source, repository)
        .run(ImportRequest {
            address: args.address,
            from: args.from,
            to: args.to,
        })
        .await
        .map_err(anyhow::Error::msg)
        .map_err(RunError::Runtime)?;
    let stdout = std::io::stdout();
    let mut output = std::io::BufWriter::new(stdout.lock());
    for fact in summary.facts {
        serde_json::to_writer(&mut output, &fact)
            .map_err(anyhow::Error::from)
            .map_err(RunError::Runtime)?;
        writeln!(output)
            .map_err(anyhow::Error::from)
            .map_err(RunError::Runtime)?;
    }
    output
        .flush()
        .map_err(anyhow::Error::from)
        .map_err(RunError::Runtime)?;
    Ok(())
}
