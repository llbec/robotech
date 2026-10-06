#![allow(dead_code)]
use async_trait::async_trait;
use hyperliquid::parser::HyperliquidParser;
use protocol_api::QueryKind;
use serde_json::{Value, json};
use shared_types::Network;
use sqlx::{Connection, PgConnection, postgres::PgConnectOptions};
use std::{str::FromStr, sync::Arc};
use trade_log::{
    acquisition::{SourceReader, SourceResponse},
    query::{QueryError, QueryResult, QueryService},
    stored_query::StoredRequest,
    validation::QueryRequest,
};
use trade_log_adapters::postgres::Postgres;
pub const ACCOUNT: &str = "0x0000000000000000000000000000000000000001";
pub const FILLS: &str = include_str!("../../../../../tests/fixtures/hyperliquid/fills.json");
pub const META: &str = include_str!("../../../../../tests/fixtures/hyperliquid/meta.json");
pub const SPOT: &str = include_str!("../../../../../tests/fixtures/hyperliquid/spot-meta.json");
pub async fn database() -> Postgres {
    let url = std::env::var("ROBOTECH_TEST_DATABASE_URL")
        .expect("Set ROBOTECH_TEST_DATABASE_URL to an isolated PostgreSQL administrator URL");
    let options = PgConnectOptions::from_str(&url).unwrap();
    if std::env::var("ROBOTECH_TEST_DATABASE_EXACT").as_deref() == Ok("1") {
        let pool = sqlx::postgres::PgPoolOptions::new()
            .max_connections(8)
            .connect_with(options)
            .await
            .unwrap();
        let store = Postgres {
            pool,
            network: Network::Mainnet,
            mirror: None,
        };
        store.migrate().await.unwrap();
        store.check_schema().await.unwrap();
        return store;
    }
    let mut admin = PgConnection::connect_with(&options).await.unwrap();
    let name = format!("robotech_test_{}", uuid::Uuid::new_v4().simple());
    sqlx::query(&format!("CREATE DATABASE {name}"))
        .execute(&mut admin)
        .await
        .unwrap();
    let pool = sqlx::postgres::PgPoolOptions::new()
        .max_connections(8)
        .connect_with(options.database(&name))
        .await
        .unwrap();
    let store = Postgres {
        pool,
        network: Network::Mainnet,
        mirror: None,
    };
    store.migrate().await.unwrap();
    store.migrate().await.unwrap();
    store.check_schema().await.unwrap();
    store
}
pub struct Source(pub Value);
#[async_trait]
impl SourceReader for Source {
    async fn wait_retry(&self, _: u64) {}
    async fn fetch(&self, kind: QueryKind, _: &str) -> Result<SourceResponse, QueryError> {
        let body = match kind {
            QueryKind::UserFills | QueryKind::UserFillsByTime => {
                serde_json::to_vec(&self.0).unwrap()
            }
            QueryKind::Meta => META.as_bytes().to_vec(),
            QueryKind::SpotMeta => SPOT.as_bytes().to_vec(),
        };
        Ok(SourceResponse {
            status: 200,
            body,
            received_at: shared_types::now(),
            retry_after_seconds: None,
        })
    }
}
pub fn fills(n: usize) -> Value {
    let original: Value = serde_json::from_str(FILLS).unwrap();
    json!(
        (0..n)
            .map(|i| {
                let mut f = original[0].clone();
                f["tid"] = json!(9007199254740993_u64 + i as u64);
                f["time"] = json!(1791097200000_i64 + i as i64);
                f
            })
            .collect::<Vec<_>>()
    )
}
pub async fn collect(
    store: &Postgres,
    id: &str,
    source: Value,
    limit: usize,
) -> Result<QueryResult, QueryError> {
    QueryService {
        network: store.network.clone(),
        source: Arc::new(Source(source)),
        parser: Arc::new(HyperliquidParser),
        evidence: Arc::new(store.clone()),
    }
    .execute(
        id,
        &QueryRequest {
            account: ACCOUNT.into(),
            limit,
        },
        "trace_0123456789abcdef0123456789abcdef",
    )
    .await
}
pub fn request(limit: usize) -> StoredRequest {
    StoredRequest {
        account: ACCOUNT.into(),
        limit,
        start_time: None,
        end_time: None,
        cursor: None,
    }
}
pub async fn count(store: &Postgres, table: &str) -> i64 {
    sqlx::query_scalar(&format!("SELECT count(*) FROM trade_log.{table}"))
        .fetch_one(&store.pool)
        .await
        .unwrap()
}
pub fn collection_config() -> trade_log::collection::CollectionConfig {
    serde_json::from_str(include_str!(
        "../../../../../tests/fixtures/v0.3/collection-cases.json"
    ))
    .unwrap()
}
pub struct RangeSource(pub Value);
#[async_trait]
impl SourceReader for RangeSource {
    async fn wait_retry(&self, _: u64) {}
    async fn fetch(&self, kind: QueryKind, account: &str) -> Result<SourceResponse, QueryError> {
        Source(self.0.clone()).fetch(kind, account).await
    }
    async fn fetch_range(
        &self,
        _: &str,
        start: i64,
        end: i64,
    ) -> Result<SourceResponse, QueryError> {
        let records: Vec<_> = self
            .0
            .as_array()
            .unwrap()
            .iter()
            .filter(|f| f["time"].as_i64().is_some_and(|v| v >= start && v < end))
            .take(2000)
            .cloned()
            .collect();
        Ok(SourceResponse {
            status: 200,
            body: serde_json::to_vec(&records).unwrap(),
            received_at: shared_types::now(),
            retry_after_seconds: None,
        })
    }
}
pub fn runtime(
    db: &Postgres,
    c: trade_log::collection::CollectionConfig,
    source: Value,
) -> trade_log_adapters::collector_runtime::CollectorRuntime {
    trade_log_adapters::collector_runtime::CollectorRuntime {
        store: db.clone(),
        config: c,
        source: Arc::new(RangeSource(source)),
        parser: Arc::new(HyperliquidParser),
    }
}
pub async fn no_schedule(db: &Postgres) {
    sqlx::query("UPDATE trade_log.collection_checkpoints SET next_run_at=NULL")
        .execute(&db.pool)
        .await
        .unwrap();
}

// Seed the published v0.2 schema using its original columns, before running the new migrator.
// The current application deliberately refuses to operate against an outdated schema.
pub async fn legacy_collect(
    db: &Postgres,
    id: &str,
    fills: Value,
) -> Result<QueryResult, QueryError> {
    use sqlx::Row;
    use trade_log::{
        parsing::{ParseContext, ProtocolParser},
        raw_log::RawEvidenceStore,
    };
    let request = QueryRequest {
        account: ACCOUNT.into(),
        limit: 100,
    };
    db.begin(id, &request, "trace_legacy").await?;
    let source = Source(fills);
    let mut raw = String::new();
    let mut bodies = Vec::new();
    for kind in [QueryKind::UserFills, QueryKind::Meta, QueryKind::SpotMeta] {
        let response = source.fetch(kind, ACCOUNT).await?;
        let rid = db.save_response(id, kind, 1, &response, &request).await?;
        if kind == QueryKind::UserFills {
            raw = rid;
        }
        bodies.push(response.body);
    }
    let parsed = HyperliquidParser.parse(ParseContext {
        network: &db.network,
        account: ACCOUNT,
        raw_log_id: &raw,
        fills: &bodies[0],
        meta: &bodies[1],
        spot_meta: &bodies[2],
    })?;
    let result = trade_log::normalization::result(
        id,
        ACCOUNT,
        &db.network,
        &shared_types::now(),
        100,
        parsed,
    )?;
    let raw_id: uuid::Uuid =
        sqlx::query("SELECT id FROM trade_log.raw_logs WHERE source_event_id=$1")
            .bind(raw)
            .fetch_one(&db.pool)
            .await
            .unwrap()
            .get("id");
    for (index, fact) in result.trades.iter().enumerate() {
        let at = chrono::DateTime::parse_from_rfc3339(&fact.occurred_at)
            .unwrap()
            .with_timezone(&chrono::Utc);
        sqlx::query("INSERT INTO trade_log.account_fact_versions(fact_id,revision,event_id,account_key,fact_type,ordering_key,sub_index,change_type,confirmation_status,occurred_at,payload,raw_log_id,parser_version,content_hash) VALUES($1,1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,'hyperliquid-v1',$12)").bind(&fact.fact_id).bind(format!("{}:1",fact.fact_id)).bind(&fact.account_key).bind(&fact.fact_type).bind(&fact.ordering_key).bind(fact.sub_index as i32).bind(&fact.change_type).bind(&fact.confirmation_status).bind(at).bind(json!(fact)).bind(raw_id).bind(trade_log_adapters::postgres::facts::content_hash(fact)?).execute(&db.pool).await.unwrap();
        sqlx::query("INSERT INTO trade_log.account_facts_current(fact_id,current_revision,account_key,fact_type,ordering_key,sub_index,occurred_at,source_tid,ingest_seq) VALUES($1,1,$2,$3,$4,$5,$6,$7::text::numeric,$8)").bind(&fact.fact_id).bind(&fact.account_key).bind(&fact.fact_type).bind(&fact.ordering_key).bind(fact.sub_index as i32).bind(at).bind(&fact.source_ref).bind(index as i64+1).execute(&db.pool).await.unwrap();
        sqlx::query("INSERT INTO trade_log.fact_observations VALUES($1,1,$2,$3)")
            .bind(&fact.fact_id)
            .bind(raw_id)
            .bind(index as i32)
            .execute(&db.pool)
            .await
            .unwrap();
    }
    sqlx::query("UPDATE trade_log.ingestion_state SET committed_seq=$1 WHERE id=1")
        .bind(result.trades.len() as i64)
        .execute(&db.pool)
        .await
        .unwrap();
    sqlx::query(
        "UPDATE trade_log.collection_jobs SET status='COMPLETED',result=$2 WHERE query_id=$1",
    )
    .bind(id)
    .bind(json!(result))
    .execute(&db.pool)
    .await
    .unwrap();
    Ok(result)
}
