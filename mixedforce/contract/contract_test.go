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
	//contractABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"compAccrued\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"accrued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

	err := ReadABI(".abi")
	if err != nil {
		t.Fatalf(err.Error())
	}

	contract, err := NewContract("0xb74633f2022452f377403B638167b0A135DB096d", "https://http-mainnet-node.huobichain.com")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var out []interface{}
	err = contract.ContractCaller.contract.Call(nil, &out, "compAccrued", common.HexToAddress("0x7c0Fc56fE45c98cCF000DD548F457D2fCA9ECC6e"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(out...)
}
