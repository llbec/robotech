package eventhandle

import (
	"math/big"
	"robotech/applications/marginv1/config"
	EventEmitter "robotech/contract/margin/v1/event"
	"robotech/utils/ethUtils"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
	contractEventEmitter *EventEmitter.EventEmitter
)

func init() {
	contract, err := EventEmitter.NewEventEmitter(config.GetEventEmitterAddr(), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	contractEventEmitter = contract
}

func EventProc(event types.Log) error {
	switch event.Topics[0] {
	case EventEmitter.PositionEvent:
		// TODO: handle Position event
		positionEvent, err := contractEventEmitter.ParsePosition(event)
		if err != nil {
			return err
		}
		return positionEventProc(positionEvent)
	default:
	}
	return nil
}

func positionEventProc(event *EventEmitter.EventEmitterPosition) error {
	switch event.ActionType {
	case big.NewInt(EventEmitter.DEPOSIT): // TODO: handle deposit event
	case big.NewInt(EventEmitter.WITHDRAW): // TODO: handle withdraw event
	default: //
	}
	return nil
}
