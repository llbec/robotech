package liquidation

import (
	"math/big"
	EventEmitter "robotech/contract/margin/v1/event"

	"github.com/ethereum/go-ethereum/common"
)

func UpdateCollector(event *EventEmitter.EventEmitterPosition) {
	if event.BaseCollateral.Cmp(big.NewInt(0)) == 0 && event.MemeCollateral.Cmp(big.NewInt(0)) == 0 {
		if _, ok := PositionCollector[event.MemeToken]; ok {
			delete(PositionCollector[event.MemeToken], event.Account)
		}
	} else {
		if _, ok := PositionCollector[event.MemeToken]; !ok {
			PositionCollector[event.MemeToken] = make(map[common.Address]Position)
		}
		PositionCollector[event.MemeToken][event.Account] = Position{
			MemeToken:      event.MemeToken,
			PositionId:     event.PositionId,
			BaseCollateral: event.BaseCollateral,
			BaseDebt:       event.BaseDebtScaled,
			MemeCollateral: event.MemeCollateral,
			MemeDebt:       event.MemeDebtScaled,
		}
	}
}
