package main

import (
	"fmt"
	"testing"
)

func Test_GenerateCfg(t *testing.T) {
	GeneratConfig("")
}

func Test_readkey(t *testing.T) {
	k, e := getSecret()
	fmt.Println(k, e)
}
