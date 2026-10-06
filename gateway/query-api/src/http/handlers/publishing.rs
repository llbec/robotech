use crate::{
    http::{error, middleware::TraceId, response},
    state::AppState,
};
use axum::{
    Extension,
    extract::State,
    http::StatusCode,
    response::{IntoResponse, Response},
};
pub async fn handle(State(s): State<AppState>, Extension(trace): Extension<TraceId>) -> Response {
    let result = async {
        let client = s.publisher.ok_or_else(|| {
            trade_log::query::QueryError::unavailable("Webhook publishing is not enabled")
        })?;
        client.status(&trace.0).await
    }
    .await;
    match result {
        Ok(data) => response::success(data, trace).into_response(),
        Err(e) => error::response(
            StatusCode::SERVICE_UNAVAILABLE,
            "DEPENDENCY_UNAVAILABLE",
            &e.message,
            trace,
        ),
    }
}
