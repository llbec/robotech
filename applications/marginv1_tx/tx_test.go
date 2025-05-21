package main_test

import (
	"context"
	"math/big"
	EventEmitter "robotech/contract/margin/v1/event"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	txHash               common.Hash
	contractAddress      common.Address
	block_number         *big.Int
	contractEventEmitter *EventEmitter.EventEmitter
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":          "string",
		"tx_hash":          "string",
		"contract_address": "string",
		"block_number":     "int64",
	})
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	txHash = common.HexToHash(cfgMap["tx_hash"].(string))
	contractAddress = common.HexToAddress(cfgMap["contract_address"].(string))
	block_number = big.NewInt(cfgMap["block_number"].(int64))

	contract, err := EventEmitter.NewEventEmitter(contractAddress, ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	contractEventEmitter = contract
}

func Test_tx(t *testing.T) {
	tx, isPending, err := ethUtils.GetClient().TransactionByHash(context.Background(), txHash)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx)
	t.Log(isPending)
}

func Test_logs(t *testing.T) {
	query := ethereum.FilterQuery{
		FromBlock: block_number,
		ToBlock:   block_number,
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := ethUtils.GetClient().FilterLogs(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(EventEmitter.EventEmitterABI)))
	if err != nil {
		t.Fatal(err)
	}

	for _, vLog := range logs {
		switch vLog.Topics[0] {
		case contractAbi.Events["Add"].ID:
			event, err := contractEventEmitter.ParseAdd(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nAdd\n", event)
		case contractAbi.Events["Remove"].ID:
			event, err := contractEventEmitter.ParseRemove(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nRemove\n", event)
		case contractAbi.Events["Deposit"].ID:
			event, err := contractEventEmitter.ParseDeposit(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nDeposit\n", event)
		case contractAbi.Events["Withdraw"].ID:
			event, err := contractEventEmitter.ParseWithdraw(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nWithdraw\n", event)
		case contractAbi.Events["Position"].ID:
			event, err := contractEventEmitter.ParsePosition(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nPosition\n", event)
		case contractAbi.Events["Swap"].ID:
			event, err := contractEventEmitter.ParseSwap(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nSwap\n", event)
		case contractAbi.Events["Borrow"].ID:
			event, err := contractEventEmitter.ParseBorrow(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nBorrow\n", event)
		case contractAbi.Events["Repay"].ID:
			event, err := contractEventEmitter.ParseRepay(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nRepay\n", event)
		case contractAbi.Events["Liquidation"].ID:
			event, err := contractEventEmitter.ParseLiquidation(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nLiquidation\n", event)
		case contractAbi.Events["Close"].ID:
			event, err := contractEventEmitter.ParseClose(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nClose\n", event)
		case contractAbi.Events["PoolCreated"].ID:
			event, err := contractEventEmitter.ParsePoolCreated(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nPoolCreated\n", event)
		case contractAbi.Events["PoolUpdated"].ID:
			event, err := contractEventEmitter.ParsePoolUpdated(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nPoolUpdated\n", event)
		case contractAbi.Events["ClaimFees"].ID:
			event, err := contractEventEmitter.ParseClaimFees(vLog)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("\nClaimFees\n", event)
		default:
			t.Log(vLog.Topics[0].Hex())
		}
	}
}
