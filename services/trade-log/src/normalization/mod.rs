use crate::parsing::Parsed;
use crate::query::{Counts, ObservedRange, QueryError, QueryResult};
use shared_types::Network;

pub fn result(
    id: &str,
    account: &str,
    network: &Network,
    queried_at: &str,
    limit: usize,
    mut parsed: Parsed,
) -> Result<QueryResult, QueryError> {
    if parsed.invalid_records > 0
        && parsed.trades.is_empty()
        && parsed.spot_records == 0
        && parsed.unsupported_records == 0
    {
        return Err(QueryError::incomplete("No usable source records"));
    }
    parsed.trades.sort_by(|a, b| {
        b.occurred_at
            .cmp(&a.occurred_at)
            .then_with(|| {
                b.source_ref
                    .parse::<u64>()
                    .unwrap_or(0)
                    .cmp(&a.source_ref.parse::<u64>().unwrap_or(0))
            })
            .then_with(|| a.fact_id.cmp(&b.fact_id))
    });
    if parsed.source_records >= 2000 {
        parsed.warnings.push("SOURCE_RECORD_LIMIT".into());
    }
    if parsed.source_records == 0 {
        parsed.warnings.push("NO_SOURCE_RECORDS".into());
    } else if parsed.trades.is_empty() && parsed.spot_records > 0 {
        parsed.warnings.push("NO_PERPETUAL_RECORDS".into());
    }
    let limited = parsed.source_records >= 2000
        || parsed.invalid_records > 0
        || parsed.unsupported_records > 0;
    let observed_range = parsed
        .trades
        .first()
        .zip(parsed.trades.last())
        .map(|(last, first)| ObservedRange {
            first_at: first.occurred_at.clone(),
            last_at: last.occurred_at.clone(),
        });
    let perpetual_records = parsed.trades.len();
    Ok(QueryResult {
        persistence: None,
        query_id: id.into(),
        account: account.into(),
        network: network.name().into(),
        queried_at: queried_at.into(),
        query_scope: "RECENT_SOURCE_WINDOW".into(),
        coverage: if limited {
            "LIMITED"
        } else {
            "SOURCE_WINDOW_UNVERIFIED"
        }
        .into(),
        counts: Counts {
            source_records: parsed.source_records,
            duplicate_records: parsed.duplicate_records,
            spot_records: parsed.spot_records,
            unsupported_records: parsed.unsupported_records,
            invalid_records: parsed.invalid_records,
            perpetual_records,
            returned_records: perpetual_records.min(limit),
        },
        display_truncated: perpetual_records > limit,
        observed_range,
        trades: parsed.trades,
        warnings: parsed.warnings,
        evidence_ref: id.into(),
    })
}
