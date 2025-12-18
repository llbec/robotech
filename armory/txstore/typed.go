package txstore

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type TxEvent struct {
	BlockHeight int64
	BlockTime   int64
	TxIndex     int64
	TxHash      string

	TxType string

	FromAddress string
	ToAddress   string

	Logs []*types.Log

	Day    string // YYYY-MM-DD
	Hour   int
	Minute int
	Second int
}
