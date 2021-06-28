// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package compound

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

// PriceOracleABI is the input ABI used to generate the binding from.
const PriceOracleABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"contractCToken\",\"name\":\"cToken\",\"type\":\"address\"}],\"name\":\"getUnderlyingPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isPriceOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PriceOracle is an auto generated Go binding around an Ethereum contract.
type PriceOracle struct {
	PriceOracleCaller     // Read-only binding to the contract
	PriceOracleTransactor // Write-only binding to the contract
	PriceOracleFilterer   // Log filterer for contract events
}

// PriceOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type PriceOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PriceOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PriceOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PriceOracleSession struct {
	Contract     *PriceOracle      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PriceOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PriceOracleCallerSession struct {
	Contract *PriceOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PriceOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PriceOracleTransactorSession struct {
	Contract     *PriceOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PriceOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type PriceOracleRaw struct {
	Contract *PriceOracle // Generic contract binding to access the raw methods on
}

// PriceOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PriceOracleCallerRaw struct {
	Contract *PriceOracleCaller // Generic read-only contract binding to access the raw methods on
}

// PriceOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PriceOracleTransactorRaw struct {
	Contract *PriceOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPriceOracle creates a new instance of PriceOracle, bound to a specific deployed contract.
func NewPriceOracle(address common.Address, backend bind.ContractBackend) (*PriceOracle, error) {
	contract, err := bindPriceOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PriceOracle{PriceOracleCaller: PriceOracleCaller{contract: contract}, PriceOracleTransactor: PriceOracleTransactor{contract: contract}, PriceOracleFilterer: PriceOracleFilterer{contract: contract}}, nil
}

// NewPriceOracleCaller creates a new read-only instance of PriceOracle, bound to a specific deployed contract.
func NewPriceOracleCaller(address common.Address, caller bind.ContractCaller) (*PriceOracleCaller, error) {
	contract, err := bindPriceOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PriceOracleCaller{contract: contract}, nil
}

// NewPriceOracleTransactor creates a new write-only instance of PriceOracle, bound to a specific deployed contract.
func NewPriceOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*PriceOracleTransactor, error) {
	contract, err := bindPriceOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PriceOracleTransactor{contract: contract}, nil
}

// NewPriceOracleFilterer creates a new log filterer instance of PriceOracle, bound to a specific deployed contract.
func NewPriceOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*PriceOracleFilterer, error) {
	contract, err := bindPriceOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PriceOracleFilterer{contract: contract}, nil
}

// bindPriceOracle binds a generic wrapper to an already deployed contract.
func bindPriceOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PriceOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PriceOracle *PriceOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceOracle.Contract.PriceOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PriceOracle *PriceOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceOracle.Contract.PriceOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PriceOracle *PriceOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceOracle.Contract.PriceOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PriceOracle *PriceOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PriceOracle *PriceOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PriceOracle *PriceOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceOracle.Contract.contract.Transact(opts, method, params...)
}

// GetUnderlyingPrice is a free data retrieval call binding the contract method 0xfc57d4df.
//
// Solidity: function getUnderlyingPrice(address cToken) view returns(uint256)
func (_PriceOracle *PriceOracleCaller) GetUnderlyingPrice(opts *bind.CallOpts, cToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PriceOracle.contract.Call(opts, &out, "getUnderlyingPrice", cToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnderlyingPrice is a free data retrieval call binding the contract method 0xfc57d4df.
//
// Solidity: function getUnderlyingPrice(address cToken) view returns(uint256)
func (_PriceOracle *PriceOracleSession) GetUnderlyingPrice(cToken common.Address) (*big.Int, error) {
	return _PriceOracle.Contract.GetUnderlyingPrice(&_PriceOracle.CallOpts, cToken)
}

// GetUnderlyingPrice is a free data retrieval call binding the contract method 0xfc57d4df.
//
// Solidity: function getUnderlyingPrice(address cToken) view returns(uint256)
func (_PriceOracle *PriceOracleCallerSession) GetUnderlyingPrice(cToken common.Address) (*big.Int, error) {
	return _PriceOracle.Contract.GetUnderlyingPrice(&_PriceOracle.CallOpts, cToken)
}

// IsPriceOracle is a free data retrieval call binding the contract method 0x66331bba.
//
// Solidity: function isPriceOracle() view returns(bool)
func (_PriceOracle *PriceOracleCaller) IsPriceOracle(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PriceOracle.contract.Call(opts, &out, "isPriceOracle")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPriceOracle is a free data retrieval call binding the contract method 0x66331bba.
//
// Solidity: function isPriceOracle() view returns(bool)
func (_PriceOracle *PriceOracleSession) IsPriceOracle() (bool, error) {
	return _PriceOracle.Contract.IsPriceOracle(&_PriceOracle.CallOpts)
}

// IsPriceOracle is a free data retrieval call binding the contract method 0x66331bba.
//
// Solidity: function isPriceOracle() view returns(bool)
func (_PriceOracle *PriceOracleCallerSession) IsPriceOracle() (bool, error) {
	return _PriceOracle.Contract.IsPriceOracle(&_PriceOracle.CallOpts)
}
