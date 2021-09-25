package lendingpoolevent

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/llbec/robotech/armory/aave"
	"github.com/llbec/robotech/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type parseFunc func(eventLog types.Log) (interface{}, error)

var (
	contractAbi = abi.ABI{}
	handleList  = make(map[common.Hash]EventHandle)
	parseList   = make(map[common.Hash]parseFunc)
)

func init() {
	cabi, err := abi.JSON(strings.NewReader(string(aave.LendingPoolABI)))
	if err != nil {
		log.Fatalf("Parsing contract abi error: %v\n", err)
	}
	contractAbi = cabi
	parseInit()
}

func parseInit() {
	parseList[BorrowEvent] = borrowEventParse
	parseList[DepositEvent] = depositEventParse
	parseList[FlashLoanEvent] = flashLoanEventParse
	parseList[LiquidationEvent] = liquidationEventParse
	parseList[PausedEvent] = pausedEventParse
	parseList[RebalanceStableRateEvent] = rebalanceStableRateEventParse
	parseList[RepayEvent] = repayEventParse
	parseList[ReserveUpdatedEvent] = reserveUpdatedEventParse
	parseList[CollateralDisableEvent] = collateralDisableEventParse
	parseList[CollateralEnableEvent] = collateralEnableEventParse
	parseList[SwapEvent] = swapEventParse
	parseList[UnpausedEvent] = unpausedEventParse
	parseList[WithdrawEvent] = withdrawEventParse
}

func RegisterHandle(evt common.Hash, handle EventHandle) {
	handleList[evt] = handle
}

func LendingPoolEventHandle(eventLog types.Log) error {
	retVar, err := parseList[eventLog.Topics[0]](eventLog)
	if err != nil {
		return err
	}
	if handleList[eventLog.Topics[0]] != nil {
		handleList[eventLog.Topics[0]](eventLog.BlockNumber, eventLog.Index, retVar)
	}

	return nil
}

func ParseLendingPoolEvent(eventLog types.Log, evt common.Hash) (interface{}, error) {
	check := func() bool {
		return eventLog.Topics[0] == evt
	}
	switch eventLog.Topics[0] {
	case BorrowEvent:
		if !check() {
			break
		}
		return borrowEventParse(eventLog)
	case DepositEvent:
		if !check() {
			break
		}
		return depositEventParse(eventLog)
	case FlashLoanEvent:
		if !check() {
			break
		}
		return flashLoanEventParse(eventLog)
	case LiquidationEvent:
		if !check() {
			break
		}
		return liquidationEventParse(eventLog)
	case PausedEvent:
		if !check() {
			break
		}
		return pausedEventParse(eventLog)
	case RebalanceStableRateEvent:
		if !check() {
			break
		}
		return rebalanceStableRateEventParse(eventLog)
	case RepayEvent:
		if !check() {
			break
		}
		return repayEventParse(eventLog)
	case ReserveUpdatedEvent:
		if !check() {
			break
		}
		return reserveUpdatedEventParse(eventLog)
	case CollateralDisableEvent:
		if !check() {
			break
		}
		return collateralDisableEventParse(eventLog)
	case CollateralEnableEvent:
		if !check() {
			break
		}
		return collateralEnableEventParse(eventLog)
	case SwapEvent:
		if !check() {
			break
		}
		return swapEventParse(eventLog)
	case UnpausedEvent:
		if !check() {
			break
		}
		return unpausedEventParse(eventLog)
	case WithdrawEvent:
		if !check() {
			break
		}
		return withdrawEventParse(eventLog)
	default:
		return nil, fmt.Errorf("can not Parse event(%v) in block(%v) TX(%v) index(%v)", eventLog.Topics[0], eventLog.BlockNumber, eventLog.TxHash, eventLog.Index)
	}
	return nil, nil
}

func borrowEventParse(eventLog types.Log) (interface{}, error) {
	datas, err := contractAbi.Unpack("Borrow", eventLog.Data)
	if err != nil {
		return nil, fmt.Errorf("borrowEventParse unpack: %v", err)
	}
	retval := BorrowEventData{}
	retval.Reserve = common.BytesToAddress(eventLog.Topics[1].Bytes())
	retval.User = datas[0].(common.Address)
	retval.OnBehalfOf = common.BytesToAddress(eventLog.Topics[2].Bytes())
	retval.Amount = datas[1].(*big.Int)
	retval.BorrowRateMode = datas[2].(*big.Int)
	retval.BorrowRate = datas[3].(*big.Int)
	retval.Referral = utils.BytesToUInt16(eventLog.Topics[3].Bytes())
	//fmt.Print(retval.User)
	return retval, nil
}
func depositEventParse(eventLog types.Log) (interface{}, error) {
	datas, err := contractAbi.Unpack("Deposit", eventLog.Data)
	if err != nil {
		return nil, fmt.Errorf("depositEventParse unpack: %v", err)
	}
	retval := DepositEventData{}
	retval.Reserve = common.BytesToAddress(eventLog.Topics[1].Bytes())
	retval.User = datas[0].(common.Address)
	retval.OnBehalfOf = common.BytesToAddress(eventLog.Topics[2].Bytes())
	retval.Amount = datas[1].(*big.Int)
	retval.Referral = utils.BytesToUInt16(eventLog.Topics[3].Bytes())
	//fmt.Print(retval)
	return retval, nil
}
func flashLoanEventParse(eventLog types.Log) (interface{}, error) {
	retval := FlashLoanEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func liquidationEventParse(eventLog types.Log) (interface{}, error) {
	retval := LiquidationCallEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func pausedEventParse(eventLog types.Log) (interface{}, error) {
	return nil, nil
}
func rebalanceStableRateEventParse(eventLog types.Log) (interface{}, error) {
	retval := RebalanceStableBorrowRateEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func repayEventParse(eventLog types.Log) (interface{}, error) {
	datas, err := contractAbi.Unpack("Repay", eventLog.Data)
	if err != nil {
		return nil, fmt.Errorf("repayEventParse unpack: %v", err)
	}
	retval := RepayEventData{}
	retval.Reserve = common.BytesToAddress(eventLog.Topics[1].Bytes())
	retval.User = common.BytesToAddress(eventLog.Topics[2].Bytes())
	retval.Repayer = common.BytesToAddress(eventLog.Topics[3].Bytes())
	retval.Amount = datas[0].(*big.Int)
	//fmt.Print(retval)
	return retval, nil
}
func reserveUpdatedEventParse(eventLog types.Log) (interface{}, error) {
	retval := ReserveDataUpdatedEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func collateralDisableEventParse(eventLog types.Log) (interface{}, error) {
	retval := ReserveUsedAsCollateralDisabledEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func collateralEnableEventParse(eventLog types.Log) (interface{}, error) {
	retval := ReserveUsedAsCollateralEnabledEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func swapEventParse(eventLog types.Log) (interface{}, error) {
	retval := SwapEventData{}
	//fmt.Print(retval)
	return retval, nil
}
func unpausedEventParse(eventLog types.Log) (interface{}, error) {
	return nil, nil
}
func withdrawEventParse(eventLog types.Log) (interface{}, error) {
	datas, err := contractAbi.Unpack("Withdraw", eventLog.Data)
	if err != nil {
		return nil, fmt.Errorf("withdrawEventParse unpack: %v", err)
	}
	retval := WithdrawEventData{}
	retval.Reserve = common.BytesToAddress(eventLog.Topics[1].Bytes())
	retval.User = common.BytesToAddress(eventLog.Topics[2].Bytes())
	retval.To = common.BytesToAddress(eventLog.Topics[3].Bytes())
	retval.Amount = datas[0].(*big.Int)
	//fmt.Print(retval)
	return retval, nil
}
