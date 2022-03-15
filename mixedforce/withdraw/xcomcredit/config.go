package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/armory/aave"
	"github.com/llbec/robotech/armory/client"
	"github.com/spf13/viper"
)

var (
	contract    string
	rpcURL      string
	asset       string
	usrAddress  string
	blockPeroid int64
	gasLimit    uint64
)

const (
	CONTRACT = "contract"
	RPC      = "rpc"
	ASSET    = "asset"
	USR      = "address"
	PEROID   = "peroid"
	GASLIMIT = "gaslimit"
	maxEnd   = 1000
)

//state variable
var (
	cfg           *viper.Viper
	lendingPool   *aave.LendingPool
	atoken        *aave.AToken
	rpcClient     client.Client
	skey          string
	atAddr        common.Address
	currentHeight uint64
)

func readConfig() {
	var lock sync.Mutex
	err := cfg.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal("ReadInConfig ", err)
	}
	lock.Lock()
	defer lock.Unlock()
	contract = cfg.GetString(CONTRACT)
	rpcURL = cfg.GetString(RPC)
	asset = cfg.GetString(ASSET)
	usrAddress = cfg.GetString(USR)
	blockPeroid = cfg.GetInt64(PEROID)
	gasLimit = cfg.GetUint64(GASLIMIT)
}

func envInit() (err error) {
	rpcClient = client.NewEthChain(rpcURL, rpcURL)

	if lendingPool, err = aave.NewLendingPool(
		common.HexToAddress(contract),
		rpcClient.GetHTTPClient(),
	); err == nil {
		if data, err := lendingPool.GetReserveData(nil, common.HexToAddress(asset)); err == nil {
			atAddr = data.ATokenAddress
			atoken, err = aave.NewAToken(atAddr, rpcClient.GetHTTPClient())
			return err
		}
	}
	return
}

func generatConfig(cfgpath string) error {
	viper.Set(RPC, "192.168.11.51:1235")
	viper.Set(CONTRACT, "wallet.token")
	viper.Set(ASSET, "")
	viper.Set(USR, "")
	viper.Set(PEROID, 60)
	viper.Set(GASLIMIT, 4400000)
	path := filepath.Join(cfgpath, "credit.yaml")
	return viper.WriteConfigAs(path)
}

func showCfg() {
	fmt.Printf("Running:\n\tRPC:%v\n\tContract address:%v\n\tusr address:%v\n\tcheck period:%v\n",
		rpcURL,
		contract,
		usrAddress,
		blockPeroid)
}

func loadConfig(cfgpath string) {
	cfg = viper.New()
	cfg.AddConfigPath(cfgpath)
	cfg.SetConfigName("credit")
	cfg.SetConfigType("yaml")

	//default value
	//cfg.SetDefault(rpcURL, "")

	readConfig()
	showCfg()
}
