package eventhandle

import (
	"context"
	"log"
	"math/big"
	"robotech/applications/marginv1/config"
	"robotech/applications/marginv1/liquidation"
	EventEmitter "robotech/contract/margin/v1/event"
	"robotech/utils/ethUtils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	contractEventEmitter *EventEmitter.EventEmitter
	check_pool_map       = map[common.Address]bool{}
)

const (
	maxEnd = 100
)

func init() {
	contract, err := EventEmitter.NewEventEmitter(config.GetEventEmitterAddr(), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	contractEventEmitter = contract
	log.Printf("EventEmitter initialized")
}

func SyncBlock(from, to int64) (int64, []common.Address) {
	//log.Printf("block sync task from %v to %v, step %v", from, to, maxEnd)
	start := from
	end := start + maxEnd
	for {
		if end > to {
			end = to
		}
		//log.Printf("sync block from %v to %v\n", start, end)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{config.GetEventEmitterAddr()},
			FromBlock: big.NewInt(start),
			ToBlock:   big.NewInt(end),
		}
		logs, err := ethUtils.GetClient().FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("FilterLogs(%v-%v): %v\n", start, end, err)
			return start, []common.Address{}
		}
		for _, vLog := range logs {
			err = eventProc(vLog)
			if err != nil {
				log.Printf("HandleLog(%v-%v): %v\n", vLog.BlockNumber, vLog.Index, err)
			}
		}
		if end == to {
			break
		}
		start = end + 1
		end = start + maxEnd
	}

	var checkList []common.Address
	for k, v := range check_pool_map {
		if v {
			checkList = append(checkList, k)
		}
		delete(check_pool_map, k)
	}

	return end, checkList
}

func eventProc(event types.Log) error {
	switch event.Topics[0] {
	case EventEmitter.PositionEvent:
		// TODO: handle Position event
		positionEvent, err := contractEventEmitter.ParsePosition(event)
		if err != nil {
			return err
		}
		return positionEventProc(positionEvent)
	case EventEmitter.SwapEvent:
		swapEvent, err := contractEventEmitter.ParseSwap(event)
		if err != nil {
			return err
		}

		memeToken := swapEvent.TokenIn
		if swapEvent.TokenIn == config.GetBaseTokenAddr() {
			memeToken = (swapEvent.TokenOut)
		}
		//liquidation.Liquidation(memeToken)
		check_pool_map[memeToken] = true
	case EventEmitter.RemoveEvent:
		removeEvent, err := contractEventEmitter.ParseRemove(event)
		if err != nil {
			return err
		}
		//liquidation.Liquidation(removeEvent.MemeToken)
		check_pool_map[removeEvent.MemeToken] = true
	default:
	}
	return nil
}

func positionEventProc(event *EventEmitter.EventEmitterPosition) error {
	/*switch event.ActionType {
	case big.NewInt(EventEmitter.DEPOSIT): // TODO: handle deposit event
	case big.NewInt(EventEmitter.WITHDRAW): // TODO: handle withdraw event
	default: //
	}*/
	liquidation.UpdateCollector(event)
	if event.ActionType.Cmp(big.NewInt(EventEmitter.SWAP)) == 0 {
		//liquidation.Liquidation(event.MemeToken)
		check_pool_map[event.MemeToken] = true
	}
	return nil
}
