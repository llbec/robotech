use crate::{
    http::{
        middleware::TraceId,
        response::{self, Envelope},
    },
    state::AppState,
};
use axum::{Extension, Json, extract::State};
use serde::Serialize;

#[derive(Serialize)]
pub struct Health {
    status: &'static str,
    service: &'static str,
    started_at: String,
}

pub async fn handle(
    State(state): State<AppState>,
    Extension(trace): Extension<TraceId>,
) -> Json<Envelope<Health>> {
    response::success(
        Health {
            status: "ok",
            service: crate::SERVICE,
            started_at: state.started_at,
        },
        trace,
    )
}
