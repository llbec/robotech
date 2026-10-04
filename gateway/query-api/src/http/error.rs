use super::middleware::TraceId;
use axum::{
    Extension, Json,
    http::{StatusCode, header},
    response::{IntoResponse, Response},
};
use serde::Serialize;

#[derive(Serialize)]
struct ErrorBody {
    code: &'static str,
    message: &'static str,
    trace_id: String,
}

pub fn response(
    status: StatusCode,
    code: &'static str,
    message: &'static str,
    trace: TraceId,
) -> Response {
    (
        status,
        Json(ErrorBody {
            code,
            message,
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

pub async fn method_not_allowed(Extension(trace): Extension<TraceId>) -> Response {
    let mut result = response(
        StatusCode::METHOD_NOT_ALLOWED,
        "METHOD_NOT_ALLOWED",
        "Method not allowed",
        trace,
    );
    result
        .headers_mut()
        .insert(header::ALLOW, "GET, HEAD".parse().expect("static header"));
    result
}
