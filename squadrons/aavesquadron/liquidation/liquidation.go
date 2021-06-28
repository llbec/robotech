package liquidation

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/llbec/robotech/logistics/db"
)

type AAVELiquidation struct {
	HttpClient           *ethclient.Client
	WSClient             *ethclient.Client
	AccountDB            db.DB
	LiquidationDB        db.DB
	LendingPool          *aave.LendingPool
	ProtocaldataProvider *aave.AaveProtocolDataProvider
	AaveOracle           *aave.AaveOracle
	FlashLiquidationAdp  *flashloan.AaveLiquidationAdapter
}

func (al *AAVELiquidation) LoadDB() {
	if al.AccountDB == nil || al.LiquidationDB == nil {
		panic("DB is null")
	}
}
