package model

import "time"

type Transaction struct {
	ProjectID   string
	TxHash      string
	BlockNumber uint64
	BlockTime   time.Time
	TxIndex     uint

	FromAddress string
	ToAddress   string
	Value       string

	Gas      uint64
	GasPrice string
	Nonce    uint64

	InputData []byte
	Status    int
	CreatedAt time.Time
}
