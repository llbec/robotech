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

// CompflashLiquidationAdapterLiquidationParams is an auto generated low-level Go binding around an user-defined struct.
type CompflashLiquidationAdapterLiquidationParams struct {
	CollateralCtoken common.Address
	BorrowedCtoken   common.Address
	User             common.Address
	DebtToCover      *big.Int
	UseEthPath       bool
}

// CompoundLiquidationABI is the input ABI used to generate the binding from.
const CompoundLiquidationABI = "[{\"inputs\":[{\"internalType\":\"contractILendingPoolAddressesProvider\",\"name\":\"addressesProvider\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV2Router02\",\"name\":\"uniswapRouter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wethAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAsset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAsset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"receivedAmount\",\"type\":\"uint256\"}],\"name\":\"Swapped\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractILendingPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"collateralCtoken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrowedCtoken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"useEthPath\",\"type\":\"bool\"}],\"internalType\":\"structCompflashLiquidationAdapter.LiquidationParams\",\"name\":\"arg\",\"type\":\"tuple\"}],\"name\":\"Execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FLASHLOAN_PREMIUM_TOTAL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LENDING_POOL\",\"outputs\":[{\"internalType\":\"contractILendingPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_SLIPPAGE_PERCENT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE\",\"outputs\":[{\"internalType\":\"contractIPriceOracleGetter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNISWAP_ROUTER\",\"outputs\":[{\"internalType\":\"contractIUniswapV2Router02\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WETH_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"premiums\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"executeOperation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"reserveOut\",\"type\":\"address\"}],\"name\":\"getAmountsIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"reserveOut\",\"type\":\"address\"}],\"name\":\"getAmountsOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"rescueTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// CompoundLiquidation is an auto generated Go binding around an Ethereum contract.
type CompoundLiquidation struct {
	CompoundLiquidationCaller     // Read-only binding to the contract
	CompoundLiquidationTransactor // Write-only binding to the contract
	CompoundLiquidationFilterer   // Log filterer for contract events
}

// CompoundLiquidationCaller is an auto generated read-only Go binding around an Ethereum contract.
type CompoundLiquidationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompoundLiquidationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CompoundLiquidationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompoundLiquidationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CompoundLiquidationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompoundLiquidationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CompoundLiquidationSession struct {
	Contract     *CompoundLiquidation // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CompoundLiquidationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CompoundLiquidationCallerSession struct {
	Contract *CompoundLiquidationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CompoundLiquidationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CompoundLiquidationTransactorSession struct {
	Contract     *CompoundLiquidationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CompoundLiquidationRaw is an auto generated low-level Go binding around an Ethereum contract.
type CompoundLiquidationRaw struct {
	Contract *CompoundLiquidation // Generic contract binding to access the raw methods on
}

// CompoundLiquidationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CompoundLiquidationCallerRaw struct {
	Contract *CompoundLiquidationCaller // Generic read-only contract binding to access the raw methods on
}

// CompoundLiquidationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CompoundLiquidationTransactorRaw struct {
	Contract *CompoundLiquidationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCompoundLiquidation creates a new instance of CompoundLiquidation, bound to a specific deployed contract.
func NewCompoundLiquidation(address common.Address, backend bind.ContractBackend) (*CompoundLiquidation, error) {
	contract, err := bindCompoundLiquidation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CompoundLiquidation{CompoundLiquidationCaller: CompoundLiquidationCaller{contract: contract}, CompoundLiquidationTransactor: CompoundLiquidationTransactor{contract: contract}, CompoundLiquidationFilterer: CompoundLiquidationFilterer{contract: contract}}, nil
}

// NewCompoundLiquidationCaller creates a new read-only instance of CompoundLiquidation, bound to a specific deployed contract.
func NewCompoundLiquidationCaller(address common.Address, caller bind.ContractCaller) (*CompoundLiquidationCaller, error) {
	contract, err := bindCompoundLiquidation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CompoundLiquidationCaller{contract: contract}, nil
}

// NewCompoundLiquidationTransactor creates a new write-only instance of CompoundLiquidation, bound to a specific deployed contract.
func NewCompoundLiquidationTransactor(address common.Address, transactor bind.ContractTransactor) (*CompoundLiquidationTransactor, error) {
	contract, err := bindCompoundLiquidation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CompoundLiquidationTransactor{contract: contract}, nil
}

// NewCompoundLiquidationFilterer creates a new log filterer instance of CompoundLiquidation, bound to a specific deployed contract.
func NewCompoundLiquidationFilterer(address common.Address, filterer bind.ContractFilterer) (*CompoundLiquidationFilterer, error) {
	contract, err := bindCompoundLiquidation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CompoundLiquidationFilterer{contract: contract}, nil
}

// bindCompoundLiquidation binds a generic wrapper to an already deployed contract.
func bindCompoundLiquidation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CompoundLiquidationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompoundLiquidation *CompoundLiquidationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CompoundLiquidation.Contract.CompoundLiquidationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompoundLiquidation *CompoundLiquidationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.CompoundLiquidationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompoundLiquidation *CompoundLiquidationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.CompoundLiquidationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompoundLiquidation *CompoundLiquidationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CompoundLiquidation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompoundLiquidation *CompoundLiquidationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompoundLiquidation *CompoundLiquidationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.contract.Transact(opts, method, params...)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCaller) ADDRESSESPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "ADDRESSES_PROVIDER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _CompoundLiquidation.Contract.ADDRESSESPROVIDER(&_CompoundLiquidation.CallOpts)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _CompoundLiquidation.Contract.ADDRESSESPROVIDER(&_CompoundLiquidation.CallOpts)
}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint256)
func (_CompoundLiquidation *CompoundLiquidationCaller) FLASHLOANPREMIUMTOTAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "FLASHLOAN_PREMIUM_TOTAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint256)
func (_CompoundLiquidation *CompoundLiquidationSession) FLASHLOANPREMIUMTOTAL() (*big.Int, error) {
	return _CompoundLiquidation.Contract.FLASHLOANPREMIUMTOTAL(&_CompoundLiquidation.CallOpts)
}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint256)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) FLASHLOANPREMIUMTOTAL() (*big.Int, error) {
	return _CompoundLiquidation.Contract.FLASHLOANPREMIUMTOTAL(&_CompoundLiquidation.CallOpts)
}

// LENDINGPOOL is a free data retrieval call binding the contract method 0xb4dcfc77.
//
// Solidity: function LENDING_POOL() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCaller) LENDINGPOOL(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "LENDING_POOL")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LENDINGPOOL is a free data retrieval call binding the contract method 0xb4dcfc77.
//
// Solidity: function LENDING_POOL() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationSession) LENDINGPOOL() (common.Address, error) {
	return _CompoundLiquidation.Contract.LENDINGPOOL(&_CompoundLiquidation.CallOpts)
}

// LENDINGPOOL is a free data retrieval call binding the contract method 0xb4dcfc77.
//
// Solidity: function LENDING_POOL() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) LENDINGPOOL() (common.Address, error) {
	return _CompoundLiquidation.Contract.LENDINGPOOL(&_CompoundLiquidation.CallOpts)
}

// MAXSLIPPAGEPERCENT is a free data retrieval call binding the contract method 0x32e4b286.
//
// Solidity: function MAX_SLIPPAGE_PERCENT() view returns(uint256)
func (_CompoundLiquidation *CompoundLiquidationCaller) MAXSLIPPAGEPERCENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "MAX_SLIPPAGE_PERCENT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSLIPPAGEPERCENT is a free data retrieval call binding the contract method 0x32e4b286.
//
// Solidity: function MAX_SLIPPAGE_PERCENT() view returns(uint256)
func (_CompoundLiquidation *CompoundLiquidationSession) MAXSLIPPAGEPERCENT() (*big.Int, error) {
	return _CompoundLiquidation.Contract.MAXSLIPPAGEPERCENT(&_CompoundLiquidation.CallOpts)
}

// MAXSLIPPAGEPERCENT is a free data retrieval call binding the contract method 0x32e4b286.
//
// Solidity: function MAX_SLIPPAGE_PERCENT() view returns(uint256)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) MAXSLIPPAGEPERCENT() (*big.Int, error) {
	return _CompoundLiquidation.Contract.MAXSLIPPAGEPERCENT(&_CompoundLiquidation.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x38013f02.
//
// Solidity: function ORACLE() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCaller) ORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ORACLE is a free data retrieval call binding the contract method 0x38013f02.
//
// Solidity: function ORACLE() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationSession) ORACLE() (common.Address, error) {
	return _CompoundLiquidation.Contract.ORACLE(&_CompoundLiquidation.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x38013f02.
//
// Solidity: function ORACLE() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) ORACLE() (common.Address, error) {
	return _CompoundLiquidation.Contract.ORACLE(&_CompoundLiquidation.CallOpts)
}

// UNISWAPROUTER is a free data retrieval call binding the contract method 0xd8264920.
//
// Solidity: function UNISWAP_ROUTER() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCaller) UNISWAPROUTER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "UNISWAP_ROUTER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UNISWAPROUTER is a free data retrieval call binding the contract method 0xd8264920.
//
// Solidity: function UNISWAP_ROUTER() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationSession) UNISWAPROUTER() (common.Address, error) {
	return _CompoundLiquidation.Contract.UNISWAPROUTER(&_CompoundLiquidation.CallOpts)
}

// UNISWAPROUTER is a free data retrieval call binding the contract method 0xd8264920.
//
// Solidity: function UNISWAP_ROUTER() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) UNISWAPROUTER() (common.Address, error) {
	return _CompoundLiquidation.Contract.UNISWAPROUTER(&_CompoundLiquidation.CallOpts)
}

// WETHADDRESS is a free data retrieval call binding the contract method 0x040141e5.
//
// Solidity: function WETH_ADDRESS() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCaller) WETHADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "WETH_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETHADDRESS is a free data retrieval call binding the contract method 0x040141e5.
//
// Solidity: function WETH_ADDRESS() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationSession) WETHADDRESS() (common.Address, error) {
	return _CompoundLiquidation.Contract.WETHADDRESS(&_CompoundLiquidation.CallOpts)
}

// WETHADDRESS is a free data retrieval call binding the contract method 0x040141e5.
//
// Solidity: function WETH_ADDRESS() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) WETHADDRESS() (common.Address, error) {
	return _CompoundLiquidation.Contract.WETHADDRESS(&_CompoundLiquidation.CallOpts)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0xcdf58cd6.
//
// Solidity: function getAmountsIn(uint256 amountOut, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_CompoundLiquidation *CompoundLiquidationCaller) GetAmountsIn(opts *bind.CallOpts, amountOut *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "getAmountsIn", amountOut, reserveIn, reserveOut)

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
func (_CompoundLiquidation *CompoundLiquidationSession) GetAmountsIn(amountOut *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _CompoundLiquidation.Contract.GetAmountsIn(&_CompoundLiquidation.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0xcdf58cd6.
//
// Solidity: function getAmountsIn(uint256 amountOut, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_CompoundLiquidation *CompoundLiquidationCallerSession) GetAmountsIn(amountOut *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _CompoundLiquidation.Contract.GetAmountsIn(&_CompoundLiquidation.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xbaf7fa99.
//
// Solidity: function getAmountsOut(uint256 amountIn, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_CompoundLiquidation *CompoundLiquidationCaller) GetAmountsOut(opts *bind.CallOpts, amountIn *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "getAmountsOut", amountIn, reserveIn, reserveOut)

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
func (_CompoundLiquidation *CompoundLiquidationSession) GetAmountsOut(amountIn *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _CompoundLiquidation.Contract.GetAmountsOut(&_CompoundLiquidation.CallOpts, amountIn, reserveIn, reserveOut)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xbaf7fa99.
//
// Solidity: function getAmountsOut(uint256 amountIn, address reserveIn, address reserveOut) view returns(uint256, uint256, uint256, uint256, address[])
func (_CompoundLiquidation *CompoundLiquidationCallerSession) GetAmountsOut(amountIn *big.Int, reserveIn common.Address, reserveOut common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _CompoundLiquidation.Contract.GetAmountsOut(&_CompoundLiquidation.CallOpts, amountIn, reserveIn, reserveOut)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundLiquidation.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationSession) Owner() (common.Address, error) {
	return _CompoundLiquidation.Contract.Owner(&_CompoundLiquidation.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CompoundLiquidation *CompoundLiquidationCallerSession) Owner() (common.Address, error) {
	return _CompoundLiquidation.Contract.Owner(&_CompoundLiquidation.CallOpts)
}

// Execute is a paid mutator transaction binding the contract method 0x22cf310a.
//
// Solidity: function Execute((address,address,address,uint256,bool) arg) returns()
func (_CompoundLiquidation *CompoundLiquidationTransactor) Execute(opts *bind.TransactOpts, arg CompflashLiquidationAdapterLiquidationParams) (*types.Transaction, error) {
	return _CompoundLiquidation.contract.Transact(opts, "Execute", arg)
}

// Execute is a paid mutator transaction binding the contract method 0x22cf310a.
//
// Solidity: function Execute((address,address,address,uint256,bool) arg) returns()
func (_CompoundLiquidation *CompoundLiquidationSession) Execute(arg CompflashLiquidationAdapterLiquidationParams) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.Execute(&_CompoundLiquidation.TransactOpts, arg)
}

// Execute is a paid mutator transaction binding the contract method 0x22cf310a.
//
// Solidity: function Execute((address,address,address,uint256,bool) arg) returns()
func (_CompoundLiquidation *CompoundLiquidationTransactorSession) Execute(arg CompflashLiquidationAdapterLiquidationParams) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.Execute(&_CompoundLiquidation.TransactOpts, arg)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_CompoundLiquidation *CompoundLiquidationTransactor) ExecuteOperation(opts *bind.TransactOpts, assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _CompoundLiquidation.contract.Transact(opts, "executeOperation", assets, amounts, premiums, initiator, params)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_CompoundLiquidation *CompoundLiquidationSession) ExecuteOperation(assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.ExecuteOperation(&_CompoundLiquidation.TransactOpts, assets, amounts, premiums, initiator, params)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_CompoundLiquidation *CompoundLiquidationTransactorSession) ExecuteOperation(assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.ExecuteOperation(&_CompoundLiquidation.TransactOpts, assets, amounts, premiums, initiator, params)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CompoundLiquidation *CompoundLiquidationTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundLiquidation.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CompoundLiquidation *CompoundLiquidationSession) RenounceOwnership() (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.RenounceOwnership(&_CompoundLiquidation.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CompoundLiquidation *CompoundLiquidationTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.RenounceOwnership(&_CompoundLiquidation.TransactOpts)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x00ae3bf8.
//
// Solidity: function rescueTokens(address token) returns()
func (_CompoundLiquidation *CompoundLiquidationTransactor) RescueTokens(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _CompoundLiquidation.contract.Transact(opts, "rescueTokens", token)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x00ae3bf8.
//
// Solidity: function rescueTokens(address token) returns()
func (_CompoundLiquidation *CompoundLiquidationSession) RescueTokens(token common.Address) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.RescueTokens(&_CompoundLiquidation.TransactOpts, token)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x00ae3bf8.
//
// Solidity: function rescueTokens(address token) returns()
func (_CompoundLiquidation *CompoundLiquidationTransactorSession) RescueTokens(token common.Address) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.RescueTokens(&_CompoundLiquidation.TransactOpts, token)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CompoundLiquidation *CompoundLiquidationTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CompoundLiquidation.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CompoundLiquidation *CompoundLiquidationSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.TransferOwnership(&_CompoundLiquidation.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CompoundLiquidation *CompoundLiquidationTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CompoundLiquidation.Contract.TransferOwnership(&_CompoundLiquidation.TransactOpts, newOwner)
}

// CompoundLiquidationOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CompoundLiquidation contract.
type CompoundLiquidationOwnershipTransferredIterator struct {
	Event *CompoundLiquidationOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CompoundLiquidationOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundLiquidationOwnershipTransferred)
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
		it.Event = new(CompoundLiquidationOwnershipTransferred)
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
func (it *CompoundLiquidationOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundLiquidationOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundLiquidationOwnershipTransferred represents a OwnershipTransferred event raised by the CompoundLiquidation contract.
type CompoundLiquidationOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CompoundLiquidation *CompoundLiquidationFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CompoundLiquidationOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CompoundLiquidation.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CompoundLiquidationOwnershipTransferredIterator{contract: _CompoundLiquidation.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CompoundLiquidation *CompoundLiquidationFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CompoundLiquidationOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CompoundLiquidation.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundLiquidationOwnershipTransferred)
				if err := _CompoundLiquidation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CompoundLiquidation *CompoundLiquidationFilterer) ParseOwnershipTransferred(log types.Log) (*CompoundLiquidationOwnershipTransferred, error) {
	event := new(CompoundLiquidationOwnershipTransferred)
	if err := _CompoundLiquidation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundLiquidationSwappedIterator is returned from FilterSwapped and is used to iterate over the raw logs and unpacked data for Swapped events raised by the CompoundLiquidation contract.
type CompoundLiquidationSwappedIterator struct {
	Event *CompoundLiquidationSwapped // Event containing the contract specifics and raw log

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
func (it *CompoundLiquidationSwappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundLiquidationSwapped)
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
		it.Event = new(CompoundLiquidationSwapped)
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
func (it *CompoundLiquidationSwappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundLiquidationSwappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundLiquidationSwapped represents a Swapped event raised by the CompoundLiquidation contract.
type CompoundLiquidationSwapped struct {
	FromAsset      common.Address
	ToAsset        common.Address
	FromAmount     *big.Int
	ReceivedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSwapped is a free log retrieval operation binding the contract event 0xa078c4190abe07940190effc1846be0ccf03ad6007bc9e93f9697d0b460befbb.
//
// Solidity: event Swapped(address fromAsset, address toAsset, uint256 fromAmount, uint256 receivedAmount)
func (_CompoundLiquidation *CompoundLiquidationFilterer) FilterSwapped(opts *bind.FilterOpts) (*CompoundLiquidationSwappedIterator, error) {

	logs, sub, err := _CompoundLiquidation.contract.FilterLogs(opts, "Swapped")
	if err != nil {
		return nil, err
	}
	return &CompoundLiquidationSwappedIterator{contract: _CompoundLiquidation.contract, event: "Swapped", logs: logs, sub: sub}, nil
}

// WatchSwapped is a free log subscription operation binding the contract event 0xa078c4190abe07940190effc1846be0ccf03ad6007bc9e93f9697d0b460befbb.
//
// Solidity: event Swapped(address fromAsset, address toAsset, uint256 fromAmount, uint256 receivedAmount)
func (_CompoundLiquidation *CompoundLiquidationFilterer) WatchSwapped(opts *bind.WatchOpts, sink chan<- *CompoundLiquidationSwapped) (event.Subscription, error) {

	logs, sub, err := _CompoundLiquidation.contract.WatchLogs(opts, "Swapped")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundLiquidationSwapped)
				if err := _CompoundLiquidation.contract.UnpackLog(event, "Swapped", log); err != nil {
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
func (_CompoundLiquidation *CompoundLiquidationFilterer) ParseSwapped(log types.Log) (*CompoundLiquidationSwapped, error) {
	event := new(CompoundLiquidationSwapped)
	if err := _CompoundLiquidation.contract.UnpackLog(event, "Swapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
