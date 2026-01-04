package main

import (
	"context"

	"ether-indexer/internal/storage"
	"ether-indexer/internal/syncer"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

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
