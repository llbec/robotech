package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"ether-indexer/internal/storage"
	"ether-indexer/internal/syncer"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	var dataPath string
	flag.StringVar(&dataPath, "p", "data", "project db directory")
	flag.Parse()

	dbPath := filepath.Join(dataPath, "indexer.db")
	fmt.Println("Using DB:", dbPath)

	proDB, _ := storage.OpenDB(dbPath)
	defer proDB.Close()
	_ = storage.ProjectMigrate(proDB)

	repo := storage.NewProjectRepo(proDB)
	projects, err := repo.ListAll()
	if err != nil {
		log.Fatal("list projects failed:", err)
	}
	for _, p := range projects {
		fmt.Println(p)
	}

	shards := storage.NewShardManager(storage.ShardPolicy{
		BlockRange: 5_000_000,
		BasePath:   "./data",
	})

	metaDB, _ := storage.OpenDB("./meta.db")
	_ = storage.Migrate(metaDB)

	txRepo := storage.NewTxRepository(shards)
	cpRepo := storage.NewCheckpointRepo(metaDB)

	client, _ := ethclient.Dial("https://YOUR_RPC")

	s := syncer.NewScheduler(client, txRepo, cpRepo)
	s.Run(context.Background(), "demo-project")
}
