CREATE TABLE nansen_import_runs (
    id UUID PRIMARY KEY, address TEXT NOT NULL, from_time TIMESTAMPTZ NOT NULL, to_time TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('RUNNING','COMPLETED','FAILED')), pages_fetched INTEGER NOT NULL DEFAULT 0,
    raw_count BIGINT NOT NULL DEFAULT 0, normalized_count BIGINT NOT NULL DEFAULT 0, failed_count BIGINT NOT NULL DEFAULT 0,
    last_error TEXT, started_at TIMESTAMPTZ NOT NULL DEFAULT now(), finished_at TIMESTAMPTZ
);
CREATE TABLE nansen_raw_pages (
    run_id UUID NOT NULL REFERENCES nansen_import_runs(id), page_number INTEGER NOT NULL, response JSONB NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY (run_id,page_number)
);
CREATE TABLE raw_logs (
    id UUID PRIMARY KEY, run_id UUID NOT NULL REFERENCES nansen_import_runs(id), page_number INTEGER NOT NULL, record_index INTEGER NOT NULL,
    source_id TEXT NOT NULL, content_hash TEXT NOT NULL UNIQUE, payload JSONB NOT NULL, transaction_ref TEXT, ordering_key TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL, parse_status TEXT NOT NULL CHECK (parse_status IN ('PENDING','PARSED','FAILED')), parse_error TEXT,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX raw_logs_ordering_idx ON raw_logs(ordering_key);
CREATE TABLE account_fact_versions (
    fact_id TEXT NOT NULL, revision INTEGER NOT NULL, event_id TEXT NOT NULL UNIQUE, fact_type TEXT NOT NULL, account_key TEXT NOT NULL,
    ordering_key TEXT NOT NULL, sub_index INTEGER NOT NULL, change_type TEXT NOT NULL, confirmation_status TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL, payload JSONB NOT NULL, raw_log_id UUID REFERENCES raw_logs(id), schema_version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY (fact_id,revision)
);
CREATE INDEX account_fact_versions_account_order_idx ON account_fact_versions(account_key,ordering_key,sub_index);
CREATE INDEX account_fact_versions_type_time_idx ON account_fact_versions(fact_type,occurred_at);

