package shard

import (
	"fmt"
	"path/filepath"
	"time"
)

// BuildShardDir 构建 shard 根目录路径
// basePath: 项目 base path
// projectID: 项目 id
// month: YYYY-MM
func BuildShardDir(basePath, projectID, month string) string {
	return filepath.Join(basePath, projectID, month)
}

// BuildShardFile 构建 shard 文件路径
func BuildShardFile(basePath, projectID, month string, shardIndex int) string {
	dir := BuildShardDir(basePath, projectID, month)
	return filepath.Join(dir, fmt.Sprintf("shard_%03d.db", shardIndex))
}

// MonthFromTime 从 timestamp 获取 YYYY-MM
func MonthFromTime(t time.Time) string {
	return t.Format("2006-01")
}
