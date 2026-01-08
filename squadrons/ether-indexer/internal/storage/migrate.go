package storage

import "database/sql"

func dbExecArray(db *sql.DB, list []string) error {
	for _, s := range list {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func ProjectMigrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			project_id    TEXT PRIMARY KEY,
			active        BOOLEAN NOT NULL DEFAULT 0,
			rpc_endpoint  TEXT NOT NULL,
			description   TEXT,
			block_range   INTEGER NOT NULL,
			created_at    INTEGER,
			updated_at    INTEGER
		);`,
	}
	return dbExecArray(db, stmts)
}

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS checkpoints (
			project_id TEXT PRIMARY KEY,
			start_block INTEGER,
			current_block INTEGER,
			step INTEGER NOT NULL DEFAULT 9,
			interval_sec INTEGER NOT NULL DEFAULT 3,
			addresses JSON,
			topics JSON,
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
		`CREATE INDEX IF NOT EXISTS idx_tx_block
         ON transactions(project_id, block_number);`,
		`CREATE INDEX IF NOT EXISTS idx_tx_time
         ON transactions(project_id, block_time);`,
	}
	return dbExecArray(db, stmts)
}
