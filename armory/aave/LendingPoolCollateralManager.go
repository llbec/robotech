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

// LendingPoolCollateralManagerABI is the input ABI used to generate the binding from.
const LendingPoolCollateralManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"collateral\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"principal\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidatedCollateralAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"receiveAToken\",\"type\":\"bool\"}],\"name\":\"LiquidationCall\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"ReserveUsedAsCollateralDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"ReserveUsedAsCollateralEnabled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collateralAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"debtAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"receiveAToken\",\"type\":\"bool\"}],\"name\":\"liquidationCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// LendingPoolCollateralManager is an auto generated Go binding around an Ethereum contract.
type LendingPoolCollateralManager struct {
	LendingPoolCollateralManagerCaller     // Read-only binding to the contract
	LendingPoolCollateralManagerTransactor // Write-only binding to the contract
	LendingPoolCollateralManagerFilterer   // Log filterer for contract events
}

// LendingPoolCollateralManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type LendingPoolCollateralManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingPoolCollateralManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LendingPoolCollateralManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingPoolCollateralManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LendingPoolCollateralManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingPoolCollateralManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LendingPoolCollateralManagerSession struct {
	Contract     *LendingPoolCollateralManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LendingPoolCollateralManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LendingPoolCollateralManagerCallerSession struct {
	Contract *LendingPoolCollateralManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// LendingPoolCollateralManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LendingPoolCollateralManagerTransactorSession struct {
	Contract     *LendingPoolCollateralManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// LendingPoolCollateralManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type LendingPoolCollateralManagerRaw struct {
	Contract *LendingPoolCollateralManager // Generic contract binding to access the raw methods on
}

// LendingPoolCollateralManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LendingPoolCollateralManagerCallerRaw struct {
	Contract *LendingPoolCollateralManagerCaller // Generic read-only contract binding to access the raw methods on
}

// LendingPoolCollateralManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LendingPoolCollateralManagerTransactorRaw struct {
	Contract *LendingPoolCollateralManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLendingPoolCollateralManager creates a new instance of LendingPoolCollateralManager, bound to a specific deployed contract.
func NewLendingPoolCollateralManager(address common.Address, backend bind.ContractBackend) (*LendingPoolCollateralManager, error) {
	contract, err := bindLendingPoolCollateralManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManager{LendingPoolCollateralManagerCaller: LendingPoolCollateralManagerCaller{contract: contract}, LendingPoolCollateralManagerTransactor: LendingPoolCollateralManagerTransactor{contract: contract}, LendingPoolCollateralManagerFilterer: LendingPoolCollateralManagerFilterer{contract: contract}}, nil
}

// NewLendingPoolCollateralManagerCaller creates a new read-only instance of LendingPoolCollateralManager, bound to a specific deployed contract.
func NewLendingPoolCollateralManagerCaller(address common.Address, caller bind.ContractCaller) (*LendingPoolCollateralManagerCaller, error) {
	contract, err := bindLendingPoolCollateralManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManagerCaller{contract: contract}, nil
}

// NewLendingPoolCollateralManagerTransactor creates a new write-only instance of LendingPoolCollateralManager, bound to a specific deployed contract.
func NewLendingPoolCollateralManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*LendingPoolCollateralManagerTransactor, error) {
	contract, err := bindLendingPoolCollateralManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManagerTransactor{contract: contract}, nil
}

// NewLendingPoolCollateralManagerFilterer creates a new log filterer instance of LendingPoolCollateralManager, bound to a specific deployed contract.
func NewLendingPoolCollateralManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*LendingPoolCollateralManagerFilterer, error) {
	contract, err := bindLendingPoolCollateralManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManagerFilterer{contract: contract}, nil
}

// bindLendingPoolCollateralManager binds a generic wrapper to an already deployed contract.
func bindLendingPoolCollateralManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LendingPoolCollateralManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LendingPoolCollateralManager.Contract.LendingPoolCollateralManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.Contract.LendingPoolCollateralManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.Contract.LendingPoolCollateralManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LendingPoolCollateralManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.Contract.contract.Transact(opts, method, params...)
}

// LiquidationCall is a paid mutator transaction binding the contract method 0x00a718a9.
//
// Solidity: function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover, bool receiveAToken) returns(uint256, string)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerTransactor) LiquidationCall(opts *bind.TransactOpts, collateralAsset common.Address, debtAsset common.Address, user common.Address, debtToCover *big.Int, receiveAToken bool) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.contract.Transact(opts, "liquidationCall", collateralAsset, debtAsset, user, debtToCover, receiveAToken)
}

// LiquidationCall is a paid mutator transaction binding the contract method 0x00a718a9.
//
// Solidity: function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover, bool receiveAToken) returns(uint256, string)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerSession) LiquidationCall(collateralAsset common.Address, debtAsset common.Address, user common.Address, debtToCover *big.Int, receiveAToken bool) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.Contract.LiquidationCall(&_LendingPoolCollateralManager.TransactOpts, collateralAsset, debtAsset, user, debtToCover, receiveAToken)
}

// LiquidationCall is a paid mutator transaction binding the contract method 0x00a718a9.
//
// Solidity: function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover, bool receiveAToken) returns(uint256, string)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerTransactorSession) LiquidationCall(collateralAsset common.Address, debtAsset common.Address, user common.Address, debtToCover *big.Int, receiveAToken bool) (*types.Transaction, error) {
	return _LendingPoolCollateralManager.Contract.LiquidationCall(&_LendingPoolCollateralManager.TransactOpts, collateralAsset, debtAsset, user, debtToCover, receiveAToken)
}

// LendingPoolCollateralManagerLiquidationCallIterator is returned from FilterLiquidationCall and is used to iterate over the raw logs and unpacked data for LiquidationCall events raised by the LendingPoolCollateralManager contract.
type LendingPoolCollateralManagerLiquidationCallIterator struct {
	Event *LendingPoolCollateralManagerLiquidationCall // Event containing the contract specifics and raw log

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
func (it *LendingPoolCollateralManagerLiquidationCallIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LendingPoolCollateralManagerLiquidationCall)
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
		it.Event = new(LendingPoolCollateralManagerLiquidationCall)
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
func (it *LendingPoolCollateralManagerLiquidationCallIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LendingPoolCollateralManagerLiquidationCallIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LendingPoolCollateralManagerLiquidationCall represents a LiquidationCall event raised by the LendingPoolCollateralManager contract.
type LendingPoolCollateralManagerLiquidationCall struct {
	Collateral                 common.Address
	Principal                  common.Address
	User                       common.Address
	DebtToCover                *big.Int
	LiquidatedCollateralAmount *big.Int
	Liquidator                 common.Address
	ReceiveAToken              bool
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterLiquidationCall is a free log retrieval operation binding the contract event 0xe413a321e8681d831f4dbccbca790d2952b56f977908e45be37335533e005286.
//
// Solidity: event LiquidationCall(address indexed collateral, address indexed principal, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) FilterLiquidationCall(opts *bind.FilterOpts, collateral []common.Address, principal []common.Address, user []common.Address) (*LendingPoolCollateralManagerLiquidationCallIterator, error) {

	var collateralRule []interface{}
	for _, collateralItem := range collateral {
		collateralRule = append(collateralRule, collateralItem)
	}
	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _LendingPoolCollateralManager.contract.FilterLogs(opts, "LiquidationCall", collateralRule, principalRule, userRule)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManagerLiquidationCallIterator{contract: _LendingPoolCollateralManager.contract, event: "LiquidationCall", logs: logs, sub: sub}, nil
}

// WatchLiquidationCall is a free log subscription operation binding the contract event 0xe413a321e8681d831f4dbccbca790d2952b56f977908e45be37335533e005286.
//
// Solidity: event LiquidationCall(address indexed collateral, address indexed principal, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) WatchLiquidationCall(opts *bind.WatchOpts, sink chan<- *LendingPoolCollateralManagerLiquidationCall, collateral []common.Address, principal []common.Address, user []common.Address) (event.Subscription, error) {

	var collateralRule []interface{}
	for _, collateralItem := range collateral {
		collateralRule = append(collateralRule, collateralItem)
	}
	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _LendingPoolCollateralManager.contract.WatchLogs(opts, "LiquidationCall", collateralRule, principalRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LendingPoolCollateralManagerLiquidationCall)
				if err := _LendingPoolCollateralManager.contract.UnpackLog(event, "LiquidationCall", log); err != nil {
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

// ParseLiquidationCall is a log parse operation binding the contract event 0xe413a321e8681d831f4dbccbca790d2952b56f977908e45be37335533e005286.
//
// Solidity: event LiquidationCall(address indexed collateral, address indexed principal, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) ParseLiquidationCall(log types.Log) (*LendingPoolCollateralManagerLiquidationCall, error) {
	event := new(LendingPoolCollateralManagerLiquidationCall)
	if err := _LendingPoolCollateralManager.contract.UnpackLog(event, "LiquidationCall", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator is returned from FilterReserveUsedAsCollateralDisabled and is used to iterate over the raw logs and unpacked data for ReserveUsedAsCollateralDisabled events raised by the LendingPoolCollateralManager contract.
type LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator struct {
	Event *LendingPoolCollateralManagerReserveUsedAsCollateralDisabled // Event containing the contract specifics and raw log

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
func (it *LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LendingPoolCollateralManagerReserveUsedAsCollateralDisabled)
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
		it.Event = new(LendingPoolCollateralManagerReserveUsedAsCollateralDisabled)
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
func (it *LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LendingPoolCollateralManagerReserveUsedAsCollateralDisabled represents a ReserveUsedAsCollateralDisabled event raised by the LendingPoolCollateralManager contract.
type LendingPoolCollateralManagerReserveUsedAsCollateralDisabled struct {
	Reserve common.Address
	User    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveUsedAsCollateralDisabled is a free log retrieval operation binding the contract event 0x44c58d81365b66dd4b1a7f36c25aa97b8c71c361ee4937adc1a00000227db5dd.
//
// Solidity: event ReserveUsedAsCollateralDisabled(address indexed reserve, address indexed user)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) FilterReserveUsedAsCollateralDisabled(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _LendingPoolCollateralManager.contract.FilterLogs(opts, "ReserveUsedAsCollateralDisabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManagerReserveUsedAsCollateralDisabledIterator{contract: _LendingPoolCollateralManager.contract, event: "ReserveUsedAsCollateralDisabled", logs: logs, sub: sub}, nil
}

// WatchReserveUsedAsCollateralDisabled is a free log subscription operation binding the contract event 0x44c58d81365b66dd4b1a7f36c25aa97b8c71c361ee4937adc1a00000227db5dd.
//
// Solidity: event ReserveUsedAsCollateralDisabled(address indexed reserve, address indexed user)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) WatchReserveUsedAsCollateralDisabled(opts *bind.WatchOpts, sink chan<- *LendingPoolCollateralManagerReserveUsedAsCollateralDisabled, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _LendingPoolCollateralManager.contract.WatchLogs(opts, "ReserveUsedAsCollateralDisabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LendingPoolCollateralManagerReserveUsedAsCollateralDisabled)
				if err := _LendingPoolCollateralManager.contract.UnpackLog(event, "ReserveUsedAsCollateralDisabled", log); err != nil {
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

// ParseReserveUsedAsCollateralDisabled is a log parse operation binding the contract event 0x44c58d81365b66dd4b1a7f36c25aa97b8c71c361ee4937adc1a00000227db5dd.
//
// Solidity: event ReserveUsedAsCollateralDisabled(address indexed reserve, address indexed user)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) ParseReserveUsedAsCollateralDisabled(log types.Log) (*LendingPoolCollateralManagerReserveUsedAsCollateralDisabled, error) {
	event := new(LendingPoolCollateralManagerReserveUsedAsCollateralDisabled)
	if err := _LendingPoolCollateralManager.contract.UnpackLog(event, "ReserveUsedAsCollateralDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator is returned from FilterReserveUsedAsCollateralEnabled and is used to iterate over the raw logs and unpacked data for ReserveUsedAsCollateralEnabled events raised by the LendingPoolCollateralManager contract.
type LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator struct {
	Event *LendingPoolCollateralManagerReserveUsedAsCollateralEnabled // Event containing the contract specifics and raw log

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
func (it *LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LendingPoolCollateralManagerReserveUsedAsCollateralEnabled)
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
		it.Event = new(LendingPoolCollateralManagerReserveUsedAsCollateralEnabled)
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
func (it *LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LendingPoolCollateralManagerReserveUsedAsCollateralEnabled represents a ReserveUsedAsCollateralEnabled event raised by the LendingPoolCollateralManager contract.
type LendingPoolCollateralManagerReserveUsedAsCollateralEnabled struct {
	Reserve common.Address
	User    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveUsedAsCollateralEnabled is a free log retrieval operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) FilterReserveUsedAsCollateralEnabled(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _LendingPoolCollateralManager.contract.FilterLogs(opts, "ReserveUsedAsCollateralEnabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCollateralManagerReserveUsedAsCollateralEnabledIterator{contract: _LendingPoolCollateralManager.contract, event: "ReserveUsedAsCollateralEnabled", logs: logs, sub: sub}, nil
}

// WatchReserveUsedAsCollateralEnabled is a free log subscription operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) WatchReserveUsedAsCollateralEnabled(opts *bind.WatchOpts, sink chan<- *LendingPoolCollateralManagerReserveUsedAsCollateralEnabled, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _LendingPoolCollateralManager.contract.WatchLogs(opts, "ReserveUsedAsCollateralEnabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LendingPoolCollateralManagerReserveUsedAsCollateralEnabled)
				if err := _LendingPoolCollateralManager.contract.UnpackLog(event, "ReserveUsedAsCollateralEnabled", log); err != nil {
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

// ParseReserveUsedAsCollateralEnabled is a log parse operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_LendingPoolCollateralManager *LendingPoolCollateralManagerFilterer) ParseReserveUsedAsCollateralEnabled(log types.Log) (*LendingPoolCollateralManagerReserveUsedAsCollateralEnabled, error) {
	event := new(LendingPoolCollateralManagerReserveUsedAsCollateralEnabled)
	if err := _LendingPoolCollateralManager.contract.UnpackLog(event, "ReserveUsedAsCollateralEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
