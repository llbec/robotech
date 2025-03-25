// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package role

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

// RoleStoreMetaData contains all meta data concerning the RoleStore contract.
var RoleStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ThereMustBeAtLeastOneRoleAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ThereMustBeAtLeastOneTimelockMultiSig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"role\",\"type\":\"string\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"getRoleCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"roleKey\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"roleKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getRoleMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getRoles\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"roleKey\",\"type\":\"bytes32\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"roleKey\",\"type\":\"bytes32\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"roleKey\",\"type\":\"bytes32\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RoleStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use RoleStoreMetaData.ABI instead.
var RoleStoreABI = RoleStoreMetaData.ABI

// RoleStore is an auto generated Go binding around an Ethereum contract.
type RoleStore struct {
	RoleStoreCaller     // Read-only binding to the contract
	RoleStoreTransactor // Write-only binding to the contract
	RoleStoreFilterer   // Log filterer for contract events
}

// RoleStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type RoleStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RoleStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RoleStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RoleStoreSession struct {
	Contract     *RoleStore        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RoleStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RoleStoreCallerSession struct {
	Contract *RoleStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RoleStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RoleStoreTransactorSession struct {
	Contract     *RoleStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RoleStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type RoleStoreRaw struct {
	Contract *RoleStore // Generic contract binding to access the raw methods on
}

// RoleStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RoleStoreCallerRaw struct {
	Contract *RoleStoreCaller // Generic read-only contract binding to access the raw methods on
}

// RoleStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RoleStoreTransactorRaw struct {
	Contract *RoleStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRoleStore creates a new instance of RoleStore, bound to a specific deployed contract.
func NewRoleStore(address common.Address, backend bind.ContractBackend) (*RoleStore, error) {
	contract, err := bindRoleStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RoleStore{RoleStoreCaller: RoleStoreCaller{contract: contract}, RoleStoreTransactor: RoleStoreTransactor{contract: contract}, RoleStoreFilterer: RoleStoreFilterer{contract: contract}}, nil
}

// NewRoleStoreCaller creates a new read-only instance of RoleStore, bound to a specific deployed contract.
func NewRoleStoreCaller(address common.Address, caller bind.ContractCaller) (*RoleStoreCaller, error) {
	contract, err := bindRoleStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RoleStoreCaller{contract: contract}, nil
}

// NewRoleStoreTransactor creates a new write-only instance of RoleStore, bound to a specific deployed contract.
func NewRoleStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*RoleStoreTransactor, error) {
	contract, err := bindRoleStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RoleStoreTransactor{contract: contract}, nil
}

// NewRoleStoreFilterer creates a new log filterer instance of RoleStore, bound to a specific deployed contract.
func NewRoleStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*RoleStoreFilterer, error) {
	contract, err := bindRoleStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RoleStoreFilterer{contract: contract}, nil
}

// bindRoleStore binds a generic wrapper to an already deployed contract.
func bindRoleStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RoleStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RoleStore *RoleStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RoleStore.Contract.RoleStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RoleStore *RoleStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoleStore.Contract.RoleStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RoleStore *RoleStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RoleStore.Contract.RoleStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RoleStore *RoleStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RoleStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RoleStore *RoleStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoleStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RoleStore *RoleStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RoleStore.Contract.contract.Transact(opts, method, params...)
}

// GetRoleCount is a free data retrieval call binding the contract method 0x83d33319.
//
// Solidity: function getRoleCount() view returns(uint256)
func (_RoleStore *RoleStoreCaller) GetRoleCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RoleStore.contract.Call(opts, &out, "getRoleCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleCount is a free data retrieval call binding the contract method 0x83d33319.
//
// Solidity: function getRoleCount() view returns(uint256)
func (_RoleStore *RoleStoreSession) GetRoleCount() (*big.Int, error) {
	return _RoleStore.Contract.GetRoleCount(&_RoleStore.CallOpts)
}

// GetRoleCount is a free data retrieval call binding the contract method 0x83d33319.
//
// Solidity: function getRoleCount() view returns(uint256)
func (_RoleStore *RoleStoreCallerSession) GetRoleCount() (*big.Int, error) {
	return _RoleStore.Contract.GetRoleCount(&_RoleStore.CallOpts)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 roleKey) view returns(uint256)
func (_RoleStore *RoleStoreCaller) GetRoleMemberCount(opts *bind.CallOpts, roleKey [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _RoleStore.contract.Call(opts, &out, "getRoleMemberCount", roleKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 roleKey) view returns(uint256)
func (_RoleStore *RoleStoreSession) GetRoleMemberCount(roleKey [32]byte) (*big.Int, error) {
	return _RoleStore.Contract.GetRoleMemberCount(&_RoleStore.CallOpts, roleKey)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 roleKey) view returns(uint256)
func (_RoleStore *RoleStoreCallerSession) GetRoleMemberCount(roleKey [32]byte) (*big.Int, error) {
	return _RoleStore.Contract.GetRoleMemberCount(&_RoleStore.CallOpts, roleKey)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0x2a861f57.
//
// Solidity: function getRoleMembers(bytes32 roleKey, uint256 start, uint256 end) view returns(address[])
func (_RoleStore *RoleStoreCaller) GetRoleMembers(opts *bind.CallOpts, roleKey [32]byte, start *big.Int, end *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _RoleStore.contract.Call(opts, &out, "getRoleMembers", roleKey, start, end)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRoleMembers is a free data retrieval call binding the contract method 0x2a861f57.
//
// Solidity: function getRoleMembers(bytes32 roleKey, uint256 start, uint256 end) view returns(address[])
func (_RoleStore *RoleStoreSession) GetRoleMembers(roleKey [32]byte, start *big.Int, end *big.Int) ([]common.Address, error) {
	return _RoleStore.Contract.GetRoleMembers(&_RoleStore.CallOpts, roleKey, start, end)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0x2a861f57.
//
// Solidity: function getRoleMembers(bytes32 roleKey, uint256 start, uint256 end) view returns(address[])
func (_RoleStore *RoleStoreCallerSession) GetRoleMembers(roleKey [32]byte, start *big.Int, end *big.Int) ([]common.Address, error) {
	return _RoleStore.Contract.GetRoleMembers(&_RoleStore.CallOpts, roleKey, start, end)
}

// GetRoles is a free data retrieval call binding the contract method 0x821c1898.
//
// Solidity: function getRoles(uint256 start, uint256 end) view returns(bytes32[])
func (_RoleStore *RoleStoreCaller) GetRoles(opts *bind.CallOpts, start *big.Int, end *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _RoleStore.contract.Call(opts, &out, "getRoles", start, end)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetRoles is a free data retrieval call binding the contract method 0x821c1898.
//
// Solidity: function getRoles(uint256 start, uint256 end) view returns(bytes32[])
func (_RoleStore *RoleStoreSession) GetRoles(start *big.Int, end *big.Int) ([][32]byte, error) {
	return _RoleStore.Contract.GetRoles(&_RoleStore.CallOpts, start, end)
}

// GetRoles is a free data retrieval call binding the contract method 0x821c1898.
//
// Solidity: function getRoles(uint256 start, uint256 end) view returns(bytes32[])
func (_RoleStore *RoleStoreCallerSession) GetRoles(start *big.Int, end *big.Int) ([][32]byte, error) {
	return _RoleStore.Contract.GetRoles(&_RoleStore.CallOpts, start, end)
}

// HasRole is a free data retrieval call binding the contract method 0xac4ab3fb.
//
// Solidity: function hasRole(address account, bytes32 roleKey) view returns(bool)
func (_RoleStore *RoleStoreCaller) HasRole(opts *bind.CallOpts, account common.Address, roleKey [32]byte) (bool, error) {
	var out []interface{}
	err := _RoleStore.contract.Call(opts, &out, "hasRole", account, roleKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0xac4ab3fb.
//
// Solidity: function hasRole(address account, bytes32 roleKey) view returns(bool)
func (_RoleStore *RoleStoreSession) HasRole(account common.Address, roleKey [32]byte) (bool, error) {
	return _RoleStore.Contract.HasRole(&_RoleStore.CallOpts, account, roleKey)
}

// HasRole is a free data retrieval call binding the contract method 0xac4ab3fb.
//
// Solidity: function hasRole(address account, bytes32 roleKey) view returns(bool)
func (_RoleStore *RoleStoreCallerSession) HasRole(account common.Address, roleKey [32]byte) (bool, error) {
	return _RoleStore.Contract.HasRole(&_RoleStore.CallOpts, account, roleKey)
}

// GrantRole is a paid mutator transaction binding the contract method 0xab2742dc.
//
// Solidity: function grantRole(address account, bytes32 roleKey) returns()
func (_RoleStore *RoleStoreTransactor) GrantRole(opts *bind.TransactOpts, account common.Address, roleKey [32]byte) (*types.Transaction, error) {
	return _RoleStore.contract.Transact(opts, "grantRole", account, roleKey)
}

// GrantRole is a paid mutator transaction binding the contract method 0xab2742dc.
//
// Solidity: function grantRole(address account, bytes32 roleKey) returns()
func (_RoleStore *RoleStoreSession) GrantRole(account common.Address, roleKey [32]byte) (*types.Transaction, error) {
	return _RoleStore.Contract.GrantRole(&_RoleStore.TransactOpts, account, roleKey)
}

// GrantRole is a paid mutator transaction binding the contract method 0xab2742dc.
//
// Solidity: function grantRole(address account, bytes32 roleKey) returns()
func (_RoleStore *RoleStoreTransactorSession) GrantRole(account common.Address, roleKey [32]byte) (*types.Transaction, error) {
	return _RoleStore.Contract.GrantRole(&_RoleStore.TransactOpts, account, roleKey)
}

// RevokeRole is a paid mutator transaction binding the contract method 0x208dd1ff.
//
// Solidity: function revokeRole(address account, bytes32 roleKey) returns()
func (_RoleStore *RoleStoreTransactor) RevokeRole(opts *bind.TransactOpts, account common.Address, roleKey [32]byte) (*types.Transaction, error) {
	return _RoleStore.contract.Transact(opts, "revokeRole", account, roleKey)
}

// RevokeRole is a paid mutator transaction binding the contract method 0x208dd1ff.
//
// Solidity: function revokeRole(address account, bytes32 roleKey) returns()
func (_RoleStore *RoleStoreSession) RevokeRole(account common.Address, roleKey [32]byte) (*types.Transaction, error) {
	return _RoleStore.Contract.RevokeRole(&_RoleStore.TransactOpts, account, roleKey)
}

// RevokeRole is a paid mutator transaction binding the contract method 0x208dd1ff.
//
// Solidity: function revokeRole(address account, bytes32 roleKey) returns()
func (_RoleStore *RoleStoreTransactorSession) RevokeRole(account common.Address, roleKey [32]byte) (*types.Transaction, error) {
	return _RoleStore.Contract.RevokeRole(&_RoleStore.TransactOpts, account, roleKey)
}
