package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// check interface
var _ Client = (*ChainClient)(nil)

func NewEthChain(http, ws string) Client {
	cr, _ := ethclient.Dial(http)
	cs, _ := ethclient.Dial(ws)
	return &ChainClient{
		readClient: cr,
		subsClient: cs,
	}
}

func (c *ChainClient) GetAuther(secretkey string, gasLimit uint64) (*bind.TransactOpts, error) {
	secret, err := crypto.HexToECDSA(secretkey)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key %v", err)
	}
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.readClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("read nonce failed %v", err)
	}
	chainid, err := c.readClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read chain ID failed %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(secret, chainid)
	if err != nil {
		return nil, fmt.Errorf("new transaction failed %v", err)
	}
	gasPrice, err := c.readClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read suggest Gas price failed %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   //uint64(4400000) // in units
	auth.GasPrice = gasPrice
	return auth, nil
}

func (c *ChainClient) GetHTTPClient() *ethclient.Client {
	return c.readClient
}

func (c *ChainClient) GetWSClient() *ethclient.Client {
	return c.subsClient
}
