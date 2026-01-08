package shard

import (
	"sync"
)

// ShardManager 管理多个项目 shardRouter
type ShardManager struct {
	mu      sync.Mutex
	routers map[string]*ShardRouter // projectID -> router
}

func NewShardManager() *ShardManager {
	return &ShardManager{
		routers: make(map[string]*ShardRouter),
	}
}

// 获取或创建项目 router
func (m *ShardManager) GetRouter(projectID, basePath string) *ShardRouter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if router, ok := m.routers[projectID]; ok {
		return router
	}

	router := NewRouter(basePath, projectID)
	m.routers[projectID] = router
	return router
}

// 关闭某个项目 shard
func (m *ShardManager) CloseRouter(projectID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if router, ok := m.routers[projectID]; ok {
		_ = router.Close()
		delete(m.routers, projectID)
	}
}

// 全部关闭
func (m *ShardManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for pid, router := range m.routers {
		_ = router.Close()
		delete(m.routers, pid)
	}
}
