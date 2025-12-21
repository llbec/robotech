package api

import (
	"encoding/json"
	"net/http"
	"txstore/model"
	"txstore/service"
)

func ConsumerScan(svc *service.ConsumerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		consumer := r.URL.Query().Get("consumer")
		limit := 100

		txs, cp, err := svc.Scan(project, consumer, limit)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"checkpoint": cp,
			"txs":        txs,
		})
	}
}

func ConsumerCommit(svc *service.ConsumerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cp model.Checkpoint
		json.NewDecoder(r.Body).Decode(&cp)
		if err := svc.Commit(&cp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("ok"))
	}
}
