package contract

import (
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	contractABI string
)

type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	Contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	Contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	Contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

func loadABI(abiFile string) error {
	f, err := ioutil.ReadFile(abiFile)
	if err != nil {
		return err
	}
	contractABI = string(f)
	return nil
}

func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func ConnectByABI(address, url, abiFile string) (*Contract, error) {
	err := loadABI(abiFile)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	backend := bind.ContractBackend(client)
	contract, err := bindContract(common.HexToAddress(address), backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{
		ContractCaller:     ContractCaller{Contract: contract},
		ContractTransactor: ContractTransactor{Contract: contract},
		ContractFilterer:   ContractFilterer{Contract: contract},
	}, nil
}

func ConnectByInterface(address, url, ifFile string) (*Contract, error) {
	err := ReadInterfaces(ifFile)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	backend := bind.ContractBackend(client)
	contract, err := bindContract(common.HexToAddress(address), backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{
		ContractCaller:     ContractCaller{Contract: contract},
		ContractTransactor: ContractTransactor{Contract: contract},
		ContractFilterer:   ContractFilterer{Contract: contract},
	}, nil
}
