package sqlite

import (
	"encoding/json"
	"fmt"
	"txstore/model"
)

// SplitArchiveMove 多 project 支持，并在迁移完成后删除源库区块
func SplitArchiveMove(src, dst *Store, fromBlock, toBlock int64) error {
	tx, err := src.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 获取所有 project
	projRows, err := tx.Query(`SELECT DISTINCT project_id FROM transactions`)
	if err != nil {
		return err
	}
	var projects []string
	for projRows.Next() {
		var p string
		projRows.Scan(&p)
		projects = append(projects, p)
	}
	projRows.Close()

	totalTx := 0
	for _, project := range projects {
		// 查询交易
		rows, err := tx.Query(`
            SELECT tx_hash, block_height, tx_index, block_time,
                   sender, receiver, value, input, status, tx_type, tx_action
            FROM transactions
            WHERE project_id=? AND block_height BETWEEN ? AND ?
            ORDER BY block_height, tx_index
        `, project, fromBlock, toBlock)
		if err != nil {
			return err
		}

		var txs []*model.Transaction
		for rows.Next() {
			t := &model.Transaction{ProjectID: project}
			if err_1 := rows.Scan(
				&t.TxHash, &t.BlockHeight, &t.TxIndex, &t.BlockTime,
				&t.From, &t.To, &t.Value, &t.Input, &t.Status,
				&t.TxType, &t.TxAction,
			); err_1 != nil {
				rows.Close()
				return err_1
			}
			// 读取日志
			logRows, err_2 := tx.Query(`
                SELECT log_index, contract, topics, data
                FROM logs WHERE project_id=? AND tx_hash=? ORDER BY log_index
            `, project, t.TxHash)
			if err_2 != nil {
				rows.Close()
				return err_2
			}
			for logRows.Next() {
				var l model.Log
				var topicsStr string
				logRows.Scan(&l.LogIndex, &l.Contract, &topicsStr, &l.Data)
				l.Topics = []string{}
				json.Unmarshal([]byte(topicsStr), &l.Topics)
				t.Logs = append(t.Logs, l)
			}
			logRows.Close()
			txs = append(txs, t)
		}
		rows.Close()

		if len(txs) == 0 {
			continue
		}

		// 插入 archive
		if err_3 := dst.Insert(project, txs); err_3 != nil {
			return err_3
		}

		// 删除源库中的交易和日志
		hashes := ""
		for i, t := range txs {
			if i > 0 {
				hashes += ","
			}
			hashes += fmt.Sprintf("'%s'", t.TxHash)
		}
		_, err = tx.Exec(fmt.Sprintf(`DELETE FROM logs WHERE project_id=? AND tx_hash IN (%s)`, hashes), project)
		if err != nil {
			return err
		}
		_, err = tx.Exec(fmt.Sprintf(`DELETE FROM transactions WHERE project_id=? AND tx_hash IN (%s)`, hashes), project)
		if err != nil {
			return err
		}

		totalTx += len(txs)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("迁移完成: 区块 %d-%d, 交易总数 %d, project数 %d\n", fromBlock, toBlock, totalTx, len(projects))
	return nil
}
