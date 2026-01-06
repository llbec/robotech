package storage

import (
	"database/sql"
	"ether-indexer/internal/model"
	"time"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(id, rpc, desc string) error {
	_, err := r.db.Exec(`
		INSERT INTO projects
		(project_id, active, rpc_endpoint, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		id,
		true,
		rpc,
		desc,
		time.Now().Unix(),
		time.Now().Unix(),
	)
	return err
}

func (r *ProjectRepo) ListAll() ([]model.Project, error) {
	rows, err := r.db.Query(`
		SELECT project_id, active, rpc_endpoint, description, created_at, updated_at
		FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Project
	for rows.Next() {
		var p model.Project
		rows.Scan(
			&p.ProjectID,
			&p.Active,
			&p.RPCEndpoint,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		res = append(res, p)
	}
	return res, nil
}

func (r *ProjectRepo) UpdateActive(id string, active bool) error {
	_, err := r.db.Exec(`
		UPDATE projects
		SET active = ?, updated_at = ?
		WHERE project_id = ?`,
		active, time.Now().Unix(), id)
	return err
}

func (r *ProjectRepo) UpdateRPC(id, rpc string) error {
	_, err := r.db.Exec(`
		UPDATE projects
		SET rpc_endpoint = ?, updated_at = ?
		WHERE project_id = ?`,
		rpc, time.Now().Unix(), id)
	return err
}

func (r *ProjectRepo) UpdateDescription(id, desc string) error {
	_, err := r.db.Exec(`
		UPDATE projects
		SET description = ?, updated_at = ?
		WHERE project_id = ?`,
		desc, time.Now().Unix(), id)
	return err
}
