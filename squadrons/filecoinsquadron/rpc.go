package filecoinsquadron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Rpc struct {
	urlHost   string
	urlPath   string
	urlScheme string
	authToken string
}

//scheme: http or https
func NewRpc(host, path, scheme, token string) *Rpc {
	return &Rpc{
		urlHost:   host,
		urlPath:   path,
		urlScheme: scheme,
		authToken: token,
	}
}

/*
curl -X POST -H "Content-Type: application/json"
  --user <PROJECT_ID>:<PROJECT_SECRET>
  --url https://filecoin.infura.io
  --data '{ "id": 0, "jsonrpc": "2.0", "method": "Filecoin.ChainHead", "params": [] }'
*/

func (rpc *Rpc) Post(method string, params map[string]interface{}) ([]byte, error) {
	body := make(map[string]interface{})
	body["jsonrpc"] = "2.0"
	body["method"] = method
	body["id"] = 3
	for k, v := range params {
		body[k] = v
	}

	bytesData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bytesData)

	u := url.URL{Scheme: rpc.urlScheme, Host: rpc.urlHost, Path: rpc.urlPath}
	request, err := http.NewRequest("POST", u.String(), reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", rpc.authToken))
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
