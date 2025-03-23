package liquidation

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Position struct {
	MemeToken      common.Address
	PositionId     *big.Int
	BaseCollateral *big.Int
	BaseDebt       *big.Int
	MemeCollateral *big.Int
	MemeDebt       *big.Int
}

var (
	PositionCollector = map[common.Address]map[common.Address]Position{}
)
