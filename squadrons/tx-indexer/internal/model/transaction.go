package model

import (
	"math/big"
	"time"
)

type Transaction struct {
	TxHash []byte

	BlockNumber uint64

	BlockHash []byte

	BlockTime time.Time

	TxIndex int

	FromAddress []byte

	ToAddress []byte

	Value *big.Int

	Input []byte

	GasUsed uint64

	GasPrice *big.Int

	Status uint8

	IsContractCall bool

	MethodID []byte

	CreatedAt time.Time
}
