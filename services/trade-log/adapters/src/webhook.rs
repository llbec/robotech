use crate::postgres::publishing::Claim;
use chrono::{DateTime, Utc};
use std::time::Duration;
#[derive(Clone)]
pub struct Webhook {
    client: reqwest::Client,
    url: String,
    authorization: reqwest::header::HeaderValue,
    max_response_bytes: usize,
}
pub use trade_log::publishing::DeliveryOutcome as Delivery;
pub fn classify(
    http: Option<u16>,
    retry_after: Option<&str>,
    attempt: i32,
    base: u64,
    cap: u64,
    now: DateTime<Utc>,
) -> Delivery {
    let backoff = base
        .saturating_mul(
            1_u64
                .checked_shl((attempt - 1).clamp(0, 63) as u32)
                .unwrap_or(u64::MAX),
        )
        .min(cap);
    // Deterministic bounded jitter avoids a shared random-number dependency.
    let jitter = (attempt as u64 * 17) % ((backoff / 10).max(1));
    let mut delay = backoff.saturating_add(jitter).min(cap);
    if matches!(http, Some(429 | 503))
        && let Some(value) = retry_after
    {
        let value = value.trim();
        if !value.is_empty()
            && value.bytes().all(|b| b.is_ascii_digit())
            && value.parse::<u64>().is_err()
        {
            return Delivery {
                status: "BLOCKED",
                http,
                error: Some("RETRY_AFTER_EXCEEDS_LIMIT"),
                delay: 0,
            };
        }
        let wait = value.parse::<u64>().ok().or_else(|| {
            DateTime::parse_from_rfc2822(value)
                .ok()
                .map(|date| (date.with_timezone(&Utc) - now).num_seconds().max(0) as u64)
        });
        if let Some(wait) = wait {
            if wait > 86400 {
                return Delivery {
                    status: "BLOCKED",
                    http,
                    error: Some("RETRY_AFTER_EXCEEDS_LIMIT"),
                    delay: 0,
                };
            }
            delay = delay.max(wait);
        }
    }
    match http {
        Some(200..=299) => Delivery {
            status: "DELIVERED",
            http,
            error: None,
            delay: 0,
        },
        None | Some(408 | 425 | 429 | 500..=599) => Delivery {
            status: "RETRY_WAIT",
            http,
            error: Some(if http.is_none() {
                "NETWORK_OR_TIMEOUT"
            } else {
                "HTTP_TRANSIENT"
            }),
            delay,
        },
        _ => Delivery {
            status: "BLOCKED",
            http,
            error: Some("HTTP_PERMANENT_REJECTION"),
            delay: 0,
        },
    }
}
impl Webhook {
    pub fn new(
        url: &str,
        token: &str,
        timeout: u64,
        max_response_bytes: usize,
    ) -> Result<Self, String> {
        let mut authorization = reqwest::header::HeaderValue::from_str(&format!("Bearer {token}"))
            .map_err(|_| "Invalid webhook credential")?;
        authorization.set_sensitive(true);
        let client = reqwest::Client::builder()
            .timeout(Duration::from_secs(timeout))
            .connect_timeout(Duration::from_secs(timeout.min(5)))
            .redirect(reqwest::redirect::Policy::none())
            .build()
            .map_err(|_| "Cannot initialize webhook client")?;
        Ok(Self {
            client,
            url: url.into(),
            authorization,
            max_response_bytes,
        })
    }
    pub async fn send(&self, c: &Claim, base: u64, cap: u64) -> Delivery {
        self.send_bytes(&c.event_id, c.attempt, &c.body, base, cap)
            .await
    }
    async fn send_bytes(
        &self,
        event_id: &str,
        attempt: i32,
        body: &[u8],
        base: u64,
        cap: u64,
    ) -> Delivery {
        let response = self
            .client
            .post(&self.url)
            .header("authorization", self.authorization.clone())
            .header("content-type", "application/json")
            .header("x-robotech-event-id", event_id)
            .header("x-robotech-delivery-attempt", attempt)
            .body(body.to_vec())
            .send()
            .await;
        match response {
            Err(_) => classify(None, None, attempt, base, cap, Utc::now()),
            Ok(mut response) => {
                let status = response.status().as_u16();
                let after = response
                    .headers()
                    .get("retry-after")
                    .and_then(|h| h.to_str().ok())
                    .map(str::to_owned);
                // Response text is never logged. Only drain bounded bytes, within the client timeout.
                let mut size = 0;
                while let Ok(Some(chunk)) = response.chunk().await {
                    size += chunk.len();
                    if size > self.max_response_bytes {
                        break;
                    }
                }
                classify(
                    Some(status),
                    after.as_deref(),
                    attempt,
                    base,
                    cap,
                    Utc::now(),
                )
            }
        }
    }
}

#[async_trait::async_trait]
impl trade_log::publishing::PublicationTarget for Webhook {
    async fn deliver(
        &self,
        event_id: &str,
        attempt: i32,
        body: &[u8],
        base: u64,
        cap: u64,
    ) -> Delivery {
        self.send_bytes(event_id, attempt, body, base, cap).await
    }
}
