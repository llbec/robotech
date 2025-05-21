package EventEmitter

import (
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	AddEvent         = crypto.Keccak256Hash([]byte("Add(address,address,address,address,uint256,uint256,uint256)"))
	RemoveEvent      = crypto.Keccak256Hash([]byte("Remove(address,address,address,uint256,address,uint256,uint256)"))
	DepositEvent     = crypto.Keccak256Hash([]byte("Deposit(address,address,address,uint256,uint256,uint256,uint256,uint256,uint256)"))
	WithdrawEvent    = crypto.Keccak256Hash([]byte("Withdraw(address,address,address,uint256,uint256,address,uint256,uint256,uint256,uint256)"))
	PositionEvent    = crypto.Keccak256Hash([]byte("Position(address,uint256,address,address,uint256,uint256,uint256,uint256,uint256)"))
	SwapEvent        = crypto.Keccak256Hash([]byte("Swap(address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)"))
	BorrowEvent      = crypto.Keccak256Hash([]byte("Borrow(address,address,address,uint256,uint8,uint256,uint256,uint256,uint256,uint256,uint256)"))
	RepayEvent       = crypto.Keccak256Hash([]byte("Repay(address,address,address,uint256,uint8,uint256,uint256,uint256,uint256,uint256)"))
	LiquidationEvent = crypto.Keccak256Hash([]byte("Liquidation(address,address,uint256,uint256,uint256,uint256,uint256,uint256)"))
	CloseEvent       = crypto.Keccak256Hash([]byte("Close(address,uint256)"))
	PoolCreatedEvent = crypto.Keccak256Hash([]byte("PoolCreated(address,address,address,uint256,uint256,uint256)"))
	PoolUpdatedEvent = crypto.Keccak256Hash([]byte("PoolUpdated(address,uint256,uint256,uint256,uint256)"))
	ClaimFeesEvent   = crypto.Keccak256Hash([]byte("ClaimFees(address,uint256,uint256,uint256)"))
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

var (
	EvenNames = []string{
		"Add",
		"Remove",
		"Deposit",
		"Withdraw",
		"Position",
		"Swap",
		"Borrow",
		"Repay",
		"Liquidation",
		"Close",
		"PoolCreated",
		"PoolUpdated",
		"ClaimFees",
	}
)
