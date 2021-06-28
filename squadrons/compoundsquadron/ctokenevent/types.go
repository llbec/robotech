package ctokenevent

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	AccrueInterestEvent             = crypto.Keccak256Hash([]byte("AccrueInterest(uint256,uint256,uint256,uint256)"))
	ApprovalEvent                   = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
	BorrowEvent                     = crypto.Keccak256Hash([]byte("Borrow(address,uint256,uint256,uint256)"))
	FailureEvent                    = crypto.Keccak256Hash([]byte("Failure(uint256,uint256,uint256)"))
	LiquidateBorrowEvent            = crypto.Keccak256Hash([]byte("LiquidateBorrow(address,address,uint256,address,uint256)"))
	MintEvent                       = crypto.Keccak256Hash([]byte("Mint(address,uint256,uint256)"))
	NewAdminEvent                   = crypto.Keccak256Hash([]byte("NewAdmin(address,address)"))
	NewComptrollerEvent             = crypto.Keccak256Hash([]byte("NewComptroller(address,address)"))
	NewMarketInterestRateModelEvent = crypto.Keccak256Hash([]byte("NewMarketInterestRateModel(address,address)"))
	NewPendingAdminEvent            = crypto.Keccak256Hash([]byte("NewPendingAdmin(address,address)"))
	NewReserveFactorEvent           = crypto.Keccak256Hash([]byte("NewReserveFactor(uint256,uint256)"))
	RedeemEvent                     = crypto.Keccak256Hash([]byte("Redeem(address,uint256,uint256)"))
	RepayBorrowEvent                = crypto.Keccak256Hash([]byte("RepayBorrow(address,address,uint256,uint256,uint256)"))
	ReservesAddedEvent              = crypto.Keccak256Hash([]byte("ReservesAdded(address,uint256,uint256)"))
	ReservesReducedEvent            = crypto.Keccak256Hash([]byte("ReservesReduced(address,uint256,uint256)"))
	TransferEvent                   = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
)

type BorrowEventData struct {
	Borrower       common.Address
	BorrowAmount   *big.Int
	AccountBorrows *big.Int
	TotalBorrows   *big.Int
}
