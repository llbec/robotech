package laf

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LafAgent struct {
	client           *ethclient.Client
	lafContract      common.Address
	stakingContract  common.Address
	referralContract common.Address
	usdtContract     common.Address
	swapContract     common.Address
}

type LafAgentConfig struct {
	RpcUrl           string `json:"rpcUrl"`
	LafContract      string `json:"lafContract"`
	StakingContract  string `json:"stakingContract"`
	ReferralContract string `json:"referralContract"`
	USDTContract     string `json:"usdtContract"`
	SwapContract     string `json:"swapContract"`
}
