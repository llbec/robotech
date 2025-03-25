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
	account_1      common.Address
	account_2      common.Address
	threshold      *big.Int
)

func init() {
	configMap := map[string]string{
		"rpc_url":        "string",
		"reader_addr":    "string",
		"dataStore_addr": "string",
		"account_1":      "string",
		"account_2":      "string",
	}

	envMap, err := utils.ReadConfig(configMap)
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(envMap["rpc_url"].(string))
	readerContract, err = Reader.NewReader(common.HexToAddress(envMap["reader_addr"].(string)), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	dataStoreAddr = common.HexToAddress(envMap["dataStore_addr"].(string))
	account_1 = common.HexToAddress(envMap["account_1"].(string))
	account_2 = common.HexToAddress(envMap["account_2"].(string))

	threshold, _ = big.NewInt(0).SetString("1100000000000000000000000000", 10)
}

func Test_GetPosition(t *testing.T) {
	positions, err := readerContract.GetPositions(nil, dataStoreAddr, account_2)
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
	positions1, err := readerContract.GetPositions(nil, dataStoreAddr, account_1)
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
