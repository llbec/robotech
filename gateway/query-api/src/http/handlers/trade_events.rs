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
use trade_log::stored_query::StoredRequest;
use trade_log::{query::QueryError, validation::QueryRequest};
#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Parameters {
    #[serde(default)]
    instrument_type: trade_log::instrument::InstrumentType,
    account: String,
    limit: Option<usize>,
    source: Option<String>,
    start_time: Option<String>,
    end_time: Option<String>,
    cursor: Option<String>,
}

pub async fn handle(
    State(state): State<AppState>,
    Extension(trace): Extension<TraceId>,
    parameters: Result<Query<Parameters>, QueryRejection>,
) -> Response {
    let result = async {
        let Query(p) =
            parameters.map_err(|_| QueryError::validation("Invalid query parameters"))?;
        let source = p.source.as_deref().unwrap_or("live");
        match source {
            "live" => {
                if p.start_time.is_some() || p.end_time.is_some() || p.cursor.is_some() {
                    return Err(QueryError::validation(
                        "Time and cursor require source=stored",
                    ));
                }
                let request = QueryRequest {
                    instrument_type: p.instrument_type,
                    account: p.account,
                    limit: p.limit.unwrap_or(100),
                };
                request.validated()?;
                let client = state
                    .trade_log
                    .ok_or_else(|| QueryError::unavailable("Trade queries are not enabled"))?;
                Ok(serde_json::json!(client.query(&request, &trace.0).await?))
            }
            "stored" => {
                let request = StoredRequest {
                    instrument_type: p.instrument_type,
                    account: p.account,
                    limit: p.limit.unwrap_or(100),
                    start_time: p.start_time,
                    end_time: p.end_time,
                    cursor: p.cursor,
                }
                .validated()?;
                // Reject empty cursors here; decoding and request binding belong to the owner service.
                if request.cursor.as_ref().is_some_and(|c| c.is_empty()) {
                    return Err(QueryError::validation("Invalid cursor"));
                }
                let client = state
                    .trade_log
                    .ok_or_else(|| QueryError::unavailable("Trade queries are not enabled"))?;
                Ok(serde_json::json!(client.stored(&request, &trace.0).await?))
            }
            _ => Err(QueryError::validation("source must be live or stored")),
        }
    }
    .await;
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

pub async fn live(
    State(state): State<AppState>,
    Extension(trace): Extension<TraceId>,
    parameters: Result<Query<Parameters>, QueryRejection>,
) -> Response {
    by_function(state, trace, parameters, "live").await
}
pub async fn stored(
    State(state): State<AppState>,
    Extension(trace): Extension<TraceId>,
    parameters: Result<Query<Parameters>, QueryRejection>,
) -> Response {
    by_function(state, trace, parameters, "stored").await
}
async fn by_function(
    state: AppState,
    trace: TraceId,
    parameters: Result<Query<Parameters>, QueryRejection>,
    source: &str,
) -> Response {
    match parameters {
        Ok(Query(mut p)) => {
            if p.source.is_some() {
                return error::response(
                    StatusCode::BAD_REQUEST,
                    "VALIDATION_ERROR",
                    "source is not accepted; select /trade-events/live or /trade-events/stored",
                    trace,
                );
            }
            p.source = Some(source.into());
            handle(State(state), Extension(trace), Ok(Query(p))).await
        }
        Err(e) => handle(State(state), Extension(trace), Err(e)).await,
    }
}
