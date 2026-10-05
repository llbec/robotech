use chrono::{SecondsFormat, Utc};

#[derive(Clone)]
pub struct AppState {
    pub collector: Option<std::sync::Arc<crate::clients::collector::CollectorClient>>,
    pub started_at: String,
    pub trade_log: Option<std::sync::Arc<crate::clients::trade_log::TradeLogClient>>,
}

impl AppState {
    pub fn new() -> Self {
        Self {
            trade_log: None,
            collector: None,
            started_at: Utc::now().to_rfc3339_opts(SecondsFormat::Millis, true),
        }
    }
}

impl Default for AppState {
    fn default() -> Self {
        Self::new()
    }
}
