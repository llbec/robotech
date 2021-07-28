package aavesquadron

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
)

type AaveSquadron struct {
	rpcClient   *ethclient.Client
	LendingPool *aave.LendingPool
}
