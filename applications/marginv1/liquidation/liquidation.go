package liquidation

import (
	"fmt"
	"log"
	"math/big"
	"robotech/applications/marginv1/config"
	"robotech/applications/marginv1/storage"
	"robotech/applications/marginv1/types"
	v1 "robotech/contract/margin/v1"
	EventEmitter "robotech/contract/margin/v1/event"
	Reader "robotech/contract/margin/v1/reader"
	Router "robotech/contract/margin/v1/router"
	"robotech/utils/ethUtils"

	"github.com/ethereum/go-ethereum/common"
)

func UpdateCollector(event *EventEmitter.EventEmitterPosition) {
	if event.BaseCollateral.Cmp(big.NewInt(0)) == 0 && event.MemeCollateral.Cmp(big.NewInt(0)) == 0 {
		/*if _, ok := PositionCollector[event.MemeToken]; ok {
			delete(PositionCollector[event.MemeToken], event.Account)
		}*/
		storage.RemovePosition(event.MemeToken, event.Account)
	} else {
		/*if _, ok := PositionCollector[event.MemeToken]; !ok {
			PositionCollector[event.MemeToken] = make(map[common.Address]Position)
		}
		PositionCollector[event.MemeToken][event.Account] = Position{
			//MemeToken:      event.MemeToken,
			PositionId: event.PositionId,
			//BaseCollateral: event.BaseCollateral,
			//BaseDebt:       event.BaseDebtScaled,
			//MemeCollateral: event.MemeCollateral,
			//MemeDebt:       event.MemeDebtScaled,
		}*/
		storage.AddPosition(event.MemeToken, event.Account, types.Position{
			PositionId: event.PositionId,
		})
	}
}

func Liquidation(pools []common.Address) {
	for _, pool := range pools {
		err := liquidation(pool)
		if err != nil {
			log.Printf("Liquidation pool(%v) failed: %v", pool, err)
		}
	}
}

func liquidation(mToken common.Address) error {
	storeList := storage.GetPositions(mToken)

	log.Printf("Liquidation pool(%v) total count: %v\n%v", mToken, len(storeList), storeList)

	var positionKeys [][32]byte
	for _, position := range storeList {
		positionKey := v1.CalcPositionKey(position.Account, position.PositionId)
		var hash [32]byte
		copy(hash[:], positionKey)
		positionKeys = append(positionKeys, hash)
	}

	positions, err := CallGetPositions2(positionKeys)
	if err != nil {
		return fmt.Errorf("CallGetPositions2 failed: %v\n%v", err, positionKeys)
	}

	var params []Router.LiquidationUtilsLiquidationParams
	for _, position := range positions {
		log.Printf("Liquidation pool(%v) position(%v-%v) margin level: %v, threshold: %v\n", mToken, position.Account, position.Id, position.MarginLevel, MarginLevelThreshold())
		if position.MarginLevel.Cmp(MarginLevelThreshold()) < 1 {
			params = append(params, Router.LiquidationUtilsLiquidationParams{
				Account:    position.Account,
				PositionId: position.Id,
			})
		}
	}

	if len(params) == 0 {
		return nil
	}

	log.Printf("Liquidation pool(%v) liquidation position count: %v\n", mToken, len(params))

	return CallLiquidation(params)
}

func CallGetPositions(account common.Address) ([]Reader.ReaderPositionUtilsGetPosition, error) {
	readerContract, err := Reader.NewReader(config.GetReaderAddr(), ethUtils.GetClient())
	if err != nil {
		return []Reader.ReaderPositionUtilsGetPosition{}, fmt.Errorf("NewReader failed: %v", err)
	}

	return readerContract.GetPositions(nil, config.GetDataStoreAddr(), account)
}

func CallGetPositions2(positionKeys [][32]byte) ([]Reader.ReaderPositionUtilsGetPosition, error) {
	readerContract, err := Reader.NewReader(config.GetReaderAddr(), ethUtils.GetClient())
	if err != nil {
		return []Reader.ReaderPositionUtilsGetPosition{}, fmt.Errorf("NewReader failed: %v", err)
	}

	return readerContract.GetPositions2(nil, config.GetDataStoreAddr(), positionKeys)
}

func CallLiquidation(params []Router.LiquidationUtilsLiquidationParams) error {
	log.Printf("Liquidation params: %v\n", (params))

	exchangeRouterContract, err := Router.NewRouter(config.GetExchangeRouterAddr(), ethUtils.GetClient())
	if err != nil {
		return fmt.Errorf("NewRouter failed: %v", err)
	}

	opt, err := ethUtils.GetOpt(config.GetSKey())
	if err != nil {
		return fmt.Errorf("GetOpt failed: %v", err)
	}

	tx, err := exchangeRouterContract.ExecuteLiquidationBatch(opt, params)
	if err != nil {
		return fmt.Errorf("ExecuteLiquidationBatch failed: %v", err)
	}
	log.Printf("Liquidation tx: %s\n", tx.Hash().Hex())

	return nil
}

func MarginLevelThreshold() *big.Int {
	threshold, _ := big.NewInt(0).SetString("1100000000000000000000000000", 10)
	return threshold
}
