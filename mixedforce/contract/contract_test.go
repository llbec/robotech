package main

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_ReadABI(t *testing.T) {
	err := ReadABI(".abi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(contractABI)
}

func Test_ViewFunction(t *testing.T) {
	err := ReadABI(".abi")
	if err != nil {
		t.Fatalf(err.Error())
	}

	contract, err := NewContract("0x6537d6307ca40231939985BCF7D83096Dd1B4C09", "https://http-mainnet-node.huobichain.com")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var out []interface{}
	err = contract.ContractCaller.contract.Call(nil, &out, "compAccrued", common.HexToAddress("0x744808B56E9B470CE901ABaa9A6261E92928cE82"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(out...)
}
