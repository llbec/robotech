CREATE TABLE trade_log.instrument_activation (
 account_key text NOT NULL, instrument_type text NOT NULL CHECK(instrument_type='SPOT'),
 enabled boolean NOT NULL, activation_epoch bigint NOT NULL DEFAULT 0,
 activated_at timestamptz NOT NULL DEFAULT clock_timestamp(), monitoring_generation bigint NOT NULL,
 policy jsonb NOT NULL, PRIMARY KEY(account_key,instrument_type)
);
CREATE TABLE trade_log.spot_backfill_jobs (
 id uuid PRIMARY KEY, account_key text NOT NULL REFERENCES trade_log.watch_accounts(account_key),
 start_ms bigint NOT NULL CHECK(start_ms>=0), end_ms bigint NOT NULL CHECK(end_ms>start_ms),
 scanned_ms bigint NOT NULL, window_ms bigint NOT NULL DEFAULT 3600000, evidence_cursor uuid, evidence_complete boolean NOT NULL DEFAULT false,
 status text NOT NULL DEFAULT 'PENDING' CHECK(status IN ('PENDING','RUNNING','RETRY_WAIT','BLOCKED','COMPLETED')),
 lease_owner uuid, lease_until timestamptz, attempts integer NOT NULL DEFAULT 0,
 retry_at timestamptz, last_error text, source_coverage text NOT NULL DEFAULT 'SOURCE_HISTORY_NOT_VERIFIED',
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(), updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 UNIQUE(account_key,start_ms,end_ms)
);
CREATE UNIQUE INDEX spot_backfill_active ON trade_log.spot_backfill_jobs(account_key) WHERE status<>'COMPLETED';
CREATE INDEX facts_instrument_type ON trade_log.account_fact_versions ((payload->'payload'->>'instrument_type'),account_key);

ALTER TABLE trade_log.raw_logs ALTER COLUMN parser_version SET DEFAULT 'hyperliquid-v2';
