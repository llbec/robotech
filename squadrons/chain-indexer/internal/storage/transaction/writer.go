package transaction

import (
	"database/sql"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) (*Writer, error) {
	if _, err := db.Exec(CreateTableSQL); err != nil {
		return nil, err
	}
	return &Writer{db: db}, nil
}

func (w *Writer) BatchInsert(txs []*Transaction) error {
	tx, err := w.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
INSERT OR IGNORE INTO transactions (
	tx_hash, project_id, block_number, block_time, tx_index,
	from_address, to_address, value,
	gas, gas_price, nonce,
	input_data, status, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, t := range txs {
		_, err = stmt.Exec(
			t.TxHash, t.ProjectID, t.BlockNumber, t.BlockTime, t.TxIndex,
			t.FromAddress, t.ToAddress, t.Value,
			t.Gas, t.GasPrice, t.Nonce,
			t.InputData, t.Status, t.CreatedAt,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
