package main

import (
	"crypto/ecdsa"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/llbec/robotech/utils/ethUtis"
	"github.com/spf13/viper"
)

const (
	RPC  = "rpc"
	SKEY = "secret"
)

var (
	sKey *ecdsa.PrivateKey
)

func GeneratConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("os.Getwd: %v\n", wd)

	viper.Set(RPC, "https://chain-rpc.com")
	viper.Set(SKEY, "")

	path := filepath.Join(wd, "config.yaml")
	return viper.WriteConfigAs(path)
}

func InitEnv() {
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

	ethUtis.SetClient(viper.GetString(RPC))
	sKey, err = crypto.HexToECDSA(viper.GetString(SKEY))
	if err != nil {
		panic(err)
	}
}
