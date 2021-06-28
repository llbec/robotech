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

var (
	contractAbi = abi.ABI{}
)

func init() {
	cabi, err := abi.JSON(strings.NewReader(string(aave.LendingPoolABI)))
	if err != nil {
		log.Fatalf("Parsing contract abi error: %v\n", err)
	}
	contractAbi = cabi
}

func HandleLendingPoolEvent(eventLog types.Log, evt common.Hash) (interface{}, error) {
	check := func() bool {
		return eventLog.Topics[0] == evt
	}
	switch eventLog.Topics[0] {
	case BorrowEvent:
		if !check() {
			break
		}
		return borrowEventHandle(eventLog)
	case DepositEvent:
		if !check() {
			break
		}
		return depositEventHandle(eventLog)
	case FlashLoanEvent:
		if !check() {
			break
		}
		return flashLoanEventHandle(eventLog)
	case LiquidationEvent:
		if !check() {
			break
		}
		return liquidationEventHandle(eventLog)
	case PausedEvent:
		if !check() {
			break
		}
		return pausedEventHandle(eventLog)
	case RebalanceStableRateEvent:
		if !check() {
			break
		}
		return rebalanceStableRateEventHandle(eventLog)
	case RepayEvent:
		if !check() {
			break
		}
		return repayEventHandle(eventLog)
	case ReserveUpdatedEvent:
		if !check() {
			break
		}
		return reserveUpdatedEventHandle(eventLog)
	case CollateralDisableEvent:
		if !check() {
			break
		}
		return collateralDisableEventHandle(eventLog)
	case CollateralEnableEvent:
		if !check() {
			break
		}
		return collateralEnableEventHandle(eventLog)
	case SwapEvent:
		if !check() {
			break
		}
		return swapEventHandle(eventLog)
	case UnpausedEvent:
		if !check() {
			break
		}
		return unpausedEventHandle(eventLog)
	case WithdrawEvent:
		if !check() {
			break
		}
		return withdrawEventHandle(eventLog)
	default:
		return nil, fmt.Errorf("unknow event(%v) in block(%v) TX(%v) index(%v)", eventLog.Topics[0], eventLog.BlockNumber, eventLog.TxHash, eventLog.Index)
	}
	return nil, nil
}

func borrowEventHandle(eventLog types.Log) (interface{}, error) {
	datas, err := contractAbi.Unpack("Borrow", eventLog.Data)
	if err != nil {
		return nil, fmt.Errorf("borrowEventHandle unpack: %v", err)
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
func depositEventHandle(eventLog types.Log) (interface{}, error) {
	retval := DepositEventData{}
	fmt.Print(retval)
	return retval, nil
}
func flashLoanEventHandle(eventLog types.Log) (interface{}, error) {
	retval := FlashLoanEventData{}
	fmt.Print(retval)
	return retval, nil
}
func liquidationEventHandle(eventLog types.Log) (interface{}, error) {
	retval := LiquidationCallEventData{}
	fmt.Print(retval)
	return retval, nil
}
func pausedEventHandle(eventLog types.Log) (interface{}, error) {
	return nil, nil
}
func rebalanceStableRateEventHandle(eventLog types.Log) (interface{}, error) {
	retval := RebalanceStableBorrowRateEventData{}
	fmt.Print(retval)
	return retval, nil
}
func repayEventHandle(eventLog types.Log) (interface{}, error) {
	retval := RepayEventData{}
	fmt.Print(retval)
	return retval, nil
}
func reserveUpdatedEventHandle(eventLog types.Log) (interface{}, error) {
	retval := ReserveDataUpdatedEventData{}
	fmt.Print(retval)
	return retval, nil
}
func collateralDisableEventHandle(eventLog types.Log) (interface{}, error) {
	retval := ReserveUsedAsCollateralDisabledEventData{}
	fmt.Print(retval)
	return retval, nil
}
func collateralEnableEventHandle(eventLog types.Log) (interface{}, error) {
	retval := ReserveUsedAsCollateralEnabledEventData{}
	fmt.Print(retval)
	return retval, nil
}
func swapEventHandle(eventLog types.Log) (interface{}, error) {
	retval := SwapEventData{}
	fmt.Print(retval)
	return retval, nil
}
func unpausedEventHandle(eventLog types.Log) (interface{}, error) {
	return nil, nil
}
func withdrawEventHandle(eventLog types.Log) (interface{}, error) {
	retval := WithdrawEventData{}
	fmt.Print(retval)
	return retval, nil
}
