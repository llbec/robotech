package storage

import (
	"database/sql"
	"os"
	"sync"
)

type ShardManager struct {
	mu      sync.Mutex
	policy  ShardPolicy
	current *sql.DB
	start   uint64
}

func NewShardManager(p ShardPolicy) *ShardManager {
	return &ShardManager{policy: p}
}

func (m *ShardManager) GetDB(block uint64) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	start := (block / m.policy.BlockRange) * m.policy.BlockRange
	if m.current != nil && start == m.start {
		return m.current, nil
	}

	if m.current != nil {
		_ = m.current.Close()
	}

	if err := os.MkdirAll(m.policy.BasePath, 0755); err != nil {
		return nil, err
	}

	db, err := OpenDB(m.policy.DBPath(start))
	if err != nil {
		return nil, err
	}
	if err := Migrate(db); err != nil {
		return nil, err
	}

	m.current = db
	m.start = start
	return db, nil
}
