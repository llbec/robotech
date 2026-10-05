ALTER TABLE trade_log.collection_jobs ADD COLUMN job_origin TEXT NOT NULL DEFAULT 'MANUAL' CHECK(job_origin IN ('MANUAL','COLLECTOR'));
ALTER TABLE trade_log.collection_jobs ADD COLUMN checkpoint_key TEXT;
ALTER TABLE trade_log.collection_jobs ADD COLUMN lease_epoch BIGINT;
ALTER TABLE trade_log.raw_logs ADD COLUMN page_no INT NOT NULL DEFAULT 0 CHECK(page_no>=0);
ALTER TABLE trade_log.raw_logs DROP CONSTRAINT raw_logs_collection_job_id_kind_attempt_key;
ALTER TABLE trade_log.raw_logs ADD UNIQUE(collection_job_id,kind,page_no,attempt);
CREATE TABLE trade_log.collection_checkpoints (
 chain_id TEXT NOT NULL, source_id TEXT NOT NULL DEFAULT 'official_http', partition_key TEXT NOT NULL,
 position JSONB NOT NULL, cursor TEXT, finality_proof JSONB,
 status TEXT NOT NULL DEFAULT 'STARTING' CHECK(status IN ('STARTING','RUNNING','WAITING','RETRY_WAIT','FAILED','STOPPED')),
 pending_work JSONB,
 last_attempt_at TIMESTAMPTZ, last_success_at TIMESTAMPTZ,
 last_query_id TEXT REFERENCES trade_log.collection_jobs(query_id),
 last_success_query_id TEXT REFERENCES trade_log.collection_jobs(query_id),
 consecutive_failures INT NOT NULL DEFAULT 0 CHECK(consecutive_failures>=0),
 next_run_at TIMESTAMPTZ, last_error JSONB,
 lease_owner UUID, lease_epoch BIGINT NOT NULL DEFAULT 0 CHECK(lease_epoch>=0),
 heartbeat_at TIMESTAMPTZ, lease_expires_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 PRIMARY KEY(chain_id,source_id,partition_key)
);
CREATE INDEX checkpoints_schedule ON trade_log.collection_checkpoints(status,next_run_at);
