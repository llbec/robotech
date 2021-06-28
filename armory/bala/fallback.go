// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bala

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FallBackABI is the input ABI used to generate the binding from.
const FallBackABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getAssetPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"setAssetPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// FallBack is an auto generated Go binding around an Ethereum contract.
type FallBack struct {
	FallBackCaller     // Read-only binding to the contract
	FallBackTransactor // Write-only binding to the contract
	FallBackFilterer   // Log filterer for contract events
}

// FallBackCaller is an auto generated read-only Go binding around an Ethereum contract.
type FallBackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FallBackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FallBackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FallBackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FallBackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FallBackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FallBackSession struct {
	Contract     *FallBack         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FallBackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FallBackCallerSession struct {
	Contract *FallBackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// FallBackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FallBackTransactorSession struct {
	Contract     *FallBackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FallBackRaw is an auto generated low-level Go binding around an Ethereum contract.
type FallBackRaw struct {
	Contract *FallBack // Generic contract binding to access the raw methods on
}

// FallBackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FallBackCallerRaw struct {
	Contract *FallBackCaller // Generic read-only contract binding to access the raw methods on
}

// FallBackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FallBackTransactorRaw struct {
	Contract *FallBackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFallBack creates a new instance of FallBack, bound to a specific deployed contract.
func NewFallBack(address common.Address, backend bind.ContractBackend) (*FallBack, error) {
	contract, err := bindFallBack(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FallBack{FallBackCaller: FallBackCaller{contract: contract}, FallBackTransactor: FallBackTransactor{contract: contract}, FallBackFilterer: FallBackFilterer{contract: contract}}, nil
}

// NewFallBackCaller creates a new read-only instance of FallBack, bound to a specific deployed contract.
func NewFallBackCaller(address common.Address, caller bind.ContractCaller) (*FallBackCaller, error) {
	contract, err := bindFallBack(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FallBackCaller{contract: contract}, nil
}

// NewFallBackTransactor creates a new write-only instance of FallBack, bound to a specific deployed contract.
func NewFallBackTransactor(address common.Address, transactor bind.ContractTransactor) (*FallBackTransactor, error) {
	contract, err := bindFallBack(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FallBackTransactor{contract: contract}, nil
}

// NewFallBackFilterer creates a new log filterer instance of FallBack, bound to a specific deployed contract.
func NewFallBackFilterer(address common.Address, filterer bind.ContractFilterer) (*FallBackFilterer, error) {
	contract, err := bindFallBack(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FallBackFilterer{contract: contract}, nil
}

// bindFallBack binds a generic wrapper to an already deployed contract.
func bindFallBack(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FallBackABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FallBack *FallBackRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FallBack.Contract.FallBackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FallBack *FallBackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FallBack.Contract.FallBackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FallBack *FallBackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FallBack.Contract.FallBackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FallBack *FallBackCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FallBack.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FallBack *FallBackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FallBack.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FallBack *FallBackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FallBack.Contract.contract.Transact(opts, method, params...)
}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address asset) view returns(uint256)
func (_FallBack *FallBackCaller) GetAssetPrice(opts *bind.CallOpts, asset common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FallBack.contract.Call(opts, &out, "getAssetPrice", asset)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address asset) view returns(uint256)
func (_FallBack *FallBackSession) GetAssetPrice(asset common.Address) (*big.Int, error) {
	return _FallBack.Contract.GetAssetPrice(&_FallBack.CallOpts, asset)
}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address asset) view returns(uint256)
func (_FallBack *FallBackCallerSession) GetAssetPrice(asset common.Address) (*big.Int, error) {
	return _FallBack.Contract.GetAssetPrice(&_FallBack.CallOpts, asset)
}

// SetAssetPrice is a paid mutator transaction binding the contract method 0x51323f72.
//
// Solidity: function setAssetPrice(address asset, uint256 price) returns()
func (_FallBack *FallBackTransactor) SetAssetPrice(opts *bind.TransactOpts, asset common.Address, price *big.Int) (*types.Transaction, error) {
	return _FallBack.contract.Transact(opts, "setAssetPrice", asset, price)
}

// SetAssetPrice is a paid mutator transaction binding the contract method 0x51323f72.
//
// Solidity: function setAssetPrice(address asset, uint256 price) returns()
func (_FallBack *FallBackSession) SetAssetPrice(asset common.Address, price *big.Int) (*types.Transaction, error) {
	return _FallBack.Contract.SetAssetPrice(&_FallBack.TransactOpts, asset, price)
}

// SetAssetPrice is a paid mutator transaction binding the contract method 0x51323f72.
//
// Solidity: function setAssetPrice(address asset, uint256 price) returns()
func (_FallBack *FallBackTransactorSession) SetAssetPrice(asset common.Address, price *big.Int) (*types.Transaction, error) {
	return _FallBack.Contract.SetAssetPrice(&_FallBack.TransactOpts, asset, price)
}
