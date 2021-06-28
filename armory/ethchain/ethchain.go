package ethchain

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthChain struct {
	ReadClient *ethclient.Client
	SubsClient *ethclient.Client
}

func NewEthChain(http, ws string) *EthChain {
	cr, err := ethclient.Dial(http)
	if err != nil {
		log.Panic(err)
	}
	cs, err := ethclient.Dial(ws)
	if err != nil {
		log.Panic(err)
	}
	return &EthChain{
		ReadClient: cr,
		SubsClient: cs,
	}
}
