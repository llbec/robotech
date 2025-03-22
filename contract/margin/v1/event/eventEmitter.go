// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EventEmitter

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

// EventLiquidation is an auto generated low-level Go binding around an user-defined struct.
type EventLiquidation struct {
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
}

// EventEmitterMetaData contains all meta data concerning the EventEmitter contract.
var EventEmitterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"_roleStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"role\",\"type\":\"string\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"adder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"name\":\"Add\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"borrowAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"Borrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"scaledUnclaimedFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidityIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unclaimedFee\",\"type\":\"uint256\"}],\"name\":\"ClaimFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"}],\"name\":\"Close\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"marginLevel\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"marginLevelLiquidationThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalCollateralUsd\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalDebtUsd\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memePrice\",\"type\":\"uint256\"}],\"name\":\"Liquidation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"createdTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDecimals\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDecimals\",\"type\":\"uint256\"}],\"name\":\"PoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidityRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidityIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"}],\"name\":\"PoolUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actionType\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"Position\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeAmount\",\"type\":\"uint256\"}],\"name\":\"Remove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"repayer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"repayAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"Repay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"withdrawer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"supplier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"baseAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"name\":\"emitAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"borrowAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"internalType\":\"structEvent.Liquidation\",\"name\":\"liquidation\",\"type\":\"tuple\"}],\"name\":\"emitBorrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"scaledUnclaimedFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unclaimedFee\",\"type\":\"uint256\"}],\"name\":\"emitClaimFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"}],\"name\":\"emitClose\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"emitDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"marginLevel\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"marginLevelLiquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateralUsd\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtUsd\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memePrice\",\"type\":\"uint256\"}],\"name\":\"emitLiquidation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"createdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDecimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDecimals\",\"type\":\"uint256\"}],\"name\":\"emitPoolCreated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidityRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowIndex\",\"type\":\"uint256\"}],\"name\":\"emitPoolUpdated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"actionType\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"emitPosition\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"baseAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeAmount\",\"type\":\"uint256\"}],\"name\":\"emitRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"repayer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tokenIndex\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"repayAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"internalType\":\"structEvent.Liquidation\",\"name\":\"liquidation\",\"type\":\"tuple\"}],\"name\":\"emitRepay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"internalType\":\"structEvent.Liquidation\",\"name\":\"liquidation\",\"type\":\"tuple\"}],\"name\":\"emitSwap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"withdrawer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"memeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"positionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"baseCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseDebtScaled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"memeDebtScaled\",\"type\":\"uint256\"}],\"name\":\"emitWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleStore\",\"outputs\":[{\"internalType\":\"contractRoleStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EventEmitterABI is the input ABI used to generate the binding from.
// Deprecated: Use EventEmitterMetaData.ABI instead.
var EventEmitterABI = EventEmitterMetaData.ABI

// EventEmitter is an auto generated Go binding around an Ethereum contract.
type EventEmitter struct {
	EventEmitterCaller     // Read-only binding to the contract
	EventEmitterTransactor // Write-only binding to the contract
	EventEmitterFilterer   // Log filterer for contract events
}

// EventEmitterCaller is an auto generated read-only Go binding around an Ethereum contract.
type EventEmitterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EventEmitterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EventEmitterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EventEmitterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EventEmitterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EventEmitterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EventEmitterSession struct {
	Contract     *EventEmitter     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EventEmitterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EventEmitterCallerSession struct {
	Contract *EventEmitterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// EventEmitterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EventEmitterTransactorSession struct {
	Contract     *EventEmitterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// EventEmitterRaw is an auto generated low-level Go binding around an Ethereum contract.
type EventEmitterRaw struct {
	Contract *EventEmitter // Generic contract binding to access the raw methods on
}

// EventEmitterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EventEmitterCallerRaw struct {
	Contract *EventEmitterCaller // Generic read-only contract binding to access the raw methods on
}

// EventEmitterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EventEmitterTransactorRaw struct {
	Contract *EventEmitterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEventEmitter creates a new instance of EventEmitter, bound to a specific deployed contract.
func NewEventEmitter(address common.Address, backend bind.ContractBackend) (*EventEmitter, error) {
	contract, err := bindEventEmitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EventEmitter{EventEmitterCaller: EventEmitterCaller{contract: contract}, EventEmitterTransactor: EventEmitterTransactor{contract: contract}, EventEmitterFilterer: EventEmitterFilterer{contract: contract}}, nil
}

// NewEventEmitterCaller creates a new read-only instance of EventEmitter, bound to a specific deployed contract.
func NewEventEmitterCaller(address common.Address, caller bind.ContractCaller) (*EventEmitterCaller, error) {
	contract, err := bindEventEmitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EventEmitterCaller{contract: contract}, nil
}

// NewEventEmitterTransactor creates a new write-only instance of EventEmitter, bound to a specific deployed contract.
func NewEventEmitterTransactor(address common.Address, transactor bind.ContractTransactor) (*EventEmitterTransactor, error) {
	contract, err := bindEventEmitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EventEmitterTransactor{contract: contract}, nil
}

// NewEventEmitterFilterer creates a new log filterer instance of EventEmitter, bound to a specific deployed contract.
func NewEventEmitterFilterer(address common.Address, filterer bind.ContractFilterer) (*EventEmitterFilterer, error) {
	contract, err := bindEventEmitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EventEmitterFilterer{contract: contract}, nil
}

// bindEventEmitter binds a generic wrapper to an already deployed contract.
func bindEventEmitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EventEmitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EventEmitter *EventEmitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EventEmitter.Contract.EventEmitterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EventEmitter *EventEmitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EventEmitter.Contract.EventEmitterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EventEmitter *EventEmitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EventEmitter.Contract.EventEmitterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EventEmitter *EventEmitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EventEmitter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EventEmitter *EventEmitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EventEmitter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EventEmitter *EventEmitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EventEmitter.Contract.contract.Transact(opts, method, params...)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_EventEmitter *EventEmitterCaller) RoleStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EventEmitter.contract.Call(opts, &out, "roleStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_EventEmitter *EventEmitterSession) RoleStore() (common.Address, error) {
	return _EventEmitter.Contract.RoleStore(&_EventEmitter.CallOpts)
}

// RoleStore is a free data retrieval call binding the contract method 0x4a4a7b04.
//
// Solidity: function roleStore() view returns(address)
func (_EventEmitter *EventEmitterCallerSession) RoleStore() (common.Address, error) {
	return _EventEmitter.Contract.RoleStore(&_EventEmitter.CallOpts)
}

// EmitAdd is a paid mutator transaction binding the contract method 0x33b365b4.
//
// Solidity: function emitAdd(address supplier, address baseToken, address memeToken, address to, uint256 baseAmount, uint256 memeAmount, uint256 liquidity) returns()
func (_EventEmitter *EventEmitterTransactor) EmitAdd(opts *bind.TransactOpts, supplier common.Address, baseToken common.Address, memeToken common.Address, to common.Address, baseAmount *big.Int, memeAmount *big.Int, liquidity *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitAdd", supplier, baseToken, memeToken, to, baseAmount, memeAmount, liquidity)
}

// EmitAdd is a paid mutator transaction binding the contract method 0x33b365b4.
//
// Solidity: function emitAdd(address supplier, address baseToken, address memeToken, address to, uint256 baseAmount, uint256 memeAmount, uint256 liquidity) returns()
func (_EventEmitter *EventEmitterSession) EmitAdd(supplier common.Address, baseToken common.Address, memeToken common.Address, to common.Address, baseAmount *big.Int, memeAmount *big.Int, liquidity *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitAdd(&_EventEmitter.TransactOpts, supplier, baseToken, memeToken, to, baseAmount, memeAmount, liquidity)
}

// EmitAdd is a paid mutator transaction binding the contract method 0x33b365b4.
//
// Solidity: function emitAdd(address supplier, address baseToken, address memeToken, address to, uint256 baseAmount, uint256 memeAmount, uint256 liquidity) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitAdd(supplier common.Address, baseToken common.Address, memeToken common.Address, to common.Address, baseAmount *big.Int, memeAmount *big.Int, liquidity *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitAdd(&_EventEmitter.TransactOpts, supplier, baseToken, memeToken, to, baseAmount, memeAmount, liquidity)
}

// EmitBorrow is a paid mutator transaction binding the contract method 0xea34a577.
//
// Solidity: function emitBorrow(address borrower, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 borrowAmount, uint256 borrowRate, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterTransactor) EmitBorrow(opts *bind.TransactOpts, borrower common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, tokenIndex uint8, borrowAmount *big.Int, borrowRate *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitBorrow", borrower, baseToken, memeToken, positionId, tokenIndex, borrowAmount, borrowRate, liquidation)
}

// EmitBorrow is a paid mutator transaction binding the contract method 0xea34a577.
//
// Solidity: function emitBorrow(address borrower, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 borrowAmount, uint256 borrowRate, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterSession) EmitBorrow(borrower common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, tokenIndex uint8, borrowAmount *big.Int, borrowRate *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitBorrow(&_EventEmitter.TransactOpts, borrower, baseToken, memeToken, positionId, tokenIndex, borrowAmount, borrowRate, liquidation)
}

// EmitBorrow is a paid mutator transaction binding the contract method 0xea34a577.
//
// Solidity: function emitBorrow(address borrower, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 borrowAmount, uint256 borrowRate, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitBorrow(borrower common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, tokenIndex uint8, borrowAmount *big.Int, borrowRate *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitBorrow(&_EventEmitter.TransactOpts, borrower, baseToken, memeToken, positionId, tokenIndex, borrowAmount, borrowRate, liquidation)
}

// EmitClaimFees is a paid mutator transaction binding the contract method 0x9c845792.
//
// Solidity: function emitClaimFees(address underlyingAsset, uint256 scaledUnclaimedFee, uint256 liquidityIndex, uint256 unclaimedFee) returns()
func (_EventEmitter *EventEmitterTransactor) EmitClaimFees(opts *bind.TransactOpts, underlyingAsset common.Address, scaledUnclaimedFee *big.Int, liquidityIndex *big.Int, unclaimedFee *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitClaimFees", underlyingAsset, scaledUnclaimedFee, liquidityIndex, unclaimedFee)
}

// EmitClaimFees is a paid mutator transaction binding the contract method 0x9c845792.
//
// Solidity: function emitClaimFees(address underlyingAsset, uint256 scaledUnclaimedFee, uint256 liquidityIndex, uint256 unclaimedFee) returns()
func (_EventEmitter *EventEmitterSession) EmitClaimFees(underlyingAsset common.Address, scaledUnclaimedFee *big.Int, liquidityIndex *big.Int, unclaimedFee *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitClaimFees(&_EventEmitter.TransactOpts, underlyingAsset, scaledUnclaimedFee, liquidityIndex, unclaimedFee)
}

// EmitClaimFees is a paid mutator transaction binding the contract method 0x9c845792.
//
// Solidity: function emitClaimFees(address underlyingAsset, uint256 scaledUnclaimedFee, uint256 liquidityIndex, uint256 unclaimedFee) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitClaimFees(underlyingAsset common.Address, scaledUnclaimedFee *big.Int, liquidityIndex *big.Int, unclaimedFee *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitClaimFees(&_EventEmitter.TransactOpts, underlyingAsset, scaledUnclaimedFee, liquidityIndex, unclaimedFee)
}

// EmitClose is a paid mutator transaction binding the contract method 0x5a7a3776.
//
// Solidity: function emitClose(address account, uint256 positionId) returns()
func (_EventEmitter *EventEmitterTransactor) EmitClose(opts *bind.TransactOpts, account common.Address, positionId *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitClose", account, positionId)
}

// EmitClose is a paid mutator transaction binding the contract method 0x5a7a3776.
//
// Solidity: function emitClose(address account, uint256 positionId) returns()
func (_EventEmitter *EventEmitterSession) EmitClose(account common.Address, positionId *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitClose(&_EventEmitter.TransactOpts, account, positionId)
}

// EmitClose is a paid mutator transaction binding the contract method 0x5a7a3776.
//
// Solidity: function emitClose(address account, uint256 positionId) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitClose(account common.Address, positionId *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitClose(&_EventEmitter.TransactOpts, account, positionId)
}

// EmitDeposit is a paid mutator transaction binding the contract method 0x119c6c83.
//
// Solidity: function emitDeposit(address depositor, address baseToken, address memeToken, uint256 positionId, uint256 depositAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterTransactor) EmitDeposit(opts *bind.TransactOpts, depositor common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, depositAmount *big.Int, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitDeposit", depositor, baseToken, memeToken, positionId, depositAmount, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitDeposit is a paid mutator transaction binding the contract method 0x119c6c83.
//
// Solidity: function emitDeposit(address depositor, address baseToken, address memeToken, uint256 positionId, uint256 depositAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterSession) EmitDeposit(depositor common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, depositAmount *big.Int, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitDeposit(&_EventEmitter.TransactOpts, depositor, baseToken, memeToken, positionId, depositAmount, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitDeposit is a paid mutator transaction binding the contract method 0x119c6c83.
//
// Solidity: function emitDeposit(address depositor, address baseToken, address memeToken, uint256 positionId, uint256 depositAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitDeposit(depositor common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, depositAmount *big.Int, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitDeposit(&_EventEmitter.TransactOpts, depositor, baseToken, memeToken, positionId, depositAmount, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitLiquidation is a paid mutator transaction binding the contract method 0x42ff64a5.
//
// Solidity: function emitLiquidation(address liquidator, address account, uint256 positionId, uint256 marginLevel, uint256 marginLevelLiquidationThreshold, uint256 totalCollateralUsd, uint256 totalDebtUsd, uint256 memePrice) returns()
func (_EventEmitter *EventEmitterTransactor) EmitLiquidation(opts *bind.TransactOpts, liquidator common.Address, account common.Address, positionId *big.Int, marginLevel *big.Int, marginLevelLiquidationThreshold *big.Int, totalCollateralUsd *big.Int, totalDebtUsd *big.Int, memePrice *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitLiquidation", liquidator, account, positionId, marginLevel, marginLevelLiquidationThreshold, totalCollateralUsd, totalDebtUsd, memePrice)
}

// EmitLiquidation is a paid mutator transaction binding the contract method 0x42ff64a5.
//
// Solidity: function emitLiquidation(address liquidator, address account, uint256 positionId, uint256 marginLevel, uint256 marginLevelLiquidationThreshold, uint256 totalCollateralUsd, uint256 totalDebtUsd, uint256 memePrice) returns()
func (_EventEmitter *EventEmitterSession) EmitLiquidation(liquidator common.Address, account common.Address, positionId *big.Int, marginLevel *big.Int, marginLevelLiquidationThreshold *big.Int, totalCollateralUsd *big.Int, totalDebtUsd *big.Int, memePrice *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitLiquidation(&_EventEmitter.TransactOpts, liquidator, account, positionId, marginLevel, marginLevelLiquidationThreshold, totalCollateralUsd, totalDebtUsd, memePrice)
}

// EmitLiquidation is a paid mutator transaction binding the contract method 0x42ff64a5.
//
// Solidity: function emitLiquidation(address liquidator, address account, uint256 positionId, uint256 marginLevel, uint256 marginLevelLiquidationThreshold, uint256 totalCollateralUsd, uint256 totalDebtUsd, uint256 memePrice) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitLiquidation(liquidator common.Address, account common.Address, positionId *big.Int, marginLevel *big.Int, marginLevelLiquidationThreshold *big.Int, totalCollateralUsd *big.Int, totalDebtUsd *big.Int, memePrice *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitLiquidation(&_EventEmitter.TransactOpts, liquidator, account, positionId, marginLevel, marginLevelLiquidationThreshold, totalCollateralUsd, totalDebtUsd, memePrice)
}

// EmitPoolCreated is a paid mutator transaction binding the contract method 0x55ac84ba.
//
// Solidity: function emitPoolCreated(address baseToken, address memeToken, address source, uint256 createdTimestamp, uint256 baseDecimals, uint256 memeDecimals) returns()
func (_EventEmitter *EventEmitterTransactor) EmitPoolCreated(opts *bind.TransactOpts, baseToken common.Address, memeToken common.Address, source common.Address, createdTimestamp *big.Int, baseDecimals *big.Int, memeDecimals *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitPoolCreated", baseToken, memeToken, source, createdTimestamp, baseDecimals, memeDecimals)
}

// EmitPoolCreated is a paid mutator transaction binding the contract method 0x55ac84ba.
//
// Solidity: function emitPoolCreated(address baseToken, address memeToken, address source, uint256 createdTimestamp, uint256 baseDecimals, uint256 memeDecimals) returns()
func (_EventEmitter *EventEmitterSession) EmitPoolCreated(baseToken common.Address, memeToken common.Address, source common.Address, createdTimestamp *big.Int, baseDecimals *big.Int, memeDecimals *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitPoolCreated(&_EventEmitter.TransactOpts, baseToken, memeToken, source, createdTimestamp, baseDecimals, memeDecimals)
}

// EmitPoolCreated is a paid mutator transaction binding the contract method 0x55ac84ba.
//
// Solidity: function emitPoolCreated(address baseToken, address memeToken, address source, uint256 createdTimestamp, uint256 baseDecimals, uint256 memeDecimals) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitPoolCreated(baseToken common.Address, memeToken common.Address, source common.Address, createdTimestamp *big.Int, baseDecimals *big.Int, memeDecimals *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitPoolCreated(&_EventEmitter.TransactOpts, baseToken, memeToken, source, createdTimestamp, baseDecimals, memeDecimals)
}

// EmitPoolUpdated is a paid mutator transaction binding the contract method 0x7c24dab7.
//
// Solidity: function emitPoolUpdated(address underlyingAsset, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex) returns()
func (_EventEmitter *EventEmitterTransactor) EmitPoolUpdated(opts *bind.TransactOpts, underlyingAsset common.Address, liquidityRate *big.Int, borrowRate *big.Int, liquidityIndex *big.Int, borrowIndex *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitPoolUpdated", underlyingAsset, liquidityRate, borrowRate, liquidityIndex, borrowIndex)
}

// EmitPoolUpdated is a paid mutator transaction binding the contract method 0x7c24dab7.
//
// Solidity: function emitPoolUpdated(address underlyingAsset, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex) returns()
func (_EventEmitter *EventEmitterSession) EmitPoolUpdated(underlyingAsset common.Address, liquidityRate *big.Int, borrowRate *big.Int, liquidityIndex *big.Int, borrowIndex *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitPoolUpdated(&_EventEmitter.TransactOpts, underlyingAsset, liquidityRate, borrowRate, liquidityIndex, borrowIndex)
}

// EmitPoolUpdated is a paid mutator transaction binding the contract method 0x7c24dab7.
//
// Solidity: function emitPoolUpdated(address underlyingAsset, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitPoolUpdated(underlyingAsset common.Address, liquidityRate *big.Int, borrowRate *big.Int, liquidityIndex *big.Int, borrowIndex *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitPoolUpdated(&_EventEmitter.TransactOpts, underlyingAsset, liquidityRate, borrowRate, liquidityIndex, borrowIndex)
}

// EmitPosition is a paid mutator transaction binding the contract method 0x09cd7ba2.
//
// Solidity: function emitPosition(address account, uint256 actionType, address baseToken, address memeToken, uint256 positionId, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterTransactor) EmitPosition(opts *bind.TransactOpts, account common.Address, actionType *big.Int, baseToken common.Address, memeToken common.Address, positionId *big.Int, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitPosition", account, actionType, baseToken, memeToken, positionId, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitPosition is a paid mutator transaction binding the contract method 0x09cd7ba2.
//
// Solidity: function emitPosition(address account, uint256 actionType, address baseToken, address memeToken, uint256 positionId, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterSession) EmitPosition(account common.Address, actionType *big.Int, baseToken common.Address, memeToken common.Address, positionId *big.Int, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitPosition(&_EventEmitter.TransactOpts, account, actionType, baseToken, memeToken, positionId, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitPosition is a paid mutator transaction binding the contract method 0x09cd7ba2.
//
// Solidity: function emitPosition(address account, uint256 actionType, address baseToken, address memeToken, uint256 positionId, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitPosition(account common.Address, actionType *big.Int, baseToken common.Address, memeToken common.Address, positionId *big.Int, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitPosition(&_EventEmitter.TransactOpts, account, actionType, baseToken, memeToken, positionId, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitRemove is a paid mutator transaction binding the contract method 0x292ae722.
//
// Solidity: function emitRemove(address remover, address baseToken, address memeToken, uint256 liquidity, address to, uint256 baseAmount, uint256 memeAmount) returns()
func (_EventEmitter *EventEmitterTransactor) EmitRemove(opts *bind.TransactOpts, remover common.Address, baseToken common.Address, memeToken common.Address, liquidity *big.Int, to common.Address, baseAmount *big.Int, memeAmount *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitRemove", remover, baseToken, memeToken, liquidity, to, baseAmount, memeAmount)
}

// EmitRemove is a paid mutator transaction binding the contract method 0x292ae722.
//
// Solidity: function emitRemove(address remover, address baseToken, address memeToken, uint256 liquidity, address to, uint256 baseAmount, uint256 memeAmount) returns()
func (_EventEmitter *EventEmitterSession) EmitRemove(remover common.Address, baseToken common.Address, memeToken common.Address, liquidity *big.Int, to common.Address, baseAmount *big.Int, memeAmount *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitRemove(&_EventEmitter.TransactOpts, remover, baseToken, memeToken, liquidity, to, baseAmount, memeAmount)
}

// EmitRemove is a paid mutator transaction binding the contract method 0x292ae722.
//
// Solidity: function emitRemove(address remover, address baseToken, address memeToken, uint256 liquidity, address to, uint256 baseAmount, uint256 memeAmount) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitRemove(remover common.Address, baseToken common.Address, memeToken common.Address, liquidity *big.Int, to common.Address, baseAmount *big.Int, memeAmount *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitRemove(&_EventEmitter.TransactOpts, remover, baseToken, memeToken, liquidity, to, baseAmount, memeAmount)
}

// EmitRepay is a paid mutator transaction binding the contract method 0x8262009e.
//
// Solidity: function emitRepay(address repayer, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 repayAmount, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterTransactor) EmitRepay(opts *bind.TransactOpts, repayer common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, tokenIndex uint8, repayAmount *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitRepay", repayer, baseToken, memeToken, positionId, tokenIndex, repayAmount, liquidation)
}

// EmitRepay is a paid mutator transaction binding the contract method 0x8262009e.
//
// Solidity: function emitRepay(address repayer, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 repayAmount, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterSession) EmitRepay(repayer common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, tokenIndex uint8, repayAmount *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitRepay(&_EventEmitter.TransactOpts, repayer, baseToken, memeToken, positionId, tokenIndex, repayAmount, liquidation)
}

// EmitRepay is a paid mutator transaction binding the contract method 0x8262009e.
//
// Solidity: function emitRepay(address repayer, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 repayAmount, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitRepay(repayer common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, tokenIndex uint8, repayAmount *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitRepay(&_EventEmitter.TransactOpts, repayer, baseToken, memeToken, positionId, tokenIndex, repayAmount, liquidation)
}

// EmitSwap is a paid mutator transaction binding the contract method 0x11ccb21d.
//
// Solidity: function emitSwap(address account, address tokenIn, address tokenOut, uint256 positionId, uint256 amountIn, uint256 amountOut, uint256 fee, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterTransactor) EmitSwap(opts *bind.TransactOpts, account common.Address, tokenIn common.Address, tokenOut common.Address, positionId *big.Int, amountIn *big.Int, amountOut *big.Int, fee *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitSwap", account, tokenIn, tokenOut, positionId, amountIn, amountOut, fee, liquidation)
}

// EmitSwap is a paid mutator transaction binding the contract method 0x11ccb21d.
//
// Solidity: function emitSwap(address account, address tokenIn, address tokenOut, uint256 positionId, uint256 amountIn, uint256 amountOut, uint256 fee, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterSession) EmitSwap(account common.Address, tokenIn common.Address, tokenOut common.Address, positionId *big.Int, amountIn *big.Int, amountOut *big.Int, fee *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitSwap(&_EventEmitter.TransactOpts, account, tokenIn, tokenOut, positionId, amountIn, amountOut, fee, liquidation)
}

// EmitSwap is a paid mutator transaction binding the contract method 0x11ccb21d.
//
// Solidity: function emitSwap(address account, address tokenIn, address tokenOut, uint256 positionId, uint256 amountIn, uint256 amountOut, uint256 fee, (uint256,uint256,uint256,uint256) liquidation) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitSwap(account common.Address, tokenIn common.Address, tokenOut common.Address, positionId *big.Int, amountIn *big.Int, amountOut *big.Int, fee *big.Int, liquidation EventLiquidation) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitSwap(&_EventEmitter.TransactOpts, account, tokenIn, tokenOut, positionId, amountIn, amountOut, fee, liquidation)
}

// EmitWithdraw is a paid mutator transaction binding the contract method 0x82fcd8ca.
//
// Solidity: function emitWithdraw(address withdrawer, address baseToken, address memeToken, uint256 positionId, uint256 withdrawAmount, address to, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterTransactor) EmitWithdraw(opts *bind.TransactOpts, withdrawer common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, withdrawAmount *big.Int, to common.Address, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.contract.Transact(opts, "emitWithdraw", withdrawer, baseToken, memeToken, positionId, withdrawAmount, to, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitWithdraw is a paid mutator transaction binding the contract method 0x82fcd8ca.
//
// Solidity: function emitWithdraw(address withdrawer, address baseToken, address memeToken, uint256 positionId, uint256 withdrawAmount, address to, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterSession) EmitWithdraw(withdrawer common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, withdrawAmount *big.Int, to common.Address, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitWithdraw(&_EventEmitter.TransactOpts, withdrawer, baseToken, memeToken, positionId, withdrawAmount, to, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EmitWithdraw is a paid mutator transaction binding the contract method 0x82fcd8ca.
//
// Solidity: function emitWithdraw(address withdrawer, address baseToken, address memeToken, uint256 positionId, uint256 withdrawAmount, address to, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled) returns()
func (_EventEmitter *EventEmitterTransactorSession) EmitWithdraw(withdrawer common.Address, baseToken common.Address, memeToken common.Address, positionId *big.Int, withdrawAmount *big.Int, to common.Address, baseCollateral *big.Int, baseDebtScaled *big.Int, memeCollateral *big.Int, memeDebtScaled *big.Int) (*types.Transaction, error) {
	return _EventEmitter.Contract.EmitWithdraw(&_EventEmitter.TransactOpts, withdrawer, baseToken, memeToken, positionId, withdrawAmount, to, baseCollateral, baseDebtScaled, memeCollateral, memeDebtScaled)
}

// EventEmitterAddIterator is returned from FilterAdd and is used to iterate over the raw logs and unpacked data for Add events raised by the EventEmitter contract.
type EventEmitterAddIterator struct {
	Event *EventEmitterAdd // Event containing the contract specifics and raw log

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
func (it *EventEmitterAddIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterAdd)
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
		it.Event = new(EventEmitterAdd)
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
func (it *EventEmitterAddIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterAdd represents a Add event raised by the EventEmitter contract.
type EventEmitterAdd struct {
	Adder      common.Address
	BaseToken  common.Address
	MemeToken  common.Address
	To         common.Address
	BaseAmount *big.Int
	MemeAmount *big.Int
	Liquidity  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAdd is a free log retrieval operation binding the contract event 0xfb437c87369d4cc94da483ca04941a9d959a63ff6d87fdbea74b8072c1ff9861.
//
// Solidity: event Add(address indexed adder, address baseToken, address memeToken, address to, uint256 baseAmount, uint256 memeAmount, uint256 liquidity)
func (_EventEmitter *EventEmitterFilterer) FilterAdd(opts *bind.FilterOpts, adder []common.Address) (*EventEmitterAddIterator, error) {

	var adderRule []interface{}
	for _, adderItem := range adder {
		adderRule = append(adderRule, adderItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Add", adderRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterAddIterator{contract: _EventEmitter.contract, event: "Add", logs: logs, sub: sub}, nil
}

// WatchAdd is a free log subscription operation binding the contract event 0xfb437c87369d4cc94da483ca04941a9d959a63ff6d87fdbea74b8072c1ff9861.
//
// Solidity: event Add(address indexed adder, address baseToken, address memeToken, address to, uint256 baseAmount, uint256 memeAmount, uint256 liquidity)
func (_EventEmitter *EventEmitterFilterer) WatchAdd(opts *bind.WatchOpts, sink chan<- *EventEmitterAdd, adder []common.Address) (event.Subscription, error) {

	var adderRule []interface{}
	for _, adderItem := range adder {
		adderRule = append(adderRule, adderItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Add", adderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterAdd)
				if err := _EventEmitter.contract.UnpackLog(event, "Add", log); err != nil {
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

// ParseAdd is a log parse operation binding the contract event 0xfb437c87369d4cc94da483ca04941a9d959a63ff6d87fdbea74b8072c1ff9861.
//
// Solidity: event Add(address indexed adder, address baseToken, address memeToken, address to, uint256 baseAmount, uint256 memeAmount, uint256 liquidity)
func (_EventEmitter *EventEmitterFilterer) ParseAdd(log types.Log) (*EventEmitterAdd, error) {
	event := new(EventEmitterAdd)
	if err := _EventEmitter.contract.UnpackLog(event, "Add", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterBorrowIterator is returned from FilterBorrow and is used to iterate over the raw logs and unpacked data for Borrow events raised by the EventEmitter contract.
type EventEmitterBorrowIterator struct {
	Event *EventEmitterBorrow // Event containing the contract specifics and raw log

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
func (it *EventEmitterBorrowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterBorrow)
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
		it.Event = new(EventEmitterBorrow)
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
func (it *EventEmitterBorrowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterBorrowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterBorrow represents a Borrow event raised by the EventEmitter contract.
type EventEmitterBorrow struct {
	Borrower       common.Address
	BaseToken      common.Address
	MemeToken      common.Address
	PositionId     *big.Int
	TokenIndex     uint8
	BorrowAmount   *big.Int
	BorrowRate     *big.Int
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBorrow is a free log retrieval operation binding the contract event 0x2ef1e8737d096826c9abef1201bb205ea380555780a54bd904cf67ca7dba8c5f.
//
// Solidity: event Borrow(address indexed borrower, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 borrowAmount, uint256 borrowRate, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) FilterBorrow(opts *bind.FilterOpts, borrower []common.Address) (*EventEmitterBorrowIterator, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Borrow", borrowerRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterBorrowIterator{contract: _EventEmitter.contract, event: "Borrow", logs: logs, sub: sub}, nil
}

// WatchBorrow is a free log subscription operation binding the contract event 0x2ef1e8737d096826c9abef1201bb205ea380555780a54bd904cf67ca7dba8c5f.
//
// Solidity: event Borrow(address indexed borrower, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 borrowAmount, uint256 borrowRate, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) WatchBorrow(opts *bind.WatchOpts, sink chan<- *EventEmitterBorrow, borrower []common.Address) (event.Subscription, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Borrow", borrowerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterBorrow)
				if err := _EventEmitter.contract.UnpackLog(event, "Borrow", log); err != nil {
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

// ParseBorrow is a log parse operation binding the contract event 0x2ef1e8737d096826c9abef1201bb205ea380555780a54bd904cf67ca7dba8c5f.
//
// Solidity: event Borrow(address indexed borrower, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 borrowAmount, uint256 borrowRate, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) ParseBorrow(log types.Log) (*EventEmitterBorrow, error) {
	event := new(EventEmitterBorrow)
	if err := _EventEmitter.contract.UnpackLog(event, "Borrow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterClaimFeesIterator is returned from FilterClaimFees and is used to iterate over the raw logs and unpacked data for ClaimFees events raised by the EventEmitter contract.
type EventEmitterClaimFeesIterator struct {
	Event *EventEmitterClaimFees // Event containing the contract specifics and raw log

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
func (it *EventEmitterClaimFeesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterClaimFees)
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
		it.Event = new(EventEmitterClaimFees)
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
func (it *EventEmitterClaimFeesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterClaimFeesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterClaimFees represents a ClaimFees event raised by the EventEmitter contract.
type EventEmitterClaimFees struct {
	Pool               common.Address
	ScaledUnclaimedFee *big.Int
	LiquidityIndex     *big.Int
	UnclaimedFee       *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterClaimFees is a free log retrieval operation binding the contract event 0xd2e6085315c6e1c1c7406a47c7d006a8c1f931396d868c16046dea71365ff031.
//
// Solidity: event ClaimFees(address indexed pool, uint256 scaledUnclaimedFee, uint256 liquidityIndex, uint256 unclaimedFee)
func (_EventEmitter *EventEmitterFilterer) FilterClaimFees(opts *bind.FilterOpts, pool []common.Address) (*EventEmitterClaimFeesIterator, error) {

	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "ClaimFees", poolRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterClaimFeesIterator{contract: _EventEmitter.contract, event: "ClaimFees", logs: logs, sub: sub}, nil
}

// WatchClaimFees is a free log subscription operation binding the contract event 0xd2e6085315c6e1c1c7406a47c7d006a8c1f931396d868c16046dea71365ff031.
//
// Solidity: event ClaimFees(address indexed pool, uint256 scaledUnclaimedFee, uint256 liquidityIndex, uint256 unclaimedFee)
func (_EventEmitter *EventEmitterFilterer) WatchClaimFees(opts *bind.WatchOpts, sink chan<- *EventEmitterClaimFees, pool []common.Address) (event.Subscription, error) {

	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "ClaimFees", poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterClaimFees)
				if err := _EventEmitter.contract.UnpackLog(event, "ClaimFees", log); err != nil {
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

// ParseClaimFees is a log parse operation binding the contract event 0xd2e6085315c6e1c1c7406a47c7d006a8c1f931396d868c16046dea71365ff031.
//
// Solidity: event ClaimFees(address indexed pool, uint256 scaledUnclaimedFee, uint256 liquidityIndex, uint256 unclaimedFee)
func (_EventEmitter *EventEmitterFilterer) ParseClaimFees(log types.Log) (*EventEmitterClaimFees, error) {
	event := new(EventEmitterClaimFees)
	if err := _EventEmitter.contract.UnpackLog(event, "ClaimFees", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterCloseIterator is returned from FilterClose and is used to iterate over the raw logs and unpacked data for Close events raised by the EventEmitter contract.
type EventEmitterCloseIterator struct {
	Event *EventEmitterClose // Event containing the contract specifics and raw log

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
func (it *EventEmitterCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterClose)
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
		it.Event = new(EventEmitterClose)
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
func (it *EventEmitterCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterClose represents a Close event raised by the EventEmitter contract.
type EventEmitterClose struct {
	Account    common.Address
	PositionId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClose is a free log retrieval operation binding the contract event 0x684222b0069d4a2e5e0d986611cc5182d543904c4e4264bf770d4e51faefc822.
//
// Solidity: event Close(address indexed account, uint256 positionId)
func (_EventEmitter *EventEmitterFilterer) FilterClose(opts *bind.FilterOpts, account []common.Address) (*EventEmitterCloseIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Close", accountRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterCloseIterator{contract: _EventEmitter.contract, event: "Close", logs: logs, sub: sub}, nil
}

// WatchClose is a free log subscription operation binding the contract event 0x684222b0069d4a2e5e0d986611cc5182d543904c4e4264bf770d4e51faefc822.
//
// Solidity: event Close(address indexed account, uint256 positionId)
func (_EventEmitter *EventEmitterFilterer) WatchClose(opts *bind.WatchOpts, sink chan<- *EventEmitterClose, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Close", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterClose)
				if err := _EventEmitter.contract.UnpackLog(event, "Close", log); err != nil {
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

// ParseClose is a log parse operation binding the contract event 0x684222b0069d4a2e5e0d986611cc5182d543904c4e4264bf770d4e51faefc822.
//
// Solidity: event Close(address indexed account, uint256 positionId)
func (_EventEmitter *EventEmitterFilterer) ParseClose(log types.Log) (*EventEmitterClose, error) {
	event := new(EventEmitterClose)
	if err := _EventEmitter.contract.UnpackLog(event, "Close", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the EventEmitter contract.
type EventEmitterDepositIterator struct {
	Event *EventEmitterDeposit // Event containing the contract specifics and raw log

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
func (it *EventEmitterDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterDeposit)
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
		it.Event = new(EventEmitterDeposit)
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
func (it *EventEmitterDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterDeposit represents a Deposit event raised by the EventEmitter contract.
type EventEmitterDeposit struct {
	Depositor      common.Address
	BaseToken      common.Address
	MemeToken      common.Address
	PositionId     *big.Int
	DepositAmount  *big.Int
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x8a28c8f95aa6ea7a5133fc0d3b124f64fecc7c3c53414fca4db4c02cc53e2ad6.
//
// Solidity: event Deposit(address indexed depositor, address baseToken, address memeToken, uint256 positionId, uint256 depositAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) FilterDeposit(opts *bind.FilterOpts, depositor []common.Address) (*EventEmitterDepositIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Deposit", depositorRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterDepositIterator{contract: _EventEmitter.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x8a28c8f95aa6ea7a5133fc0d3b124f64fecc7c3c53414fca4db4c02cc53e2ad6.
//
// Solidity: event Deposit(address indexed depositor, address baseToken, address memeToken, uint256 positionId, uint256 depositAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *EventEmitterDeposit, depositor []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Deposit", depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterDeposit)
				if err := _EventEmitter.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x8a28c8f95aa6ea7a5133fc0d3b124f64fecc7c3c53414fca4db4c02cc53e2ad6.
//
// Solidity: event Deposit(address indexed depositor, address baseToken, address memeToken, uint256 positionId, uint256 depositAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) ParseDeposit(log types.Log) (*EventEmitterDeposit, error) {
	event := new(EventEmitterDeposit)
	if err := _EventEmitter.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterLiquidationIterator is returned from FilterLiquidation and is used to iterate over the raw logs and unpacked data for Liquidation events raised by the EventEmitter contract.
type EventEmitterLiquidationIterator struct {
	Event *EventEmitterLiquidation // Event containing the contract specifics and raw log

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
func (it *EventEmitterLiquidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterLiquidation)
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
		it.Event = new(EventEmitterLiquidation)
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
func (it *EventEmitterLiquidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterLiquidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterLiquidation represents a Liquidation event raised by the EventEmitter contract.
type EventEmitterLiquidation struct {
	Liquidator                      common.Address
	Account                         common.Address
	PositionId                      *big.Int
	MarginLevel                     *big.Int
	MarginLevelLiquidationThreshold *big.Int
	TotalCollateralUsd              *big.Int
	TotalDebtUsd                    *big.Int
	MemePrice                       *big.Int
	Raw                             types.Log // Blockchain specific contextual infos
}

// FilterLiquidation is a free log retrieval operation binding the contract event 0x0211bad4030ad74b54fadd2bfc0ca77ac2478c90fe05d0b58357e7ce701ad2d4.
//
// Solidity: event Liquidation(address indexed liquidator, address indexed account, uint256 positionId, uint256 marginLevel, uint256 marginLevelLiquidationThreshold, uint256 totalCollateralUsd, uint256 totalDebtUsd, uint256 memePrice)
func (_EventEmitter *EventEmitterFilterer) FilterLiquidation(opts *bind.FilterOpts, liquidator []common.Address, account []common.Address) (*EventEmitterLiquidationIterator, error) {

	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Liquidation", liquidatorRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterLiquidationIterator{contract: _EventEmitter.contract, event: "Liquidation", logs: logs, sub: sub}, nil
}

// WatchLiquidation is a free log subscription operation binding the contract event 0x0211bad4030ad74b54fadd2bfc0ca77ac2478c90fe05d0b58357e7ce701ad2d4.
//
// Solidity: event Liquidation(address indexed liquidator, address indexed account, uint256 positionId, uint256 marginLevel, uint256 marginLevelLiquidationThreshold, uint256 totalCollateralUsd, uint256 totalDebtUsd, uint256 memePrice)
func (_EventEmitter *EventEmitterFilterer) WatchLiquidation(opts *bind.WatchOpts, sink chan<- *EventEmitterLiquidation, liquidator []common.Address, account []common.Address) (event.Subscription, error) {

	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Liquidation", liquidatorRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterLiquidation)
				if err := _EventEmitter.contract.UnpackLog(event, "Liquidation", log); err != nil {
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

// ParseLiquidation is a log parse operation binding the contract event 0x0211bad4030ad74b54fadd2bfc0ca77ac2478c90fe05d0b58357e7ce701ad2d4.
//
// Solidity: event Liquidation(address indexed liquidator, address indexed account, uint256 positionId, uint256 marginLevel, uint256 marginLevelLiquidationThreshold, uint256 totalCollateralUsd, uint256 totalDebtUsd, uint256 memePrice)
func (_EventEmitter *EventEmitterFilterer) ParseLiquidation(log types.Log) (*EventEmitterLiquidation, error) {
	event := new(EventEmitterLiquidation)
	if err := _EventEmitter.contract.UnpackLog(event, "Liquidation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterPoolCreatedIterator is returned from FilterPoolCreated and is used to iterate over the raw logs and unpacked data for PoolCreated events raised by the EventEmitter contract.
type EventEmitterPoolCreatedIterator struct {
	Event *EventEmitterPoolCreated // Event containing the contract specifics and raw log

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
func (it *EventEmitterPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterPoolCreated)
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
		it.Event = new(EventEmitterPoolCreated)
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
func (it *EventEmitterPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterPoolCreated represents a PoolCreated event raised by the EventEmitter contract.
type EventEmitterPoolCreated struct {
	BaseToken        common.Address
	MemeToken        common.Address
	Source           common.Address
	CreatedTimestamp *big.Int
	BaseDecimals     *big.Int
	MemeDecimals     *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPoolCreated is a free log retrieval operation binding the contract event 0xa8c38247fd3f092e3bc806fb7dff597b939f4ec6ffad1514a35eb2776e2a61b8.
//
// Solidity: event PoolCreated(address baseToken, address memeToken, address source, uint256 createdTimestamp, uint256 baseDecimals, uint256 memeDecimals)
func (_EventEmitter *EventEmitterFilterer) FilterPoolCreated(opts *bind.FilterOpts) (*EventEmitterPoolCreatedIterator, error) {

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "PoolCreated")
	if err != nil {
		return nil, err
	}
	return &EventEmitterPoolCreatedIterator{contract: _EventEmitter.contract, event: "PoolCreated", logs: logs, sub: sub}, nil
}

// WatchPoolCreated is a free log subscription operation binding the contract event 0xa8c38247fd3f092e3bc806fb7dff597b939f4ec6ffad1514a35eb2776e2a61b8.
//
// Solidity: event PoolCreated(address baseToken, address memeToken, address source, uint256 createdTimestamp, uint256 baseDecimals, uint256 memeDecimals)
func (_EventEmitter *EventEmitterFilterer) WatchPoolCreated(opts *bind.WatchOpts, sink chan<- *EventEmitterPoolCreated) (event.Subscription, error) {

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "PoolCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterPoolCreated)
				if err := _EventEmitter.contract.UnpackLog(event, "PoolCreated", log); err != nil {
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

// ParsePoolCreated is a log parse operation binding the contract event 0xa8c38247fd3f092e3bc806fb7dff597b939f4ec6ffad1514a35eb2776e2a61b8.
//
// Solidity: event PoolCreated(address baseToken, address memeToken, address source, uint256 createdTimestamp, uint256 baseDecimals, uint256 memeDecimals)
func (_EventEmitter *EventEmitterFilterer) ParsePoolCreated(log types.Log) (*EventEmitterPoolCreated, error) {
	event := new(EventEmitterPoolCreated)
	if err := _EventEmitter.contract.UnpackLog(event, "PoolCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterPoolUpdatedIterator is returned from FilterPoolUpdated and is used to iterate over the raw logs and unpacked data for PoolUpdated events raised by the EventEmitter contract.
type EventEmitterPoolUpdatedIterator struct {
	Event *EventEmitterPoolUpdated // Event containing the contract specifics and raw log

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
func (it *EventEmitterPoolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterPoolUpdated)
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
		it.Event = new(EventEmitterPoolUpdated)
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
func (it *EventEmitterPoolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterPoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterPoolUpdated represents a PoolUpdated event raised by the EventEmitter contract.
type EventEmitterPoolUpdated struct {
	Pool           common.Address
	LiquidityRate  *big.Int
	BorrowRate     *big.Int
	LiquidityIndex *big.Int
	BorrowIndex    *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPoolUpdated is a free log retrieval operation binding the contract event 0xc320a8529b15b851aaa68519919ac344e1caceaf4f47d15df6e58181dfec6319.
//
// Solidity: event PoolUpdated(address indexed pool, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex)
func (_EventEmitter *EventEmitterFilterer) FilterPoolUpdated(opts *bind.FilterOpts, pool []common.Address) (*EventEmitterPoolUpdatedIterator, error) {

	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "PoolUpdated", poolRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterPoolUpdatedIterator{contract: _EventEmitter.contract, event: "PoolUpdated", logs: logs, sub: sub}, nil
}

// WatchPoolUpdated is a free log subscription operation binding the contract event 0xc320a8529b15b851aaa68519919ac344e1caceaf4f47d15df6e58181dfec6319.
//
// Solidity: event PoolUpdated(address indexed pool, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex)
func (_EventEmitter *EventEmitterFilterer) WatchPoolUpdated(opts *bind.WatchOpts, sink chan<- *EventEmitterPoolUpdated, pool []common.Address) (event.Subscription, error) {

	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "PoolUpdated", poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterPoolUpdated)
				if err := _EventEmitter.contract.UnpackLog(event, "PoolUpdated", log); err != nil {
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

// ParsePoolUpdated is a log parse operation binding the contract event 0xc320a8529b15b851aaa68519919ac344e1caceaf4f47d15df6e58181dfec6319.
//
// Solidity: event PoolUpdated(address indexed pool, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex)
func (_EventEmitter *EventEmitterFilterer) ParsePoolUpdated(log types.Log) (*EventEmitterPoolUpdated, error) {
	event := new(EventEmitterPoolUpdated)
	if err := _EventEmitter.contract.UnpackLog(event, "PoolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterPositionIterator is returned from FilterPosition and is used to iterate over the raw logs and unpacked data for Position events raised by the EventEmitter contract.
type EventEmitterPositionIterator struct {
	Event *EventEmitterPosition // Event containing the contract specifics and raw log

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
func (it *EventEmitterPositionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterPosition)
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
		it.Event = new(EventEmitterPosition)
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
func (it *EventEmitterPositionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterPositionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterPosition represents a Position event raised by the EventEmitter contract.
type EventEmitterPosition struct {
	Account        common.Address
	ActionType     *big.Int
	BaseToken      common.Address
	MemeToken      common.Address
	PositionId     *big.Int
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPosition is a free log retrieval operation binding the contract event 0xe72b15b63ab2516a28bddccae211fb56af88f7072bcebe4b560beb9189a6492d.
//
// Solidity: event Position(address indexed account, uint256 actionType, address baseToken, address memeToken, uint256 positionId, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) FilterPosition(opts *bind.FilterOpts, account []common.Address) (*EventEmitterPositionIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Position", accountRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterPositionIterator{contract: _EventEmitter.contract, event: "Position", logs: logs, sub: sub}, nil
}

// WatchPosition is a free log subscription operation binding the contract event 0xe72b15b63ab2516a28bddccae211fb56af88f7072bcebe4b560beb9189a6492d.
//
// Solidity: event Position(address indexed account, uint256 actionType, address baseToken, address memeToken, uint256 positionId, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) WatchPosition(opts *bind.WatchOpts, sink chan<- *EventEmitterPosition, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Position", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterPosition)
				if err := _EventEmitter.contract.UnpackLog(event, "Position", log); err != nil {
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

// ParsePosition is a log parse operation binding the contract event 0xe72b15b63ab2516a28bddccae211fb56af88f7072bcebe4b560beb9189a6492d.
//
// Solidity: event Position(address indexed account, uint256 actionType, address baseToken, address memeToken, uint256 positionId, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) ParsePosition(log types.Log) (*EventEmitterPosition, error) {
	event := new(EventEmitterPosition)
	if err := _EventEmitter.contract.UnpackLog(event, "Position", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterRemoveIterator is returned from FilterRemove and is used to iterate over the raw logs and unpacked data for Remove events raised by the EventEmitter contract.
type EventEmitterRemoveIterator struct {
	Event *EventEmitterRemove // Event containing the contract specifics and raw log

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
func (it *EventEmitterRemoveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterRemove)
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
		it.Event = new(EventEmitterRemove)
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
func (it *EventEmitterRemoveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterRemove represents a Remove event raised by the EventEmitter contract.
type EventEmitterRemove struct {
	Remover    common.Address
	BaseToken  common.Address
	MemeToken  common.Address
	Liquidity  *big.Int
	To         common.Address
	BaseAmount *big.Int
	MemeAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemove is a free log retrieval operation binding the contract event 0xe10a339dd5329df14a8ec13eb4115b75ab032fd40e2ff2594a4a5bf80e497a41.
//
// Solidity: event Remove(address indexed remover, address baseToken, address memeToken, uint256 liquidity, address to, uint256 baseAmount, uint256 memeAmount)
func (_EventEmitter *EventEmitterFilterer) FilterRemove(opts *bind.FilterOpts, remover []common.Address) (*EventEmitterRemoveIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Remove", removerRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterRemoveIterator{contract: _EventEmitter.contract, event: "Remove", logs: logs, sub: sub}, nil
}

// WatchRemove is a free log subscription operation binding the contract event 0xe10a339dd5329df14a8ec13eb4115b75ab032fd40e2ff2594a4a5bf80e497a41.
//
// Solidity: event Remove(address indexed remover, address baseToken, address memeToken, uint256 liquidity, address to, uint256 baseAmount, uint256 memeAmount)
func (_EventEmitter *EventEmitterFilterer) WatchRemove(opts *bind.WatchOpts, sink chan<- *EventEmitterRemove, remover []common.Address) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Remove", removerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterRemove)
				if err := _EventEmitter.contract.UnpackLog(event, "Remove", log); err != nil {
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

// ParseRemove is a log parse operation binding the contract event 0xe10a339dd5329df14a8ec13eb4115b75ab032fd40e2ff2594a4a5bf80e497a41.
//
// Solidity: event Remove(address indexed remover, address baseToken, address memeToken, uint256 liquidity, address to, uint256 baseAmount, uint256 memeAmount)
func (_EventEmitter *EventEmitterFilterer) ParseRemove(log types.Log) (*EventEmitterRemove, error) {
	event := new(EventEmitterRemove)
	if err := _EventEmitter.contract.UnpackLog(event, "Remove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterRepayIterator is returned from FilterRepay and is used to iterate over the raw logs and unpacked data for Repay events raised by the EventEmitter contract.
type EventEmitterRepayIterator struct {
	Event *EventEmitterRepay // Event containing the contract specifics and raw log

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
func (it *EventEmitterRepayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterRepay)
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
		it.Event = new(EventEmitterRepay)
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
func (it *EventEmitterRepayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterRepayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterRepay represents a Repay event raised by the EventEmitter contract.
type EventEmitterRepay struct {
	Repayer        common.Address
	BaseToken      common.Address
	MemeToken      common.Address
	PositionId     *big.Int
	TokenIndex     uint8
	RepayAmount    *big.Int
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRepay is a free log retrieval operation binding the contract event 0x4fc3cf57f1c587a9f0808812dd72668a2b82f54411b0737d06c4043ff044cf9a.
//
// Solidity: event Repay(address indexed repayer, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 repayAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) FilterRepay(opts *bind.FilterOpts, repayer []common.Address) (*EventEmitterRepayIterator, error) {

	var repayerRule []interface{}
	for _, repayerItem := range repayer {
		repayerRule = append(repayerRule, repayerItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Repay", repayerRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterRepayIterator{contract: _EventEmitter.contract, event: "Repay", logs: logs, sub: sub}, nil
}

// WatchRepay is a free log subscription operation binding the contract event 0x4fc3cf57f1c587a9f0808812dd72668a2b82f54411b0737d06c4043ff044cf9a.
//
// Solidity: event Repay(address indexed repayer, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 repayAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) WatchRepay(opts *bind.WatchOpts, sink chan<- *EventEmitterRepay, repayer []common.Address) (event.Subscription, error) {

	var repayerRule []interface{}
	for _, repayerItem := range repayer {
		repayerRule = append(repayerRule, repayerItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Repay", repayerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterRepay)
				if err := _EventEmitter.contract.UnpackLog(event, "Repay", log); err != nil {
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

// ParseRepay is a log parse operation binding the contract event 0x4fc3cf57f1c587a9f0808812dd72668a2b82f54411b0737d06c4043ff044cf9a.
//
// Solidity: event Repay(address indexed repayer, address baseToken, address memeToken, uint256 positionId, uint8 tokenIndex, uint256 repayAmount, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) ParseRepay(log types.Log) (*EventEmitterRepay, error) {
	event := new(EventEmitterRepay)
	if err := _EventEmitter.contract.UnpackLog(event, "Repay", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the EventEmitter contract.
type EventEmitterSwapIterator struct {
	Event *EventEmitterSwap // Event containing the contract specifics and raw log

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
func (it *EventEmitterSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterSwap)
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
		it.Event = new(EventEmitterSwap)
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
func (it *EventEmitterSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterSwap represents a Swap event raised by the EventEmitter contract.
type EventEmitterSwap struct {
	Account        common.Address
	TokenIn        common.Address
	TokenOut       common.Address
	PositionId     *big.Int
	AmountIn       *big.Int
	AmountOut      *big.Int
	Fee            *big.Int
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0x2fe197cd2fd0cdf1ee4a888ae1fdf65897620942dd65441d363cd479527a8568.
//
// Solidity: event Swap(address indexed account, address tokenIn, address tokenOut, uint256 positionId, uint256 amountIn, uint256 amountOut, uint256 fee, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) FilterSwap(opts *bind.FilterOpts, account []common.Address) (*EventEmitterSwapIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Swap", accountRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterSwapIterator{contract: _EventEmitter.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0x2fe197cd2fd0cdf1ee4a888ae1fdf65897620942dd65441d363cd479527a8568.
//
// Solidity: event Swap(address indexed account, address tokenIn, address tokenOut, uint256 positionId, uint256 amountIn, uint256 amountOut, uint256 fee, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *EventEmitterSwap, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Swap", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterSwap)
				if err := _EventEmitter.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0x2fe197cd2fd0cdf1ee4a888ae1fdf65897620942dd65441d363cd479527a8568.
//
// Solidity: event Swap(address indexed account, address tokenIn, address tokenOut, uint256 positionId, uint256 amountIn, uint256 amountOut, uint256 fee, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) ParseSwap(log types.Log) (*EventEmitterSwap, error) {
	event := new(EventEmitterSwap)
	if err := _EventEmitter.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventEmitterWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the EventEmitter contract.
type EventEmitterWithdrawIterator struct {
	Event *EventEmitterWithdraw // Event containing the contract specifics and raw log

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
func (it *EventEmitterWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventEmitterWithdraw)
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
		it.Event = new(EventEmitterWithdraw)
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
func (it *EventEmitterWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventEmitterWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventEmitterWithdraw represents a Withdraw event raised by the EventEmitter contract.
type EventEmitterWithdraw struct {
	Withdrawer     common.Address
	BaseToken      common.Address
	MemeToken      common.Address
	PositionId     *big.Int
	WithdrawAmount *big.Int
	To             common.Address
	BaseCollateral *big.Int
	BaseDebtScaled *big.Int
	MemeCollateral *big.Int
	MemeDebtScaled *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x8be7f12253faa74082f51d8d13cc0e76176d2ab1442fad134c53068eb958374c.
//
// Solidity: event Withdraw(address indexed withdrawer, address baseToken, address memeToken, uint256 positionId, uint256 withdrawAmount, address to, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) FilterWithdraw(opts *bind.FilterOpts, withdrawer []common.Address) (*EventEmitterWithdrawIterator, error) {

	var withdrawerRule []interface{}
	for _, withdrawerItem := range withdrawer {
		withdrawerRule = append(withdrawerRule, withdrawerItem)
	}

	logs, sub, err := _EventEmitter.contract.FilterLogs(opts, "Withdraw", withdrawerRule)
	if err != nil {
		return nil, err
	}
	return &EventEmitterWithdrawIterator{contract: _EventEmitter.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x8be7f12253faa74082f51d8d13cc0e76176d2ab1442fad134c53068eb958374c.
//
// Solidity: event Withdraw(address indexed withdrawer, address baseToken, address memeToken, uint256 positionId, uint256 withdrawAmount, address to, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *EventEmitterWithdraw, withdrawer []common.Address) (event.Subscription, error) {

	var withdrawerRule []interface{}
	for _, withdrawerItem := range withdrawer {
		withdrawerRule = append(withdrawerRule, withdrawerItem)
	}

	logs, sub, err := _EventEmitter.contract.WatchLogs(opts, "Withdraw", withdrawerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventEmitterWithdraw)
				if err := _EventEmitter.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x8be7f12253faa74082f51d8d13cc0e76176d2ab1442fad134c53068eb958374c.
//
// Solidity: event Withdraw(address indexed withdrawer, address baseToken, address memeToken, uint256 positionId, uint256 withdrawAmount, address to, uint256 baseCollateral, uint256 baseDebtScaled, uint256 memeCollateral, uint256 memeDebtScaled)
func (_EventEmitter *EventEmitterFilterer) ParseWithdraw(log types.Log) (*EventEmitterWithdraw, error) {
	event := new(EventEmitterWithdraw)
	if err := _EventEmitter.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
