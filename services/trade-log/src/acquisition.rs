use crate::domain::HyperliquidFill;
use anyhow::{bail, Context, Result};
use futures_util::{SinkExt, StreamExt};
use reqwest::{Client, StatusCode};
use serde_json::{json, Value};
use std::{collections::HashSet, time::Duration};
use tokio_tungstenite::{connect_async, tungstenite::Message, MaybeTlsStream, WebSocketStream};

pub type UserSocket = WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>;

#[derive(Clone)]
pub struct HyperliquidClient {
    http: Client,
    http_url: String,
    ws_url: String,
}

pub struct FillRange {
    pub fills: Vec<HyperliquidFill>,
    pub history_complete: bool,
}

impl HyperliquidClient {
    pub fn new(http_url: String, ws_url: String) -> Result<Self> {
        Ok(Self {
            http: Client::builder().timeout(Duration::from_secs(20)).build()?,
            http_url,
            ws_url,
        })
    }

    pub async fn metadata(&self) -> Result<(Value, Value)> {
        let perp = self.post_info(json!({"type":"meta"})).await?;
        let spot = self.post_info(json!({"type":"spotMeta"})).await?;
        Ok((perp, spot))
    }

    pub async fn fills_by_time(
        &self,
        address: &str,
        start_ms: i64,
        end_ms: i64,
    ) -> Result<FillRange> {
        let mut cursor = start_ms;
        let mut fills = Vec::new();
        let mut seen = HashSet::new();
        let mut pages = 0;
        let mut capped = false;
        loop {
            pages += 1;
            let value = self
                .post_info(json!({
                    "type":"userFillsByTime",
                    "user":address,
                    "startTime":cursor,
                    "endTime":end_ms,
                    "aggregateByTime":false
                }))
                .await?;
            let page: Vec<HyperliquidFill> =
                serde_json::from_value(value).context("invalid userFillsByTime response")?;
            let page_len = page.len();
            let mut max_time = cursor;
            for fill in page {
                max_time = max_time.max(fill.time);
                if seen.insert(fill.tid) {
                    fills.push(fill);
                }
            }
            if page_len < 2000 || max_time >= end_ms {
                break;
            }
            if max_time <= cursor {
                bail!("HYPERLIQUID_PAGINATION_STALLED at {cursor}");
            }
            cursor = max_time;
            if pages >= 5 {
                capped = true;
                break;
            }
        }
        fills.sort_by_key(|fill| (fill.time, fill.tid));
        Ok(FillRange {
            history_complete: !capped,
            fills,
        })
    }

    pub async fn connect_user_fills(&self, address: &str) -> Result<UserSocket> {
        let (mut socket, _) = connect_async(&self.ws_url)
            .await
            .context("connect Hyperliquid websocket")?;
        socket
            .send(Message::Text(
                json!({"method":"subscribe","subscription":{"type":"userFills","user":address,"aggregateByTime":false}})
                    .to_string()
                    .into(),
            ))
            .await?;
        Ok(socket)
    }

    async fn post_info(&self, body: Value) -> Result<Value> {
        let mut delay = 250_u64;
        for attempt in 1..=4 {
            match self.http.post(&self.http_url).json(&body).send().await {
                Ok(response) if response.status().is_success() => {
                    return response.json().await.context("decode Hyperliquid response")
                }
                Ok(response)
                    if response.status() == StatusCode::TOO_MANY_REQUESTS
                        || response.status().is_server_error() =>
                {
                    if attempt == 4 {
                        bail!("Hyperliquid retry exhausted: HTTP {}", response.status())
                    }
                }
                Ok(response) => bail!("Hyperliquid HTTP error: {}", response.status()),
                Err(error) => {
                    if attempt == 4 {
                        return Err(error).context("Hyperliquid network error");
                    }
                }
            }
            tokio::time::sleep(Duration::from_millis(delay)).await;
            delay *= 2;
        }
        unreachable!()
    }
}

pub async fn next_text(socket: &mut UserSocket) -> Result<Option<String>> {
    while let Some(message) = socket.next().await {
        match message? {
            Message::Text(value) => return Ok(Some(value.to_string())),
            Message::Ping(value) => socket.send(Message::Pong(value)).await?,
            Message::Close(_) => return Ok(None),
            _ => {}
        }
    }
    Ok(None)
}
