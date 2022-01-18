package client

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client interface {
	GetAuther(secretkey string, gasLimit uint64) (*bind.TransactOpts, error)
	GetHTTPClient() *ethclient.Client
	GetWSClient() *ethclient.Client
}

type ChainClient struct {
	readClient *ethclient.Client
	subsClient *ethclient.Client
}
