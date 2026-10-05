use axum::{
    Json, Router,
    extract::{MatchedPath, Request, State, rejection::JsonRejection},
    http::{HeaderMap, StatusCode},
    middleware::{self, Next},
    response::{IntoResponse, Response},
    routing::{get, post},
};
use serde_json::json;
use sha2::{Digest, Sha256};
use std::{
    sync::Arc,
    time::{Duration, Instant},
};
use subtle::ConstantTimeEq;
use tokio::sync::Semaphore;
use trade_log::{
    query::{QueryError, QueryService},
    validation::QueryRequest,
};

#[derive(Clone)]
pub struct InternalState {
    pub service: Arc<QueryService>,
    pub stored: Option<Arc<dyn trade_log::stored_query::StoredQuery>>,
    pub credential: Arc<String>,
    pub permits: Arc<Semaphore>,
    pub timeout: Duration,
}
pub fn credential_matches(headers: &HeaderMap, token: &str) -> bool {
    let Some(value) = headers
        .get("authorization")
        .and_then(|v| v.to_str().ok())
        .and_then(|v| v.strip_prefix("Bearer "))
    else {
        return false;
    };
    let provided = Sha256::digest(value.as_bytes());
    let expected = Sha256::digest(token.as_bytes());
    bool::from(provided.as_slice().ct_eq(expected.as_slice()))
}
fn trace(headers: &HeaderMap) -> String {
    headers
        .get("x-trace-id")
        .and_then(|v| v.to_str().ok())
        .filter(|v| {
            v.len() == 38
                && v.starts_with("trace_")
                && v[6..].bytes().all(|b| b.is_ascii_hexdigit())
        })
        .map(str::to_owned)
        .unwrap_or_else(|| format!("trace_{}", uuid::Uuid::new_v4().simple()))
}
fn error_response(error: QueryError, trace: &str) -> Response {
    (
        StatusCode::from_u16(error.status()).expect("valid status"),
        Json(json!({"code":error.code,"message":error.message,"trace_id":trace})),
    )
        .into_response()
}
async fn authenticated(
    State(state): State<InternalState>,
    mut request: Request,
    next: Next,
) -> Response {
    let authorized = credential_matches(request.headers(), &state.credential);
    // Only authenticated callers can propagate a trace ID.
    let trace = if authorized {
        trace(request.headers())
    } else {
        format!("trace_{}", uuid::Uuid::new_v4().simple())
    };
    request.extensions_mut().insert(trace.clone());
    let method = request.method().clone();
    let route = request
        .extensions()
        .get::<MatchedPath>()
        .map_or("unmatched", |v| v.as_str())
        .to_owned();
    let started = Instant::now();
    let mut response = if authorized {
        next.run(request).await
    } else {
        (StatusCode::UNAUTHORIZED,Json(json!({"code":"AUTHENTICATION_REQUIRED","message":"Service authentication required","trace_id":trace}))).into_response()
    };
    response
        .headers_mut()
        .insert("x-trace-id", trace.parse().expect("trace ASCII"));
    let status = response.status().as_u16();
    let duration_ms = started.elapsed().as_secs_f64() * 1000.;
    if status >= 500 {
        tracing::error!(service="trade-log-query",%trace,trace_id=%trace,%method,%route,status,duration_ms,"request_completed");
    } else if status >= 400 {
        tracing::warn!(service="trade-log-query",trace_id=%trace,%method,%route,status,duration_ms,"request_completed");
    } else {
        tracing::info!(service="trade-log-query",trace_id=%trace,%method,%route,status,duration_ms,"request_completed");
    }
    response
}
async fn query(
    State(state): State<InternalState>,
    axum::Extension(trace): axum::Extension<String>,
    input: Result<Json<QueryRequest>, JsonRejection>,
) -> Response {
    let request = match input {
        Ok(Json(request)) => match request.validated() {
            Ok(_) => request,
            Err(e) => return error_response(e, &trace),
        },
        Err(_) => return error_response(QueryError::validation("Invalid query body"), &trace),
    };
    let Ok(_permit) = state.permits.try_acquire() else {
        return error_response(QueryError::limited(), &trace);
    };
    let id = format!("query_{}", uuid::Uuid::new_v4().simple());
    let result =
        tokio::time::timeout(state.timeout, state.service.execute(&id, &request, &trace)).await;
    match result {
        Ok(Ok(result)) => Json(json!({"data":result,"meta":{"trace_id":trace,"schema_version":1}}))
            .into_response(),
        Ok(Err(error)) => error_response(error, &trace),
        Err(_) => {
            let error = QueryError::unavailable("Query timed out");
            // Best-effort bounded diagnostic closure after cancellation; no detached query work.
            let _ = tokio::time::timeout(
                Duration::from_millis(250),
                state.service.evidence.finish(&id, Err(&error)),
            )
            .await;
            error_response(error, &trace)
        }
    }
}
async fn stored_query(
    State(state): State<InternalState>,
    axum::Extension(trace): axum::Extension<String>,
    input: Result<Json<trade_log::stored_query::StoredRequest>, JsonRejection>,
) -> Response {
    let result = async {
        let Json(request) = input.map_err(|_| QueryError::validation("Invalid query body"))?;
        let request = request.validated()?;
        let store = state
            .stored
            .ok_or_else(|| QueryError::unavailable("Stored queries unavailable"))?;
        tokio::time::timeout(state.timeout, store.stored(&request))
            .await
            .map_err(|_| QueryError::unavailable("Stored query timed out"))?
    }
    .await;
    match result {
        Ok(data) => {
            Json(json!({"data":data,"meta":{"trace_id":trace,"schema_version":1}})).into_response()
        }
        Err(e) => error_response(e, &trace),
    }
}
async fn health(axum::Extension(trace): axum::Extension<String>) -> Response {
    Json(json!({"data":{"status":"ok","service":"trade-log-query","version":env!("CARGO_PKG_VERSION")},"meta":{"trace_id":trace,"schema_version":1}})).into_response()
}
async fn not_found(axum::Extension(trace): axum::Extension<String>) -> Response {
    (
        StatusCode::NOT_FOUND,
        Json(json!({"code":"RESOURCE_NOT_FOUND","message":"Route not found","trace_id":trace})),
    )
        .into_response()
}
async fn method_not_allowed(axum::Extension(trace): axum::Extension<String>) -> Response {
    (
        StatusCode::METHOD_NOT_ALLOWED,
        Json(json!({"code":"METHOD_NOT_ALLOWED","message":"Method not allowed","trace_id":trace})),
    )
        .into_response()
}
pub fn router(state: InternalState) -> Router {
    Router::new()
        .route("/internal/v1/trade-queries", post(query))
        .route("/internal/v1/stored-trade-queries", post(stored_query))
        .route("/internal/v1/health", get(health))
        .fallback(not_found)
        .method_not_allowed_fallback(method_not_allowed)
        .layer(axum::extract::DefaultBodyLimit::max(16384))
        .layer(middleware::from_fn_with_state(state.clone(), authenticated))
        .with_state(state)
}
