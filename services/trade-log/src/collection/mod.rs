use crate::{checkpoint::ScanRange, parsing::Parsed, query::QueryError};
use account_facts::AccountFact;
use serde::Deserialize;
use std::collections::BTreeMap;

#[derive(Clone, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct CollectionConfig {
    pub account: String,
    pub start_time: String,
    pub interval_seconds: u64,
    pub overlap_seconds: u64,
    pub safety_delay_seconds: u64,
    pub max_window_seconds: u64,
    pub max_requests_per_round: usize,
    pub round_timeout_seconds: u64,
    pub lease_seconds: u64,
    pub retry_base_seconds: u64,
    pub retry_max_seconds: u64,
}
pub fn error(code: &str, message: &str, retryable: bool) -> QueryError {
    QueryError {
        code: code.into(),
        message: message.into(),
        retryable,
    }
}
pub fn time(ms: i64) -> String {
    chrono::DateTime::from_timestamp_millis(ms)
        .expect("validated milliseconds")
        .to_rfc3339_opts(chrono::SecondsFormat::Millis, true)
}
pub fn start_ms(text: &str) -> Result<i64, QueryError> {
    let request = crate::stored_query::StoredRequest {
        account: "0x0000000000000000000000000000000000000001".into(),
        limit: 1,
        start_time: Some(text.into()),
        end_time: None,
        cursor: None,
    }
    .validated()?;
    Ok(
        chrono::DateTime::parse_from_rfc3339(
            request.start_time.as_deref().expect("provided start"),
        )
        .map_err(|_| QueryError::validation("Invalid start time"))?
        .timestamp_millis(),
    )
}
pub fn scan_range(
    initial: i64,
    watermark: Option<i64>,
    now: i64,
    c: &CollectionConfig,
) -> Result<Option<ScanRange>, QueryError> {
    let w = watermark.unwrap_or(initial);
    if now < w {
        return Err(error(
            "CLOCK_BEHIND_WATERMARK",
            "Clock is behind the committed watermark",
            true,
        ));
    }
    let end = now
        .saturating_sub(c.safety_delay_seconds as i64 * 1000)
        .min(w.saturating_add(c.max_window_seconds as i64 * 1000));
    let start = initial.max(w.saturating_sub(c.overlap_seconds as i64 * 1000));
    // Replaying overlap is only useful when the scan can move forward.
    Ok((end > w).then_some(ScanRange {
        start_ms: start,
        end_ms: end,
    }))
}
// Keep per-page observations as well as unique facts so cross-page duplicates retain lineage.
pub fn merge(parts: Vec<Parsed>) -> Result<(Parsed, Vec<AccountFact>), QueryError> {
    let mut facts = BTreeMap::new();
    let mut observations = Vec::new();
    let mut merged = Parsed {
        source_records: 0,
        duplicate_records: 0,
        spot_records: 0,
        unsupported_records: 0,
        invalid_records: 0,
        trades: Vec::new(),
        warnings: Vec::new(),
    };
    for part in parts {
        if part.invalid_records > 0 || part.unsupported_records > 0 {
            return Err(QueryError::incomplete(
                "Collection contains invalid or unsupported records",
            ));
        }
        merged.source_records += part.source_records;
        merged.duplicate_records += part.duplicate_records;
        merged.spot_records += part.spot_records;
        merged.warnings.extend(part.warnings);
        for fact in part.trades {
            let mut value = serde_json::to_value(&fact).map_err(|_| QueryError::storage())?;
            value
                .as_object_mut()
                .ok_or_else(QueryError::storage)?
                .remove("raw_log_id");
            value["payload"]["extension"]
                .as_object_mut()
                .ok_or_else(QueryError::storage)?
                .remove("source_indices");
            if let Some((previous, _)) = facts.get(&fact.fact_id) {
                if previous != &value {
                    return Err(QueryError::conflict());
                }
                merged.duplicate_records += 1;
            } else {
                facts.insert(fact.fact_id.clone(), (value, fact.clone()));
            }
            observations.push(fact);
        }
    }
    merged.trades = facts.into_values().map(|(_, f)| f).collect();
    merged.warnings.sort();
    merged.warnings.dedup();
    Ok((merged, observations))
}
