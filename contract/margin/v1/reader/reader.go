// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Reader

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ReaderPoolUtilsAsset is an auto generated low-level Go binding around an user-defined struct.
type ReaderPoolUtilsAsset struct {
	Token                   common.Address
	Symbol                  string
	Decimals                *big.Int
	BorrowIndex             *big.Int
	BorrowApy               *big.Int
	TotalCollateral         *big.Int
	TotalCollateralWithDebt *big.Int
	TotalDebtScaled         *big.Int
	PoolBalance             *big.Int
	PriceLiquidity          *big.Int
	AvaiableLoan            *big.Int
	ActualTradable          *big.Int
}

// ReaderPoolUtilsAssetInfo is an auto generated low-level Go binding around an user-defined struct.
type ReaderPoolUtilsAssetInfo struct {
	Token       common.Address
	Symbol      string
	Decimals    *big.Int
	BorrowIndex *big.Int
}

// ReaderPoolUtilsGetPool is an auto generated low-level Go binding around an user-defined struct.
type ReaderPoolUtilsGetPool struct {
	Assets               [2]ReaderPoolUtilsAsset
	Bank                 common.Address
	InterestRateStrategy common.Address
	Configuration        *big.Int
	LastUpdateTimestamp  *big.Int
	PriceDecimals        *big.Int
	Price                *big.Int
	Source               common.Address
	ShortEnabled         bool
	CreatedTimestamp     *big.Int
	UnclaimedFee         *big.Int
	MemeMaxDepositAmount *big.Int
	AveragePrice         *big.Int
}

// ReaderPoolUtilsGetPoolInfo is an auto generated low-level Go binding around an user-defined struct.
type ReaderPoolUtilsGetPoolInfo struct {
	Assets        [2]ReaderPoolUtilsAssetInfo
	PriceDecimals *big.Int
	Price         *big.Int
}

// ReaderPositionUtilsAsset is an auto generated low-level Go binding around an user-defined struct.
type ReaderPositionUtilsAsset struct {
	Token           common.Address
	Symbol          string
	Decimals        *big.Int
	Balance         *big.Int
	Debt            *big.Int
	NetWorth        *big.Int
	MaxRedeemAmount *big.Int
	BorrowApy       *big.Int
	Equity          *big.Int
	EquityValue     *big.Int
}

// ReaderPositionUtilsGetPosition is an auto generated low-level Go binding around an user-defined struct.
type ReaderPositionUtilsGetPosition struct {
	Assets             [2]ReaderPositionUtilsAsset
	Id                 *big.Int
	Account            common.Address
	MarginLevel        *big.Int
	EntryPrice         *big.Int
	IndexPrice         *big.Int
	Pnl                *big.Int
	LiquidationPrice   *big.Int
	ToLiquidationPrice *big.Int
}

// ReaderMetaData contains all meta data concerning the Reader contract.
var ReaderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenOutIndex\",\"type\":\"uint8\"}],\"name\":\"calcAmountIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenInIndex\",\"type\":\"uint8\"}],\"name\":\"calcAmountOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"calcLiquidityOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"name\":\"calcTokenPairOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getDefaultInterestRateStrategy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getDefaultPoolConfiguration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getMarginLevelThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getPools\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowApy\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateralWithDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"poolBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"avaiableLoan\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualTradable\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.Asset[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"address\",\"name\":\"bank\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"configuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shortEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"createdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unclaimedFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeMaxDepositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"averagePrice\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.GetPool[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getPools\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowApy\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateralWithDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"poolBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"avaiableLoan\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualTradable\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.Asset[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"address\",\"name\":\"bank\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"configuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shortEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"createdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unclaimedFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeMaxDepositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"averagePrice\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.GetPool[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"poolKeys\",\"type\":\"bytes32[]\"}],\"name\":\"getPools2\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowApy\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateralWithDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"poolBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"avaiableLoan\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualTradable\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.Asset[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"address\",\"name\":\"bank\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"configuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shortEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"createdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unclaimedFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeMaxDepositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"averagePrice\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.GetPool[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getPoolsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getPoolsInfo\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.AssetInfo[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint256\",\"name\":\"priceDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.GetPoolInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getPoolsInfo\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.AssetInfo[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint256\",\"name\":\"priceDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.GetPoolInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"poolKeys\",\"type\":\"bytes32[]\"}],\"name\":\"getPoolsInfo\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.AssetInfo[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint256\",\"name\":\"priceDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"internalType\":\"structReaderPoolUtils.GetPoolInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getPositions\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"netWorth\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"maxRedeemAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowApy\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"equity\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"equityValue\",\"type\":\"int256\"}],\"internalType\":\"structReaderPositionUtils.Asset[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"marginLevel\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"entryPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"IndexPrice\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"pnl\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationPrice\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"toLiquidationPrice\",\"type\":\"int256\"}],\"internalType\":\"structReaderPositionUtils.GetPosition[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"positionKeys\",\"type\":\"bytes32[]\"}],\"name\":\"getPositions2\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"netWorth\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"maxRedeemAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowApy\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"equity\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"equityValue\",\"type\":\"int256\"}],\"internalType\":\"structReaderPositionUtils.Asset[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"marginLevel\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"entryPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"IndexPrice\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"pnl\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationPrice\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"toLiquidationPrice\",\"type\":\"int256\"}],\"internalType\":\"structReaderPositionUtils.GetPosition[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getTokenBase\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dataStore\",\"type\":\"address\"}],\"name\":\"getTreasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ReaderABI is the input ABI used to generate the binding from.
// Deprecated: Use ReaderMetaData.ABI instead.
var ReaderABI = ReaderMetaData.ABI

// Reader is an auto generated Go binding around an Ethereum contract.
type Reader struct {
	ReaderCaller     // Read-only binding to the contract
	ReaderTransactor // Write-only binding to the contract
	ReaderFilterer   // Log filterer for contract events
}

// ReaderCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReaderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReaderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReaderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReaderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReaderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReaderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReaderSession struct {
	Contract     *Reader           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReaderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReaderCallerSession struct {
	Contract *ReaderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ReaderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReaderTransactorSession struct {
	Contract     *ReaderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReaderRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReaderRaw struct {
	Contract *Reader // Generic contract binding to access the raw methods on
}

// ReaderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReaderCallerRaw struct {
	Contract *ReaderCaller // Generic read-only contract binding to access the raw methods on
}

// ReaderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReaderTransactorRaw struct {
	Contract *ReaderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReader creates a new instance of Reader, bound to a specific deployed contract.
func NewReader(address common.Address, backend bind.ContractBackend) (*Reader, error) {
	contract, err := bindReader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Reader{ReaderCaller: ReaderCaller{contract: contract}, ReaderTransactor: ReaderTransactor{contract: contract}, ReaderFilterer: ReaderFilterer{contract: contract}}, nil
}

// NewReaderCaller creates a new read-only instance of Reader, bound to a specific deployed contract.
func NewReaderCaller(address common.Address, caller bind.ContractCaller) (*ReaderCaller, error) {
	contract, err := bindReader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReaderCaller{contract: contract}, nil
}

// NewReaderTransactor creates a new write-only instance of Reader, bound to a specific deployed contract.
func NewReaderTransactor(address common.Address, transactor bind.ContractTransactor) (*ReaderTransactor, error) {
	contract, err := bindReader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReaderTransactor{contract: contract}, nil
}

// NewReaderFilterer creates a new log filterer instance of Reader, bound to a specific deployed contract.
func NewReaderFilterer(address common.Address, filterer bind.ContractFilterer) (*ReaderFilterer, error) {
	contract, err := bindReader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReaderFilterer{contract: contract}, nil
}

// bindReader binds a generic wrapper to an already deployed contract.
func bindReader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReaderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reader *ReaderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reader.Contract.ReaderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reader *ReaderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reader.Contract.ReaderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reader *ReaderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reader.Contract.ReaderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reader *ReaderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reader *ReaderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reader *ReaderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reader.Contract.contract.Transact(opts, method, params...)
}

// CalcAmountIn is a free data retrieval call binding the contract method 0xc2bdeda1.
//
// Solidity: function calcAmountIn(address dataStore, address token0, address token1, uint256 amountOut, uint8 tokenOutIndex) view returns(uint256, uint256, int256)
func (_Reader *ReaderCaller) CalcAmountIn(opts *bind.CallOpts, dataStore common.Address, token0 common.Address, token1 common.Address, amountOut *big.Int, tokenOutIndex uint8) (*big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "calcAmountIn", dataStore, token0, token1, amountOut, tokenOutIndex)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// CalcAmountIn is a free data retrieval call binding the contract method 0xc2bdeda1.
//
// Solidity: function calcAmountIn(address dataStore, address token0, address token1, uint256 amountOut, uint8 tokenOutIndex) view returns(uint256, uint256, int256)
func (_Reader *ReaderSession) CalcAmountIn(dataStore common.Address, token0 common.Address, token1 common.Address, amountOut *big.Int, tokenOutIndex uint8) (*big.Int, *big.Int, *big.Int, error) {
	return _Reader.Contract.CalcAmountIn(&_Reader.CallOpts, dataStore, token0, token1, amountOut, tokenOutIndex)
}

// CalcAmountIn is a free data retrieval call binding the contract method 0xc2bdeda1.
//
// Solidity: function calcAmountIn(address dataStore, address token0, address token1, uint256 amountOut, uint8 tokenOutIndex) view returns(uint256, uint256, int256)
func (_Reader *ReaderCallerSession) CalcAmountIn(dataStore common.Address, token0 common.Address, token1 common.Address, amountOut *big.Int, tokenOutIndex uint8) (*big.Int, *big.Int, *big.Int, error) {
	return _Reader.Contract.CalcAmountIn(&_Reader.CallOpts, dataStore, token0, token1, amountOut, tokenOutIndex)
}

// CalcAmountOut is a free data retrieval call binding the contract method 0xd28b0a15.
//
// Solidity: function calcAmountOut(address dataStore, address token0, address token1, uint256 amountIn, uint8 tokenInIndex) view returns(uint256, uint256, int256)
func (_Reader *ReaderCaller) CalcAmountOut(opts *bind.CallOpts, dataStore common.Address, token0 common.Address, token1 common.Address, amountIn *big.Int, tokenInIndex uint8) (*big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "calcAmountOut", dataStore, token0, token1, amountIn, tokenInIndex)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// CalcAmountOut is a free data retrieval call binding the contract method 0xd28b0a15.
//
// Solidity: function calcAmountOut(address dataStore, address token0, address token1, uint256 amountIn, uint8 tokenInIndex) view returns(uint256, uint256, int256)
func (_Reader *ReaderSession) CalcAmountOut(dataStore common.Address, token0 common.Address, token1 common.Address, amountIn *big.Int, tokenInIndex uint8) (*big.Int, *big.Int, *big.Int, error) {
	return _Reader.Contract.CalcAmountOut(&_Reader.CallOpts, dataStore, token0, token1, amountIn, tokenInIndex)
}

// CalcAmountOut is a free data retrieval call binding the contract method 0xd28b0a15.
//
// Solidity: function calcAmountOut(address dataStore, address token0, address token1, uint256 amountIn, uint8 tokenInIndex) view returns(uint256, uint256, int256)
func (_Reader *ReaderCallerSession) CalcAmountOut(dataStore common.Address, token0 common.Address, token1 common.Address, amountIn *big.Int, tokenInIndex uint8) (*big.Int, *big.Int, *big.Int, error) {
	return _Reader.Contract.CalcAmountOut(&_Reader.CallOpts, dataStore, token0, token1, amountIn, tokenInIndex)
}

// CalcLiquidityOut is a free data retrieval call binding the contract method 0xf68a7131.
//
// Solidity: function calcLiquidityOut(address dataStore, address token0, address token1, uint256 amount0, uint256 amount1) view returns(uint256)
func (_Reader *ReaderCaller) CalcLiquidityOut(opts *bind.CallOpts, dataStore common.Address, token0 common.Address, token1 common.Address, amount0 *big.Int, amount1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "calcLiquidityOut", dataStore, token0, token1, amount0, amount1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcLiquidityOut is a free data retrieval call binding the contract method 0xf68a7131.
//
// Solidity: function calcLiquidityOut(address dataStore, address token0, address token1, uint256 amount0, uint256 amount1) view returns(uint256)
func (_Reader *ReaderSession) CalcLiquidityOut(dataStore common.Address, token0 common.Address, token1 common.Address, amount0 *big.Int, amount1 *big.Int) (*big.Int, error) {
	return _Reader.Contract.CalcLiquidityOut(&_Reader.CallOpts, dataStore, token0, token1, amount0, amount1)
}

// CalcLiquidityOut is a free data retrieval call binding the contract method 0xf68a7131.
//
// Solidity: function calcLiquidityOut(address dataStore, address token0, address token1, uint256 amount0, uint256 amount1) view returns(uint256)
func (_Reader *ReaderCallerSession) CalcLiquidityOut(dataStore common.Address, token0 common.Address, token1 common.Address, amount0 *big.Int, amount1 *big.Int) (*big.Int, error) {
	return _Reader.Contract.CalcLiquidityOut(&_Reader.CallOpts, dataStore, token0, token1, amount0, amount1)
}

// CalcTokenPairOut is a free data retrieval call binding the contract method 0x317b50ec.
//
// Solidity: function calcTokenPairOut(address dataStore, address token0, address token1, uint256 liquidity) view returns(uint256, uint256)
func (_Reader *ReaderCaller) CalcTokenPairOut(opts *bind.CallOpts, dataStore common.Address, token0 common.Address, token1 common.Address, liquidity *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "calcTokenPairOut", dataStore, token0, token1, liquidity)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// CalcTokenPairOut is a free data retrieval call binding the contract method 0x317b50ec.
//
// Solidity: function calcTokenPairOut(address dataStore, address token0, address token1, uint256 liquidity) view returns(uint256, uint256)
func (_Reader *ReaderSession) CalcTokenPairOut(dataStore common.Address, token0 common.Address, token1 common.Address, liquidity *big.Int) (*big.Int, *big.Int, error) {
	return _Reader.Contract.CalcTokenPairOut(&_Reader.CallOpts, dataStore, token0, token1, liquidity)
}

// CalcTokenPairOut is a free data retrieval call binding the contract method 0x317b50ec.
//
// Solidity: function calcTokenPairOut(address dataStore, address token0, address token1, uint256 liquidity) view returns(uint256, uint256)
func (_Reader *ReaderCallerSession) CalcTokenPairOut(dataStore common.Address, token0 common.Address, token1 common.Address, liquidity *big.Int) (*big.Int, *big.Int, error) {
	return _Reader.Contract.CalcTokenPairOut(&_Reader.CallOpts, dataStore, token0, token1, liquidity)
}

// GetDefaultInterestRateStrategy is a free data retrieval call binding the contract method 0xe335adb7.
//
// Solidity: function getDefaultInterestRateStrategy(address dataStore) view returns(address)
func (_Reader *ReaderCaller) GetDefaultInterestRateStrategy(opts *bind.CallOpts, dataStore common.Address) (common.Address, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getDefaultInterestRateStrategy", dataStore)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDefaultInterestRateStrategy is a free data retrieval call binding the contract method 0xe335adb7.
//
// Solidity: function getDefaultInterestRateStrategy(address dataStore) view returns(address)
func (_Reader *ReaderSession) GetDefaultInterestRateStrategy(dataStore common.Address) (common.Address, error) {
	return _Reader.Contract.GetDefaultInterestRateStrategy(&_Reader.CallOpts, dataStore)
}

// GetDefaultInterestRateStrategy is a free data retrieval call binding the contract method 0xe335adb7.
//
// Solidity: function getDefaultInterestRateStrategy(address dataStore) view returns(address)
func (_Reader *ReaderCallerSession) GetDefaultInterestRateStrategy(dataStore common.Address) (common.Address, error) {
	return _Reader.Contract.GetDefaultInterestRateStrategy(&_Reader.CallOpts, dataStore)
}

// GetDefaultPoolConfiguration is a free data retrieval call binding the contract method 0x50ed592d.
//
// Solidity: function getDefaultPoolConfiguration(address dataStore) view returns(uint256)
func (_Reader *ReaderCaller) GetDefaultPoolConfiguration(opts *bind.CallOpts, dataStore common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getDefaultPoolConfiguration", dataStore)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDefaultPoolConfiguration is a free data retrieval call binding the contract method 0x50ed592d.
//
// Solidity: function getDefaultPoolConfiguration(address dataStore) view returns(uint256)
func (_Reader *ReaderSession) GetDefaultPoolConfiguration(dataStore common.Address) (*big.Int, error) {
	return _Reader.Contract.GetDefaultPoolConfiguration(&_Reader.CallOpts, dataStore)
}

// GetDefaultPoolConfiguration is a free data retrieval call binding the contract method 0x50ed592d.
//
// Solidity: function getDefaultPoolConfiguration(address dataStore) view returns(uint256)
func (_Reader *ReaderCallerSession) GetDefaultPoolConfiguration(dataStore common.Address) (*big.Int, error) {
	return _Reader.Contract.GetDefaultPoolConfiguration(&_Reader.CallOpts, dataStore)
}

// GetMarginLevelThreshold is a free data retrieval call binding the contract method 0x57b91ca6.
//
// Solidity: function getMarginLevelThreshold(address dataStore) view returns(uint256)
func (_Reader *ReaderCaller) GetMarginLevelThreshold(opts *bind.CallOpts, dataStore common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getMarginLevelThreshold", dataStore)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMarginLevelThreshold is a free data retrieval call binding the contract method 0x57b91ca6.
//
// Solidity: function getMarginLevelThreshold(address dataStore) view returns(uint256)
func (_Reader *ReaderSession) GetMarginLevelThreshold(dataStore common.Address) (*big.Int, error) {
	return _Reader.Contract.GetMarginLevelThreshold(&_Reader.CallOpts, dataStore)
}

// GetMarginLevelThreshold is a free data retrieval call binding the contract method 0x57b91ca6.
//
// Solidity: function getMarginLevelThreshold(address dataStore) view returns(uint256)
func (_Reader *ReaderCallerSession) GetMarginLevelThreshold(dataStore common.Address) (*big.Int, error) {
	return _Reader.Contract.GetMarginLevelThreshold(&_Reader.CallOpts, dataStore)
}

// GetPools is a free data retrieval call binding the contract method 0x5c39f467.
//
// Solidity: function getPools(address dataStore) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderCaller) GetPools(opts *bind.CallOpts, dataStore common.Address) ([]ReaderPoolUtilsGetPool, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPools", dataStore)

	if err != nil {
		return *new([]ReaderPoolUtilsGetPool), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPoolUtilsGetPool)).(*[]ReaderPoolUtilsGetPool)

	return out0, err

}

// GetPools is a free data retrieval call binding the contract method 0x5c39f467.
//
// Solidity: function getPools(address dataStore) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderSession) GetPools(dataStore common.Address) ([]ReaderPoolUtilsGetPool, error) {
	return _Reader.Contract.GetPools(&_Reader.CallOpts, dataStore)
}

// GetPools is a free data retrieval call binding the contract method 0x5c39f467.
//
// Solidity: function getPools(address dataStore) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderCallerSession) GetPools(dataStore common.Address) ([]ReaderPoolUtilsGetPool, error) {
	return _Reader.Contract.GetPools(&_Reader.CallOpts, dataStore)
}

// GetPools0 is a free data retrieval call binding the contract method 0x8f6c7a3c.
//
// Solidity: function getPools(address dataStore, uint256 start, uint256 end) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderCaller) GetPools0(opts *bind.CallOpts, dataStore common.Address, start *big.Int, end *big.Int) ([]ReaderPoolUtilsGetPool, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPools0", dataStore, start, end)

	if err != nil {
		return *new([]ReaderPoolUtilsGetPool), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPoolUtilsGetPool)).(*[]ReaderPoolUtilsGetPool)

	return out0, err

}

// GetPools0 is a free data retrieval call binding the contract method 0x8f6c7a3c.
//
// Solidity: function getPools(address dataStore, uint256 start, uint256 end) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderSession) GetPools0(dataStore common.Address, start *big.Int, end *big.Int) ([]ReaderPoolUtilsGetPool, error) {
	return _Reader.Contract.GetPools0(&_Reader.CallOpts, dataStore, start, end)
}

// GetPools0 is a free data retrieval call binding the contract method 0x8f6c7a3c.
//
// Solidity: function getPools(address dataStore, uint256 start, uint256 end) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderCallerSession) GetPools0(dataStore common.Address, start *big.Int, end *big.Int) ([]ReaderPoolUtilsGetPool, error) {
	return _Reader.Contract.GetPools0(&_Reader.CallOpts, dataStore, start, end)
}

// GetPools2 is a free data retrieval call binding the contract method 0x50376aaa.
//
// Solidity: function getPools2(address dataStore, bytes32[] poolKeys) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderCaller) GetPools2(opts *bind.CallOpts, dataStore common.Address, poolKeys [][32]byte) ([]ReaderPoolUtilsGetPool, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPools2", dataStore, poolKeys)

	if err != nil {
		return *new([]ReaderPoolUtilsGetPool), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPoolUtilsGetPool)).(*[]ReaderPoolUtilsGetPool)

	return out0, err

}

// GetPools2 is a free data retrieval call binding the contract method 0x50376aaa.
//
// Solidity: function getPools2(address dataStore, bytes32[] poolKeys) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderSession) GetPools2(dataStore common.Address, poolKeys [][32]byte) ([]ReaderPoolUtilsGetPool, error) {
	return _Reader.Contract.GetPools2(&_Reader.CallOpts, dataStore, poolKeys)
}

// GetPools2 is a free data retrieval call binding the contract method 0x50376aaa.
//
// Solidity: function getPools2(address dataStore, bytes32[] poolKeys) view returns(((address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256,uint256,uint256,address,bool,uint256,uint256,uint256,uint256)[])
func (_Reader *ReaderCallerSession) GetPools2(dataStore common.Address, poolKeys [][32]byte) ([]ReaderPoolUtilsGetPool, error) {
	return _Reader.Contract.GetPools2(&_Reader.CallOpts, dataStore, poolKeys)
}

// GetPoolsCount is a free data retrieval call binding the contract method 0x5a6f5710.
//
// Solidity: function getPoolsCount(address dataStore) view returns(uint256)
func (_Reader *ReaderCaller) GetPoolsCount(opts *bind.CallOpts, dataStore common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPoolsCount", dataStore)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPoolsCount is a free data retrieval call binding the contract method 0x5a6f5710.
//
// Solidity: function getPoolsCount(address dataStore) view returns(uint256)
func (_Reader *ReaderSession) GetPoolsCount(dataStore common.Address) (*big.Int, error) {
	return _Reader.Contract.GetPoolsCount(&_Reader.CallOpts, dataStore)
}

// GetPoolsCount is a free data retrieval call binding the contract method 0x5a6f5710.
//
// Solidity: function getPoolsCount(address dataStore) view returns(uint256)
func (_Reader *ReaderCallerSession) GetPoolsCount(dataStore common.Address) (*big.Int, error) {
	return _Reader.Contract.GetPoolsCount(&_Reader.CallOpts, dataStore)
}

// GetPoolsInfo is a free data retrieval call binding the contract method 0x1a308175.
//
// Solidity: function getPoolsInfo(address dataStore) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderCaller) GetPoolsInfo(opts *bind.CallOpts, dataStore common.Address) ([]ReaderPoolUtilsGetPoolInfo, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPoolsInfo", dataStore)

	if err != nil {
		return *new([]ReaderPoolUtilsGetPoolInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPoolUtilsGetPoolInfo)).(*[]ReaderPoolUtilsGetPoolInfo)

	return out0, err

}

// GetPoolsInfo is a free data retrieval call binding the contract method 0x1a308175.
//
// Solidity: function getPoolsInfo(address dataStore) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderSession) GetPoolsInfo(dataStore common.Address) ([]ReaderPoolUtilsGetPoolInfo, error) {
	return _Reader.Contract.GetPoolsInfo(&_Reader.CallOpts, dataStore)
}

// GetPoolsInfo is a free data retrieval call binding the contract method 0x1a308175.
//
// Solidity: function getPoolsInfo(address dataStore) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderCallerSession) GetPoolsInfo(dataStore common.Address) ([]ReaderPoolUtilsGetPoolInfo, error) {
	return _Reader.Contract.GetPoolsInfo(&_Reader.CallOpts, dataStore)
}

// GetPoolsInfo0 is a free data retrieval call binding the contract method 0x382fc72e.
//
// Solidity: function getPoolsInfo(address dataStore, uint256 start, uint256 end) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderCaller) GetPoolsInfo0(opts *bind.CallOpts, dataStore common.Address, start *big.Int, end *big.Int) ([]ReaderPoolUtilsGetPoolInfo, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPoolsInfo0", dataStore, start, end)

	if err != nil {
		return *new([]ReaderPoolUtilsGetPoolInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPoolUtilsGetPoolInfo)).(*[]ReaderPoolUtilsGetPoolInfo)

	return out0, err

}

// GetPoolsInfo0 is a free data retrieval call binding the contract method 0x382fc72e.
//
// Solidity: function getPoolsInfo(address dataStore, uint256 start, uint256 end) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderSession) GetPoolsInfo0(dataStore common.Address, start *big.Int, end *big.Int) ([]ReaderPoolUtilsGetPoolInfo, error) {
	return _Reader.Contract.GetPoolsInfo0(&_Reader.CallOpts, dataStore, start, end)
}

// GetPoolsInfo0 is a free data retrieval call binding the contract method 0x382fc72e.
//
// Solidity: function getPoolsInfo(address dataStore, uint256 start, uint256 end) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderCallerSession) GetPoolsInfo0(dataStore common.Address, start *big.Int, end *big.Int) ([]ReaderPoolUtilsGetPoolInfo, error) {
	return _Reader.Contract.GetPoolsInfo0(&_Reader.CallOpts, dataStore, start, end)
}

// GetPoolsInfo1 is a free data retrieval call binding the contract method 0xeed07428.
//
// Solidity: function getPoolsInfo(address dataStore, bytes32[] poolKeys) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderCaller) GetPoolsInfo1(opts *bind.CallOpts, dataStore common.Address, poolKeys [][32]byte) ([]ReaderPoolUtilsGetPoolInfo, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPoolsInfo1", dataStore, poolKeys)

	if err != nil {
		return *new([]ReaderPoolUtilsGetPoolInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPoolUtilsGetPoolInfo)).(*[]ReaderPoolUtilsGetPoolInfo)

	return out0, err

}

// GetPoolsInfo1 is a free data retrieval call binding the contract method 0xeed07428.
//
// Solidity: function getPoolsInfo(address dataStore, bytes32[] poolKeys) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderSession) GetPoolsInfo1(dataStore common.Address, poolKeys [][32]byte) ([]ReaderPoolUtilsGetPoolInfo, error) {
	return _Reader.Contract.GetPoolsInfo1(&_Reader.CallOpts, dataStore, poolKeys)
}

// GetPoolsInfo1 is a free data retrieval call binding the contract method 0xeed07428.
//
// Solidity: function getPoolsInfo(address dataStore, bytes32[] poolKeys) view returns(((address,string,uint256,uint256)[2],uint256,uint256)[])
func (_Reader *ReaderCallerSession) GetPoolsInfo1(dataStore common.Address, poolKeys [][32]byte) ([]ReaderPoolUtilsGetPoolInfo, error) {
	return _Reader.Contract.GetPoolsInfo1(&_Reader.CallOpts, dataStore, poolKeys)
}

// GetPositions is a free data retrieval call binding the contract method 0x739118a4.
//
// Solidity: function getPositions(address dataStore, address account) view returns(((address,string,uint256,uint256,uint256,int256,uint256,uint256,int256,int256)[2],uint256,address,uint256,uint256,uint256,int256,uint256,int256)[])
func (_Reader *ReaderCaller) GetPositions(opts *bind.CallOpts, dataStore common.Address, account common.Address) ([]ReaderPositionUtilsGetPosition, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPositions", dataStore, account)

	if err != nil {
		return *new([]ReaderPositionUtilsGetPosition), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPositionUtilsGetPosition)).(*[]ReaderPositionUtilsGetPosition)

	return out0, err

}

// GetPositions is a free data retrieval call binding the contract method 0x739118a4.
//
// Solidity: function getPositions(address dataStore, address account) view returns(((address,string,uint256,uint256,uint256,int256,uint256,uint256,int256,int256)[2],uint256,address,uint256,uint256,uint256,int256,uint256,int256)[])
func (_Reader *ReaderSession) GetPositions(dataStore common.Address, account common.Address) ([]ReaderPositionUtilsGetPosition, error) {
	return _Reader.Contract.GetPositions(&_Reader.CallOpts, dataStore, account)
}

// GetPositions is a free data retrieval call binding the contract method 0x739118a4.
//
// Solidity: function getPositions(address dataStore, address account) view returns(((address,string,uint256,uint256,uint256,int256,uint256,uint256,int256,int256)[2],uint256,address,uint256,uint256,uint256,int256,uint256,int256)[])
func (_Reader *ReaderCallerSession) GetPositions(dataStore common.Address, account common.Address) ([]ReaderPositionUtilsGetPosition, error) {
	return _Reader.Contract.GetPositions(&_Reader.CallOpts, dataStore, account)
}

// GetPositions2 is a free data retrieval call binding the contract method 0xa6a100a0.
//
// Solidity: function getPositions2(address dataStore, bytes32[] positionKeys) view returns(((address,string,uint256,uint256,uint256,int256,uint256,uint256,int256,int256)[2],uint256,address,uint256,uint256,uint256,int256,uint256,int256)[])
func (_Reader *ReaderCaller) GetPositions2(opts *bind.CallOpts, dataStore common.Address, positionKeys [][32]byte) ([]ReaderPositionUtilsGetPosition, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getPositions2", dataStore, positionKeys)

	if err != nil {
		return *new([]ReaderPositionUtilsGetPosition), err
	}

	out0 := *abi.ConvertType(out[0], new([]ReaderPositionUtilsGetPosition)).(*[]ReaderPositionUtilsGetPosition)

	return out0, err

}

// GetPositions2 is a free data retrieval call binding the contract method 0xa6a100a0.
//
// Solidity: function getPositions2(address dataStore, bytes32[] positionKeys) view returns(((address,string,uint256,uint256,uint256,int256,uint256,uint256,int256,int256)[2],uint256,address,uint256,uint256,uint256,int256,uint256,int256)[])
func (_Reader *ReaderSession) GetPositions2(dataStore common.Address, positionKeys [][32]byte) ([]ReaderPositionUtilsGetPosition, error) {
	return _Reader.Contract.GetPositions2(&_Reader.CallOpts, dataStore, positionKeys)
}

// GetPositions2 is a free data retrieval call binding the contract method 0xa6a100a0.
//
// Solidity: function getPositions2(address dataStore, bytes32[] positionKeys) view returns(((address,string,uint256,uint256,uint256,int256,uint256,uint256,int256,int256)[2],uint256,address,uint256,uint256,uint256,int256,uint256,int256)[])
func (_Reader *ReaderCallerSession) GetPositions2(dataStore common.Address, positionKeys [][32]byte) ([]ReaderPositionUtilsGetPosition, error) {
	return _Reader.Contract.GetPositions2(&_Reader.CallOpts, dataStore, positionKeys)
}

// GetTokenBase is a free data retrieval call binding the contract method 0x28a0ccf4.
//
// Solidity: function getTokenBase(address dataStore) view returns(address)
func (_Reader *ReaderCaller) GetTokenBase(opts *bind.CallOpts, dataStore common.Address) (common.Address, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getTokenBase", dataStore)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenBase is a free data retrieval call binding the contract method 0x28a0ccf4.
//
// Solidity: function getTokenBase(address dataStore) view returns(address)
func (_Reader *ReaderSession) GetTokenBase(dataStore common.Address) (common.Address, error) {
	return _Reader.Contract.GetTokenBase(&_Reader.CallOpts, dataStore)
}

// GetTokenBase is a free data retrieval call binding the contract method 0x28a0ccf4.
//
// Solidity: function getTokenBase(address dataStore) view returns(address)
func (_Reader *ReaderCallerSession) GetTokenBase(dataStore common.Address) (common.Address, error) {
	return _Reader.Contract.GetTokenBase(&_Reader.CallOpts, dataStore)
}

// GetTreasury is a free data retrieval call binding the contract method 0x78f212d1.
//
// Solidity: function getTreasury(address dataStore) view returns(address)
func (_Reader *ReaderCaller) GetTreasury(opts *bind.CallOpts, dataStore common.Address) (common.Address, error) {
	var out []interface{}
	err := _Reader.contract.Call(opts, &out, "getTreasury", dataStore)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTreasury is a free data retrieval call binding the contract method 0x78f212d1.
//
// Solidity: function getTreasury(address dataStore) view returns(address)
func (_Reader *ReaderSession) GetTreasury(dataStore common.Address) (common.Address, error) {
	return _Reader.Contract.GetTreasury(&_Reader.CallOpts, dataStore)
}

// GetTreasury is a free data retrieval call binding the contract method 0x78f212d1.
//
// Solidity: function getTreasury(address dataStore) view returns(address)
func (_Reader *ReaderCallerSession) GetTreasury(dataStore common.Address) (common.Address, error) {
	return _Reader.Contract.GetTreasury(&_Reader.CallOpts, dataStore)
}
