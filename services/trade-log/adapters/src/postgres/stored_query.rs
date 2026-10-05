use super::*;
use account_facts::AccountFact;
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use serde_json::Value;
use sqlx::Row;
use trade_log::{
    query::ObservedRange,
    stored_query::{Cursor, StoredQuery, StoredRequest, StoredResult},
};
fn date(s: Option<&str>) -> Result<Option<DateTime<Utc>>, QueryError> {
    s.map(|v| {
        DateTime::parse_from_rfc3339(v)
            .map(|d| d.with_timezone(&Utc))
            .map_err(|_| QueryError::validation("Invalid time"))
    })
    .transpose()
}
#[async_trait]
impl StoredQuery for Postgres {
    async fn stored(&self, request: &StoredRequest) -> Result<StoredResult, QueryError> {
        let request = request.validated()?;
        let cursor = request.decode_cursor(self.network.name())?;
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        sqlx::query("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ READ ONLY")
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let current: i64 =
            sqlx::query_scalar("SELECT committed_seq FROM trade_log.ingestion_state WHERE id=1")
                .fetch_one(&mut *tx)
                .await
                .map_err(db_error)?;
        let snapshot = cursor.as_ref().map_or(current, |c| c.snapshot);
        if snapshot > current {
            return Err(QueryError::validation("Invalid future cursor"));
        }
        let account_key = format!(
            "hyperliquid:{}:hyperliquid:{}",
            self.network.name(),
            request.account
        );
        let start = date(request.start_time.as_deref())?;
        let end = date(request.end_time.as_deref())?;
        let filter = "c.account_key=$1 AND c.fact_type='TRADE' AND NOT c.is_retracted AND c.ingest_seq<=$2 AND ($3::timestamptz IS NULL OR c.occurred_at >= $3) AND ($4::timestamptz IS NULL OR c.occurred_at < $4)";
        let stat=sqlx::query(&format!("SELECT count(*) AS n,min(c.occurred_at) AS first,max(c.occurred_at) AS last FROM trade_log.account_facts_current c WHERE {filter}"))
            .bind(&account_key).bind(snapshot).bind(start).bind(end).fetch_one(&mut *tx).await.map_err(db_error)?;
        let page=sqlx::query(&format!("SELECT v.payload FROM trade_log.account_facts_current c JOIN trade_log.account_fact_versions v ON (v.fact_id,v.revision)=(c.fact_id,c.current_revision) WHERE {filter} AND ($5::timestamptz IS NULL OR c.occurred_at<$5 OR (c.occurred_at=$5 AND (c.source_tid<$6::text::numeric OR (c.source_tid=$6::text::numeric AND c.fact_id>$7)))) ORDER BY c.occurred_at DESC,c.source_tid DESC,c.fact_id ASC LIMIT $8"))
            .bind(&account_key).bind(snapshot).bind(start).bind(end).bind(date(cursor.as_ref().map(|c|c.time.as_str()))?).bind(cursor.as_ref().map(|c|&c.tid)).bind(cursor.as_ref().map(|c|&c.fact_id)).bind((request.limit+1) as i64)
            .fetch_all(&mut *tx).await.map_err(db_error)?;
        let mut trades: Vec<AccountFact> = page
            .into_iter()
            .map(|r| {
                serde_json::from_value(r.get::<Value, _>("payload"))
                    .map_err(|_| QueryError::storage())
            })
            .collect::<Result<_, _>>()?;
        let has_more = trades.len() > request.limit;
        trades.truncate(request.limit);
        let next_cursor = if has_more {
            let last = trades.last().ok_or_else(QueryError::storage)?;
            Some(
                Cursor {
                    version: 1,
                    account: request.account.clone(),
                    network: self.network.name().into(),
                    range: request.range(),
                    snapshot,
                    time: last.occurred_at.clone(),
                    tid: last.source_ref.clone(),
                    fact_id: last.fact_id.clone(),
                }
                .encode()?,
            )
        } else {
            None
        };
        let format = |t: DateTime<Utc>| t.to_rfc3339_opts(chrono::SecondsFormat::Millis, true);
        let observed_range = stat
            .get::<Option<DateTime<Utc>>, _>("first")
            .zip(stat.get::<Option<DateTime<Utc>>, _>("last"))
            .map(|(a, b)| ObservedRange {
                first_at: format(a),
                last_at: format(b),
            });
        tx.commit().await.map_err(db_error)?;
        Ok(StoredResult {
            account: request.account.clone(),
            network: self.network.name().into(),
            query_scope: "STORED_TIME_RANGE".into(),
            coverage: "STORED_RECORDS_ONLY".into(),
            request_range: request.range(),
            snapshot_seq: snapshot.to_string(),
            matched_records: stat.get::<i64, _>("n") as usize,
            returned_records: trades.len(),
            observed_range,
            trades,
            has_more,
            next_cursor,
            warnings: vec!["STORED_HISTORY_NOT_VERIFIED".into()],
        })
    }
}
