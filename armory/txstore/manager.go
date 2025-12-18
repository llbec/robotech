package txstore

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

type StoreManager struct {
	dir    string
	mu     sync.Mutex
	stores map[string]*SQLiteStore
}

func NewStoreManager(dir string) *StoreManager {
	return &StoreManager{
		dir:    dir,
		stores: make(map[string]*SQLiteStore),
	}
}

func (m *StoreManager) getMonthKey(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func (m *StoreManager) GetStoreByTime(ts int64) (*SQLiteStore, error) {
	key := m.getMonthKey(ts)

	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.stores[key]; ok {
		return s, nil
	}

	path := filepath.Join(m.dir, fmt.Sprintf("tx-%s.db", key))
	store, err := OpenSQLite(path)
	if err != nil {
		return nil, err
	}

	m.stores[key] = store
	return store, nil
}
