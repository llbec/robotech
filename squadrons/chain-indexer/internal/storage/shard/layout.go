package shard

import (
	"fmt"
	"time"
)

func MonthKey(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	return fmt.Sprintf("%04d-%02d", t.Year(), int(t.Month()))
}
