pub mod bootstrap;
pub mod config;
pub mod http;
pub mod lifecycle;
pub mod logging;
pub mod state;

pub const SERVICE: &str = "query-api";
pub const VERSION: &str = env!("CARGO_PKG_VERSION");
pub mod clients;
