package aavesquadron

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
)

type AaveSquadron struct {
	RpcClient   *ethclient.Client
	LendingPool *aave.LendingPool
}
