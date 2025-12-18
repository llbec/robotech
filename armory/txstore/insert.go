package txstore

import (
	"context"
	"encoding/json"
)

func (s *SQLiteStore) InsertTxEvent(ctx context.Context, ev *TxEvent) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	INSERT OR IGNORE INTO tx_events (
		block_height, block_time, tx_index, tx_hash,
		tx_type, from_address, to_address,
		day, hour, minute, second
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ev.BlockHeight,
		ev.BlockTime,
		ev.TxIndex,
		ev.TxHash,
		ev.TxType,
		ev.FromAddress,
		ev.ToAddress,
		ev.Day,
		ev.Hour,
		ev.Minute,
		ev.Second,
	)
	if err != nil {
		return err
	}

	for _, lg := range ev.Logs {
		topics, _ := json.Marshal(lg.Topics)

		_, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO tx_logs (
			tx_hash, log_index,
			address, topics, data,
			block_number, tx_index,
			block_hash, block_timestamp,
			removed
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			lg.TxHash.Hex(),
			lg.Index,
			lg.Address.Hex(),
			string(topics),
			lg.Data,
			lg.BlockNumber,
			lg.TxIndex,
			lg.BlockHash.Hex(),
			lg.BlockTimestamp,
			boolToInt(lg.Removed),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (m *StoreManager) InsertTx(ctx context.Context, ev *TxEvent) error {
	store, err := m.GetStoreByTime(ev.BlockTime)
	if err != nil {
		return err
	}
	return store.InsertTxEvent(ctx, ev)
}
