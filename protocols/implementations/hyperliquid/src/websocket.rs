use async_trait::async_trait;
use futures_util::{SinkExt, StreamExt};
use serde_json::{Value, json};
use tokio::net::TcpStream;
use tokio_tungstenite::{
    MaybeTlsStream, WebSocketStream,
    tungstenite::{Message, protocol::WebSocketConfig},
};
use trade_log::{
    query::QueryError,
    realtime::{RealtimeConnection, RealtimeSource, StreamEvent, StreamMessage, stream_error},
};

pub struct WebsocketSource {
    endpoint: String,
}
impl WebsocketSource {
    pub fn new(network: &shared_types::Network) -> Self {
        Self {
            endpoint: match network {
                shared_types::Network::Mainnet => "wss://api.hyperliquid.xyz/ws",
                shared_types::Network::Testnet => "wss://api.hyperliquid-testnet.xyz/ws",
            }
            .into(),
        }
    }
    /// Allows isolated local protocol tests; production uses network-derived URLs.
    pub fn with_endpoint(endpoint: String) -> Self {
        Self { endpoint }
    }
}
struct Connection {
    socket: WebSocketStream<MaybeTlsStream<TcpStream>>,
    account: String,
}

pub fn subscription(account: &str) -> Value {
    json!({"method":"subscribe","subscription":{"type":"userFills","user":account,"aggregateByTime":false}})
}
pub fn decode(body: &[u8], account: &str) -> Result<StreamEvent, QueryError> {
    let value: Value = serde_json::from_slice(body)
        .map_err(|_| stream_error("INVALID_WS_MESSAGE", "Invalid websocket JSON"))?;
    match value["channel"].as_str() {
        Some("subscriptionResponse") => {
            let expected = subscription(account);
            if value["data"] != expected {
                return Err(stream_error(
                    "SUBSCRIPTION_REJECTED",
                    "Unexpected subscription confirmation",
                ));
            }
            Ok(StreamEvent::Subscribed)
        }
        Some("pong") => Ok(StreamEvent::Pong),
        Some("error") => Err(stream_error(
            "SUBSCRIPTION_REJECTED",
            "Source websocket error",
        )),
        Some("userFills") => {
            let user = value["data"]["user"]
                .as_str()
                .and_then(|s| shared_types::account(s).ok());
            if user.as_deref() != Some(account) {
                return Err(stream_error(
                    "WS_ACCOUNT_MISMATCH",
                    "Websocket account mismatch",
                ));
            }
            if !value["data"]["fills"].is_array() {
                return Err(stream_error("INVALID_WS_MESSAGE", "Invalid fills array"));
            }
            let mode = match value["data"]["isSnapshot"].as_bool() {
                Some(true) => "SNAPSHOT",
                Some(false) => "LIVE_UPDATE",
                None => "UNKNOWN",
            };
            Ok(StreamEvent::Data(StreamMessage {
                body: body.to_vec(),
                received_at: shared_types::now(),
                mode: mode.into(),
                fills: serde_json::to_vec(&value["data"]["fills"])
                    .map_err(|_| QueryError::storage())?,
            }))
        }
        Some(_) => Ok(StreamEvent::Ignore),
        None => Err(stream_error(
            "INVALID_WS_MESSAGE",
            "Missing websocket channel",
        )),
    }
}
#[async_trait]
impl RealtimeSource for WebsocketSource {
    async fn connect(
        &self,
        account: &str,
        max_bytes: usize,
    ) -> Result<Box<dyn RealtimeConnection>, QueryError> {
        let config = WebSocketConfig::default()
            .max_message_size(Some(max_bytes))
            .max_frame_size(Some(max_bytes));
        let (mut socket, _) =
            tokio_tungstenite::connect_async_with_config(&self.endpoint, Some(config), false)
                .await
                .map_err(ws_error)?;
        socket
            .send(Message::Text(subscription(account).to_string().into()))
            .await
            .map_err(ws_error)?;
        Ok(Box::new(Connection {
            socket,
            account: account.into(),
        }))
    }
}
fn ws_error(e: tokio_tungstenite::tungstenite::Error) -> QueryError {
    match e {
        tokio_tungstenite::tungstenite::Error::Capacity(_) => trade_log::collection::error(
            "MESSAGE_TOO_LARGE",
            "Websocket message exceeds limit",
            false,
        ),
        _ => stream_error("WS_IO_ERROR", "Websocket connection failed or closed"),
    }
}
#[async_trait]
impl RealtimeConnection for Connection {
    async fn next(&mut self) -> Result<StreamEvent, QueryError> {
        match self.socket.next().await {
            Some(Ok(Message::Text(text))) => match decode(text.as_bytes(), &self.account) {
                Err(e) if e.code != "SUBSCRIPTION_REJECTED" => Ok(StreamEvent::Rejected(
                    StreamMessage {
                        body: text.as_bytes().to_vec(),
                        received_at: shared_types::now(),
                        mode: "UNKNOWN".into(),
                        fills: b"[]".to_vec(),
                    },
                    e,
                )),
                result => result,
            },
            Some(Ok(Message::Ping(_))) => {
                self.socket.flush().await.map_err(ws_error)?;
                Ok(StreamEvent::Ignore)
            }
            Some(Ok(Message::Pong(_))) => Ok(StreamEvent::Ignore), // Control pong is not the application heartbeat.
            Some(Ok(Message::Close(_))) | None => {
                Err(stream_error("DISCONNECT", "Websocket disconnected"))
            }
            Some(Ok(_)) => Err(stream_error(
                "INVALID_WS_MESSAGE",
                "Unsupported websocket data frame",
            )),
            Some(Err(e)) => Err(ws_error(e)),
        }
    }
    async fn ping(&mut self) -> Result<(), QueryError> {
        self.socket
            .send(Message::Text("{\"method\":\"ping\"}".into()))
            .await
            .map_err(ws_error)
    }
    async fn close(&mut self) {
        let _ = self.socket.close(None).await;
    }
}
