package config

import (
	"crypto/ecdsa"
	"log"

	"robotech/utils"
	"robotech/utils/ethUtils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	rpc               = "rpc"
	secret            = "secret"
	eventEmitter      = "eventEmitter"
	reader            = "reader"
	dataStore         = "dataStore"
	exchangeRouter    = "exchangeRouter"
	baseToken         = "baseToken"
	daemon            = "daemonPort"
	last_block_number = "lastBlockNumber"
	read_block_period = "readBlockPeriod"
)

var (
	sKey               *ecdsa.PrivateKey
	eventEmitterAddr   common.Address
	readerAddr         common.Address
	dataStoreAddr      common.Address
	exchangeRouterAddr common.Address
	baseTokenAddr      common.Address
	daemonPort         int
	workDir            string
	blockNumber        uint64
	readBlockPeriod    int
)

func init() {
	cfgMap, err := utils.ReadConfig(map[string]string{
		rpc:               "string",
		secret:            "string",
		eventEmitter:      "string",
		reader:            "string",
		dataStore:         "string",
		exchangeRouter:    "string",
		baseToken:         "string",
		daemon:            "int",
		last_block_number: "uint64",
		read_block_period: "int",
	})
	if err != nil {
		panic(err)
	}

	ethUtils.SetClient(cfgMap[rpc].(string))
	sKey, err = crypto.HexToECDSA(cfgMap[secret].(string))
	if err != nil {
		panic(err)
	}

	eventEmitterAddr = common.HexToAddress(cfgMap[eventEmitter].(string))
	readerAddr = common.HexToAddress(cfgMap[reader].(string))
	dataStoreAddr = common.HexToAddress(cfgMap[dataStore].(string))
	exchangeRouterAddr = common.HexToAddress(cfgMap[exchangeRouter].(string))
	baseTokenAddr = common.HexToAddress(cfgMap[baseToken].(string))
	daemonPort = cfgMap[daemon].(int)
	blockNumber = cfgMap[last_block_number].(uint64)
	readBlockPeriod = cfgMap[read_block_period].(int)

	log.Printf("config: %v\n", cfgMap)
}

func GetSKey() *ecdsa.PrivateKey {
	return sKey
}

func GetEventEmitterAddr() common.Address {
	return eventEmitterAddr
}

func GetReaderAddr() common.Address {
	return readerAddr
}

func GetDataStoreAddr() common.Address {
	return dataStoreAddr
}

func GetExchangeRouterAddr() common.Address {
	return exchangeRouterAddr
}

func GetBaseTokenAddr() common.Address {
	return baseTokenAddr
}

func GetDaemonPort() int {
	return daemonPort
}

func GetWorkDir() string {
	return workDir
}

func GetBlockNumber() uint64 {
	return blockNumber
}

func GetReadBlockPeriod() int {
	return readBlockPeriod
}
