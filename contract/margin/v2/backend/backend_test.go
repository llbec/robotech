package backend_test

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"robotech/contract/margin/v2/backend"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
		"secret":           "string",
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

func Test_verifySignature(t *testing.T) {
	t.Log("Test_verifyRemoveSignature")
	// 1. 获取签名
	// 2. 调用 verifyRemoveSignature 函数验证签名
	// 3. 打印结果

	maxAmount0 := big.NewInt(0)
	price := big.NewInt(1000000000000)
	height := big.NewInt(24999012)

	hashBytes, err := backendContract.GetRemoveHash(nil, maxAmount0, price, height)
	if err != nil {
		t.Fatal(err)
	}
	msgHash := common.BytesToHash(hashBytes[:])
	t.Log("msgHash:", msgHash.Hex())

	eip191Msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msgHash), msgHash)
	eip191Hash := crypto.Keccak256Hash([]byte(eip191Msg))

	signature, err := crypto.Sign(eip191Hash.Bytes(), secret)
	if err != nil {
		t.Fatal(err)
	}
	signature[crypto.RecoveryIDOffset] += 27
	t.Log("signature:", hexutil.Encode(signature))
	signer, err := backendContract.RecoverSigner(nil, hashBytes, signature)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signer:", signer)
	t.Log("wallet:", crypto.PubkeyToAddress(secret.PublicKey).Hex())
}
