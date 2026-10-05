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
