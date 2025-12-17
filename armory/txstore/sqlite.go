package txstore

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db   *sql.DB
	path string
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	if err := os.MkdirAll(getDir(path), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA cache_size = -200000;", // ~200MB
		"PRAGMA foreign_keys = ON;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, err
		}
	}

	return &SQLiteStore{
		db:   db,
		path: path,
	}, nil
}

func (s *SQLiteStore) Init() error {
	return initSchema(s.db)
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) InsertTx(tx *TxEvent) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO tx_events (
			block_height, block_time, tx_index, tx_hash,
			tx_type,
			from_address, to_address,
			raw_tx, parsed_tx,
			day, hour, minute, second
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		tx.BlockHeight,
		tx.BlockTime,
		tx.TxIndex,
		tx.TxHash,
		tx.TxType,
		nullIfEmpty(tx.FromAddress),
		nullIfEmpty(tx.ToAddress),
		tx.RawTx,
		tx.ParsedTx,
		tx.Day,
		tx.Hour,
		tx.Minute,
		tx.Second,
	)

	return err
}

func (s *SQLiteStore) InsertBatch(txs []*TxEvent) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO tx_events (
			block_height, block_time, tx_index, tx_hash,
			tx_type,
			from_address, to_address,
			raw_tx, parsed_tx,
			day, hour, minute, second
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range txs {
		if _, err := stmt.Exec(
			e.BlockHeight, e.BlockTime, e.TxIndex, e.TxHash,
			e.TxType,
			nullIfEmpty(e.FromAddress),
			nullIfEmpty(e.ToAddress),
			e.RawTx, e.ParsedTx,
			e.Day, e.Hour, e.Minute, e.Second,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteStore) ScanFrom(
	fromBlock int64,
	fn func(*TxEvent) error,
) error {
	rows, err := s.db.Query(`
		SELECT
			block_height, block_time, tx_index, tx_hash,
			tx_type,
			from_address, to_address,
			raw_tx, parsed_tx,
			day, hour, minute, second
		FROM tx_events
		WHERE block_height >= ?
		ORDER BY block_height ASC, tx_index ASC
	`, fromBlock)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tx TxEvent
		if err := rows.Scan(
			&tx.BlockHeight,
			&tx.BlockTime,
			&tx.TxIndex,
			&tx.TxHash,
			&tx.TxType,
			&tx.FromAddress,
			&tx.ToAddress,
			&tx.RawTx,
			&tx.ParsedTx,
			&tx.Day,
			&tx.Hour,
			&tx.Minute,
			&tx.Second,
		); err != nil {
			return err
		}

		if err := fn(&tx); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (s *SQLiteStore) QueryByTimeAndType(
	start, end int64,
	txType string,
	fn func(*TxEvent) error,
) error {

	rows, err := s.db.Query(`
		SELECT
			block_height, block_time, tx_index, tx_hash,
			tx_type,
			from_address, to_address,
			raw_tx, parsed_tx,
			day, hour, minute, second
		FROM tx_events
		WHERE block_time BETWEEN ?
		  AND ?
		  AND tx_type = ?
		ORDER BY block_time ASC
	`, start, end, txType)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tx TxEvent
		if err := rows.Scan(
			&tx.BlockHeight,
			&tx.BlockTime,
			&tx.TxIndex,
			&tx.TxHash,
			&tx.TxType,
			&tx.FromAddress,
			&tx.ToAddress,
			&tx.RawTx,
			&tx.ParsedTx,
			&tx.Day,
			&tx.Hour,
			&tx.Minute,
			&tx.Second,
		); err != nil {
			return err
		}

		if err := fn(&tx); err != nil {
			return err
		}
	}

	return rows.Err()
}

func getDir(path string) string {
	if idx := len(path) - len("/"); idx > 0 {
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				return path[:i]
			}
		}
	}
	return "."
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
