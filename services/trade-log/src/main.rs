use anyhow::{Context, Result};
use std::sync::Arc;
use tracing::info;
use tracing_subscriber::EnvFilter;
use trade_log::{
    acquisition::HyperliquidClient,
    api::{router, AppState},
    market_metadata::MetadataStore,
    publishing, Config, MonitorManager, Repository,
};

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::try_from_default_env().unwrap_or_else(|_| "info".into()))
        .init();
    let config = Config::from_env()?;
    let repo = Repository::connect(&config.database_url).await?;
    repo.migrate().await?;
    let client = HyperliquidClient::new(
        config.hyperliquid_http_url.clone(),
        config.hyperliquid_ws_url.clone(),
    )?;
    let metadata = Arc::new(MetadataStore::empty());
    let monitors = MonitorManager::new(repo.clone(), client, metadata, config.clone());
    if let Err(error) = monitors.refresh_metadata().await {
        tracing::warn!(%error, "live market metadata unavailable; trying persisted metadata");
        if !monitors.load_persisted_metadata().await? {
            return Err(error).context("no persisted market metadata is available");
        }
    }
    monitors.restore().await?;

    let publisher_repo = repo.clone();
    let publisher_config = config.clone();
    tokio::spawn(async move { publishing::run(publisher_repo, publisher_config).await });

    let metadata_manager = monitors.clone();
    tokio::spawn(async move {
        let mut interval = tokio::time::interval(std::time::Duration::from_secs(3600));
        interval.tick().await;
        loop {
            interval.tick().await;
            if let Err(error) = metadata_manager.refresh_metadata().await {
                tracing::warn!(%error, "market metadata refresh failed");
            }
        }
    });

    let state = Arc::new(AppState {
        repo,
        monitors,
        config: config.clone(),
    });
    let listener = tokio::net::TcpListener::bind(config.listen_addr).await?;
    info!(address = %config.listen_addr, "trade-log-server listening");
    axum::serve(listener, router(state))
        .with_graceful_shutdown(shutdown())
        .await?;
    Ok(())
}

async fn shutdown() {
    #[cfg(unix)]
    {
        use tokio::signal::unix::{signal, SignalKind};
        let mut terminate = signal(SignalKind::terminate()).expect("install SIGTERM handler");
        tokio::select! {
            _ = tokio::signal::ctrl_c() => {},
            _ = terminate.recv() => {},
        }
    }
    #[cfg(not(unix))]
    let _ = tokio::signal::ctrl_c().await;
    info!("shutdown signal received");
}
