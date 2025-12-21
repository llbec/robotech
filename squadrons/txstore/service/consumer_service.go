package service

import "txstore/model"

type ConsumerService struct {
	tx *TxService
	cp *DBManager
}

func NewConsumerService(tx *TxService) *ConsumerService {
	return &ConsumerService{tx: tx}
}

func (s *ConsumerService) Scan(project, consumer string, limit int) ([]*model.Transaction, *model.Checkpoint, error) {
	cp, _ := s.tx.db.Active().Get(project, consumer)
	if cp == nil {
		cp = &model.Checkpoint{ProjectID: project, Consumer: consumer}
	}
	txs, err := s.tx.Scan(project, cp.BlockHeight, cp.TxIndex, limit)
	return txs, cp, err
}

func (s *ConsumerService) Commit(cp *model.Checkpoint) error {
	return s.tx.db.Active().Save(cp)
}
