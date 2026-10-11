use crate::{
    http::{error, middleware::TraceId, response},
    state::AppState,
};
use axum::{
    Extension, Json,
    extract::{
        Path, Query, State,
        rejection::{JsonRejection, QueryRejection},
    },
    http::{HeaderMap, StatusCode},
    response::{IntoResponse, Response},
};
use subtle::ConstantTimeEq;
use trade_log::{
    query::QueryError,
    watch::{AddAccount, ListParams, invalid},
};
fn result(
    r: Result<serde_json::Value, QueryError>,
    trace: TraceId,
    status: StatusCode,
) -> Response {
    match r {
        Ok(data) => (status, response::success(data, trace)).into_response(),
        Err(e) => error::response(
            StatusCode::from_u16(e.status()).unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
            &e.code,
            &e.message,
            trace,
        ),
    }
}
fn authorize(s: &AppState, h: &HeaderMap) -> bool {
    let Some(token) = &s.management_credential else {
        return false;
    };
    let expected = format!("Bearer {token}");
    let actual = h
        .get("authorization")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("");
    actual.as_bytes().ct_eq(expected.as_bytes()).into()
}
fn unauthorized(trace: TraceId) -> Response {
    error::response(
        StatusCode::UNAUTHORIZED,
        "AUTHENTICATION_REQUIRED",
        "Management authentication required",
        trace,
    )
}
fn client(s: &AppState) -> Result<&crate::clients::collector::CollectorClient, QueryError> {
    s.collector
        .as_deref()
        .ok_or_else(|| QueryError::unavailable("Automatic collection is not enabled"))
}
pub async fn handle(
    State(s): State<AppState>,
    Extension(trace): Extension<TraceId>,
    p: Result<Query<ListParams>, QueryRejection>,
) -> Response {
    let r = async {
        let Query(p) = p.map_err(|_| invalid("Invalid list parameters"))?;
        let mut url = reqwest::Url::parse("http://localhost/").expect("literal URL");
        {
            let mut q = url.query_pairs_mut();
            if let Some(state) = p.state {
                q.append_pair("state", &state);
            }
            if let Some(limit) = p.limit {
                q.append_pair("limit", &limit.to_string());
            }
            if let Some(cursor) = p.cursor {
                q.append_pair("cursor", &cursor);
            }
        }
        let query = url
            .query()
            .filter(|v| !v.is_empty())
            .map(|v| format!("?{v}"))
            .unwrap_or_default();
        client(&s)?
            .watch(reqwest::Method::GET, "", &query, None, &trace.0)
            .await
    }
    .await;
    result(r, trace, StatusCode::OK)
}
pub async fn add(
    State(s): State<AppState>,
    Extension(trace): Extension<TraceId>,
    h: HeaderMap,
    p: Result<Json<AddAccount>, JsonRejection>,
) -> Response {
    if !authorize(&s, &h) {
        return unauthorized(trace);
    }
    let r = async {
        let Json(mut p) = p.map_err(|_| invalid("Invalid JSON request"))?;
        p.validate()?;
        client(&s)?
            .watch(
                reqwest::Method::POST,
                "",
                "",
                Some(serde_json::json!(p)),
                &trace.0,
            )
            .await
    }
    .await;
    result(r, trace, StatusCode::CREATED)
}
pub async fn detail(
    State(s): State<AppState>,
    Extension(trace): Extension<TraceId>,
    Path(id): Path<String>,
) -> Response {
    let r = async {
        uuid::Uuid::parse_str(&id).map_err(|_| invalid("Invalid account ID"))?;
        client(&s)?
            .watch(reqwest::Method::GET, &format!("/{id}"), "", None, &trace.0)
            .await
    }
    .await;
    result(r, trace, StatusCode::OK)
}
pub async fn pause(
    State(s): State<AppState>,
    Extension(trace): Extension<TraceId>,
    Path(id): Path<String>,
    h: HeaderMap,
) -> Response {
    change(s, trace, id, h, false).await
}
pub async fn resume(
    State(s): State<AppState>,
    Extension(trace): Extension<TraceId>,
    Path(id): Path<String>,
    h: HeaderMap,
) -> Response {
    change(s, trace, id, h, true).await
}
async fn change(s: AppState, trace: TraceId, id: String, h: HeaderMap, active: bool) -> Response {
    if !authorize(&s, &h) {
        return unauthorized(trace);
    }
    let r = async {
        uuid::Uuid::parse_str(&id).map_err(|_| invalid("Invalid account ID"))?;
        client(&s)?
            .watch(
                if active {
                    reqwest::Method::POST
                } else {
                    reqwest::Method::DELETE
                },
                &format!("/{id}{}", if active { "/resume" } else { "" }),
                "",
                None,
                &trace.0,
            )
            .await
    }
    .await;
    result(r, trace, StatusCode::ACCEPTED)
}
