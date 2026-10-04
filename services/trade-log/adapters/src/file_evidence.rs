use async_trait::async_trait;
use protocol_api::QueryKind;
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use shared_types::Network;
use std::path::{Path, PathBuf};
use tokio::{fs, io::AsyncWriteExt};
use trade_log::{
    acquisition::SourceResponse,
    query::{QueryError, QueryResult},
    raw_log::RawEvidenceStore,
    validation::QueryRequest,
};

pub struct FileEvidence {
    pub directory: PathBuf,
    pub network: Network,
}
impl FileEvidence {
    pub async fn check(&self) -> Result<(), QueryError> {
        fs::create_dir_all(&self.directory)
            .await
            .map_err(|_| QueryError::storage())?;
        let probe = self
            .directory
            .join(format!(".probe_{}", uuid::Uuid::new_v4()));
        atomic(&probe, b"write probe").await?;
        fs::remove_file(probe)
            .await
            .map_err(|_| QueryError::storage())
    }
    fn path(&self, id: &str) -> Result<PathBuf, QueryError> {
        if !id.starts_with("query_") || !id.bytes().all(|b| b.is_ascii_alphanumeric() || b == b'_')
        {
            return Err(QueryError::storage());
        }
        Ok(self.directory.join(id))
    }
}
async fn atomic(path: &Path, bytes: &[u8]) -> Result<(), QueryError> {
    let temporary = path.with_extension(format!("tmp_{}", uuid::Uuid::new_v4().simple()));
    let mut file = fs::OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(&temporary)
        .await
        .map_err(|_| QueryError::storage())?;
    file.write_all(bytes)
        .await
        .map_err(|_| QueryError::storage())?;
    file.sync_all().await.map_err(|_| QueryError::storage())?;
    drop(file);
    fs::rename(&temporary, path)
        .await
        .map_err(|_| QueryError::storage())
}
async fn save(path: &Path, value: &Value) -> Result<(), QueryError> {
    atomic(
        path,
        &serde_json::to_vec_pretty(value).map_err(|_| QueryError::storage())?,
    )
    .await
}
#[async_trait]
impl RawEvidenceStore for FileEvidence {
    async fn begin(&self, id: &str, request: &QueryRequest, trace: &str) -> Result<(), QueryError> {
        let directory = self.path(id)?;
        fs::create_dir(&directory)
            .await
            .map_err(|_| QueryError::storage())?;
        for child in ["requests", "responses", "metadata"] {
            fs::create_dir(directory.join(child))
                .await
                .map_err(|_| QueryError::storage())?;
        }
        save(&directory.join("manifest.json"),&json!({"query_id":id,"request":request,"normalized_account":shared_types::account(&request.account).ok(),"trace_id":trace,"network":self.network,"started_at":shared_types::now(),"status":"RUNNING","tool_version":env!("CARGO_PKG_VERSION"),"parser_version":"hyperliquid-v1","schema_version":1,"identity_version":"hl-fill-v1"})).await
    }
    async fn save_response(
        &self,
        id: &str,
        kind: QueryKind,
        attempt: usize,
        response: &SourceResponse,
        request: &QueryRequest,
    ) -> Result<String, QueryError> {
        let directory = self.path(id)?;
        let raw_id = format!("raw_{id}_{}_{}", kind.api_name(), attempt);
        let mut body = json!({"type":kind.api_name()});
        if kind == QueryKind::UserFills {
            body["user"] = request.account.clone().into();
            body["aggregateByTime"] = false.into();
        }
        atomic(
            &directory.join("responses").join(format!("{raw_id}.body")),
            &response.body,
        )
        .await?;
        save(
            &directory.join("requests").join(format!("{raw_id}.json")),
            &json!({"endpoint":self.network.endpoint(),"method":"POST","body":body}),
        )
        .await?;
        save(&directory.join("metadata").join(format!("{raw_id}.json")),&json!({"raw_log_id":raw_id,"kind":kind,"attempt":attempt,"http_status":response.status,"received_at":response.received_at,"sha256":format!("{:x}",Sha256::digest(&response.body)),"retry_after_seconds":response.retry_after_seconds})).await?;
        Ok(raw_id)
    }
    async fn finish(
        &self,
        id: &str,
        result: Result<&QueryResult, &QueryError>,
    ) -> Result<(), QueryError> {
        let directory = self.path(id)?;
        let manifest = directory.join("manifest.json");
        let bytes = fs::read(&manifest)
            .await
            .map_err(|_| QueryError::storage())?;
        let mut value: Value = serde_json::from_slice(&bytes).map_err(|_| QueryError::storage())?;
        value["finished_at"] = shared_types::now().into();
        let mut entries = fs::read_dir(directory.join("metadata"))
            .await
            .map_err(|_| QueryError::storage())?;
        let mut raw_logs = Vec::new();
        while let Some(entry) = entries
            .next_entry()
            .await
            .map_err(|_| QueryError::storage())?
        {
            if entry.path().extension().is_some_and(|ext| ext == "json") {
                raw_logs.push(entry.file_name().to_string_lossy().into_owned());
            }
        }
        raw_logs.sort();
        value["raw_log_metadata"] = json!(raw_logs);
        match result {
            Ok(result) => {
                value["status"] = "COMPLETED".into();
                value["coverage"] = result.coverage.clone().into();
                value["counts"] = json!(result.counts);
                // Persist the complete normalized set before the HTTP display limit is applied.
                let mut full = result.clone();
                full.counts.returned_records = full.trades.len();
                full.display_truncated = false;
                save(&directory.join("result.json"), &json!(full)).await?;
            }
            Err(error) => {
                value["status"] = "FAILED".into();
                value["error"] = json!(error);
            }
        }
        save(&manifest, &value).await
    }
}
