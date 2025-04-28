// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package backendCheck

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

// BackendCheckMetaData contains all meta data concerning the BackendCheck contract.
var BackendCheckMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"getRemoveHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"badDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"list\",\"type\":\"bytes32[]\"}],\"name\":\"getSwapHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_msgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_msgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"verifySignature\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// BackendCheckABI is the input ABI used to generate the binding from.
// Deprecated: Use BackendCheckMetaData.ABI instead.
var BackendCheckABI = BackendCheckMetaData.ABI

// BackendCheck is an auto generated Go binding around an Ethereum contract.
type BackendCheck struct {
	BackendCheckCaller     // Read-only binding to the contract
	BackendCheckTransactor // Write-only binding to the contract
	BackendCheckFilterer   // Log filterer for contract events
}

// BackendCheckCaller is an auto generated read-only Go binding around an Ethereum contract.
type BackendCheckCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BackendCheckTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BackendCheckTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BackendCheckFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BackendCheckFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BackendCheckSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BackendCheckSession struct {
	Contract     *BackendCheck     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BackendCheckCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BackendCheckCallerSession struct {
	Contract *BackendCheckCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BackendCheckTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BackendCheckTransactorSession struct {
	Contract     *BackendCheckTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BackendCheckRaw is an auto generated low-level Go binding around an Ethereum contract.
type BackendCheckRaw struct {
	Contract *BackendCheck // Generic contract binding to access the raw methods on
}

// BackendCheckCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BackendCheckCallerRaw struct {
	Contract *BackendCheckCaller // Generic read-only contract binding to access the raw methods on
}

// BackendCheckTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BackendCheckTransactorRaw struct {
	Contract *BackendCheckTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBackendCheck creates a new instance of BackendCheck, bound to a specific deployed contract.
func NewBackendCheck(address common.Address, backend bind.ContractBackend) (*BackendCheck, error) {
	contract, err := bindBackendCheck(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BackendCheck{BackendCheckCaller: BackendCheckCaller{contract: contract}, BackendCheckTransactor: BackendCheckTransactor{contract: contract}, BackendCheckFilterer: BackendCheckFilterer{contract: contract}}, nil
}

// NewBackendCheckCaller creates a new read-only instance of BackendCheck, bound to a specific deployed contract.
func NewBackendCheckCaller(address common.Address, caller bind.ContractCaller) (*BackendCheckCaller, error) {
	contract, err := bindBackendCheck(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BackendCheckCaller{contract: contract}, nil
}

// NewBackendCheckTransactor creates a new write-only instance of BackendCheck, bound to a specific deployed contract.
func NewBackendCheckTransactor(address common.Address, transactor bind.ContractTransactor) (*BackendCheckTransactor, error) {
	contract, err := bindBackendCheck(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BackendCheckTransactor{contract: contract}, nil
}

// NewBackendCheckFilterer creates a new log filterer instance of BackendCheck, bound to a specific deployed contract.
func NewBackendCheckFilterer(address common.Address, filterer bind.ContractFilterer) (*BackendCheckFilterer, error) {
	contract, err := bindBackendCheck(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BackendCheckFilterer{contract: contract}, nil
}

// bindBackendCheck binds a generic wrapper to an already deployed contract.
func bindBackendCheck(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BackendCheckMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BackendCheck *BackendCheckRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BackendCheck.Contract.BackendCheckCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BackendCheck *BackendCheckRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BackendCheck.Contract.BackendCheckTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BackendCheck *BackendCheckRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BackendCheck.Contract.BackendCheckTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BackendCheck *BackendCheckCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BackendCheck.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BackendCheck *BackendCheckTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BackendCheck.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BackendCheck *BackendCheckTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BackendCheck.Contract.contract.Transact(opts, method, params...)
}

// GetRemoveHash is a free data retrieval call binding the contract method 0x2831b038.
//
// Solidity: function getRemoveHash(uint256 maxAmount, uint256 price, uint256 height) pure returns(bytes32)
func (_BackendCheck *BackendCheckCaller) GetRemoveHash(opts *bind.CallOpts, maxAmount *big.Int, price *big.Int, height *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BackendCheck.contract.Call(opts, &out, "getRemoveHash", maxAmount, price, height)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRemoveHash is a free data retrieval call binding the contract method 0x2831b038.
//
// Solidity: function getRemoveHash(uint256 maxAmount, uint256 price, uint256 height) pure returns(bytes32)
func (_BackendCheck *BackendCheckSession) GetRemoveHash(maxAmount *big.Int, price *big.Int, height *big.Int) ([32]byte, error) {
	return _BackendCheck.Contract.GetRemoveHash(&_BackendCheck.CallOpts, maxAmount, price, height)
}

// GetRemoveHash is a free data retrieval call binding the contract method 0x2831b038.
//
// Solidity: function getRemoveHash(uint256 maxAmount, uint256 price, uint256 height) pure returns(bytes32)
func (_BackendCheck *BackendCheckCallerSession) GetRemoveHash(maxAmount *big.Int, price *big.Int, height *big.Int) ([32]byte, error) {
	return _BackendCheck.Contract.GetRemoveHash(&_BackendCheck.CallOpts, maxAmount, price, height)
}

// GetSwapHash is a free data retrieval call binding the contract method 0xc274abeb.
//
// Solidity: function getSwapHash(uint256 badDebt, uint256 height, uint256 price, bytes32[] list) pure returns(bytes32)
func (_BackendCheck *BackendCheckCaller) GetSwapHash(opts *bind.CallOpts, badDebt *big.Int, height *big.Int, price *big.Int, list [][32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BackendCheck.contract.Call(opts, &out, "getSwapHash", badDebt, height, price, list)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetSwapHash is a free data retrieval call binding the contract method 0xc274abeb.
//
// Solidity: function getSwapHash(uint256 badDebt, uint256 height, uint256 price, bytes32[] list) pure returns(bytes32)
func (_BackendCheck *BackendCheckSession) GetSwapHash(badDebt *big.Int, height *big.Int, price *big.Int, list [][32]byte) ([32]byte, error) {
	return _BackendCheck.Contract.GetSwapHash(&_BackendCheck.CallOpts, badDebt, height, price, list)
}

// GetSwapHash is a free data retrieval call binding the contract method 0xc274abeb.
//
// Solidity: function getSwapHash(uint256 badDebt, uint256 height, uint256 price, bytes32[] list) pure returns(bytes32)
func (_BackendCheck *BackendCheckCallerSession) GetSwapHash(badDebt *big.Int, height *big.Int, price *big.Int, list [][32]byte) ([32]byte, error) {
	return _BackendCheck.Contract.GetSwapHash(&_BackendCheck.CallOpts, badDebt, height, price, list)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 _msgHash, bytes _signature) pure returns(address)
func (_BackendCheck *BackendCheckCaller) RecoverSigner(opts *bind.CallOpts, _msgHash [32]byte, _signature []byte) (common.Address, error) {
	var out []interface{}
	err := _BackendCheck.contract.Call(opts, &out, "recoverSigner", _msgHash, _signature)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 _msgHash, bytes _signature) pure returns(address)
func (_BackendCheck *BackendCheckSession) RecoverSigner(_msgHash [32]byte, _signature []byte) (common.Address, error) {
	return _BackendCheck.Contract.RecoverSigner(&_BackendCheck.CallOpts, _msgHash, _signature)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 _msgHash, bytes _signature) pure returns(address)
func (_BackendCheck *BackendCheckCallerSession) RecoverSigner(_msgHash [32]byte, _signature []byte) (common.Address, error) {
	return _BackendCheck.Contract.RecoverSigner(&_BackendCheck.CallOpts, _msgHash, _signature)
}

// VerifySignature is a free data retrieval call binding the contract method 0xdaca6f78.
//
// Solidity: function verifySignature(bytes32 _msgHash, bytes _signature) pure returns(address)
func (_BackendCheck *BackendCheckCaller) VerifySignature(opts *bind.CallOpts, _msgHash [32]byte, _signature []byte) (common.Address, error) {
	var out []interface{}
	err := _BackendCheck.contract.Call(opts, &out, "verifySignature", _msgHash, _signature)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VerifySignature is a free data retrieval call binding the contract method 0xdaca6f78.
//
// Solidity: function verifySignature(bytes32 _msgHash, bytes _signature) pure returns(address)
func (_BackendCheck *BackendCheckSession) VerifySignature(_msgHash [32]byte, _signature []byte) (common.Address, error) {
	return _BackendCheck.Contract.VerifySignature(&_BackendCheck.CallOpts, _msgHash, _signature)
}

// VerifySignature is a free data retrieval call binding the contract method 0xdaca6f78.
//
// Solidity: function verifySignature(bytes32 _msgHash, bytes _signature) pure returns(address)
func (_BackendCheck *BackendCheckCallerSession) VerifySignature(_msgHash [32]byte, _signature []byte) (common.Address, error) {
	return _BackendCheck.Contract.VerifySignature(&_BackendCheck.CallOpts, _msgHash, _signature)
}
