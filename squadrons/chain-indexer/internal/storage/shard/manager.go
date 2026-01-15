package shard

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

type Manager struct {
	mu    sync.Mutex
	cache map[string]*sql.DB
}

func NewManager() *Manager {
	return &Manager{
		cache: make(map[string]*sql.DB),
	}
}

func (m *Manager) OpenShard(path string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, ok := m.cache[path]; ok {
		return db, nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(path, "shard.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	m.cache[path] = db
	return db, nil
}
