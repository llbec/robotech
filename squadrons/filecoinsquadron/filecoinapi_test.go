package filecoinsquadron_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"unsafe"

	"github.com/bitly/go-simplejson"
	"github.com/llbec/robotech/squadrons/filecoinsquadron"
	"github.com/spf13/viper"
)

func Test_paramsforChainGetTipSetByHeight(t *testing.T) {
	height := 666666
	body := `
			{ "jsonrpc": "2.0", "method":"Filecoin.ChainGetTipSetByHeight", "params": [` + strconv.Itoa(height) + `, null], "id": 3 }
		`
	fmt.Println(body)

	bodyMap := make(map[string]interface{})
	bodyMap["jsonrpc"] = "2.0"
	bodyMap["method"] = "Filecoin.ChainGetTipSetByHeight"
	bodyMap["params"] = []interface{}{height, nil}
	bodyMap["id"] = 3
	bytesData, err := json.Marshal(bodyMap)
	if err != nil {
		t.Fatal(err)
	}

	str := (*string)(unsafe.Pointer(&bytesData))
	fmt.Println("rpc params:\n", *str)
}

func loadCfg() (cfg *viper.Viper, filAPI *filecoinsquadron.FileCoinAPI) {
	cfg = viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigName("test")
	cfg.SetConfigType("yaml")
	if err := cfg.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(err)
	}
	rpc := filecoinsquadron.NewRpc(cfg.GetString("nodeserver"), "/rpc/v0", "http", cfg.GetString("nodetoken"))
	filAPI = filecoinsquadron.NewFileCoinAPI(rpc, nil)
	return
}

func Test_BodyParse(t *testing.T) {
	cfg, filecoinAPI := loadCfg()
	tipsetBuf, err := filecoinAPI.GetTipsetByHeight(cfg.GetInt64("height"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(tipsetBuf))
}

func Test_BlockMessage(t *testing.T) {
	cfg, api := loadCfg()

	rpcBytes, err := api.ChainGetBlockMessages(cfg.GetString("block"))
	if err != nil {
		t.Fatal("ChainGetBlockMessages ", err)
	}

	fmt.Println(string(rpcBytes))

	rpcJson, err := simplejson.NewJson(rpcBytes)
	if err != nil {
		t.Fatal("NewJson ", err)
	}

	fmt.Println("json ", rpcJson)
}
