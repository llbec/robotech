package DataStore_test

import (
	"crypto/ecdsa"
	"fmt"
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

func Test_Bytes32(t *testing.T) {
	t.Log("Test_Bytes32")

	backend := v1.ABIEncode("BackendContract")
	hex := utils.Bytes2Hex(backend[:])
	fmt.Println(hex)

	positionHandle := v1.ABIEncode("PositionHandleContract")
	hex = utils.Bytes2Hex(positionHandle[:])
	fmt.Println(hex)
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

func Test_getFunctions(t *testing.T) {
	t.Log("Test_getFunctions")
	// 测试获取函数
	backend, err := dataStoreContract.GetBackend(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("backend: %s", backend.Hex())

	positionHandle, err := dataStoreContract.GetPositionHandle(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("position handle: %s", positionHandle.Hex())
}
