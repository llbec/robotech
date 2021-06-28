package main

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	HTTPHOST  = "http"
	WSHOST    = "ws"
	TROLLER   = "troller"
	FLASHLOAN = "flashloan"
	PATH      = "path"
	PROGRAMDB = "programdb"
	DEBTORDB  = "debtordb"
)

var (
	rpc          string
	wsHost       string
	compTroller  string
	flashContact string
	dbPath       string
	proDB        string
	debtorDB     string
	cfg          *viper.Viper
	cfgLock      sync.Mutex
)

func LoadConfig(cfgpath string) {
	cfg = viper.New()
	cfg.AddConfigPath(cfgpath)
	cfg.SetConfigName("compliquidation")
	cfg.SetConfigType("yaml")

	cfg.WatchConfig()
	cfg.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file:%s Op:%s\nreload file:%s\n", e.Name, e.Op, e.Name)
		readConfig()
		showCfg()
	})
	readConfig()
	showCfg()
}

func readConfig() {
	err := cfg.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(err)
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	rpc = cfg.GetString(HTTPHOST)
	wsHost = cfg.GetString(WSHOST)
	compTroller = cfg.GetString(TROLLER)
	flashContact = cfg.GetString(FLASHLOAN)
	dbPath = cfg.GetString(PATH)
	proDB = cfg.GetString(PROGRAMDB)
	debtorDB = cfg.GetString(DEBTORDB)
}

func GeneratConfig(cfgpath string) error {
	viper.Set(HTTPHOST, "https://http-mainnet-node.huobichain.com")
	viper.Set(WSHOST, "wss://heco.apron.network/ws/v1/39407a13d0e54fdc99dde78ca1c41900")
	viper.Set(TROLLER, "0xb74633f2022452f377403B638167b0A135DB096d")
	viper.Set(FLASHLOAN, "0x6524a87301129E3419ad8E1eC8d6A5Ef133d59B2")
	viper.Set(PATH, ".")
	viper.Set(PROGRAMDB, "comp.db")
	viper.Set(DEBTORDB, "debtor.db")
	path := filepath.Join(cfgpath, "compliquidation.yaml")
	return viper.WriteConfigAs(path)
}

func showCfg() {
	fmt.Printf("Running:\n\thttp host: %v\n\twebsocket host: %v\n\tCompTroller: %v\n\tflashloan: %v\n",
		rpc, wsHost, compTroller, flashContact)
}
