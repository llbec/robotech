package log

import "database/sql"

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) (*Writer, error) {
	if _, err := db.Exec(CreateTableSQL); err != nil {
		return nil, err
	}
	return &Writer{db: db}, nil
}

func (w *Writer) BatchInsert(logs []*TransactionLog) error {
	tx, err := w.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
INSERT OR IGNORE INTO transaction_logs (
	project_id, tx_hash, log_index, address, topics, data
) VALUES (?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for _, l := range logs {
		_, err := stmt.Exec(
			l.ProjectID, l.TxHash, l.LogIndex,
			l.Address, l.Topics, l.Data,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
