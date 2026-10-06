use crate::config::TradeLogConfig;
use std::{sync::Arc, time::Duration};
use trade_log::publishing::PublishingStatus;
use trade_log::query::QueryError;
#[derive(Clone)]
pub struct PublisherClient {
    client: reqwest::Client,
    url: String,
    authorization: reqwest::header::HeaderValue,
}
impl PublisherClient {
    pub fn from_config(c: &TradeLogConfig) -> Result<Option<Arc<Self>>, String> {
        if !c.enabled {
            return Ok(None);
        }
        let token = std::fs::read_to_string(
            c.credential_file
                .as_ref()
                .ok_or("publisher.credential_file required")?,
        )
        .map_err(|_| "Cannot read publisher credential file")?;
        let token = token.trim();
        if token.len() < 32 || !token.bytes().all(|b| b.is_ascii_alphanumeric()) {
            return Err("Invalid publisher credential".into());
        }
        let mut authorization = reqwest::header::HeaderValue::from_str(&format!("Bearer {token}"))
            .map_err(|_| "Invalid publisher credential")?;
        authorization.set_sensitive(true);
        let client = reqwest::Client::builder()
            .timeout(Duration::from_secs(
                c.request_timeout_seconds
                    .ok_or("Publisher timeout required")?,
            ))
            .connect_timeout(Duration::from_secs(3))
            .redirect(reqwest::redirect::Policy::none())
            .build()
            .map_err(|_| "Cannot initialize publisher client")?;
        Ok(Some(Arc::new(Self {
            client,
            url: format!(
                "{}/internal/v1/publishing-status",
                c.base_url
                    .as_deref()
                    .ok_or("Publisher URL required")?
                    .trim_end_matches('/')
            ),
            authorization,
        })))
    }
    pub async fn status(&self, trace: &str) -> Result<PublishingStatus, QueryError> {
        let mut response = self
            .client
            .get(&self.url)
            .header("authorization", self.authorization.clone())
            .header("x-trace-id", trace)
            .send()
            .await
            .map_err(|_| QueryError::unavailable("Publisher unavailable or timed out"))?;
        if response.status() != 200 {
            return Err(QueryError::unavailable("Publisher status unavailable"));
        }
        let mut bytes = Vec::new();
        while let Some(chunk) = response
            .chunk()
            .await
            .map_err(|_| QueryError::unavailable("Publisher response failed"))?
        {
            if bytes.len() + chunk.len() > 65536 {
                return Err(QueryError::unavailable("Publisher response too large"));
            }
            bytes.extend_from_slice(&chunk);
        }
        #[derive(serde::Deserialize)]
        struct Success {
            data: PublishingStatus,
        }
        serde_json::from_slice::<Success>(&bytes)
            .map(|v| v.data)
            .map_err(|_| QueryError::unavailable("Invalid publisher status response"))
    }
}
