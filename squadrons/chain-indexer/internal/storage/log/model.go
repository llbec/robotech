package log

type TransactionLog struct {
	ProjectID string
	TxHash    string
	LogIndex  uint

	Address string
	Topics  string
	Data    []byte
}
