use chrono::{SecondsFormat, Utc};

#[derive(Clone)]
pub struct AppState {
    pub started_at: String,
}

impl AppState {
    pub fn new() -> Self {
        Self {
            started_at: Utc::now().to_rfc3339_opts(SecondsFormat::Millis, true),
        }
    }
}

impl Default for AppState {
    fn default() -> Self {
        Self::new()
    }
}
