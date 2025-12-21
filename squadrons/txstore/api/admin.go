package api

import (
	"net/http"
	"txstore/service"
	"txstore/storage/sqlite"
)

func AddArchive(dbm *service.DBManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		st, err := sqlite.Open(path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		dbm.AddArchive(st)
		w.Write([]byte("ok"))
	}
}
