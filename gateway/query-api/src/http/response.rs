use super::middleware::TraceId;
use axum::Json;
use serde::Serialize;

#[derive(Serialize)]
pub struct Meta {
    pub trace_id: String,
    pub schema_version: u32,
}

#[derive(Serialize)]
pub struct Envelope<T> {
    pub data: T,
    pub meta: Meta,
}

pub fn success<T: Serialize>(data: T, trace: TraceId) -> Json<Envelope<T>> {
    Json(Envelope {
        data,
        meta: Meta {
            trace_id: trace.0,
            schema_version: 1,
        },
    })
}
