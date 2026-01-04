package txstore

import "ether-indexer/internal/model"

type Store interface {
	GetTransaction(projectID, hash string) (*model.Transaction, error)
}
