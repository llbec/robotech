// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package aave

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

// ATokenABI is the input ABI used to generate the binding from.
const ATokenABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"BalanceTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"aTokenDecimals\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"aTokenName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"aTokenSymbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ATOKEN_REVISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EIP712_REVISION\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POOL\",\"outputs\":[{\"internalType\":\"contractILendingPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESERVE_TREASURY_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNDERLYING_ASSET_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiverOfUnderlying\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getIncentivesController\",\"outputs\":[{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getScaledUserBalanceAndSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"handleRepayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractILendingPool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"aTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"aTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"aTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"mintToTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"scaledBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scaledTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferOnLiquidation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferUnderlyingTo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AToken is an auto generated Go binding around an Ethereum contract.
type AToken struct {
	ATokenCaller     // Read-only binding to the contract
	ATokenTransactor // Write-only binding to the contract
	ATokenFilterer   // Log filterer for contract events
}

// ATokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type ATokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ATokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ATokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ATokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ATokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ATokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ATokenSession struct {
	Contract     *AToken           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ATokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ATokenCallerSession struct {
	Contract *ATokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ATokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ATokenTransactorSession struct {
	Contract     *ATokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ATokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type ATokenRaw struct {
	Contract *AToken // Generic contract binding to access the raw methods on
}

// ATokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ATokenCallerRaw struct {
	Contract *ATokenCaller // Generic read-only contract binding to access the raw methods on
}

// ATokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ATokenTransactorRaw struct {
	Contract *ATokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAToken creates a new instance of AToken, bound to a specific deployed contract.
func NewAToken(address common.Address, backend bind.ContractBackend) (*AToken, error) {
	contract, err := bindAToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AToken{ATokenCaller: ATokenCaller{contract: contract}, ATokenTransactor: ATokenTransactor{contract: contract}, ATokenFilterer: ATokenFilterer{contract: contract}}, nil
}

// NewATokenCaller creates a new read-only instance of AToken, bound to a specific deployed contract.
func NewATokenCaller(address common.Address, caller bind.ContractCaller) (*ATokenCaller, error) {
	contract, err := bindAToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ATokenCaller{contract: contract}, nil
}

// NewATokenTransactor creates a new write-only instance of AToken, bound to a specific deployed contract.
func NewATokenTransactor(address common.Address, transactor bind.ContractTransactor) (*ATokenTransactor, error) {
	contract, err := bindAToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ATokenTransactor{contract: contract}, nil
}

// NewATokenFilterer creates a new log filterer instance of AToken, bound to a specific deployed contract.
func NewATokenFilterer(address common.Address, filterer bind.ContractFilterer) (*ATokenFilterer, error) {
	contract, err := bindAToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ATokenFilterer{contract: contract}, nil
}

// bindAToken binds a generic wrapper to an already deployed contract.
func bindAToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ATokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AToken *ATokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AToken.Contract.ATokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AToken *ATokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AToken.Contract.ATokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AToken *ATokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AToken.Contract.ATokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AToken *ATokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AToken *ATokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AToken *ATokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AToken.Contract.contract.Transact(opts, method, params...)
}

// ATOKENREVISION is a free data retrieval call binding the contract method 0x0bd7ad3b.
//
// Solidity: function ATOKEN_REVISION() view returns(uint256)
func (_AToken *ATokenCaller) ATOKENREVISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "ATOKEN_REVISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ATOKENREVISION is a free data retrieval call binding the contract method 0x0bd7ad3b.
//
// Solidity: function ATOKEN_REVISION() view returns(uint256)
func (_AToken *ATokenSession) ATOKENREVISION() (*big.Int, error) {
	return _AToken.Contract.ATOKENREVISION(&_AToken.CallOpts)
}

// ATOKENREVISION is a free data retrieval call binding the contract method 0x0bd7ad3b.
//
// Solidity: function ATOKEN_REVISION() view returns(uint256)
func (_AToken *ATokenCallerSession) ATOKENREVISION() (*big.Int, error) {
	return _AToken.Contract.ATOKENREVISION(&_AToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_AToken *ATokenCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_AToken *ATokenSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _AToken.Contract.DOMAINSEPARATOR(&_AToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_AToken *ATokenCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _AToken.Contract.DOMAINSEPARATOR(&_AToken.CallOpts)
}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_AToken *ATokenCaller) EIP712REVISION(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "EIP712_REVISION")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_AToken *ATokenSession) EIP712REVISION() ([]byte, error) {
	return _AToken.Contract.EIP712REVISION(&_AToken.CallOpts)
}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_AToken *ATokenCallerSession) EIP712REVISION() ([]byte, error) {
	return _AToken.Contract.EIP712REVISION(&_AToken.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_AToken *ATokenCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "PERMIT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_AToken *ATokenSession) PERMITTYPEHASH() ([32]byte, error) {
	return _AToken.Contract.PERMITTYPEHASH(&_AToken.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_AToken *ATokenCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _AToken.Contract.PERMITTYPEHASH(&_AToken.CallOpts)
}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_AToken *ATokenCaller) POOL(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "POOL")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_AToken *ATokenSession) POOL() (common.Address, error) {
	return _AToken.Contract.POOL(&_AToken.CallOpts)
}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_AToken *ATokenCallerSession) POOL() (common.Address, error) {
	return _AToken.Contract.POOL(&_AToken.CallOpts)
}

// RESERVETREASURYADDRESS is a free data retrieval call binding the contract method 0xae167335.
//
// Solidity: function RESERVE_TREASURY_ADDRESS() view returns(address)
func (_AToken *ATokenCaller) RESERVETREASURYADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "RESERVE_TREASURY_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RESERVETREASURYADDRESS is a free data retrieval call binding the contract method 0xae167335.
//
// Solidity: function RESERVE_TREASURY_ADDRESS() view returns(address)
func (_AToken *ATokenSession) RESERVETREASURYADDRESS() (common.Address, error) {
	return _AToken.Contract.RESERVETREASURYADDRESS(&_AToken.CallOpts)
}

// RESERVETREASURYADDRESS is a free data retrieval call binding the contract method 0xae167335.
//
// Solidity: function RESERVE_TREASURY_ADDRESS() view returns(address)
func (_AToken *ATokenCallerSession) RESERVETREASURYADDRESS() (common.Address, error) {
	return _AToken.Contract.RESERVETREASURYADDRESS(&_AToken.CallOpts)
}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_AToken *ATokenCaller) UNDERLYINGASSETADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "UNDERLYING_ASSET_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_AToken *ATokenSession) UNDERLYINGASSETADDRESS() (common.Address, error) {
	return _AToken.Contract.UNDERLYINGASSETADDRESS(&_AToken.CallOpts)
}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_AToken *ATokenCallerSession) UNDERLYINGASSETADDRESS() (common.Address, error) {
	return _AToken.Contract.UNDERLYINGASSETADDRESS(&_AToken.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0xb9844d8d.
//
// Solidity: function _nonces(address ) view returns(uint256)
func (_AToken *ATokenCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "_nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0xb9844d8d.
//
// Solidity: function _nonces(address ) view returns(uint256)
func (_AToken *ATokenSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AToken.Contract.Nonces(&_AToken.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0xb9844d8d.
//
// Solidity: function _nonces(address ) view returns(uint256)
func (_AToken *ATokenCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AToken.Contract.Nonces(&_AToken.CallOpts, arg0)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AToken *ATokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AToken *ATokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AToken.Contract.Allowance(&_AToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AToken *ATokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AToken.Contract.Allowance(&_AToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_AToken *ATokenCaller) BalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "balanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_AToken *ATokenSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _AToken.Contract.BalanceOf(&_AToken.CallOpts, user)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_AToken *ATokenCallerSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _AToken.Contract.BalanceOf(&_AToken.CallOpts, user)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AToken *ATokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AToken *ATokenSession) Decimals() (uint8, error) {
	return _AToken.Contract.Decimals(&_AToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AToken *ATokenCallerSession) Decimals() (uint8, error) {
	return _AToken.Contract.Decimals(&_AToken.CallOpts)
}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_AToken *ATokenCaller) GetIncentivesController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "getIncentivesController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_AToken *ATokenSession) GetIncentivesController() (common.Address, error) {
	return _AToken.Contract.GetIncentivesController(&_AToken.CallOpts)
}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_AToken *ATokenCallerSession) GetIncentivesController() (common.Address, error) {
	return _AToken.Contract.GetIncentivesController(&_AToken.CallOpts)
}

// GetScaledUserBalanceAndSupply is a free data retrieval call binding the contract method 0x0afbcdc9.
//
// Solidity: function getScaledUserBalanceAndSupply(address user) view returns(uint256, uint256)
func (_AToken *ATokenCaller) GetScaledUserBalanceAndSupply(opts *bind.CallOpts, user common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "getScaledUserBalanceAndSupply", user)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetScaledUserBalanceAndSupply is a free data retrieval call binding the contract method 0x0afbcdc9.
//
// Solidity: function getScaledUserBalanceAndSupply(address user) view returns(uint256, uint256)
func (_AToken *ATokenSession) GetScaledUserBalanceAndSupply(user common.Address) (*big.Int, *big.Int, error) {
	return _AToken.Contract.GetScaledUserBalanceAndSupply(&_AToken.CallOpts, user)
}

// GetScaledUserBalanceAndSupply is a free data retrieval call binding the contract method 0x0afbcdc9.
//
// Solidity: function getScaledUserBalanceAndSupply(address user) view returns(uint256, uint256)
func (_AToken *ATokenCallerSession) GetScaledUserBalanceAndSupply(user common.Address) (*big.Int, *big.Int, error) {
	return _AToken.Contract.GetScaledUserBalanceAndSupply(&_AToken.CallOpts, user)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AToken *ATokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AToken *ATokenSession) Name() (string, error) {
	return _AToken.Contract.Name(&_AToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AToken *ATokenCallerSession) Name() (string, error) {
	return _AToken.Contract.Name(&_AToken.CallOpts)
}

// ScaledBalanceOf is a free data retrieval call binding the contract method 0x1da24f3e.
//
// Solidity: function scaledBalanceOf(address user) view returns(uint256)
func (_AToken *ATokenCaller) ScaledBalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "scaledBalanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScaledBalanceOf is a free data retrieval call binding the contract method 0x1da24f3e.
//
// Solidity: function scaledBalanceOf(address user) view returns(uint256)
func (_AToken *ATokenSession) ScaledBalanceOf(user common.Address) (*big.Int, error) {
	return _AToken.Contract.ScaledBalanceOf(&_AToken.CallOpts, user)
}

// ScaledBalanceOf is a free data retrieval call binding the contract method 0x1da24f3e.
//
// Solidity: function scaledBalanceOf(address user) view returns(uint256)
func (_AToken *ATokenCallerSession) ScaledBalanceOf(user common.Address) (*big.Int, error) {
	return _AToken.Contract.ScaledBalanceOf(&_AToken.CallOpts, user)
}

// ScaledTotalSupply is a free data retrieval call binding the contract method 0xb1bf962d.
//
// Solidity: function scaledTotalSupply() view returns(uint256)
func (_AToken *ATokenCaller) ScaledTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "scaledTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScaledTotalSupply is a free data retrieval call binding the contract method 0xb1bf962d.
//
// Solidity: function scaledTotalSupply() view returns(uint256)
func (_AToken *ATokenSession) ScaledTotalSupply() (*big.Int, error) {
	return _AToken.Contract.ScaledTotalSupply(&_AToken.CallOpts)
}

// ScaledTotalSupply is a free data retrieval call binding the contract method 0xb1bf962d.
//
// Solidity: function scaledTotalSupply() view returns(uint256)
func (_AToken *ATokenCallerSession) ScaledTotalSupply() (*big.Int, error) {
	return _AToken.Contract.ScaledTotalSupply(&_AToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AToken *ATokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AToken *ATokenSession) Symbol() (string, error) {
	return _AToken.Contract.Symbol(&_AToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AToken *ATokenCallerSession) Symbol() (string, error) {
	return _AToken.Contract.Symbol(&_AToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AToken *ATokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AToken *ATokenSession) TotalSupply() (*big.Int, error) {
	return _AToken.Contract.TotalSupply(&_AToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AToken *ATokenCallerSession) TotalSupply() (*big.Int, error) {
	return _AToken.Contract.TotalSupply(&_AToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AToken *ATokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AToken *ATokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Approve(&_AToken.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AToken *ATokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Approve(&_AToken.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xd7020d0a.
//
// Solidity: function burn(address user, address receiverOfUnderlying, uint256 amount, uint256 index) returns()
func (_AToken *ATokenTransactor) Burn(opts *bind.TransactOpts, user common.Address, receiverOfUnderlying common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "burn", user, receiverOfUnderlying, amount, index)
}

// Burn is a paid mutator transaction binding the contract method 0xd7020d0a.
//
// Solidity: function burn(address user, address receiverOfUnderlying, uint256 amount, uint256 index) returns()
func (_AToken *ATokenSession) Burn(user common.Address, receiverOfUnderlying common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Burn(&_AToken.TransactOpts, user, receiverOfUnderlying, amount, index)
}

// Burn is a paid mutator transaction binding the contract method 0xd7020d0a.
//
// Solidity: function burn(address user, address receiverOfUnderlying, uint256 amount, uint256 index) returns()
func (_AToken *ATokenTransactorSession) Burn(user common.Address, receiverOfUnderlying common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Burn(&_AToken.TransactOpts, user, receiverOfUnderlying, amount, index)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_AToken *ATokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_AToken *ATokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.DecreaseAllowance(&_AToken.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_AToken *ATokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.DecreaseAllowance(&_AToken.TransactOpts, spender, subtractedValue)
}

// HandleRepayment is a paid mutator transaction binding the contract method 0x88dd91a1.
//
// Solidity: function handleRepayment(address user, uint256 amount) returns()
func (_AToken *ATokenTransactor) HandleRepayment(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "handleRepayment", user, amount)
}

// HandleRepayment is a paid mutator transaction binding the contract method 0x88dd91a1.
//
// Solidity: function handleRepayment(address user, uint256 amount) returns()
func (_AToken *ATokenSession) HandleRepayment(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.HandleRepayment(&_AToken.TransactOpts, user, amount)
}

// HandleRepayment is a paid mutator transaction binding the contract method 0x88dd91a1.
//
// Solidity: function handleRepayment(address user, uint256 amount) returns()
func (_AToken *ATokenTransactorSession) HandleRepayment(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.HandleRepayment(&_AToken.TransactOpts, user, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_AToken *ATokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_AToken *ATokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.IncreaseAllowance(&_AToken.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_AToken *ATokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.IncreaseAllowance(&_AToken.TransactOpts, spender, addedValue)
}

// Initialize is a paid mutator transaction binding the contract method 0x183fb413.
//
// Solidity: function initialize(address pool, address treasury, address underlyingAsset, address incentivesController, uint8 aTokenDecimals, string aTokenName, string aTokenSymbol, bytes params) returns()
func (_AToken *ATokenTransactor) Initialize(opts *bind.TransactOpts, pool common.Address, treasury common.Address, underlyingAsset common.Address, incentivesController common.Address, aTokenDecimals uint8, aTokenName string, aTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "initialize", pool, treasury, underlyingAsset, incentivesController, aTokenDecimals, aTokenName, aTokenSymbol, params)
}

// Initialize is a paid mutator transaction binding the contract method 0x183fb413.
//
// Solidity: function initialize(address pool, address treasury, address underlyingAsset, address incentivesController, uint8 aTokenDecimals, string aTokenName, string aTokenSymbol, bytes params) returns()
func (_AToken *ATokenSession) Initialize(pool common.Address, treasury common.Address, underlyingAsset common.Address, incentivesController common.Address, aTokenDecimals uint8, aTokenName string, aTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _AToken.Contract.Initialize(&_AToken.TransactOpts, pool, treasury, underlyingAsset, incentivesController, aTokenDecimals, aTokenName, aTokenSymbol, params)
}

// Initialize is a paid mutator transaction binding the contract method 0x183fb413.
//
// Solidity: function initialize(address pool, address treasury, address underlyingAsset, address incentivesController, uint8 aTokenDecimals, string aTokenName, string aTokenSymbol, bytes params) returns()
func (_AToken *ATokenTransactorSession) Initialize(pool common.Address, treasury common.Address, underlyingAsset common.Address, incentivesController common.Address, aTokenDecimals uint8, aTokenName string, aTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _AToken.Contract.Initialize(&_AToken.TransactOpts, pool, treasury, underlyingAsset, incentivesController, aTokenDecimals, aTokenName, aTokenSymbol, params)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address user, uint256 amount, uint256 index) returns(bool)
func (_AToken *ATokenTransactor) Mint(opts *bind.TransactOpts, user common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "mint", user, amount, index)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address user, uint256 amount, uint256 index) returns(bool)
func (_AToken *ATokenSession) Mint(user common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Mint(&_AToken.TransactOpts, user, amount, index)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address user, uint256 amount, uint256 index) returns(bool)
func (_AToken *ATokenTransactorSession) Mint(user common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Mint(&_AToken.TransactOpts, user, amount, index)
}

// MintToTreasury is a paid mutator transaction binding the contract method 0x7df5bd3b.
//
// Solidity: function mintToTreasury(uint256 amount, uint256 index) returns()
func (_AToken *ATokenTransactor) MintToTreasury(opts *bind.TransactOpts, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "mintToTreasury", amount, index)
}

// MintToTreasury is a paid mutator transaction binding the contract method 0x7df5bd3b.
//
// Solidity: function mintToTreasury(uint256 amount, uint256 index) returns()
func (_AToken *ATokenSession) MintToTreasury(amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.MintToTreasury(&_AToken.TransactOpts, amount, index)
}

// MintToTreasury is a paid mutator transaction binding the contract method 0x7df5bd3b.
//
// Solidity: function mintToTreasury(uint256 amount, uint256 index) returns()
func (_AToken *ATokenTransactorSession) MintToTreasury(amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.MintToTreasury(&_AToken.TransactOpts, amount, index)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AToken *ATokenTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AToken *ATokenSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AToken.Contract.Permit(&_AToken.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AToken *ATokenTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AToken.Contract.Permit(&_AToken.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_AToken *ATokenTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_AToken *ATokenSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Transfer(&_AToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_AToken *ATokenTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.Transfer(&_AToken.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_AToken *ATokenTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_AToken *ATokenSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.TransferFrom(&_AToken.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_AToken *ATokenTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.TransferFrom(&_AToken.TransactOpts, sender, recipient, amount)
}

// TransferOnLiquidation is a paid mutator transaction binding the contract method 0xf866c319.
//
// Solidity: function transferOnLiquidation(address from, address to, uint256 value) returns()
func (_AToken *ATokenTransactor) TransferOnLiquidation(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "transferOnLiquidation", from, to, value)
}

// TransferOnLiquidation is a paid mutator transaction binding the contract method 0xf866c319.
//
// Solidity: function transferOnLiquidation(address from, address to, uint256 value) returns()
func (_AToken *ATokenSession) TransferOnLiquidation(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.TransferOnLiquidation(&_AToken.TransactOpts, from, to, value)
}

// TransferOnLiquidation is a paid mutator transaction binding the contract method 0xf866c319.
//
// Solidity: function transferOnLiquidation(address from, address to, uint256 value) returns()
func (_AToken *ATokenTransactorSession) TransferOnLiquidation(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.TransferOnLiquidation(&_AToken.TransactOpts, from, to, value)
}

// TransferUnderlyingTo is a paid mutator transaction binding the contract method 0x4efecaa5.
//
// Solidity: function transferUnderlyingTo(address target, uint256 amount) returns(uint256)
func (_AToken *ATokenTransactor) TransferUnderlyingTo(opts *bind.TransactOpts, target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.contract.Transact(opts, "transferUnderlyingTo", target, amount)
}

// TransferUnderlyingTo is a paid mutator transaction binding the contract method 0x4efecaa5.
//
// Solidity: function transferUnderlyingTo(address target, uint256 amount) returns(uint256)
func (_AToken *ATokenSession) TransferUnderlyingTo(target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.TransferUnderlyingTo(&_AToken.TransactOpts, target, amount)
}

// TransferUnderlyingTo is a paid mutator transaction binding the contract method 0x4efecaa5.
//
// Solidity: function transferUnderlyingTo(address target, uint256 amount) returns(uint256)
func (_AToken *ATokenTransactorSession) TransferUnderlyingTo(target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AToken.Contract.TransferUnderlyingTo(&_AToken.TransactOpts, target, amount)
}

// ATokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AToken contract.
type ATokenApprovalIterator struct {
	Event *ATokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ATokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ATokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ATokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ATokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ATokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ATokenApproval represents a Approval event raised by the AToken contract.
type ATokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AToken *ATokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ATokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ATokenApprovalIterator{contract: _AToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AToken *ATokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ATokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ATokenApproval)
				if err := _AToken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AToken *ATokenFilterer) ParseApproval(log types.Log) (*ATokenApproval, error) {
	event := new(ATokenApproval)
	if err := _AToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ATokenBalanceTransferIterator is returned from FilterBalanceTransfer and is used to iterate over the raw logs and unpacked data for BalanceTransfer events raised by the AToken contract.
type ATokenBalanceTransferIterator struct {
	Event *ATokenBalanceTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ATokenBalanceTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ATokenBalanceTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ATokenBalanceTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ATokenBalanceTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ATokenBalanceTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ATokenBalanceTransfer represents a BalanceTransfer event raised by the AToken contract.
type ATokenBalanceTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Index *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterBalanceTransfer is a free log retrieval operation binding the contract event 0x4beccb90f994c31aced7a23b5611020728a23d8ec5cddd1a3e9d97b96fda8666.
//
// Solidity: event BalanceTransfer(address indexed from, address indexed to, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) FilterBalanceTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ATokenBalanceTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AToken.contract.FilterLogs(opts, "BalanceTransfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ATokenBalanceTransferIterator{contract: _AToken.contract, event: "BalanceTransfer", logs: logs, sub: sub}, nil
}

// WatchBalanceTransfer is a free log subscription operation binding the contract event 0x4beccb90f994c31aced7a23b5611020728a23d8ec5cddd1a3e9d97b96fda8666.
//
// Solidity: event BalanceTransfer(address indexed from, address indexed to, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) WatchBalanceTransfer(opts *bind.WatchOpts, sink chan<- *ATokenBalanceTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AToken.contract.WatchLogs(opts, "BalanceTransfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ATokenBalanceTransfer)
				if err := _AToken.contract.UnpackLog(event, "BalanceTransfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBalanceTransfer is a log parse operation binding the contract event 0x4beccb90f994c31aced7a23b5611020728a23d8ec5cddd1a3e9d97b96fda8666.
//
// Solidity: event BalanceTransfer(address indexed from, address indexed to, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) ParseBalanceTransfer(log types.Log) (*ATokenBalanceTransfer, error) {
	event := new(ATokenBalanceTransfer)
	if err := _AToken.contract.UnpackLog(event, "BalanceTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ATokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the AToken contract.
type ATokenBurnIterator struct {
	Event *ATokenBurn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ATokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ATokenBurn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ATokenBurn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ATokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ATokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ATokenBurn represents a Burn event raised by the AToken contract.
type ATokenBurn struct {
	From   common.Address
	Target common.Address
	Value  *big.Int
	Index  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2.
//
// Solidity: event Burn(address indexed from, address indexed target, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) FilterBurn(opts *bind.FilterOpts, from []common.Address, target []common.Address) (*ATokenBurnIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _AToken.contract.FilterLogs(opts, "Burn", fromRule, targetRule)
	if err != nil {
		return nil, err
	}
	return &ATokenBurnIterator{contract: _AToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2.
//
// Solidity: event Burn(address indexed from, address indexed target, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *ATokenBurn, from []common.Address, target []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _AToken.contract.WatchLogs(opts, "Burn", fromRule, targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ATokenBurn)
				if err := _AToken.contract.UnpackLog(event, "Burn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBurn is a log parse operation binding the contract event 0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2.
//
// Solidity: event Burn(address indexed from, address indexed target, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) ParseBurn(log types.Log) (*ATokenBurn, error) {
	event := new(ATokenBurn)
	if err := _AToken.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ATokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AToken contract.
type ATokenInitializedIterator struct {
	Event *ATokenInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ATokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ATokenInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ATokenInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ATokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ATokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ATokenInitialized represents a Initialized event raised by the AToken contract.
type ATokenInitialized struct {
	UnderlyingAsset      common.Address
	Pool                 common.Address
	Treasury             common.Address
	IncentivesController common.Address
	ATokenDecimals       uint8
	ATokenName           string
	ATokenSymbol         string
	Params               []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xb19e051f8af41150ccccb3fc2c2d8d15f4a4cf434f32a559ba75fe73d6eea20b.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address treasury, address incentivesController, uint8 aTokenDecimals, string aTokenName, string aTokenSymbol, bytes params)
func (_AToken *ATokenFilterer) FilterInitialized(opts *bind.FilterOpts, underlyingAsset []common.Address, pool []common.Address) (*ATokenInitializedIterator, error) {

	var underlyingAssetRule []interface{}
	for _, underlyingAssetItem := range underlyingAsset {
		underlyingAssetRule = append(underlyingAssetRule, underlyingAssetItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _AToken.contract.FilterLogs(opts, "Initialized", underlyingAssetRule, poolRule)
	if err != nil {
		return nil, err
	}
	return &ATokenInitializedIterator{contract: _AToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xb19e051f8af41150ccccb3fc2c2d8d15f4a4cf434f32a559ba75fe73d6eea20b.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address treasury, address incentivesController, uint8 aTokenDecimals, string aTokenName, string aTokenSymbol, bytes params)
func (_AToken *ATokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ATokenInitialized, underlyingAsset []common.Address, pool []common.Address) (event.Subscription, error) {

	var underlyingAssetRule []interface{}
	for _, underlyingAssetItem := range underlyingAsset {
		underlyingAssetRule = append(underlyingAssetRule, underlyingAssetItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _AToken.contract.WatchLogs(opts, "Initialized", underlyingAssetRule, poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ATokenInitialized)
				if err := _AToken.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xb19e051f8af41150ccccb3fc2c2d8d15f4a4cf434f32a559ba75fe73d6eea20b.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address treasury, address incentivesController, uint8 aTokenDecimals, string aTokenName, string aTokenSymbol, bytes params)
func (_AToken *ATokenFilterer) ParseInitialized(log types.Log) (*ATokenInitialized, error) {
	event := new(ATokenInitialized)
	if err := _AToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ATokenMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the AToken contract.
type ATokenMintIterator struct {
	Event *ATokenMint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ATokenMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ATokenMint)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ATokenMint)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ATokenMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ATokenMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ATokenMint represents a Mint event raised by the AToken contract.
type ATokenMint struct {
	From  common.Address
	Value *big.Int
	Index *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed from, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) FilterMint(opts *bind.FilterOpts, from []common.Address) (*ATokenMintIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AToken.contract.FilterLogs(opts, "Mint", fromRule)
	if err != nil {
		return nil, err
	}
	return &ATokenMintIterator{contract: _AToken.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed from, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *ATokenMint, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AToken.contract.WatchLogs(opts, "Mint", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ATokenMint)
				if err := _AToken.contract.UnpackLog(event, "Mint", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMint is a log parse operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed from, uint256 value, uint256 index)
func (_AToken *ATokenFilterer) ParseMint(log types.Log) (*ATokenMint, error) {
	event := new(ATokenMint)
	if err := _AToken.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ATokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AToken contract.
type ATokenTransferIterator struct {
	Event *ATokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ATokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ATokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ATokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ATokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ATokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ATokenTransfer represents a Transfer event raised by the AToken contract.
type ATokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AToken *ATokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ATokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ATokenTransferIterator{contract: _AToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AToken *ATokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ATokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ATokenTransfer)
				if err := _AToken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AToken *ATokenFilterer) ParseTransfer(log types.Log) (*ATokenTransfer, error) {
	event := new(ATokenTransfer)
	if err := _AToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
