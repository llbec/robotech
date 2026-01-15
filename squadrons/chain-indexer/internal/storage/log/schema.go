package log

const CreateTableSQL = `
CREATE TABLE IF NOT EXISTS transaction_logs (
	project_id TEXT,
	tx_hash TEXT,
	log_index INTEGER,
	address TEXT,
	topics TEXT,
	data BLOB,
	PRIMARY KEY (tx_hash, log_index)
);
`
