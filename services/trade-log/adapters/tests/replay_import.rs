mod common;
use common::*;
use hyperliquid::parser::HyperliquidParser;
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn reparse_is_readonly_and_detects_corruption() {
    let db = database().await;
    collect(&db, "query_reparse", fills(3), 1).await.unwrap();
    let report = db
        .reparse("query_reparse", &HyperliquidParser)
        .await
        .unwrap();
    assert_eq!(report.comparison, "SAME");
    assert_eq!(report.fact_records, 3);
    assert_eq!(count(&db, "account_fact_versions").await, 3);
    sqlx::query("UPDATE trade_log.raw_logs SET body='corrupted'::bytea WHERE kind='userFills'")
        .execute(&db.pool)
        .await
        .unwrap();
    assert_eq!(
        db.reparse("query_reparse", &HyperliquidParser)
            .await
            .unwrap_err()
            .code,
        "INCOMPLETE_DATA"
    );
    db.pool.close().await;
}
#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn legacy_fixture_import_is_idempotent_and_bad_checksum_rejected() {
    let db = database().await;
    let fixture = std::path::Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("../../../tests/fixtures/v0.2/import-evidence/query_fixture");
    assert!(
        db.import_directory(&fixture, &HyperliquidParser)
            .await
            .unwrap()
    );
    assert!(
        !db.import_directory(&fixture, &HyperliquidParser)
            .await
            .unwrap()
    );
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert_eq!(
        db.reparse("query_fixture", &HyperliquidParser)
            .await
            .unwrap()
            .comparison,
        "SAME"
    );
    let temp = std::env::temp_dir().join(format!("robotech-import-{}", uuid::Uuid::new_v4()));
    let copied = temp.join("query_fixture");
    fn copy(from: &std::path::Path, to: &std::path::Path) {
        std::fs::create_dir_all(to).unwrap();
        for entry in std::fs::read_dir(from).unwrap() {
            let entry = entry.unwrap();
            let target = to.join(entry.file_name());
            if entry.file_type().unwrap().is_dir() {
                copy(&entry.path(), &target)
            } else {
                std::fs::copy(entry.path(), target).unwrap();
            }
        }
    }
    copy(&fixture, &copied);
    let mut manifest: serde_json::Value =
        serde_json::from_slice(&std::fs::read(copied.join("manifest.json")).unwrap()).unwrap();
    manifest
        .as_object_mut()
        .unwrap()
        .remove("normalized_account");
    std::fs::write(
        copied.join("manifest.json"),
        serde_json::to_vec(&manifest).unwrap(),
    )
    .unwrap();
    assert!(
        !db.import_directory(&copied, &HyperliquidParser)
            .await
            .unwrap()
    );
    std::fs::write(
        copied.join("responses/raw_query_fixture_userFills_1.body"),
        b"[]",
    )
    .unwrap();
    assert_eq!(
        db.import_directory(&copied, &HyperliquidParser)
            .await
            .unwrap_err()
            .code,
        "INCOMPLETE_DATA"
    );
    assert_eq!(count(&db, "account_facts_current").await, 1);
    std::fs::remove_dir_all(temp).unwrap();
    db.pool.close().await;
}

#[tokio::test]
#[ignore = "requires isolated PostgreSQL via ROBOTECH_TEST_DATABASE_URL"]
async fn interrupted_partial_import_resumes_without_duplicate_raw_logs() {
    use protocol_api::QueryKind;
    use trade_log::{
        acquisition::SourceResponse, raw_log::RawEvidenceStore, validation::QueryRequest,
    };
    let db = database().await;
    let fixture = std::path::Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("../../../tests/fixtures/v0.2/import-evidence/query_fixture");
    let request = QueryRequest {
        account: ACCOUNT.into(),
        limit: 100,
    };
    db.begin(
        "query_fixture",
        &request,
        "trace_0123456789abcdef0123456789abcdef",
    )
    .await
    .unwrap();
    db.save_response(
        "query_fixture",
        QueryKind::UserFills,
        1,
        &SourceResponse {
            status: 200,
            body: std::fs::read(fixture.join("responses/raw_query_fixture_userFills_1.body"))
                .unwrap(),
            received_at: "2026-10-04T12:00:00.000Z".into(),
            retry_after_seconds: None,
        },
        &request,
    )
    .await
    .unwrap();
    db.interrupt().await.unwrap();
    assert!(
        db.import_directory(&fixture, &HyperliquidParser)
            .await
            .unwrap()
    );
    assert_eq!(count(&db, "raw_logs").await, 3);
    assert_eq!(count(&db, "account_facts_current").await, 1);
    assert!(
        !db.import_directory(&fixture, &HyperliquidParser)
            .await
            .unwrap()
    );
    db.pool.close().await;
}
