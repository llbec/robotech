package filecoinsquadron_test

import (
	"encoding/json"
	"fmt"
	"math/big"
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

func Test_PayeeRadarInBlock(t *testing.T) {
	cfg, filecoinAPI := loadCfg()
	targets := cfg.GetStringMapString("target")
	payeeMap := make(map[string]bool)
	for t := range targets {
		payeeMap[t] = true
	}
	tipsetBuf, err := filecoinAPI.GetTipsetByHeight(cfg.GetInt64("height"))
	if err != nil {
		t.Fatal(err)
	}
	tipset, err := filecoinAPI.ReadTipSet(tipsetBuf)
	if err != nil {
		t.Fatal(fmt.Errorf("ReadTipSet Height(%v): %v", cfg.GetInt64("height"), err))
	}
	ts := make(map[string]filecoinsquadron.MsgInfo)
	for _, b := range tipset.Blocks {
		msgs, err := filecoinAPI.PayeeRadarInBlock(payeeMap, b)
		if err != nil {
			t.Fatal(fmt.Errorf("PayeeRadarInBlock Height(%v-%s): %v", cfg.GetInt64("height"), b, err))
		}
		for _, m := range msgs {
			ts[m.Cid] = m
			fmt.Printf("Target(%v@%v): %v@%v\n", b, cfg.GetInt64("height"), m.To, m.Cid)
		}
	}
	fmt.Println("result: ", ts)
}

func Test_TipsetParse(t *testing.T) {
	cfg, filecoinAPI := loadCfg()
	tipsetBuf, err := filecoinAPI.GetTipsetByHeight(cfg.GetInt64("height"))
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Print(string(tipsetBuf))
	tipset, err := filecoinAPI.ReadTipSet(tipsetBuf)
	if err != nil {
		t.Fatal(fmt.Errorf("ReadTipSet Height(%v): %v", cfg.GetInt64("height"), err))
	}
	//fmt.Print(tipset)
	targets := cfg.GetStringMapString("target")
	payeeMap := make(map[string]bool)
	for t := range targets {
		payeeMap[t] = true
	}
	//fmt.Println(payeeMap)
	for _, b := range tipset.Blocks {
		rpcBytes, err := filecoinAPI.ChainGetBlockMessages(b)
		if err != nil {
			t.Fatal("ChainGetBlockMessages ", err)
		}
		rpcJson, err := simplejson.NewJson(rpcBytes)
		if err != nil {
			t.Fatal("NewJson ", err)
		}
		blsTxsJson := rpcJson.GetPath("result", "BlsMessages")
		blsTxs, err := blsTxsJson.Array()
		if err != nil {
			t.Fatal(fmt.Errorf("BlsMessages array: %v", err))
		}
		for i := 0; i < len(blsTxs); i++ {
			to, err := blsTxsJson.GetIndex(i).GetPath("To").String()
			if err != nil {
				t.Fatal(fmt.Errorf("BlsMessage[%d] To: %v", i, err))
			}
			//fmt.Printf("block(%v): to(%v)\n", b, to)
			if payeeMap[to] {
				txhash, err := blsTxsJson.GetIndex(i).GetPath("CID", "/").String()
				if err != nil {
					t.Fatal(fmt.Errorf("BlsMessage[%d] cid: %v", i, err))
				}
				from, err := blsTxsJson.GetIndex(i).GetPath("From").String()
				if err != nil {
					t.Fatal(fmt.Errorf("BlsMessage[%d] From: %v", i, err))
				}
				value, err := blsTxsJson.GetIndex(i).GetPath("Value").Bytes()
				if err != nil {
					t.Fatal(fmt.Errorf("BlsMessage[%d] Value: %v", i, err))
				}
				v, err := strconv.ParseFloat(string(value), 64)
				if err != nil {
					t.Fatal(fmt.Errorf("BlsMessage[%d] ParseFloat: %v", i, err))
				}
				vBig := big.NewFloat(v)
				vBig.Quo(vBig, big.NewFloat(1e18))
				fmt.Printf("BlsMessages: %v\t%v: %v -> %v\n", txhash, vBig, from, to)
			}
		}

		//SecpkMessages
		secpkTxsJson := rpcJson.GetPath("result", "SecpkMessages")
		txs, err := secpkTxsJson.Array()
		if err != nil {
			t.Fatal(fmt.Errorf("secpkMessages array: %v", err))
		}

		for i := 0; i < len(txs); i++ {
			to, err := secpkTxsJson.GetIndex(i).GetPath("Message", "To").String()
			if err != nil {
				t.Fatal(fmt.Errorf("secpkMessage[%d] To: %v", i, err))
			}
			if payeeMap[to] {
				txhash, err := secpkTxsJson.GetIndex(i).GetPath("CID", "/").String()
				if err != nil {
					t.Fatal(fmt.Errorf("secpkMessage[%d] cid: %v", i, err))
				}
				value, err := secpkTxsJson.GetIndex(i).GetPath("Message", "Value").Bytes()
				if err != nil {
					t.Fatal(fmt.Errorf("secpkMessage[%d] Value: %v", i, err))
				}
				v, err := strconv.ParseFloat(string(value), 64)
				if err != nil {
					t.Fatal(fmt.Errorf("secpkMessage[%d] ParseFloat: %v", i, err))
				}
				vBig := big.NewFloat(v)
				vBig.Quo(vBig, big.NewFloat(1e18))
				from, err := secpkTxsJson.GetIndex(i).GetPath("Message", "From").String()
				if err != nil {
					t.Fatal(fmt.Errorf("secpkMessage[%d] From: %v", i, err))
				}
				fmt.Printf("SecpkMessages: %v\t%v: %v -> %v\n", txhash, vBig, from, to)
			}
		}
	}
}

func Test_BlockMessage(t *testing.T) {
	cfg, api := loadCfg()

	rpcBytes, err := api.ChainGetBlockMessages(cfg.GetString("block"))
	if err != nil {
		t.Fatal("ChainGetBlockMessages ", err)
	}

	//fmt.Println(string(rpcBytes))

	rpcJson, err := simplejson.NewJson(rpcBytes)
	if err != nil {
		t.Fatal("NewJson ", err)
	}

	//fmt.Println("json ", rpcJson)

	blsTxsJson := rpcJson.GetPath("result", "BlsMessages")
	blsTxs, err := blsTxsJson.Array()
	if err != nil {
		t.Fatal(fmt.Errorf("BlsMessages array: %v", err))
	}
	for i := 0; i < len(blsTxs); i++ {
		to, err := blsTxsJson.GetIndex(i).GetPath("To").String()
		if err != nil {
			t.Fatal(fmt.Errorf("BlsMessage[%d] To: %v", i, err))
		}
		fmt.Println(to)
	}
}
