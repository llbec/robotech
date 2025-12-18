CREATE TABLE IF NOT EXISTS tx_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    block_height INTEGER NOT NULL,
    block_time   INTEGER NOT NULL,
    tx_index     INTEGER NOT NULL,
    tx_hash      TEXT NOT NULL,

    tx_type      TEXT NOT NULL,

    from_address TEXT,
    to_address   TEXT,

    parsed_tx    TEXT NOT NULL,

    day          TEXT NOT NULL,
    hour         INTEGER NOT NULL,
    minute       INTEGER NOT NULL,
    second       INTEGER NOT NULL
);

-- 严格链上顺序
CREATE UNIQUE INDEX IF NOT EXISTS idx_tx_order
ON tx_events(block_height, tx_index);

-- 交易唯一性
CREATE UNIQUE INDEX IF NOT EXISTS idx_tx_hash
ON tx_events(tx_hash);

-- 时间 + 类型筛选
CREATE INDEX IF NOT EXISTS idx_tx_time_type
ON tx_events(block_time, tx_type, tx_action);

-- 报表（按天）
CREATE INDEX IF NOT EXISTS idx_tx_day_type
ON tx_events(day, tx_type);
