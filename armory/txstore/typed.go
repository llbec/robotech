package txstore

type TxEvent struct {
	BlockHeight int64
	BlockTime   int64
	TxIndex     int64
	TxHash      string

	TxType string

	FromAddress string
	ToAddress   string

	RawTx    string
	ParsedTx string

	Day    string // YYYY-MM-DD
	Hour   int
	Minute int
	Second int
}
