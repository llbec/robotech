package ethsquadron

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

var (
	EthConfig = params.MainnetChainConfig
)

func NewEthSquadron(url string) (*EthSquadron, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthSquadron{client}, nil
}
