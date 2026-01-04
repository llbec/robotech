package storage

import (
	"database/sql"

	"ether-indexer/internal/model"
)

type TxRepository struct {
	shards *ShardManager
}

func NewTxRepository(sm *ShardManager) *TxRepository {
	return &TxRepository{shards: sm}
}

func (r *TxRepository) SaveBatch(
	block uint64,
	txs []model.Transaction,
	logs []model.TransactionLog,
	updateCP func(tx *sql.Tx) error,
) error {

	db, err := r.shards.GetDB(block)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txStmt, _ := tx.Prepare(`
	INSERT OR IGNORE INTO transactions
	(project_id, tx_hash, block_number, block_time, tx_index,
	 from_address, to_address, value, gas, gas_price, nonce,
	 input_data, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	defer txStmt.Close()

	for _, t := range txs {
		if _, err := txStmt.Exec(
			t.ProjectID, t.TxHash, t.BlockNumber, t.BlockTime, t.TxIndex,
			t.FromAddress, t.ToAddress, t.Value,
			t.Gas, t.GasPrice, t.Nonce,
			t.InputData, t.Status, t.CreatedAt,
		); err != nil {
			return err
		}
	}

	logStmt, _ := tx.Prepare(`
	INSERT OR IGNORE INTO transaction_logs
	(project_id, tx_hash, log_index, address, topics, data)
	VALUES (?, ?, ?, ?, ?, ?)
	`)
	defer logStmt.Close()

	for _, l := range logs {
		if _, err := logStmt.Exec(
			l.ProjectID, l.TxHash, l.LogIndex,
			l.Address, l.Topics, l.Data,
		); err != nil {
			return err
		}
	}

	if err := updateCP(tx); err != nil {
		return err
	}

	return tx.Commit()
}
