package Reader_test

import (
	"math/big"
	v1 "robotech/contract/margin/v1"
	Reader "robotech/contract/margin/v1/reader"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	readerContract *Reader.Reader
	dataStoreAddr  common.Address
	account        common.Address
	usdt           common.Address
	meme           common.Address
	threshold      *big.Int
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":        "string",
		"reader_addr":    "string",
		"dataStore_addr": "string",
		"account":        "string",
		"usdt":           "string",
		"meme":           "string",
	})
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	readerContract, err = Reader.NewReader(common.HexToAddress(cfgMap["reader_addr"].(string)), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	dataStoreAddr = common.HexToAddress(cfgMap["dataStore_addr"].(string))
	account = common.HexToAddress(cfgMap["account"].(string))
	usdt = common.HexToAddress(cfgMap["usdt"].(string))
	meme = common.HexToAddress(cfgMap["meme"].(string))

	threshold, _ = big.NewInt(0).SetString("1100000000000000000000000000", 10)
}

func Test_calcAmountOut(t *testing.T) {
	t.Logf("calcAmountOut params:\n\tdataStore: %v\n\tusdt: %v\n\tmeme: %v\n\tamountIn: %v",
		dataStoreAddr.Hex(), usdt.Hex(), meme.Hex(), big.NewInt(10000000000),
	)
	amountOut, amountFee, priceImpact, err := readerContract.CalcAmountOut(
		nil,
		dataStoreAddr,
		usdt,
		meme,
		big.NewInt(10000000000),
		0,
	)
	if err != nil {
		t.Fatalf("CalcAmountOut error: %v", err)
	}
	t.Logf("\tamountOut:%v\n\tamountFee:%v\npriceImpact:%v", amountOut, amountFee, priceImpact)
}

func Test_GetPosition(t *testing.T) {
	positions, err := readerContract.GetPositions(nil, dataStoreAddr, account)
	if err != nil {
		t.Fatalf("GetPositions error: %v", err)
	}
	for _, position := range positions {
		if position.MarginLevel.Cmp(threshold) < 1 {
			t.Logf("MarginLevel(%s)(%d): %v", position.Account, position.Id, position.MarginLevel)
		}
		t.Logf("Position%v", position)
	}
}

func Test_GetPositions(t *testing.T) {
	positions1, err := readerContract.GetPositions(nil, dataStoreAddr, account)
	if err != nil {
		t.Fatalf("GetPositions1 error: %v", err)
	}
	//t.Logf("GetPositions1: %v", positions1)

	positionKeys := [][32]byte{}
	for _, position := range positions1 {
		positionKey := [32]byte(v1.CalcPositionKey(position.Account, position.Id))
		positionKeys = append(positionKeys, positionKey)
		t.Logf("positionKey(%s)(%d): %v", position.Account, position.Id, positionKey)
	}
	t.Logf("positionKeys: %v", positionKeys)
	// TODO: implement
	positions, err := readerContract.GetPositions2(nil, dataStoreAddr, positionKeys)
	if err != nil {
		t.Fatalf("GetPositions2 error: %v", err)
	}
	//t.Logf("GetPositions2: %v", positions)
	if len(positions) != len(positions1) {
		t.Fatalf("GetPositions2 length mismatch: %v != %v", len(positions), len(positions1))
	}
	for i := range positions {
		if positions[i].MarginLevel.Cmp(positions1[i].MarginLevel) != 0 {
			t.Fatalf("GetPositions2 mismatch: \n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v",
				positions[i].Assets[1].Symbol, positions1[i].Assets[1].Symbol,
				positions1[i].Account, positions[i].Account,
				positions1[i].Id, positions[i].Id,
				positions1[i].MarginLevel, positions[i].MarginLevel,
				positions1[i], positions[i],
			)
		}
	}
}
