package txstore

import (
	"context"
	"time"
)

type Checkpoint struct {
	Name        string
	BlockHeight int64
	TxIndex     int64
	BlockTime   int64
}

func (s *SQLiteStore) LoadCheckpoint(
	ctx context.Context,
	name string,
) (*Checkpoint, error) {

	row := s.db.QueryRowContext(ctx, `
	SELECT name, block_height, tx_index, block_time
	FROM checkpoints WHERE name = ?`, name)

	cp := &Checkpoint{}
	err := row.Scan(
		&cp.Name,
		&cp.BlockHeight,
		&cp.TxIndex,
		&cp.BlockTime,
	)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (s *SQLiteStore) SaveCheckpoint(
	ctx context.Context,
	cp *Checkpoint,
) error {

	_, err := s.db.ExecContext(ctx, `
	INSERT INTO checkpoints (name, block_height, tx_index, block_time, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(name) DO UPDATE SET
		block_height = excluded.block_height,
		tx_index = excluded.tx_index,
		block_time = excluded.block_time,
		updated_at = excluded.updated_at
	`,
		cp.Name,
		cp.BlockHeight,
		cp.TxIndex,
		cp.BlockTime,
		time.Now().Unix(),
	)
	return err
}
