package model

type Log struct {
	TxHash   string   `json:"tx_hash"`
	LogIndex int64    `json:"log_index"`
	Contract string   `json:"contract"`
	Topics   []string `json:"topics"`
	Data     []byte   `json:"data"`
}
