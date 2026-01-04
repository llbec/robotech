package model

type TransactionLog struct {
	ProjectID string
	TxHash    string
	LogIndex  uint

	Address string
	Topics  string // JSON
	Data    []byte
}
