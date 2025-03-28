package config_test

import (
	"crypto/ecdsa"
	"robotech/contract/margin/v1/config"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	configContract *config.Config
	dataStoreAddr  common.Address
	treasury       common.Address
	secret         *ecdsa.PrivateKey
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":         "string",
		"config_contract": "string",
		"data_store":      "string",
		"treasury":        "string",
		"secret":          "string",
	})
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	configContract, err = config.NewConfig(common.HexToAddress(cfgMap["config_contract"].(string)), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	dataStoreAddr = common.HexToAddress(cfgMap["data_store"].(string))
	treasury = common.HexToAddress(cfgMap["treasury"].(string))

	secret, err = crypto.HexToECDSA(cfgMap["secret"].(string))
	if err != nil {
		panic(err)
	}
}

func Test_setTreasury(t *testing.T) {
	opts, err := ethUtils.GetOpt(secret)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := configContract.SetTreasury(opts, treasury)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx: %s", tx.Hash().Hex())
	receipt, err := ethUtils.WaitTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("receipt: %+v", receipt)
}
