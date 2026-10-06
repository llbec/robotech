use account_facts::AccountFact;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
#[derive(Clone, Debug, Serialize, Deserialize, PartialEq)]
#[serde(deny_unknown_fields)]
pub struct CandidatePolicy {
    pub max_event_age_seconds: i64,
    pub signal_ttl_seconds: i64,
    pub clock_skew_tolerance_seconds: i64,
    pub version: String,
}
pub struct CandidateContext {
    pub enabled: bool,
    pub account_key: String,
    pub activated_at: DateTime<Utc>,
    pub now: DateTime<Utc>,
    pub received_at: DateTime<Utc>,
    pub subscribed_at: Option<DateTime<Utc>>,
    pub snapshot_sequence: Option<i64>,
    pub sequence: i64,
    pub message_mode: String,
    pub flag_absent: bool,
    pub metadata_stale: bool,
}
/// Returns a realtime reason on admission, otherwise a durable suppression reason.
pub fn decide(
    f: &AccountFact,
    c: &CandidateContext,
    p: &CandidatePolicy,
) -> Result<&'static str, &'static str> {
    if !c.enabled {
        return Err("DISABLED");
    }
    if c.message_mode == "SNAPSHOT" {
        return Err("SNAPSHOT");
    }
    if c.subscribed_at.is_none_or(|t| t > c.received_at) {
        return Err("SESSION_UNCONFIRMED");
    }
    let reason = match c.message_mode.as_str() {
        "LIVE_UPDATE" => "EXPLICIT_LIVE_UPDATE",
        "UNKNOWN" if c.flag_absent && c.snapshot_sequence.is_some_and(|s| s < c.sequence) => {
            "POST_SNAPSHOT_UNFLAGGED"
        }
        _ => return Err("MODE_UNCONFIRMED"),
    };
    if c.metadata_stale {
        return Err("METADATA_STALE");
    }
    if f.account_key != c.account_key
        || f.fact_type != "TRADE"
        || f.revision != 1
        || f.change_type != "UPSERT"
        || f.payload.instrument_type != "PERPETUAL"
        || f.payload.trigger_type != "USER"
        || !f.payload.copy_eligible
        || !["OPEN", "INCREASE", "DECREASE", "CLOSE", "REVERSE"]
            .contains(&f.payload.position_effect.as_str())
    {
        return Err("NOT_COPY_ELIGIBLE");
    }
    let at = DateTime::parse_from_rfc3339(&f.occurred_at)
        .map_err(|_| "STALE_EVENT")?
        .with_timezone(&Utc);
    if at < c.activated_at || c.received_at < c.activated_at {
        return Err("BEFORE_ACTIVATION");
    }
    if at > c.now + chrono::Duration::seconds(p.clock_skew_tolerance_seconds) {
        return Err("FUTURE_EVENT");
    }
    if (c.now - at).num_milliseconds() > p.max_event_age_seconds * 1000
        || (c.received_at - at).num_milliseconds() > p.max_event_age_seconds * 1000
    {
        return Err("STALE_EVENT");
    }
    Ok(reason)
}
