CREATE SCHEMA trade_log;
CREATE TABLE trade_log.collection_jobs (
 id UUID PRIMARY KEY, query_id TEXT UNIQUE NOT NULL,
 chain_id TEXT NOT NULL, protocol TEXT NOT NULL DEFAULT 'hyperliquid', source_id TEXT NOT NULL DEFAULT 'official_http',
 account TEXT NOT NULL, network TEXT NOT NULL, mode TEXT NOT NULL DEFAULT 'HISTORICAL',
 status TEXT NOT NULL CHECK (status IN ('RUNNING','COMPLETED','FAILED','INTERRUPTED')),
 request JSONB NOT NULL, trace_id TEXT NOT NULL, result JSONB, error JSONB,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE trade_log.raw_logs (
 id UUID PRIMARY KEY, collection_job_id UUID NOT NULL REFERENCES trade_log.collection_jobs(id),
 chain_id TEXT NOT NULL, protocol TEXT NOT NULL DEFAULT 'hyperliquid', source_id TEXT NOT NULL DEFAULT 'official_http',
 source_event_id TEXT NOT NULL, kind TEXT NOT NULL, attempt INT NOT NULL CHECK(attempt > 0),
 chain_position JSONB NOT NULL, ordering_key TEXT NOT NULL,
 payload JSONB, body BYTEA NOT NULL, http_status INT NOT NULL, request JSONB NOT NULL,
 sha256 TEXT NOT NULL, observed_at TIMESTAMPTZ NOT NULL,
 parse_status TEXT NOT NULL DEFAULT 'RECEIVED', parser_version TEXT NOT NULL DEFAULT 'hyperliquid-v1',
 UNIQUE(chain_id,source_id,source_event_id), UNIQUE(collection_job_id,kind,attempt)
);
CREATE TABLE trade_log.account_fact_versions (
 fact_id TEXT NOT NULL, revision INT NOT NULL CHECK(revision > 0), event_id TEXT UNIQUE NOT NULL,
 account_key TEXT NOT NULL, fact_type TEXT NOT NULL, ordering_key TEXT NOT NULL, sub_index INT NOT NULL,
 change_type TEXT NOT NULL, confirmation_status TEXT NOT NULL, occurred_at TIMESTAMPTZ NOT NULL,
 payload JSONB NOT NULL, raw_log_id UUID NOT NULL REFERENCES trade_log.raw_logs(id),
 parser_version TEXT NOT NULL, content_hash TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 PRIMARY KEY(fact_id,revision)
);
CREATE TABLE trade_log.account_facts_current (
 fact_id TEXT PRIMARY KEY, current_revision INT NOT NULL,
 account_key TEXT NOT NULL, fact_type TEXT NOT NULL, ordering_key TEXT NOT NULL, sub_index INT NOT NULL,
 occurred_at TIMESTAMPTZ NOT NULL, source_tid NUMERIC(20,0) NOT NULL CHECK(source_tid >= 0),
 is_retracted BOOLEAN NOT NULL DEFAULT false, ingest_seq BIGINT UNIQUE NOT NULL,
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 FOREIGN KEY(fact_id,current_revision) REFERENCES trade_log.account_fact_versions(fact_id,revision)
);
CREATE INDEX facts_time ON trade_log.account_facts_current(account_key,occurred_at DESC,source_tid DESC,fact_id);
CREATE TABLE trade_log.fact_observations (
 fact_id TEXT NOT NULL, revision INT NOT NULL, raw_log_id UUID NOT NULL REFERENCES trade_log.raw_logs(id),
 source_index INT NOT NULL CHECK(source_index >= 0),
 PRIMARY KEY(fact_id,revision,raw_log_id,source_index),
 FOREIGN KEY(fact_id,revision) REFERENCES trade_log.account_fact_versions(fact_id,revision)
);
CREATE TABLE trade_log.ingestion_state (id INT PRIMARY KEY CHECK(id=1), committed_seq BIGINT NOT NULL CHECK(committed_seq>=0));
INSERT INTO trade_log.ingestion_state VALUES(1,0);
