package filecoinsquadron

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

type TipSet struct {
	Height    int64
	Timestamp int64
	Blocks    []string
}

type MsgInfo struct {
	Cid   string
	From  string
	To    string
	Value *big.Float
}

type FileCoinAPI struct {
	rpc    *Rpc
	subs   *Subscribe
	logger log.Logger
}

func NewFileCoinAPI(rpcnode *Rpc, subsnode *Subscribe) *FileCoinAPI {
	return &FileCoinAPI{
		rpc:  rpcnode,
		subs: subsnode,
	}
}

func (api *FileCoinAPI) GetTipsetByHeight(height int64) ([]byte, error) {
	//{ "jsonrpc": "2.0", "method":"Filecoin.ChainGetTipSetByHeight", "params": [height, null], "id": 3 }
	params := make(map[string]interface{})
	params["params"] = []interface{}{height, nil}
	return api.rpc.Post("Filecoin.ChainGetTipSetByHeight", params)
}

func (api *FileCoinAPI) ChainGetBlockMessages(blockid string) ([]byte, error) {
	//{ "jsonrpc": "2.0", "method":"Filecoin.ChainGetBlockMessages", "params": [{"/":"blockid"}], "id": 3 }
	params := make(map[string]interface{})
	cid := make(map[string]string)
	cid["/"] = blockid
	params["params"] = []map[string]string{cid}
	return api.rpc.Post("Filecoin.ChainGetBlockMessages", params)
}

func (api *FileCoinAPI) ChainNotify() (*websocket.Conn, error) {
	body := make(map[string]interface{})
	body["jsonrpc"] = "2.0"
	body["method"] = "Filecoin.ChainNotify"
	body["id"] = 3
	bytesBody, err := json.Marshal(body)
	if err != nil {
		api.logger.Panic(err)
	}

	c, err := api.subs.Connect()
	if err != nil {
		return nil, err
	}

	err = c.WriteMessage(websocket.BinaryMessage, bytesBody)
	if err != nil {
		return nil, err
	}

	// first response is {"jsonrpc":"2.0","result":1,"id":3}
	_, resp, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	respJson, err := simplejson.NewJson(resp)
	if err != nil {
		api.logger.Panic(err)
	}
	r, ok := respJson.CheckGet("result")
	i, _ := r.Int()
	if !ok || i != 1 {
		api.logger.Panic(err)
	}

	return c, nil
}

func (api *FileCoinAPI) ReadHeightFromTipSet(rpcBytes []byte) (int64, error) {
	rpcJson, err := simplejson.NewJson(rpcBytes)
	if err != nil {
		return 0, fmt.Errorf("simplejson new: %v", err)
	}
	tipSetJson := rpcJson.Get("params").GetIndex(1).GetIndex(0)
	ts := tipSetJson.GetPath("Val", "Blocks").GetIndex(0).Get("Height")
	height, err := ts.Int64()
	if err != nil {
		return 0, fmt.Errorf("timstamp: %v", err)
	}
	return height, nil
}

func (api *FileCoinAPI) ReadTipSet(rpcBytes []byte) (*TipSet, error) {
	tipSetJson, err := simplejson.NewJson(rpcBytes)
	if err != nil {
		return nil, err
	}

	tipSet := &TipSet{}

	b0 := tipSetJson.GetPath("result", "Blocks").GetIndex(0)
	ts := b0.Get("Timestamp")
	tipSet.Timestamp, err = ts.Int64()
	if err != nil {
		return nil, err
	}
	hs := b0.Get("Height")
	tipSet.Height, err = hs.Int64()
	if err != nil {
		return nil, err
	}

	cidsJson := tipSetJson.GetPath("result", "Cids")
	cidArray, err := cidsJson.Array()
	if err != nil {
		log.Fatalf("read Cids: %v\n", err)
		return nil, err
	}
	for i := range cidArray {
		cidstr, err := cidsJson.GetIndex(i).Get("/").String()
		if err != nil {
			return nil, err
		}
		tipSet.Blocks = append(tipSet.Blocks, cidstr)
	}

	return tipSet, nil
}

func (api *FileCoinAPI) ReadBlock(blockid string) ([]MsgInfo, error) {
	targets := []MsgInfo{}

	rpcBytes, err := api.ChainGetBlockMessages(blockid)
	if err != nil {
		return nil, err
	}

	rpcJson, err := simplejson.NewJson(rpcBytes)
	if err != nil {
		return nil, fmt.Errorf("simplejson new: %v", err)
	}
	txsJson := rpcJson.GetPath("result", "SecpkMessages")
	txs, err := txsJson.Array()
	if err != nil {
		return nil, fmt.Errorf("secpkMessages array: %v", err)
	}

	for i := 0; i < len(txs); i++ {
		to, err := txsJson.GetIndex(i).GetPath("Message", "To").String()
		if err != nil {
			return nil, fmt.Errorf("message[%d] To: %v", i, err)
		}
		txhash, err := txsJson.GetIndex(i).GetPath("CID", "/").String()
		if err != nil {
			return nil, fmt.Errorf("message[%d] cid: %v", i, err)
		}
		value, err := txsJson.GetIndex(i).GetPath("Message", "Value").Bytes()
		if err != nil {
			return nil, fmt.Errorf("message[%d] Value: %v", i, err)
		}
		v, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return nil, fmt.Errorf("message[%d] ParseFloat: %v", i, err)
		}
		vBig := big.NewFloat(v)
		vBig.Quo(vBig, big.NewFloat(1e18))
		from, err := txsJson.GetIndex(i).GetPath("Message", "From").String()
		if err != nil {
			return nil, fmt.Errorf("message[%d] From: %v", i, err)
		}
		targets = append(targets, MsgInfo{Cid: txhash, From: from, To: to, Value: vBig})
	}

	return targets, nil
}

func (api *FileCoinAPI) PayeeRadarInBlock(payee string, blockid string) ([]MsgInfo, error) {
	targets := []MsgInfo{}

	rpcBytes, err := api.ChainGetBlockMessages(blockid)
	if err != nil {
		return nil, err
	}

	rpcJson, err := simplejson.NewJson(rpcBytes)
	if err != nil {
		return nil, fmt.Errorf("simplejson new: %v", err)
	}
	txsJson := rpcJson.GetPath("result", "SecpkMessages")
	txs, err := txsJson.Array()
	if err != nil {
		return nil, fmt.Errorf("secpkMessages array: %v", err)
	}

	for i := 0; i < len(txs); i++ {
		to, err := txsJson.GetIndex(i).GetPath("Message", "To").String()
		if err != nil {
			return nil, fmt.Errorf("message[%d] To: %v", i, err)
		}
		if to == payee {
			txhash, err := txsJson.GetIndex(i).GetPath("CID", "/").String()
			if err != nil {
				return nil, fmt.Errorf("message[%d] cid: %v", i, err)
			}
			value, err := txsJson.GetIndex(i).GetPath("Message", "Value").Bytes()
			if err != nil {
				return nil, fmt.Errorf("message[%d] Value: %v", i, err)
			}
			v, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return nil, fmt.Errorf("message[%d] ParseFloat: %v", i, err)
			}
			vBig := big.NewFloat(v)
			vBig.Quo(vBig, big.NewFloat(1e18))
			from, err := txsJson.GetIndex(i).GetPath("Message", "From").String()
			if err != nil {
				return nil, fmt.Errorf("message[%d] From: %v", i, err)
			}
			targets = append(targets, MsgInfo{Cid: txhash, From: from, To: to, Value: vBig})
		}
	}

	return targets, nil
}
