// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package backend

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

// PoolAsset is an auto generated low-level Go binding around an user-defined struct.
type PoolAsset struct {
	Token                   common.Address
	BorrowIndex             *big.Int
	BorrowRate              *big.Int
	TotalCollateral         *big.Int
	TotalCollateralWithDebt *big.Int
	TotalDebtScaled         *big.Int
	UnclaimedFee            *big.Int
}

// PoolProps is an auto generated low-level Go binding around an user-defined struct.
type PoolProps struct {
	Assets               [2]PoolAsset
	Bank                 common.Address
	InterestRateStrategy common.Address
	Configuration        *big.Int
	LastUpdateTimestamp  *big.Int
}

// BackendMetaData contains all meta data concerning the Backend contract.
var BackendMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"verifyHeight\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateralWithDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unclaimedFee\",\"type\":\"uint256\"}],\"internalType\":\"structPool.Asset[2]\",\"name\":\"assets\",\"type\":\"tuple[2]\"},{\"internalType\":\"address\",\"name\":\"bank\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"configuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structPool.Props\",\"name\":\"pool\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"verifyPrice\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxAmount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"verifyRemoveSignature\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"badDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"list\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"verifySwapSignature\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BackendABI is the input ABI used to generate the binding from.
// Deprecated: Use BackendMetaData.ABI instead.
var BackendABI = BackendMetaData.ABI

// Backend is an auto generated Go binding around an Ethereum contract.
type Backend struct {
	BackendCaller     // Read-only binding to the contract
	BackendTransactor // Write-only binding to the contract
	BackendFilterer   // Log filterer for contract events
}

// BackendCaller is an auto generated read-only Go binding around an Ethereum contract.
type BackendCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BackendTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BackendTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BackendFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BackendFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BackendSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BackendSession struct {
	Contract     *Backend          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BackendCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BackendCallerSession struct {
	Contract *BackendCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BackendTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BackendTransactorSession struct {
	Contract     *BackendTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BackendRaw is an auto generated low-level Go binding around an Ethereum contract.
type BackendRaw struct {
	Contract *Backend // Generic contract binding to access the raw methods on
}

// BackendCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BackendCallerRaw struct {
	Contract *BackendCaller // Generic read-only contract binding to access the raw methods on
}

// BackendTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BackendTransactorRaw struct {
	Contract *BackendTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBackend creates a new instance of Backend, bound to a specific deployed contract.
func NewBackend(address common.Address, backend bind.ContractBackend) (*Backend, error) {
	contract, err := bindBackend(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Backend{BackendCaller: BackendCaller{contract: contract}, BackendTransactor: BackendTransactor{contract: contract}, BackendFilterer: BackendFilterer{contract: contract}}, nil
}

// NewBackendCaller creates a new read-only instance of Backend, bound to a specific deployed contract.
func NewBackendCaller(address common.Address, caller bind.ContractCaller) (*BackendCaller, error) {
	contract, err := bindBackend(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BackendCaller{contract: contract}, nil
}

// NewBackendTransactor creates a new write-only instance of Backend, bound to a specific deployed contract.
func NewBackendTransactor(address common.Address, transactor bind.ContractTransactor) (*BackendTransactor, error) {
	contract, err := bindBackend(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BackendTransactor{contract: contract}, nil
}

// NewBackendFilterer creates a new log filterer instance of Backend, bound to a specific deployed contract.
func NewBackendFilterer(address common.Address, filterer bind.ContractFilterer) (*BackendFilterer, error) {
	contract, err := bindBackend(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BackendFilterer{contract: contract}, nil
}

// bindBackend binds a generic wrapper to an already deployed contract.
func bindBackend(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BackendMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Backend *BackendRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Backend.Contract.BackendCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Backend *BackendRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Backend.Contract.BackendTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Backend *BackendRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Backend.Contract.BackendTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Backend *BackendCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Backend.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Backend *BackendTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Backend.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Backend *BackendTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Backend.Contract.contract.Transact(opts, method, params...)
}

// VerifyHeight is a free data retrieval call binding the contract method 0x7db7ab1e.
//
// Solidity: function verifyHeight(uint256 height) view returns()
func (_Backend *BackendCaller) VerifyHeight(opts *bind.CallOpts, height *big.Int) error {
	var out []interface{}
	err := _Backend.contract.Call(opts, &out, "verifyHeight", height)

	if err != nil {
		return err
	}

	return err

}

// VerifyHeight is a free data retrieval call binding the contract method 0x7db7ab1e.
//
// Solidity: function verifyHeight(uint256 height) view returns()
func (_Backend *BackendSession) VerifyHeight(height *big.Int) error {
	return _Backend.Contract.VerifyHeight(&_Backend.CallOpts, height)
}

// VerifyHeight is a free data retrieval call binding the contract method 0x7db7ab1e.
//
// Solidity: function verifyHeight(uint256 height) view returns()
func (_Backend *BackendCallerSession) VerifyHeight(height *big.Int) error {
	return _Backend.Contract.VerifyHeight(&_Backend.CallOpts, height)
}

// VerifyPrice is a free data retrieval call binding the contract method 0x0f8ef2c3.
//
// Solidity: function verifyPrice(((address,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256) pool, uint256 price) view returns()
func (_Backend *BackendCaller) VerifyPrice(opts *bind.CallOpts, pool PoolProps, price *big.Int) error {
	var out []interface{}
	err := _Backend.contract.Call(opts, &out, "verifyPrice", pool, price)

	if err != nil {
		return err
	}

	return err

}

// VerifyPrice is a free data retrieval call binding the contract method 0x0f8ef2c3.
//
// Solidity: function verifyPrice(((address,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256) pool, uint256 price) view returns()
func (_Backend *BackendSession) VerifyPrice(pool PoolProps, price *big.Int) error {
	return _Backend.Contract.VerifyPrice(&_Backend.CallOpts, pool, price)
}

// VerifyPrice is a free data retrieval call binding the contract method 0x0f8ef2c3.
//
// Solidity: function verifyPrice(((address,uint256,uint256,uint256,uint256,uint256,uint256)[2],address,address,uint256,uint256) pool, uint256 price) view returns()
func (_Backend *BackendCallerSession) VerifyPrice(pool PoolProps, price *big.Int) error {
	return _Backend.Contract.VerifyPrice(&_Backend.CallOpts, pool, price)
}

// VerifyRemoveSignature is a free data retrieval call binding the contract method 0xbb3b027c.
//
// Solidity: function verifyRemoveSignature(uint256 maxAmount0, uint256 height, uint256 price, bytes signature) view returns()
func (_Backend *BackendCaller) VerifyRemoveSignature(opts *bind.CallOpts, maxAmount0 *big.Int, height *big.Int, price *big.Int, signature []byte) error {
	var out []interface{}
	err := _Backend.contract.Call(opts, &out, "verifyRemoveSignature", maxAmount0, height, price, signature)

	if err != nil {
		return err
	}

	return err

}

// VerifyRemoveSignature is a free data retrieval call binding the contract method 0xbb3b027c.
//
// Solidity: function verifyRemoveSignature(uint256 maxAmount0, uint256 height, uint256 price, bytes signature) view returns()
func (_Backend *BackendSession) VerifyRemoveSignature(maxAmount0 *big.Int, height *big.Int, price *big.Int, signature []byte) error {
	return _Backend.Contract.VerifyRemoveSignature(&_Backend.CallOpts, maxAmount0, height, price, signature)
}

// VerifyRemoveSignature is a free data retrieval call binding the contract method 0xbb3b027c.
//
// Solidity: function verifyRemoveSignature(uint256 maxAmount0, uint256 height, uint256 price, bytes signature) view returns()
func (_Backend *BackendCallerSession) VerifyRemoveSignature(maxAmount0 *big.Int, height *big.Int, price *big.Int, signature []byte) error {
	return _Backend.Contract.VerifyRemoveSignature(&_Backend.CallOpts, maxAmount0, height, price, signature)
}

// VerifySwapSignature is a free data retrieval call binding the contract method 0x0f600ca5.
//
// Solidity: function verifySwapSignature(uint256 badDebt, uint256 price, uint256 height, bytes32[] list, bytes signature) view returns()
func (_Backend *BackendCaller) VerifySwapSignature(opts *bind.CallOpts, badDebt *big.Int, price *big.Int, height *big.Int, list [][32]byte, signature []byte) error {
	var out []interface{}
	err := _Backend.contract.Call(opts, &out, "verifySwapSignature", badDebt, price, height, list, signature)

	if err != nil {
		return err
	}

	return err

}

// VerifySwapSignature is a free data retrieval call binding the contract method 0x0f600ca5.
//
// Solidity: function verifySwapSignature(uint256 badDebt, uint256 price, uint256 height, bytes32[] list, bytes signature) view returns()
func (_Backend *BackendSession) VerifySwapSignature(badDebt *big.Int, price *big.Int, height *big.Int, list [][32]byte, signature []byte) error {
	return _Backend.Contract.VerifySwapSignature(&_Backend.CallOpts, badDebt, price, height, list, signature)
}

// VerifySwapSignature is a free data retrieval call binding the contract method 0x0f600ca5.
//
// Solidity: function verifySwapSignature(uint256 badDebt, uint256 price, uint256 height, bytes32[] list, bytes signature) view returns()
func (_Backend *BackendCallerSession) VerifySwapSignature(badDebt *big.Int, price *big.Int, height *big.Int, list [][32]byte, signature []byte) error {
	return _Backend.Contract.VerifySwapSignature(&_Backend.CallOpts, badDebt, price, height, list, signature)
}
