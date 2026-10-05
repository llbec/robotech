use crate::{collection::error, query::QueryError};
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::collections::VecDeque;

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct WebsocketConfig {
    pub enabled: bool,
    pub connect_timeout_seconds: u64,
    pub subscribe_timeout_seconds: u64,
    pub ping_interval_seconds: u64,
    pub pong_timeout_seconds: u64,
    pub reconnect_base_seconds: u64,
    pub reconnect_max_seconds: u64,
    pub reconnect_jitter_percent: u64,
    pub reconnect_reset_after_seconds: u64,
    pub max_message_bytes: usize,
    pub max_pending_messages: usize,
    pub max_pending_bytes: usize,
    pub metadata_refresh_seconds: u64,
    pub commit_timeout_seconds: u64,
}
impl Default for WebsocketConfig {
    fn default() -> Self {
        Self {
            enabled: false,
            connect_timeout_seconds: 10,
            subscribe_timeout_seconds: 10,
            ping_interval_seconds: 20,
            pong_timeout_seconds: 10,
            reconnect_base_seconds: 5,
            reconnect_max_seconds: 60,
            reconnect_jitter_percent: 20,
            reconnect_reset_after_seconds: 60,
            max_message_bytes: 16777216,
            max_pending_messages: 256,
            max_pending_bytes: 33554432,
            metadata_refresh_seconds: 300,
            commit_timeout_seconds: 10,
        }
    }
}
impl WebsocketConfig {
    pub fn validate(&self, lease_seconds: u64) -> Result<(), String> {
        if !(1..=60).contains(&self.connect_timeout_seconds)
            || !(1..=60).contains(&self.subscribe_timeout_seconds)
            || self.ping_interval_seconds == 0
            || self.pong_timeout_seconds == 0
            || self
                .ping_interval_seconds
                .saturating_add(self.pong_timeout_seconds)
                >= 60
            || self.reconnect_base_seconds == 0
            || self.reconnect_max_seconds < self.reconnect_base_seconds
            || self.reconnect_max_seconds > 60
            || self.reconnect_jitter_percent > 20
            || !(60..=3600).contains(&self.reconnect_reset_after_seconds)
            || !(1024..=67108864).contains(&self.max_message_bytes)
            || self.max_pending_bytes < self.max_message_bytes
            || self.max_pending_bytes > 268435456
            || !(1..=4096).contains(&self.max_pending_messages)
            || !(1..=86400).contains(&self.metadata_refresh_seconds)
            || self.commit_timeout_seconds == 0
            || self.commit_timeout_seconds > 60
            || self.commit_timeout_seconds >= lease_seconds
        {
            return Err("Invalid websocket timing, queue or message limits".into());
        }
        Ok(())
    }
}
#[derive(Clone, Debug)]
pub struct StreamMessage {
    pub body: Vec<u8>,
    pub received_at: String,
    pub mode: String,
    pub fills: Vec<u8>,
}
pub enum StreamEvent {
    Subscribed,
    Pong,
    Data(StreamMessage),
    Rejected(StreamMessage, QueryError),
    Ignore,
}
#[async_trait]
pub trait RealtimeConnection: Send {
    async fn next(&mut self) -> Result<StreamEvent, QueryError>;
    async fn ping(&mut self) -> Result<(), QueryError>;
    async fn close(&mut self);
}
#[async_trait]
pub trait RealtimeSource: Send + Sync {
    async fn connect(
        &self,
        account: &str,
        max_bytes: usize,
    ) -> Result<Box<dyn RealtimeConnection>, QueryError>;
}

/// Pure scheduler, driven by a monotonic clock. Failed handshakes count as attempts.
#[derive(Default)]
pub struct ReconnectSchedule {
    failures: u32,
    attempts: VecDeque<u64>,
    next_attempt_ms: u64,
}
impl ReconnectSchedule {
    pub fn wait_ms(&mut self, now_ms: u64, c: &WebsocketConfig, random: u64) -> u64 {
        let base = c
            .reconnect_base_seconds
            .saturating_mul(1_u64 << self.failures.min(20))
            .min(c.reconnect_max_seconds)
            * 1000;
        self.failures = self.failures.saturating_add(1);
        let extra = random % (base * c.reconnect_jitter_percent / 100 + 1);
        let wait = (base + extra)
            .min(c.reconnect_max_seconds * 1000)
            .max(self.rate_wait_ms(now_ms));
        self.next_attempt_ms = now_ms.saturating_add(wait);
        wait
    }
    pub fn pending_wait_ms(&mut self, now_ms: u64) -> u64 {
        self.next_attempt_ms
            .saturating_sub(now_ms)
            .max(self.rate_wait_ms(now_ms))
    }
    pub fn rate_wait_ms(&mut self, now_ms: u64) -> u64 {
        while self
            .attempts
            .front()
            .is_some_and(|t| now_ms.saturating_sub(*t) >= 60000)
        {
            self.attempts.pop_front();
        }
        if self.attempts.len() >= 10 {
            self.attempts[0]
                .saturating_add(60000)
                .saturating_sub(now_ms)
        } else {
            0
        }
    }
    pub fn attempted(&mut self, now_ms: u64) {
        self.attempts.push_back(now_ms);
    }
    pub fn stable(&mut self, elapsed_seconds: u64, heartbeat_ok: bool, c: &WebsocketConfig) {
        if heartbeat_ok && elapsed_seconds >= c.reconnect_reset_after_seconds {
            self.failures = 0;
        }
    }
}
pub fn stream_error(code: &str, message: &str) -> QueryError {
    error(code, message, true)
}
