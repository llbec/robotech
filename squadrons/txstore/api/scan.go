package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"txstore/service"
)

func Scan(svc *service.TxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		block, _ := strconv.ParseInt(r.URL.Query().Get("block"), 10, 64)
		index, _ := strconv.ParseInt(r.URL.Query().Get("index"), 10, 64)
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		txs, err := svc.Scan(project, block, index, limit)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(txs)
	}
}
