use account_facts::AccountFactEnvelope;
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use uuid::Uuid;

use crate::{
    coverage::{unavailable, CoverageResult, CoverageStatus},
    normalization::{normalize_trade, raw_hash, FailedTrade},
    raw_log::{NansenPage, NansenSnapshot},
    repository::{PageWrite, TradeRepository},
    validation::{validate_address, validate_range},
};

#[derive(Debug, Clone)]
pub struct ImportRequest {
    pub address: String,
    pub from: DateTime<Utc>,
    pub to: DateTime<Utc>,
}

#[derive(Debug, Clone, Default)]
pub struct ImportSummary {
    pub run_id: Uuid,
    pub pages: u32,
    pub raw: u64,
    pub normalized: u64,
    pub failed: u64,
    pub facts: Vec<AccountFactEnvelope>,
    pub coverage: Vec<CoverageResult>,
}

#[async_trait]
pub trait AccountDataSource: Send + Sync {
    async fn fetch_page(
        &self,
        address: &str,
        from: DateTime<Utc>,
        to: DateTime<Utc>,
        page: u32,
    ) -> Result<NansenPage, String>;
    async fn fetch_positions(&self, address: &str) -> Result<NansenSnapshot, String>;
}

pub struct ImportJob<S, R> {
    source: S,
    repository: R,
}
impl<S, R> ImportJob<S, R> {
    pub fn new(source: S, repository: R) -> Self {
        Self { source, repository }
    }
}

impl<S: AccountDataSource, R: TradeRepository> ImportJob<S, R> {
    pub async fn run(&self, request: ImportRequest) -> Result<ImportSummary, String> {
        let address = validate_address(&request.address).map_err(|e| e.to_string())?;
        validate_range(request.from, request.to).map_err(|e| e.to_string())?;
        let run_id = self
            .repository
            .start_run(&address, request.from, request.to)
            .await?;
        let result = self
            .run_started(run_id, &address, request.from, request.to)
            .await;
        if let Err(error) = &result {
            let _ = self.repository.fail_run(run_id, error).await;
        }
        result
    }

    async fn run_started(
        &self,
        run_id: Uuid,
        address: &str,
        from: DateTime<Utc>,
        to: DateTime<Utc>,
    ) -> Result<ImportSummary, String> {
        let mut summary = ImportSummary {
            run_id,
            ..Default::default()
        };
        let mut trade_parse_failed = false;
        for page_number in 1.. {
            let page = self
                .source
                .fetch_page(address, from, to, page_number)
                .await?;
            let observed_at = Utc::now();
            let mut facts = Vec::new();
            let mut failures = Vec::new();
            for trade in &page.trades {
                match normalize_trade(trade, observed_at, &run_id.to_string()) {
                    Ok(fact) => facts.push(fact),
                    Err(error) => failures.push(FailedTrade {
                        raw_hash: raw_hash(trade),
                        reason: error.to_string(),
                    }),
                }
            }
            let saved = self
                .repository
                .persist_page(
                    run_id,
                    PageWrite {
                        page: page.page,
                        response: &page.response,
                        trades: &page.trades,
                        facts: &facts,
                        failures: &failures,
                    },
                )
                .await?;
            summary.pages += 1;
            summary.raw += page.trades.len() as u64;
            summary.normalized += facts.len() as u64;
            summary.failed += failures.len() as u64;
            trade_parse_failed |= !failures.is_empty();
            summary.facts.extend(saved);
            if page.is_last_page {
                break;
            }
        }
        summary.coverage.push(CoverageResult {
            category: "PERPETUAL_TRADES".into(),
            endpoint: Some("POST /api/v1/profiler/perp-trades".into()),
            status: if trade_parse_failed {
                CoverageStatus::Partial
            } else {
                CoverageStatus::Complete
            },
            record_count: summary.raw,
            range_start: Some(from),
            range_end: Some(to),
            missing_fields: Vec::new(),
            evidence: serde_json::json!({"pages": summary.pages, "conversion_failures": summary.failed}),
        });

        let positions = match self.source.fetch_positions(address).await {
            Ok(positions) => {
                self.repository.persist_snapshot(run_id, &positions).await?;
                summary.coverage.push(CoverageResult {
                    category: "CURRENT_PERPETUAL_POSITIONS".into(),
                    endpoint: Some(positions.endpoint.clone()),
                    status: if positions.missing_fields.is_empty() {
                        CoverageStatus::Complete
                    } else {
                        CoverageStatus::Partial
                    },
                    record_count: positions.record_count,
                    range_start: None,
                    range_end: Some(positions.observed_at),
                    missing_fields: positions.missing_fields.clone(),
                    evidence: serde_json::json!({"observed_at": positions.observed_at}),
                });
                Some(positions)
            }
            Err(error)
                if ["status=402", "status=403", "status=404", "status=422"]
                    .iter()
                    .any(|status| error.contains(status)) =>
            {
                let mut result = unavailable("CURRENT_PERPETUAL_POSITIONS", &error);
                result.endpoint = Some("POST /api/v1/profiler/perp-positions".into());
                summary.coverage.push(result);
                None
            }
            Err(error) => return Err(error),
        };
        let position_count = positions.as_ref().map_or(0, |item| item.record_count);
        let position_observed_at = positions.as_ref().map(|item| item.observed_at);
        summary.coverage.push(CoverageResult {
            category: "FUNDING_AND_FEES".into(),
            endpoint: Some("perp-trades + perp-positions".into()),
            status: CoverageStatus::Partial,
            record_count: summary.raw + position_count,
            range_start: Some(from),
            range_end: position_observed_at.or(Some(to)),
            missing_fields: vec!["event-level funding history".into()],
            evidence: serde_json::json!({"reason":"fees and cumulative funding are available, but complete event history is not proven"}),
        });
        summary.coverage.push(CoverageResult {
            category: "ACCOUNT_VALUE_SNAPSHOTS".into(),
            endpoint: positions.as_ref().map(|item| item.endpoint.clone()),
            status: if positions.is_some() {
                CoverageStatus::Partial
            } else {
                CoverageStatus::Unavailable
            },
            record_count: u64::from(positions.is_some()),
            range_start: None,
            range_end: position_observed_at,
            missing_fields: vec!["complete historical account-value snapshots".into()],
            evidence: serde_json::json!({"reason":"current perp state is available, but historical account value is not proven"}),
        });
        for (category, reason) in [
            (
                "CAPITAL_TRANSFERS",
                "no confirmed Nansen Hyperliquid Core address-activity endpoint",
            ),
            (
                "LIQUIDATIONS_AND_REWARDS",
                "no confirmed complete address event-history endpoint",
            ),
            (
                "MANAGEMENT_ACTIVITIES",
                "no confirmed endpoint for ApproveAgent-style activity",
            ),
        ] {
            summary.coverage.push(unavailable(category, reason));
        }
        self.repository
            .persist_coverage(run_id, &summary.coverage)
            .await?;
        summary.facts.sort_by(|a, b| {
            (&a.occurred_at, &a.fact.ordering_key, a.fact.sub_index).cmp(&(
                &b.occurred_at,
                &b.fact.ordering_key,
                b.fact.sub_index,
            ))
        });
        self.repository
            .complete_run(
                run_id,
                summary.pages,
                summary.raw,
                summary.normalized,
                summary.failed,
            )
            .await?;
        Ok(summary)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use async_trait::async_trait;
    use serde_json::json;
    use std::sync::Mutex;

    struct FixtureSource;
    #[async_trait]
    impl AccountDataSource for FixtureSource {
        async fn fetch_page(
            &self,
            address: &str,
            _from: DateTime<Utc>,
            _to: DateTime<Utc>,
            page: u32,
        ) -> Result<NansenPage, String> {
            let trade = |timestamp: &str, side: &str, oid: u32| {
                serde_json::from_value(json!({
                "user":address,"timestamp":timestamp,"block_number":page,"transaction_hash":format!("tx{oid}"),"oid":oid,
                "token_symbol":"ETH","side":side,"action":"Open","price":"10","size":"2","value_usd":"20","fee_usd":"0.1","fee_token_symbol":"USDC"
            })).unwrap()
            };
            let trades = if page == 1 {
                vec![
                    trade("2026-09-01T02:00:00Z", "Buy", 2),
                    trade("2026-09-01T01:00:00Z", "unknown", 9),
                ]
            } else {
                vec![trade("2026-09-01T00:00:00Z", "Sell", 1)]
            };
            Ok(NansenPage {
                page,
                is_last_page: page == 2,
                response: json!({"page":page}),
                trades,
            })
        }
        async fn fetch_positions(&self, _: &str) -> Result<NansenSnapshot, String> {
            Ok(NansenSnapshot {
                endpoint: "POST /api/v1/profiler/perp-positions".into(),
                observed_at: Utc::now(),
                response: json!({"data":{"asset_positions":[]}}),
                record_count: 0,
                missing_fields: Vec::new(),
            })
        }
    }

    #[derive(Default)]
    struct MemoryRepository {
        writes: Mutex<Vec<u32>>,
    }
    #[async_trait]
    impl TradeRepository for MemoryRepository {
        async fn start_run(
            &self,
            _: &str,
            _: DateTime<Utc>,
            _: DateTime<Utc>,
        ) -> Result<Uuid, String> {
            Ok(Uuid::nil())
        }
        async fn persist_page(
            &self,
            _: Uuid,
            write: PageWrite<'_>,
        ) -> Result<Vec<AccountFactEnvelope>, String> {
            self.writes.lock().unwrap().push(write.page);
            Ok(write.facts.to_vec())
        }
        async fn persist_snapshot(&self, _: Uuid, _: &NansenSnapshot) -> Result<(), String> {
            Ok(())
        }
        async fn persist_coverage(&self, _: Uuid, _: &[CoverageResult]) -> Result<(), String> {
            Ok(())
        }
        async fn complete_run(
            &self,
            _: Uuid,
            _: u32,
            _: u64,
            _: u64,
            _: u64,
        ) -> Result<(), String> {
            Ok(())
        }
        async fn fail_run(&self, _: Uuid, _: &str) -> Result<(), String> {
            Ok(())
        }
    }

    #[tokio::test]
    async fn imports_all_pages_continues_after_bad_record_and_sorts_output() {
        let from = "2026-09-01T00:00:00Z".parse().unwrap();
        let to = "2026-09-02T00:00:00Z".parse().unwrap();
        let result = ImportJob::new(FixtureSource, MemoryRepository::default())
            .run(ImportRequest {
                address: "0x0000000000000000000000000000000000000000".into(),
                from,
                to,
            })
            .await
            .unwrap();
        assert_eq!(
            (result.pages, result.raw, result.normalized, result.failed),
            (2, 3, 2, 1)
        );
        assert!(result.facts[0].occurred_at < result.facts[1].occurred_at);
        assert_eq!(result.coverage.len(), 7);
        assert_eq!(result.coverage[0].status, CoverageStatus::Partial);
    }
}
