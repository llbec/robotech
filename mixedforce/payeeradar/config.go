package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

const (
	SERVERIP = "server"
	PAYEE    = "target"
	TOKEN    = "token"
)

var (
	Payee    string
	Token    string
	Server   string
	CtrlChan chan int
	cfg      *viper.Viper
	targets  map[string]string
)

func LoadConfig(cfgpath string) {
	cfg = viper.New()
	cfg.AddConfigPath(cfgpath)
	cfg.SetConfigName("xland")
	cfg.SetConfigType("yaml")

	/*cfg.WatchConfig()
	cfg.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file:%s Op:%s\nreload file:%s\n", e.Name, e.Op, e.Name)
		readConfig()
		showCfg()
	})*/
	readConfig()
	showCfg()
}

func readConfig() {
	var lock sync.Mutex

	err := cfg.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal("ReadInConfig ", err)
	}

	f, err := os.Open(cfg.GetString(TOKEN))
	if err != nil {
		log.Panicf("Open token file: %v", err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Panicf("Read token file: %v", err)
	}

	lock.Lock()
	defer lock.Unlock()
	Server = cfg.GetString(SERVERIP)
	Token = string(content)

	targets = cfg.GetStringMapString(PAYEE)
}

func GeneratConfig(cfgpath string) error {
	viper.Set(SERVERIP, "192.168.11.51:1235")
	viper.Set(TOKEN, "wallet.token")
	ts := map[string]string{}
	ts["d0214E7B831E814a3151196e1c0818873486D4B"] = "183.238.69.213:10002"
	ts["E5C585eDE6ae6527Edf4cce5223715FCD5d76d07"] = "183.238.69.213:10002"
	viper.Set(PAYEE, ts)
	path := filepath.Join(cfgpath, "xland.yaml")
	return viper.WriteConfigAs(path)
}

func showCfg() {
	fmt.Printf("Running:\n\tfilecoin node Token:%v\n\tfilecoin node host:%v\n",
		Token,
		Server)
	for t, xl := range targets {
		fmt.Printf("\t%v @ %v\n", t, xl)
	}
}
