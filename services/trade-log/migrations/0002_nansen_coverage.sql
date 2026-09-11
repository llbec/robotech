ALTER TABLE nansen_import_runs
    ADD COLUMN endpoints_attempted INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN endpoints_succeeded INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN endpoints_failed INTEGER NOT NULL DEFAULT 0;

CREATE TABLE nansen_raw_responses (
    id UUID PRIMARY KEY,
    run_id UUID NOT NULL REFERENCES nansen_import_runs(id),
    endpoint TEXT NOT NULL,
    response_kind TEXT NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL,
    response JSONB NOT NULL
);

CREATE INDEX nansen_raw_responses_run_endpoint_idx
    ON nansen_raw_responses(run_id, endpoint);

CREATE TABLE nansen_coverage_results (
    run_id UUID NOT NULL REFERENCES nansen_import_runs(id),
    category TEXT NOT NULL,
    endpoint TEXT,
    status TEXT NOT NULL CHECK (status IN ('COMPLETE','PARTIAL','UNAVAILABLE','UNVERIFIED')),
    record_count BIGINT NOT NULL DEFAULT 0,
    range_start TIMESTAMPTZ,
    range_end TIMESTAMPTZ,
    missing_fields TEXT[] NOT NULL DEFAULT '{}',
    evidence JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (run_id, category)
);
