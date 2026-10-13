use serde::{Deserialize, Serialize};
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum InstrumentType {
    #[default]
    Perpetual,
    Spot,
    All,
}
impl InstrumentType {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Perpetual => "PERPETUAL",
            Self::Spot => "SPOT",
            Self::All => "ALL",
        }
    }
    pub fn accepts(self, value: &str) -> bool {
        self == Self::All || self.as_str() == value
    }
}
