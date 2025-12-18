package txstore

const createTxEventTable = `
CREATE TABLE IF NOT EXISTS tx_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,

	block_height INTEGER NOT NULL,
	block_time   INTEGER NOT NULL,
	tx_index     INTEGER NOT NULL,
	tx_hash      TEXT    NOT NULL,

	tx_type TEXT,

	from_address TEXT,
	to_address   TEXT,

	day    TEXT,
	hour   INTEGER,
	minute INTEGER,
	second INTEGER,

	UNIQUE(block_height, tx_index),
	UNIQUE(tx_hash)
);

CREATE INDEX IF NOT EXISTS idx_tx_events_time 
ON tx_events(block_time);

CREATE INDEX IF NOT EXISTS idx_tx_events_day_type 
ON tx_events(day, tx_type);
`

const createTxLogTable = `
CREATE TABLE IF NOT EXISTS tx_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,

	tx_hash TEXT NOT NULL,
	log_index INTEGER NOT NULL,

	address TEXT,
	topics  TEXT,    -- JSON
	data    BLOB,

	block_number INTEGER,
	tx_index INTEGER,
	block_hash TEXT,
	block_timestamp INTEGER,
	removed INTEGER,

	UNIQUE(tx_hash, log_index)
);

CREATE INDEX IF NOT EXISTS idx_tx_logs_txhash 
ON tx_logs(tx_hash);

CREATE INDEX IF NOT EXISTS idx_tx_logs_address 
ON tx_logs(address);
`

const createCheckpointTable = `
CREATE TABLE IF NOT EXISTS checkpoints (
	name TEXT PRIMARY KEY,   -- e.g. "scanner", "state_replay"
	block_height INTEGER,
	tx_index INTEGER,
	block_time INTEGER,
	updated_at INTEGER
);
`
