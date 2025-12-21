package model

type Checkpoint struct {
	ProjectID   string `json:"project_id"`
	Consumer    string `json:"consumer"`
	BlockHeight int64  `json:"block_height"`
	TxIndex     int64  `json:"tx_index"`
	UpdatedAt   int64  `json:"updated_at"`
}
