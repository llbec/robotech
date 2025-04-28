package ethUtils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
// maxEnd = 1000
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

func WaitTransaction(tx *types.Transaction) (*types.Receipt, error) {
	return bind.WaitMined(context.Background(), GetClient(), tx)
}

func SignMessage(privateKey *ecdsa.PrivateKey, message string) ([]byte, error) {
	hash := crypto.Keccak256Hash([]byte(message))
	//return SignHash(privateKey, accounts.TextHash([]byte(message)))
	return SignHash(privateKey, hash.Bytes())
}

func SignHash(privateKey *ecdsa.PrivateKey, messageHash []byte) ([]byte, error) {
	//messageHash := accounts.TextHash([]byte(message))

	signature, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		return []byte{}, err
	}

	//signature[crypto.RecoveryIDOffset] += 27

	//return hexutil.Encode(signature), nil
	return signature, nil
}

func RecoverSignature(signature, messageHash []byte) (common.Address, error) {
	/*signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		return common.MaxAddress, err
	}*/

	//signature[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	//messageHash := accounts.TextHash([]byte(message))

	/*sigPublicKey, err := crypto.Ecrecover(messageHash, signature)
	if err != nil {
		log.Fatal(err)
	}*/
	sigPublicKeyECDSA, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return common.MaxAddress, err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA), nil
}
