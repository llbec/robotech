package ctokenevent

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/llbec/robotech/armory/compound"
)

var (
	contractAbi = abi.ABI{}
)

func init() {
	cabi, err := abi.JSON(strings.NewReader(string(compound.CTokenABI)))
	if err != nil {
		log.Fatalf("Parsing contract abi error: %v\n", err)
	}
	contractAbi = cabi
}

func HandleEvent(eventLog types.Log, evt common.Hash) (interface{}, error) {
	check := func() bool {
		return eventLog.Topics[0] == evt
	}
	switch eventLog.Topics[0] {
	case BorrowEvent:
		if !check() {
			break
		}
		return borrowEventHandle(eventLog)
	default:
		//return nil, fmt.Errorf("unknow event(%v) in block(%v) TX(%v) index(%v)", eventLog.Topics[0], eventLog.BlockNumber, eventLog.TxHash, eventLog.Index)
		return nil, nil
	}
	return nil, nil
}

func borrowEventHandle(eventLog types.Log) (interface{}, error) {
	datas, err := contractAbi.Unpack("Borrow", eventLog.Data)
	if err != nil {
		return nil, fmt.Errorf("borrowEventHandle unpack: %v", err)
	}
	retval := BorrowEventData{}
	retval.Borrower = datas[0].(common.Address)
	retval.BorrowAmount = datas[1].(*big.Int)
	retval.AccountBorrows = datas[2].(*big.Int)
	retval.TotalBorrows = datas[3].(*big.Int)
	//fmt.Print(retval.User)
	return retval, nil
}
