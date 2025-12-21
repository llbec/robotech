package service

import (
	"sort"
	"sync"
	"txstore/storage/sqlite"
)

type DBManager struct {
	mu      sync.RWMutex
	active  *sqlite.Store
	archive []*sqlite.Store
}

func NewDBManager(active *sqlite.Store) *DBManager {
	return &DBManager{active: active}
}

func (m *DBManager) Active() *sqlite.Store {
	return m.active
}

func (m *DBManager) All() []*sqlite.Store {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := append([]*sqlite.Store{}, m.archive...)
	res = append(res, m.active)
	return res
}

func (m *DBManager) AddArchive(store *sqlite.Store) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.archive = append(m.archive, store)
	sort.Slice(m.archive, func(i, j int) bool {
		return m.archive[i].Path < m.archive[j].Path
	})
}
