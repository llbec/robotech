use super::*;
use account_facts::{AccountFact, AccountFactEnvelope};
use chrono::{DateTime, Utc};
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use sqlx::{Postgres as Pg, Row, Transaction};
use trade_log::publishing::{CandidateContext, CandidatePolicy, decide};
use uuid::Uuid;

pub struct PublishingObservation {
    control: Option<(String, i64, CandidatePolicy)>,
    context: CandidateContext,
    query_id: String,
    session_id: Uuid,
}
impl Postgres {
    pub async fn configure_publishing(
        &self,
        key: &str,
        url: &str,
        enabled: bool,
        policy: &CandidatePolicy,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        // Serialize configuration even when the account has no existing control row.
        sqlx::query("SELECT pg_advisory_xact_lock(hashtextextended($1,0))")
            .bind(key)
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        let old = sqlx::query(
            "SELECT * FROM trade_log.publishing_control WHERE account_key=$1 FOR UPDATE",
        )
        .bind(key)
        .fetch_optional(&mut *tx)
        .await
        .map_err(db_error)?;
        let target = format!("webhook_{:x}", Sha256::digest(url.as_bytes()));
        if let Some(old) = old {
            let old_url: String = old.get("target_url");
            if !old_url.is_empty() && !url.is_empty() && old_url != url {
                return Err(trade_log::collection::error(
                    "TARGET_CONFIGURATION_CHANGED",
                    "Existing publication target cannot be changed",
                    false,
                ));
            }
            let activate = enabled
                && (!old.get::<bool, _>("enabled")
                    || old.get::<Value, _>("policy") != json!(policy));
            sqlx::query("UPDATE trade_log.publishing_control SET enabled=$2,target_url=CASE WHEN target_url='' THEN $3 ELSE target_url END,target_id=CASE WHEN target_url='' THEN $4 ELSE target_id END,activation_epoch=activation_epoch+CASE WHEN $5 THEN 1 ELSE 0 END,activated_at=CASE WHEN $5 THEN clock_timestamp() ELSE activated_at END,policy=$6,heartbeat_at=clock_timestamp(),last_error=NULL,updated_at=clock_timestamp() WHERE account_key=$1")
                .bind(key).bind(enabled).bind(url).bind(target).bind(activate).bind(json!(policy)).execute(&mut *tx).await.map_err(db_error)?;
        } else {
            sqlx::query("INSERT INTO trade_log.publishing_control(account_key,target_url,target_id,enabled,activation_epoch,policy,heartbeat_at) VALUES($1,$2,$3,$4,$5,$6,clock_timestamp())")
                .bind(key).bind(url).bind(target).bind(enabled).bind(if enabled {1_i64}else{0}).bind(json!(policy)).execute(&mut *tx).await.map_err(db_error)?;
        }
        tx.commit().await.map_err(db_error)
    }
    pub async fn publishing_observation(
        &self,
        tx: &mut Transaction<'_, Pg>,
        query: &str,
        c: &super::realtime::StreamCommit<'_>,
    ) -> Result<PublishingObservation, QueryError> {
        let control = sqlx::query(
            "SELECT * FROM trade_log.publishing_control WHERE account_key=$1 FOR SHARE",
        )
        .bind(&c.lease.key)
        .fetch_optional(&mut **tx)
        .await
        .map_err(db_error)?;
        let raw=sqlx::query("SELECT r.observed_at,r.payload,r.message_mode,s.subscribed_at,s.snapshot_sequence,COALESCE((k.websocket_state->>'metadata_stale')::boolean,false) AS metadata_stale,clock_timestamp() AS decided FROM trade_log.raw_logs r JOIN trade_log.collection_jobs j ON j.id=r.collection_job_id JOIN trade_log.collection_stream_sessions s ON s.session_id=r.session_id JOIN trade_log.collection_checkpoints k ON k.partition_key=s.checkpoint_key AND k.source_id='official_http' WHERE j.query_id=$1").bind(query).fetch_one(&mut **tx).await.map_err(db_error)?;
        let payload: Value = raw.get("payload");
        let context = CandidateContext {
            enabled: control.as_ref().is_some_and(|r| r.get("enabled")),
            account_key: c.lease.key.clone(),
            activated_at: control
                .as_ref()
                .map(|r| r.get("activated_at"))
                .unwrap_or_else(Utc::now),
            now: raw.get("decided"),
            received_at: raw.get("observed_at"),
            subscribed_at: raw.get("subscribed_at"),
            snapshot_sequence: raw.get("snapshot_sequence"),
            sequence: c.sequence,
            message_mode: raw.get("message_mode"),
            flag_absent: payload["data"]
                .as_object()
                .is_some_and(|o| !o.contains_key("isSnapshot")),
            metadata_stale: raw.get("metadata_stale"),
        };
        let control = control
            .map(|r| {
                Ok((
                    r.get("target_id"),
                    r.get("activation_epoch"),
                    serde_json::from_value(r.get("policy")).map_err(|_| QueryError::storage())?,
                ))
            })
            .transpose()?;
        Ok(PublishingObservation {
            control,
            context,
            query_id: query.into(),
            session_id: c.session_id,
        })
    }
    pub async fn enqueue_observation(
        &self,
        tx: &mut Transaction<'_, Pg>,
        o: &PublishingObservation,
        f: &AccountFact,
        raw: Uuid,
        index: i32,
    ) -> Result<(), QueryError> {
        let fallback = CandidatePolicy {
            max_event_age_seconds: 30,
            signal_ttl_seconds: 60,
            clock_skew_tolerance_seconds: 5,
            version: "candidate-v1".into(),
        };
        let p = o.control.as_ref().map(|v| &v.2).unwrap_or(&fallback);
        let (result, reason) = match decide(f, &o.context, p) {
            Err(reason) => ("SUPPRESSED", reason),
            Ok(reason) => {
                let canonical:Value=sqlx::query_scalar("SELECT payload FROM trade_log.account_fact_versions WHERE fact_id=$1 AND revision=1").bind(&f.fact_id).fetch_one(&mut **tx).await.map_err(db_error)?;
                let fact: AccountFact =
                    serde_json::from_value(canonical).map_err(|_| QueryError::storage())?;
                let at = DateTime::parse_from_rfc3339(&f.occurred_at)
                    .map_err(|_| QueryError::storage())?
                    .with_timezone(&Utc);
                let event = AccountFactEnvelope {
                    schema_version: 1,
                    event_type: "account.fact.v1".into(),
                    event_id: format!("{}:1", f.fact_id),
                    partition_key: f.account_key.clone(),
                    occurred_at: f.occurred_at.clone(),
                    received_at: o
                        .context
                        .received_at
                        .to_rfc3339_opts(chrono::SecondsFormat::Millis, true),
                    stored_at: o
                        .context
                        .now
                        .to_rfc3339_opts(chrono::SecondsFormat::Millis, true),
                    published_at: None,
                    expires_at: (at + chrono::Duration::seconds(p.signal_ttl_seconds))
                        .to_rfc3339_opts(chrono::SecondsFormat::Millis, true),
                    observation: json!({"transport":"WEBSOCKET","query_id":o.query_id,"session_id":o.session_id,"message_sequence":o.context.sequence,"raw_log_id":f.raw_log_id,"message_mode":o.context.message_mode,"realtime_reason":reason,"publishing_policy_version":p.version}),
                    fact,
                };
                let rows=sqlx::query("INSERT INTO trade_log.outbox_events(event_id,fact_id,revision,topic,partition_key,target_id,payload) VALUES($1,$2,1,'account.fact.v1',$3,$4,$5) ON CONFLICT DO NOTHING").bind(&event.event_id).bind(&f.fact_id).bind(&f.account_key).bind(&o.control.as_ref().ok_or_else(QueryError::storage)?.0).bind(json!(event)).execute(&mut **tx).await.map_err(db_error)?.rows_affected();
                (
                    if rows == 1 {
                        "ELIGIBLE"
                    } else {
                        "ALREADY_ENQUEUED"
                    },
                    reason,
                )
            }
        };
        sqlx::query("INSERT INTO trade_log.publication_decisions(fact_id,revision,raw_log_id,source_index,session_id,message_sequence,activation_epoch,result,reason,decided_at) VALUES($1,1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT DO NOTHING").bind(&f.fact_id).bind(raw).bind(index).bind(o.session_id).bind(o.context.sequence).bind(o.control.as_ref().map(|v|v.1).unwrap_or(0)).bind(result).bind(reason).bind(o.context.now).execute(&mut **tx).await.map_err(db_error)?;
        Ok(())
    }
    pub async fn claim_publication(
        &self,
        key: &str,
        owner: Uuid,
        lease_seconds: u64,
    ) -> Result<Option<Claim>, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        // Control lock precedes outbox lock, shared with activation/disable semantics.
        let enabled: Option<bool> = sqlx::query_scalar(
            "SELECT enabled FROM trade_log.publishing_control WHERE account_key=$1 FOR SHARE",
        )
        .bind(key)
        .fetch_optional(&mut *tx)
        .await
        .map_err(db_error)?;
        if enabled != Some(true) {
            return Ok(None);
        }
        let row=sqlx::query("SELECT * FROM trade_log.outbox_events WHERE partition_key=$1 AND ((status IN ('PENDING','RETRY_WAIT') AND next_retry_at<=clock_timestamp()) OR (status='SENDING' AND lease_expires_at<=clock_timestamp())) ORDER BY created_at,event_id LIMIT 1 FOR UPDATE SKIP LOCKED").bind(key).fetch_optional(&mut *tx).await.map_err(db_error)?;
        let Some(row) = row else { return Ok(None) };
        let event: String = row.get("event_id");
        let attempt: i32 = row.get::<i32, _>("attempts") + 1;
        let epoch: i64 = row.get::<i64, _>("lease_epoch") + 1;
        let mut payload: Value = row.get("payload");
        let published: DateTime<Utc> = sqlx::query_scalar("SELECT clock_timestamp()")
            .fetch_one(&mut *tx)
            .await
            .map_err(db_error)?;
        let body = if let Some(body) = row.get::<Option<Vec<u8>>, _>("wire_body") {
            if format!("{:x}", Sha256::digest(&body))
                != row
                    .get::<Option<String>, _>("body_sha256")
                    .unwrap_or_default()
            {
                return Err(QueryError::incomplete("Outbox checksum mismatch"));
            }
            body
        } else {
            payload["published_at"] =
                json!(published.to_rfc3339_opts(chrono::SecondsFormat::Millis, true));
            serde_json::to_vec(&payload).map_err(|_| QueryError::storage())?
        };
        sqlx::query("UPDATE trade_log.outbox_events SET status='SENDING',attempts=$2,lease_owner=$3,lease_epoch=$4,lease_expires_at=clock_timestamp()+make_interval(secs=>$5),published_at=COALESCE(published_at,$6),payload=$7,wire_body=$8,body_sha256=$9,updated_at=clock_timestamp() WHERE event_id=$1").bind(&event).bind(attempt).bind(owner).bind(epoch).bind(lease_seconds as f64).bind(published).bind(payload).bind(&body).bind(format!("{:x}",Sha256::digest(&body))).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("INSERT INTO trade_log.delivery_attempts(event_id,attempt_number,lease_epoch) VALUES($1,$2,$3)").bind(&event).bind(attempt).bind(epoch).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        Ok(Some(Claim {
            event_id: event,
            attempt,
            epoch,
            owner,
            body,
        }))
    }
    pub async fn finish_publication(
        &self,
        c: &Claim,
        status: &str,
        http: Option<u16>,
        error: Option<&str>,
        delay: u64,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        let rows=sqlx::query("UPDATE trade_log.outbox_events SET status=$4,delivered_at=CASE WHEN $4='DELIVERED' THEN clock_timestamp() ELSE delivered_at END,next_retry_at=clock_timestamp()+make_interval(secs=>$5),last_error=$6,lease_owner=NULL,lease_expires_at=NULL,updated_at=clock_timestamp() WHERE event_id=$1 AND lease_owner=$2 AND lease_epoch=$3 AND status='SENDING' AND lease_expires_at>clock_timestamp()")
            .bind(&c.event_id).bind(c.owner).bind(c.epoch).bind(status).bind(delay as f64).bind(error).execute(&mut *tx).await.map_err(db_error)?.rows_affected();
        if rows != 1 {
            return Err(trade_log::collection::error(
                "LEASE_LOST",
                "Publication lease lost",
                true,
            ));
        }
        sqlx::query("UPDATE trade_log.delivery_attempts SET finished_at=clock_timestamp(),http_status=$3,result=$4,error=$5 WHERE event_id=$1 AND attempt_number=$2 AND lease_epoch=$6").bind(&c.event_id).bind(c.attempt).bind(http.map(i32::from)).bind(status).bind(error).bind(c.epoch).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn retry_publication(&self, key: &str, event: &str) -> Result<(), QueryError> {
        let changed=sqlx::query("UPDATE trade_log.outbox_events SET status='RETRY_WAIT',next_retry_at=clock_timestamp(),last_error=NULL,updated_at=clock_timestamp() WHERE event_id=$1 AND partition_key=$2 AND status='BLOCKED'").bind(event).bind(key).execute(&self.pool).await.map_err(db_error)?;
        if changed.rows_affected() != 1 {
            return Err(QueryError::conflict());
        }
        Ok(())
    }
    pub async fn touch_publishing_heartbeat(&self, key: &str) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.publishing_control SET heartbeat_at=clock_timestamp() WHERE account_key=$1").bind(key).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
    pub async fn publishing_heartbeat(
        &self,
        key: &str,
        error: Option<&str>,
    ) -> Result<(), QueryError> {
        sqlx::query("UPDATE trade_log.publishing_control SET heartbeat_at=clock_timestamp(),last_error=$2 WHERE account_key=$1").bind(key).bind(error).execute(&self.pool).await.map_err(db_error)?;
        Ok(())
    }
    pub async fn publishing_status(&self, key: &str) -> Result<Value, QueryError> {
        let c=sqlx::query("SELECT *,clock_timestamp() AS current_time FROM trade_log.publishing_control WHERE account_key=$1").bind(key).fetch_one(&self.pool).await.map_err(db_error)?;
        let counts:Value=sqlx::query_scalar("SELECT COALESCE(jsonb_object_agg(status,n),'{}') FROM (SELECT status,count(*) n FROM trade_log.outbox_events WHERE partition_key=$1 GROUP BY status) q").bind(key).fetch_one(&self.pool).await.map_err(db_error)?;
        let q=sqlx::query("SELECT min(created_at) FILTER(WHERE status<>'DELIVERED') AS oldest,max(published_at) AS sent,max(delivered_at) AS acknowledged,count(*) FILTER(WHERE status<>'DELIVERED' AND (payload->>'expires_at')::timestamptz<clock_timestamp()) AS expired FROM trade_log.outbox_events WHERE partition_key=$1").bind(key).fetch_one(&self.pool).await.map_err(db_error)?;
        let suppressed:Value=sqlx::query_scalar("SELECT COALESCE(jsonb_object_agg(reason,n),'{}') FROM (SELECT reason,count(*) n FROM trade_log.publication_decisions d JOIN trade_log.account_fact_versions f USING(fact_id,revision) WHERE f.account_key=$1 AND d.result='SUPPRESSED' GROUP BY reason) q").bind(key).fetch_one(&self.pool).await.map_err(db_error)?;
        let last:Option<String>=sqlx::query_scalar("SELECT last_error FROM trade_log.outbox_events WHERE partition_key=$1 AND last_error IS NOT NULL ORDER BY updated_at DESC LIMIT 1").bind(key).fetch_optional(&self.pool).await.map_err(db_error)?.flatten();
        let enabled: bool = c.get("enabled");
        let heartbeat: Option<DateTime<Utc>> = c.get("heartbeat_at");
        let now: DateTime<Utc> = c.get("current_time");
        let status = if !enabled {
            "DISABLED"
        } else if heartbeat.is_none() {
            "STARTING"
        } else if heartbeat.is_some_and(|h| (now - h).num_seconds() > 30) {
            "STOPPED"
        } else if counts["BLOCKED"].as_i64().unwrap_or(0) > 0 {
            "BLOCKED"
        } else if counts["RETRY_WAIT"].as_i64().unwrap_or(0) > 0
            || c.get::<Option<String>, _>("last_error").is_some()
        {
            "DEGRADED"
        } else {
            "RUNNING"
        };
        Ok(
            json!({"enabled":enabled,"status":status,"account":key.rsplit(':').next(),"account_key":key,"target_id":c.get::<String,_>("target_id"),"activation_epoch":c.get::<i64,_>("activation_epoch"),"activated_at":c.get::<DateTime<Utc>,_>("activated_at"),"heartbeat_at":heartbeat,"policy":c.get::<Value,_>("policy"),"outbox":counts,"oldest_pending_at":q.get::<Option<DateTime<Utc>>,_>("oldest"),"last_published_at":q.get::<Option<DateTime<Utc>>,_>("sent"),"last_delivered_at":q.get::<Option<DateTime<Utc>>,_>("acknowledged"),"last_error":last.or(c.get("last_error")),"expired_pending":q.get::<i64,_>("expired"),"suppression_counts":suppressed}),
        )
    }
}
pub struct Claim {
    pub event_id: String,
    pub attempt: i32,
    pub epoch: i64,
    pub owner: Uuid,
    pub body: Vec<u8>,
}
