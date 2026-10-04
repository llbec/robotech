use serde::{Deserialize, Serialize};
#[derive(Clone, Copy, Debug, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "snake_case")]
pub enum QueryKind {
    UserFills,
    Meta,
    SpotMeta,
}
impl QueryKind {
    pub fn api_name(self) -> &'static str {
        match self {
            Self::UserFills => "userFills",
            Self::Meta => "meta",
            Self::SpotMeta => "spotMeta",
        }
    }
}
