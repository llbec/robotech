package model

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type Checkpoint struct {
	ProjectID     string
	StartBlock    uint64
	CurrentBlock  uint64
	Step          uint64
	IntervalSec   int
	AddressesJSON string
	TopicsJSON    string
}

func (c *Checkpoint) Addresses() []string {
	var a []string
	_ = json.Unmarshal([]byte(c.AddressesJSON), &a)
	return a
}

func (c *Checkpoint) Topics() [][]common.Hash {
	var raw [][]string
	_ = json.Unmarshal([]byte(c.TopicsJSON), &raw)

	var out [][]common.Hash
	for _, layer := range raw {
		var l []common.Hash
		for _, v := range layer {
			l = append(l, common.HexToHash(v))
		}
		out = append(out, l)
	}
	return out
}
