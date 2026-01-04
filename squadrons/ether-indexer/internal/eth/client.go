package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

func Dial(rpc string) (*ethclient.Client, error) {
	return ethclient.Dial(rpc)
}
