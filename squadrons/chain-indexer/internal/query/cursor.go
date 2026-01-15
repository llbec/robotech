package query

import (
	"encoding/base64"
	"encoding/json"
)

type Cursor struct {
	Month     string `json:"month"`
	Offset    int    `json:"offset"`
	BlockTime int64  `json:"block_time"`
	TxIndex   uint   `json:"tx_index"`
}

func EncodeCursor(c *Cursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeCursor(s string) (*Cursor, error) {
	if s == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
