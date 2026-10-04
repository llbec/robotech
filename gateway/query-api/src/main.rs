use std::process::ExitCode;

#[tokio::main]
async fn main() -> ExitCode {
    match query_api::bootstrap::run().await {
        Ok(()) => ExitCode::SUCCESS,
        Err(error) => {
            eprintln!("{error}");
            ExitCode::from(error.exit_code())
        }
    }
}
