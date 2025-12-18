package txstore

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func OpenSQLite(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// 强烈建议的性能配置
	pragmas := []string{
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA synchronous = NORMAL;`,
		`PRAGMA temp_store = MEMORY;`,
		`PRAGMA cache_size = -200000;`, // ~200MB
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, err
		}
	}

	if _, err := db.Exec(createTxEventTable); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createTxLogTable); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}
