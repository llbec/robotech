package storage

import "fmt"

type ShardPolicy struct {
	BlockRange uint64
	BasePath   string
}

func (p ShardPolicy) DBPath(start uint64) string {
	end := start + p.BlockRange
	return fmt.Sprintf("%s/tx_%d_%d.db", p.BasePath, start, end)
}
