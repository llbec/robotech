package laf_test

import (
	"encoding/json"
	"os"
	"robotech/squadrons/laf"
	"testing"
)

type TestConfig struct {
	ReferralStartBlock uint64 `json:"referral_create"`
	StakingStartBlock  uint64 `json:"staking_create"`
	LafStartBlock      uint64 `json:"laf_create"`
}

var TestEnv = "../../.env"

func loadConfig() (*TestConfig, error) {
	file, err := os.Open(TestEnv)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg TestConfig
	err = json.NewDecoder(file).Decode(&cfg)
	return &cfg, err
}

func Test_FilterTxs_referral(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	agent := laf.NewLAFAgent(TestEnv)
	txs, err := agent.FilterTxs(cfg.ReferralStartBlock, cfg.ReferralStartBlock+9)
	if err != nil {
		t.Errorf("FilterTxs failed: %v", err)
	}
	if len(txs) == 0 {
		t.Errorf("FilterTxs returned empty txs")
	}
	t.Logf("FilterTxs returned %d txs:\n%v", len(txs), txs)
}

func Test_FilterTxs_staking(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	agent := laf.NewLAFAgent(TestEnv)
	txs, err := agent.FilterTxs(cfg.StakingStartBlock, cfg.StakingStartBlock+9)
	if err != nil {
		t.Errorf("FilterTxs failed: %v", err)
	}
	if len(txs) == 0 {
		t.Errorf("FilterTxs returned empty txs")
	}
	t.Logf("FilterTxs returned %d txs:\n%v", len(txs), txs)
}

func Test_FilterTxs_laf(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	agent := laf.NewLAFAgent(TestEnv)
	txs, err := agent.FilterTxs(cfg.LafStartBlock, cfg.LafStartBlock+9)
	if err != nil {
		t.Errorf("FilterTxs failed: %v", err)
	}
	if len(txs) == 0 {
		t.Errorf("FilterTxs returned empty txs")
	}
	t.Logf("FilterTxs returned %d txs:\n%v", len(txs), txs)
}
