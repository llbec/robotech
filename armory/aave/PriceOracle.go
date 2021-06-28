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

// PriceOracleABI is the input ABI used to generate the binding from.
const PriceOracleABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AssetPriceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"EthPriceUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"}],\"name\":\"getAssetPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEthUsdPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setAssetPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setEthUsdPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address _asset) view returns(uint256)
func (_PriceOracle *PriceOracleCaller) GetAssetPrice(opts *bind.CallOpts, _asset common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PriceOracle.contract.Call(opts, &out, "getAssetPrice", _asset)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address _asset) view returns(uint256)
func (_PriceOracle *PriceOracleSession) GetAssetPrice(_asset common.Address) (*big.Int, error) {
	return _PriceOracle.Contract.GetAssetPrice(&_PriceOracle.CallOpts, _asset)
}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address _asset) view returns(uint256)
func (_PriceOracle *PriceOracleCallerSession) GetAssetPrice(_asset common.Address) (*big.Int, error) {
	return _PriceOracle.Contract.GetAssetPrice(&_PriceOracle.CallOpts, _asset)
}

// GetEthUsdPrice is a free data retrieval call binding the contract method 0xa0a8045e.
//
// Solidity: function getEthUsdPrice() view returns(uint256)
func (_PriceOracle *PriceOracleCaller) GetEthUsdPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceOracle.contract.Call(opts, &out, "getEthUsdPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEthUsdPrice is a free data retrieval call binding the contract method 0xa0a8045e.
//
// Solidity: function getEthUsdPrice() view returns(uint256)
func (_PriceOracle *PriceOracleSession) GetEthUsdPrice() (*big.Int, error) {
	return _PriceOracle.Contract.GetEthUsdPrice(&_PriceOracle.CallOpts)
}

// GetEthUsdPrice is a free data retrieval call binding the contract method 0xa0a8045e.
//
// Solidity: function getEthUsdPrice() view returns(uint256)
func (_PriceOracle *PriceOracleCallerSession) GetEthUsdPrice() (*big.Int, error) {
	return _PriceOracle.Contract.GetEthUsdPrice(&_PriceOracle.CallOpts)
}

// SetAssetPrice is a paid mutator transaction binding the contract method 0x51323f72.
//
// Solidity: function setAssetPrice(address _asset, uint256 _price) returns()
func (_PriceOracle *PriceOracleTransactor) SetAssetPrice(opts *bind.TransactOpts, _asset common.Address, _price *big.Int) (*types.Transaction, error) {
	return _PriceOracle.contract.Transact(opts, "setAssetPrice", _asset, _price)
}

// SetAssetPrice is a paid mutator transaction binding the contract method 0x51323f72.
//
// Solidity: function setAssetPrice(address _asset, uint256 _price) returns()
func (_PriceOracle *PriceOracleSession) SetAssetPrice(_asset common.Address, _price *big.Int) (*types.Transaction, error) {
	return _PriceOracle.Contract.SetAssetPrice(&_PriceOracle.TransactOpts, _asset, _price)
}

// SetAssetPrice is a paid mutator transaction binding the contract method 0x51323f72.
//
// Solidity: function setAssetPrice(address _asset, uint256 _price) returns()
func (_PriceOracle *PriceOracleTransactorSession) SetAssetPrice(_asset common.Address, _price *big.Int) (*types.Transaction, error) {
	return _PriceOracle.Contract.SetAssetPrice(&_PriceOracle.TransactOpts, _asset, _price)
}

// SetEthUsdPrice is a paid mutator transaction binding the contract method 0xb951883a.
//
// Solidity: function setEthUsdPrice(uint256 _price) returns()
func (_PriceOracle *PriceOracleTransactor) SetEthUsdPrice(opts *bind.TransactOpts, _price *big.Int) (*types.Transaction, error) {
	return _PriceOracle.contract.Transact(opts, "setEthUsdPrice", _price)
}

// SetEthUsdPrice is a paid mutator transaction binding the contract method 0xb951883a.
//
// Solidity: function setEthUsdPrice(uint256 _price) returns()
func (_PriceOracle *PriceOracleSession) SetEthUsdPrice(_price *big.Int) (*types.Transaction, error) {
	return _PriceOracle.Contract.SetEthUsdPrice(&_PriceOracle.TransactOpts, _price)
}

// SetEthUsdPrice is a paid mutator transaction binding the contract method 0xb951883a.
//
// Solidity: function setEthUsdPrice(uint256 _price) returns()
func (_PriceOracle *PriceOracleTransactorSession) SetEthUsdPrice(_price *big.Int) (*types.Transaction, error) {
	return _PriceOracle.Contract.SetEthUsdPrice(&_PriceOracle.TransactOpts, _price)
}

// PriceOracleAssetPriceUpdatedIterator is returned from FilterAssetPriceUpdated and is used to iterate over the raw logs and unpacked data for AssetPriceUpdated events raised by the PriceOracle contract.
type PriceOracleAssetPriceUpdatedIterator struct {
	Event *PriceOracleAssetPriceUpdated // Event containing the contract specifics and raw log

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
func (it *PriceOracleAssetPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceOracleAssetPriceUpdated)
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
		it.Event = new(PriceOracleAssetPriceUpdated)
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
func (it *PriceOracleAssetPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceOracleAssetPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceOracleAssetPriceUpdated represents a AssetPriceUpdated event raised by the PriceOracle contract.
type PriceOracleAssetPriceUpdated struct {
	Asset     common.Address
	Price     *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAssetPriceUpdated is a free log retrieval operation binding the contract event 0xce6e0b57367bae95ca7198e1172f653ea64a645c16ab586b4cefa9237bfc2d92.
//
// Solidity: event AssetPriceUpdated(address _asset, uint256 _price, uint256 timestamp)
func (_PriceOracle *PriceOracleFilterer) FilterAssetPriceUpdated(opts *bind.FilterOpts) (*PriceOracleAssetPriceUpdatedIterator, error) {

	logs, sub, err := _PriceOracle.contract.FilterLogs(opts, "AssetPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &PriceOracleAssetPriceUpdatedIterator{contract: _PriceOracle.contract, event: "AssetPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchAssetPriceUpdated is a free log subscription operation binding the contract event 0xce6e0b57367bae95ca7198e1172f653ea64a645c16ab586b4cefa9237bfc2d92.
//
// Solidity: event AssetPriceUpdated(address _asset, uint256 _price, uint256 timestamp)
func (_PriceOracle *PriceOracleFilterer) WatchAssetPriceUpdated(opts *bind.WatchOpts, sink chan<- *PriceOracleAssetPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _PriceOracle.contract.WatchLogs(opts, "AssetPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceOracleAssetPriceUpdated)
				if err := _PriceOracle.contract.UnpackLog(event, "AssetPriceUpdated", log); err != nil {
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

// ParseAssetPriceUpdated is a log parse operation binding the contract event 0xce6e0b57367bae95ca7198e1172f653ea64a645c16ab586b4cefa9237bfc2d92.
//
// Solidity: event AssetPriceUpdated(address _asset, uint256 _price, uint256 timestamp)
func (_PriceOracle *PriceOracleFilterer) ParseAssetPriceUpdated(log types.Log) (*PriceOracleAssetPriceUpdated, error) {
	event := new(PriceOracleAssetPriceUpdated)
	if err := _PriceOracle.contract.UnpackLog(event, "AssetPriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceOracleEthPriceUpdatedIterator is returned from FilterEthPriceUpdated and is used to iterate over the raw logs and unpacked data for EthPriceUpdated events raised by the PriceOracle contract.
type PriceOracleEthPriceUpdatedIterator struct {
	Event *PriceOracleEthPriceUpdated // Event containing the contract specifics and raw log

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
func (it *PriceOracleEthPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceOracleEthPriceUpdated)
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
		it.Event = new(PriceOracleEthPriceUpdated)
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
func (it *PriceOracleEthPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceOracleEthPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceOracleEthPriceUpdated represents a EthPriceUpdated event raised by the PriceOracle contract.
type PriceOracleEthPriceUpdated struct {
	Price     *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEthPriceUpdated is a free log retrieval operation binding the contract event 0xb4f35977939fa8b5ffe552d517a8ff5223046b1fdd3ee0068ae38d1e2b8d0016.
//
// Solidity: event EthPriceUpdated(uint256 _price, uint256 timestamp)
func (_PriceOracle *PriceOracleFilterer) FilterEthPriceUpdated(opts *bind.FilterOpts) (*PriceOracleEthPriceUpdatedIterator, error) {

	logs, sub, err := _PriceOracle.contract.FilterLogs(opts, "EthPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &PriceOracleEthPriceUpdatedIterator{contract: _PriceOracle.contract, event: "EthPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchEthPriceUpdated is a free log subscription operation binding the contract event 0xb4f35977939fa8b5ffe552d517a8ff5223046b1fdd3ee0068ae38d1e2b8d0016.
//
// Solidity: event EthPriceUpdated(uint256 _price, uint256 timestamp)
func (_PriceOracle *PriceOracleFilterer) WatchEthPriceUpdated(opts *bind.WatchOpts, sink chan<- *PriceOracleEthPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _PriceOracle.contract.WatchLogs(opts, "EthPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceOracleEthPriceUpdated)
				if err := _PriceOracle.contract.UnpackLog(event, "EthPriceUpdated", log); err != nil {
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

// ParseEthPriceUpdated is a log parse operation binding the contract event 0xb4f35977939fa8b5ffe552d517a8ff5223046b1fdd3ee0068ae38d1e2b8d0016.
//
// Solidity: event EthPriceUpdated(uint256 _price, uint256 timestamp)
func (_PriceOracle *PriceOracleFilterer) ParseEthPriceUpdated(log types.Log) (*PriceOracleEthPriceUpdated, error) {
	event := new(PriceOracleEthPriceUpdated)
	if err := _PriceOracle.contract.UnpackLog(event, "EthPriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
