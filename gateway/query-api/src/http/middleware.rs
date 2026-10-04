use axum::{
    extract::{MatchedPath, Request},
    middleware::Next,
    response::Response,
};
use std::time::Instant;
use uuid::Uuid;

#[derive(Clone)]
pub struct TraceId(pub String);

pub async fn trace_request(mut request: Request, next: Next) -> Response {
    let trace = TraceId(format!("trace_{}", Uuid::new_v4().simple()));
    let method = request.method().clone();
    let path = request
        .extensions()
        .get::<MatchedPath>()
        .map_or("unmatched", |p| p.as_str())
        .to_owned();
    request.extensions_mut().insert(trace.clone());
    let started = Instant::now();
    let mut response = next.run(request).await;
    response.headers_mut().insert(
        "x-trace-id",
        trace.0.parse().expect("generated ASCII trace ID"),
    );
    let status = response.status().as_u16();
    let duration_ms = started.elapsed().as_secs_f64() * 1000.0;
    macro_rules! log_request {
        ($level:ident) => { tracing::$level!(service = crate::SERVICE, trace_id = %trace.0,
            method = %method, route = %path, status, duration_ms, "request_completed") };
    }
    if status >= 500 {
        log_request!(error);
    } else if status >= 400 {
        log_request!(warn);
    } else {
        log_request!(info);
    }
    response
}
