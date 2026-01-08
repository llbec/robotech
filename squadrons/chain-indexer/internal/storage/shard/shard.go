package shard

import (
	"database/sql"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Shard struct {
	FilePath string
	db       *sql.DB
	mu       sync.Mutex // 避免并发写入冲突
}

func NewShard(filePath string) (*Shard, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, err
	}

	s := &Shard{
		FilePath: filePath,
		db:       db,
	}

	if err := s.init(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// 初始化 shard 内表
func (s *Shard) init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sqlStmt := `
CREATE TABLE IF NOT EXISTS transactions (
	tx_hash TEXT PRIMARY KEY,
	block_number INTEGER,
	block_time DATETIME,
	from_address TEXT,
	to_address TEXT,
	value TEXT,
	log_index INTEGER
);`

	_, err := s.db.Exec(sqlStmt)
	return err
}

// 写入单条 transaction
func (s *Shard) InsertTransaction(txHash string, blockNumber int64, blockTime time.Time,
	from, to, value string, logIndex int) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
INSERT OR IGNORE INTO transactions
(tx_hash, block_number, block_time, from_address, to_address, value, log_index)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		txHash, blockNumber, blockTime, from, to, value, logIndex,
	)
	return err
}

// 查询示例（分页）
// limit, offset 分页
func (s *Shard) QueryTransactions(startTime, endTime time.Time, limit, offset int) (*sql.Rows, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Query(`
SELECT tx_hash, block_number, block_time, from_address, to_address, value, log_index
FROM transactions
WHERE block_time BETWEEN ? AND ?
ORDER BY block_time ASC
LIMIT ? OFFSET ?`,
		startTime, endTime, limit, offset,
	)
}

// 关闭 shard
func (s *Shard) Close() error {
	return s.db.Close()
}
