package compliquidation

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Asset struct {
	Symbol  string
	Address common.Address
	Balance *big.Float
	Value   *big.Float
}

type Account struct {
	Address     string
	Collaterals []*Asset
	Debts       []*Asset
	Collateral  *big.Float
	Debt        *big.Float
	Health      *big.Float
}

func (a *Account) String() string {
	bs, _ := json.Marshal(a)
	return string(bs)
}
