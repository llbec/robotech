package query

import (
	"fmt"
	"time"
)

func monthsBetween(start, end time.Time) []string {
	var res []string

	cur := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !cur.After(last) {
		res = append(res, fmt.Sprintf("%04d-%02d", cur.Year(), cur.Month()))
		cur = cur.AddDate(0, 1, 0)
	}
	return res
}
