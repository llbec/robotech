package xlandteam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bitly/go-simplejson"
)

type XlandTeam struct {
	host string
}

func NewXlandTeam(h string) *XlandTeam {
	return &XlandTeam{host: h}
}

func (xl *XlandTeam) XlandSaveValue(timstamp int64, value string) error {
	body := make(map[string]interface{})
	body["time"] = timstamp
	body["value"] = value
	bytesData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json Marshal: %v", err)
	}
	reader := bytes.NewReader(bytesData)

	u := url.URL{Scheme: "http", Host: xl.host, Path: "/common/saveValue"}
	request, err := http.NewRequest("POST", u.String(), reader)
	if err != nil {
		return fmt.Errorf("http NewRequest: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("IO ReadAll: %v", err)
	}
	retJson, err := simplejson.NewJson(respBytes)
	if err != nil {
		return fmt.Errorf("simplejson new: %v", err)
	}
	msg, err := retJson.Get("message").String()
	if err != nil || msg != "ok" {
		if err != nil {
			return fmt.Errorf("response message: %v", err)
		}
		return fmt.Errorf("response invalid message: %v", msg)
	}
	return nil
}
