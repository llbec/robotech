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

// RoleMetaData contains all meta data concerning the Role contract.
var RoleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CONFIG_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CONTROLLER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_DISTRIBUTION_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOV_TOKEN_CONTROLLER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIMITED_CONFIG_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIQUIDATION_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POOL_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRICING_KEEPER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROLE_ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROUTER_PLUGIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMELOCK_ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMELOCK_MULTISIG\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RoleABI is the input ABI used to generate the binding from.
// Deprecated: Use RoleMetaData.ABI instead.
var RoleABI = RoleMetaData.ABI

// Role is an auto generated Go binding around an Ethereum contract.
type Role struct {
	RoleCaller     // Read-only binding to the contract
	RoleTransactor // Write-only binding to the contract
	RoleFilterer   // Log filterer for contract events
}

// RoleCaller is an auto generated read-only Go binding around an Ethereum contract.
type RoleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RoleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RoleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RoleSession struct {
	Contract     *Role             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RoleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RoleCallerSession struct {
	Contract *RoleCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RoleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RoleTransactorSession struct {
	Contract     *RoleTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RoleRaw is an auto generated low-level Go binding around an Ethereum contract.
type RoleRaw struct {
	Contract *Role // Generic contract binding to access the raw methods on
}

// RoleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RoleCallerRaw struct {
	Contract *RoleCaller // Generic read-only contract binding to access the raw methods on
}

// RoleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RoleTransactorRaw struct {
	Contract *RoleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRole creates a new instance of Role, bound to a specific deployed contract.
func NewRole(address common.Address, backend bind.ContractBackend) (*Role, error) {
	contract, err := bindRole(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Role{RoleCaller: RoleCaller{contract: contract}, RoleTransactor: RoleTransactor{contract: contract}, RoleFilterer: RoleFilterer{contract: contract}}, nil
}

// NewRoleCaller creates a new read-only instance of Role, bound to a specific deployed contract.
func NewRoleCaller(address common.Address, caller bind.ContractCaller) (*RoleCaller, error) {
	contract, err := bindRole(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RoleCaller{contract: contract}, nil
}

// NewRoleTransactor creates a new write-only instance of Role, bound to a specific deployed contract.
func NewRoleTransactor(address common.Address, transactor bind.ContractTransactor) (*RoleTransactor, error) {
	contract, err := bindRole(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RoleTransactor{contract: contract}, nil
}

// NewRoleFilterer creates a new log filterer instance of Role, bound to a specific deployed contract.
func NewRoleFilterer(address common.Address, filterer bind.ContractFilterer) (*RoleFilterer, error) {
	contract, err := bindRole(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RoleFilterer{contract: contract}, nil
}

// bindRole binds a generic wrapper to an already deployed contract.
func bindRole(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RoleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Role *RoleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Role.Contract.RoleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Role *RoleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Role.Contract.RoleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Role *RoleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Role.Contract.RoleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Role *RoleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Role.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Role *RoleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Role.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Role *RoleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Role.Contract.contract.Transact(opts, method, params...)
}

// CONFIGKEEPER is a free data retrieval call binding the contract method 0xc66aade1.
//
// Solidity: function CONFIG_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) CONFIGKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "CONFIG_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CONFIGKEEPER is a free data retrieval call binding the contract method 0xc66aade1.
//
// Solidity: function CONFIG_KEEPER() view returns(bytes32)
func (_Role *RoleSession) CONFIGKEEPER() ([32]byte, error) {
	return _Role.Contract.CONFIGKEEPER(&_Role.CallOpts)
}

// CONFIGKEEPER is a free data retrieval call binding the contract method 0xc66aade1.
//
// Solidity: function CONFIG_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) CONFIGKEEPER() ([32]byte, error) {
	return _Role.Contract.CONFIGKEEPER(&_Role.CallOpts)
}

// CONTROLLER is a free data retrieval call binding the contract method 0xee0fc121.
//
// Solidity: function CONTROLLER() view returns(bytes32)
func (_Role *RoleCaller) CONTROLLER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "CONTROLLER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CONTROLLER is a free data retrieval call binding the contract method 0xee0fc121.
//
// Solidity: function CONTROLLER() view returns(bytes32)
func (_Role *RoleSession) CONTROLLER() ([32]byte, error) {
	return _Role.Contract.CONTROLLER(&_Role.CallOpts)
}

// CONTROLLER is a free data retrieval call binding the contract method 0xee0fc121.
//
// Solidity: function CONTROLLER() view returns(bytes32)
func (_Role *RoleCallerSession) CONTROLLER() ([32]byte, error) {
	return _Role.Contract.CONTROLLER(&_Role.CallOpts)
}

// FEEDISTRIBUTIONKEEPER is a free data retrieval call binding the contract method 0x75d3adbb.
//
// Solidity: function FEE_DISTRIBUTION_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) FEEDISTRIBUTIONKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "FEE_DISTRIBUTION_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FEEDISTRIBUTIONKEEPER is a free data retrieval call binding the contract method 0x75d3adbb.
//
// Solidity: function FEE_DISTRIBUTION_KEEPER() view returns(bytes32)
func (_Role *RoleSession) FEEDISTRIBUTIONKEEPER() ([32]byte, error) {
	return _Role.Contract.FEEDISTRIBUTIONKEEPER(&_Role.CallOpts)
}

// FEEDISTRIBUTIONKEEPER is a free data retrieval call binding the contract method 0x75d3adbb.
//
// Solidity: function FEE_DISTRIBUTION_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) FEEDISTRIBUTIONKEEPER() ([32]byte, error) {
	return _Role.Contract.FEEDISTRIBUTIONKEEPER(&_Role.CallOpts)
}

// FEEKEEPER is a free data retrieval call binding the contract method 0xcf7fb1d0.
//
// Solidity: function FEE_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) FEEKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "FEE_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FEEKEEPER is a free data retrieval call binding the contract method 0xcf7fb1d0.
//
// Solidity: function FEE_KEEPER() view returns(bytes32)
func (_Role *RoleSession) FEEKEEPER() ([32]byte, error) {
	return _Role.Contract.FEEKEEPER(&_Role.CallOpts)
}

// FEEKEEPER is a free data retrieval call binding the contract method 0xcf7fb1d0.
//
// Solidity: function FEE_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) FEEKEEPER() ([32]byte, error) {
	return _Role.Contract.FEEKEEPER(&_Role.CallOpts)
}

// GOVTOKENCONTROLLER is a free data retrieval call binding the contract method 0xe0fde20c.
//
// Solidity: function GOV_TOKEN_CONTROLLER() view returns(bytes32)
func (_Role *RoleCaller) GOVTOKENCONTROLLER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "GOV_TOKEN_CONTROLLER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GOVTOKENCONTROLLER is a free data retrieval call binding the contract method 0xe0fde20c.
//
// Solidity: function GOV_TOKEN_CONTROLLER() view returns(bytes32)
func (_Role *RoleSession) GOVTOKENCONTROLLER() ([32]byte, error) {
	return _Role.Contract.GOVTOKENCONTROLLER(&_Role.CallOpts)
}

// GOVTOKENCONTROLLER is a free data retrieval call binding the contract method 0xe0fde20c.
//
// Solidity: function GOV_TOKEN_CONTROLLER() view returns(bytes32)
func (_Role *RoleCallerSession) GOVTOKENCONTROLLER() ([32]byte, error) {
	return _Role.Contract.GOVTOKENCONTROLLER(&_Role.CallOpts)
}

// LIMITEDCONFIGKEEPER is a free data retrieval call binding the contract method 0x774fb4b8.
//
// Solidity: function LIMITED_CONFIG_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) LIMITEDCONFIGKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "LIMITED_CONFIG_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LIMITEDCONFIGKEEPER is a free data retrieval call binding the contract method 0x774fb4b8.
//
// Solidity: function LIMITED_CONFIG_KEEPER() view returns(bytes32)
func (_Role *RoleSession) LIMITEDCONFIGKEEPER() ([32]byte, error) {
	return _Role.Contract.LIMITEDCONFIGKEEPER(&_Role.CallOpts)
}

// LIMITEDCONFIGKEEPER is a free data retrieval call binding the contract method 0x774fb4b8.
//
// Solidity: function LIMITED_CONFIG_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) LIMITEDCONFIGKEEPER() ([32]byte, error) {
	return _Role.Contract.LIMITEDCONFIGKEEPER(&_Role.CallOpts)
}

// LIQUIDATIONKEEPER is a free data retrieval call binding the contract method 0xf13c5a4b.
//
// Solidity: function LIQUIDATION_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) LIQUIDATIONKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "LIQUIDATION_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LIQUIDATIONKEEPER is a free data retrieval call binding the contract method 0xf13c5a4b.
//
// Solidity: function LIQUIDATION_KEEPER() view returns(bytes32)
func (_Role *RoleSession) LIQUIDATIONKEEPER() ([32]byte, error) {
	return _Role.Contract.LIQUIDATIONKEEPER(&_Role.CallOpts)
}

// LIQUIDATIONKEEPER is a free data retrieval call binding the contract method 0xf13c5a4b.
//
// Solidity: function LIQUIDATION_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) LIQUIDATIONKEEPER() ([32]byte, error) {
	return _Role.Contract.LIQUIDATIONKEEPER(&_Role.CallOpts)
}

// POOLKEEPER is a free data retrieval call binding the contract method 0xe629d48c.
//
// Solidity: function POOL_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) POOLKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "POOL_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// POOLKEEPER is a free data retrieval call binding the contract method 0xe629d48c.
//
// Solidity: function POOL_KEEPER() view returns(bytes32)
func (_Role *RoleSession) POOLKEEPER() ([32]byte, error) {
	return _Role.Contract.POOLKEEPER(&_Role.CallOpts)
}

// POOLKEEPER is a free data retrieval call binding the contract method 0xe629d48c.
//
// Solidity: function POOL_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) POOLKEEPER() ([32]byte, error) {
	return _Role.Contract.POOLKEEPER(&_Role.CallOpts)
}

// PRICINGKEEPER is a free data retrieval call binding the contract method 0x9ecff617.
//
// Solidity: function PRICING_KEEPER() view returns(bytes32)
func (_Role *RoleCaller) PRICINGKEEPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "PRICING_KEEPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PRICINGKEEPER is a free data retrieval call binding the contract method 0x9ecff617.
//
// Solidity: function PRICING_KEEPER() view returns(bytes32)
func (_Role *RoleSession) PRICINGKEEPER() ([32]byte, error) {
	return _Role.Contract.PRICINGKEEPER(&_Role.CallOpts)
}

// PRICINGKEEPER is a free data retrieval call binding the contract method 0x9ecff617.
//
// Solidity: function PRICING_KEEPER() view returns(bytes32)
func (_Role *RoleCallerSession) PRICINGKEEPER() ([32]byte, error) {
	return _Role.Contract.PRICINGKEEPER(&_Role.CallOpts)
}

// ROLEADMIN is a free data retrieval call binding the contract method 0xd391014b.
//
// Solidity: function ROLE_ADMIN() view returns(bytes32)
func (_Role *RoleCaller) ROLEADMIN(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "ROLE_ADMIN")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ROLEADMIN is a free data retrieval call binding the contract method 0xd391014b.
//
// Solidity: function ROLE_ADMIN() view returns(bytes32)
func (_Role *RoleSession) ROLEADMIN() ([32]byte, error) {
	return _Role.Contract.ROLEADMIN(&_Role.CallOpts)
}

// ROLEADMIN is a free data retrieval call binding the contract method 0xd391014b.
//
// Solidity: function ROLE_ADMIN() view returns(bytes32)
func (_Role *RoleCallerSession) ROLEADMIN() ([32]byte, error) {
	return _Role.Contract.ROLEADMIN(&_Role.CallOpts)
}

// ROUTERPLUGIN is a free data retrieval call binding the contract method 0x9b8b49f8.
//
// Solidity: function ROUTER_PLUGIN() view returns(bytes32)
func (_Role *RoleCaller) ROUTERPLUGIN(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "ROUTER_PLUGIN")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ROUTERPLUGIN is a free data retrieval call binding the contract method 0x9b8b49f8.
//
// Solidity: function ROUTER_PLUGIN() view returns(bytes32)
func (_Role *RoleSession) ROUTERPLUGIN() ([32]byte, error) {
	return _Role.Contract.ROUTERPLUGIN(&_Role.CallOpts)
}

// ROUTERPLUGIN is a free data retrieval call binding the contract method 0x9b8b49f8.
//
// Solidity: function ROUTER_PLUGIN() view returns(bytes32)
func (_Role *RoleCallerSession) ROUTERPLUGIN() ([32]byte, error) {
	return _Role.Contract.ROUTERPLUGIN(&_Role.CallOpts)
}

// TIMELOCKADMIN is a free data retrieval call binding the contract method 0xe2ff47b3.
//
// Solidity: function TIMELOCK_ADMIN() view returns(bytes32)
func (_Role *RoleCaller) TIMELOCKADMIN(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "TIMELOCK_ADMIN")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TIMELOCKADMIN is a free data retrieval call binding the contract method 0xe2ff47b3.
//
// Solidity: function TIMELOCK_ADMIN() view returns(bytes32)
func (_Role *RoleSession) TIMELOCKADMIN() ([32]byte, error) {
	return _Role.Contract.TIMELOCKADMIN(&_Role.CallOpts)
}

// TIMELOCKADMIN is a free data retrieval call binding the contract method 0xe2ff47b3.
//
// Solidity: function TIMELOCK_ADMIN() view returns(bytes32)
func (_Role *RoleCallerSession) TIMELOCKADMIN() ([32]byte, error) {
	return _Role.Contract.TIMELOCKADMIN(&_Role.CallOpts)
}

// TIMELOCKMULTISIG is a free data retrieval call binding the contract method 0x4479d97b.
//
// Solidity: function TIMELOCK_MULTISIG() view returns(bytes32)
func (_Role *RoleCaller) TIMELOCKMULTISIG(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Role.contract.Call(opts, &out, "TIMELOCK_MULTISIG")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TIMELOCKMULTISIG is a free data retrieval call binding the contract method 0x4479d97b.
//
// Solidity: function TIMELOCK_MULTISIG() view returns(bytes32)
func (_Role *RoleSession) TIMELOCKMULTISIG() ([32]byte, error) {
	return _Role.Contract.TIMELOCKMULTISIG(&_Role.CallOpts)
}

// TIMELOCKMULTISIG is a free data retrieval call binding the contract method 0x4479d97b.
//
// Solidity: function TIMELOCK_MULTISIG() view returns(bytes32)
func (_Role *RoleCallerSession) TIMELOCKMULTISIG() ([32]byte, error) {
	return _Role.Contract.TIMELOCKMULTISIG(&_Role.CallOpts)
}
