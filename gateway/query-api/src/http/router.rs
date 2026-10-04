use super::{error, handlers, middleware::trace_request};
use crate::state::AppState;
use axum::{Router, middleware, routing::get};

pub fn build(state: AppState) -> Router {
    Router::new()
        .route("/api/v1/health", get(handlers::health::handle))
        .route("/api/v1/version", get(handlers::version::handle))
        .fallback(error::not_found)
        .method_not_allowed_fallback(error::method_not_allowed)
        .layer(middleware::from_fn(trace_request))
        .with_state(state)
}
