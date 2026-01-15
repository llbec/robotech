package query

import (
	"time"
)

type QueryRequest struct {
	ProjectID string

	StartTime time.Time
	EndTime   time.Time

	PageSize int
	Cursor   string
}

type QueryResult[T any] struct {
	Data       []T
	NextCursor string
}
