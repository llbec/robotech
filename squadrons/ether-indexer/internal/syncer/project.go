package syncer

import (
	"context"
	"ether-indexer/internal/model"
	"ether-indexer/internal/storage"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	mapProjectIDToScheduler = map[string]*Scheduler{}
)

func GetScheduler(projectID string) *Scheduler {
	return mapProjectIDToScheduler[projectID]
}

func RegisterScheduler(p *model.Project) {
	shards := storage.NewShardManager(storage.ShardPolicy{
		BlockRange: p.BlockRange,
		BasePath:   p.BasePath,
	})

	metaDB, _ := storage.OpenDB(p.BasePath + "/" + p.ProjectID + ".db")
	_ = storage.Migrate(metaDB)

	txRepo := storage.NewTxRepository(shards)
	cpRepo := storage.NewCheckpointRepo(metaDB)

	client, _ := ethclient.Dial(p.RPCEndpoint)

	s := NewScheduler(client, txRepo, cpRepo)
	s.Run(context.Background(), p.ProjectID)

	mapProjectIDToScheduler[p.ProjectID] = s
}
