package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Position struct {
	MemeToken  common.Address `json:"memeToken"`
	Account    common.Address `json:"account"`
	PositionId *big.Int       `json:"positionId"`
	//BaseCollateral *big.Int       `json:"baseCollateral"`
	//BaseDebt       *big.Int       `json:"baseDebt"`
	//MemeCollateral *big.Int       `json:"memeCollateral"`
	//MemeDebt       *big.Int       `json:"MemeDebt"`
}

/*var (
	PositionCollector = map[common.Address]map[common.Address]Position{}
)*/
