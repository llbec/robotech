package main

import (
	"testing"
)

func Test_generateCfg(t *testing.T) {
	err := GeneratConfig(".")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_LoadCfg(t *testing.T) {
	LoadConfig(".")
}
