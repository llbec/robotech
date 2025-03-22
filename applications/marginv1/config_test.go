package main

import "testing"

func Test_generateCfg(t *testing.T) {
	err := GeneratConfig()
	if err != nil {
		t.Fatalf("GeneratConfig error: %v", err)
	}
}
