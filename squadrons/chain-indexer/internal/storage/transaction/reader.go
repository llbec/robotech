package transaction

import "database/sql"

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

// cursor: tx_hash + block_time 组合（由 query 层控制）
func (r *Reader) QueryByTime(
	projectID string,
	start, end int64,
	limit int,
) ([]*Transaction, error) {

	rows, err := r.db.Query(`
SELECT
	tx_hash, project_id, block_number, block_time, tx_index,
	from_address, to_address, value,
	gas, gas_price, nonce,
	input_data, status, created_at
FROM transactions
WHERE project_id = ?
  AND block_time >= ?
  AND block_time < ?
ORDER BY block_time, tx_index
LIMIT ?
`, projectID, start, end, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.TxHash, &t.ProjectID, &t.BlockNumber, &t.BlockTime, &t.TxIndex,
			&t.FromAddress, &t.ToAddress, &t.Value,
			&t.Gas, &t.GasPrice, &t.Nonce,
			&t.InputData, &t.Status, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, &t)
	}
	return result, nil
}
