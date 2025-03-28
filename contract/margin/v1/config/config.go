// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package config

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

// ConfigMetaData contains all meta data concerning the Config contract.
var ConfigMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"_roleStore\",\"type\":\"address\"},{\"internalType\":\"contractDataStore\",\"name\":\"_dataStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"EmptyPool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyTokenBase\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyTreasury\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"MaxValidFeeFactor\",\"type\":\"uint256\"}],\"name\":\"InvalidFeeFactor\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"role\",\"type\":\"string\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"dataStore\",\"outputs\":[{\"internalType\":\"contractDataStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleStore\",\"outputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"debtSaftyFactor\",\"type\":\"uint256\"}],\"name\":\"setDebtSafetyFactor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"interestRateStrategy\",\"type\":\"address\"}],\"name\":\"setDefaultInterestRateStrategy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"configuration\",\"type\":\"uint256\"}],\"name\":\"setDefaultPoolConfiguration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"liquidationFee\",\"type\":\"uint256\"}],\"name\":\"setLiquidationFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setMarginLevelThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxBorrowRate\",\"type\":\"uint256\"}],\"name\":\"setMaxBorrowRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxDepositRate\",\"type\":\"uint256\"}],\"name\":\"setMaxDepositRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shortEnabled\",\"type\":\"bool\"}],\"name\":\"setShortEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"shortLiquidityThreshold\",\"type\":\"uint256\"}],\"name\":\"setShortLiquidityThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeFactor\",\"type\":\"uint256\"}],\"name\":\"setSwapFeeFactor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"setTokenBase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"multipierFactor\",\"type\":\"uint256\"}],\"name\":\"setTradableDebtMultipierFactor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"}],\"name\":\"setTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeFactor\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFeeFactor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setTwapPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ConfigABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfigMetaData.ABI instead.
var ConfigABI = ConfigMetaData.ABI

// Config is an auto generated Go binding around an Ethereum contract.
type Config struct {
	ConfigCaller     // Read-only binding to the contract
	ConfigTransactor // Write-only binding to the contract
	ConfigFilterer   // Log filterer for contract events
}

// ConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfigSession struct {
	Contract     *Config           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfigCallerSession struct {
	Contract *ConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfigTransactorSession struct {
	Contract     *ConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfigRaw struct {
	Contract *Config // Generic contract binding to access the raw methods on
}

// ConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfigCallerRaw struct {
	Contract *ConfigCaller // Generic read-only contract binding to access the raw methods on
}

// ConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfigTransactorRaw struct {
	Contract *ConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfig creates a new instance of Config, bound to a specific deployed contract.
func NewConfig(address common.Address, backend bind.ContractBackend) (*Config, error) {
	contract, err := bindConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Config{ConfigCaller: ConfigCaller{contract: contract}, ConfigTransactor: ConfigTransactor{contract: contract}, ConfigFilterer: ConfigFilterer{contract: contract}}, nil
}

// NewConfigCaller creates a new read-only instance of Config, bound to a specific deployed contract.
func NewConfigCaller(address common.Address, caller bind.ContractCaller) (*ConfigCaller, error) {
	contract, err := bindConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigCaller{contract: contract}, nil
}

// NewConfigTransactor creates a new write-only instance of Config, bound to a specific deployed contract.
func NewConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfigTransactor, error) {
	contract, err := bindConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigTransactor{contract: contract}, nil
}

// NewConfigFilterer creates a new log filterer instance of Config, bound to a specific deployed contract.
func NewConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfigFilterer, error) {
	contract, err := bindConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfigFilterer{contract: contract}, nil
}

// bindConfig binds a generic wrapper to an already deployed contract.
func bindConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Config *ConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Config.Contract.ConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Config *ConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Config.Contract.ConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Config *ConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Config.Contract.ConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Config *ConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Config.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Config *ConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Config.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Config *ConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Config.Contract.contract.Transact(opts, method, params...)
}

// DataStore is a free data retrieval call binding the contract method 0x660d0d67.
//
// Solidity: function dataStore() view returns(address)
func (_Config *ConfigCaller) DataStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Config.contract.Call(opts, &out, "dataStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DataStore is a free data retrieval call binding the contract method 0x660d0d67.
//
// Solidity: function dataStore() view returns(address)
func (_Config *ConfigSession) DataStore() (common.Address, error) {
	return _Config.Contract.DataStore(&_Config.CallOpts)
}

// DataStore is a free data retrieval call binding the contract method 0x660d0d67.
//
// Solidity: function dataStore() view returns(address)
func (_Config *ConfigCallerSession) DataStore() (common.Address, error) {
	return _Config.Contract.DataStore(&_Config.CallOpts)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_Config *ConfigCaller) RoleStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Config.contract.Call(opts, &out, "roleStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_Config *ConfigSession) RoleStore() (common.Address, error) {
	return _Config.Contract.RoleStore(&_Config.CallOpts)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_Config *ConfigCallerSession) RoleStore() (common.Address, error) {
	return _Config.Contract.RoleStore(&_Config.CallOpts)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Config *ConfigTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Config *ConfigSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Config.Contract.Multicall(&_Config.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Config *ConfigTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Config.Contract.Multicall(&_Config.TransactOpts, data)
}

// SetDebtSafetyFactor is a paid mutator transaction binding the contract method 0xbb789bd7.
//
// Solidity: function setDebtSafetyFactor(uint256 debtSaftyFactor) returns()
func (_Config *ConfigTransactor) SetDebtSafetyFactor(opts *bind.TransactOpts, debtSaftyFactor *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setDebtSafetyFactor", debtSaftyFactor)
}

// SetDebtSafetyFactor is a paid mutator transaction binding the contract method 0xbb789bd7.
//
// Solidity: function setDebtSafetyFactor(uint256 debtSaftyFactor) returns()
func (_Config *ConfigSession) SetDebtSafetyFactor(debtSaftyFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetDebtSafetyFactor(&_Config.TransactOpts, debtSaftyFactor)
}

// SetDebtSafetyFactor is a paid mutator transaction binding the contract method 0xbb789bd7.
//
// Solidity: function setDebtSafetyFactor(uint256 debtSaftyFactor) returns()
func (_Config *ConfigTransactorSession) SetDebtSafetyFactor(debtSaftyFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetDebtSafetyFactor(&_Config.TransactOpts, debtSaftyFactor)
}

// SetDefaultInterestRateStrategy is a paid mutator transaction binding the contract method 0x69ea8682.
//
// Solidity: function setDefaultInterestRateStrategy(address interestRateStrategy) returns()
func (_Config *ConfigTransactor) SetDefaultInterestRateStrategy(opts *bind.TransactOpts, interestRateStrategy common.Address) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setDefaultInterestRateStrategy", interestRateStrategy)
}

// SetDefaultInterestRateStrategy is a paid mutator transaction binding the contract method 0x69ea8682.
//
// Solidity: function setDefaultInterestRateStrategy(address interestRateStrategy) returns()
func (_Config *ConfigSession) SetDefaultInterestRateStrategy(interestRateStrategy common.Address) (*types.Transaction, error) {
	return _Config.Contract.SetDefaultInterestRateStrategy(&_Config.TransactOpts, interestRateStrategy)
}

// SetDefaultInterestRateStrategy is a paid mutator transaction binding the contract method 0x69ea8682.
//
// Solidity: function setDefaultInterestRateStrategy(address interestRateStrategy) returns()
func (_Config *ConfigTransactorSession) SetDefaultInterestRateStrategy(interestRateStrategy common.Address) (*types.Transaction, error) {
	return _Config.Contract.SetDefaultInterestRateStrategy(&_Config.TransactOpts, interestRateStrategy)
}

// SetDefaultPoolConfiguration is a paid mutator transaction binding the contract method 0x5eadca7e.
//
// Solidity: function setDefaultPoolConfiguration(uint256 configuration) returns()
func (_Config *ConfigTransactor) SetDefaultPoolConfiguration(opts *bind.TransactOpts, configuration *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setDefaultPoolConfiguration", configuration)
}

// SetDefaultPoolConfiguration is a paid mutator transaction binding the contract method 0x5eadca7e.
//
// Solidity: function setDefaultPoolConfiguration(uint256 configuration) returns()
func (_Config *ConfigSession) SetDefaultPoolConfiguration(configuration *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetDefaultPoolConfiguration(&_Config.TransactOpts, configuration)
}

// SetDefaultPoolConfiguration is a paid mutator transaction binding the contract method 0x5eadca7e.
//
// Solidity: function setDefaultPoolConfiguration(uint256 configuration) returns()
func (_Config *ConfigTransactorSession) SetDefaultPoolConfiguration(configuration *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetDefaultPoolConfiguration(&_Config.TransactOpts, configuration)
}

// SetLiquidationFee is a paid mutator transaction binding the contract method 0xfef0bec8.
//
// Solidity: function setLiquidationFee(uint256 liquidationFee) returns()
func (_Config *ConfigTransactor) SetLiquidationFee(opts *bind.TransactOpts, liquidationFee *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setLiquidationFee", liquidationFee)
}

// SetLiquidationFee is a paid mutator transaction binding the contract method 0xfef0bec8.
//
// Solidity: function setLiquidationFee(uint256 liquidationFee) returns()
func (_Config *ConfigSession) SetLiquidationFee(liquidationFee *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetLiquidationFee(&_Config.TransactOpts, liquidationFee)
}

// SetLiquidationFee is a paid mutator transaction binding the contract method 0xfef0bec8.
//
// Solidity: function setLiquidationFee(uint256 liquidationFee) returns()
func (_Config *ConfigTransactorSession) SetLiquidationFee(liquidationFee *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetLiquidationFee(&_Config.TransactOpts, liquidationFee)
}

// SetMarginLevelThreshold is a paid mutator transaction binding the contract method 0x88235d10.
//
// Solidity: function setMarginLevelThreshold(uint256 threshold) returns()
func (_Config *ConfigTransactor) SetMarginLevelThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setMarginLevelThreshold", threshold)
}

// SetMarginLevelThreshold is a paid mutator transaction binding the contract method 0x88235d10.
//
// Solidity: function setMarginLevelThreshold(uint256 threshold) returns()
func (_Config *ConfigSession) SetMarginLevelThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetMarginLevelThreshold(&_Config.TransactOpts, threshold)
}

// SetMarginLevelThreshold is a paid mutator transaction binding the contract method 0x88235d10.
//
// Solidity: function setMarginLevelThreshold(uint256 threshold) returns()
func (_Config *ConfigTransactorSession) SetMarginLevelThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetMarginLevelThreshold(&_Config.TransactOpts, threshold)
}

// SetMaxBorrowRate is a paid mutator transaction binding the contract method 0x57c25c66.
//
// Solidity: function setMaxBorrowRate(uint256 maxBorrowRate) returns()
func (_Config *ConfigTransactor) SetMaxBorrowRate(opts *bind.TransactOpts, maxBorrowRate *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setMaxBorrowRate", maxBorrowRate)
}

// SetMaxBorrowRate is a paid mutator transaction binding the contract method 0x57c25c66.
//
// Solidity: function setMaxBorrowRate(uint256 maxBorrowRate) returns()
func (_Config *ConfigSession) SetMaxBorrowRate(maxBorrowRate *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetMaxBorrowRate(&_Config.TransactOpts, maxBorrowRate)
}

// SetMaxBorrowRate is a paid mutator transaction binding the contract method 0x57c25c66.
//
// Solidity: function setMaxBorrowRate(uint256 maxBorrowRate) returns()
func (_Config *ConfigTransactorSession) SetMaxBorrowRate(maxBorrowRate *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetMaxBorrowRate(&_Config.TransactOpts, maxBorrowRate)
}

// SetMaxDepositRate is a paid mutator transaction binding the contract method 0x956228c6.
//
// Solidity: function setMaxDepositRate(uint256 maxDepositRate) returns()
func (_Config *ConfigTransactor) SetMaxDepositRate(opts *bind.TransactOpts, maxDepositRate *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setMaxDepositRate", maxDepositRate)
}

// SetMaxDepositRate is a paid mutator transaction binding the contract method 0x956228c6.
//
// Solidity: function setMaxDepositRate(uint256 maxDepositRate) returns()
func (_Config *ConfigSession) SetMaxDepositRate(maxDepositRate *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetMaxDepositRate(&_Config.TransactOpts, maxDepositRate)
}

// SetMaxDepositRate is a paid mutator transaction binding the contract method 0x956228c6.
//
// Solidity: function setMaxDepositRate(uint256 maxDepositRate) returns()
func (_Config *ConfigTransactorSession) SetMaxDepositRate(maxDepositRate *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetMaxDepositRate(&_Config.TransactOpts, maxDepositRate)
}

// SetShortEnabled is a paid mutator transaction binding the contract method 0x2ffaf559.
//
// Solidity: function setShortEnabled(address token0, address token1, bool shortEnabled) returns()
func (_Config *ConfigTransactor) SetShortEnabled(opts *bind.TransactOpts, token0 common.Address, token1 common.Address, shortEnabled bool) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setShortEnabled", token0, token1, shortEnabled)
}

// SetShortEnabled is a paid mutator transaction binding the contract method 0x2ffaf559.
//
// Solidity: function setShortEnabled(address token0, address token1, bool shortEnabled) returns()
func (_Config *ConfigSession) SetShortEnabled(token0 common.Address, token1 common.Address, shortEnabled bool) (*types.Transaction, error) {
	return _Config.Contract.SetShortEnabled(&_Config.TransactOpts, token0, token1, shortEnabled)
}

// SetShortEnabled is a paid mutator transaction binding the contract method 0x2ffaf559.
//
// Solidity: function setShortEnabled(address token0, address token1, bool shortEnabled) returns()
func (_Config *ConfigTransactorSession) SetShortEnabled(token0 common.Address, token1 common.Address, shortEnabled bool) (*types.Transaction, error) {
	return _Config.Contract.SetShortEnabled(&_Config.TransactOpts, token0, token1, shortEnabled)
}

// SetShortLiquidityThreshold is a paid mutator transaction binding the contract method 0x92547c38.
//
// Solidity: function setShortLiquidityThreshold(uint256 shortLiquidityThreshold) returns()
func (_Config *ConfigTransactor) SetShortLiquidityThreshold(opts *bind.TransactOpts, shortLiquidityThreshold *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setShortLiquidityThreshold", shortLiquidityThreshold)
}

// SetShortLiquidityThreshold is a paid mutator transaction binding the contract method 0x92547c38.
//
// Solidity: function setShortLiquidityThreshold(uint256 shortLiquidityThreshold) returns()
func (_Config *ConfigSession) SetShortLiquidityThreshold(shortLiquidityThreshold *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetShortLiquidityThreshold(&_Config.TransactOpts, shortLiquidityThreshold)
}

// SetShortLiquidityThreshold is a paid mutator transaction binding the contract method 0x92547c38.
//
// Solidity: function setShortLiquidityThreshold(uint256 shortLiquidityThreshold) returns()
func (_Config *ConfigTransactorSession) SetShortLiquidityThreshold(shortLiquidityThreshold *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetShortLiquidityThreshold(&_Config.TransactOpts, shortLiquidityThreshold)
}

// SetSwapFeeFactor is a paid mutator transaction binding the contract method 0x6a97c0fe.
//
// Solidity: function setSwapFeeFactor(address token0, address token1, uint256 feeFactor) returns()
func (_Config *ConfigTransactor) SetSwapFeeFactor(opts *bind.TransactOpts, token0 common.Address, token1 common.Address, feeFactor *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setSwapFeeFactor", token0, token1, feeFactor)
}

// SetSwapFeeFactor is a paid mutator transaction binding the contract method 0x6a97c0fe.
//
// Solidity: function setSwapFeeFactor(address token0, address token1, uint256 feeFactor) returns()
func (_Config *ConfigSession) SetSwapFeeFactor(token0 common.Address, token1 common.Address, feeFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetSwapFeeFactor(&_Config.TransactOpts, token0, token1, feeFactor)
}

// SetSwapFeeFactor is a paid mutator transaction binding the contract method 0x6a97c0fe.
//
// Solidity: function setSwapFeeFactor(address token0, address token1, uint256 feeFactor) returns()
func (_Config *ConfigTransactorSession) SetSwapFeeFactor(token0 common.Address, token1 common.Address, feeFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetSwapFeeFactor(&_Config.TransactOpts, token0, token1, feeFactor)
}

// SetTokenBase is a paid mutator transaction binding the contract method 0x1dee7f3e.
//
// Solidity: function setTokenBase(address token) returns()
func (_Config *ConfigTransactor) SetTokenBase(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setTokenBase", token)
}

// SetTokenBase is a paid mutator transaction binding the contract method 0x1dee7f3e.
//
// Solidity: function setTokenBase(address token) returns()
func (_Config *ConfigSession) SetTokenBase(token common.Address) (*types.Transaction, error) {
	return _Config.Contract.SetTokenBase(&_Config.TransactOpts, token)
}

// SetTokenBase is a paid mutator transaction binding the contract method 0x1dee7f3e.
//
// Solidity: function setTokenBase(address token) returns()
func (_Config *ConfigTransactorSession) SetTokenBase(token common.Address) (*types.Transaction, error) {
	return _Config.Contract.SetTokenBase(&_Config.TransactOpts, token)
}

// SetTradableDebtMultipierFactor is a paid mutator transaction binding the contract method 0x6a72c63f.
//
// Solidity: function setTradableDebtMultipierFactor(uint256 multipierFactor) returns()
func (_Config *ConfigTransactor) SetTradableDebtMultipierFactor(opts *bind.TransactOpts, multipierFactor *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setTradableDebtMultipierFactor", multipierFactor)
}

// SetTradableDebtMultipierFactor is a paid mutator transaction binding the contract method 0x6a72c63f.
//
// Solidity: function setTradableDebtMultipierFactor(uint256 multipierFactor) returns()
func (_Config *ConfigSession) SetTradableDebtMultipierFactor(multipierFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetTradableDebtMultipierFactor(&_Config.TransactOpts, multipierFactor)
}

// SetTradableDebtMultipierFactor is a paid mutator transaction binding the contract method 0x6a72c63f.
//
// Solidity: function setTradableDebtMultipierFactor(uint256 multipierFactor) returns()
func (_Config *ConfigTransactorSession) SetTradableDebtMultipierFactor(multipierFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetTradableDebtMultipierFactor(&_Config.TransactOpts, multipierFactor)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury) returns()
func (_Config *ConfigTransactor) SetTreasury(opts *bind.TransactOpts, treasury common.Address) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setTreasury", treasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury) returns()
func (_Config *ConfigSession) SetTreasury(treasury common.Address) (*types.Transaction, error) {
	return _Config.Contract.SetTreasury(&_Config.TransactOpts, treasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury) returns()
func (_Config *ConfigTransactorSession) SetTreasury(treasury common.Address) (*types.Transaction, error) {
	return _Config.Contract.SetTreasury(&_Config.TransactOpts, treasury)
}

// SetTreasuryFeeFactor is a paid mutator transaction binding the contract method 0x30a767c5.
//
// Solidity: function setTreasuryFeeFactor(address token0, address token1, uint256 feeFactor) returns()
func (_Config *ConfigTransactor) SetTreasuryFeeFactor(opts *bind.TransactOpts, token0 common.Address, token1 common.Address, feeFactor *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setTreasuryFeeFactor", token0, token1, feeFactor)
}

// SetTreasuryFeeFactor is a paid mutator transaction binding the contract method 0x30a767c5.
//
// Solidity: function setTreasuryFeeFactor(address token0, address token1, uint256 feeFactor) returns()
func (_Config *ConfigSession) SetTreasuryFeeFactor(token0 common.Address, token1 common.Address, feeFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetTreasuryFeeFactor(&_Config.TransactOpts, token0, token1, feeFactor)
}

// SetTreasuryFeeFactor is a paid mutator transaction binding the contract method 0x30a767c5.
//
// Solidity: function setTreasuryFeeFactor(address token0, address token1, uint256 feeFactor) returns()
func (_Config *ConfigTransactorSession) SetTreasuryFeeFactor(token0 common.Address, token1 common.Address, feeFactor *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetTreasuryFeeFactor(&_Config.TransactOpts, token0, token1, feeFactor)
}

// SetTwapPeriod is a paid mutator transaction binding the contract method 0x5e657adf.
//
// Solidity: function setTwapPeriod(uint256 period) returns()
func (_Config *ConfigTransactor) SetTwapPeriod(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Config.contract.Transact(opts, "setTwapPeriod", period)
}

// SetTwapPeriod is a paid mutator transaction binding the contract method 0x5e657adf.
//
// Solidity: function setTwapPeriod(uint256 period) returns()
func (_Config *ConfigSession) SetTwapPeriod(period *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetTwapPeriod(&_Config.TransactOpts, period)
}

// SetTwapPeriod is a paid mutator transaction binding the contract method 0x5e657adf.
//
// Solidity: function setTwapPeriod(uint256 period) returns()
func (_Config *ConfigTransactorSession) SetTwapPeriod(period *big.Int) (*types.Transaction, error) {
	return _Config.Contract.SetTwapPeriod(&_Config.TransactOpts, period)
}
