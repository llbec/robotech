package laf

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// -------- LafAgent --------
type LafAgent struct {
	client           *ethclient.Client
	lafContract      common.Address
	stakingContract  common.Address
	referralContract common.Address
	usdtContract     common.Address
	swapContract     common.Address
}

// -------- LafAgentConfig --------
type LafAgentConfig struct {
	RpcUrl           string `json:"rpc_url"`
	LafContract      string `json:"laf_contract"`
	StakingContract  string `json:"staking_contract"`
	ReferralContract string `json:"referral_contract"`
	USDTContract     string `json:"usdt_contract"`
	SwapContract     string `json:"swap_contract"`
}

// -------- LafAgent DB struct --------
type LafAgentData struct {
	Type        string `json:"type"`
	Transaction string `json:"transaction"`
}

// -------- LafAgent transaction types --------
const (
	TransferEvent = "Transfer"
)
