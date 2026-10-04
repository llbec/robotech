use crate::{
    http::{error, middleware::TraceId, response},
    state::AppState,
};
use axum::{
    Extension,
    extract::{Query, State, rejection::QueryRejection},
    http::StatusCode,
    response::{IntoResponse, Response},
};
use serde::Deserialize;
use trade_log::{query::QueryError, validation::QueryRequest};
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Parameters {
    account: String,
    limit: Option<usize>,
}

pub async fn handle(
    State(state): State<AppState>,
    Extension(trace): Extension<TraceId>,
    parameters: Result<Query<Parameters>, QueryRejection>,
) -> Response {
    let request = match parameters {
        Ok(Query(p)) => {
            let request = QueryRequest {
                account: p.account,
                limit: p.limit.unwrap_or(100),
            };
            request.validated().map(|_| request)
        }
        Err(_) => Err(QueryError::validation("Invalid query parameters")),
    };
    let result = match request {
        Err(e) => Err(e),
        Ok(request) => match state.trade_log {
            Some(client) => client.query(&request, &trace.0).await,
            None => Err(QueryError::unavailable("Trade queries are not enabled")),
        },
    };
    match result {
        Ok(data) => response::success(data, trace).into_response(),
        Err(e) => {
            let code = match e.code.as_str() {
                "VALIDATION_ERROR" => "VALIDATION_ERROR",
                "INCOMPLETE_DATA" => "INCOMPLETE_DATA",
                "VERSION_CONFLICT" => "VERSION_CONFLICT",
                "RATE_LIMITED" => "RATE_LIMITED",
                "INTERNAL_INVARIANT_VIOLATION" => "INTERNAL_INVARIANT_VIOLATION",
                _ => "DEPENDENCY_UNAVAILABLE",
            };
            error::response(
                StatusCode::from_u16(e.status()).expect("error status"),
                code,
                &e.message,
                trace,
            )
        }
    }
}
