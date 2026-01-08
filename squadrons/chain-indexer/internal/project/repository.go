package project

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Init() error {
	projectSQL := `
CREATE TABLE IF NOT EXISTS projects (
	project_id TEXT PRIMARY KEY,
	active BOOLEAN,
	rpc_endpoint TEXT,
	start_month TEXT,
	end_month TEXT,
	base_path TEXT,
	step INTEGER,
	interval_sec INTEGER,
	address TEXT,
	topics TEXT,
	created_at DATETIME,
	updated_at DATETIME
);`

	checkpointSQL := `
CREATE TABLE IF NOT EXISTS checkpoint (
	project_id TEXT PRIMARY KEY,
	start_block INTEGER,
	current_block INTEGER,
	updated_at DATETIME
);`

	if _, err := r.db.Exec(projectSQL); err != nil {
		return err
	}
	if _, err := r.db.Exec(checkpointSQL); err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateProject(p *Project, startBlock int64) error {
	if err := ValidateProjectCreate(p); err != nil {
		return err
	}

	addr, _ := json.Marshal(p.Address)
	topics, _ := json.Marshal(p.Topics)

	now := time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
INSERT INTO projects (
	project_id, active, rpc_endpoint,
	start_month, end_month,
	base_path, step, interval_sec,
	address, topics,
	created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ProjectID, p.Active, p.RPCEndpoint,
		p.StartMonth, p.EndMonth,
		p.BasePath, p.Step, p.IntervalSec,
		string(addr), string(topics),
		now, now,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO checkpoint (
	project_id, start_block, current_block, updated_at
) VALUES (?, ?, ?, ?)`,
		p.ProjectID, startBlock, startBlock, now,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) ListProjects() ([]*Project, error) {
	rows, err := r.db.Query(`SELECT
	project_id, active, rpc_endpoint,
	start_month, end_month,
	base_path, step, interval_sec,
	address, topics,
	created_at, updated_at
	FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*Project
	for rows.Next() {
		var p Project
		var addr, topics string

		if err := rows.Scan(
			&p.ProjectID, &p.Active, &p.RPCEndpoint,
			&p.StartMonth, &p.EndMonth,
			&p.BasePath, &p.Step, &p.IntervalSec,
			&addr, &topics,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(addr), &p.Address)
		_ = json.Unmarshal([]byte(topics), &p.Topics)

		res = append(res, &p)
	}
	return res, nil
}

func (r *Repository) GetProject(id string) (*Project, error) {
	row := r.db.QueryRow(`
SELECT project_id, active, rpc_endpoint,
	start_month, end_month,
	base_path, step, interval_sec,
	address, topics,
	created_at, updated_at
FROM projects WHERE project_id = ?`, id)

	var p Project
	var addr, topics string

	if err := row.Scan(
		&p.ProjectID, &p.Active, &p.RPCEndpoint,
		&p.StartMonth, &p.EndMonth,
		&p.BasePath, &p.Step, &p.IntervalSec,
		&addr, &topics,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(addr), &p.Address)
	_ = json.Unmarshal([]byte(topics), &p.Topics)

	return &p, nil
}

func (r *Repository) UpdateProject(p *Project) error {
	old, err := r.GetProject(p.ProjectID)
	if err != nil {
		return err
	}

	if err = ValidateProjectUpdate(old, p); err != nil {
		return err
	}

	addr, _ := json.Marshal(p.Address)
	topics, _ := json.Marshal(p.Topics)

	_, err = r.db.Exec(`
UPDATE projects SET
	active = ?,
	rpc_endpoint = ?,
	base_path = ?,
	step = ?,
	interval_sec = ?,
	address = ?,
	topics = ?,
	updated_at = ?
WHERE project_id = ?`,
		p.Active, p.RPCEndpoint, p.BasePath,
		p.Step, p.IntervalSec,
		string(addr), string(topics),
		time.Now(), p.ProjectID,
	)
	return err
}

func (r *Repository) GetCheckpoint(projectID string) (*Checkpoint, error) {
	row := r.db.QueryRow(`
SELECT project_id, start_block, current_block, updated_at
FROM checkpoint WHERE project_id = ?`, projectID)

	var c Checkpoint
	if err := row.Scan(
		&c.ProjectID,
		&c.StartBlock,
		&c.CurrentBlock,
		&c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) SaveCheckpoint(c *Checkpoint) error {
	if c.ProjectID == "" {
		return errors.New("project_id required")
	}

	_, err := r.db.Exec(`
UPDATE checkpoint SET
	current_block = ?,
	updated_at = ?
WHERE project_id = ?`,
		c.CurrentBlock,
		time.Now(),
		c.ProjectID,
	)
	return err
}
