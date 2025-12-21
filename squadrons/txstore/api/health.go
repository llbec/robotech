package api

import "net/http"

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}
