package backend_test

import (
	"crypto/ecdsa"
	"math/big"
	"robotech/contract/margin/v2/backend"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	backendContract *backend.Backend
	secret          *ecdsa.PrivateKey
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":          "string",
		"backend_contract": "string",
	})
	if err != nil {
		panic(err)
	}
	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	backendContract, err = backend.NewBackend(
		common.HexToAddress(cfgMap["backend_contract"].(string)),
		ethUtils.GetClient(),
	)
	if err != nil {
		panic(err)
	}

	secret, err = crypto.HexToECDSA(cfgMap["secret"].(string))
	if err != nil {
		panic(err)
	}
}

func Test_verifySwapSignature(t *testing.T) {
	t.Log("Test_verifySwapSignature")
	// 1. 获取签名
	// 2. 调用 verifySwapSignature 函数验证签名
	// 3. 打印结果

	bedDebt := big.NewInt(0)
	price := big.NewInt(1000000000000)
	height := big.NewInt(24904376)
	list := [][32]byte{}
	sigStr := "0x6b4daf72ab9c99f13200fb725906b17d42244685410d2c4a33ded46bbc7949ae"
	signature := common.Hex2Bytes(sigStr[2:])

	err := backendContract.VerifySwapSignature(nil, bedDebt, price, height, list, signature)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func Test_verifyRemoveSignature(t *testing.T) {
	t.Log("Test_verifyRemoveSignature")
	// 1. 获取签名
	// 2. 调用 verifyRemoveSignature 函数验证签名
	// 3. 打印结果

	maxAmount0 := big.NewInt(0)
	price := big.NewInt(1000000000000)
	height := big.NewInt(24904376)
	sigStr := "0x6b4daf72ab9c99f13200fb725906b17d42244685410d2c4a33ded46bbc7949ae"
	signature := common.Hex2Bytes(sigStr[2:])

	err := backendContract.VerifyRemoveSignature(nil, maxAmount0, height, price, signature)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}
