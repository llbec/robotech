package main

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	lastupdate  = int64(0)
	assets      []AssetFeeds
	cfgfile     *viper.Viper
	cfgLock     sync.Mutex
	cfg         *config
	gCfgUpdated = false
)

type config struct {
	assetscontract    map[string]string
	assetschanlinknet map[string]string
	rpcs              map[string]string
	interval          int64
	threhold          int
	price             string
	testnet           bool
	fallbackaddr      string
	fallbacknet       string
	tickime           int64
}

const (
	ASSETS     = "assets"
	NETS       = "nets"
	NETDEFAULT = "netdefault"
	RPCS       = "rpcs"
	INTERVAL   = "interval"
	THREHOLD   = "threhold"
	PRICE      = "price"
	TESTNET    = "testnet"
	FBADDR     = "fallbackcontract"
	FBNET      = "fallbacknet"
	TICK       = "ticktime"
)

func GeneratConfig(dir string) error {
	viper.Set(ASSETS, map[string]string{
		"USDT": "0xD16bAbe52980554520F6Da505dF4d1b124c815a7",
		"USDC": "0x92a0bD4584c147D1B0e8F9185dB0BDa10B05Ed7e",
		"BTC":  "0xAad9654a4df6973A92C1fd3e95281F0B37960CCd",
		"BCH":  "0x6705b383dEF2D9f3f93bc00C5FDe402613d2D695",
		"FIL":  "0x0bf85D3B0C9ebCC282FDe0591882d12E57E700B3",
		"ETH":  "0xA1588dC914e236bB5AE4208Ce3081246f7A00193",
		"LTC":  "0x13e93721DC992b3E14333dBdb48C0e7Ec55431c3",
		"DOT":  "0xA71e7ae7A5B154f1ED12476a4C54c5Ec6e3426AC"})
	viper.Set(NETS, map[string]string{"USDT": "heco", "USDC": "bsc"})
	viper.Set(NETDEFAULT, "heco")
	viper.Set(RPCS, map[string]string{
		"heco": "https://http-mainnet-node.huobichain.com",
		"bsc":  "https://bsc-dataseed1.defibit.io/",
		"hsc":  "https://http-mainnet.hoosmartchain.com"})
	viper.Set(INTERVAL, int64(1800))
	viper.Set(THREHOLD, 10)
	viper.Set(PRICE, "USD")
	viper.Set(TESTNET, false)
	viper.Set(FBADDR, "0x2C2882158DC6774BF0cF2168065494A70135D40d")
	viper.Set(FBNET, "hsc")
	viper.Set(TICK, int64(300))

	path := filepath.Join(dir, "config.yaml")
	return viper.WriteConfigAs(path)
}

func LoadConfig(dir string) error {
	cfgfile = viper.New()
	cfgfile.AddConfigPath(dir)
	cfgfile.SetConfigName("config")
	cfgfile.SetConfigType("yaml")

	cfgfile.WatchConfig()
	cfgfile.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file:%s Op:%s\nreload file:%s\n", e.Name, e.Op, e.Name)
		readConfig()
		//showCfg()
	})
	return readConfig()
}

func readConfig() error {
	err := cfgfile.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		return err
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()

	cfg = &config{
		assetscontract:    map[string]string{},
		assetschanlinknet: map[string]string{},
		rpcs:              map[string]string{},
	}

	cfg.assetscontract = cfgfile.GetStringMapString(ASSETS)
	netdefault := cfgfile.GetString(NETDEFAULT)
	nets := cfgfile.GetStringMapString(NETS)
	for a := range cfg.assetscontract {
		if n, ok := nets[a]; ok {
			cfg.assetschanlinknet[a] = n
		} else {
			cfg.assetschanlinknet[a] = netdefault
		}
	}
	cfg.rpcs = cfgfile.GetStringMapString(RPCS)
	cfg.fallbackaddr = cfgfile.GetString(FBADDR)
	cfg.interval = cfgfile.GetInt64(INTERVAL)
	cfg.price = cfgfile.GetString(PRICE)
	cfg.testnet = cfgfile.GetBool(TESTNET)
	cfg.threhold = cfgfile.GetInt(THREHOLD)
	cfg.fallbacknet = cfgfile.GetString(FBNET)
	cfg.tickime = cfgfile.GetInt64(TICK)

	gCfgUpdated = true

	return nil
}
