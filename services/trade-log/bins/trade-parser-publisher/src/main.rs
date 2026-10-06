#[tokio::main]
async fn main() -> std::process::ExitCode {
    match trade_parser_publisher::bootstrap::run().await {
        Ok(()) => std::process::ExitCode::SUCCESS,
        Err((code, message)) => {
            eprintln!("{message}");
            std::process::ExitCode::from(code)
        }
    }
}
