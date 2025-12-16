package lafabi

/*
Referral
	event BindReferral(address indexed user,address parent);
	event SetOperator(address indexed operators, bool status);
*/
const REFERRALABI = `
[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"user","type":"address"},
			{"indexed":false,"internalType":"address","name":"parent","type":"address"}
		],
		"name":"BindReferral",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"operators","type":"address"},
			{"indexed":true,"internalType":"bool","name":"status","type":"bool"}
		],
		"name":"SetOperator",
		"type":"event"
	}
]`

/*
Staking
	event Staked(address indexed user, uint256 amount, uint256 timestamp, uint256 index, uint256 stakeTime);
	event RewardPaid(address indexed user, uint256 reward, uint40 timestamp, uint256 index);
	event Transfer(address indexed from, address indexed to, uint256 amount);
	event OwnershipTransferred(address indexed user, address indexed newOwner);
*/
const STAKINGABI = `
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
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"user","type":"address"},
			{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}
		],
		"name":"OwnershipTransferred",
		"type":"event"
	}
]`

/*
LAF
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
