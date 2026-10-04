use crate::http::{
    middleware::TraceId,
    response::{self, Envelope},
};
use axum::{Extension, Json};
use serde::Serialize;

#[derive(Serialize)]
pub struct Version {
    service: &'static str,
    version: &'static str,
}

pub async fn handle(Extension(trace): Extension<TraceId>) -> Json<Envelope<Version>> {
    response::success(
        Version {
            service: crate::SERVICE,
            version: crate::VERSION,
        },
        trace,
    )
}
