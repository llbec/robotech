use super::*;
use chrono::{DateTime, Utc};
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use sqlx::{Postgres as Pg, Row, Transaction};
use trade_log::{
    checkpoint::Lease,
    collection::CollectionConfig,
    parsing::{ParseContext, ProtocolParser},
    query::QueryResult,
    realtime::StreamMessage,
};
use uuid::Uuid;

pub struct StreamCommit<'a> {
    pub lease: &'a Lease,
    pub session_id: Uuid,
    pub sequence: i64,
}
impl Postgres {
    pub async fn stream_ready(&self, l: &Lease) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        sqlx::query("SELECT committed_seq FROM trade_log.ingestion_state WHERE id=1 FOR UPDATE")
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn stream_state(&self, l: &Lease, patch: Value) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2,updated_at=now() WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(patch).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn prepare_stream(
        &self,
        l: &Lease,
        c: &CollectionConfig,
        enabled: bool,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        let interrupted:Vec<Uuid> = sqlx::query_scalar("UPDATE trade_log.collection_stream_sessions SET status='INTERRUPTED',closed_at=now(),close_reason='PROCESS_INTERRUPTED' WHERE checkpoint_key=$1 AND closed_at IS NULL RETURNING session_id").bind(&l.key).fetch_all(&mut *tx).await.map_err(db_error)?;
        for id in interrupted {
            self.open_gap_tx(&mut tx, l, c, id, "PROCESS_INTERRUPTED")
                .await?;
        }
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2 WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(json!({"enabled":enabled,"status":if enabled{"CONNECTING"}else{"DISABLED"},"pending_messages":0,"pending_bytes":0,"next_retry_at":null})).execute(&mut *tx).await.map_err(db_error)?;
        if !enabled {
            sqlx::query("UPDATE trade_log.collection_gaps SET end_ms=COALESCE(end_ms,$2),status=CASE WHEN status='BLOCKED' THEN status ELSE 'SCANNING' END WHERE checkpoint_key=$1 AND status IN ('OPEN','SCANNING','BLOCKED')").bind(&l.key).bind(Utc::now().timestamp_millis()).execute(&mut *tx).await.map_err(db_error)?;
        }
        refresh_recovery(&mut tx, &l.key).await?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn begin_stream(
        &self,
        l: &Lease,
        c: &CollectionConfig,
        id: Uuid,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("INSERT INTO trade_log.collection_stream_sessions(session_id,checkpoint_key,lease_epoch,account,network,status) VALUES($1,$2,$3,$4,$5,'CONNECTING')").bind(id).bind(&l.key).bind(l.epoch).bind(&c.account).bind(self.network.name()).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2 WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(json!({"enabled":true,"status":"CONNECTING","session_id":id,"connected_at":null,"subscribed_at":null,"last_pong_at":null,"next_retry_at":null})).execute(&mut *tx).await.map_err(db_error)?;
        self.open_gap_tx(&mut tx, l, c, id, "STARTUP").await?;
        refresh_recovery(&mut tx, &l.key).await?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn stream_connected(&self, l: &Lease, id: Uuid) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("UPDATE trade_log.collection_stream_sessions SET connected_at=now(),status='SUBSCRIBING' WHERE session_id=$1 AND lease_epoch=$2 AND closed_at IS NULL").bind(id).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || jsonb_build_object('status','SUBSCRIBING','connected_at',$2::text,'connection_count',COALESCE((websocket_state->>'connection_count')::bigint,0)+1,'reconnect_count',COALESCE((websocket_state->>'connection_count')::bigint,0)) WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(shared_types::now()).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn stream_subscribed(&self, l: &Lease, id: Uuid) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        let now = Utc::now();
        sqlx::query("UPDATE trade_log.collection_stream_sessions SET subscribed_at=$2,status='LIVE' WHERE session_id=$1 AND lease_epoch=$3").bind(id).bind(now).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2 WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(json!({"status":"LIVE","subscribed_at":shared_types::now(),"last_error":null})).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("UPDATE trade_log.collection_gaps SET end_ms=GREATEST(COALESCE(end_ms,$2),$2),status=CASE WHEN status='BLOCKED' THEN status ELSE 'SCANNING' END WHERE checkpoint_key=$1 AND status<>'HTTP_SCANNED'").bind(&l.key).bind(now.timestamp_millis()).execute(&mut *tx).await.map_err(db_error)?;
        refresh_recovery(&mut tx, &l.key).await?;
        tx.commit().await.map_err(db_error)
    }
    async fn stream_session_patch(
        &self,
        l: &Lease,
        id: Uuid,
        column: &str,
        status: &str,
        patch: Value,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        // column is selected only by the two internal callers, never by input.
        sqlx::query(&format!("UPDATE trade_log.collection_stream_sessions SET {column}=now(),status=$2 WHERE session_id=$1 AND lease_epoch=$3 AND closed_at IS NULL")).bind(id).bind(status).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2 WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(patch).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn stream_pong(&self, l: &Lease, id: Uuid) -> Result<(), QueryError> {
        self.stream_session_patch(
            l,
            id,
            "last_pong_at",
            "LIVE",
            json!({"last_pong_at":shared_types::now()}),
        )
        .await
    }
    pub async fn end_stream(
        &self,
        l: &Lease,
        c: &CollectionConfig,
        id: Uuid,
        e: &QueryError,
        stopped: bool,
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("UPDATE trade_log.collection_stream_sessions SET closed_at=now(),status=$2,close_reason=$3 WHERE session_id=$1 AND lease_epoch=$4").bind(id).bind(if stopped{"STOPPED"}else{"CLOSED"}).bind(&e.code).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
        self.open_gap_tx(&mut tx, l, c, id, &e.code).await?;
        let fatal = !e.retryable;
        if fatal {
            sqlx::query("UPDATE trade_log.collection_gaps SET status='BLOCKED',last_error=$2 WHERE checkpoint_key=$1 AND status<>'HTTP_SCANNED'").bind(&l.key).bind(json!(e)).execute(&mut *tx).await.map_err(db_error)?;
        }
        sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2 WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(json!({"status":if stopped{"STOPPED"}else if fatal{"FAILED"}else{"RECONNECT_WAIT"},"last_error":{"code":e.code,"message":e.message,"occurred_at":shared_types::now()},"pending_messages":0,"pending_bytes":0})).execute(&mut *tx).await.map_err(db_error)?;
        refresh_recovery(&mut tx, &l.key).await?;
        tx.commit().await.map_err(db_error)
    }
    async fn open_gap_tx(
        &self,
        tx: &mut Transaction<'_, Pg>,
        l: &Lease,
        c: &CollectionConfig,
        id: Uuid,
        reason: &str,
    ) -> Result<(), QueryError> {
        let position:Value=sqlx::query_scalar("SELECT position FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).fetch_one(&mut **tx).await.map_err(db_error)?;
        let initial = position["initial_start_ms"]
            .as_i64()
            .ok_or_else(QueryError::storage)?;
        let start = initial.max(
            position["scanned_through_ms"]
                .as_i64()
                .unwrap_or(initial)
                .saturating_sub((c.overlap_seconds * 1000) as i64),
        );
        // All pending recovery ranges end at a future subscription boundary, so coalesce them.
        let existing:Option<Uuid>=sqlx::query_scalar("SELECT gap_id FROM trade_log.collection_gaps WHERE checkpoint_key=$1 AND status<>'HTTP_SCANNED' ORDER BY detected_at LIMIT 1 FOR UPDATE").bind(&l.key).fetch_optional(&mut **tx).await.map_err(db_error)?;
        if let Some(gap) = existing {
            sqlx::query("UPDATE trade_log.collection_gaps SET start_ms=LEAST(start_ms,$2),end_ms=NULL,session_ids=CASE WHEN session_ids @> $3 THEN session_ids ELSE session_ids || $3 END,status=CASE WHEN status='BLOCKED' THEN status ELSE 'OPEN' END,reason=$4 WHERE gap_id=$1").bind(gap).bind(start).bind(json!([id])).bind(reason).execute(&mut **tx).await.map_err(db_error)?;
        } else {
            sqlx::query("INSERT INTO trade_log.collection_gaps(gap_id,checkpoint_key,session_ids,reason,start_ms,status) VALUES($1,$2,$3,$4,$5,'OPEN')").bind(Uuid::new_v4()).bind(&l.key).bind(json!([id])).bind(reason).bind(start).execute(&mut **tx).await.map_err(db_error)?;
        }
        Ok(())
    }
    pub async fn archive_stream(
        &self,
        l: &Lease,
        id: Uuid,
        seq: i64,
        message: &StreamMessage,
        meta: &[u8],
        spot: &[u8],
    ) -> Result<String, QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        let query = format!("query_{}_{}", id.simple(), seq);
        let request = json!({"account":l.key.rsplit(':').next(),"meta":serde_json::from_slice::<Value>(meta).map_err(|_|QueryError::storage())?,"spot_meta":serde_json::from_slice::<Value>(spot).map_err(|_|QueryError::storage())?,"source_path":"data.fills","meta_snapshots":[metadata_snapshot(meta,spot)?],"active_meta_snapshot":0});
        sqlx::query("INSERT INTO trade_log.collection_jobs(id,query_id,chain_id,source_id,account,network,mode,status,request,trace_id,job_origin,checkpoint_key,lease_epoch,transport,session_id,message_sequence,message_mode) VALUES($1,$2,$3,'official_ws',$4,$5,'REALTIME','RUNNING',$6,$7,'COLLECTOR',$8,$9,'WEBSOCKET',$10,$11,$12) ON CONFLICT(query_id) DO NOTHING")
            .bind(Uuid::new_v4()).bind(&query).bind(format!("hyperliquid:{}",self.network.name())).bind(l.key.rsplit(':').next().ok_or_else(QueryError::storage)?).bind(self.network.name()).bind(request).bind(format!("trace_{}",Uuid::new_v4().simple())).bind(&l.key).bind(l.epoch).bind(id).bind(seq).bind(&message.mode).execute(&mut *tx).await.map_err(db_error)?;
        let raw = format!("raw_{query}_userFills_1");
        let observed = DateTime::parse_from_rfc3339(&message.received_at)
            .map_err(|_| QueryError::storage())?
            .with_timezone(&Utc);
        let inserted=sqlx::query("INSERT INTO trade_log.raw_logs(id,collection_job_id,chain_id,source_id,source_event_id,kind,attempt,chain_position,ordering_key,payload,body,http_status,request,sha256,observed_at,transport,session_id,message_sequence,message_mode) SELECT $1,id,chain_id,'official_ws',$3,'userFills',1,$4,$5,$6,$7,NULL,$8,$9,$10,'WEBSOCKET',$11,$12,$13 FROM trade_log.collection_jobs WHERE query_id=$2 ON CONFLICT(session_id,message_sequence) WHERE session_id IS NOT NULL DO NOTHING")
            .bind(Uuid::new_v4()).bind(&query).bind(&raw).bind(json!({"session_id":id,"message_sequence":seq,"source_path":"data.fills"})).bind(format!("{}:{raw}",message.received_at)).bind(serde_json::from_slice::<Value>(&message.body).ok()).bind(&message.body).bind(json!({"subscription":{"type":"userFills","user":l.key.rsplit(':').next(),"aggregateByTime":false}})).bind(format!("{:x}",Sha256::digest(&message.body))).bind(observed).bind(id).bind(seq).bind(&message.mode).execute(&mut *tx).await.map_err(db_error)?;
        if inserted.rows_affected() == 0 {
            let existing:(Vec<u8>,String)=sqlx::query_as("SELECT body,message_mode FROM trade_log.raw_logs WHERE session_id=$1 AND message_sequence=$2").bind(id).bind(seq).fetch_one(&mut *tx).await.map_err(db_error)?;
            if existing != (message.body.clone(), message.mode.clone()) {
                return Err(QueryError::conflict());
            }
        }
        if inserted.rows_affected() > 0 {
            sqlx::query("UPDATE trade_log.collection_stream_sessions SET last_received_at=$2,received_messages=received_messages+1 WHERE session_id=$1").bind(id).bind(observed).execute(&mut *tx).await.map_err(db_error)?;
            sqlx::query("UPDATE trade_log.collection_checkpoints SET websocket_state=websocket_state || $2 WHERE partition_key=$1 AND source_id='official_http'").bind(&l.key).bind(json!({"last_received_at":message.received_at})).execute(&mut *tx).await.map_err(db_error)?;
        }
        tx.commit().await.map_err(db_error)?;
        if let Some(m) = &self.mirror
            && let Err(e) = m
                .stream_response(&query, &raw, id, seq, message, (meta, spot))
                .await
        {
            tracing::warn!(code=%e.code,query_id=%query,"evidence_mirror_failed");
        }
        Ok(query)
    }
    pub async fn stream_result(
        &self,
        id: &str,
        parser: &dyn ProtocolParser,
    ) -> Result<QueryResult, QueryError> {
        let row=sqlx::query("SELECT j.account,j.request,r.source_event_id,r.body,r.sha256 FROM trade_log.collection_jobs j JOIN trade_log.raw_logs r ON r.collection_job_id=j.id WHERE j.query_id=$1 AND j.transport='WEBSOCKET'").bind(id).fetch_one(&self.pool).await.map_err(db_error)?;
        let bytes: Vec<u8> = row.get("body");
        if format!("{:x}", Sha256::digest(&bytes)) != row.get::<String, _>("sha256") {
            return Err(QueryError::incomplete("Stored response checksum mismatch"));
        }
        let envelope: Value = serde_json::from_slice(&bytes)
            .map_err(|_| QueryError::incomplete("Invalid websocket JSON"))?;
        let account: String = row.get("account");
        if envelope["channel"] != "userFills"
            || envelope["data"]["user"]
                .as_str()
                .and_then(|s| shared_types::account(s).ok())
                .as_deref()
                != Some(&account)
            || !envelope["data"]["fills"].is_array()
        {
            return Err(QueryError::incomplete("Invalid websocket account or fills"));
        }
        let req: Value = row.get("request");
        let snapshot =
            &req["meta_snapshots"][req["active_meta_snapshot"].as_u64().unwrap_or(0) as usize];
        let meta = verified_metadata(snapshot, "meta")?;
        let spot = verified_metadata(snapshot, "spot_meta")?;

        let parsed = parser.parse(ParseContext {
            network: &self.network,
            account: &account,
            raw_log_id: &row.get::<String, _>("source_event_id"),
            fills: &serde_json::to_vec(&envelope["data"]["fills"])
                .map_err(|_| QueryError::storage())?,
            meta: &meta,
            spot_meta: &spot,
        })?;
        let (parsed, _) = trade_log::collection::merge(vec![parsed])?;
        let mut result = trade_log::normalization::result(
            id,
            &account,
            &self.network,
            &shared_types::now(),
            usize::MAX,
            parsed,
        )?;
        result.query_scope = "REALTIME_MESSAGE".into();
        result.coverage = "SOURCE_HISTORY_NOT_VERIFIED".into();
        result.warnings.push("SOURCE_HISTORY_NOT_VERIFIED".into());
        Ok(result)
    }
    pub async fn update_stream_metadata(
        &self,
        l: &Lease,
        id: &str,
        meta: &[u8],
        spot: &[u8],
    ) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        let row:Value=sqlx::query_scalar("SELECT request FROM trade_log.collection_jobs WHERE query_id=$1 AND checkpoint_key=$2 AND transport='WEBSOCKET' AND status<>'COMPLETED' FOR UPDATE").bind(id).bind(&l.key).fetch_one(&mut *tx).await.map_err(db_error)?;
        let mut request = row;
        let snapshots = request["meta_snapshots"]
            .as_array_mut()
            .ok_or_else(QueryError::storage)?;
        let snapshot = metadata_snapshot(meta, spot)?;
        let index = if let Some(index) = snapshots.iter().position(|s| s == &snapshot) {
            index
        } else {
            snapshots.push(snapshot);
            snapshots.len() - 1
        };
        request["active_meta_snapshot"] = json!(index);
        sqlx::query("UPDATE trade_log.collection_jobs SET request=$2 WHERE query_id=$1")
            .bind(id)
            .bind(request)
            .execute(&mut *tx)
            .await
            .map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
    pub async fn stream_fatal_error(&self, l: &Lease) -> Result<Option<QueryError>, QueryError> {
        let value:Option<Value>=sqlx::query_scalar("SELECT error FROM trade_log.collection_jobs WHERE checkpoint_key=$1 AND transport='WEBSOCKET' AND status='FAILED' AND error->>'code'='VERSION_CONFLICT' ORDER BY created_at DESC LIMIT 1").bind(&l.key).fetch_optional(&self.pool).await.map_err(db_error)?.flatten();
        value
            .map(|v| serde_json::from_value(v).map_err(|_| QueryError::storage()))
            .transpose()
    }
    pub async fn fail_stream(&self, l: &Lease, id: &str, e: &QueryError) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("UPDATE trade_log.collection_jobs SET status='FAILED',error=$2,updated_at=now() WHERE query_id=$1 AND status='RUNNING' AND checkpoint_key=$3 AND lease_epoch=$4").bind(id).bind(json!(e)).bind(&l.key).bind(l.epoch).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)?;
        if let Some(m) = &self.mirror {
            use trade_log::raw_log::RawEvidenceStore;
            let _ = m.finish(id, Err(e)).await;
        }
        Ok(())
    }
    pub async fn pending_stream_jobs(
        &self,
        l: &Lease,
    ) -> Result<Vec<(String, Uuid, i64)>, QueryError> {
        sqlx::query_as("SELECT query_id,session_id,message_sequence FROM trade_log.collection_jobs WHERE checkpoint_key=$1 AND transport='WEBSOCKET' AND status<>'COMPLETED' AND (error IS NULL OR error->>'code' IN ('DEPENDENCY_UNAVAILABLE','LEASE_LOST','COMMIT_TIMEOUT','INCOMPLETE_DATA')) ORDER BY created_at LIMIT 256").bind(&l.key).fetch_all(&self.pool).await.map_err(db_error)
    }
    pub async fn resume_stream_job(&self, l: &Lease, id: &str) -> Result<(), QueryError> {
        let mut tx = self.pool.begin().await.map_err(db_error)?;
        self.lock_collection(&mut tx, l).await?;
        sqlx::query("UPDATE trade_log.collection_jobs SET status='RUNNING',lease_epoch=$2,error=NULL WHERE query_id=$1 AND status<>'COMPLETED' AND checkpoint_key=$3 AND transport='WEBSOCKET'").bind(id).bind(l.epoch).bind(&l.key).execute(&mut *tx).await.map_err(db_error)?;
        tx.commit().await.map_err(db_error)
    }
}
/// Called in the very same transaction as the successful HTTP watermark commit.
pub async fn http_recovered(
    tx: &mut Transaction<'_, Pg>,
    key: &str,
    end: i64,
) -> Result<(), QueryError> {
    sqlx::query("UPDATE trade_log.collection_gaps SET status='HTTP_SCANNED',scanned_at=now(),last_error=NULL WHERE checkpoint_key=$1 AND (status IN ('OPEN','SCANNING') OR (status='BLOCKED' AND last_error->>'code'='DEPENDENCY_UNAVAILABLE')) AND end_ms IS NOT NULL AND end_ms<=$2").bind(key).bind(end).execute(&mut **tx).await.map_err(db_error)?;
    refresh_recovery(tx, key).await
}
pub async fn refresh_recovery(tx: &mut Transaction<'_, Pg>, key: &str) -> Result<(), QueryError> {
    let row=sqlx::query("SELECT count(*) FILTER(WHERE status<>'HTTP_SCANNED') AS open_count,max(end_ms) FILTER(WHERE status<>'HTTP_SCANNED') AS target,bool_or(status='BLOCKED') AS blocked,max(scanned_at) AS scanned FROM trade_log.collection_gaps WHERE checkpoint_key=$1").bind(key).fetch_one(&mut **tx).await.map_err(db_error)?;
    let enabled:bool=sqlx::query_scalar("SELECT COALESCE((websocket_state->>'enabled')::boolean,false) FROM trade_log.collection_checkpoints WHERE partition_key=$1 AND source_id='official_http'").bind(key).fetch_one(&mut **tx).await.map_err(db_error)?;
    let count: i64 = row.get("open_count");
    let target: Option<i64> = row.get("target");
    let scanned: Option<DateTime<Utc>> = row.get("scanned");
    let state = if row.get::<Option<bool>, _>("blocked") == Some(true) {
        "BLOCKED"
    } else if count > 0 {
        if target.is_some() {
            "SCANNING"
        } else {
            "NOT_STARTED"
        }
    } else if enabled {
        "HTTP_SCANNED"
    } else {
        "DISABLED"
    };
    let error:Option<Value>=sqlx::query_scalar("SELECT last_error FROM trade_log.collection_gaps WHERE checkpoint_key=$1 AND status='BLOCKED' ORDER BY detected_at DESC LIMIT 1").bind(key).fetch_optional(&mut **tx).await.map_err(db_error)?.flatten();
    sqlx::query("UPDATE trade_log.collection_checkpoints SET recovery_state=$2 WHERE partition_key=$1 AND source_id='official_http'").bind(key).bind(json!({"status":state,"target_through":target.map(trade_log::collection::time),"open_gap_count":count,"last_scanned_at":scanned.map(|t|t.to_rfc3339_opts(chrono::SecondsFormat::Millis,true)),"last_error":error})).execute(&mut **tx).await.map_err(db_error)?;
    Ok(())
}

fn metadata_snapshot(meta: &[u8], spot: &[u8]) -> Result<Value, QueryError> {
    Ok(
        json!({"meta":std::str::from_utf8(meta).map_err(|_|QueryError::storage())?,"spot_meta":std::str::from_utf8(spot).map_err(|_|QueryError::storage())?,"meta_sha256":format!("{:x}",Sha256::digest(meta)),"spot_meta_sha256":format!("{:x}",Sha256::digest(spot))}),
    )
}
fn verified_metadata(snapshot: &Value, kind: &str) -> Result<Vec<u8>, QueryError> {
    let bytes = snapshot[kind]
        .as_str()
        .ok_or_else(|| QueryError::incomplete("Missing archived market metadata"))?
        .as_bytes()
        .to_vec();
    if snapshot[format!("{kind}_sha256")].as_str() != Some(&format!("{:x}", Sha256::digest(&bytes)))
    {
        return Err(QueryError::incomplete("Metadata checksum mismatch"));
    }
    Ok(bytes)
}
