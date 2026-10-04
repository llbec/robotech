use async_trait::async_trait;
use protocol_api::QueryKind;
use std::time::Duration;
use trade_log::{
    acquisition::{SourceReader, SourceResponse},
    query::QueryError,
};

pub struct HttpSource {
    client: reqwest::Client,
    endpoint: String,
    max_bytes: usize,
}
impl HttpSource {
    pub fn new(
        endpoint: &str,
        connect_timeout: Duration,
        request_timeout: Duration,
        max_bytes: usize,
    ) -> Result<Self, QueryError> {
        let client = reqwest::Client::builder()
            .connect_timeout(connect_timeout)
            .timeout(request_timeout)
            .redirect(reqwest::redirect::Policy::none())
            .build()
            .map_err(|_| QueryError::unavailable("Source client initialization failed"))?;
        Ok(Self {
            client,
            endpoint: endpoint.into(),
            max_bytes,
        })
    }
}
#[async_trait]
impl SourceReader for HttpSource {
    async fn fetch(&self, kind: QueryKind, account: &str) -> Result<SourceResponse, QueryError> {
        let mut body = serde_json::json!({"type":kind.api_name()});
        if kind == QueryKind::UserFills {
            body["user"] = account.into();
            body["aggregateByTime"] = false.into();
        }
        let mut response = self
            .client
            .post(&self.endpoint)
            .json(&body)
            .send()
            .await
            .map_err(|_| {
                QueryError::unavailable("Source request failed or timed out").with_retry()
            })?;
        let status = response.status().as_u16();
        let retry_after_seconds = response
            .headers()
            .get("retry-after")
            .and_then(|v| v.to_str().ok())
            .and_then(|v| v.parse().ok());
        if response
            .content_length()
            .is_some_and(|n| n > self.max_bytes as u64)
        {
            return Err(QueryError::unavailable(
                "Source response exceeds size limit",
            ));
        }
        let mut bytes = Vec::new();
        while let Some(chunk) = response.chunk().await.map_err(|_| {
            QueryError::unavailable("Source response read failed or timed out").with_retry()
        })? {
            if bytes.len().saturating_add(chunk.len()) > self.max_bytes {
                return Err(QueryError::unavailable(
                    "Source response exceeds size limit",
                ));
            }
            bytes.extend_from_slice(&chunk);
        }
        Ok(SourceResponse {
            status,
            body: bytes,
            received_at: shared_types::now(),
            retry_after_seconds,
        })
    }
    async fn wait_retry(&self, seconds: u64) {
        // A huge Retry-After must not overflow Tokio's clock. The query deadline
        // is at most 120s, so waits beyond it can be capped above that deadline.
        let jitter = (retry_jitter() % 251) as u64;
        tokio::time::sleep(Duration::from_secs(seconds.min(121)) + Duration::from_millis(jitter))
            .await;
    }
}

fn retry_jitter() -> u128 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap_or_default()
        .as_nanos()
}
