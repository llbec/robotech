package model

import "time"

type Project struct {
	ProjectID   string
	Name        string
	ChainID     int64
	RPCEndpoint string
	CreatedAt   time.Time
}
