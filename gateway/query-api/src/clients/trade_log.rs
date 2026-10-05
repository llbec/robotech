use crate::config::TradeLogConfig;
use serde::{Deserialize, de::DeserializeOwned};
use std::{sync::Arc, time::Duration};
use trade_log::stored_query::{StoredRequest, StoredResult};
use trade_log::{
    query::{QueryError, QueryResult},
    validation::QueryRequest,
};

#[derive(Clone)]
pub struct TradeLogClient {
    client: reqwest::Client,
    url: String,
    authorization: reqwest::header::HeaderValue,
}
impl TradeLogClient {
    pub fn from_config(config: &TradeLogConfig) -> Result<Option<Arc<Self>>, String> {
        if !config.enabled {
            return Ok(None);
        }
        let file = config
            .credential_file
            .as_ref()
            .ok_or("trade_log.credential_file required")?;
        let token =
            std::fs::read_to_string(file).map_err(|_| "cannot read trade_log credential file")?;
        let token = token.trim();
        if token.len() < 32 || !token.bytes().all(|b| b.is_ascii_alphanumeric()) {
            return Err(
                "service credential must have at least 32 ASCII alphanumeric characters".into(),
            );
        }
        let mut authorization = reqwest::header::HeaderValue::from_str(&format!("Bearer {token}"))
            .map_err(|_| "invalid service credential")?;
        authorization.set_sensitive(true);
        let client = reqwest::Client::builder()
            .timeout(Duration::from_secs(
                config
                    .request_timeout_seconds
                    .ok_or("trade_log timeout required")?,
            ))
            .connect_timeout(Duration::from_secs(5))
            .redirect(reqwest::redirect::Policy::none())
            .build()
            .map_err(|_| "internal client initialization failed")?;
        Ok(Some(Arc::new(Self {
            client,
            url: format!(
                "{}/internal/v1/trade-queries",
                config
                    .base_url
                    .as_deref()
                    .ok_or("trade_log base URL required")?
                    .trim_end_matches('/')
            ),
            authorization,
        })))
    }
    pub async fn query(
        &self,
        request: &QueryRequest,
        trace: &str,
    ) -> Result<QueryResult, QueryError> {
        self.send(&self.url, &serde_json::json!(request), trace)
            .await
    }
    pub async fn stored(
        &self,
        request: &StoredRequest,
        trace: &str,
    ) -> Result<StoredResult, QueryError> {
        let url = self
            .url
            .trim_end_matches("/internal/v1/trade-queries")
            .to_owned()
            + "/internal/v1/stored-trade-queries";
        self.send(&url, &serde_json::json!(request), trace).await
    }
    async fn send<T: DeserializeOwned>(
        &self,
        url: &str,
        request: &serde_json::Value,
        trace: &str,
    ) -> Result<T, QueryError> {
        let mut response = self
            .client
            .post(url)
            .header("authorization", self.authorization.clone())
            .header("x-trace-id", trace)
            .json(request)
            .send()
            .await
            .map_err(|_| QueryError::unavailable("Trade query service unavailable or timed out"))?;
        let status = response.status().as_u16();
        let mut bytes = Vec::new();
        while let Some(chunk) = response
            .chunk()
            .await
            .map_err(|_| QueryError::unavailable("Trade query service response failed"))?
        {
            if bytes.len().saturating_add(chunk.len()) > 32 * 1024 * 1024 {
                return Err(QueryError::unavailable("Trade query response too large"));
            }
            bytes.extend_from_slice(&chunk);
        }
        if status == 200 {
            #[derive(Deserialize)]
            struct Success<T> {
                data: T,
            }
            return serde_json::from_slice::<Success<T>>(&bytes)
                .map(|v| v.data)
                .map_err(|_| QueryError::unavailable("Invalid trade query service response"));
        }
        let code = serde_json::from_slice::<serde_json::Value>(&bytes)
            .ok()
            .and_then(|v| v["code"].as_str().map(str::to_owned));
        let error = match (status, code.as_deref()) {
            (400, Some("VALIDATION_ERROR")) => QueryError::validation("Invalid trade query"),
            (422, Some("INCOMPLETE_DATA")) => {
                QueryError::incomplete("Source records could not be reliably parsed")
            }
            (409, Some("VERSION_CONFLICT")) => QueryError::conflict(),
            (429, Some("RATE_LIMITED")) => QueryError::limited(),
            (500, Some("INTERNAL_INVARIANT_VIOLATION")) => QueryError::storage(),
            _ => QueryError::unavailable("Trade query service unavailable"),
        };
        Err(error)
    }
}
