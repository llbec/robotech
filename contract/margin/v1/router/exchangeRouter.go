// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Router

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

// BorrowUtilsBorrowParams is an auto generated low-level Go binding around an user-defined struct.
type BorrowUtilsBorrowParams struct {
	PositionId   *big.Int
	TokenIndex   uint8
	BorrowAmount *big.Int
}

// CloseUtilsCloseParams is an auto generated low-level Go binding around an user-defined struct.
type CloseUtilsCloseParams struct {
	PositionId *big.Int
}

// DepositUtilsDepositParams is an auto generated low-level Go binding around an user-defined struct.
type DepositUtilsDepositParams struct {
	PositionId *big.Int
	Token0     common.Address
	Token1     common.Address
	TokenIndex uint8
}

// LiquidationUtilsLiquidationParams is an auto generated low-level Go binding around an user-defined struct.
type LiquidationUtilsLiquidationParams struct {
	Account    common.Address
	PositionId *big.Int
}

// LiquidityUtilsAddParams is an auto generated low-level Go binding around an user-defined struct.
type LiquidityUtilsAddParams struct {
	Token0 common.Address
	Token1 common.Address
	To     common.Address
}

// LiquidityUtilsRemoveParams is an auto generated low-level Go binding around an user-defined struct.
type LiquidityUtilsRemoveParams struct {
	Token0    common.Address
	Token1    common.Address
	Liquidity *big.Int
	To        common.Address
}

// RepayUtilsRepayParams is an auto generated low-level Go binding around an user-defined struct.
type RepayUtilsRepayParams struct {
	PositionId  *big.Int
	TokenIndex  uint8
	RepayAmount *big.Int
	RepayOption uint8
}

// SwapUtilsSwapInPositionParams is an auto generated low-level Go binding around an user-defined struct.
type SwapUtilsSwapInPositionParams struct {
	PositionId *big.Int
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
}

// SwapUtilsSwapParams is an auto generated low-level Go binding around an user-defined struct.
type SwapUtilsSwapParams struct {
	Token0     common.Address
	Token1     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
}

// WithdrawUtilsWithdrawParams is an auto generated low-level Go binding around an user-defined struct.
type WithdrawUtilsWithdrawParams struct {
	PositionId     *big.Int
	TokenIndex     uint8
	WithdrawAmount *big.Int
	To             common.Address
}

// RouterMetaData contains all meta data concerning the Router contract.
var RouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractRouter\",\"name\":\"_router\",\"type\":\"address\"},{\"internalType\":\"contractRoleStore\",\"name\":\"_roleStore\",\"type\":\"address\"},{\"internalType\":\"contractDataStore\",\"name\":\"_dataStore\",\"type\":\"address\"},{\"internalType\":\"contractILiquidityHandler\",\"name\":\"_liquidityHandler\",\"type\":\"address\"},{\"internalType\":\"contractISwapHandler\",\"name\":\"_swapHandler\",\"type\":\"address\"},{\"internalType\":\"contractIDepositHandler\",\"name\":\"_depositHandler\",\"type\":\"address\"},{\"internalType\":\"contractIBorrowHandler\",\"name\":\"_borrowHandler\",\"type\":\"address\"},{\"internalType\":\"contractIRepayHandler\",\"name\":\"_repayHandler\",\"type\":\"address\"},{\"internalType\":\"contractIWithdrawHandler\",\"name\":\"_withdrawHandler\",\"type\":\"address\"},{\"internalType\":\"contractILiquidationHandler\",\"name\":\"_liquidationHandler\",\"type\":\"address\"},{\"internalType\":\"contractICloseHandler\",\"name\":\"_closeHandler\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyReceiver\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"role\",\"type\":\"string\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"borrowHandler\",\"outputs\":[{\"internalType\":\"contractIBorrowHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"closeHandler\",\"outputs\":[{\"internalType\":\"contractICloseHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dataStore\",\"outputs\":[{\"internalType\":\"contractDataStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositHandler\",\"outputs\":[{\"internalType\":\"contractIDepositHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"internalType\":\"structLiquidityUtils.AddParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"borrowAmount\",\"type\":\"uint256\"}],\"internalType\":\"structBorrowUtils.BorrowParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeBorrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"}],\"internalType\":\"structCloseUtils.CloseParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeClose\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"}],\"internalType\":\"structDepositUtils.DepositParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"}],\"internalType\":\"structLiquidationUtils.LiquidationParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeLiquidation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"}],\"internalType\":\"structLiquidationUtils.LiquidationParams[]\",\"name\":\"params\",\"type\":\"tuple[]\"}],\"name\":\"executeLiquidationBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"internalType\":\"structLiquidityUtils.RemoveParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"repayAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"repayOption\",\"type\":\"uint8\"}],\"internalType\":\"structRepayUtils.RepayParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeRepay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"internalType\":\"structSwapUtils.SwapParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeSwap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"}],\"internalType\":\"structSwapUtils.SwapInPositionParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeSwapInPosition\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"internalType\":\"structWithdrawUtils.WithdrawParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"executeWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidationHandler\",\"outputs\":[{\"internalType\":\"contractILiquidationHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidityHandler\",\"outputs\":[{\"internalType\":\"contractILiquidityHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"repayHandler\",\"outputs\":[{\"internalType\":\"contractIRepayHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleStore\",\"outputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"contractRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendTokens\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapHandler\",\"outputs\":[{\"internalType\":\"contractISwapHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawHandler\",\"outputs\":[{\"internalType\":\"contractIWithdrawHandler\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RouterABI is the input ABI used to generate the binding from.
// Deprecated: Use RouterMetaData.ABI instead.
var RouterABI = RouterMetaData.ABI

// Router is an auto generated Go binding around an Ethereum contract.
type Router struct {
	RouterCaller     // Read-only binding to the contract
	RouterTransactor // Write-only binding to the contract
	RouterFilterer   // Log filterer for contract events
}

// RouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type RouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RouterSession struct {
	Contract     *Router           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RouterCallerSession struct {
	Contract *RouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RouterTransactorSession struct {
	Contract     *RouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type RouterRaw struct {
	Contract *Router // Generic contract binding to access the raw methods on
}

// RouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RouterCallerRaw struct {
	Contract *RouterCaller // Generic read-only contract binding to access the raw methods on
}

// RouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RouterTransactorRaw struct {
	Contract *RouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRouter creates a new instance of Router, bound to a specific deployed contract.
func NewRouter(address common.Address, backend bind.ContractBackend) (*Router, error) {
	contract, err := bindRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Router{RouterCaller: RouterCaller{contract: contract}, RouterTransactor: RouterTransactor{contract: contract}, RouterFilterer: RouterFilterer{contract: contract}}, nil
}

// NewRouterCaller creates a new read-only instance of Router, bound to a specific deployed contract.
func NewRouterCaller(address common.Address, caller bind.ContractCaller) (*RouterCaller, error) {
	contract, err := bindRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RouterCaller{contract: contract}, nil
}

// NewRouterTransactor creates a new write-only instance of Router, bound to a specific deployed contract.
func NewRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*RouterTransactor, error) {
	contract, err := bindRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RouterTransactor{contract: contract}, nil
}

// NewRouterFilterer creates a new log filterer instance of Router, bound to a specific deployed contract.
func NewRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*RouterFilterer, error) {
	contract, err := bindRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RouterFilterer{contract: contract}, nil
}

// bindRouter binds a generic wrapper to an already deployed contract.
func bindRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Router *RouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Router.Contract.RouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Router *RouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Router.Contract.RouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Router *RouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Router.Contract.RouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Router *RouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Router.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Router *RouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Router.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Router *RouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Router.Contract.contract.Transact(opts, method, params...)
}

// BorrowHandler is a free data retrieval call binding the contract method 0xa82ed4ce.
//
// Solidity: function borrowHandler() view returns(address)
func (_Router *RouterCaller) BorrowHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "borrowHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BorrowHandler is a free data retrieval call binding the contract method 0xa82ed4ce.
//
// Solidity: function borrowHandler() view returns(address)
func (_Router *RouterSession) BorrowHandler() (common.Address, error) {
	return _Router.Contract.BorrowHandler(&_Router.CallOpts)
}

// BorrowHandler is a free data retrieval call binding the contract method 0xa82ed4ce.
//
// Solidity: function borrowHandler() view returns(address)
func (_Router *RouterCallerSession) BorrowHandler() (common.Address, error) {
	return _Router.Contract.BorrowHandler(&_Router.CallOpts)
}

// CloseHandler is a free data retrieval call binding the contract method 0xed27afaf.
//
// Solidity: function closeHandler() view returns(address)
func (_Router *RouterCaller) CloseHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "closeHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CloseHandler is a free data retrieval call binding the contract method 0xed27afaf.
//
// Solidity: function closeHandler() view returns(address)
func (_Router *RouterSession) CloseHandler() (common.Address, error) {
	return _Router.Contract.CloseHandler(&_Router.CallOpts)
}

// CloseHandler is a free data retrieval call binding the contract method 0xed27afaf.
//
// Solidity: function closeHandler() view returns(address)
func (_Router *RouterCallerSession) CloseHandler() (common.Address, error) {
	return _Router.Contract.CloseHandler(&_Router.CallOpts)
}

// DataStore is a free data retrieval call binding the contract method 0x660d0d67.
//
// Solidity: function dataStore() view returns(address)
func (_Router *RouterCaller) DataStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "dataStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DataStore is a free data retrieval call binding the contract method 0x660d0d67.
//
// Solidity: function dataStore() view returns(address)
func (_Router *RouterSession) DataStore() (common.Address, error) {
	return _Router.Contract.DataStore(&_Router.CallOpts)
}

// DataStore is a free data retrieval call binding the contract method 0x660d0d67.
//
// Solidity: function dataStore() view returns(address)
func (_Router *RouterCallerSession) DataStore() (common.Address, error) {
	return _Router.Contract.DataStore(&_Router.CallOpts)
}

// DepositHandler is a free data retrieval call binding the contract method 0x9c8b2cfb.
//
// Solidity: function depositHandler() view returns(address)
func (_Router *RouterCaller) DepositHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "depositHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DepositHandler is a free data retrieval call binding the contract method 0x9c8b2cfb.
//
// Solidity: function depositHandler() view returns(address)
func (_Router *RouterSession) DepositHandler() (common.Address, error) {
	return _Router.Contract.DepositHandler(&_Router.CallOpts)
}

// DepositHandler is a free data retrieval call binding the contract method 0x9c8b2cfb.
//
// Solidity: function depositHandler() view returns(address)
func (_Router *RouterCallerSession) DepositHandler() (common.Address, error) {
	return _Router.Contract.DepositHandler(&_Router.CallOpts)
}

// LiquidationHandler is a free data retrieval call binding the contract method 0xd25adeb3.
//
// Solidity: function liquidationHandler() view returns(address)
func (_Router *RouterCaller) LiquidationHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "liquidationHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LiquidationHandler is a free data retrieval call binding the contract method 0xd25adeb3.
//
// Solidity: function liquidationHandler() view returns(address)
func (_Router *RouterSession) LiquidationHandler() (common.Address, error) {
	return _Router.Contract.LiquidationHandler(&_Router.CallOpts)
}

// LiquidationHandler is a free data retrieval call binding the contract method 0xd25adeb3.
//
// Solidity: function liquidationHandler() view returns(address)
func (_Router *RouterCallerSession) LiquidationHandler() (common.Address, error) {
	return _Router.Contract.LiquidationHandler(&_Router.CallOpts)
}

// LiquidityHandler is a free data retrieval call binding the contract method 0xd9e08811.
//
// Solidity: function liquidityHandler() view returns(address)
func (_Router *RouterCaller) LiquidityHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "liquidityHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LiquidityHandler is a free data retrieval call binding the contract method 0xd9e08811.
//
// Solidity: function liquidityHandler() view returns(address)
func (_Router *RouterSession) LiquidityHandler() (common.Address, error) {
	return _Router.Contract.LiquidityHandler(&_Router.CallOpts)
}

// LiquidityHandler is a free data retrieval call binding the contract method 0xd9e08811.
//
// Solidity: function liquidityHandler() view returns(address)
func (_Router *RouterCallerSession) LiquidityHandler() (common.Address, error) {
	return _Router.Contract.LiquidityHandler(&_Router.CallOpts)
}

// RepayHandler is a free data retrieval call binding the contract method 0x9d451d0c.
//
// Solidity: function repayHandler() view returns(address)
func (_Router *RouterCaller) RepayHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "repayHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RepayHandler is a free data retrieval call binding the contract method 0x9d451d0c.
//
// Solidity: function repayHandler() view returns(address)
func (_Router *RouterSession) RepayHandler() (common.Address, error) {
	return _Router.Contract.RepayHandler(&_Router.CallOpts)
}

// RepayHandler is a free data retrieval call binding the contract method 0x9d451d0c.
//
// Solidity: function repayHandler() view returns(address)
func (_Router *RouterCallerSession) RepayHandler() (common.Address, error) {
	return _Router.Contract.RepayHandler(&_Router.CallOpts)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_Router *RouterCaller) RoleStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "roleStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_Router *RouterSession) RoleStore() (common.Address, error) {
	return _Router.Contract.RoleStore(&_Router.CallOpts)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_Router *RouterCallerSession) RoleStore() (common.Address, error) {
	return _Router.Contract.RoleStore(&_Router.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_Router *RouterCaller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_Router *RouterSession) Router() (common.Address, error) {
	return _Router.Contract.Router(&_Router.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_Router *RouterCallerSession) Router() (common.Address, error) {
	return _Router.Contract.Router(&_Router.CallOpts)
}

// SwapHandler is a free data retrieval call binding the contract method 0x8a53aaac.
//
// Solidity: function swapHandler() view returns(address)
func (_Router *RouterCaller) SwapHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "swapHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapHandler is a free data retrieval call binding the contract method 0x8a53aaac.
//
// Solidity: function swapHandler() view returns(address)
func (_Router *RouterSession) SwapHandler() (common.Address, error) {
	return _Router.Contract.SwapHandler(&_Router.CallOpts)
}

// SwapHandler is a free data retrieval call binding the contract method 0x8a53aaac.
//
// Solidity: function swapHandler() view returns(address)
func (_Router *RouterCallerSession) SwapHandler() (common.Address, error) {
	return _Router.Contract.SwapHandler(&_Router.CallOpts)
}

// WithdrawHandler is a free data retrieval call binding the contract method 0x083473ef.
//
// Solidity: function withdrawHandler() view returns(address)
func (_Router *RouterCaller) WithdrawHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "withdrawHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WithdrawHandler is a free data retrieval call binding the contract method 0x083473ef.
//
// Solidity: function withdrawHandler() view returns(address)
func (_Router *RouterSession) WithdrawHandler() (common.Address, error) {
	return _Router.Contract.WithdrawHandler(&_Router.CallOpts)
}

// WithdrawHandler is a free data retrieval call binding the contract method 0x083473ef.
//
// Solidity: function withdrawHandler() view returns(address)
func (_Router *RouterCallerSession) WithdrawHandler() (common.Address, error) {
	return _Router.Contract.WithdrawHandler(&_Router.CallOpts)
}

// ExecuteAdd is a paid mutator transaction binding the contract method 0x21b108cc.
//
// Solidity: function executeAdd((address,address,address) params) returns()
func (_Router *RouterTransactor) ExecuteAdd(opts *bind.TransactOpts, params LiquidityUtilsAddParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeAdd", params)
}

// ExecuteAdd is a paid mutator transaction binding the contract method 0x21b108cc.
//
// Solidity: function executeAdd((address,address,address) params) returns()
func (_Router *RouterSession) ExecuteAdd(params LiquidityUtilsAddParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteAdd(&_Router.TransactOpts, params)
}

// ExecuteAdd is a paid mutator transaction binding the contract method 0x21b108cc.
//
// Solidity: function executeAdd((address,address,address) params) returns()
func (_Router *RouterTransactorSession) ExecuteAdd(params LiquidityUtilsAddParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteAdd(&_Router.TransactOpts, params)
}

// ExecuteBorrow is a paid mutator transaction binding the contract method 0x907e7454.
//
// Solidity: function executeBorrow((uint256,uint8,uint256) params) returns()
func (_Router *RouterTransactor) ExecuteBorrow(opts *bind.TransactOpts, params BorrowUtilsBorrowParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeBorrow", params)
}

// ExecuteBorrow is a paid mutator transaction binding the contract method 0x907e7454.
//
// Solidity: function executeBorrow((uint256,uint8,uint256) params) returns()
func (_Router *RouterSession) ExecuteBorrow(params BorrowUtilsBorrowParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteBorrow(&_Router.TransactOpts, params)
}

// ExecuteBorrow is a paid mutator transaction binding the contract method 0x907e7454.
//
// Solidity: function executeBorrow((uint256,uint8,uint256) params) returns()
func (_Router *RouterTransactorSession) ExecuteBorrow(params BorrowUtilsBorrowParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteBorrow(&_Router.TransactOpts, params)
}

// ExecuteClose is a paid mutator transaction binding the contract method 0xc56d4a2a.
//
// Solidity: function executeClose((uint256) params) returns()
func (_Router *RouterTransactor) ExecuteClose(opts *bind.TransactOpts, params CloseUtilsCloseParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeClose", params)
}

// ExecuteClose is a paid mutator transaction binding the contract method 0xc56d4a2a.
//
// Solidity: function executeClose((uint256) params) returns()
func (_Router *RouterSession) ExecuteClose(params CloseUtilsCloseParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteClose(&_Router.TransactOpts, params)
}

// ExecuteClose is a paid mutator transaction binding the contract method 0xc56d4a2a.
//
// Solidity: function executeClose((uint256) params) returns()
func (_Router *RouterTransactorSession) ExecuteClose(params CloseUtilsCloseParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteClose(&_Router.TransactOpts, params)
}

// ExecuteDeposit is a paid mutator transaction binding the contract method 0x47570cdc.
//
// Solidity: function executeDeposit((uint256,address,address,uint8) params) returns()
func (_Router *RouterTransactor) ExecuteDeposit(opts *bind.TransactOpts, params DepositUtilsDepositParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeDeposit", params)
}

// ExecuteDeposit is a paid mutator transaction binding the contract method 0x47570cdc.
//
// Solidity: function executeDeposit((uint256,address,address,uint8) params) returns()
func (_Router *RouterSession) ExecuteDeposit(params DepositUtilsDepositParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteDeposit(&_Router.TransactOpts, params)
}

// ExecuteDeposit is a paid mutator transaction binding the contract method 0x47570cdc.
//
// Solidity: function executeDeposit((uint256,address,address,uint8) params) returns()
func (_Router *RouterTransactorSession) ExecuteDeposit(params DepositUtilsDepositParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteDeposit(&_Router.TransactOpts, params)
}

// ExecuteLiquidation is a paid mutator transaction binding the contract method 0xb6103215.
//
// Solidity: function executeLiquidation((address,uint256) params) returns()
func (_Router *RouterTransactor) ExecuteLiquidation(opts *bind.TransactOpts, params LiquidationUtilsLiquidationParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeLiquidation", params)
}

// ExecuteLiquidation is a paid mutator transaction binding the contract method 0xb6103215.
//
// Solidity: function executeLiquidation((address,uint256) params) returns()
func (_Router *RouterSession) ExecuteLiquidation(params LiquidationUtilsLiquidationParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteLiquidation(&_Router.TransactOpts, params)
}

// ExecuteLiquidation is a paid mutator transaction binding the contract method 0xb6103215.
//
// Solidity: function executeLiquidation((address,uint256) params) returns()
func (_Router *RouterTransactorSession) ExecuteLiquidation(params LiquidationUtilsLiquidationParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteLiquidation(&_Router.TransactOpts, params)
}

// ExecuteLiquidationBatch is a paid mutator transaction binding the contract method 0xa76718d2.
//
// Solidity: function executeLiquidationBatch((address,uint256)[] params) returns()
func (_Router *RouterTransactor) ExecuteLiquidationBatch(opts *bind.TransactOpts, params []LiquidationUtilsLiquidationParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeLiquidationBatch", params)
}

// ExecuteLiquidationBatch is a paid mutator transaction binding the contract method 0xa76718d2.
//
// Solidity: function executeLiquidationBatch((address,uint256)[] params) returns()
func (_Router *RouterSession) ExecuteLiquidationBatch(params []LiquidationUtilsLiquidationParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteLiquidationBatch(&_Router.TransactOpts, params)
}

// ExecuteLiquidationBatch is a paid mutator transaction binding the contract method 0xa76718d2.
//
// Solidity: function executeLiquidationBatch((address,uint256)[] params) returns()
func (_Router *RouterTransactorSession) ExecuteLiquidationBatch(params []LiquidationUtilsLiquidationParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteLiquidationBatch(&_Router.TransactOpts, params)
}

// ExecuteRemove is a paid mutator transaction binding the contract method 0xf3e01a4e.
//
// Solidity: function executeRemove((address,address,uint256,address) params) returns()
func (_Router *RouterTransactor) ExecuteRemove(opts *bind.TransactOpts, params LiquidityUtilsRemoveParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeRemove", params)
}

// ExecuteRemove is a paid mutator transaction binding the contract method 0xf3e01a4e.
//
// Solidity: function executeRemove((address,address,uint256,address) params) returns()
func (_Router *RouterSession) ExecuteRemove(params LiquidityUtilsRemoveParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteRemove(&_Router.TransactOpts, params)
}

// ExecuteRemove is a paid mutator transaction binding the contract method 0xf3e01a4e.
//
// Solidity: function executeRemove((address,address,uint256,address) params) returns()
func (_Router *RouterTransactorSession) ExecuteRemove(params LiquidityUtilsRemoveParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteRemove(&_Router.TransactOpts, params)
}

// ExecuteRepay is a paid mutator transaction binding the contract method 0x7b123bc8.
//
// Solidity: function executeRepay((uint256,uint8,uint256,uint8) params) returns()
func (_Router *RouterTransactor) ExecuteRepay(opts *bind.TransactOpts, params RepayUtilsRepayParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeRepay", params)
}

// ExecuteRepay is a paid mutator transaction binding the contract method 0x7b123bc8.
//
// Solidity: function executeRepay((uint256,uint8,uint256,uint8) params) returns()
func (_Router *RouterSession) ExecuteRepay(params RepayUtilsRepayParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteRepay(&_Router.TransactOpts, params)
}

// ExecuteRepay is a paid mutator transaction binding the contract method 0x7b123bc8.
//
// Solidity: function executeRepay((uint256,uint8,uint256,uint8) params) returns()
func (_Router *RouterTransactorSession) ExecuteRepay(params RepayUtilsRepayParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteRepay(&_Router.TransactOpts, params)
}

// ExecuteSwap is a paid mutator transaction binding the contract method 0x34f5db7e.
//
// Solidity: function executeSwap((address,address,uint256,uint256,uint256,uint256,address) params) returns()
func (_Router *RouterTransactor) ExecuteSwap(opts *bind.TransactOpts, params SwapUtilsSwapParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeSwap", params)
}

// ExecuteSwap is a paid mutator transaction binding the contract method 0x34f5db7e.
//
// Solidity: function executeSwap((address,address,uint256,uint256,uint256,uint256,address) params) returns()
func (_Router *RouterSession) ExecuteSwap(params SwapUtilsSwapParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteSwap(&_Router.TransactOpts, params)
}

// ExecuteSwap is a paid mutator transaction binding the contract method 0x34f5db7e.
//
// Solidity: function executeSwap((address,address,uint256,uint256,uint256,uint256,address) params) returns()
func (_Router *RouterTransactorSession) ExecuteSwap(params SwapUtilsSwapParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteSwap(&_Router.TransactOpts, params)
}

// ExecuteSwapInPosition is a paid mutator transaction binding the contract method 0xb1840675.
//
// Solidity: function executeSwapInPosition((uint256,uint256,uint256,uint256,uint256) params) returns()
func (_Router *RouterTransactor) ExecuteSwapInPosition(opts *bind.TransactOpts, params SwapUtilsSwapInPositionParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeSwapInPosition", params)
}

// ExecuteSwapInPosition is a paid mutator transaction binding the contract method 0xb1840675.
//
// Solidity: function executeSwapInPosition((uint256,uint256,uint256,uint256,uint256) params) returns()
func (_Router *RouterSession) ExecuteSwapInPosition(params SwapUtilsSwapInPositionParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteSwapInPosition(&_Router.TransactOpts, params)
}

// ExecuteSwapInPosition is a paid mutator transaction binding the contract method 0xb1840675.
//
// Solidity: function executeSwapInPosition((uint256,uint256,uint256,uint256,uint256) params) returns()
func (_Router *RouterTransactorSession) ExecuteSwapInPosition(params SwapUtilsSwapInPositionParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteSwapInPosition(&_Router.TransactOpts, params)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0xad0c9934.
//
// Solidity: function executeWithdraw((uint256,uint8,uint256,address) params) returns()
func (_Router *RouterTransactor) ExecuteWithdraw(opts *bind.TransactOpts, params WithdrawUtilsWithdrawParams) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "executeWithdraw", params)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0xad0c9934.
//
// Solidity: function executeWithdraw((uint256,uint8,uint256,address) params) returns()
func (_Router *RouterSession) ExecuteWithdraw(params WithdrawUtilsWithdrawParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteWithdraw(&_Router.TransactOpts, params)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0xad0c9934.
//
// Solidity: function executeWithdraw((uint256,uint8,uint256,address) params) returns()
func (_Router *RouterTransactorSession) ExecuteWithdraw(params WithdrawUtilsWithdrawParams) (*types.Transaction, error) {
	return _Router.Contract.ExecuteWithdraw(&_Router.TransactOpts, params)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_Router *RouterTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_Router *RouterSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Router.Contract.Multicall(&_Router.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_Router *RouterTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Router.Contract.Multicall(&_Router.TransactOpts, data)
}

// SendTokens is a paid mutator transaction binding the contract method 0xe6d66ac8.
//
// Solidity: function sendTokens(address token, address receiver, uint256 amount) payable returns()
func (_Router *RouterTransactor) SendTokens(opts *bind.TransactOpts, token common.Address, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "sendTokens", token, receiver, amount)
}

// SendTokens is a paid mutator transaction binding the contract method 0xe6d66ac8.
//
// Solidity: function sendTokens(address token, address receiver, uint256 amount) payable returns()
func (_Router *RouterSession) SendTokens(token common.Address, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Router.Contract.SendTokens(&_Router.TransactOpts, token, receiver, amount)
}

// SendTokens is a paid mutator transaction binding the contract method 0xe6d66ac8.
//
// Solidity: function sendTokens(address token, address receiver, uint256 amount) payable returns()
func (_Router *RouterTransactorSession) SendTokens(token common.Address, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Router.Contract.SendTokens(&_Router.TransactOpts, token, receiver, amount)
}
