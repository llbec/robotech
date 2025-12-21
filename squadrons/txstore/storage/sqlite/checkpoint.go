package sqlite

import (
	"database/sql"
	"time"
	"txstore/model"
)

func (s *Store) Get(project, consumer string) (*model.Checkpoint, error) {
	cp := &model.Checkpoint{}
	err := s.DB.QueryRow(`
        SELECT project_id, consumer, block_height, tx_index, updated_at
        FROM checkpoints WHERE project_id=? AND consumer=?
    `, project, consumer).Scan(
		&cp.ProjectID, &cp.Consumer,
		&cp.BlockHeight, &cp.TxIndex, &cp.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cp, err
}

func (s *Store) Save(cp *model.Checkpoint) error {
	cp.UpdatedAt = time.Now().Unix()
	_, err := s.DB.Exec(`
        INSERT INTO checkpoints VALUES (?, ?, ?, ?, ?)
        ON CONFLICT(project_id, consumer)
        DO UPDATE SET
            block_height=excluded.block_height,
            tx_index=excluded.tx_index,
            updated_at=excluded.updated_at
    `, cp.ProjectID, cp.Consumer, cp.BlockHeight, cp.TxIndex, cp.UpdatedAt)
	return err
}
