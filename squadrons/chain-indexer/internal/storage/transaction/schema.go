package transaction

const CreateTableSQL = `
CREATE TABLE IF NOT EXISTS transactions (
	tx_hash TEXT PRIMARY KEY,
	project_id TEXT,
	block_number INTEGER,
	block_time INTEGER,
	tx_index INTEGER,

	from_address TEXT,
	to_address TEXT,
	value TEXT,

	gas INTEGER,
	gas_price TEXT,
	nonce INTEGER,

	input_data BLOB,
	status INTEGER,
	created_at INTEGER
);

CREATE INDEX IF NOT EXISTS idx_tx_time
ON transactions(project_id, block_time);
`
