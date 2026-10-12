use super::{error, handlers, middleware::trace_request};
use crate::state::AppState;
use axum::{Router, middleware, routing::get};

pub fn build(state: AppState) -> Router {
    Router::new()
        .route(
            "/api/v1/trade-events/live",
            get(handlers::trade_events::live),
        )
        .route(
            "/api/v1/trade-events/stored",
            get(handlers::trade_events::stored),
        )
        .route(
            "/api/v1/monitoring/accounts/add",
            axum::routing::post(handlers::watch_accounts::add),
        )
        .route(
            "/api/v1/monitoring/accounts/list",
            get(handlers::watch_accounts::handle),
        )
        .route(
            "/api/v1/monitoring/accounts/detail",
            get(handlers::watch_accounts::detail_by_id),
        )
        .route(
            "/api/v1/monitoring/accounts/pause",
            axum::routing::post(handlers::watch_accounts::pause_by_id),
        )
        .route(
            "/api/v1/monitoring/accounts/resume",
            axum::routing::post(handlers::watch_accounts::resume_by_id),
        )
        .route(
            "/api/v1/monitoring/accounts/collection-status",
            get(handlers::watch_accounts::collection_status),
        )
        .route(
            "/api/v1/monitoring/accounts/publishing-status",
            get(handlers::watch_accounts::publishing_status),
        )
        .route(
            "/api/v1/publishing/status",
            get(handlers::publishing::summary),
        )
        // Legacy routes remain compatible; new clients use function routes above.
        .route("/api/v1/health", get(handlers::health::handle))
        .route("/api/v1/version", get(handlers::version::handle))
        .route("/api/v1/trade-events", get(handlers::trade_events::handle))
        .route(
            "/api/v1/watch-accounts",
            get(handlers::watch_accounts::handle).post(handlers::watch_accounts::add),
        )
        .route(
            "/api/v1/watch-accounts/{id}",
            get(handlers::watch_accounts::detail).delete(handlers::watch_accounts::pause),
        )
        .route(
            "/api/v1/watch-accounts/{id}/resume",
            axum::routing::post(handlers::watch_accounts::resume),
        )
        .route(
            "/api/v1/publishing-status",
            get(handlers::publishing::handle),
        )
        .fallback(error::not_found)
        .method_not_allowed_fallback(error::method_not_allowed)
        .layer(middleware::from_fn(trace_request))
        .with_state(state)
}
