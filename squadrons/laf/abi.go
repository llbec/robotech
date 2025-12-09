package laf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// -------- Referral events --------
// event BindReferral(address indexed user,address parent);
const referralABI = `
[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"user","type":"address"},
			{"indexed":false,"internalType":"address","name":"parent","type":"address"}
		],
		"name":"BindReferral",
		"type":"event"
	}
]`

type BindReferralEvent struct {
	User   common.Address `json:"user"`
	Parent common.Address `json:"parent"`
}

// -------- Staking events --------
// event Staked(address indexed user, uint256 amount, uint256 timestamp, uint256 index, uint256 stakeTime);
// event RewardPaid(address indexed user, uint256 reward, uint40 timestamp, uint256 index);
// event Transfer(address indexed from, address indexed to, uint256 amount);
const stakingABI = `
[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"user","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"index","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"stakeTime","type":"uint256"}
		],
		"name":"Staked",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"user","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"reward","type":"uint256"},
			{"indexed":false,"internalType":"uint40","name":"timestamp","type":"uint40"},
			{"indexed":false,"internalType":"uint256","name":"index","type":"uint256"}
		],
		"name":"RewardPaid",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"from","type":"address"},
			{"indexed":true,"internalType":"address","name":"to","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}
		],
		"name":"Transfer",
		"type":"event"
	}
]`

type StakedEvent struct {
	User      common.Address `json:"user"`
	Amount    *big.Int       `json:"amount"`
	Timestamp *big.Int       `json:"timestamp"`
	Index     *big.Int       `json:"index"`
	StakeTime *big.Int       `json:"stakeTime"`
}

type RewardPaidEvent struct {
	User      common.Address `json:"user"`
	Reward    *big.Int       `json:"reward"`
	Timestamp *big.Int       `json:"timestamp"`
	Index     *big.Int       `json:"index"`
}

type StakeingTransferEvent struct {
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Amount *big.Int       `json:"amount"`
}

// -------- LAF events --------
/*
	event Approval(address indexed owner, address indexed spender, uint256 amount);
    event ExcludedFromFee(address account);
    event IncludedToFee(address account);
    event OwnershipTransferred(address indexed user, address indexed newOwner);
    event Transfer(address indexed from, address indexed to, uint256 amount);
*/
const LAFABI = `
[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"owner","type":"address"},
			{"indexed":true,"internalType":"address","name":"spender","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}
		],
		"name":"Approval",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":false,"internalType":"address","name":"account","type":"address"}
		],
		"name":"ExcludedFromFee",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":false,"internalType":"address","name":"account","type":"address"}
		],
		"name":"IncludedToFee",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"user","type":"address"},
			{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}
		],
		"name":"OwnershipTransferred",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"from","type":"address"},
			{"indexed":true,"internalType":"address","name":"to","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}
		],
		"name":"Transfer",
		"type":"event"
	}
]`

type ApprovalEvent struct {
	Owner   common.Address `json:"owner"`
	Spender common.Address `json:"spender"`
	Amount  *big.Int       `json:"amount"`
}
type ExcludedFromFeeEvent struct {
	Account common.Address `json:"account"`
}
type IncludedToFeeEvent struct {
	Account common.Address `json:"account"`
}
type OwnershipTransferredEvent struct {
	User     common.Address `json:"user"`
	NewOwner common.Address `json:"newOwner"`
}
type LAFTransferEvent struct {
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Amount *big.Int       `json:"amount"`
}

// -------- Swap events --------
// Sync (uint112 reserve0, uint112 reserve1)View Source
const SwapABI = `
[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},
			{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}
		],
		"name":"Sync",
		"type":"event"
	}
]`

type SyncEvent struct {
	Reserve0 *big.Int `json:"reserve0"`
	Reserve1 *big.Int `json:"reserve1"`
}
