package storage_test

import (
	"math/big"
	"robotech/applications/marginv1/storage"
	"robotech/applications/marginv1/types"
	"robotech/utils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	test_account common.Address
	test_token   common.Address
	test_id      *big.Int
	test_height  uint64
)

func init() {
	rMap, err := utils.ReadConfig(map[string]string{
		"account":         "string",
		"token":           "string",
		"id":              "int64",
		"lastBlockNumber": "uint64",
	})
	if err != nil {
		panic(err)
	}
	test_account = common.HexToAddress(rMap["account"].(string))
	test_token = common.HexToAddress(rMap["token"].(string))
	test_id = big.NewInt(rMap["id"].(int64))
	test_height = rMap["lastBlockNumber"].(uint64)
}

func Test_SetLastBlockNumber(t *testing.T) {
	t.Logf("Test_SetLastBlockNumber testing at %v => %v", storage.GetLastBlockNumber(), test_height)
	storage.SetLastBlockNumber(test_height)
	storageNum := storage.GetLastBlockNumber()
	if storageNum != test_height {
		t.Errorf("Expected %d, got %d", test_height, storageNum)
	}
	t.Logf("Last block number: %d", storageNum)
}

func Test_Storage(t *testing.T) {
	t.Logf("Test_Storage testing at %v => %v => %v", test_token, test_account, test_id)
	err := storage.AddPool(test_token)
	if err != nil {
		t.Errorf("AddPool error: %v", err)
	}
	ok, err := storage.HasPool(test_token)
	if err != nil {
		t.Errorf("HasPool error: %v", err)
	}
	if !ok {
		t.Errorf("Expected true, got false")
	}
	ok, err = storage.HasPool(test_account)
	if err != nil {
		t.Errorf("HasPool error: %v", err)
	}
	if ok {
		t.Errorf("Expected false, got true")
	}
	err = storage.AddPosition(test_token, test_account, types.Position{
		MemeToken:  test_token,
		Account:    test_account,
		PositionId: test_id,
	})
	if err != nil {
		t.Errorf("AddPosition error: %v", err)
	}
	positions := storage.GetPositions(test_token)
	if len(positions) != 1 {
		t.Errorf("Expected 1, got %d", len(positions))
	}
	if positions[0].MemeToken != test_token ||
		positions[0].Account != test_account ||
		positions[0].PositionId.Cmp(test_id) != 0 {
		t.Errorf("Expected %v, got %v", types.Position{
			MemeToken:  test_token,
			Account:    test_account,
			PositionId: test_id,
		}, positions[0])
	}
}
