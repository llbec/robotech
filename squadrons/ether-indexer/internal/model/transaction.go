package model

type Transaction struct {
	ProjectID   string
	TxHash      string
	BlockNumber uint64
	BlockTime   int64
	TxIndex     uint

	FromAddress string
	ToAddress   string
	Value       string

	Gas      uint64
	GasPrice string
	Nonce    uint64

	InputData []byte
	Status    int
	CreatedAt int64
}
