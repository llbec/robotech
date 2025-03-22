package ethUtils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client *ethclient.Client
)

func SetClient(url string) (err error) {
	client, err = ethclient.Dial(url) // 连接以太坊节点
	return
}

func GetOpt(sk string) (*bind.TransactOpts, error) {
	secret, err := crypto.HexToECDSA(sk)
	if err != nil {
		return nil, err
	}
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
