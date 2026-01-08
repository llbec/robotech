package storage

import (
	"database/sql"

	"ether-indexer/internal/model"
)

type CheckpointRepo struct {
	db *sql.DB
}

func NewCheckpointRepo(db *sql.DB) *CheckpointRepo {
	return &CheckpointRepo{db: db}
}

func (r *CheckpointRepo) AddCheckpoint(
	projectID string,
	startBlock uint64,
	currentBlock uint64,
	step uint64,
	intervalSec int,
	addressesJSON string,
	topicsJSON string,
) error {
	_, err := r.db.Exec(`
		INSERT INTO checkpoints
			(project_id, start_block, current_block, step,
		       interval_sec, addresses, topics)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, projectID, startBlock, currentBlock, step,
		intervalSec, addressesJSON, topicsJSON)
	return err
}

func (r *CheckpointRepo) Load(projectID string) (*model.Checkpoint, error) {
	row := r.db.QueryRow(`
		SELECT project_id, start_block, current_block, step,
		       interval_sec, addresses, topics
		FROM checkpoints WHERE project_id = ?
	`, projectID)

	cp := &model.Checkpoint{}
	err := row.Scan(
		&cp.ProjectID,
		&cp.StartBlock,
		&cp.CurrentBlock,
		&cp.Step,
		&cp.IntervalSec,
		&cp.AddressesJSON,
		&cp.TopicsJSON,
	)
	return cp, err
}

func (r *CheckpointRepo) UpdateTx(
	tx *sql.Tx,
	projectID string,
	next uint64,
) error {
	_, err := tx.Exec(`
		UPDATE checkpoints
		SET current_block = ?
		WHERE project_id = ?
	`, next, projectID)
	return err
}
