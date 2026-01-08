package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ether-indexer/internal/storage"
	"ether-indexer/internal/syncer"
)

func main() {

	var dataPath string
	flag.StringVar(&dataPath, "p", ".data", "project db directory")
	flag.Parse()

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		log.Fatal("create data directory failed:", err)
	}

	dbPath := filepath.Join(dataPath, "indexer.db")
	fmt.Println("Using DB:", dbPath)

	proDB, _ := storage.OpenDB(dbPath)
	defer func() {
		if err := proDB.Close(); err != nil {
			log.Printf("close db error: %v", err)
		}
	}()
	_ = storage.ProjectMigrate(proDB)

	repo := storage.NewProjectRepo(proDB, dataPath)
	projects, err := repo.ListAll()
	if err != nil {
		log.Fatal("list projects failed:", err)
	}
	for _, p := range projects {
		//fmt.Println(p)
		syncer.RegisterScheduler(&p)
	}
}
