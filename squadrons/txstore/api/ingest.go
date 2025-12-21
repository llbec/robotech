package api

import (
	"encoding/json"
	"net/http"
	"txstore/model"
	"txstore/service"
)

func BatchInsert(svc *service.TxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		var txs []*model.Transaction
		if err := json.NewDecoder(r.Body).Decode(&txs); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := svc.Insert(project, txs); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("ok"))
	}
}
