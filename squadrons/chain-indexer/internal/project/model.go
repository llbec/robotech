package project

import "time"

type Project struct {
	ProjectID   string
	Active      bool
	RPCEndpoint string

	StartMonth string // YYYY-MM
	EndMonth   string // 可为空

	BasePath    string
	Step        int64
	IntervalSec int64

	Address []string
	Topics  []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Checkpoint struct {
	ProjectID    string
	StartBlock   int64
	CurrentBlock int64
	UpdatedAt    time.Time
}
