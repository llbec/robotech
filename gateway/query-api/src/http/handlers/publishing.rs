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
#[derive(serde::Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Params {
    account: Option<String>,
}
pub async fn handle(
    State(s): State<AppState>,
    Extension(trace): Extension<TraceId>,
    p: Result<axum::extract::Query<Params>, axum::extract::rejection::QueryRejection>,
) -> Response {
    let result = async {
        let client = s.publisher.ok_or_else(|| {
            trade_log::query::QueryError::unavailable("Webhook publishing is not enabled")
        })?;
        let p = p
            .map_err(|_| trade_log::watch::invalid("Invalid query parameters"))?
            .0;
        let account = p
            .account
            .map(|a| shared_types::account(&a).map_err(trade_log::watch::invalid))
            .transpose()?;
        client.status(&trace.0, account.as_deref()).await
    }
    .await;
    match result {
        Ok(data) => response::success(data, trace).into_response(),
        Err(e) => error::response(
            StatusCode::from_u16(e.status()).unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
            &e.code,
            &e.message,
            trace,
        ),
    }
}
