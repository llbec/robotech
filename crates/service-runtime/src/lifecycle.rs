use axum::{
    Router,
    extract::Request,
    middleware::{self, Next},
    response::IntoResponse,
};
use std::{future::Future, time::Duration};
use tokio::{net::TcpListener, sync::oneshot};
use tokio_util::sync::CancellationToken;

#[derive(Debug, thiserror::Error)]
pub enum RunError {
    #[error("HTTP service failed: {0}")]
    Io(#[from] std::io::Error),
    #[error("shutdown timeout")]
    ShutdownTimeout,
}

// Signal handlers are installed before binding, so registration errors are startup failures.
#[cfg(unix)]
pub fn shutdown_signal() -> std::io::Result<impl Future<Output = ()> + Send> {
    use tokio::signal::unix::{SignalKind, signal};
    let mut interrupt = signal(SignalKind::interrupt())?;
    let mut terminate = signal(SignalKind::terminate())?;
    Ok(async move {
        tokio::select! { _ = interrupt.recv() => {}, _ = terminate.recv() => {} }
    })
}

#[cfg(not(unix))]
pub fn shutdown_signal() -> std::io::Result<impl Future<Output = ()> + Send> {
    Ok(async {
        let _ = tokio::signal::ctrl_c().await;
    })
}

pub async fn serve<F>(
    service: &'static str,
    listener: TcpListener,
    router: Router,
    shutdown: F,
    timeout: Duration,
) -> Result<(), RunError>
where
    F: Future<Output = ()> + Send + 'static,
{
    let cancel = CancellationToken::new();
    let request_cancel = cancel.clone();
    // Axum connections are spawned independently. Cancel handlers as well as dropping
    // the serve future on timeout, so in-flight work cannot continue indefinitely.
    let router = router.layer(middleware::from_fn(move |request: Request, next: Next| {
        let cancel = request_cancel.clone();
        async move {
            tokio::select! {
                response = next.run(request) => response,
                _ = cancel.cancelled() => axum::http::StatusCode::SERVICE_UNAVAILABLE.into_response(),
            }
        }
    }));
    let (started_tx, mut started_rx) = oneshot::channel();
    let graceful = async move {
        shutdown.await;
        tracing::info!(service, "shutdown_started");
        let _ = started_tx.send(());
    };
    let server = axum::serve(listener, router)
        .with_graceful_shutdown(graceful)
        .into_future();
    tokio::pin!(server);
    tokio::select! {
        result = &mut server => { result?; }
        _ = &mut started_rx => {
            match tokio::time::timeout(timeout, &mut server).await {
                Ok(result) => { result?; }
                Err(_) => {
                    cancel.cancel();
                    tracing::error!(service, "shutdown_timeout");
                    return Err(RunError::ShutdownTimeout);
                }
            }
        }
    }
    tracing::info!(service, "shutdown_completed");
    Ok(())
}
