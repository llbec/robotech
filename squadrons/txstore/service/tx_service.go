package service

import (
	"sort"
	"txstore/model"
)

type TxService struct {
	db *DBManager
}

func NewTxService(db *DBManager) *TxService {
	return &TxService{db: db}
}

func (s *TxService) Insert(project string, txs []*model.Transaction) error {
	return s.db.Active().Insert(project, txs)
}

func (s *TxService) Scan(project string, block, index int64, limit int) ([]*model.Transaction, error) {
	var all []*model.Transaction
	for _, st := range s.db.All() {
		txs, _ := st.Scan(project, block, index, limit)
		all = append(all, txs...)
	}
	sortTxs(all)
	if len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func sortTxs(txs []*model.Transaction) {
	sort.Slice(txs, func(i, j int) bool {
		if txs[i].BlockHeight == txs[j].BlockHeight {
			return txs[i].TxIndex < txs[j].TxIndex
		}
		return txs[i].BlockHeight < txs[j].BlockHeight
	})
}
