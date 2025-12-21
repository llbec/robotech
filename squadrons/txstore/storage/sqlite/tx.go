package sqlite

import (
	"encoding/json"
	"txstore/model"
)

func (s *Store) Insert(project string, txs []*model.Transaction) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txStmt, _ := tx.Prepare(`
        INSERT OR IGNORE INTO transactions
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)

	logStmt, _ := tx.Prepare(`
        INSERT OR IGNORE INTO logs
        VALUES (?, ?, ?, ?, ?, ?)
    `)

	for _, t := range txs {
		_, err := txStmt.Exec(
			project, t.TxHash, t.BlockHeight, t.TxIndex, t.BlockTime,
			t.From, t.To, t.Value, t.Input, t.Status,
			t.TxType, t.TxAction,
		)
		if err != nil {
			return err
		}

		for _, l := range t.Logs {
			topics, _ := json.Marshal(l.Topics)
			_, err := logStmt.Exec(
				project, t.TxHash, l.LogIndex, l.Contract, string(topics), l.Data,
			)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Store) Scan(project string, fromBlock, fromIndex int64, limit int) ([]*model.Transaction, error) {
	rows, err := s.DB.Query(`
        SELECT tx_hash, block_height, tx_index, block_time,
               sender, receiver, value, input, status, tx_type, tx_action
        FROM transactions
        WHERE project_id=?
          AND (block_height > ? OR (block_height=? AND tx_index>?))
        ORDER BY block_height, tx_index
        LIMIT ?
    `, project, fromBlock, fromBlock, fromIndex, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*model.Transaction
	for rows.Next() {
		t := &model.Transaction{ProjectID: project}
		rows.Scan(
			&t.TxHash, &t.BlockHeight, &t.TxIndex, &t.BlockTime,
			&t.From, &t.To, &t.Value, &t.Input,
			&t.Status, &t.TxType, &t.TxAction,
		)
		res = append(res, t)
	}
	return res, nil
}
