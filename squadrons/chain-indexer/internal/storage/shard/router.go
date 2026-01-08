package shard

import (
	"fmt"
	"sync"
	"time"
)

// ShardRouter 管理同一项目下按月份 shard 的打开 / 缓存
type ShardRouter struct {
	basePath  string
	projectID string

	shards map[string][]*Shard // month -> shard 列表
	mu     sync.Mutex
}

// NewRouter 创建 router
func NewRouter(basePath, projectID string) *ShardRouter {
	return &ShardRouter{
		basePath:  basePath,
		projectID: projectID,
		shards:    make(map[string][]*Shard),
	}
}

// GetShardForTime 返回时间对应的 shard（这里简单返回第 0 个 shard，可后续按容量扩展）
func (r *ShardRouter) GetShardForTime(t time.Time) (*Shard, error) {
	month := MonthFromTime(t)

	r.mu.Lock()
	defer r.mu.Unlock()

	// 已打开
	if shards, ok := r.shards[month]; ok && len(shards) > 0 {
		return shards[0], nil
	}

	// 新建 shard
	filePath := BuildShardFile(r.basePath, r.projectID, month, 0)
	shard, err := NewShard(filePath)
	if err != nil {
		return nil, fmt.Errorf("create shard failed: %w", err)
	}

	r.shards[month] = []*Shard{shard}
	return shard, nil
}

// Close 关闭所有 shard
func (r *ShardRouter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, shards := range r.shards {
		for _, s := range shards {
			_ = s.Close()
		}
	}
	return nil
}
