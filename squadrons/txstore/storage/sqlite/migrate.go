package sqlite

import "database/sql"

func migrate(db *sql.DB) error {
	stmts := []string{
		`PRAGMA journal_mode=WAL;`,
		`PRAGMA synchronous=NORMAL;`,

		`CREATE TABLE IF NOT EXISTS transactions (
            project_id TEXT,
            tx_hash TEXT,
            block_height INTEGER,
            tx_index INTEGER,
            block_time INTEGER,
            sender TEXT,
            receiver TEXT,
            value TEXT,
            input BLOB,
            status INTEGER,
            tx_type TEXT,
            tx_action TEXT,
            PRIMARY KEY (project_id, tx_hash)
        );`,

		`CREATE TABLE IF NOT EXISTS logs (
            project_id TEXT,
            tx_hash TEXT,
            log_index INTEGER,
            contract TEXT,
            topics TEXT,
            data BLOB,
            PRIMARY KEY (project_id, tx_hash, log_index)
        );`,

		`CREATE TABLE IF NOT EXISTS checkpoints (
            project_id TEXT,
            consumer TEXT,
            block_height INTEGER,
            tx_index INTEGER,
            updated_at INTEGER,
            PRIMARY KEY (project_id, consumer)
        );`,

		`CREATE INDEX IF NOT EXISTS idx_scan
         ON transactions(project_id, block_height, tx_index);`,
	}

	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
