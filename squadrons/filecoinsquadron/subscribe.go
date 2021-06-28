package filecoinsquadron

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Subscribe struct {
	urlHost   string
	urlPath   string
	urlScheme string
	authToken string
}

// scheme: ws or wss
func NewSubscribe(host, path, scheme, token string) *Subscribe {
	return &Subscribe{
		urlHost:   host,
		urlPath:   path,
		urlScheme: scheme,
		authToken: token,
	}
}

func (subs *Subscribe) Connect() (*websocket.Conn, error) {
	head := make(http.Header)
	head.Set("Content-Type", "application/json;charset=UTF-8")
	head.Set("Authorization", fmt.Sprintf("Bearer %v", subs.authToken))

	u := url.URL{Scheme: subs.urlScheme, Host: subs.urlHost, Path: subs.urlPath}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), head)
	if err != nil {
		return nil, err
	}
	return c, nil
}
