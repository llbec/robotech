package main

import (
	"fmt"
	"testing"

	"github.com/llbec/robotech/squadrons/filecoinsquadron"
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

func Test_Radar(t *testing.T) {
	LoadConfig(".")
	initEnv()
	start := int64(1095520)
	current := int64(1115240)

	for last := start; last <= current; last++ {
		targets := make(map[string]filecoinsquadron.MsgInfo)
		tipsetBytes, err := filecoinAPI.GetTipsetByHeight(last)
		if err != nil {
			t.Fatalf("GetTipsetByHeight Height(%v): %v", last, err)
		}

		tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
		if err != nil {
			t.Fatalf("ReadTipSet Height(%v): %v", last, err)
		}

		for _, b := range tipset.Blocks {
			msgs, err := filecoinAPI.PayeeRadarInBlock(payeeMap, b)
			if err != nil {
				t.Fatalf("PayeeRadarInBlock Height(%v): %v", last, err)
			}
			for _, m := range msgs {
				targets[m.Cid] = m
				fmt.Println("got one ", m.From, m.To, m.Value)
			}
		}
	}
}
