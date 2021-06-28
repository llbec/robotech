package fallbackteam

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/bala"
)

type FallBack struct {
	rpc       *ethclient.Client
	contract  *bala.FallBack
	secretkey string
}

func NewFallback(url, addr string) *FallBack {
	node, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	ct, err := bala.NewFallBack(common.HexToAddress(addr), node)
	if err != nil {
		panic(err)
	}

	return &FallBack{rpc: node, contract: ct}
}

func (fback *FallBack) SetKey(key string) {
	fback.secretkey = key
}

func (fback *FallBack) GetAuther() (*bind.TransactOpts, error) {
	secret, err := crypto.HexToECDSA(fback.secretkey)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key %v", err)
	}
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := fback.rpc.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("read nonce failed %v", err)
	}
	chainid, err := fback.rpc.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read chain ID failed %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(secret, chainid)
	if err != nil {
		return nil, fmt.Errorf("new transaction failed %v", err)
	}
	gasPrice, err := fback.rpc.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read suggest Gas price failed %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4400000) // in units
	auth.GasPrice = gasPrice
	return auth, nil
}

func (fback *FallBack) GetPrice(address string) (*big.Int, error) {
	p, e := fback.contract.GetAssetPrice(nil, common.HexToAddress(address))
	if e != nil {
		return big.NewInt(0), e
	}
	return p, nil
}

func (fback *FallBack) SetPrice(address string, price *big.Int) error {
	auther, err := fback.GetAuther()
	if err != nil {
		return fmt.Errorf("make transaction failed: %v", err)
	}
	_, err = fback.contract.SetAssetPrice(auther, common.HexToAddress(address), price)
	if err != nil {
		return err
	}
	return nil
}
