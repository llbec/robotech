package storage

import "txstore/model"

type TxStore interface {
	Insert(project string, txs []*model.Transaction) error
	Scan(project string, fromBlock, fromIndex int64, limit int) ([]*model.Transaction, error)
}

type CheckpointStore interface {
	Get(project, consumer string) (*model.Checkpoint, error)
	Save(cp *model.Checkpoint) error
}
