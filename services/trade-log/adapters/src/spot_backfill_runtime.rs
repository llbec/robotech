use crate::postgres::Postgres;
use protocol_api::QueryKind;
use serde_json::{Value, json};
use sqlx::Row;
use std::{sync::Arc, time::Duration};
use tokio_util::sync::CancellationToken;
use trade_log::{
    acquisition::SourceReader,
    instrument::InstrumentType,
    parsing::{ParseContext, ProtocolParser},
    query::{QueryError, QueryResult},
    raw_log::RawEvidenceStore,
    validation::QueryRequest,
};
use uuid::Uuid;
pub struct SpotBackfillRuntime {
    pub store: Postgres,
    pub source: Arc<dyn SourceReader>,
    pub parser: Arc<dyn ProtocolParser>,
}
impl SpotBackfillRuntime {
    pub async fn run(self, cancel: CancellationToken) {
        loop {
            tokio::select! {_=cancel.cancelled()=>break,_=tokio::time::sleep(Duration::from_secs(2))=>{}}
            match self.store.claim_spot_backfill(Uuid::new_v4()).await {
                Ok(Some(job)) => {
                    let r = tokio::select! {_=cancel.cancelled()=>break,r=tokio::time::timeout(Duration::from_secs(100),self.batch(&job))=>r.unwrap_or_else(|_|Err(QueryError::unavailable("Backfill round timeout").with_retry()))};
                    if let Err(e) = r {
                        tracing::warn!(account_key=job["account_key"].as_str(),job_id=job["id"].as_str(),code=%e.code,"spot_backfill_failed");
                        self.store.fail_spot_batch(&job, &e).await;
                    }
                }
                Ok(None) => {}
                Err(e) => tracing::warn!(code=%e.code,"spot_backfill_store_unavailable"),
            }
        }
    }
    async fn mark(&self, id: &str, job: &Value, origin: Option<&str>) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.collection_jobs SET request=request || $2 WHERE query_id=$1").bind(id).bind(json!({"spot_backfill":{"job_id":job["id"],"owner":job["lease_owner"],"account_key":job["account_key"],"generation":job["generation"]},"origin_query_id":origin})).execute(&self.store.pool).await.map_err(|_|QueryError::storage())?;
        Ok(())
    }
    pub async fn batch(&self, job: &Value) -> Result<(), QueryError> {
        let account = job["account_key"]
            .as_str()
            .ok_or_else(QueryError::storage)?
            .rsplit(':')
            .next()
            .ok_or_else(QueryError::storage)?;
        let start = job["start_ms"].as_i64().ok_or_else(QueryError::storage)?;
        let end = job["end_ms"].as_i64().ok_or_else(QueryError::storage)?;
        let scanned = job["scanned_ms"].as_i64().ok_or_else(QueryError::storage)?;
        let request = QueryRequest {
            account: account.into(),
            limit: 2000,
            instrument_type: InstrumentType::Spot,
        };
        let id = format!("query_spot_{}", Uuid::new_v4().simple());
        if job["evidence_complete"] != true {
            let cursor = job["evidence_cursor"]
                .as_str()
                .map(Uuid::parse_str)
                .transpose()
                .map_err(|_| QueryError::storage())?;
            let archived=sqlx::query("SELECT id,query_id FROM trade_log.collection_jobs WHERE account=$1 AND network=$2 AND status='COMPLETED' AND created_at<$3::text::timestamptz AND NOT (request ? 'spot_backfill') AND ($4::uuid IS NULL OR id>$4) ORDER BY id LIMIT 1")
                .bind(account).bind(self.store.network.name()).bind(job["created_at"].as_str()).bind(cursor).fetch_optional(&self.store.pool).await.map_err(|_|QueryError::unavailable("Cannot read archived evidence").with_retry())?;
            if let Some(row) = archived {
                let origin: String = row.get("query_id");
                let mut result = self
                    .store
                    .reparsed_result(&origin, self.parser.as_ref())
                    .await?;
                result.trades.retain(|f| {
                    f.payload.instrument_type == "SPOT"
                        && chrono::DateTime::parse_from_rfc3339(&f.occurred_at).is_ok_and(|t| {
                            t.timestamp_millis() >= start && t.timestamp_millis() < end
                        })
                });
                result.query_id = id.clone();
                result.instrument_type = InstrumentType::Spot;
                self.store.begin(&id, &request, "spot_backfill").await?;
                self.mark(&id, job, Some(&origin)).await?;
                self.persist(&id, result).await?;
                return self
                    .store
                    .finish_spot_batch(job, scanned, Some(row.get("id")), false)
                    .await;
            }
            return self.store.finish_spot_batch(job, scanned, None, true).await;
        }
        let until = (scanned + job["window_ms"].as_i64().unwrap_or(3600000)).min(end);
        self.store.begin(&id, &request, "spot_backfill").await?;
        self.mark(&id, job, None).await?;
        let r = self.http_batch(&id, &request, scanned, until).await;
        match r {
            Ok(result) => {
                self.persist(&id, result).await?;
                self.store.finish_spot_batch(job, until, None, true).await
            }
            Err(e) => {
                let _ = self.store.finish(&id, Err(&e)).await;
                Err(e)
            }
        }
    }
    async fn persist(&self, id: &str, mut result: QueryResult) -> Result<(), QueryError> {
        result.query_scope = "SPOT_BACKFILL".into();
        result.coverage = "SOURCE_HISTORY_NOT_VERIFIED".into();
        result.counts.matched_records = result.trades.len();
        result.counts.returned_records = result.trades.len();
        result.display_truncated = false;
        self.store.finish(id, Ok(&result)).await
    }
    async fn http_batch(
        &self,
        id: &str,
        request: &QueryRequest,
        start: i64,
        end: i64,
    ) -> Result<QueryResult, QueryError> {
        let response = self
            .source
            .fetch_range(&request.account, start, end)
            .await?;
        let raw = self
            .store
            .save_response(id, QueryKind::UserFills, 1, &response, request)
            .await?;
        sqlx::query("UPDATE trade_log.raw_logs SET request=jsonb_set(request,'{body}',$2) WHERE source_event_id=$1").bind(&raw).bind(json!({"type":"userFillsByTime","user":request.account,"startTime":start,"endTime":end-1,"aggregateByTime":false})).execute(&self.store.pool).await.map_err(|_|QueryError::storage())?;
        if !(200..300).contains(&response.status) {
            if response.status == 429 || response.status >= 500 {
                if let Some(delay) = response.retry_after_seconds {
                    sqlx::query("UPDATE trade_log.spot_backfill_jobs SET retry_at=clock_timestamp()+$2::double precision*interval '1 second' WHERE id=(SELECT (request->'spot_backfill'->>'job_id')::uuid FROM trade_log.collection_jobs WHERE query_id=$1)").bind(id).bind(delay.min(86400) as f64).execute(&self.store.pool).await.map_err(|_|QueryError::storage())?;
                }
                return Err(QueryError::unavailable("Backfill source HTTP error").with_retry());
            }
            return Err(QueryError::incomplete(
                "Backfill source rejected the request",
            ));
        }
        let fills: Vec<Value> = serde_json::from_slice(&response.body)
            .map_err(|_| QueryError::incomplete("Invalid fills"))?;
        if fills.len() >= 2000 {
            return Err(trade_log::collection::error(
                "SOURCE_WINDOW_SATURATED",
                "Reduce spot backfill window",
                true,
            ));
        }
        let mut metadata = vec![];
        for kind in [QueryKind::Meta, QueryKind::SpotMeta] {
            let r = self.source.fetch(kind, &request.account).await?;
            self.store.save_response(id, kind, 1, &r, request).await?;
            if !(200..300).contains(&r.status) {
                return Err(QueryError::unavailable("Backfill metadata unavailable").with_retry());
            }
            metadata.push(r.body);
        }
        let parsed = self.parser.parse(ParseContext {
            network: &self.store.network,
            account: &request.account,
            raw_log_id: &raw,
            fills: &response.body,
            meta: &metadata[0],
            spot_meta: &metadata[1],
        })?;
        if parsed.invalid_records > 0 || parsed.unsupported_records > 0 {
            return Err(QueryError::incomplete(
                "Backfill has unsupported or invalid records",
            ));
        }
        if parsed.trades.iter().any(|f| {
            chrono::DateTime::parse_from_rfc3339(&f.occurred_at).map_or(true, |t| {
                t.timestamp_millis() < start || t.timestamp_millis() >= end
            })
        }) {
            return Err(QueryError::incomplete(
                "Backfill source returned out-of-range records",
            ));
        }
        let mut result = trade_log::normalization::result(
            id,
            &request.account,
            &self.store.network,
            &shared_types::now(),
            2000,
            parsed,
        )?;
        result.instrument_type = InstrumentType::Spot;
        result
            .trades
            .retain(|f| f.payload.instrument_type == "SPOT");
        Ok(result)
    }
}
