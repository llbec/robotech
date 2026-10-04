use std::process::ExitCode;
#[tokio::main]
async fn main() -> ExitCode {
    match trade_log_query::bootstrap::run().await {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            eprintln!("{e}");
            ExitCode::from(e.code)
        }
    }
}
