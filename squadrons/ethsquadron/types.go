package ethsquadron

import "github.com/ethereum/go-ethereum/ethclient"

type EthSquadron struct {
	node *ethclient.Client
}
