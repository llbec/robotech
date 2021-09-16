package aavesquadron

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
)

func NewAaveSquadron() *AaveSquadron {
	return &AaveSquadron{}
}

func (aSquadron *AaveSquadron) SetRpcClient(url string) error {
	client, err := ethclient.Dial(url)
	if err != nil {
		return err
	}
	aSquadron.RpcClient = client
	return nil
}

func (aSquadron *AaveSquadron) SetLendingPool(addr string) error {
	pool, err := aave.NewLendingPool(common.HexToAddress(addr), aSquadron.RpcClient)
	if err != nil {
		return err
	}
	aSquadron.LendingPool = pool
	return nil
}

func (aSquadron *AaveSquadron) GetTxOpts(skey string) (*bind.TransactOpts, error) {
	secret, err := crypto.HexToECDSA(skey)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key %v", err)
	}
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := aSquadron.RpcClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("read nonce failed %v", err)
	}
	chainid, err := aSquadron.RpcClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read chain ID failed %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(secret, chainid)
	if err != nil {
		return nil, fmt.Errorf("new transaction failed %v", err)
	}
	gasPrice, err := aSquadron.RpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read suggest Gas price failed %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(9000000) // in units
	auth.GasPrice = gasPrice
	return auth, nil
}

func (aSquadron *AaveSquadron) GetAToken(addr string) (*aave.AToken, error) {
	return aave.NewAToken(common.HexToAddress(addr), aSquadron.RpcClient)
}
