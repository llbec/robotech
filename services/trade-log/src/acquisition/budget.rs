//! Conservative Hyperliquid info-request weights; network reservations live in adapters.
use protocol_api::QueryKind;
pub fn reservation_weight(kind: QueryKind) -> i32 {
    match kind {
        QueryKind::UserFills | QueryKind::UserFillsByTime => 120,
        QueryKind::Meta | QueryKind::SpotMeta => 20,
    }
}
pub fn response_weight(kind: QueryKind, count: usize) -> i32 {
    match kind {
        QueryKind::UserFills | QueryKind::UserFillsByTime => 20 + (count / 20) as i32,
        QueryKind::Meta | QueryKind::SpotMeta => 20,
    }
}
