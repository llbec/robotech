package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/spf13/viper"
)

const (
	RPC                      = "rpcs"
	SKEY                     = "secret"
	LENDINGPOOL              = "lendingpool"
	AAVEPROTOCOLDATAPROVIDER = "AaveProtocolDataProvider"
	AAVELIQUIDATIONADAPTER   = "liquidationAdapter"
	AAVEORACLE               = "AaveOracle"
	STARTHEIGHT              = "height"
	LIQUPEROID               = "liquidatePeroid"
	BLKPEROID                = "blockPeroid"
)

func GeneratConfig(dir string) error {
	viper.Set(RPC, "https://http-mainnet-node.huobichain.com")
	viper.Set(SKEY, "0x")
	viper.Set(LENDINGPOOL, "0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")
	viper.Set(AAVEPROTOCOLDATAPROVIDER, "0x18C789412E7652d094505E249d985DFE14Bd772e")
	viper.Set(AAVELIQUIDATIONADAPTER, "0xAD7958B141fd151910A67e973A4e1173d38C5801")
	viper.Set(AAVEORACLE, "0x1be62C4A97a45B04628E2B4f38F2eC71cC709e24")
	viper.Set(STARTHEIGHT, 4036695)
	//viper.Set(DMPORT, 1234)
	viper.Set(LIQUPEROID, 10)
	viper.Set(BLKPEROID, 3)

	path := filepath.Join(dir, "config.yaml")
	return viper.WriteConfigAs(path)
}

func InitEnv(dir string) {
	viper.AddConfigPath(dir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(err)
	}

	rpcURL = viper.GetString(RPC)
	Client, err = ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}

	LendingPoolAddress = viper.GetString(LENDINGPOOL)
	LendingPool, err = aave.NewLendingPool(
		common.HexToAddress(LendingPoolAddress),
		Client,
	)
	if err != nil {
		panic(err)
	}

	DataProviderAddress = viper.GetString(AAVEPROTOCOLDATAPROVIDER)
	ProtocaldataProvider, err = aave.NewAaveProtocolDataProvider(
		common.HexToAddress(DataProviderAddress),
		Client,
	)
	if err != nil {
		panic(err)
	}

	FlashAdapterAddress = viper.GetString(AAVELIQUIDATIONADAPTER)
	FlashLiquidationAdp, err = flashloan.NewAaveLiquidationAdapter(
		common.HexToAddress(FlashAdapterAddress),
		Client,
	)
	if err != nil {
		panic(err)
	}

	OracleAddress = viper.GetString(AAVEORACLE)
	AaveOracle, err = aave.NewAaveOracle(
		common.HexToAddress(OracleAddress),
		Client,
	)
	if err != nil {
		panic(err)
	}

	StartHeight = viper.GetUint64(STARTHEIGHT)
	//DmPort = viper.GetInt(DMPORT)
	LiquidatePeroid = viper.GetInt64(LIQUPEROID)
	BlockPeroid = viper.GetInt64(BLKPEROID)
	sKey = viper.GetString(SKEY)
}

func GetAuther(secretkey string) (*bind.TransactOpts, error) {
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
	nonce, err := Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("read nonce failed %v", err)
	}
	chainid, err := Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read chain ID failed %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(secret, chainid)
	if err != nil {
		return nil, fmt.Errorf("new transaction failed %v", err)
	}
	gasPrice, err := Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read suggest Gas price failed %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4400000) // in units
	auth.GasPrice = gasPrice
	return auth, nil
}
