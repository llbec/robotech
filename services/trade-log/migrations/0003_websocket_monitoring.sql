ALTER TABLE trade_log.collection_checkpoints ADD COLUMN websocket_state JSONB NOT NULL DEFAULT '{"enabled":false,"status":"DISABLED","reconnect_count":0}';
ALTER TABLE trade_log.collection_checkpoints ADD COLUMN recovery_state JSONB NOT NULL DEFAULT '{"status":"DISABLED"}';
CREATE TABLE trade_log.collection_stream_sessions (
 session_id UUID PRIMARY KEY, checkpoint_key TEXT NOT NULL, lease_epoch BIGINT NOT NULL,
 account TEXT NOT NULL, network TEXT NOT NULL,
 status TEXT NOT NULL CHECK(status IN ('CONNECTING','SUBSCRIBING','LIVE','CLOSED','STOPPED','INTERRUPTED')),
 connected_at TIMESTAMPTZ, subscribed_at TIMESTAMPTZ, closed_at TIMESTAMPTZ,
 close_reason TEXT, last_received_at TIMESTAMPTZ, last_committed_at TIMESTAMPTZ, last_pong_at TIMESTAMPTZ,
 received_messages BIGINT NOT NULL DEFAULT 0, committed_messages BIGINT NOT NULL DEFAULT 0,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX stream_sessions_checkpoint ON trade_log.collection_stream_sessions(checkpoint_key,created_at);
CREATE TABLE trade_log.collection_gaps (
 gap_id UUID PRIMARY KEY, checkpoint_key TEXT NOT NULL, session_ids JSONB NOT NULL DEFAULT '[]', reason TEXT NOT NULL,
 start_ms BIGINT NOT NULL, end_ms BIGINT, status TEXT NOT NULL CHECK(status IN ('OPEN','SCANNING','HTTP_SCANNED','BLOCKED')),
 detected_at TIMESTAMPTZ NOT NULL DEFAULT now(), scanned_at TIMESTAMPTZ, last_error JSONB
);
CREATE INDEX gaps_checkpoint ON trade_log.collection_gaps(checkpoint_key,status);
ALTER TABLE trade_log.collection_jobs ADD COLUMN transport TEXT NOT NULL DEFAULT 'HTTP' CHECK(transport IN ('HTTP','WEBSOCKET'));
ALTER TABLE trade_log.collection_jobs ADD COLUMN session_id UUID REFERENCES trade_log.collection_stream_sessions(session_id);
ALTER TABLE trade_log.collection_jobs ADD COLUMN message_sequence BIGINT;
ALTER TABLE trade_log.collection_jobs ADD COLUMN message_mode TEXT CHECK(message_mode IN ('SNAPSHOT','LIVE_UPDATE','UNKNOWN'));
ALTER TABLE trade_log.raw_logs ADD COLUMN transport TEXT NOT NULL DEFAULT 'HTTP' CHECK(transport IN ('HTTP','WEBSOCKET'));
ALTER TABLE trade_log.raw_logs ADD COLUMN session_id UUID REFERENCES trade_log.collection_stream_sessions(session_id);
ALTER TABLE trade_log.raw_logs ADD COLUMN message_sequence BIGINT;
ALTER TABLE trade_log.raw_logs ADD COLUMN message_mode TEXT CHECK(message_mode IN ('SNAPSHOT','LIVE_UPDATE','UNKNOWN'));
ALTER TABLE trade_log.raw_logs ALTER COLUMN http_status DROP NOT NULL;
ALTER TABLE trade_log.raw_logs ADD CONSTRAINT raw_transport_status CHECK((transport='HTTP' AND http_status IS NOT NULL) OR (transport='WEBSOCKET' AND http_status IS NULL));
CREATE UNIQUE INDEX raw_stream_sequence ON trade_log.raw_logs(session_id,message_sequence) WHERE session_id IS NOT NULL;
ALTER TABLE trade_log.account_fact_versions ADD COLUMN semantic_hash_version TEXT;
ALTER TABLE trade_log.account_fact_versions ADD COLUMN semantic_content_hash TEXT;
