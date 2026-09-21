CREATE TABLE monitored_addresses (
    address TEXT PRIMARY KEY CHECK (address ~ '^0x[0-9a-f]{40}$'),
    status TEXT NOT NULL CHECK (status IN ('PENDING','BACKFILLING','LIVE','DEGRADED','FAILED')),
    requested_start TIMESTAMPTZ NOT NULL,
    coverage_start TIMESTAMPTZ,
    history_complete BOOLEAN NOT NULL DEFAULT TRUE,
    last_event_time TIMESTAMPTZ,
    last_received_at TIMESTAMPTZ,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE collection_runs (
    id UUID PRIMARY KEY,
    address TEXT NOT NULL REFERENCES monitored_addresses(address),
    kind TEXT NOT NULL CHECK (kind IN ('BACKFILL','RECONCILE')),
    range_start TIMESTAMPTZ NOT NULL,
    range_end TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('RUNNING','COMPLETED','FAILED')),
    records_seen BIGINT NOT NULL DEFAULT 0,
    records_inserted BIGINT NOT NULL DEFAULT 0,
    error TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ
);

CREATE INDEX collection_runs_address_started_idx ON collection_runs(address, started_at DESC);

CREATE TABLE raw_trade_logs (
    id UUID PRIMARY KEY,
    address TEXT NOT NULL REFERENCES monitored_addresses(address),
    source TEXT NOT NULL,
    transport TEXT NOT NULL CHECK (transport IN ('HTTP','WEBSOCKET')),
    source_event_id TEXT NOT NULL,
    event_time TIMESTAMPTZ NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL,
    is_snapshot BOOLEAN NOT NULL,
    payload JSONB NOT NULL,
    parse_status TEXT NOT NULL CHECK (parse_status IN ('PARSED','FAILED')),
    parse_error TEXT,
    UNIQUE (source, source_event_id)
);

CREATE INDEX raw_trade_logs_account_order_idx
    ON raw_trade_logs(address, event_time, source_event_id);

CREATE TABLE market_metadata_versions (
    id BIGSERIAL PRIMARY KEY,
    observed_at TIMESTAMPTZ NOT NULL,
    perp_payload JSONB NOT NULL,
    spot_payload JSONB NOT NULL
);

CREATE TABLE trade_fact_versions (
    fact_id TEXT NOT NULL,
    revision INT NOT NULL,
    event_id TEXT NOT NULL UNIQUE,
    account_key TEXT NOT NULL,
    account TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL,
    ordering_key TEXT NOT NULL,
    source TEXT NOT NULL,
    source_ref TEXT NOT NULL,
    raw_log_id UUID NOT NULL REFERENCES raw_trade_logs(id),
    schema_version INT NOT NULL,
    payload JSONB NOT NULL,
    PRIMARY KEY (fact_id, revision)
);

CREATE INDEX trade_fact_versions_account_order_idx
    ON trade_fact_versions(account, occurred_at, ordering_key, fact_id);

CREATE TABLE trade_facts_current (
    fact_id TEXT PRIMARY KEY,
    current_revision INT NOT NULL,
    account TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    ordering_key TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (fact_id, current_revision)
        REFERENCES trade_fact_versions(fact_id, revision)
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_id TEXT NOT NULL UNIQUE,
    fact_id TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING','DELIVERING','RETRY_WAIT','DELIVERED','FAILED')),
    attempts INT NOT NULL DEFAULT 0,
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_at TIMESTAMPTZ,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    delivered_at TIMESTAMPTZ
);

CREATE INDEX outbox_ready_idx ON outbox_events(status, available_at);

CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY,
    outbox_id UUID NOT NULL REFERENCES outbox_events(id),
    attempt INT NOT NULL,
    status TEXT NOT NULL,
    http_status INT,
    response_summary TEXT,
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (outbox_id, attempt)
);

