package api

import (
	"net/http"
	"strconv"
	"txstore/service"
)

func SplitArchiveMoveAPI(adminSvc *service.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fromBlock, _ := strconv.ParseInt(r.URL.Query().Get("from_block"), 10, 64)
		toBlock, _ := strconv.ParseInt(r.URL.Query().Get("to_block"), 10, 64)
		archivePath := r.URL.Query().Get("archive_path")
		if archivePath == "" {
			http.Error(w, "archive_path required", 400)
			return
		}

		if err := adminSvc.SplitMove(fromBlock, toBlock, archivePath); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("ok"))
	}
}
