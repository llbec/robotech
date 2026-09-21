use crate::{
    config::Config,
    domain::validate_address,
    monitoring::MonitorManager,
    repository::{RawTradeView, TradeFactView},
    Repository,
};
use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::{get, post},
    Json, Router,
};
use base64::{engine::general_purpose::URL_SAFE_NO_PAD, Engine};
use chrono::{DateTime, TimeDelta, Utc};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use std::sync::Arc;

#[derive(Clone)]
pub struct AppState {
    pub repo: Repository,
    pub monitors: MonitorManager,
    pub config: Config,
}

pub fn router(state: Arc<AppState>) -> Router {
    Router::new()
        .route("/health/live", get(live))
        .route("/health/ready", get(ready))
        .route(
            "/api/v1/monitored-addresses",
            post(create_monitor).get(list_monitors),
        )
        .route("/api/v1/monitored-addresses/{address}", get(get_monitor))
        .route(
            "/api/v1/monitored-addresses/{address}/raw-trades",
            get(raw_trades),
        )
        .route("/api/v1/monitored-addresses/{address}/trades", get(trades))
        .route("/api/v1/webhook-deliveries", get(deliveries))
        .with_state(state)
}

async fn live() -> Json<Value> {
    Json(json!({"status":"ok"}))
}

async fn ready(State(state): State<Arc<AppState>>) -> Response {
    if state.repo.ready().await {
        (StatusCode::OK, Json(json!({"status":"ready"}))).into_response()
    } else {
        (
            StatusCode::SERVICE_UNAVAILABLE,
            Json(json!({"status":"not_ready"})),
        )
            .into_response()
    }
}

#[derive(Deserialize)]
struct CreateMonitorRequest {
    address: String,
    start_time: Option<DateTime<Utc>>,
}

async fn create_monitor(
    State(state): State<Arc<AppState>>,
    Json(request): Json<CreateMonitorRequest>,
) -> Result<Response, ApiError> {
    let address = validate_address(&request.address).ok_or_else(|| {
        ApiError::bad_request(
            "INVALID_ADDRESS",
            "address must be 0x plus 40 hex characters",
        )
    })?;
    let start = request.start_time.unwrap_or_else(|| {
        Utc::now()
            - TimeDelta::from_std(state.config.default_history_lookback)
                .expect("configured duration is valid")
    });
    if start > Utc::now() {
        return Err(ApiError::bad_request(
            "INVALID_START_TIME",
            "start_time cannot be in the future",
        ));
    }
    let result = state
        .repo
        .create_monitor(&address, start, state.config.max_monitored_addresses)
        .await;
    let (monitor, created) = match result {
        Ok(value) => value,
        Err(error) if error.to_string().contains("MONITOR_LIMIT_REACHED") => {
            return Err(ApiError::conflict(
                "MONITOR_LIMIT_REACHED",
                "V0.000 supports at most 10 monitored addresses",
            ))
        }
        Err(error) => return Err(ApiError::internal(error)),
    };
    state.monitors.start(address).await;
    Ok((
        if created {
            StatusCode::CREATED
        } else {
            StatusCode::OK
        },
        Json(monitor),
    )
        .into_response())
}

async fn list_monitors(State(state): State<Arc<AppState>>) -> Result<Json<Value>, ApiError> {
    let items = state
        .repo
        .list_monitors()
        .await
        .map_err(ApiError::internal)?;
    Ok(Json(json!({"items":items})))
}

async fn get_monitor(
    State(state): State<Arc<AppState>>,
    Path(address): Path<String>,
) -> Result<Json<Value>, ApiError> {
    let address = path_address(&address)?;
    let item = state
        .repo
        .get_monitor(&address)
        .await
        .map_err(ApiError::internal)?
        .ok_or_else(ApiError::not_found)?;
    Ok(Json(json!(item)))
}

#[derive(Deserialize)]
struct PageQuery {
    cursor: Option<String>,
    limit: Option<i64>,
}

async fn raw_trades(
    State(state): State<Arc<AppState>>,
    Path(address): Path<String>,
    Query(page): Query<PageQuery>,
) -> Result<Json<Value>, ApiError> {
    let address = path_address(&address)?;
    ensure_monitor(&state.repo, &address).await?;
    let limit = page.limit.unwrap_or(100).clamp(1, 500);
    let rows = state
        .repo
        .raw_trades(&address, decode_cursor(page.cursor)?, limit + 1)
        .await
        .map_err(ApiError::internal)?;
    paged_raw(rows, limit)
}

async fn trades(
    State(state): State<Arc<AppState>>,
    Path(address): Path<String>,
    Query(page): Query<PageQuery>,
) -> Result<Json<Value>, ApiError> {
    let address = path_address(&address)?;
    ensure_monitor(&state.repo, &address).await?;
    let limit = page.limit.unwrap_or(100).clamp(1, 500);
    let rows = state
        .repo
        .trade_facts(&address, decode_cursor(page.cursor)?, limit + 1)
        .await
        .map_err(ApiError::internal)?;
    paged_facts(rows, limit)
}

async fn deliveries(
    State(state): State<Arc<AppState>>,
    Query(page): Query<PageQuery>,
) -> Result<Json<Value>, ApiError> {
    let items = state
        .repo
        .deliveries(page.limit.unwrap_or(100).clamp(1, 500))
        .await
        .map_err(ApiError::internal)?;
    Ok(Json(json!({"items":items})))
}

async fn ensure_monitor(repo: &Repository, address: &str) -> Result<(), ApiError> {
    if repo
        .get_monitor(address)
        .await
        .map_err(ApiError::internal)?
        .is_none()
    {
        return Err(ApiError::not_found());
    }
    Ok(())
}

fn paged_raw(mut rows: Vec<RawTradeView>, limit: i64) -> Result<Json<Value>, ApiError> {
    let has_more = rows.len() > limit as usize;
    rows.truncate(limit as usize);
    let next = if has_more {
        rows.last()
            .map(|row| encode_cursor(row.event_time, &row.source_event_id))
    } else {
        None
    };
    Ok(Json(json!({"items":rows,"next_cursor":next})))
}

fn paged_facts(mut rows: Vec<TradeFactView>, limit: i64) -> Result<Json<Value>, ApiError> {
    let has_more = rows.len() > limit as usize;
    rows.truncate(limit as usize);
    let next = if has_more {
        rows.last()
            .map(|row| encode_cursor(row.occurred_at, &row.ordering_key))
    } else {
        None
    };
    Ok(Json(json!({"items":rows,"next_cursor":next})))
}

#[derive(Serialize, Deserialize)]
struct Cursor {
    time: DateTime<Utc>,
    id: String,
}

fn encode_cursor(time: DateTime<Utc>, id: &str) -> String {
    URL_SAFE_NO_PAD.encode(
        serde_json::to_vec(&Cursor {
            time,
            id: id.into(),
        })
        .expect("cursor serializes"),
    )
}

fn decode_cursor(value: Option<String>) -> Result<Option<(DateTime<Utc>, String)>, ApiError> {
    value
        .map(|value| {
            let bytes = URL_SAFE_NO_PAD.decode(value).map_err(|_| {
                ApiError::bad_request("INVALID_CURSOR", "cursor is not valid base64url")
            })?;
            let cursor: Cursor = serde_json::from_slice(&bytes).map_err(|_| {
                ApiError::bad_request("INVALID_CURSOR", "cursor payload is invalid")
            })?;
            Ok((cursor.time, cursor.id))
        })
        .transpose()
}

fn path_address(value: &str) -> Result<String, ApiError> {
    validate_address(value)
        .ok_or_else(|| ApiError::bad_request("INVALID_ADDRESS", "invalid address path"))
}

#[derive(Debug)]
struct ApiError {
    status: StatusCode,
    code: &'static str,
    message: String,
}

impl ApiError {
    fn bad_request(code: &'static str, message: &str) -> Self {
        Self {
            status: StatusCode::BAD_REQUEST,
            code,
            message: message.into(),
        }
    }

    fn conflict(code: &'static str, message: &str) -> Self {
        Self {
            status: StatusCode::CONFLICT,
            code,
            message: message.into(),
        }
    }

    fn not_found() -> Self {
        Self {
            status: StatusCode::NOT_FOUND,
            code: "MONITOR_NOT_FOUND",
            message: "monitored address not found".into(),
        }
    }

    fn internal(error: impl std::fmt::Display) -> Self {
        tracing::error!(error = %error, "request failed");
        Self {
            status: StatusCode::INTERNAL_SERVER_ERROR,
            code: "INTERNAL_ERROR",
            message: "internal server error".into(),
        }
    }
}

impl IntoResponse for ApiError {
    fn into_response(self) -> Response {
        (
            self.status,
            Json(json!({"error":{"code":self.code,"message":self.message}})),
        )
            .into_response()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn cursor_round_trip() {
        let time = Utc::now();
        let encoded = encode_cursor(time, "abc");
        let decoded = decode_cursor(Some(encoded)).unwrap().unwrap();
        assert_eq!(decoded.0, time);
        assert_eq!(decoded.1, "abc");
    }
}
