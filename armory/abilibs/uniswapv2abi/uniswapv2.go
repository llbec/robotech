package uniswapv2abi

// -------- Swap events --------
/*
	event Mint(address indexed sender, uint amount0, uint amount1);
    event Burn(address indexed sender, uint amount0, uint amount1, address indexed to);
    event Swap(
        address indexed sender,
        uint amount0In,
        uint amount1In,
        uint amount0Out,
        uint amount1Out,
        address indexed to
    );
    event Sync(uint112 reserve0, uint112 reserve1);
	event Transfer(address indexed from, address indexed to, uint value);
*/
const SWAPABI = `
[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"sender","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount0","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"amount1","type":"uint256"}
		],
		"name":"Mint",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"sender","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount0","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"amount1","type":"uint256"},
			{"indexed":true,"internalType":"address","name":"to","type":"address"}
		],
		"name":"Burn",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"sender","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},
			{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"},
			{"indexed":true,"internalType":"address","name":"to","type":"address"}
		],
		"name":"Swap",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},
			{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}
		],
		"name":"Sync",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"from","type":"address"},
			{"indexed":true,"internalType":"address","name":"to","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}
		],
		"name":"Transfer",
		"type":"event"
	}
]`
