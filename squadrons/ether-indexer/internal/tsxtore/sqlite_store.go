package txstore

import (
	"ether-indexer/internal/model"
	"ether-indexer/internal/storage"
)

type SQLiteStore struct {
	shards *storage.ShardManager
}

func NewSQLiteStore(sm *storage.ShardManager) *SQLiteStore {
	return &SQLiteStore{shards: sm}
}

func (s *SQLiteStore) GetTransaction(projectID, hash string) (*model.Transaction, error) {
	db, _ := s.shards.GetDB(0)

	row := db.QueryRow(`
		SELECT project_id, tx_hash, block_number, block_time,
		       from_address, to_address, value
		FROM transactions
		WHERE project_id = ? AND tx_hash = ?
	`, projectID, hash)

	var t model.Transaction
	err := row.Scan(
		&t.ProjectID, &t.TxHash, &t.BlockNumber, &t.BlockTime,
		&t.FromAddress, &t.ToAddress, &t.Value,
	)
	return &t, err
}
