// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package flashloan

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

// FlashLiquidationAdapterLiquidationParams is an auto generated low-level Go binding around an user-defined struct.
type FlashLiquidationAdapterLiquidationParams struct {
	CollateralAsset common.Address
	BorrowedAsset   common.Address
	User            common.Address
	DebtToCover     *big.Int
	UseEthPath      bool
}

// AaveLiquidationAdapterABI is the input ABI used to generate the binding from.
const AaveLiquidationAdapterABI = "[{\"inputs\":[{\"internalType\":\"contractILendingPoolAddressesProvider\",\"name\":\"addressesProvider\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV2Router02\",\"name\":\"uniswapRouter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wethAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAsset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAsset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"receivedAmount\",\"type\":\"uint256\"}],\"name\":\"Swapped\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractILendingPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"collateralAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrowedAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"useEthPath\",\"type\":\"bool\"}],\"internalType\":\"structFlashLiquidationAdapter.LiquidationParams\",\"name\":\"arg\",\"type\":\"tuple\"}],\"name\":\"Execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FLASHLOAN_PREMIUM_TOTAL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LENDING_POOL\",\"outputs\":[{\"internalType\":\"contractILendingPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_SLIPPAGE_PERCENT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE\",\"outputs\":[{\"internalType\":\"contractIPriceOracleGetter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNISWAP_ROUTER\",\"outputs\":[{\"internalType\":\"contractIUniswapV2Router02\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WETH_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"premiums\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"executeOperation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"reserveOut\",\"type\":\"address\"}],\"name\":\"getAmountsIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"reserveOut\",\"type\":\"address\"}],\"name\":\"getAmountsOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"rescueTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AaveLiquidationAdapter is an auto generated Go binding around an Ethereum contract.
type AaveLiquidationAdapter struct {
	AaveLiquidationAdapterCaller     // Read-only binding to the contract
	AaveLiquidationAdapterTransactor // Write-only binding to the contract
	AaveLiquidationAdapterFilterer   // Log filterer for contract events
}

// AaveLiquidationAdapterCaller is an auto generated read-only Go binding around an Ethereum contract.
type AaveLiquidationAdapterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveLiquidationAdapterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AaveLiquidationAdapterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveLiquidationAdapterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AaveLiquidationAdapterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveLiquidationAdapterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AaveLiquidationAdapterSession struct {
	Contract     *AaveLiquidationAdapter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// AaveLiquidationAdapterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AaveLiquidationAdapterCallerSession struct {
	Contract *AaveLiquidationAdapterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// AaveLiquidationAdapterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AaveLiquidationAdapterTransactorSession struct {
	Contract     *AaveLiquidationAdapterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// AaveLiquidationAdapterRaw is an auto generated low-level Go binding around an Ethereum contract.
type AaveLiquidationAdapterRaw struct {
	Contract *AaveLiquidationAdapter // Generic contract binding to access the raw methods on
}

// AaveLiquidationAdapterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AaveLiquidationAdapterCallerRaw struct {
	Contract *AaveLiquidationAdapterCaller // Generic read-only contract binding to access the raw methods on
}

// AaveLiquidationAdapterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AaveLiquidationAdapterTransactorRaw struct {
	Contract *AaveLiquidationAdapterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAaveLiquidationAdapter creates a new instance of AaveLiquidationAdapter, bound to a specific deployed contract.
func NewAaveLiquidationAdapter(address common.Address, backend bind.ContractBackend) (*AaveLiquidationAdapter, error) {
	contract, err := bindAaveLiquidationAdapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AaveLiquidationAdapter{AaveLiquidationAdapterCaller: AaveLiquidationAdapterCaller{contract: contract}, AaveLiquidationAdapterTransactor: AaveLiquidationAdapterTransactor{contract: contract}, AaveLiquidationAdapterFilterer: AaveLiquidationAdapterFilterer{contract: contract}}, nil
}

// NewAaveLiquidationAdapterCaller creates a new read-only instance of AaveLiquidationAdapter, bound to a specific deployed contract.
func NewAaveLiquidationAdapterCaller(address common.Address, caller bind.ContractCaller) (*AaveLiquidationAdapterCaller, error) {
	contract, err := bindAaveLiquidationAdapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AaveLiquidationAdapterCaller{contract: contract}, nil
}

// NewAaveLiquidationAdapterTransactor creates a new write-only instance of AaveLiquidationAdapter, bound to a specific deployed contract.
func NewAaveLiquidationAdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*AaveLiquidationAdapterTransactor, error) {
	contract, err := bindAaveLiquidationAdapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AaveLiquidationAdapterTransactor{contract: contract}, nil
}

// NewAaveLiquidationAdapterFilterer creates a new log filterer instance of AaveLiquidationAdapter, bound to a specific deployed contract.
func NewAaveLiquidationAdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*AaveLiquidationAdapterFilterer, error) {
	contract, err := bindAaveLiquidationAdapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AaveLiquidationAdapterFilterer{contract: contract}, nil
}

// bindAaveLiquidationAdapter binds a generic wrapper to an already deployed contract.
func bindAaveLiquidationAdapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AaveLiquidationAdapterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AaveLiquidationAdapter *AaveLiquidationAdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AaveLiquidationAdapter.Contract.AaveLiquidationAdapterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AaveLiquidationAdapter *AaveLiquidationAdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.AaveLiquidationAdapterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AaveLiquidationAdapter *AaveLiquidationAdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.AaveLiquidationAdapterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AaveLiquidationAdapter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.contract.Transact(opts, method, params...)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) ADDRESSESPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "ADDRESSES_PROVIDER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.ADDRESSESPROVIDER(&_AaveLiquidationAdapter.CallOpts)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.ADDRESSESPROVIDER(&_AaveLiquidationAdapter.CallOpts)
}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint256)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) FLASHLOANPREMIUMTOTAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "FLASHLOAN_PREMIUM_TOTAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint256)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) FLASHLOANPREMIUMTOTAL() (*big.Int, error) {
	return _AaveLiquidationAdapter.Contract.FLASHLOANPREMIUMTOTAL(&_AaveLiquidationAdapter.CallOpts)
}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint256)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) FLASHLOANPREMIUMTOTAL() (*big.Int, error) {
	return _AaveLiquidationAdapter.Contract.FLASHLOANPREMIUMTOTAL(&_AaveLiquidationAdapter.CallOpts)
}

// LENDINGPOOL is a free data retrieval call binding the contract method 0xb4dcfc77.
//
// Solidity: function LENDING_POOL() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) LENDINGPOOL(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "LENDING_POOL")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LENDINGPOOL is a free data retrieval call binding the contract method 0xb4dcfc77.
//
// Solidity: function LENDING_POOL() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) LENDINGPOOL() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.LENDINGPOOL(&_AaveLiquidationAdapter.CallOpts)
}

// LENDINGPOOL is a free data retrieval call binding the contract method 0xb4dcfc77.
//
// Solidity: function LENDING_POOL() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) LENDINGPOOL() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.LENDINGPOOL(&_AaveLiquidationAdapter.CallOpts)
}

// MAXSLIPPAGEPERCENT is a free data retrieval call binding the contract method 0x32e4b286.
//
// Solidity: function MAX_SLIPPAGE_PERCENT() view returns(uint256)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) MAXSLIPPAGEPERCENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "MAX_SLIPPAGE_PERCENT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSLIPPAGEPERCENT is a free data retrieval call binding the contract method 0x32e4b286.
//
// Solidity: function MAX_SLIPPAGE_PERCENT() view returns(uint256)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) MAXSLIPPAGEPERCENT() (*big.Int, error) {
	return _AaveLiquidationAdapter.Contract.MAXSLIPPAGEPERCENT(&_AaveLiquidationAdapter.CallOpts)
}

// MAXSLIPPAGEPERCENT is a free data retrieval call binding the contract method 0x32e4b286.
//
// Solidity: function MAX_SLIPPAGE_PERCENT() view returns(uint256)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) MAXSLIPPAGEPERCENT() (*big.Int, error) {
	return _AaveLiquidationAdapter.Contract.MAXSLIPPAGEPERCENT(&_AaveLiquidationAdapter.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x38013f02.
//
// Solidity: function ORACLE() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) ORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ORACLE is a free data retrieval call binding the contract method 0x38013f02.
//
// Solidity: function ORACLE() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) ORACLE() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.ORACLE(&_AaveLiquidationAdapter.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x38013f02.
//
// Solidity: function ORACLE() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) ORACLE() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.ORACLE(&_AaveLiquidationAdapter.CallOpts)
}

// UNISWAPROUTER is a free data retrieval call binding the contract method 0xd8264920.
//
// Solidity: function UNISWAP_ROUTER() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) UNISWAPROUTER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "UNISWAP_ROUTER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UNISWAPROUTER is a free data retrieval call binding the contract method 0xd8264920.
//
// Solidity: function UNISWAP_ROUTER() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) UNISWAPROUTER() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.UNISWAPROUTER(&_AaveLiquidationAdapter.CallOpts)
}

// UNISWAPROUTER is a free data retrieval call binding the contract method 0xd8264920.
//
// Solidity: function UNISWAP_ROUTER() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) UNISWAPROUTER() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.UNISWAPROUTER(&_AaveLiquidationAdapter.CallOpts)
}

// WETHADDRESS is a free data retrieval call binding the contract method 0x040141e5.
//
// Solidity: function WETH_ADDRESS() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) WETHADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "WETH_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETHADDRESS is a free data retrieval call binding the contract method 0x040141e5.
//
// Solidity: function WETH_ADDRESS() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) WETHADDRESS() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.WETHADDRESS(&_AaveLiquidationAdapter.CallOpts)
}

// WETHADDRESS is a free data retrieval call binding the contract method 0x040141e5.
//
// Solidity: function WETH_ADDRESS() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) WETHADDRESS() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.WETHADDRESS(&_AaveLiquidationAdapter.CallOpts)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0xcdf58cd6.
//
// Solidity: function getAmountsIn(uint256 amountOut, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) GetAmountsIn(opts *bind.CallOpts, amountOut *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "getAmountsIn", amountOut, reserveIn, reserveOut)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return out0, out1, out2, out3, out4, err

}

// GetAmountsIn is a free data retrieval call binding the contract method 0xcdf58cd6.
//
// Solidity: function getAmountsIn(uint256 amountOut, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) GetAmountsIn(amountOut *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _AaveLiquidationAdapter.Contract.GetAmountsIn(&_AaveLiquidationAdapter.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0xcdf58cd6.
//
// Solidity: function getAmountsIn(uint256 amountOut, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) GetAmountsIn(amountOut *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _AaveLiquidationAdapter.Contract.GetAmountsIn(&_AaveLiquidationAdapter.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xbaf7fa99.
//
// Solidity: function getAmountsOut(uint256 amountIn, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) GetAmountsOut(opts *bind.CallOpts, amountIn *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "getAmountsOut", amountIn, reserveIn, reserveOut)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return out0, out1, out2, out3, out4, err

}

// GetAmountsOut is a free data retrieval call binding the contract method 0xbaf7fa99.
//
// Solidity: function getAmountsOut(uint256 amountIn, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) GetAmountsOut(amountIn *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _AaveLiquidationAdapter.Contract.GetAmountsOut(&_AaveLiquidationAdapter.CallOpts, amountIn, reserveIn, reserveOut)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xbaf7fa99.
//
// Solidity: function getAmountsOut(uint256 amountIn, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) GetAmountsOut(amountIn *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _AaveLiquidationAdapter.Contract.GetAmountsOut(&_AaveLiquidationAdapter.CallOpts, amountIn, reserveIn, reserveOut)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveLiquidationAdapter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) Owner() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.Owner(&_AaveLiquidationAdapter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterCallerSession) Owner() (common.Address, error) {
	return _AaveLiquidationAdapter.Contract.Owner(&_AaveLiquidationAdapter.CallOpts)
}

// Execute is a paid mutator transaction binding the contract method 0x22cf310a.
//
// Solidity: function Execute((address,address,address,uint256,bool) arg) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactor) Execute(opts *bind.TransactOpts, arg FlashLiquidationAdapterLiquidationParams) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.contract.Transact(opts, "Execute", arg)
}

// Execute is a paid mutator transaction binding the contract method 0x22cf310a.
//
// Solidity: function Execute((address,address,address,uint256,bool) arg) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) Execute(arg FlashLiquidationAdapterLiquidationParams) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.Execute(&_AaveLiquidationAdapter.TransactOpts, arg)
}

// Execute is a paid mutator transaction binding the contract method 0x22cf310a.
//
// Solidity: function Execute((address,address,address,uint256,bool) arg) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorSession) Execute(arg FlashLiquidationAdapterLiquidationParams) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.Execute(&_AaveLiquidationAdapter.TransactOpts, arg)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactor) ExecuteOperation(opts *bind.TransactOpts, assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.contract.Transact(opts, "executeOperation", assets, amounts, premiums, initiator, params)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) ExecuteOperation(assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.ExecuteOperation(&_AaveLiquidationAdapter.TransactOpts, assets, amounts, premiums, initiator, params)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorSession) ExecuteOperation(assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.ExecuteOperation(&_AaveLiquidationAdapter.TransactOpts, assets, amounts, premiums, initiator, params)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) RenounceOwnership() (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.RenounceOwnership(&_AaveLiquidationAdapter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.RenounceOwnership(&_AaveLiquidationAdapter.TransactOpts)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x00ae3bf8.
//
// Solidity: function rescueTokens(address token) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactor) RescueTokens(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.contract.Transact(opts, "rescueTokens", token)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x00ae3bf8.
//
// Solidity: function rescueTokens(address token) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) RescueTokens(token common.Address) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.RescueTokens(&_AaveLiquidationAdapter.TransactOpts, token)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x00ae3bf8.
//
// Solidity: function rescueTokens(address token) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorSession) RescueTokens(token common.Address) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.RescueTokens(&_AaveLiquidationAdapter.TransactOpts, token)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.TransferOwnership(&_AaveLiquidationAdapter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AaveLiquidationAdapter *AaveLiquidationAdapterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AaveLiquidationAdapter.Contract.TransferOwnership(&_AaveLiquidationAdapter.TransactOpts, newOwner)
}

// AaveLiquidationAdapterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AaveLiquidationAdapter contract.
type AaveLiquidationAdapterOwnershipTransferredIterator struct {
	Event *AaveLiquidationAdapterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AaveLiquidationAdapterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveLiquidationAdapterOwnershipTransferred)
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
		it.Event = new(AaveLiquidationAdapterOwnershipTransferred)
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
func (it *AaveLiquidationAdapterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveLiquidationAdapterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveLiquidationAdapterOwnershipTransferred represents a OwnershipTransferred event raised by the AaveLiquidationAdapter contract.
type AaveLiquidationAdapterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AaveLiquidationAdapterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AaveLiquidationAdapter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AaveLiquidationAdapterOwnershipTransferredIterator{contract: _AaveLiquidationAdapter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AaveLiquidationAdapterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AaveLiquidationAdapter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveLiquidationAdapterOwnershipTransferred)
				if err := _AaveLiquidationAdapter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterFilterer) ParseOwnershipTransferred(log types.Log) (*AaveLiquidationAdapterOwnershipTransferred, error) {
	event := new(AaveLiquidationAdapterOwnershipTransferred)
	if err := _AaveLiquidationAdapter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveLiquidationAdapterSwappedIterator is returned from FilterSwapped and is used to iterate over the raw logs and unpacked data for Swapped events raised by the AaveLiquidationAdapter contract.
type AaveLiquidationAdapterSwappedIterator struct {
	Event *AaveLiquidationAdapterSwapped // Event containing the contract specifics and raw log

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
func (it *AaveLiquidationAdapterSwappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveLiquidationAdapterSwapped)
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
		it.Event = new(AaveLiquidationAdapterSwapped)
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
func (it *AaveLiquidationAdapterSwappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveLiquidationAdapterSwappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveLiquidationAdapterSwapped represents a Swapped event raised by the AaveLiquidationAdapter contract.
type AaveLiquidationAdapterSwapped struct {
	FromAsset      common.Address
	ToAsset        common.Address
	FromAmount     *big.Int
	ReceivedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSwapped is a free log retrieval operation binding the contract event 0xa078c4190abe07940190effc1846be0ccf03ad6007bc9e93f9697d0b460befbb.
//
// Solidity: event Swapped(address fromAsset, address toAsset, uint256 fromAmount, uint256 receivedAmount)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterFilterer) FilterSwapped(opts *bind.FilterOpts) (*AaveLiquidationAdapterSwappedIterator, error) {

	logs, sub, err := _AaveLiquidationAdapter.contract.FilterLogs(opts, "Swapped")
	if err != nil {
		return nil, err
	}
	return &AaveLiquidationAdapterSwappedIterator{contract: _AaveLiquidationAdapter.contract, event: "Swapped", logs: logs, sub: sub}, nil
}

// WatchSwapped is a free log subscription operation binding the contract event 0xa078c4190abe07940190effc1846be0ccf03ad6007bc9e93f9697d0b460befbb.
//
// Solidity: event Swapped(address fromAsset, address toAsset, uint256 fromAmount, uint256 receivedAmount)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterFilterer) WatchSwapped(opts *bind.WatchOpts, sink chan<- *AaveLiquidationAdapterSwapped) (event.Subscription, error) {

	logs, sub, err := _AaveLiquidationAdapter.contract.WatchLogs(opts, "Swapped")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveLiquidationAdapterSwapped)
				if err := _AaveLiquidationAdapter.contract.UnpackLog(event, "Swapped", log); err != nil {
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

// ParseSwapped is a log parse operation binding the contract event 0xa078c4190abe07940190effc1846be0ccf03ad6007bc9e93f9697d0b460befbb.
//
// Solidity: event Swapped(address fromAsset, address toAsset, uint256 fromAmount, uint256 receivedAmount)
func (_AaveLiquidationAdapter *AaveLiquidationAdapterFilterer) ParseSwapped(log types.Log) (*AaveLiquidationAdapterSwapped, error) {
	event := new(AaveLiquidationAdapterSwapped)
	if err := _AaveLiquidationAdapter.contract.UnpackLog(event, "Swapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
