use super::facts::content_hash;
use super::*;
use protocol_api::QueryKind;
use serde_json::Value;
use sha2::{Digest, Sha256};
use sqlx::Row;
use std::{collections::BTreeMap, path::Path};
use trade_log::{
    acquisition::SourceResponse,
    parsing::{ParseContext, ProtocolParser},
    query::QueryResult,
    raw_log::RawEvidenceStore,
    replay::ReplayReport,
    validation::QueryRequest,
};
impl Postgres {
    pub async fn reparse(
        &self,
        id: &str,
        parser: &dyn ProtocolParser,
    ) -> Result<ReplayReport, QueryError> {
        let job = sqlx::query(
            "SELECT account,network,status,result,job_origin,request,transport FROM trade_log.collection_jobs WHERE query_id=$1",
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await
        .map_err(db_error)?
        .ok_or_else(|| QueryError::incomplete("Unknown query"))?;
        let network: Network = serde_json::from_value(Value::String(job.get("network")))
            .map_err(|_| QueryError::storage())?;
        let rows=sqlx::query("SELECT source_event_id,kind,body,sha256,parser_version FROM trade_log.raw_logs WHERE collection_job_id=(SELECT id FROM trade_log.collection_jobs WHERE query_id=$1) AND (http_status BETWEEN 200 AND 299 OR transport='WEBSOCKET') ORDER BY attempt DESC").bind(id).fetch_all(&self.pool).await.map_err(db_error)?;
        let mut bodies = BTreeMap::new();
        let mut raw_id = String::new();
        let mut old = String::new();
        for row in rows {
            let kind: String = row.get("kind");
            let body: Vec<u8> = row.get("body");
            if format!("{:x}", Sha256::digest(&body)) != row.get::<String, _>("sha256") {
                return Err(QueryError::incomplete("Stored response checksum mismatch"));
            }
            if kind == "userFills" && raw_id.is_empty() {
                raw_id = row.get("source_event_id");
                old = row.get("parser_version");
            }
            bodies.entry(kind).or_insert(body);
        }
        let body = |kind: &str| {
            bodies
                .get(kind)
                .map(Vec::as_slice)
                .ok_or_else(|| QueryError::incomplete("Missing successful source response"))
        };
        let account: String = job.get("account");
        let (reparsed, source_records) = if job.get::<String, _>("transport") == "WEBSOCKET" {
            let result = self.stream_result(id, parser).await?;
            let count = result.counts.source_records;
            old = "hyperliquid-v1".into();
            (result, count)
        } else if job.get::<String, _>("job_origin") == "COLLECTOR" {
            let request: Value = job.get("request");
            let ids: Vec<String> = serde_json::from_value(request["accepted_raw_ids"].clone())
                .map_err(|_| QueryError::incomplete("Collection did not complete"))?;
            let meta = request["meta_raw_id"]
                .as_str()
                .ok_or_else(QueryError::storage)?;
            let spot = request["spot_meta_raw_id"]
                .as_str()
                .ok_or_else(QueryError::storage)?;
            let (result, _) = self
                .collection_result(id, &account, &ids, meta, spot, parser)
                .await?;
            let count = result.counts.source_records;
            old = "hyperliquid-v1".into();
            (result, count)
        } else {
            let parsed = parser.parse(ParseContext {
                network: &network,
                account: &account,
                raw_log_id: &raw_id,
                fills: body("userFills")?,
                meta: body("meta")?,
                spot_meta: body("spotMeta")?,
            })?;
            let source_records = parsed.source_records;
            let reparsed = trade_log::normalization::result(
                id,
                &account,
                &network,
                &shared_types::now(),
                2000,
                parsed,
            )?;
            (reparsed, source_records)
        };
        let saved: Option<Value> = job.get("result");
        let mut differences = Vec::new();
        let reference: QueryResult = serde_json::from_value(
            saved.ok_or_else(|| QueryError::incomplete("Query did not complete"))?,
        )
        .map_err(|_| QueryError::storage())?;
        let before: BTreeMap<_, _> = reference
            .trades
            .iter()
            .map(|f| {
                Ok((
                    f.fact_id.clone(),
                    trade_log::collection::content::semantic_hash(f)?,
                ))
            })
            .collect::<Result<_, QueryError>>()?;
        let after: BTreeMap<_, _> = reparsed
            .trades
            .iter()
            .map(|f| {
                Ok((
                    f.fact_id.clone(),
                    trade_log::collection::content::semantic_hash(f)?,
                ))
            })
            .collect::<Result<_, QueryError>>()?;
        for key in before.keys().chain(after.keys()) {
            if before.get(key) != after.get(key) {
                differences.push(key.clone());
            }
        }
        // Check the saved job against the immutable database versions and observations, not just its result JSON.
        let facts=sqlx::query("SELECT DISTINCT v.fact_id,v.content_hash,v.payload,v.semantic_hash_version,v.semantic_content_hash FROM trade_log.fact_observations o JOIN trade_log.account_fact_versions v ON(v.fact_id,v.revision)=(o.fact_id,o.revision) JOIN trade_log.raw_logs r ON r.id=o.raw_log_id JOIN trade_log.collection_jobs j ON j.id=r.collection_job_id WHERE j.query_id=$1").bind(id).fetch_all(&self.pool).await.map_err(db_error)?;
        let mut associated = BTreeMap::new();
        for row in facts {
            let f: account_facts::AccountFact =
                serde_json::from_value(row.get("payload")).map_err(|_| QueryError::storage())?;
            let hash = content_hash(&f)?;
            let stored: String = row.get("content_hash");
            if hash != stored {
                differences.push(f.fact_id.clone());
            }
            let semantic = trade_log::collection::content::semantic_hash(&f)?;
            if row
                .get::<Option<String>, _>("semantic_content_hash")
                .is_some_and(|v| v != semantic)
                || row
                    .get::<Option<String>, _>("semantic_hash_version")
                    .is_some_and(|v| v != trade_log::collection::content::HASH_VERSION)
            {
                differences.push(f.fact_id.clone());
            }
            associated.insert(f.fact_id, semantic);
        }
        for key in before.keys().chain(associated.keys()) {
            if before.get(key) != associated.get(key) {
                differences.push(key.clone());
            }
        }
        differences.sort();
        differences.dedup();
        Ok(ReplayReport {
            query_id: id.into(),
            old_parser_version: old,
            parser_version: "hyperliquid-v1".into(),
            source_records,
            fact_records: reparsed.trades.len(),
            comparison: if differences.is_empty() {
                "SAME"
            } else {
                "DIFFERENT"
            }
            .into(),
            differences,
        })
    }
    pub async fn import_directory(
        &self,
        directory: &Path,
        parser: &dyn ProtocolParser,
    ) -> Result<bool, QueryError> {
        let manifest = read_json(&directory.join("manifest.json")).await?;
        let id = manifest["query_id"]
            .as_str()
            .ok_or_else(QueryError::storage)?;
        if !id.starts_with("query_")
            || !id.bytes().all(|b| b.is_ascii_alphanumeric() || b == b'_')
            || directory.file_name().and_then(|n| n.to_str()) != Some(id)
        {
            return Err(QueryError::validation("Invalid evidence query identity"));
        }
        if manifest["status"] != "COMPLETED"
            || manifest["schema_version"] != 1
            || manifest["parser_version"] != "hyperliquid-v1"
            || manifest["identity_version"] != "hl-fill-v1"
        {
            return Err(QueryError::incomplete(
                "Unsupported or incomplete evidence manifest",
            ));
        }
        let network: Network = serde_json::from_value(manifest["network"].clone())
            .map_err(|_| QueryError::incomplete("Invalid evidence network"))?;
        if network != self.network {
            return Err(QueryError::validation(
                "Evidence network does not match configuration",
            ));
        }
        let request: QueryRequest = serde_json::from_value(manifest["request"].clone())
            .map_err(|_| QueryError::incomplete("Invalid evidence request"))?;
        let normalized = request.validated()?;
        if !manifest["normalized_account"].is_null()
            && manifest["normalized_account"] != normalized.account
        {
            return Err(QueryError::conflict());
        }
        let finished_at = manifest["finished_at"]
            .as_str()
            .ok_or_else(QueryError::storage)?;
        let finished_at = chrono::DateTime::parse_from_rfc3339(finished_at)
            .map_err(|_| QueryError::storage())?
            .with_timezone(&chrono::Utc);
        let trace = manifest["trace_id"]
            .as_str()
            .ok_or_else(QueryError::storage)?;
        let mut responses = Vec::new();
        let mut bodies = BTreeMap::new();
        let names = manifest["raw_log_metadata"]
            .as_array()
            .ok_or_else(QueryError::storage)?;
        let mut names_seen = std::collections::BTreeSet::new();
        for name in names {
            let name = name.as_str().ok_or_else(QueryError::storage)?;
            if Path::new(name).file_name().and_then(|n| n.to_str()) != Some(name)
                || !names_seen.insert(name.to_owned())
            {
                return Err(QueryError::storage());
            }
            let metadata = read_json(&directory.join("metadata").join(name)).await?;
            let kind: QueryKind = serde_json::from_value(metadata["kind"].clone())
                .map_err(|_| QueryError::storage())?;
            if kind == QueryKind::UserFillsByTime {
                return Err(QueryError::incomplete(
                    "Only legacy single-page evidence can be imported",
                ));
            }
            let attempt = metadata["attempt"]
                .as_u64()
                .filter(|v| (1..=2).contains(v))
                .ok_or_else(QueryError::storage)? as usize;
            let raw_id = format!("raw_{id}_{}_{}", kind.api_name(), attempt);
            if metadata["raw_log_id"] != raw_id || name != format!("{raw_id}.json") {
                return Err(QueryError::conflict());
            }
            let bytes =
                read_bytes(&directory.join("responses").join(format!("{raw_id}.body"))).await?;
            if metadata["sha256"] != format!("{:x}", Sha256::digest(&bytes)) {
                return Err(QueryError::incomplete("Evidence checksum mismatch"));
            }
            let r = read_json(&directory.join("requests").join(format!("{raw_id}.json"))).await?;
            if r["method"] != "POST"
                || r["endpoint"] != network.endpoint()
                || r["body"]["type"] != kind.api_name()
            {
                return Err(QueryError::conflict());
            }
            if kind == QueryKind::UserFills
                && (r["body"]["user"] != normalized.account
                    || r["body"]["aggregateByTime"] != false)
            {
                return Err(QueryError::conflict());
            }
            let status = metadata["http_status"]
                .as_u64()
                .and_then(|v| u16::try_from(v).ok())
                .filter(|v| (100..=599).contains(v))
                .ok_or_else(QueryError::storage)?;
            let received = metadata["received_at"]
                .as_str()
                .ok_or_else(QueryError::storage)?
                .to_owned();
            chrono::DateTime::parse_from_rfc3339(&received).map_err(|_| QueryError::storage())?;
            if (200..300).contains(&status) {
                bodies.insert(kind.api_name(), (attempt, raw_id.clone(), bytes.clone()));
            }
            responses.push((
                kind,
                attempt,
                raw_id,
                SourceResponse {
                    body: bytes,
                    status,
                    received_at: received,
                    retry_after_seconds: metadata["retry_after_seconds"].as_u64(),
                },
            ));
        }
        let body = |kind: &str| {
            bodies
                .get(kind)
                .map(|v| v.2.as_slice())
                .ok_or_else(|| QueryError::incomplete("Missing successful source response"))
        };
        let raw_id = bodies
            .get("userFills")
            .ok_or_else(|| QueryError::incomplete("Missing fills"))?
            .1
            .clone();
        let parsed = parser.parse(ParseContext {
            network: &network,
            account: &normalized.account,
            raw_log_id: &raw_id,
            fills: body("userFills")?,
            meta: body("meta")?,
            spot_meta: body("spotMeta")?,
        })?;
        let queried_at = manifest["started_at"]
            .as_str()
            .ok_or_else(QueryError::storage)?;
        chrono::DateTime::parse_from_rfc3339(queried_at).map_err(|_| QueryError::storage())?;
        let result = trade_log::normalization::result(
            id,
            &normalized.account,
            &network,
            queried_at,
            request.limit,
            parsed,
        )?;
        let existing = sqlx::query(
            "SELECT id,status,request,network FROM trade_log.collection_jobs WHERE query_id=$1",
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await
        .map_err(db_error)?;
        let mut resume = false;
        if let Some(row) = existing {
            if row.get::<Value, _>("request") != serde_json::json!(request)
                || row.get::<String, _>("network") != network.name()
            {
                return Err(QueryError::conflict());
            }
            let raw: Vec<(String, String)> = sqlx::query_as(
                "SELECT source_event_id,sha256 FROM trade_log.raw_logs WHERE collection_job_id=$1",
            )
            .bind(row.get::<uuid::Uuid, _>("id"))
            .fetch_all(&self.pool)
            .await
            .map_err(db_error)?;
            let expected: BTreeMap<_, _> = responses
                .iter()
                .map(|(_, _, rid, r)| (rid.clone(), format!("{:x}", Sha256::digest(&r.body))))
                .collect();
            let stored: BTreeMap<_, _> = raw.into_iter().collect();
            if stored
                .iter()
                .any(|(rid, hash)| expected.get(rid) != Some(hash))
            {
                return Err(QueryError::conflict());
            }
            let status: String = row.get("status");
            if status == "COMPLETED" {
                if stored != expected {
                    return Err(QueryError::conflict());
                }
                let report = self.reparse(id, parser).await?;
                if report.comparison != "SAME" {
                    return Err(QueryError::conflict());
                }
                return Ok(false);
            }
            if !matches!(status.as_str(), "FAILED" | "INTERRUPTED") {
                return Err(QueryError::conflict());
            }
            let changed=sqlx::query("UPDATE trade_log.collection_jobs SET status='RUNNING',error=NULL,updated_at=now() WHERE query_id=$1 AND status IN ('FAILED','INTERRUPTED')").bind(id).execute(&self.pool).await.map_err(db_error)?;
            if changed.rows_affected() != 1 {
                return Err(QueryError::conflict());
            }
            resume = true;
        }
        let mut store = self.clone();
        store.mirror = None;
        if !resume {
            store.begin(id, &request, trace).await?;
        }
        let result_saved = async {
            for (kind, attempt, _, response) in responses {
                store
                    .save_response(id, kind, attempt, &response, &normalized)
                    .await?;
            }
            store
                .commit_facts_at(id, &result, Some(finished_at))
                .await?;
            Ok::<_, QueryError>(())
        }
        .await;
        if let Err(e) = result_saved {
            let _ = store.fail(id, &e).await;
            return Err(e);
        }
        Ok(true)
    }
}
async fn read_bytes(path: &Path) -> Result<Vec<u8>, QueryError> {
    let meta = tokio::fs::symlink_metadata(path)
        .await
        .map_err(|_| QueryError::incomplete("Missing evidence file"))?;
    if !meta.is_file() || meta.len() > 64 * 1024 * 1024 {
        return Err(QueryError::incomplete("Invalid or oversized evidence file"));
    }
    tokio::fs::read(path)
        .await
        .map_err(|_| QueryError::storage())
}
async fn read_json(path: &Path) -> Result<Value, QueryError> {
    serde_json::from_slice(&read_bytes(path).await?)
        .map_err(|_| QueryError::incomplete("Invalid evidence JSON"))
}
