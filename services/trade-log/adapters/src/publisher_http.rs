use crate::postgres::Postgres;
use axum::{
    Json, Router,
    extract::{Request, State},
    http::StatusCode,
    middleware::{self, Next},
    response::{IntoResponse, Response},
    routing::get,
};
use serde_json::json;
use std::{sync::Arc, time::Duration};
#[derive(Clone)]
pub struct PublisherState {
    pub reader: Arc<Postgres>,
    pub account_key: String,
    pub credential: Arc<String>,
}
async fn authentication(State(s): State<PublisherState>, mut r: Request, next: Next) -> Response {
    let authorized = crate::internal_http::credential_matches(r.headers(), &s.credential);
    let trace = if authorized {
        r.headers()
            .get("x-trace-id")
            .and_then(|v| v.to_str().ok())
            .filter(|v| {
                v.len() == 38
                    && v.starts_with("trace_")
                    && v[6..].bytes().all(|b| b.is_ascii_hexdigit())
            })
            .map(str::to_owned)
    } else {
        None
    }
    .unwrap_or_else(|| format!("trace_{}", uuid::Uuid::new_v4().simple()));
    r.extensions_mut().insert(trace.clone());
    let method = r.method().clone();
    let started = std::time::Instant::now();
    let mut result = if authorized {
        next.run(r).await
    } else {
        (StatusCode::UNAUTHORIZED,Json(json!({"code":"AUTHENTICATION_REQUIRED","message":"Service authentication required","trace_id":trace}))).into_response()
    };
    result
        .headers_mut()
        .insert("x-trace-id", trace.parse().expect("trace ASCII"));
    tracing::info!(service="trade-parser-publisher",trace_id=%trace,%method,status=result.status().as_u16(),duration_ms=started.elapsed().as_secs_f64()*1000.,"request_completed");
    result
}
async fn status(
    State(s): State<PublisherState>,
    axum::Extension(trace): axum::Extension<String>,
    p: Result<axum::extract::Query<StatusParams>, axum::extract::rejection::QueryRejection>,
) -> Response {
    let p=match p{Ok(p)=>p.0,Err(_)=>return (StatusCode::BAD_REQUEST,Json(json!({"code":"INVALID_ARGUMENT","message":"Invalid query parameters","trace_id":trace}))).into_response()};
    let result = async {
        if let Some(account) = p.account {
            let account = shared_types::account(&account).map_err(trade_log::watch::invalid)?;
            s.reader
                .publishing_status(&s.reader.checkpoint_key(&account))
                .await
        } else if s.account_key.is_empty() {
            s.reader.publishing_list().await
        } else {
            s.reader.publishing_status(&s.account_key).await
        }
    };
    match tokio::time::timeout(Duration::from_secs(10),result).await {
  Ok(Ok(data))=>Json(json!({"data":data,"meta":{"trace_id":trace,"schema_version":1}})).into_response(),
  Ok(Err(e)) if e.code=="INVALID_ARGUMENT" => (StatusCode::BAD_REQUEST,Json(json!({"code":e.code,"message":e.message,"trace_id":trace}))).into_response(),
  _=>(StatusCode::SERVICE_UNAVAILABLE,Json(json!({"code":"DEPENDENCY_UNAVAILABLE","message":"Publishing status unavailable","trace_id":trace}))).into_response(),
 }
}
async fn health(axum::Extension(trace): axum::Extension<String>) -> Response {
    Json(json!({"data":{"status":"ok","service":"trade-parser-publisher","version":env!("CARGO_PKG_VERSION")},"meta":{"trace_id":trace,"schema_version":1}})).into_response()
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
pub fn router(state: PublisherState) -> Router {
    Router::new()
        .route("/internal/v1/health", get(health))
        .route("/internal/v1/publishing-status", get(status))
        .fallback(not_found)
        .method_not_allowed_fallback(method_not_allowed)
        .layer(middleware::from_fn_with_state(
            state.clone(),
            authentication,
        ))
        .with_state(state)
}

#[derive(serde::Deserialize)]
#[serde(deny_unknown_fields)]
struct StatusParams {
    account: Option<String>,
}
