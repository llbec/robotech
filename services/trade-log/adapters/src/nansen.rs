use async_trait::async_trait;
use chrono::{DateTime, SecondsFormat, Utc};
use reqwest::{Client, StatusCode};
use serde_json::{json, Value};
use std::time::Duration;
use trade_log::{import_job::AccountDataSource, raw_log::NansenSnapshot, NansenPage, RawTrade};

pub struct NansenClient {
    client: Client,
    api_key: String,
    base_url: String,
    max_attempts: u32,
}

impl NansenClient {
    pub fn new(api_key: String) -> Result<Self, String> {
        Self::with_base_url(api_key, "https://api.nansen.ai".into())
    }

    pub fn with_base_url(api_key: String, base_url: String) -> Result<Self, String> {
        if api_key.trim().is_empty() {
            return Err("NANSEN_API_KEY is empty".into());
        }
        let client = Client::builder()
            .timeout(Duration::from_secs(30))
            .build()
            .map_err(|e| e.to_string())?;
        Ok(Self {
            client,
            api_key,
            base_url,
            max_attempts: 4,
        })
    }

    fn decode_page(page: u32, response: Value) -> Result<NansenPage, String> {
        let records = response
            .pointer("/data/data")
            .or_else(|| response.get("data").filter(|v| v.is_array()))
            .or_else(|| response.get("trades"))
            .and_then(Value::as_array)
            .ok_or_else(|| "NANSEN_INVALID_RESPONSE: missing trade array".to_string())?;
        let trades = records
            .iter()
            .cloned()
            .map(serde_json::from_value::<RawTrade>)
            .collect::<Result<Vec<_>, _>>()
            .map_err(|e| format!("NANSEN_INVALID_RESPONSE: {e}"))?;
        let is_last_page = response
            .pointer("/pagination/is_last_page")
            .or_else(|| response.pointer("/data/pagination/is_last_page"))
            .or_else(|| response.get("is_last_page"))
            .and_then(Value::as_bool)
            .ok_or_else(|| "NANSEN_INVALID_RESPONSE: missing is_last_page".to_string())?;
        Ok(NansenPage {
            page,
            is_last_page,
            response,
            trades,
        })
    }

    async fn post_json(&self, path: &str, body: &Value, context: &str) -> Result<Value, String> {
        for attempt in 1..=self.max_attempts {
            let result = self
                .client
                .post(format!("{}{}", self.base_url.trim_end_matches('/'), path))
                .header("apikey", &self.api_key)
                .json(body)
                .send()
                .await;
            match result {
                Ok(response) if response.status().is_success() => {
                    return response
                        .json::<Value>()
                        .await
                        .map_err(|error| format!("NANSEN_INVALID_RESPONSE {context}: {error}"));
                }
                Ok(response)
                    if matches!(
                        response.status(),
                        StatusCode::UNAUTHORIZED
                            | StatusCode::FORBIDDEN
                            | StatusCode::PAYMENT_REQUIRED
                    ) =>
                {
                    return Err(format!(
                        "NANSEN_AUTH_ERROR {context} status={}",
                        response.status().as_u16()
                    ));
                }
                Ok(response)
                    if response.status() == StatusCode::TOO_MANY_REQUESTS
                        || response.status().is_server_error() =>
                {
                    if attempt == self.max_attempts {
                        return Err(format!(
                            "NANSEN_RETRY_EXHAUSTED {context} status={}",
                            response.status().as_u16()
                        ));
                    }
                }
                Ok(response) => {
                    return Err(format!(
                        "NANSEN_HTTP_ERROR {context} status={}",
                        response.status().as_u16()
                    ))
                }
                Err(error) if attempt == self.max_attempts => {
                    return Err(format!("NANSEN_NETWORK_ERROR {context}: {error}"))
                }
                Err(_) => {}
            }
            tokio::time::sleep(Duration::from_secs(2u64.pow(attempt - 1))).await;
        }
        unreachable!()
    }
}

#[async_trait]
impl AccountDataSource for NansenClient {
    async fn fetch_page(
        &self,
        address: &str,
        from: DateTime<Utc>,
        to: DateTime<Utc>,
        page: u32,
    ) -> Result<NansenPage, String> {
        let body = json!({"address":address,"date":{"from":from.to_rfc3339_opts(SecondsFormat::Secs,true),"to":to.to_rfc3339_opts(SecondsFormat::Secs,true)},"pagination":{"page":page,"per_page":100},"order_by":[{"field":"timestamp","direction":"ASC"}]});
        let value = self
            .post_json(
                "/api/v1/profiler/perp-trades",
                &body,
                &format!("page={page}"),
            )
            .await?;
        Self::decode_page(page, value)
    }

    async fn fetch_positions(&self, address: &str) -> Result<NansenSnapshot, String> {
        let observed_at = Utc::now();
        let response = self
            .post_json(
                "/api/v1/profiler/perp-positions",
                &json!({"address": address}),
                "endpoint=perp-positions",
            )
            .await?;
        let positions = response
            .pointer("/data/asset_positions")
            .and_then(Value::as_array)
            .ok_or_else(|| {
                "NANSEN_INVALID_RESPONSE endpoint=perp-positions missing data.asset_positions"
                    .to_string()
            })?;
        let record_count = positions.len() as u64;
        let mut missing_fields = Vec::new();
        for field in ["margin_summary", "withdrawable"] {
            if response.pointer(&format!("/data/{field}")).is_none() {
                missing_fields.push(field.into());
            }
        }
        Ok(NansenSnapshot {
            endpoint: "POST /api/v1/profiler/perp-positions".into(),
            observed_at,
            response,
            record_count,
            missing_fields,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn decodes_nested_fixture() {
        let value = json!({"data":{"data":[{"user":"0x0000000000000000000000000000000000000000","timestamp":"2026-09-01T00:00:00Z","block_number":1,"transaction_hash":"tx","oid":1,"token_symbol":"ETH","side":"Buy","action":"Open","price":"1","size":"2","value_usd":"2","fee_usd":"0.01","fee_token_symbol":"USDC"}],"pagination":{"is_last_page":true}}});
        let page = NansenClient::decode_page(1, value).unwrap();
        assert!(page.is_last_page);
        assert_eq!(page.trades.len(), 1);
    }

    #[test]
    fn decodes_official_root_shape_and_naive_utc_timestamp() {
        let value = json!({"pagination":{"page":1,"per_page":100,"is_last_page":true},"data":[{"user":"0x0000000000000000000000000000000000000000","timestamp":"2026-09-01T00:00:00.452000","block_number":1,"transaction_hash":"tx","oid":1,"token_symbol":"ETH","side":"Long","action":"Close","price":1,"size":2,"value_usd":2,"fee_usd":0.01,"fee_token_symbol":"USDC","crossed":true}]});
        let page = NansenClient::decode_page(1, value).unwrap();
        assert_eq!(
            page.trades[0].timestamp.to_rfc3339(),
            "2026-09-01T00:00:00.452+00:00"
        );
    }
}
