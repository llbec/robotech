package lendingpoolevent

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	BorrowEvent              = crypto.Keccak256Hash([]byte("Borrow(address,address,address,uint256,uint256,uint256,uint16)"))
	DepositEvent             = crypto.Keccak256Hash([]byte("Deposit(address,address,address,uint256,uint16)"))
	FlashLoanEvent           = crypto.Keccak256Hash([]byte("FlashLoan(address,address,address,uint256,uint256,uint16)"))
	LiquidationEvent         = crypto.Keccak256Hash([]byte("LiquidationCall(address,address,address,uint256,uint256,address,bool)"))
	PausedEvent              = crypto.Keccak256Hash([]byte("Paused()"))
	RebalanceStableRateEvent = crypto.Keccak256Hash([]byte("RebalanceStableBorrowRate(address,address)"))
	RepayEvent               = crypto.Keccak256Hash([]byte("Repay(address,address,address,uint256)"))
	ReserveUpdatedEvent      = crypto.Keccak256Hash([]byte("ReserveDataUpdated(address,uint256,uint256,uint256,uint256,uint256)"))
	CollateralDisableEvent   = crypto.Keccak256Hash([]byte("ReserveUsedAsCollateralDisabled(address,address)"))
	CollateralEnableEvent    = crypto.Keccak256Hash([]byte("ReserveUsedAsCollateralEnabled(address,address)"))
	SwapEvent                = crypto.Keccak256Hash([]byte("Swap(address,address,uint256)"))
	UnpausedEvent            = crypto.Keccak256Hash([]byte("Unpaused()"))
	WithdrawEvent            = crypto.Keccak256Hash([]byte("Withdraw(address,address,address,uint256)"))
)

type EventHandle func(height uint64, logIndex uint, eventData interface{})

type BorrowEventData struct {
	Reserve        common.Address
	User           common.Address
	OnBehalfOf     common.Address
	Amount         *big.Int
	BorrowRateMode *big.Int
	BorrowRate     *big.Int
	Referral       uint16
}

type DepositEventData struct {
	Reserve    common.Address
	User       common.Address
	OnBehalfOf common.Address
	Amount     *big.Int
	Referral   uint16
}

type FlashLoanEventData struct {
	Target    common.Address
	Initiator common.Address
	Asset     common.Address
	Amount    *big.Int
	Premium   *big.Int
	Referral  uint16
}

type LiquidationCallEventData struct {
	CollateralAsset            common.Address
	DebtAsset                  common.Address
	User                       common.Address
	DebtToCover                *big.Int
	LiquidatedCollateralAmount *big.Int
	Liquidator                 common.Address
	ReceiveAToken              bool
}

type RebalanceStableBorrowRateEventData struct {
	Reserve common.Address
	User    common.Address
}

type RepayEventData struct {
	Reserve common.Address
	User    common.Address
	Repayer common.Address
	Amount  *big.Int
}

type ReserveDataUpdatedEventData struct {
	Reserve             common.Address
	LiquidityRate       *big.Int
	StableBorrowRate    *big.Int
	VariableBorrowRate  *big.Int
	LiquidityIndex      *big.Int
	VariableBorrowIndex *big.Int
}

type ReserveUsedAsCollateralDisabledEventData struct {
	Reserve common.Address
	User    common.Address
}

type ReserveUsedAsCollateralEnabledEventData struct {
	Reserve common.Address
	User    common.Address
}

type SwapEventData struct {
	Reserve  common.Address
	User     common.Address
	RateMode *big.Int
}

type WithdrawEventData struct {
	Reserve common.Address
	User    common.Address
	To      common.Address
	Amount  *big.Int
}
