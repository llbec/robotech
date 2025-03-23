package ethUtils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ProcEventFunc func(types.Log) error

const (
	maxEnd = 1000
)

var (
	client *ethclient.Client
)

func SetClient(url string) (err error) {
	client, err = ethclient.Dial(url) // 连接以太坊节点
	return
}

func GetClient() *ethclient.Client {
	return client
}

func GetOpt(secret *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	if client == nil {
		return nil, errors.New("client is nil")
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	opt, err := bind.NewKeyedTransactorWithChainID(secret, chainId)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	opt.Nonce = big.NewInt(int64(nonce))
	opt.Value = big.NewInt(0)      // in wei
	opt.GasLimit = uint64(4400000) // in units
	opt.GasPrice = gasPrice
	return opt, nil
}

func SyncBlock(from, to int64, contractAddress string, eventFunc ProcEventFunc) int64 {
	log.Printf("block sync task from %v to %v, step %v", from, to, maxEnd)
	start := from
	end := start + maxEnd
	for {
		if end > to {
			end = to
		}
		//log.Printf("sync block from %v to %v\n", start, end)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(contractAddress)},
			FromBlock: big.NewInt(int64(start)),
			ToBlock:   big.NewInt(int64(end)),
		}
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("FilterLogs(%v-%v): %v\n", start, end, err)
			return start
		}
		for _, vLog := range logs {
			err = eventFunc(vLog)
			if err != nil {
				log.Printf("HandleLog(%v-%v): %v\n", vLog.BlockNumber, vLog.Index, err)
				return int64(vLog.BlockNumber) - 1
			}
		}
		if end == to {
			break
		}
		start = end + 1
		end = start + maxEnd
	}
	return end
}
