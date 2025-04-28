package backendCheck_test

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	backendCheck "robotech/contract/margin/v2/backendcheck"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	backendCheckContract *backendCheck.BackendCheck
	secret               *ecdsa.PrivateKey
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":               "string",
		"backendCheck_contract": "string",
		"secret":                "string",
	})
	if err != nil {
		panic(err)
	}
	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	backendCheckContract, err = backendCheck.NewBackendCheck(
		common.HexToAddress(cfgMap["backendCheck_contract"].(string)),
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

func verify(hashBytes []byte) (sig string, err error) {
	signature, err := ethUtils.SignHash(secret, hashBytes[:])
	if err != nil {
		return
	}
	//fmt.println("signature:", hexutil.Encode(signature))
	sig = hexutil.Encode(signature)

	signer, err := ethUtils.RecoverSignature(signature, hashBytes)
	if err != nil {
		return
	}
	if signer != crypto.PubkeyToAddress(secret.PublicKey) {
		err = errors.New("signer != secret.PublicKey")
		return
	}
	return
}

func decodeRemoveHash(maxAmount0, height, price *big.Int) []byte {
	// 定义参数类型
	uint256Type, _ := abi.NewType("uint256", "", nil)

	// 创建 Arguments 切片
	arguments := abi.Arguments{
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: uint256Type},
	}

	// 使用 Pack 方法进行编码
	packedData, _ := arguments.Pack(maxAmount0, price, height)

	return crypto.Keccak256Hash(packedData).Bytes()
}

func Test_SwapSign(t *testing.T) {
	t.Log("Test_Sign")
	// 1. 获取签名
	// 2. 调用 verifySwapSignature 函数验证签名
	// 3. 打印结果

	bedDebt := big.NewInt(0)
	price := big.NewInt(1000000000000)
	height := big.NewInt(24904376)
	list := [][32]byte{}
	//sigStr := "0x6b4daf72ab9c99f13200fb725906b17d42244685410d2c4a33ded46bbc7949ae"
	//signature := common.Hex2Bytes(sigStr[2:])

	hashBytes, err := backendCheckContract.GetSwapHash(nil, bedDebt, price, height, list)
	if err != nil {
		t.Fatal(err)
	}
	swapHash := common.BytesToHash(hashBytes[:])
	t.Log("swapHash:", swapHash.Hex())

	signature, err := verify(hashBytes[:])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signature:", signature)

	/*signature, err := ethUtils.SignHash(secret, hashBytes[:])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signature:", hexutil.Encode(signature))

	signer, err := ethUtils.RecoverSignature(signature, swapHash.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if signer != crypto.PubkeyToAddress(secret.PublicKey) {
		t.Fatal("signer != secret.PublicKey")
	}
	t.Log("signer:", signer.Hex())*/
}

func Test_simpleSign(t *testing.T) {
	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	// 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8
	t.Log(hash.Hex())

	signature, err := verify(hash.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	// 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301
	t.Log("signature:", signature)
}

func Test_RemoveSign(t *testing.T) {
	maxAmount0 := big.NewInt(0)
	price := big.NewInt(1000000000000)
	height := big.NewInt(24999012)

	loaclHash := decodeRemoveHash(maxAmount0, height, price)
	//t.Log("loaclHash:", hexutil.Encode(loaclHash))
	t.Log("loaclHash:", common.BytesToHash(loaclHash[:]).Hex())

	signature, err := verify(loaclHash[:])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("local signature:", signature)
	t.Log("local Signer :", crypto.PubkeyToAddress(secret.PublicKey).Hex())

	hashBytes, err := backendCheckContract.GetRemoveHash(nil, maxAmount0, price, height)
	if err != nil {
		t.Fatal(err)
	}
	removeHash := common.BytesToHash(hashBytes[:])
	t.Log("removeHash:", removeHash.Hex())

	signature, err = verify(hashBytes[:])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signature:", signature)
	t.Log("Signer :", crypto.PubkeyToAddress(secret.PublicKey).Hex())
}

func Test_verify(t *testing.T) {
	hashBytes := common.Hex2Bytes("f36be22725629f45c08643f82c136e74387f89f40a7e35922c9b786f9543c85f")
	t.Log("hashBytes:", hexutil.Encode(hashBytes))
	signature, err := ethUtils.SignHash(secret, hashBytes)
	if err != nil {
		t.Fatalf("SignHash error: %v", err)
	}
	t.Log("signature:", hexutil.Encode(signature))

	signer, err := ethUtils.RecoverSignature(signature, hashBytes)
	if err != nil {
		t.Fatalf("RecoverSignature error: %v", err)
	}
	t.Log("signer:", signer.Hex())

	hashBytes = common.Hex2Bytes("cf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd")
	signature = common.Hex2Bytes("993dab3dd91f5c6dc28e17439be475478f5635c92a56e17e82349d3fb2f166196f466c0b4e0c146f285204f0dcb13e5ae67bc33f4b888ec32dfe0a063e8f3f781b")
	signer, err = ethUtils.RecoverSignature(signature, hashBytes)
	if err != nil {
		t.Fatalf("RecoverSignature error: %v", err)
	}
	t.Log("signer:", signer.Hex())
}
