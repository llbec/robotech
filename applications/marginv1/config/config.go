package config

import (
	"crypto/ecdsa"
	"log"
	"os"
	"path/filepath"

	"robotech/utils/ethUtils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

const (
	RPC          = "rpc"
	SKey         = "secret"
	EventEmitter = "eventEmitter"
)

var (
	sKey             *ecdsa.PrivateKey
	eventEmitterAddr common.Address
)

func GenerateConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("os.Getwd: %v\n", wd)

	viper.Set(RPC, "https://chain-rpc.com")
	viper.Set(SKey, "")
	viper.Set(EventEmitter, "")

	path := filepath.Join(wd, "config.yaml")
	return viper.WriteConfigAs(path)
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Printf("os.Getwd: %v\n", wd)
	viper.AddConfigPath(wd)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(err)
	}

	ethUtils.SetClient(viper.GetString(RPC))
	sKey, err = crypto.HexToECDSA(viper.GetString(SKey))
	if err != nil {
		panic(err)
	}

	eventEmitterAddr = common.HexToAddress(viper.GetString(EventEmitter))
}

func GetEventEmitterAddr() common.Address {
	return eventEmitterAddr
}
