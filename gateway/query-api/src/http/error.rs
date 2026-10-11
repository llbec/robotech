use super::middleware::TraceId;
use axum::{
    Extension, Json,
    extract::OriginalUri,
    http::{StatusCode, header},
    response::{IntoResponse, Response},
};
use serde::Serialize;

#[derive(Serialize)]
struct ErrorBody {
    code: String,
    message: String,
    trace_id: String,
}

pub fn response(status: StatusCode, code: &str, message: &str, trace: TraceId) -> Response {
    (
        status,
        Json(ErrorBody {
            code: code.into(),
            message: message.into(),
            trace_id: trace.0,
        }),
    )
        .into_response()
}

pub async fn not_found(Extension(trace): Extension<TraceId>) -> Response {
    response(
        StatusCode::NOT_FOUND,
        "RESOURCE_NOT_FOUND",
        "Route not found",
        trace,
    )
}

pub async fn method_not_allowed(
    Extension(trace): Extension<TraceId>,
    OriginalUri(uri): OriginalUri,
) -> Response {
    let allow = if uri.path().ends_with("/resume") {
        "POST"
    } else if uri.path() == "/api/v1/watch-accounts" {
        "GET, HEAD, POST"
    } else if uri.path().starts_with("/api/v1/watch-accounts/") {
        "GET, HEAD, DELETE"
    } else {
        "GET, HEAD"
    };
    let mut result = response(
        StatusCode::METHOD_NOT_ALLOWED,
        "METHOD_NOT_ALLOWED",
        "Method not allowed",
        trace,
    );
    result
        .headers_mut()
        .insert(header::ALLOW, allow.parse().expect("static header"));
    result
}
