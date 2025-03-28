package role_test

import (
	"crypto/ecdsa"
	v1 "robotech/contract/margin/v1"
	"robotech/contract/margin/v1/role"
	"robotech/utils"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	roleStoreContract *role.RoleStore
	dataStoreAddr     common.Address
	account_1         common.Address
	roleName          string
	secret            *ecdsa.PrivateKey
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		"rpc_url":        "string",
		"roleStore_addr": "string",
		"dataStore_addr": "string",
		"account_1":      "string",
		"roleName":       "string",
		"secret":         "string",
	})
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(cfgMap["rpc_url"].(string))
	roleStoreContract, err = role.NewRoleStore(common.HexToAddress(cfgMap["roleStore_addr"].(string)), ethUtils.GetClient())
	if err != nil {
		panic(err)
	}
	dataStoreAddr = common.HexToAddress(cfgMap["dataStore_addr"].(string))
	account_1 = common.HexToAddress(cfgMap["account_1"].(string))
	roleName = cfgMap["roleName"].(string)

	secret, err = crypto.HexToECDSA(cfgMap["secret"].(string))
	if err != nil {
		panic(err)
	}
}

func Test_hasRole(t *testing.T) {
	ok, err := roleStoreContract.HasRole(nil, account_1, v1.ABIEncode(roleName))
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Logf("%v has role %s", account_1, roleName)
	} else {
		t.Logf("%v has not role %s", account_1, roleName)
	}
}

func Test_grantRole(t *testing.T) {
	opt, err := ethUtils.GetOpt(secret)
	if err != nil {
		t.Fatalf("get opt failed: %v", err)
	}
	tx, err := roleStoreContract.GrantRole(opt, account_1, v1.ABIEncode(roleName))
	if err != nil {
		t.Error(err)
	}
	t.Logf("tx hash: %s", tx.Hash().Hex())

	/*time.Sleep(time.Second * 5)*/
	receipt, err := ethUtils.WaitTransaction(tx)
	if err != nil {
		t.Error(err)
	}
	t.Logf("receipt: %+v", receipt)

	ok, err := roleStoreContract.HasRole(nil, account_1, v1.ABIEncode(roleName))
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatalf("%v has not role %s", account_1, roleName)
	}
}
