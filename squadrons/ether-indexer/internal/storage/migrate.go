package storage

import "database/sql"

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS checkpoints (
			project_id TEXT PRIMARY KEY,
			start_block INTEGER,
			current_block INTEGER,
			step INTEGER,
			interval_sec INTEGER,
			addresses JSON,
			topics JSON
		);`,
		`CREATE TABLE IF NOT EXISTS transactions (
			project_id TEXT,
			tx_hash TEXT,
			block_number INTEGER,
			block_time TIMESTAMP,
			tx_index INTEGER,
			from_address TEXT,
			to_address TEXT,
			value TEXT,
			gas INTEGER,
			gas_price TEXT,
			nonce INTEGER,
			input_data BLOB,
			status INTEGER,
			created_at TIMESTAMP,
			PRIMARY KEY (project_id, tx_hash)
		);`,
		`CREATE TABLE IF NOT EXISTS transaction_logs (
			project_id TEXT,
			tx_hash TEXT,
			log_index INTEGER,
			address TEXT,
			topics JSON,
			data BLOB,
			PRIMARY KEY (project_id, tx_hash, log_index)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_scan
         ON transactions(project_id, block_number, tx_index);`,
		`CREATE INDEX IF NOT EXISTS idx_project_block_number
         ON transactions(project_id, block_number);`,
		`CREATE INDEX IF NOT EXISTS idx_project_block_time
         ON transactions(project_id, block_time);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
