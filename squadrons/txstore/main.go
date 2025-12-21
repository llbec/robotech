package main

import (
	"log"
	"net/http"
	"txstore/api"
	"txstore/service"
	"txstore/storage/sqlite"
)

func main() {
	active, err := sqlite.Open("./active.db")
	if err != nil {
		log.Fatal(err)
	}

	dbm := service.NewDBManager(active)
	txSvc := service.NewTxService(dbm)
	consumerSvc := service.NewConsumerService(txSvc)
	adminSvc := service.NewAdminService(dbm)

	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/tx/batch_insert", api.BatchInsert(txSvc))
	http.HandleFunc("/tx/scan", api.Scan(txSvc))
	http.HandleFunc("/consumer/scan", api.ConsumerScan(consumerSvc))
	http.HandleFunc("/consumer/commit", api.ConsumerCommit(consumerSvc))
	http.HandleFunc("/admin/add_archive", api.AddArchive(dbm))
	http.HandleFunc("/admin/split_move", api.SplitArchiveMoveAPI(adminSvc))

	log.Println("txstore listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
