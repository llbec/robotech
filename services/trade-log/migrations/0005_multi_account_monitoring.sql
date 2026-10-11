CREATE TABLE trade_log.monitoring_settings (
 id integer PRIMARY KEY CHECK(id=1), max_active integer NOT NULL DEFAULT 5 CHECK(max_active BETWEEN 1 AND 10),
 max_registered integer NOT NULL DEFAULT 100 CHECK(max_registered BETWEEN 1 AND 10000),
 http_budget integer NOT NULL DEFAULT 600 CHECK(http_budget BETWEEN 120 AND 1200),
 ws_budget integer NOT NULL DEFAULT 15 CHECK(ws_budget BETWEEN 1 AND 30),
 stale_seconds integer NOT NULL DEFAULT 30, imported boolean NOT NULL DEFAULT false
);
INSERT INTO trade_log.monitoring_settings(id) VALUES(1);
CREATE TABLE trade_log.watch_accounts (
 id uuid PRIMARY KEY, account text NOT NULL, account_key text NOT NULL UNIQUE,
 network text NOT NULL, start_time timestamptz NOT NULL,
 desired_state text NOT NULL DEFAULT 'ACTIVE' CHECK(desired_state IN ('ACTIVE','PAUSED')),
 generation bigint NOT NULL DEFAULT 1, created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(), paused_at timestamptz, resumed_at timestamptz,
 last_error text, last_error_at timestamptz
);
CREATE INDEX watch_accounts_order ON trade_log.watch_accounts(created_at,id);
-- Existing checkpoints are imported once without changing checkpoints or publishing controls.
INSERT INTO trade_log.watch_accounts(id,account,account_key,network,start_time)
SELECT gen_random_uuid(),split_part(partition_key,':',4),split_part(partition_key,':',1)||':'||split_part(partition_key,':',2)||':'||split_part(partition_key,':',3)||':'||split_part(partition_key,':',4),split_part(partition_key,':',2),to_timestamp((position->>'initial_start_ms')::double precision/1000)
FROM trade_log.collection_checkpoints WHERE source_id='official_http';
UPDATE trade_log.monitoring_settings SET imported=EXISTS(SELECT FROM trade_log.watch_accounts);
CREATE TABLE trade_log.source_budget_reservations (
 id uuid PRIMARY KEY, network text NOT NULL, kind text NOT NULL CHECK(kind IN ('HTTP','WS')),
 weight integer NOT NULL, reserved_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
CREATE INDEX source_budget_window ON trade_log.source_budget_reservations(network,kind,reserved_at);
ALTER TABLE trade_log.publishing_control ADD COLUMN monitoring_generation bigint NOT NULL DEFAULT 1;

ALTER TABLE trade_log.collection_stream_sessions ADD COLUMN monitoring_generation bigint;
ALTER TABLE trade_log.monitoring_settings ADD COLUMN max_http_requests integer NOT NULL DEFAULT 2;
ALTER TABLE trade_log.source_budget_reservations ADD COLUMN active_until timestamptz;
