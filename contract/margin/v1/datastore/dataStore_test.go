package DataStore_test

import (
	"crypto/ecdsa"
	v1 "robotech/contract/margin/v1"
	DataStore "robotech/contract/margin/v1/datastore"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	dataStoreContract *DataStore.DataStore
	secret            *ecdsa.PrivateKey
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":    "string",
		"data_store": "string",
		"secret":     "string",
	})
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	dataStoreContract, err = DataStore.NewDataStore(common.HexToAddress(cfgMap["data_store"].(string)), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}

	secret, err = crypto.HexToECDSA(cfgMap["secret"].(string))
	if err != nil {
		panic(err)
	}
}

// 测试获取地址函数
func Test_getAddress(t *testing.T) {
	t.Log("Test_getAddress")

	address, err := dataStoreContract.GetAddress(nil, v1.ABIEncode("BackendContract"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address: %s", address.Hex())

	address, err = dataStoreContract.GetAddress(nil, v1.ABIEncode("PositionHandleContract"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address: %s", address.Hex())

	height, err := dataStoreContract.GetUint(nil, v1.ABIEncode("HeightLimit"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("height limit: %d", height)

	priceLimit, err := dataStoreContract.GetUint(nil, v1.ABIEncode("PriceLimit"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("price limit: %d", priceLimit)
}
