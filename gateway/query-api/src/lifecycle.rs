use axum::Router;
pub use service_runtime::lifecycle::{RunError, shutdown_signal};
use std::{future::Future, time::Duration};
use tokio::net::TcpListener;

pub async fn serve<F>(
    listener: TcpListener,
    router: Router,
    shutdown: F,
    timeout: Duration,
) -> Result<(), RunError>
where
    F: Future<Output = ()> + Send + 'static,
{
    service_runtime::lifecycle::serve(crate::SERVICE, listener, router, shutdown, timeout).await
}
