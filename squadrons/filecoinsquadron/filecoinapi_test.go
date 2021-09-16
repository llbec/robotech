package filecoinsquadron_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"unsafe"

	"github.com/llbec/robotech/squadrons/filecoinsquadron"
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

func getToken() string {
	f, err := os.Open("token")
	if err != nil {
		return ""
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(content)
}

func Test_BodyParse(t *testing.T) {
	rpc := filecoinsquadron.NewRpc("192.168.11.51:1235", "/rpc/v0", "http", getToken())
	filecoinAPI := filecoinsquadron.NewFileCoinAPI(rpc, nil)
	tipsetBuf, err := filecoinAPI.GetTipsetByHeight(1116849)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(tipsetBuf)
}
