use crate::{config::Config, Repository};
use anyhow::Result;
use chrono::{TimeDelta, Utc};
use hmac::{Hmac, Mac};
use reqwest::Client;
use serde_json::Value;
use sha2::Sha256;
use sqlx::types::Json;
use std::time::Duration;
use tracing::{error, warn};
use uuid::Uuid;

type HmacSha256 = Hmac<Sha256>;

#[derive(sqlx::FromRow)]
struct OutboxItem {
    id: Uuid,
    event_id: String,
    payload: Json<Value>,
    attempts: i32,
}

pub async fn run(repo: Repository, config: Config) {
    let client = match Client::builder().timeout(Duration::from_secs(15)).build() {
        Ok(client) => client,
        Err(error) => {
            error!(%error, "cannot build webhook client");
            return;
        }
    };
    if let Err(error) = sqlx::query("UPDATE outbox_events SET status='RETRY_WAIT',available_at=now(),locked_at=NULL WHERE status='DELIVERING'")
        .execute(repo.pool()).await
    {
        error!(%error, "cannot recover webhook outbox");
    }
    loop {
        match claim(&repo).await {
            Ok(Some(item)) => {
                if let Err(error) = deliver(&repo, &config, &client, item).await {
                    warn!(%error, "webhook delivery persistence failed");
                }
            }
            Ok(None) => tokio::time::sleep(Duration::from_millis(500)).await,
            Err(error) => {
                warn!(%error, "cannot claim webhook outbox item");
                tokio::time::sleep(Duration::from_secs(2)).await;
            }
        }
    }
}

async fn claim(repo: &Repository) -> Result<Option<OutboxItem>> {
    let mut tx = repo.pool().begin().await?;
    let item = sqlx::query_as::<_, OutboxItem>(
        "SELECT id,event_id,payload,attempts FROM outbox_events WHERE status IN ('PENDING','RETRY_WAIT') AND available_at<=now() ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1",
    )
    .fetch_optional(&mut *tx)
    .await?;
    if let Some(item) = &item {
        sqlx::query("UPDATE outbox_events SET status='DELIVERING',attempts=attempts+1,locked_at=now() WHERE id=$1")
            .bind(item.id).execute(&mut *tx).await?;
    }
    tx.commit().await?;
    Ok(item.map(|mut item| {
        item.attempts += 1;
        item
    }))
}

async fn deliver(
    repo: &Repository,
    config: &Config,
    client: &Client,
    item: OutboxItem,
) -> Result<()> {
    let body = serde_json::to_vec(&item.payload.0)?;
    let signature = sign(&config.webhook_secret, &body)?;
    let response = client
        .post(&config.webhook_url)
        .header("Content-Type", "application/json")
        .header("Idempotency-Key", &item.event_id)
        .header("X-Robotech-Signature", format!("sha256={signature}"))
        .body(body)
        .send()
        .await;

    match response {
        Ok(response) => {
            let status = response.status();
            let summary = response
                .text()
                .await
                .unwrap_or_default()
                .chars()
                .take(1000)
                .collect::<String>();
            if status.is_success() {
                finish(
                    repo,
                    &item,
                    "DELIVERED",
                    Some(status.as_u16()),
                    &summary,
                    None,
                )
                .await
            } else if status.as_u16() == 408 || status.as_u16() == 429 || status.is_server_error() {
                retry(repo, &item, Some(status.as_u16()), &summary).await
            } else {
                finish(
                    repo,
                    &item,
                    "FAILED",
                    Some(status.as_u16()),
                    &summary,
                    Some(format!("permanent HTTP {status}")),
                )
                .await
            }
        }
        Err(error) => retry(repo, &item, None, &error.to_string()).await,
    }
}

async fn finish(
    repo: &Repository,
    item: &OutboxItem,
    status: &str,
    http_status: Option<u16>,
    summary: &str,
    error: Option<String>,
) -> Result<()> {
    let mut tx = repo.pool().begin().await?;
    sqlx::query("INSERT INTO webhook_deliveries(id,outbox_id,attempt,status,http_status,response_summary) VALUES($1,$2,$3,$4,$5,$6)")
        .bind(Uuid::new_v4()).bind(item.id).bind(item.attempts).bind(status)
        .bind(http_status.map(i32::from)).bind(summary).execute(&mut *tx).await?;
    sqlx::query("UPDATE outbox_events SET status=$2,last_error=$3,delivered_at=CASE WHEN $2='DELIVERED' THEN now() ELSE delivered_at END,locked_at=NULL WHERE id=$1")
        .bind(item.id).bind(status).bind(error).execute(&mut *tx).await?;
    tx.commit().await?;
    Ok(())
}

async fn retry(
    repo: &Repository,
    item: &OutboxItem,
    http_status: Option<u16>,
    summary: &str,
) -> Result<()> {
    let delay = 2_i64.pow(item.attempts.min(8) as u32).min(300);
    let mut tx = repo.pool().begin().await?;
    sqlx::query("INSERT INTO webhook_deliveries(id,outbox_id,attempt,status,http_status,response_summary) VALUES($1,$2,$3,'RETRY_WAIT',$4,$5)")
        .bind(Uuid::new_v4()).bind(item.id).bind(item.attempts)
        .bind(http_status.map(i32::from)).bind(summary.chars().take(1000).collect::<String>())
        .execute(&mut *tx).await?;
    sqlx::query("UPDATE outbox_events SET status='RETRY_WAIT',available_at=$2,last_error=$3,locked_at=NULL WHERE id=$1")
        .bind(item.id).bind(Utc::now() + TimeDelta::seconds(delay))
        .bind(summary.chars().take(1000).collect::<String>()).execute(&mut *tx).await?;
    tx.commit().await?;
    Ok(())
}

fn sign(secret: &str, body: &[u8]) -> Result<String> {
    let mut mac = HmacSha256::new_from_slice(secret.as_bytes())?;
    mac.update(body);
    Ok(hex::encode(mac.finalize().into_bytes()))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn webhook_signature_is_stable() {
        assert_eq!(
            sign("0123456789abcdef", br#"{"ok":true}"#).unwrap(),
            "8a782523af5169f2186640bc66718bf0be9396ae14b3d22bc0b0d8a04af83e8d"
        );
    }
}
