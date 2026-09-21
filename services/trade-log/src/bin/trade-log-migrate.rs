use anyhow::Context;
use anyhow::Result;
use trade_log::Repository;

#[tokio::main]
async fn main() -> Result<()> {
    let database_url = std::env::var("DATABASE_URL").context("DATABASE_URL is required")?;
    let repo = Repository::connect(&database_url).await?;
    repo.migrate().await
}
