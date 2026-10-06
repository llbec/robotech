ALTER TABLE trade_log.collection_stream_sessions ADD COLUMN snapshot_sequence BIGINT;
CREATE TABLE trade_log.publishing_control (
 account_key TEXT PRIMARY KEY, target_id TEXT NOT NULL, target_url TEXT NOT NULL,
 enabled BOOLEAN NOT NULL, activation_epoch BIGINT NOT NULL DEFAULT 0,
 activated_at TIMESTAMPTZ NOT NULL DEFAULT now(), policy JSONB NOT NULL,
 heartbeat_at TIMESTAMPTZ, last_error TEXT, updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE trade_log.outbox_events (
 event_id TEXT PRIMARY KEY, fact_id TEXT NOT NULL, revision INT NOT NULL,
 topic TEXT NOT NULL, partition_key TEXT NOT NULL REFERENCES trade_log.publishing_control(account_key), target_id TEXT NOT NULL,
 payload JSONB NOT NULL, wire_body BYTEA, body_sha256 TEXT,
 status TEXT NOT NULL DEFAULT 'PENDING' CHECK(status IN ('PENDING','SENDING','RETRY_WAIT','DELIVERED','BLOCKED')),
 attempts INT NOT NULL DEFAULT 0, next_retry_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 published_at TIMESTAMPTZ, delivered_at TIMESTAMPTZ,
 lease_owner UUID, lease_epoch BIGINT NOT NULL DEFAULT 0, lease_expires_at TIMESTAMPTZ,
 last_error TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 FOREIGN KEY(fact_id,revision) REFERENCES trade_log.account_fact_versions(fact_id,revision),
 CHECK((wire_body IS NULL) = (body_sha256 IS NULL))
);
CREATE INDEX outbox_due ON trade_log.outbox_events(status,next_retry_at,lease_expires_at);
CREATE INDEX outbox_account ON trade_log.outbox_events(partition_key,created_at,event_id);
CREATE TABLE trade_log.publication_decisions (
 fact_id TEXT NOT NULL, revision INT NOT NULL, raw_log_id UUID NOT NULL REFERENCES trade_log.raw_logs(id), source_index INT NOT NULL,
 session_id UUID NOT NULL, message_sequence BIGINT NOT NULL, activation_epoch BIGINT NOT NULL,
 result TEXT NOT NULL CHECK(result IN ('ELIGIBLE','ALREADY_ENQUEUED','SUPPRESSED')), reason TEXT NOT NULL,
 decided_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(fact_id,revision,raw_log_id,source_index),
 FOREIGN KEY(fact_id,revision) REFERENCES trade_log.account_fact_versions(fact_id,revision)
);
CREATE TABLE trade_log.delivery_attempts (
 event_id TEXT NOT NULL REFERENCES trade_log.outbox_events(event_id), attempt_number INT NOT NULL,
 started_at TIMESTAMPTZ NOT NULL DEFAULT now(), finished_at TIMESTAMPTZ, http_status INT,
 result TEXT NOT NULL DEFAULT 'UNCONFIRMED', error TEXT, lease_epoch BIGINT NOT NULL,
 PRIMARY KEY(event_id,attempt_number)
);
