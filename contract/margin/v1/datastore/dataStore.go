// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DataStore

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

// DataStoreMetaData contains all meta data concerning the DataStore contract.
var DataStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"EmptyBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyPositionHandle\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"SafeCastOverflowedIntToUint\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"role\",\"type\":\"string\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"name\":\"addAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"addBytes32\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"addUint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addressArrayValues\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"addressValues\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"applyBoundedDeltaToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"applyDeltaToInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"applyDeltaToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"applyDeltaToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"boolArrayValues\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"boolValues\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bytes32ArrayValues\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"bytes32Values\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"name\":\"containsAddress\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"containsBytes32\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"containsUint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"decrementInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"decrementUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getAddressArray\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"}],\"name\":\"getAddressCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getAddressValuesAt\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBackend\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getBoolArray\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getBytes32Array\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"}],\"name\":\"getBytes32Count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getBytes32ValuesAt\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getIntArray\",\"outputs\":[{\"internalType\":\"int256[]\",\"name\":\"\",\"type\":\"int256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPositionHandle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getStringArray\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getUintArray\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"}],\"name\":\"getUintCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"getUintValuesAt\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"incrementInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"incrementUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"_roleStore\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intArrayValues\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"intValues\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"name\":\"removeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeAddressArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeBoolArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"removeBytes32\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeBytes32\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeBytes32Array\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeInt\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeIntArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeString\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeStringArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeUint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"setKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"removeUint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"removeUintArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleStore\",\"outputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"name\":\"setAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"value\",\"type\":\"address[]\"}],\"name\":\"setAddressArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bool[]\",\"name\":\"value\",\"type\":\"bool[]\"}],\"name\":\"setBoolArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"setBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"value\",\"type\":\"bytes32[]\"}],\"name\":\"setBytes32Array\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"setInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"int256[]\",\"name\":\"value\",\"type\":\"int256[]\"}],\"name\":\"setIntArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"setString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"string[]\",\"name\":\"value\",\"type\":\"string[]\"}],\"name\":\"setStringArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"value\",\"type\":\"uint256[]\"}],\"name\":\"setUintArray\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stringArrayValues\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stringValues\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uintArrayValues\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"uintValues\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DataStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use DataStoreMetaData.ABI instead.
var DataStoreABI = DataStoreMetaData.ABI

// DataStore is an auto generated Go binding around an Ethereum contract.
type DataStore struct {
	DataStoreCaller     // Read-only binding to the contract
	DataStoreTransactor // Write-only binding to the contract
	DataStoreFilterer   // Log filterer for contract events
}

// DataStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataStoreSession struct {
	Contract     *DataStore        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DataStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataStoreCallerSession struct {
	Contract *DataStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DataStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataStoreTransactorSession struct {
	Contract     *DataStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DataStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataStoreRaw struct {
	Contract *DataStore // Generic contract binding to access the raw methods on
}

// DataStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataStoreCallerRaw struct {
	Contract *DataStoreCaller // Generic read-only contract binding to access the raw methods on
}

// DataStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataStoreTransactorRaw struct {
	Contract *DataStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataStore creates a new instance of DataStore, bound to a specific deployed contract.
func NewDataStore(address common.Address, backend bind.ContractBackend) (*DataStore, error) {
	contract, err := bindDataStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataStore{DataStoreCaller: DataStoreCaller{contract: contract}, DataStoreTransactor: DataStoreTransactor{contract: contract}, DataStoreFilterer: DataStoreFilterer{contract: contract}}, nil
}

// NewDataStoreCaller creates a new read-only instance of DataStore, bound to a specific deployed contract.
func NewDataStoreCaller(address common.Address, caller bind.ContractCaller) (*DataStoreCaller, error) {
	contract, err := bindDataStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataStoreCaller{contract: contract}, nil
}

// NewDataStoreTransactor creates a new write-only instance of DataStore, bound to a specific deployed contract.
func NewDataStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*DataStoreTransactor, error) {
	contract, err := bindDataStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataStoreTransactor{contract: contract}, nil
}

// NewDataStoreFilterer creates a new log filterer instance of DataStore, bound to a specific deployed contract.
func NewDataStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*DataStoreFilterer, error) {
	contract, err := bindDataStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataStoreFilterer{contract: contract}, nil
}

// bindDataStore binds a generic wrapper to an already deployed contract.
func bindDataStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataStore *DataStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataStore.Contract.DataStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataStore *DataStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataStore.Contract.DataStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataStore *DataStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataStore.Contract.DataStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataStore *DataStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataStore *DataStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataStore *DataStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataStore.Contract.contract.Transact(opts, method, params...)
}

// AddressArrayValues is a free data retrieval call binding the contract method 0x22f87464.
//
// Solidity: function addressArrayValues(bytes32 , uint256 ) view returns(address)
func (_DataStore *DataStoreCaller) AddressArrayValues(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "addressArrayValues", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AddressArrayValues is a free data retrieval call binding the contract method 0x22f87464.
//
// Solidity: function addressArrayValues(bytes32 , uint256 ) view returns(address)
func (_DataStore *DataStoreSession) AddressArrayValues(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _DataStore.Contract.AddressArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// AddressArrayValues is a free data retrieval call binding the contract method 0x22f87464.
//
// Solidity: function addressArrayValues(bytes32 , uint256 ) view returns(address)
func (_DataStore *DataStoreCallerSession) AddressArrayValues(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _DataStore.Contract.AddressArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// AddressValues is a free data retrieval call binding the contract method 0x22538dae.
//
// Solidity: function addressValues(bytes32 ) view returns(address)
func (_DataStore *DataStoreCaller) AddressValues(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "addressValues", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AddressValues is a free data retrieval call binding the contract method 0x22538dae.
//
// Solidity: function addressValues(bytes32 ) view returns(address)
func (_DataStore *DataStoreSession) AddressValues(arg0 [32]byte) (common.Address, error) {
	return _DataStore.Contract.AddressValues(&_DataStore.CallOpts, arg0)
}

// AddressValues is a free data retrieval call binding the contract method 0x22538dae.
//
// Solidity: function addressValues(bytes32 ) view returns(address)
func (_DataStore *DataStoreCallerSession) AddressValues(arg0 [32]byte) (common.Address, error) {
	return _DataStore.Contract.AddressValues(&_DataStore.CallOpts, arg0)
}

// BoolArrayValues is a free data retrieval call binding the contract method 0x80aacdcd.
//
// Solidity: function boolArrayValues(bytes32 , uint256 ) view returns(bool)
func (_DataStore *DataStoreCaller) BoolArrayValues(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "boolArrayValues", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BoolArrayValues is a free data retrieval call binding the contract method 0x80aacdcd.
//
// Solidity: function boolArrayValues(bytes32 , uint256 ) view returns(bool)
func (_DataStore *DataStoreSession) BoolArrayValues(arg0 [32]byte, arg1 *big.Int) (bool, error) {
	return _DataStore.Contract.BoolArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// BoolArrayValues is a free data retrieval call binding the contract method 0x80aacdcd.
//
// Solidity: function boolArrayValues(bytes32 , uint256 ) view returns(bool)
func (_DataStore *DataStoreCallerSession) BoolArrayValues(arg0 [32]byte, arg1 *big.Int) (bool, error) {
	return _DataStore.Contract.BoolArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// BoolValues is a free data retrieval call binding the contract method 0x44a242b1.
//
// Solidity: function boolValues(bytes32 ) view returns(bool)
func (_DataStore *DataStoreCaller) BoolValues(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "boolValues", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BoolValues is a free data retrieval call binding the contract method 0x44a242b1.
//
// Solidity: function boolValues(bytes32 ) view returns(bool)
func (_DataStore *DataStoreSession) BoolValues(arg0 [32]byte) (bool, error) {
	return _DataStore.Contract.BoolValues(&_DataStore.CallOpts, arg0)
}

// BoolValues is a free data retrieval call binding the contract method 0x44a242b1.
//
// Solidity: function boolValues(bytes32 ) view returns(bool)
func (_DataStore *DataStoreCallerSession) BoolValues(arg0 [32]byte) (bool, error) {
	return _DataStore.Contract.BoolValues(&_DataStore.CallOpts, arg0)
}

// Bytes32ArrayValues is a free data retrieval call binding the contract method 0xbf498dd3.
//
// Solidity: function bytes32ArrayValues(bytes32 , uint256 ) view returns(bytes32)
func (_DataStore *DataStoreCaller) Bytes32ArrayValues(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "bytes32ArrayValues", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Bytes32ArrayValues is a free data retrieval call binding the contract method 0xbf498dd3.
//
// Solidity: function bytes32ArrayValues(bytes32 , uint256 ) view returns(bytes32)
func (_DataStore *DataStoreSession) Bytes32ArrayValues(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _DataStore.Contract.Bytes32ArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// Bytes32ArrayValues is a free data retrieval call binding the contract method 0xbf498dd3.
//
// Solidity: function bytes32ArrayValues(bytes32 , uint256 ) view returns(bytes32)
func (_DataStore *DataStoreCallerSession) Bytes32ArrayValues(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _DataStore.Contract.Bytes32ArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// Bytes32Values is a free data retrieval call binding the contract method 0xd52852af.
//
// Solidity: function bytes32Values(bytes32 ) view returns(bytes32)
func (_DataStore *DataStoreCaller) Bytes32Values(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "bytes32Values", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Bytes32Values is a free data retrieval call binding the contract method 0xd52852af.
//
// Solidity: function bytes32Values(bytes32 ) view returns(bytes32)
func (_DataStore *DataStoreSession) Bytes32Values(arg0 [32]byte) ([32]byte, error) {
	return _DataStore.Contract.Bytes32Values(&_DataStore.CallOpts, arg0)
}

// Bytes32Values is a free data retrieval call binding the contract method 0xd52852af.
//
// Solidity: function bytes32Values(bytes32 ) view returns(bytes32)
func (_DataStore *DataStoreCallerSession) Bytes32Values(arg0 [32]byte) ([32]byte, error) {
	return _DataStore.Contract.Bytes32Values(&_DataStore.CallOpts, arg0)
}

// ContainsAddress is a free data retrieval call binding the contract method 0xc769d1a1.
//
// Solidity: function containsAddress(bytes32 setKey, address value) view returns(bool)
func (_DataStore *DataStoreCaller) ContainsAddress(opts *bind.CallOpts, setKey [32]byte, value common.Address) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "containsAddress", setKey, value)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ContainsAddress is a free data retrieval call binding the contract method 0xc769d1a1.
//
// Solidity: function containsAddress(bytes32 setKey, address value) view returns(bool)
func (_DataStore *DataStoreSession) ContainsAddress(setKey [32]byte, value common.Address) (bool, error) {
	return _DataStore.Contract.ContainsAddress(&_DataStore.CallOpts, setKey, value)
}

// ContainsAddress is a free data retrieval call binding the contract method 0xc769d1a1.
//
// Solidity: function containsAddress(bytes32 setKey, address value) view returns(bool)
func (_DataStore *DataStoreCallerSession) ContainsAddress(setKey [32]byte, value common.Address) (bool, error) {
	return _DataStore.Contract.ContainsAddress(&_DataStore.CallOpts, setKey, value)
}

// ContainsBytes32 is a free data retrieval call binding the contract method 0x91d4403c.
//
// Solidity: function containsBytes32(bytes32 setKey, bytes32 value) view returns(bool)
func (_DataStore *DataStoreCaller) ContainsBytes32(opts *bind.CallOpts, setKey [32]byte, value [32]byte) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "containsBytes32", setKey, value)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ContainsBytes32 is a free data retrieval call binding the contract method 0x91d4403c.
//
// Solidity: function containsBytes32(bytes32 setKey, bytes32 value) view returns(bool)
func (_DataStore *DataStoreSession) ContainsBytes32(setKey [32]byte, value [32]byte) (bool, error) {
	return _DataStore.Contract.ContainsBytes32(&_DataStore.CallOpts, setKey, value)
}

// ContainsBytes32 is a free data retrieval call binding the contract method 0x91d4403c.
//
// Solidity: function containsBytes32(bytes32 setKey, bytes32 value) view returns(bool)
func (_DataStore *DataStoreCallerSession) ContainsBytes32(setKey [32]byte, value [32]byte) (bool, error) {
	return _DataStore.Contract.ContainsBytes32(&_DataStore.CallOpts, setKey, value)
}

// ContainsUint is a free data retrieval call binding the contract method 0x310b8882.
//
// Solidity: function containsUint(bytes32 setKey, uint256 value) view returns(bool)
func (_DataStore *DataStoreCaller) ContainsUint(opts *bind.CallOpts, setKey [32]byte, value *big.Int) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "containsUint", setKey, value)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ContainsUint is a free data retrieval call binding the contract method 0x310b8882.
//
// Solidity: function containsUint(bytes32 setKey, uint256 value) view returns(bool)
func (_DataStore *DataStoreSession) ContainsUint(setKey [32]byte, value *big.Int) (bool, error) {
	return _DataStore.Contract.ContainsUint(&_DataStore.CallOpts, setKey, value)
}

// ContainsUint is a free data retrieval call binding the contract method 0x310b8882.
//
// Solidity: function containsUint(bytes32 setKey, uint256 value) view returns(bool)
func (_DataStore *DataStoreCallerSession) ContainsUint(setKey [32]byte, value *big.Int) (bool, error) {
	return _DataStore.Contract.ContainsUint(&_DataStore.CallOpts, setKey, value)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_DataStore *DataStoreCaller) GetAddress(opts *bind.CallOpts, key [32]byte) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getAddress", key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_DataStore *DataStoreSession) GetAddress(key [32]byte) (common.Address, error) {
	return _DataStore.Contract.GetAddress(&_DataStore.CallOpts, key)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_DataStore *DataStoreCallerSession) GetAddress(key [32]byte) (common.Address, error) {
	return _DataStore.Contract.GetAddress(&_DataStore.CallOpts, key)
}

// GetAddressArray is a free data retrieval call binding the contract method 0x5948f733.
//
// Solidity: function getAddressArray(bytes32 key) view returns(address[])
func (_DataStore *DataStoreCaller) GetAddressArray(opts *bind.CallOpts, key [32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getAddressArray", key)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAddressArray is a free data retrieval call binding the contract method 0x5948f733.
//
// Solidity: function getAddressArray(bytes32 key) view returns(address[])
func (_DataStore *DataStoreSession) GetAddressArray(key [32]byte) ([]common.Address, error) {
	return _DataStore.Contract.GetAddressArray(&_DataStore.CallOpts, key)
}

// GetAddressArray is a free data retrieval call binding the contract method 0x5948f733.
//
// Solidity: function getAddressArray(bytes32 key) view returns(address[])
func (_DataStore *DataStoreCallerSession) GetAddressArray(key [32]byte) ([]common.Address, error) {
	return _DataStore.Contract.GetAddressArray(&_DataStore.CallOpts, key)
}

// GetAddressCount is a free data retrieval call binding the contract method 0x35ea8059.
//
// Solidity: function getAddressCount(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreCaller) GetAddressCount(opts *bind.CallOpts, setKey [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getAddressCount", setKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAddressCount is a free data retrieval call binding the contract method 0x35ea8059.
//
// Solidity: function getAddressCount(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreSession) GetAddressCount(setKey [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetAddressCount(&_DataStore.CallOpts, setKey)
}

// GetAddressCount is a free data retrieval call binding the contract method 0x35ea8059.
//
// Solidity: function getAddressCount(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreCallerSession) GetAddressCount(setKey [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetAddressCount(&_DataStore.CallOpts, setKey)
}

// GetAddressValuesAt is a free data retrieval call binding the contract method 0xe7e4148e.
//
// Solidity: function getAddressValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(address[])
func (_DataStore *DataStoreCaller) GetAddressValuesAt(opts *bind.CallOpts, setKey [32]byte, start *big.Int, end *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getAddressValuesAt", setKey, start, end)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAddressValuesAt is a free data retrieval call binding the contract method 0xe7e4148e.
//
// Solidity: function getAddressValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(address[])
func (_DataStore *DataStoreSession) GetAddressValuesAt(setKey [32]byte, start *big.Int, end *big.Int) ([]common.Address, error) {
	return _DataStore.Contract.GetAddressValuesAt(&_DataStore.CallOpts, setKey, start, end)
}

// GetAddressValuesAt is a free data retrieval call binding the contract method 0xe7e4148e.
//
// Solidity: function getAddressValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(address[])
func (_DataStore *DataStoreCallerSession) GetAddressValuesAt(setKey [32]byte, start *big.Int, end *big.Int) ([]common.Address, error) {
	return _DataStore.Contract.GetAddressValuesAt(&_DataStore.CallOpts, setKey, start, end)
}

// GetBackend is a free data retrieval call binding the contract method 0x594eda64.
//
// Solidity: function getBackend() view returns(address)
func (_DataStore *DataStoreCaller) GetBackend(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBackend")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBackend is a free data retrieval call binding the contract method 0x594eda64.
//
// Solidity: function getBackend() view returns(address)
func (_DataStore *DataStoreSession) GetBackend() (common.Address, error) {
	return _DataStore.Contract.GetBackend(&_DataStore.CallOpts)
}

// GetBackend is a free data retrieval call binding the contract method 0x594eda64.
//
// Solidity: function getBackend() view returns(address)
func (_DataStore *DataStoreCallerSession) GetBackend() (common.Address, error) {
	return _DataStore.Contract.GetBackend(&_DataStore.CallOpts)
}

// GetBool is a free data retrieval call binding the contract method 0x7ae1cfca.
//
// Solidity: function getBool(bytes32 key) view returns(bool)
func (_DataStore *DataStoreCaller) GetBool(opts *bind.CallOpts, key [32]byte) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBool", key)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetBool is a free data retrieval call binding the contract method 0x7ae1cfca.
//
// Solidity: function getBool(bytes32 key) view returns(bool)
func (_DataStore *DataStoreSession) GetBool(key [32]byte) (bool, error) {
	return _DataStore.Contract.GetBool(&_DataStore.CallOpts, key)
}

// GetBool is a free data retrieval call binding the contract method 0x7ae1cfca.
//
// Solidity: function getBool(bytes32 key) view returns(bool)
func (_DataStore *DataStoreCallerSession) GetBool(key [32]byte) (bool, error) {
	return _DataStore.Contract.GetBool(&_DataStore.CallOpts, key)
}

// GetBoolArray is a free data retrieval call binding the contract method 0x116bb929.
//
// Solidity: function getBoolArray(bytes32 key) view returns(bool[])
func (_DataStore *DataStoreCaller) GetBoolArray(opts *bind.CallOpts, key [32]byte) ([]bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBoolArray", key)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetBoolArray is a free data retrieval call binding the contract method 0x116bb929.
//
// Solidity: function getBoolArray(bytes32 key) view returns(bool[])
func (_DataStore *DataStoreSession) GetBoolArray(key [32]byte) ([]bool, error) {
	return _DataStore.Contract.GetBoolArray(&_DataStore.CallOpts, key)
}

// GetBoolArray is a free data retrieval call binding the contract method 0x116bb929.
//
// Solidity: function getBoolArray(bytes32 key) view returns(bool[])
func (_DataStore *DataStoreCallerSession) GetBoolArray(key [32]byte) ([]bool, error) {
	return _DataStore.Contract.GetBoolArray(&_DataStore.CallOpts, key)
}

// GetBytes32 is a free data retrieval call binding the contract method 0xa6ed563e.
//
// Solidity: function getBytes32(bytes32 key) view returns(bytes32)
func (_DataStore *DataStoreCaller) GetBytes32(opts *bind.CallOpts, key [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBytes32", key)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBytes32 is a free data retrieval call binding the contract method 0xa6ed563e.
//
// Solidity: function getBytes32(bytes32 key) view returns(bytes32)
func (_DataStore *DataStoreSession) GetBytes32(key [32]byte) ([32]byte, error) {
	return _DataStore.Contract.GetBytes32(&_DataStore.CallOpts, key)
}

// GetBytes32 is a free data retrieval call binding the contract method 0xa6ed563e.
//
// Solidity: function getBytes32(bytes32 key) view returns(bytes32)
func (_DataStore *DataStoreCallerSession) GetBytes32(key [32]byte) ([32]byte, error) {
	return _DataStore.Contract.GetBytes32(&_DataStore.CallOpts, key)
}

// GetBytes32Array is a free data retrieval call binding the contract method 0xdd031997.
//
// Solidity: function getBytes32Array(bytes32 key) view returns(bytes32[])
func (_DataStore *DataStoreCaller) GetBytes32Array(opts *bind.CallOpts, key [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBytes32Array", key)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetBytes32Array is a free data retrieval call binding the contract method 0xdd031997.
//
// Solidity: function getBytes32Array(bytes32 key) view returns(bytes32[])
func (_DataStore *DataStoreSession) GetBytes32Array(key [32]byte) ([][32]byte, error) {
	return _DataStore.Contract.GetBytes32Array(&_DataStore.CallOpts, key)
}

// GetBytes32Array is a free data retrieval call binding the contract method 0xdd031997.
//
// Solidity: function getBytes32Array(bytes32 key) view returns(bytes32[])
func (_DataStore *DataStoreCallerSession) GetBytes32Array(key [32]byte) ([][32]byte, error) {
	return _DataStore.Contract.GetBytes32Array(&_DataStore.CallOpts, key)
}

// GetBytes32Count is a free data retrieval call binding the contract method 0xf3903b9f.
//
// Solidity: function getBytes32Count(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreCaller) GetBytes32Count(opts *bind.CallOpts, setKey [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBytes32Count", setKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBytes32Count is a free data retrieval call binding the contract method 0xf3903b9f.
//
// Solidity: function getBytes32Count(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreSession) GetBytes32Count(setKey [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetBytes32Count(&_DataStore.CallOpts, setKey)
}

// GetBytes32Count is a free data retrieval call binding the contract method 0xf3903b9f.
//
// Solidity: function getBytes32Count(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreCallerSession) GetBytes32Count(setKey [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetBytes32Count(&_DataStore.CallOpts, setKey)
}

// GetBytes32ValuesAt is a free data retrieval call binding the contract method 0xf069052a.
//
// Solidity: function getBytes32ValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(bytes32[])
func (_DataStore *DataStoreCaller) GetBytes32ValuesAt(opts *bind.CallOpts, setKey [32]byte, start *big.Int, end *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getBytes32ValuesAt", setKey, start, end)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetBytes32ValuesAt is a free data retrieval call binding the contract method 0xf069052a.
//
// Solidity: function getBytes32ValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(bytes32[])
func (_DataStore *DataStoreSession) GetBytes32ValuesAt(setKey [32]byte, start *big.Int, end *big.Int) ([][32]byte, error) {
	return _DataStore.Contract.GetBytes32ValuesAt(&_DataStore.CallOpts, setKey, start, end)
}

// GetBytes32ValuesAt is a free data retrieval call binding the contract method 0xf069052a.
//
// Solidity: function getBytes32ValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(bytes32[])
func (_DataStore *DataStoreCallerSession) GetBytes32ValuesAt(setKey [32]byte, start *big.Int, end *big.Int) ([][32]byte, error) {
	return _DataStore.Contract.GetBytes32ValuesAt(&_DataStore.CallOpts, setKey, start, end)
}

// GetInt is a free data retrieval call binding the contract method 0xdc97d962.
//
// Solidity: function getInt(bytes32 key) view returns(int256)
func (_DataStore *DataStoreCaller) GetInt(opts *bind.CallOpts, key [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getInt", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetInt is a free data retrieval call binding the contract method 0xdc97d962.
//
// Solidity: function getInt(bytes32 key) view returns(int256)
func (_DataStore *DataStoreSession) GetInt(key [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetInt(&_DataStore.CallOpts, key)
}

// GetInt is a free data retrieval call binding the contract method 0xdc97d962.
//
// Solidity: function getInt(bytes32 key) view returns(int256)
func (_DataStore *DataStoreCallerSession) GetInt(key [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetInt(&_DataStore.CallOpts, key)
}

// GetIntArray is a free data retrieval call binding the contract method 0x2d2899b6.
//
// Solidity: function getIntArray(bytes32 key) view returns(int256[])
func (_DataStore *DataStoreCaller) GetIntArray(opts *bind.CallOpts, key [32]byte) ([]*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getIntArray", key)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetIntArray is a free data retrieval call binding the contract method 0x2d2899b6.
//
// Solidity: function getIntArray(bytes32 key) view returns(int256[])
func (_DataStore *DataStoreSession) GetIntArray(key [32]byte) ([]*big.Int, error) {
	return _DataStore.Contract.GetIntArray(&_DataStore.CallOpts, key)
}

// GetIntArray is a free data retrieval call binding the contract method 0x2d2899b6.
//
// Solidity: function getIntArray(bytes32 key) view returns(int256[])
func (_DataStore *DataStoreCallerSession) GetIntArray(key [32]byte) ([]*big.Int, error) {
	return _DataStore.Contract.GetIntArray(&_DataStore.CallOpts, key)
}

// GetPositionHandle is a free data retrieval call binding the contract method 0xb1eecbe8.
//
// Solidity: function getPositionHandle() view returns(address)
func (_DataStore *DataStoreCaller) GetPositionHandle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getPositionHandle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPositionHandle is a free data retrieval call binding the contract method 0xb1eecbe8.
//
// Solidity: function getPositionHandle() view returns(address)
func (_DataStore *DataStoreSession) GetPositionHandle() (common.Address, error) {
	return _DataStore.Contract.GetPositionHandle(&_DataStore.CallOpts)
}

// GetPositionHandle is a free data retrieval call binding the contract method 0xb1eecbe8.
//
// Solidity: function getPositionHandle() view returns(address)
func (_DataStore *DataStoreCallerSession) GetPositionHandle() (common.Address, error) {
	return _DataStore.Contract.GetPositionHandle(&_DataStore.CallOpts)
}

// GetString is a free data retrieval call binding the contract method 0x986e791a.
//
// Solidity: function getString(bytes32 key) view returns(string)
func (_DataStore *DataStoreCaller) GetString(opts *bind.CallOpts, key [32]byte) (string, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getString", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetString is a free data retrieval call binding the contract method 0x986e791a.
//
// Solidity: function getString(bytes32 key) view returns(string)
func (_DataStore *DataStoreSession) GetString(key [32]byte) (string, error) {
	return _DataStore.Contract.GetString(&_DataStore.CallOpts, key)
}

// GetString is a free data retrieval call binding the contract method 0x986e791a.
//
// Solidity: function getString(bytes32 key) view returns(string)
func (_DataStore *DataStoreCallerSession) GetString(key [32]byte) (string, error) {
	return _DataStore.Contract.GetString(&_DataStore.CallOpts, key)
}

// GetStringArray is a free data retrieval call binding the contract method 0x01677da2.
//
// Solidity: function getStringArray(bytes32 key) view returns(string[])
func (_DataStore *DataStoreCaller) GetStringArray(opts *bind.CallOpts, key [32]byte) ([]string, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getStringArray", key)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetStringArray is a free data retrieval call binding the contract method 0x01677da2.
//
// Solidity: function getStringArray(bytes32 key) view returns(string[])
func (_DataStore *DataStoreSession) GetStringArray(key [32]byte) ([]string, error) {
	return _DataStore.Contract.GetStringArray(&_DataStore.CallOpts, key)
}

// GetStringArray is a free data retrieval call binding the contract method 0x01677da2.
//
// Solidity: function getStringArray(bytes32 key) view returns(string[])
func (_DataStore *DataStoreCallerSession) GetStringArray(key [32]byte) ([]string, error) {
	return _DataStore.Contract.GetStringArray(&_DataStore.CallOpts, key)
}

// GetUint is a free data retrieval call binding the contract method 0xbd02d0f5.
//
// Solidity: function getUint(bytes32 key) view returns(uint256)
func (_DataStore *DataStoreCaller) GetUint(opts *bind.CallOpts, key [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getUint", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUint is a free data retrieval call binding the contract method 0xbd02d0f5.
//
// Solidity: function getUint(bytes32 key) view returns(uint256)
func (_DataStore *DataStoreSession) GetUint(key [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetUint(&_DataStore.CallOpts, key)
}

// GetUint is a free data retrieval call binding the contract method 0xbd02d0f5.
//
// Solidity: function getUint(bytes32 key) view returns(uint256)
func (_DataStore *DataStoreCallerSession) GetUint(key [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetUint(&_DataStore.CallOpts, key)
}

// GetUintArray is a free data retrieval call binding the contract method 0x86ac6bdf.
//
// Solidity: function getUintArray(bytes32 key) view returns(uint256[])
func (_DataStore *DataStoreCaller) GetUintArray(opts *bind.CallOpts, key [32]byte) ([]*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getUintArray", key)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUintArray is a free data retrieval call binding the contract method 0x86ac6bdf.
//
// Solidity: function getUintArray(bytes32 key) view returns(uint256[])
func (_DataStore *DataStoreSession) GetUintArray(key [32]byte) ([]*big.Int, error) {
	return _DataStore.Contract.GetUintArray(&_DataStore.CallOpts, key)
}

// GetUintArray is a free data retrieval call binding the contract method 0x86ac6bdf.
//
// Solidity: function getUintArray(bytes32 key) view returns(uint256[])
func (_DataStore *DataStoreCallerSession) GetUintArray(key [32]byte) ([]*big.Int, error) {
	return _DataStore.Contract.GetUintArray(&_DataStore.CallOpts, key)
}

// GetUintCount is a free data retrieval call binding the contract method 0x065f21a7.
//
// Solidity: function getUintCount(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreCaller) GetUintCount(opts *bind.CallOpts, setKey [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getUintCount", setKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUintCount is a free data retrieval call binding the contract method 0x065f21a7.
//
// Solidity: function getUintCount(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreSession) GetUintCount(setKey [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetUintCount(&_DataStore.CallOpts, setKey)
}

// GetUintCount is a free data retrieval call binding the contract method 0x065f21a7.
//
// Solidity: function getUintCount(bytes32 setKey) view returns(uint256)
func (_DataStore *DataStoreCallerSession) GetUintCount(setKey [32]byte) (*big.Int, error) {
	return _DataStore.Contract.GetUintCount(&_DataStore.CallOpts, setKey)
}

// GetUintValuesAt is a free data retrieval call binding the contract method 0x7026d42c.
//
// Solidity: function getUintValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(uint256[])
func (_DataStore *DataStoreCaller) GetUintValuesAt(opts *bind.CallOpts, setKey [32]byte, start *big.Int, end *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "getUintValuesAt", setKey, start, end)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUintValuesAt is a free data retrieval call binding the contract method 0x7026d42c.
//
// Solidity: function getUintValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(uint256[])
func (_DataStore *DataStoreSession) GetUintValuesAt(setKey [32]byte, start *big.Int, end *big.Int) ([]*big.Int, error) {
	return _DataStore.Contract.GetUintValuesAt(&_DataStore.CallOpts, setKey, start, end)
}

// GetUintValuesAt is a free data retrieval call binding the contract method 0x7026d42c.
//
// Solidity: function getUintValuesAt(bytes32 setKey, uint256 start, uint256 end) view returns(uint256[])
func (_DataStore *DataStoreCallerSession) GetUintValuesAt(setKey [32]byte, start *big.Int, end *big.Int) ([]*big.Int, error) {
	return _DataStore.Contract.GetUintValuesAt(&_DataStore.CallOpts, setKey, start, end)
}

// IntArrayValues is a free data retrieval call binding the contract method 0x6339734d.
//
// Solidity: function intArrayValues(bytes32 , uint256 ) view returns(int256)
func (_DataStore *DataStoreCaller) IntArrayValues(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "intArrayValues", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IntArrayValues is a free data retrieval call binding the contract method 0x6339734d.
//
// Solidity: function intArrayValues(bytes32 , uint256 ) view returns(int256)
func (_DataStore *DataStoreSession) IntArrayValues(arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	return _DataStore.Contract.IntArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// IntArrayValues is a free data retrieval call binding the contract method 0x6339734d.
//
// Solidity: function intArrayValues(bytes32 , uint256 ) view returns(int256)
func (_DataStore *DataStoreCallerSession) IntArrayValues(arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	return _DataStore.Contract.IntArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// IntValues is a free data retrieval call binding the contract method 0x743df325.
//
// Solidity: function intValues(bytes32 ) view returns(int256)
func (_DataStore *DataStoreCaller) IntValues(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "intValues", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IntValues is a free data retrieval call binding the contract method 0x743df325.
//
// Solidity: function intValues(bytes32 ) view returns(int256)
func (_DataStore *DataStoreSession) IntValues(arg0 [32]byte) (*big.Int, error) {
	return _DataStore.Contract.IntValues(&_DataStore.CallOpts, arg0)
}

// IntValues is a free data retrieval call binding the contract method 0x743df325.
//
// Solidity: function intValues(bytes32 ) view returns(int256)
func (_DataStore *DataStoreCallerSession) IntValues(arg0 [32]byte) (*big.Int, error) {
	return _DataStore.Contract.IntValues(&_DataStore.CallOpts, arg0)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_DataStore *DataStoreCaller) RoleStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "roleStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_DataStore *DataStoreSession) RoleStore() (common.Address, error) {
	return _DataStore.Contract.RoleStore(&_DataStore.CallOpts)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_DataStore *DataStoreCallerSession) RoleStore() (common.Address, error) {
	return _DataStore.Contract.RoleStore(&_DataStore.CallOpts)
}

// StringArrayValues is a free data retrieval call binding the contract method 0xb8320a08.
//
// Solidity: function stringArrayValues(bytes32 , uint256 ) view returns(string)
func (_DataStore *DataStoreCaller) StringArrayValues(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (string, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "stringArrayValues", arg0, arg1)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// StringArrayValues is a free data retrieval call binding the contract method 0xb8320a08.
//
// Solidity: function stringArrayValues(bytes32 , uint256 ) view returns(string)
func (_DataStore *DataStoreSession) StringArrayValues(arg0 [32]byte, arg1 *big.Int) (string, error) {
	return _DataStore.Contract.StringArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// StringArrayValues is a free data retrieval call binding the contract method 0xb8320a08.
//
// Solidity: function stringArrayValues(bytes32 , uint256 ) view returns(string)
func (_DataStore *DataStoreCallerSession) StringArrayValues(arg0 [32]byte, arg1 *big.Int) (string, error) {
	return _DataStore.Contract.StringArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// StringValues is a free data retrieval call binding the contract method 0xf15caeac.
//
// Solidity: function stringValues(bytes32 ) view returns(string)
func (_DataStore *DataStoreCaller) StringValues(opts *bind.CallOpts, arg0 [32]byte) (string, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "stringValues", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// StringValues is a free data retrieval call binding the contract method 0xf15caeac.
//
// Solidity: function stringValues(bytes32 ) view returns(string)
func (_DataStore *DataStoreSession) StringValues(arg0 [32]byte) (string, error) {
	return _DataStore.Contract.StringValues(&_DataStore.CallOpts, arg0)
}

// StringValues is a free data retrieval call binding the contract method 0xf15caeac.
//
// Solidity: function stringValues(bytes32 ) view returns(string)
func (_DataStore *DataStoreCallerSession) StringValues(arg0 [32]byte) (string, error) {
	return _DataStore.Contract.StringValues(&_DataStore.CallOpts, arg0)
}

// UintArrayValues is a free data retrieval call binding the contract method 0xc4f00fde.
//
// Solidity: function uintArrayValues(bytes32 , uint256 ) view returns(uint256)
func (_DataStore *DataStoreCaller) UintArrayValues(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "uintArrayValues", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UintArrayValues is a free data retrieval call binding the contract method 0xc4f00fde.
//
// Solidity: function uintArrayValues(bytes32 , uint256 ) view returns(uint256)
func (_DataStore *DataStoreSession) UintArrayValues(arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	return _DataStore.Contract.UintArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// UintArrayValues is a free data retrieval call binding the contract method 0xc4f00fde.
//
// Solidity: function uintArrayValues(bytes32 , uint256 ) view returns(uint256)
func (_DataStore *DataStoreCallerSession) UintArrayValues(arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	return _DataStore.Contract.UintArrayValues(&_DataStore.CallOpts, arg0, arg1)
}

// UintValues is a free data retrieval call binding the contract method 0xd38eebc7.
//
// Solidity: function uintValues(bytes32 ) view returns(uint256)
func (_DataStore *DataStoreCaller) UintValues(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "uintValues", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UintValues is a free data retrieval call binding the contract method 0xd38eebc7.
//
// Solidity: function uintValues(bytes32 ) view returns(uint256)
func (_DataStore *DataStoreSession) UintValues(arg0 [32]byte) (*big.Int, error) {
	return _DataStore.Contract.UintValues(&_DataStore.CallOpts, arg0)
}

// UintValues is a free data retrieval call binding the contract method 0xd38eebc7.
//
// Solidity: function uintValues(bytes32 ) view returns(uint256)
func (_DataStore *DataStoreCallerSession) UintValues(arg0 [32]byte) (*big.Int, error) {
	return _DataStore.Contract.UintValues(&_DataStore.CallOpts, arg0)
}

// AddAddress is a paid mutator transaction binding the contract method 0xb348e639.
//
// Solidity: function addAddress(bytes32 setKey, address value) returns()
func (_DataStore *DataStoreTransactor) AddAddress(opts *bind.TransactOpts, setKey [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "addAddress", setKey, value)
}

// AddAddress is a paid mutator transaction binding the contract method 0xb348e639.
//
// Solidity: function addAddress(bytes32 setKey, address value) returns()
func (_DataStore *DataStoreSession) AddAddress(setKey [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.AddAddress(&_DataStore.TransactOpts, setKey, value)
}

// AddAddress is a paid mutator transaction binding the contract method 0xb348e639.
//
// Solidity: function addAddress(bytes32 setKey, address value) returns()
func (_DataStore *DataStoreTransactorSession) AddAddress(setKey [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.AddAddress(&_DataStore.TransactOpts, setKey, value)
}

// AddBytes32 is a paid mutator transaction binding the contract method 0xc80f4c62.
//
// Solidity: function addBytes32(bytes32 setKey, bytes32 value) returns()
func (_DataStore *DataStoreTransactor) AddBytes32(opts *bind.TransactOpts, setKey [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "addBytes32", setKey, value)
}

// AddBytes32 is a paid mutator transaction binding the contract method 0xc80f4c62.
//
// Solidity: function addBytes32(bytes32 setKey, bytes32 value) returns()
func (_DataStore *DataStoreSession) AddBytes32(setKey [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.AddBytes32(&_DataStore.TransactOpts, setKey, value)
}

// AddBytes32 is a paid mutator transaction binding the contract method 0xc80f4c62.
//
// Solidity: function addBytes32(bytes32 setKey, bytes32 value) returns()
func (_DataStore *DataStoreTransactorSession) AddBytes32(setKey [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.AddBytes32(&_DataStore.TransactOpts, setKey, value)
}

// AddUint is a paid mutator transaction binding the contract method 0xadb353dc.
//
// Solidity: function addUint(bytes32 setKey, uint256 value) returns()
func (_DataStore *DataStoreTransactor) AddUint(opts *bind.TransactOpts, setKey [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "addUint", setKey, value)
}

// AddUint is a paid mutator transaction binding the contract method 0xadb353dc.
//
// Solidity: function addUint(bytes32 setKey, uint256 value) returns()
func (_DataStore *DataStoreSession) AddUint(setKey [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.AddUint(&_DataStore.TransactOpts, setKey, value)
}

// AddUint is a paid mutator transaction binding the contract method 0xadb353dc.
//
// Solidity: function addUint(bytes32 setKey, uint256 value) returns()
func (_DataStore *DataStoreTransactorSession) AddUint(setKey [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.AddUint(&_DataStore.TransactOpts, setKey, value)
}

// ApplyBoundedDeltaToUint is a paid mutator transaction binding the contract method 0x8ca498b0.
//
// Solidity: function applyBoundedDeltaToUint(bytes32 key, int256 value) returns(uint256)
func (_DataStore *DataStoreTransactor) ApplyBoundedDeltaToUint(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "applyBoundedDeltaToUint", key, value)
}

// ApplyBoundedDeltaToUint is a paid mutator transaction binding the contract method 0x8ca498b0.
//
// Solidity: function applyBoundedDeltaToUint(bytes32 key, int256 value) returns(uint256)
func (_DataStore *DataStoreSession) ApplyBoundedDeltaToUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyBoundedDeltaToUint(&_DataStore.TransactOpts, key, value)
}

// ApplyBoundedDeltaToUint is a paid mutator transaction binding the contract method 0x8ca498b0.
//
// Solidity: function applyBoundedDeltaToUint(bytes32 key, int256 value) returns(uint256)
func (_DataStore *DataStoreTransactorSession) ApplyBoundedDeltaToUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyBoundedDeltaToUint(&_DataStore.TransactOpts, key, value)
}

// ApplyDeltaToInt is a paid mutator transaction binding the contract method 0xe4e36c4e.
//
// Solidity: function applyDeltaToInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactor) ApplyDeltaToInt(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "applyDeltaToInt", key, value)
}

// ApplyDeltaToInt is a paid mutator transaction binding the contract method 0xe4e36c4e.
//
// Solidity: function applyDeltaToInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreSession) ApplyDeltaToInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyDeltaToInt(&_DataStore.TransactOpts, key, value)
}

// ApplyDeltaToInt is a paid mutator transaction binding the contract method 0xe4e36c4e.
//
// Solidity: function applyDeltaToInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactorSession) ApplyDeltaToInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyDeltaToInt(&_DataStore.TransactOpts, key, value)
}

// ApplyDeltaToUint is a paid mutator transaction binding the contract method 0x32f85bbf.
//
// Solidity: function applyDeltaToUint(bytes32 key, int256 value, string errorMessage) returns(uint256)
func (_DataStore *DataStoreTransactor) ApplyDeltaToUint(opts *bind.TransactOpts, key [32]byte, value *big.Int, errorMessage string) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "applyDeltaToUint", key, value, errorMessage)
}

// ApplyDeltaToUint is a paid mutator transaction binding the contract method 0x32f85bbf.
//
// Solidity: function applyDeltaToUint(bytes32 key, int256 value, string errorMessage) returns(uint256)
func (_DataStore *DataStoreSession) ApplyDeltaToUint(key [32]byte, value *big.Int, errorMessage string) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyDeltaToUint(&_DataStore.TransactOpts, key, value, errorMessage)
}

// ApplyDeltaToUint is a paid mutator transaction binding the contract method 0x32f85bbf.
//
// Solidity: function applyDeltaToUint(bytes32 key, int256 value, string errorMessage) returns(uint256)
func (_DataStore *DataStoreTransactorSession) ApplyDeltaToUint(key [32]byte, value *big.Int, errorMessage string) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyDeltaToUint(&_DataStore.TransactOpts, key, value, errorMessage)
}

// ApplyDeltaToUint0 is a paid mutator transaction binding the contract method 0x3dbacd1a.
//
// Solidity: function applyDeltaToUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactor) ApplyDeltaToUint0(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "applyDeltaToUint0", key, value)
}

// ApplyDeltaToUint0 is a paid mutator transaction binding the contract method 0x3dbacd1a.
//
// Solidity: function applyDeltaToUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreSession) ApplyDeltaToUint0(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyDeltaToUint0(&_DataStore.TransactOpts, key, value)
}

// ApplyDeltaToUint0 is a paid mutator transaction binding the contract method 0x3dbacd1a.
//
// Solidity: function applyDeltaToUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactorSession) ApplyDeltaToUint0(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.ApplyDeltaToUint0(&_DataStore.TransactOpts, key, value)
}

// DecrementInt is a paid mutator transaction binding the contract method 0x6fae54f0.
//
// Solidity: function decrementInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactor) DecrementInt(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "decrementInt", key, value)
}

// DecrementInt is a paid mutator transaction binding the contract method 0x6fae54f0.
//
// Solidity: function decrementInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreSession) DecrementInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.DecrementInt(&_DataStore.TransactOpts, key, value)
}

// DecrementInt is a paid mutator transaction binding the contract method 0x6fae54f0.
//
// Solidity: function decrementInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactorSession) DecrementInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.DecrementInt(&_DataStore.TransactOpts, key, value)
}

// DecrementUint is a paid mutator transaction binding the contract method 0xe98aabc1.
//
// Solidity: function decrementUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactor) DecrementUint(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "decrementUint", key, value)
}

// DecrementUint is a paid mutator transaction binding the contract method 0xe98aabc1.
//
// Solidity: function decrementUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreSession) DecrementUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.DecrementUint(&_DataStore.TransactOpts, key, value)
}

// DecrementUint is a paid mutator transaction binding the contract method 0xe98aabc1.
//
// Solidity: function decrementUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactorSession) DecrementUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.DecrementUint(&_DataStore.TransactOpts, key, value)
}

// IncrementInt is a paid mutator transaction binding the contract method 0xcbb093dd.
//
// Solidity: function incrementInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactor) IncrementInt(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "incrementInt", key, value)
}

// IncrementInt is a paid mutator transaction binding the contract method 0xcbb093dd.
//
// Solidity: function incrementInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreSession) IncrementInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.IncrementInt(&_DataStore.TransactOpts, key, value)
}

// IncrementInt is a paid mutator transaction binding the contract method 0xcbb093dd.
//
// Solidity: function incrementInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactorSession) IncrementInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.IncrementInt(&_DataStore.TransactOpts, key, value)
}

// IncrementUint is a paid mutator transaction binding the contract method 0x340dbab3.
//
// Solidity: function incrementUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactor) IncrementUint(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "incrementUint", key, value)
}

// IncrementUint is a paid mutator transaction binding the contract method 0x340dbab3.
//
// Solidity: function incrementUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreSession) IncrementUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.IncrementUint(&_DataStore.TransactOpts, key, value)
}

// IncrementUint is a paid mutator transaction binding the contract method 0x340dbab3.
//
// Solidity: function incrementUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactorSession) IncrementUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.IncrementUint(&_DataStore.TransactOpts, key, value)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _roleStore) returns()
func (_DataStore *DataStoreTransactor) Initialize(opts *bind.TransactOpts, _roleStore common.Address) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "initialize", _roleStore)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _roleStore) returns()
func (_DataStore *DataStoreSession) Initialize(_roleStore common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.Initialize(&_DataStore.TransactOpts, _roleStore)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _roleStore) returns()
func (_DataStore *DataStoreTransactorSession) Initialize(_roleStore common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.Initialize(&_DataStore.TransactOpts, _roleStore)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0x69721d41.
//
// Solidity: function removeAddress(bytes32 setKey, address value) returns()
func (_DataStore *DataStoreTransactor) RemoveAddress(opts *bind.TransactOpts, setKey [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeAddress", setKey, value)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0x69721d41.
//
// Solidity: function removeAddress(bytes32 setKey, address value) returns()
func (_DataStore *DataStoreSession) RemoveAddress(setKey [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveAddress(&_DataStore.TransactOpts, setKey, value)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0x69721d41.
//
// Solidity: function removeAddress(bytes32 setKey, address value) returns()
func (_DataStore *DataStoreTransactorSession) RemoveAddress(setKey [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveAddress(&_DataStore.TransactOpts, setKey, value)
}

// RemoveAddress0 is a paid mutator transaction binding the contract method 0x9faf6fb6.
//
// Solidity: function removeAddress(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveAddress0(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeAddress0", key)
}

// RemoveAddress0 is a paid mutator transaction binding the contract method 0x9faf6fb6.
//
// Solidity: function removeAddress(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveAddress0(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveAddress0(&_DataStore.TransactOpts, key)
}

// RemoveAddress0 is a paid mutator transaction binding the contract method 0x9faf6fb6.
//
// Solidity: function removeAddress(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveAddress0(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveAddress0(&_DataStore.TransactOpts, key)
}

// RemoveAddressArray is a paid mutator transaction binding the contract method 0xc1dc9182.
//
// Solidity: function removeAddressArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveAddressArray(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeAddressArray", key)
}

// RemoveAddressArray is a paid mutator transaction binding the contract method 0xc1dc9182.
//
// Solidity: function removeAddressArray(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveAddressArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveAddressArray(&_DataStore.TransactOpts, key)
}

// RemoveAddressArray is a paid mutator transaction binding the contract method 0xc1dc9182.
//
// Solidity: function removeAddressArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveAddressArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveAddressArray(&_DataStore.TransactOpts, key)
}

// RemoveBool is a paid mutator transaction binding the contract method 0x9fe7ac12.
//
// Solidity: function removeBool(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveBool(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeBool", key)
}

// RemoveBool is a paid mutator transaction binding the contract method 0x9fe7ac12.
//
// Solidity: function removeBool(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveBool(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBool(&_DataStore.TransactOpts, key)
}

// RemoveBool is a paid mutator transaction binding the contract method 0x9fe7ac12.
//
// Solidity: function removeBool(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveBool(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBool(&_DataStore.TransactOpts, key)
}

// RemoveBoolArray is a paid mutator transaction binding the contract method 0xf51fc0d9.
//
// Solidity: function removeBoolArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveBoolArray(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeBoolArray", key)
}

// RemoveBoolArray is a paid mutator transaction binding the contract method 0xf51fc0d9.
//
// Solidity: function removeBoolArray(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveBoolArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBoolArray(&_DataStore.TransactOpts, key)
}

// RemoveBoolArray is a paid mutator transaction binding the contract method 0xf51fc0d9.
//
// Solidity: function removeBoolArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveBoolArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBoolArray(&_DataStore.TransactOpts, key)
}

// RemoveBytes32 is a paid mutator transaction binding the contract method 0x9921c3cc.
//
// Solidity: function removeBytes32(bytes32 setKey, bytes32 value) returns()
func (_DataStore *DataStoreTransactor) RemoveBytes32(opts *bind.TransactOpts, setKey [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeBytes32", setKey, value)
}

// RemoveBytes32 is a paid mutator transaction binding the contract method 0x9921c3cc.
//
// Solidity: function removeBytes32(bytes32 setKey, bytes32 value) returns()
func (_DataStore *DataStoreSession) RemoveBytes32(setKey [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBytes32(&_DataStore.TransactOpts, setKey, value)
}

// RemoveBytes32 is a paid mutator transaction binding the contract method 0x9921c3cc.
//
// Solidity: function removeBytes32(bytes32 setKey, bytes32 value) returns()
func (_DataStore *DataStoreTransactorSession) RemoveBytes32(setKey [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBytes32(&_DataStore.TransactOpts, setKey, value)
}

// RemoveBytes320 is a paid mutator transaction binding the contract method 0xcf6a8722.
//
// Solidity: function removeBytes32(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveBytes320(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeBytes320", key)
}

// RemoveBytes320 is a paid mutator transaction binding the contract method 0xcf6a8722.
//
// Solidity: function removeBytes32(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveBytes320(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBytes320(&_DataStore.TransactOpts, key)
}

// RemoveBytes320 is a paid mutator transaction binding the contract method 0xcf6a8722.
//
// Solidity: function removeBytes32(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveBytes320(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBytes320(&_DataStore.TransactOpts, key)
}

// RemoveBytes32Array is a paid mutator transaction binding the contract method 0xbf7f035a.
//
// Solidity: function removeBytes32Array(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveBytes32Array(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeBytes32Array", key)
}

// RemoveBytes32Array is a paid mutator transaction binding the contract method 0xbf7f035a.
//
// Solidity: function removeBytes32Array(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveBytes32Array(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBytes32Array(&_DataStore.TransactOpts, key)
}

// RemoveBytes32Array is a paid mutator transaction binding the contract method 0xbf7f035a.
//
// Solidity: function removeBytes32Array(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveBytes32Array(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveBytes32Array(&_DataStore.TransactOpts, key)
}

// RemoveInt is a paid mutator transaction binding the contract method 0xe62461ce.
//
// Solidity: function removeInt(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveInt(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeInt", key)
}

// RemoveInt is a paid mutator transaction binding the contract method 0xe62461ce.
//
// Solidity: function removeInt(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveInt(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveInt(&_DataStore.TransactOpts, key)
}

// RemoveInt is a paid mutator transaction binding the contract method 0xe62461ce.
//
// Solidity: function removeInt(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveInt(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveInt(&_DataStore.TransactOpts, key)
}

// RemoveIntArray is a paid mutator transaction binding the contract method 0x499ea50e.
//
// Solidity: function removeIntArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveIntArray(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeIntArray", key)
}

// RemoveIntArray is a paid mutator transaction binding the contract method 0x499ea50e.
//
// Solidity: function removeIntArray(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveIntArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveIntArray(&_DataStore.TransactOpts, key)
}

// RemoveIntArray is a paid mutator transaction binding the contract method 0x499ea50e.
//
// Solidity: function removeIntArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveIntArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveIntArray(&_DataStore.TransactOpts, key)
}

// RemoveString is a paid mutator transaction binding the contract method 0xcc50eadd.
//
// Solidity: function removeString(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveString(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeString", key)
}

// RemoveString is a paid mutator transaction binding the contract method 0xcc50eadd.
//
// Solidity: function removeString(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveString(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveString(&_DataStore.TransactOpts, key)
}

// RemoveString is a paid mutator transaction binding the contract method 0xcc50eadd.
//
// Solidity: function removeString(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveString(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveString(&_DataStore.TransactOpts, key)
}

// RemoveStringArray is a paid mutator transaction binding the contract method 0xe208a70d.
//
// Solidity: function removeStringArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveStringArray(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeStringArray", key)
}

// RemoveStringArray is a paid mutator transaction binding the contract method 0xe208a70d.
//
// Solidity: function removeStringArray(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveStringArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveStringArray(&_DataStore.TransactOpts, key)
}

// RemoveStringArray is a paid mutator transaction binding the contract method 0xe208a70d.
//
// Solidity: function removeStringArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveStringArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveStringArray(&_DataStore.TransactOpts, key)
}

// RemoveUint is a paid mutator transaction binding the contract method 0x42c3bd96.
//
// Solidity: function removeUint(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveUint(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeUint", key)
}

// RemoveUint is a paid mutator transaction binding the contract method 0x42c3bd96.
//
// Solidity: function removeUint(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveUint(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveUint(&_DataStore.TransactOpts, key)
}

// RemoveUint is a paid mutator transaction binding the contract method 0x42c3bd96.
//
// Solidity: function removeUint(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveUint(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveUint(&_DataStore.TransactOpts, key)
}

// RemoveUint0 is a paid mutator transaction binding the contract method 0x93266f9a.
//
// Solidity: function removeUint(bytes32 setKey, uint256 value) returns()
func (_DataStore *DataStoreTransactor) RemoveUint0(opts *bind.TransactOpts, setKey [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeUint0", setKey, value)
}

// RemoveUint0 is a paid mutator transaction binding the contract method 0x93266f9a.
//
// Solidity: function removeUint(bytes32 setKey, uint256 value) returns()
func (_DataStore *DataStoreSession) RemoveUint0(setKey [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveUint0(&_DataStore.TransactOpts, setKey, value)
}

// RemoveUint0 is a paid mutator transaction binding the contract method 0x93266f9a.
//
// Solidity: function removeUint(bytes32 setKey, uint256 value) returns()
func (_DataStore *DataStoreTransactorSession) RemoveUint0(setKey [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveUint0(&_DataStore.TransactOpts, setKey, value)
}

// RemoveUintArray is a paid mutator transaction binding the contract method 0xbe43caa3.
//
// Solidity: function removeUintArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactor) RemoveUintArray(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "removeUintArray", key)
}

// RemoveUintArray is a paid mutator transaction binding the contract method 0xbe43caa3.
//
// Solidity: function removeUintArray(bytes32 key) returns()
func (_DataStore *DataStoreSession) RemoveUintArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveUintArray(&_DataStore.TransactOpts, key)
}

// RemoveUintArray is a paid mutator transaction binding the contract method 0xbe43caa3.
//
// Solidity: function removeUintArray(bytes32 key) returns()
func (_DataStore *DataStoreTransactorSession) RemoveUintArray(key [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.RemoveUintArray(&_DataStore.TransactOpts, key)
}

// SetAddress is a paid mutator transaction binding the contract method 0xca446dd9.
//
// Solidity: function setAddress(bytes32 key, address value) returns(address)
func (_DataStore *DataStoreTransactor) SetAddress(opts *bind.TransactOpts, key [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setAddress", key, value)
}

// SetAddress is a paid mutator transaction binding the contract method 0xca446dd9.
//
// Solidity: function setAddress(bytes32 key, address value) returns(address)
func (_DataStore *DataStoreSession) SetAddress(key [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.SetAddress(&_DataStore.TransactOpts, key, value)
}

// SetAddress is a paid mutator transaction binding the contract method 0xca446dd9.
//
// Solidity: function setAddress(bytes32 key, address value) returns(address)
func (_DataStore *DataStoreTransactorSession) SetAddress(key [32]byte, value common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.SetAddress(&_DataStore.TransactOpts, key, value)
}

// SetAddressArray is a paid mutator transaction binding the contract method 0xec672cf6.
//
// Solidity: function setAddressArray(bytes32 key, address[] value) returns()
func (_DataStore *DataStoreTransactor) SetAddressArray(opts *bind.TransactOpts, key [32]byte, value []common.Address) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setAddressArray", key, value)
}

// SetAddressArray is a paid mutator transaction binding the contract method 0xec672cf6.
//
// Solidity: function setAddressArray(bytes32 key, address[] value) returns()
func (_DataStore *DataStoreSession) SetAddressArray(key [32]byte, value []common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.SetAddressArray(&_DataStore.TransactOpts, key, value)
}

// SetAddressArray is a paid mutator transaction binding the contract method 0xec672cf6.
//
// Solidity: function setAddressArray(bytes32 key, address[] value) returns()
func (_DataStore *DataStoreTransactorSession) SetAddressArray(key [32]byte, value []common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.SetAddressArray(&_DataStore.TransactOpts, key, value)
}

// SetBool is a paid mutator transaction binding the contract method 0xabfdcced.
//
// Solidity: function setBool(bytes32 key, bool value) returns(bool)
func (_DataStore *DataStoreTransactor) SetBool(opts *bind.TransactOpts, key [32]byte, value bool) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setBool", key, value)
}

// SetBool is a paid mutator transaction binding the contract method 0xabfdcced.
//
// Solidity: function setBool(bytes32 key, bool value) returns(bool)
func (_DataStore *DataStoreSession) SetBool(key [32]byte, value bool) (*types.Transaction, error) {
	return _DataStore.Contract.SetBool(&_DataStore.TransactOpts, key, value)
}

// SetBool is a paid mutator transaction binding the contract method 0xabfdcced.
//
// Solidity: function setBool(bytes32 key, bool value) returns(bool)
func (_DataStore *DataStoreTransactorSession) SetBool(key [32]byte, value bool) (*types.Transaction, error) {
	return _DataStore.Contract.SetBool(&_DataStore.TransactOpts, key, value)
}

// SetBoolArray is a paid mutator transaction binding the contract method 0x35d4d407.
//
// Solidity: function setBoolArray(bytes32 key, bool[] value) returns()
func (_DataStore *DataStoreTransactor) SetBoolArray(opts *bind.TransactOpts, key [32]byte, value []bool) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setBoolArray", key, value)
}

// SetBoolArray is a paid mutator transaction binding the contract method 0x35d4d407.
//
// Solidity: function setBoolArray(bytes32 key, bool[] value) returns()
func (_DataStore *DataStoreSession) SetBoolArray(key [32]byte, value []bool) (*types.Transaction, error) {
	return _DataStore.Contract.SetBoolArray(&_DataStore.TransactOpts, key, value)
}

// SetBoolArray is a paid mutator transaction binding the contract method 0x35d4d407.
//
// Solidity: function setBoolArray(bytes32 key, bool[] value) returns()
func (_DataStore *DataStoreTransactorSession) SetBoolArray(key [32]byte, value []bool) (*types.Transaction, error) {
	return _DataStore.Contract.SetBoolArray(&_DataStore.TransactOpts, key, value)
}

// SetBytes32 is a paid mutator transaction binding the contract method 0x4e91db08.
//
// Solidity: function setBytes32(bytes32 key, bytes32 value) returns(bytes32)
func (_DataStore *DataStoreTransactor) SetBytes32(opts *bind.TransactOpts, key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setBytes32", key, value)
}

// SetBytes32 is a paid mutator transaction binding the contract method 0x4e91db08.
//
// Solidity: function setBytes32(bytes32 key, bytes32 value) returns(bytes32)
func (_DataStore *DataStoreSession) SetBytes32(key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.SetBytes32(&_DataStore.TransactOpts, key, value)
}

// SetBytes32 is a paid mutator transaction binding the contract method 0x4e91db08.
//
// Solidity: function setBytes32(bytes32 key, bytes32 value) returns(bytes32)
func (_DataStore *DataStoreTransactorSession) SetBytes32(key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.SetBytes32(&_DataStore.TransactOpts, key, value)
}

// SetBytes32Array is a paid mutator transaction binding the contract method 0x26004846.
//
// Solidity: function setBytes32Array(bytes32 key, bytes32[] value) returns()
func (_DataStore *DataStoreTransactor) SetBytes32Array(opts *bind.TransactOpts, key [32]byte, value [][32]byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setBytes32Array", key, value)
}

// SetBytes32Array is a paid mutator transaction binding the contract method 0x26004846.
//
// Solidity: function setBytes32Array(bytes32 key, bytes32[] value) returns()
func (_DataStore *DataStoreSession) SetBytes32Array(key [32]byte, value [][32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.SetBytes32Array(&_DataStore.TransactOpts, key, value)
}

// SetBytes32Array is a paid mutator transaction binding the contract method 0x26004846.
//
// Solidity: function setBytes32Array(bytes32 key, bytes32[] value) returns()
func (_DataStore *DataStoreTransactorSession) SetBytes32Array(key [32]byte, value [][32]byte) (*types.Transaction, error) {
	return _DataStore.Contract.SetBytes32Array(&_DataStore.TransactOpts, key, value)
}

// SetInt is a paid mutator transaction binding the contract method 0x3e49bed0.
//
// Solidity: function setInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactor) SetInt(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setInt", key, value)
}

// SetInt is a paid mutator transaction binding the contract method 0x3e49bed0.
//
// Solidity: function setInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreSession) SetInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetInt(&_DataStore.TransactOpts, key, value)
}

// SetInt is a paid mutator transaction binding the contract method 0x3e49bed0.
//
// Solidity: function setInt(bytes32 key, int256 value) returns(int256)
func (_DataStore *DataStoreTransactorSession) SetInt(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetInt(&_DataStore.TransactOpts, key, value)
}

// SetIntArray is a paid mutator transaction binding the contract method 0xa9fcf76b.
//
// Solidity: function setIntArray(bytes32 key, int256[] value) returns()
func (_DataStore *DataStoreTransactor) SetIntArray(opts *bind.TransactOpts, key [32]byte, value []*big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setIntArray", key, value)
}

// SetIntArray is a paid mutator transaction binding the contract method 0xa9fcf76b.
//
// Solidity: function setIntArray(bytes32 key, int256[] value) returns()
func (_DataStore *DataStoreSession) SetIntArray(key [32]byte, value []*big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetIntArray(&_DataStore.TransactOpts, key, value)
}

// SetIntArray is a paid mutator transaction binding the contract method 0xa9fcf76b.
//
// Solidity: function setIntArray(bytes32 key, int256[] value) returns()
func (_DataStore *DataStoreTransactorSession) SetIntArray(key [32]byte, value []*big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetIntArray(&_DataStore.TransactOpts, key, value)
}

// SetString is a paid mutator transaction binding the contract method 0x6e899550.
//
// Solidity: function setString(bytes32 key, string value) returns(string)
func (_DataStore *DataStoreTransactor) SetString(opts *bind.TransactOpts, key [32]byte, value string) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setString", key, value)
}

// SetString is a paid mutator transaction binding the contract method 0x6e899550.
//
// Solidity: function setString(bytes32 key, string value) returns(string)
func (_DataStore *DataStoreSession) SetString(key [32]byte, value string) (*types.Transaction, error) {
	return _DataStore.Contract.SetString(&_DataStore.TransactOpts, key, value)
}

// SetString is a paid mutator transaction binding the contract method 0x6e899550.
//
// Solidity: function setString(bytes32 key, string value) returns(string)
func (_DataStore *DataStoreTransactorSession) SetString(key [32]byte, value string) (*types.Transaction, error) {
	return _DataStore.Contract.SetString(&_DataStore.TransactOpts, key, value)
}

// SetStringArray is a paid mutator transaction binding the contract method 0x88021a72.
//
// Solidity: function setStringArray(bytes32 key, string[] value) returns()
func (_DataStore *DataStoreTransactor) SetStringArray(opts *bind.TransactOpts, key [32]byte, value []string) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setStringArray", key, value)
}

// SetStringArray is a paid mutator transaction binding the contract method 0x88021a72.
//
// Solidity: function setStringArray(bytes32 key, string[] value) returns()
func (_DataStore *DataStoreSession) SetStringArray(key [32]byte, value []string) (*types.Transaction, error) {
	return _DataStore.Contract.SetStringArray(&_DataStore.TransactOpts, key, value)
}

// SetStringArray is a paid mutator transaction binding the contract method 0x88021a72.
//
// Solidity: function setStringArray(bytes32 key, string[] value) returns()
func (_DataStore *DataStoreTransactorSession) SetStringArray(key [32]byte, value []string) (*types.Transaction, error) {
	return _DataStore.Contract.SetStringArray(&_DataStore.TransactOpts, key, value)
}

// SetUint is a paid mutator transaction binding the contract method 0xe2a4853a.
//
// Solidity: function setUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactor) SetUint(opts *bind.TransactOpts, key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setUint", key, value)
}

// SetUint is a paid mutator transaction binding the contract method 0xe2a4853a.
//
// Solidity: function setUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreSession) SetUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetUint(&_DataStore.TransactOpts, key, value)
}

// SetUint is a paid mutator transaction binding the contract method 0xe2a4853a.
//
// Solidity: function setUint(bytes32 key, uint256 value) returns(uint256)
func (_DataStore *DataStoreTransactorSession) SetUint(key [32]byte, value *big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetUint(&_DataStore.TransactOpts, key, value)
}

// SetUintArray is a paid mutator transaction binding the contract method 0x5eb07dbd.
//
// Solidity: function setUintArray(bytes32 key, uint256[] value) returns()
func (_DataStore *DataStoreTransactor) SetUintArray(opts *bind.TransactOpts, key [32]byte, value []*big.Int) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setUintArray", key, value)
}

// SetUintArray is a paid mutator transaction binding the contract method 0x5eb07dbd.
//
// Solidity: function setUintArray(bytes32 key, uint256[] value) returns()
func (_DataStore *DataStoreSession) SetUintArray(key [32]byte, value []*big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetUintArray(&_DataStore.TransactOpts, key, value)
}

// SetUintArray is a paid mutator transaction binding the contract method 0x5eb07dbd.
//
// Solidity: function setUintArray(bytes32 key, uint256[] value) returns()
func (_DataStore *DataStoreTransactorSession) SetUintArray(key [32]byte, value []*big.Int) (*types.Transaction, error) {
	return _DataStore.Contract.SetUintArray(&_DataStore.TransactOpts, key, value)
}

// DataStoreInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DataStore contract.
type DataStoreInitializedIterator struct {
	Event *DataStoreInitialized // Event containing the contract specifics and raw log

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
func (it *DataStoreInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataStoreInitialized)
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
		it.Event = new(DataStoreInitialized)
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
func (it *DataStoreInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataStoreInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataStoreInitialized represents a Initialized event raised by the DataStore contract.
type DataStoreInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataStore *DataStoreFilterer) FilterInitialized(opts *bind.FilterOpts) (*DataStoreInitializedIterator, error) {

	logs, sub, err := _DataStore.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DataStoreInitializedIterator{contract: _DataStore.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataStore *DataStoreFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DataStoreInitialized) (event.Subscription, error) {

	logs, sub, err := _DataStore.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataStoreInitialized)
				if err := _DataStore.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataStore *DataStoreFilterer) ParseInitialized(log types.Log) (*DataStoreInitialized, error) {
	event := new(DataStoreInitialized)
	if err := _DataStore.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
