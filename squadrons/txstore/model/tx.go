package model

type Transaction struct {
	ProjectID string `json:"project_id"`

	TxHash      string `json:"tx_hash"`
	BlockHeight int64  `json:"block_height"`
	TxIndex     int64  `json:"tx_index"`
	BlockTime   int64  `json:"block_time"`

	From   string `json:"from"`
	To     string `json:"to"`
	Value  string `json:"value"`
	Input  []byte `json:"input"`
	Status int    `json:"status"`

	TxType   string `json:"tx_type"`
	TxAction string `json:"tx_action"`

	Logs []Log `json:"logs"`
}
