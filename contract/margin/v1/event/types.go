package EventEmitter

import (
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	AddEvent      = crypto.Keccak256Hash([]byte("Add(address,address,address,address,uint256,uint256,uint256)"))
	RemoveEvent   = crypto.Keccak256Hash([]byte("Remove(address,address,address,uint256,address,uint256,uint256)"))
	DepositEvent  = crypto.Keccak256Hash([]byte("Deposit(address,address,address,uint256,uint256,uint256,uint256,uint256,uint256)"))
	WithdrawEvent = crypto.Keccak256Hash([]byte("Withdraw(address,address,address,uint256,uint256,address,uint256,uint256,uint256,uint256)"))
	PositionEvent = crypto.Keccak256Hash([]byte("Position(address,uint256,address,address,uint256,uint256,uint256,uint256,uint256)"))
	SwapEvent     = crypto.Keccak256Hash([]byte("Swap(address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)"))
)

const (
	DEPOSIT     = 0 //存
	BORROW      = 1 //借
	REPAY       = 2 //还
	WITHDRAW    = 3 //赎回
	SWAP        = 4 //交易
	LIQUIDATION = 5 //清算
	CLOSE       = 6 //关闭
)
