package laf

import (
	"math/big"

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
	routeContract    common.Address
}

// -------- LafAgentConfig --------
type LafAgentConfig struct {
	RpcUrl           string `json:"rpc_url"`
	LafContract      string `json:"laf_contract"`
	StakingContract  string `json:"staking_contract"`
	ReferralContract string `json:"referral_contract"`
	USDTContract     string `json:"usdt_contract"`
	SwapContract     string `json:"swap_contract"`
	RouteContract    string `json:"route_contract"`
}

// -------- LafAgent DB struct --------
type LafTransaction struct {
	Type        string `json:"type"`
	Transaction string `json:"transaction"`
}

// -------- LafAgent transaction types --------
const (
	LAFTransfer = "laf_transfer"
	LPTransfer  = "lp_transfer"
	Stake       = "stake"
	Unstake     = "unstake"
	Swap        = "swap"
	AddLP       = "add_lp"
)

type LafTransferTX struct {
	From  common.Address `json:"from"`
	To    common.Address `json:"to"`
	Value *big.Int       `json:"value"`
}
type LPTransferTX struct {
	From  common.Address `json:"from"`
	To    common.Address `json:"to"`
	Value *big.Int       `json:"value"`
}
type StakeTX struct {
	From      common.Address `json:"from"`
	Amount    *big.Int       `json:"amount"`
	StartTime *big.Int       `json:"start_time"`
	UserIndex *big.Int       `json:"user_index"`
	StakeTime *big.Int       `json:"stake_time"`
	SwapLAF   *big.Int       `json:"swap_laf"`
	SwapUSDT  *big.Int       `json:"swap_usdt"`
	LPLaf     *big.Int       `json:"lp_laf"`
	LPUST     *big.Int       `json:"lp_usdt"`
	SyncLaf   *big.Int       `json:"sync_laf"`
	SyncUSDT  *big.Int       `json:"sync_usdt"`
}
type Dividend struct {
	Account  common.Address `json:"account"`
	Dividend *big.Int       `json:"dividend"`
}
type UnstakeTX struct {
	From        common.Address `json:"from"`
	Amount      *big.Int       `json:"amount"`
	Reward      *big.Int       `json:"reward"`
	UserIndex   *big.Int       `json:"user_index"`
	UnstakeTime *big.Int       `json:"unstake_time"`
	SwapLAF     *big.Int       `json:"swap_laf"`
	SwapUSDT    *big.Int       `json:"swap_usdt"`
	Cycle       *big.Int       `json:"cycle"`
	SyncLaf     *big.Int       `json:"sync_laf"`
	SyncUSDT    *big.Int       `json:"sync_usdt"`
	Dividends   []Dividend     `json:"dividends"`
}
type SwapTX struct {
	From     common.Address `json:"from"`
	LafIn    *big.Int       `json:"laf_in"`
	USDTIn   *big.Int       `json:"usdt_in"`
	LafOut   *big.Int       `json:"laf_out"`
	USDTOut  *big.Int       `json:"usdt_out"`
	SyncLaf  *big.Int       `json:"sync_laf"`
	SyncUSDT *big.Int       `json:"sync_usdt"`
}
type AddLP_TX struct {
	From     common.Address `json:"from"`
	LafIn    *big.Int       `json:"laf_in"`
	USDTIn   *big.Int       `json:"usdt_in"`
	SyncLaf  *big.Int       `json:"sync_laf"`
	SyncUSDT *big.Int       `json:"sync_usdt"`
}
type RemoveLP_TX struct {
	From     common.Address `json:"from"`
	LafOut   *big.Int       `json:"laf_out"`
	USDTOut  *big.Int       `json:"usdt_out"`
	SyncLaf  *big.Int       `json:"sync_laf"`
	SyncUSDT *big.Int       `json:"sync_usdt"`
}
