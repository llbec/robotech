package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Chain struct {
	name    string
	testnet bool
}

func NewChain(name string, test bool) *Chain {
	return &Chain{
		name:    name,
		testnet: test,
	}
}

func (c *Chain) Name() string {
	return c.name
}

func (c *Chain) TestNet() bool {
	return c.testnet
}

func (c *Chain) String() string {
	if c.testnet {
		return fmt.Sprintf("%s-testnet", c.name)
	}
	return fmt.Sprintf("%s-mainnet", c.name)
}

type TxOpts struct {
	secret    string
	gasLimit  uint64
	chainid   *big.Int
	value     *big.Int
	rpcClient *ethclient.Client
}

func (opts *TxOpts) GetTransactOpts() (*bind.TransactOpts, error) {
	//privatekey
	secret, err := crypto.HexToECDSA(opts.secret)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key %v", err)
	}

	//nonce
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := opts.rpcClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("read nonce failed %v", err)
	}

	//gas price
	gasPrice, err := opts.rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read suggest Gas price failed %v", err)
	}

	tx, err := bind.NewKeyedTransactorWithChainID(secret, opts.chainid)
	if err != nil {
		return nil, fmt.Errorf("new transaction failed %v", err)
	}
	tx.Nonce = big.NewInt(int64(nonce))
	tx.Value = opts.value       // in wei
	tx.GasLimit = opts.gasLimit // in units
	tx.GasPrice = gasPrice
	return tx, nil
}
