package main

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_GenerateCfg(t *testing.T) {
	GeneratConfig("")
}

func Test_readkey(t *testing.T) {
	k, e := getSecret()
	fmt.Println(k, e)
}

func Test_percent(t *testing.T) {
	for i := 0; i <= 100; i++ {
		dpercent := deltaPercent(big.NewInt(100), big.NewInt(int64(i)))
		if dpercent.Cmp(big.NewInt(int64(100-i))) != 0 {
			t.Errorf("100 - %d -> %v\n", i, dpercent)
		}
	}
	for i := 100; i <= 200; i++ {
		dpercent := deltaPercent(big.NewInt(100), big.NewInt(int64(i)))
		if dpercent.Cmp(big.NewInt(int64(i-100))) != 0 {
			t.Errorf("100 - %d -> %v\n", i, dpercent)
		}
	}
	dpercent := deltaPercent(big.NewInt(100), big.NewInt(int64(10000)))
	if dpercent.Cmp(big.NewInt(9900)) != 0 {
		t.Errorf("100 - 10000 -> %v\n", dpercent)
	}
}
