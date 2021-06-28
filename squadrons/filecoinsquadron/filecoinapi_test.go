package filecoinsquadron_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"unsafe"
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
