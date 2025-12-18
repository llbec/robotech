package txstore

import (
	"context"
	"time"
)

func (s *SQLiteStore) ListTxEventsAfter(
	ctx context.Context,
	blockHeight int64,
	limit int,
) ([]*TxEvent, error) {

	rows, err := s.db.QueryContext(ctx, `
	SELECT block_height, block_time, tx_index, tx_hash,
	       tx_type, from_address, to_address,
	       day, hour, minute, second
	FROM tx_events
	WHERE block_height > ?
	ORDER BY block_height ASC, tx_index ASC
	LIMIT ?`,
		blockHeight, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*TxEvent
	for rows.Next() {
		ev := new(TxEvent)
		if err := rows.Scan(
			&ev.BlockHeight,
			&ev.BlockTime,
			&ev.TxIndex,
			&ev.TxHash,
			&ev.TxType,
			&ev.FromAddress,
			&ev.ToAddress,
			&ev.Day,
			&ev.Hour,
			&ev.Minute,
			&ev.Second,
		); err != nil {
			return nil, err
		}
		list = append(list, ev)
	}
	return list, nil
}

func (s *SQLiteStore) QueryTxByDayAndType(
	ctx context.Context,
	day string,
	txType string,
) ([]*TxEvent, error) {

	rows, err := s.db.QueryContext(ctx, `
	SELECT block_height, block_time, tx_index, tx_hash,
	       tx_type, from_address, to_address,
	       day, hour, minute, second
	FROM tx_events
	WHERE day = ? AND tx_type = ?
	ORDER BY block_time ASC`,
		day, txType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*TxEvent
	for rows.Next() {
		ev := new(TxEvent)
		if err := rows.Scan(
			&ev.BlockHeight,
			&ev.BlockTime,
			&ev.TxIndex,
			&ev.TxHash,
			&ev.TxType,
			&ev.FromAddress,
			&ev.ToAddress,
			&ev.Day,
			&ev.Hour,
			&ev.Minute,
			&ev.Second,
		); err != nil {
			return nil, err
		}
		list = append(list, ev)
	}
	return list, nil
}

func (m *StoreManager) ScanAfter(
	ctx context.Context,
	fromTime int64,
	limit int,
) ([]*TxEvent, error) {

	var result []*TxEvent
	curTime := fromTime

	for len(result) < limit {
		store, err := m.GetStoreByTime(curTime)
		if err != nil {
			return nil, err
		}

		txs, _ := store.ListTxEventsAfter(ctx, 0, limit-len(result))
		if len(txs) == 0 {
			// 跳到下个月
			curTime = nextMonth(curTime)
			continue
		}

		result = append(result, txs...)
		curTime = txs[len(txs)-1].BlockTime
	}

	return result, nil
}

func nextMonth(ts int64) int64 {
	t := time.Unix(ts, 0).UTC()
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC).Unix()
}
