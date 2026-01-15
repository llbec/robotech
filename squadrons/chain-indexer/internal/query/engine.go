package query

import (
	"time"

	"chain-indexer/internal/storage/shard"
	"chain-indexer/internal/storage/transaction"
)

type Engine struct {
	router  *shard.Router
	manager *shard.Manager
}

func NewEngine(router *shard.Router, manager *shard.Manager) *Engine {
	return &Engine{
		router:  router,
		manager: manager,
	}
}

func (e *Engine) QueryTransactions(
	req QueryRequest,
) (*QueryResult[*transaction.Transaction], error) {

	if req.PageSize <= 0 {
		req.PageSize = 50
	}

	cur, err := DecodeCursor(req.Cursor)
	if err != nil {
		return nil, err
	}

	months := monthsBetween(req.StartTime, req.EndTime)
	if len(months) == 0 {
		return &QueryResult[*transaction.Transaction]{}, nil
	}

	var (
		result      []*transaction.Transaction
		nextCursor  *Cursor
		remaining   = req.PageSize
		startMonth  = 0
		startOffset = 0
	)

	if cur != nil {
		for i, m := range months {
			if m == cur.Month {
				startMonth = i
				startOffset = cur.Offset
				break
			}
		}
	}

	for mi := startMonth; mi < len(months); mi++ {
		month := months[mi]

		path := e.router.Route(req.ProjectID,
			time.Date(
				parseMonth(month).Year(),
				parseMonth(month).Month(),
				1, 0, 0, 0, 0, time.UTC,
			).Unix(),
		).Path

		db, err := e.manager.OpenShard(path)
		if err != nil {
			continue
		}

		reader := transaction.NewReader(db)

		rows, err := reader.QueryByTime(
			req.ProjectID,
			req.StartTime.Unix(),
			req.EndTime.Unix(),
			remaining+startOffset,
		)
		if err != nil {
			return nil, err
		}

		if startOffset > 0 && len(rows) > startOffset {
			rows = rows[startOffset:]
		}

		for i, tx := range rows {
			if remaining == 0 {
				nextCursor = &Cursor{
					Month:     month,
					Offset:    startOffset + i,
					BlockTime: tx.BlockTime,
					TxIndex:   tx.TxIndex,
				}
				break
			}
			result = append(result, tx)
			remaining--
		}

		if remaining == 0 {
			break
		}

		startOffset = 0
	}

	var cursorStr string
	if nextCursor != nil {
		cursorStr, err = EncodeCursor(nextCursor)
		if err != nil {
			return nil, err
		}
	}

	return &QueryResult[*transaction.Transaction]{
		Data:       result,
		NextCursor: cursorStr,
	}, nil
}

func parseMonth(m string) time.Time {
	t, _ := time.Parse("2006-01", m)
	return t
}
