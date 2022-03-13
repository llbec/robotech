package main

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/llbec/robotech/utils"
)

func test_init() (err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return
	}
	loadConfig(workDir)
	err = envInit()
	if err != nil {
		return
	}
	return
}

func Test_Balance(t *testing.T) {
	e := test_init()
	if e != nil {
		t.Fatal(e.Error())
	}

	amount, e := poolBalance()
	if e != nil {
		t.Fatal(e.Error())
	}
	tar := big.NewInt(5e13)
	fmt.Printf("pool balance is %.6f\n", utils.BigToRecognizable(amount, 18))

	if amount.Cmp(tar) > 0 {
		fmt.Printf("%v\n%v\namount > target\n", amount, tar)
	} else {
		fmt.Printf("%v\n%v\namount <= target\n", amount, tar)
	}

	amount, e = usrBalance()
	if e != nil {
		t.Fatal(e.Error())
	}
	fmt.Printf("address balance is %.6f\n", utils.BigToRecognizable(amount, 18))
}
