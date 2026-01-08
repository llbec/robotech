package syncer

import (
	"context"
	"ether-indexer/internal/model"
	"ether-indexer/internal/storage"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	mapProjectIDToScheduler = map[string]*Scheduler{}
)

func init() {
	mapProjectIDToScheduler = make(map[string]*Scheduler)
}

func GetScheduler(projectID string) *Scheduler {
	return mapProjectIDToScheduler[projectID]
}

func RegisterScheduler(p *model.Project) error {
	basePath := filepath.Join(p.BasePath, p.ProjectID)

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return err
	}

	shards := storage.NewShardManager(storage.ShardPolicy{
		BlockRange: p.BlockRange,
		BasePath:   basePath,
	})

	metaDB, err := storage.OpenDB(filepath.Join(p.BasePath, "checkpoints.db"))
	if err != nil {
		return err
	}
	if err = storage.Migrate(metaDB); err != nil {
		return err
	}

	txRepo := storage.NewTxRepository(shards)
	cpRepo := storage.NewCheckpointRepo(metaDB)

	client, err := ethclient.Dial(p.RPCEndpoint)
	if err != nil {
		return err
	}

	s := NewScheduler(client, txRepo, cpRepo)
	if p.Active {
		s.Run(context.Background(), p.ProjectID)
	}

	mapProjectIDToScheduler[p.ProjectID] = s
	return nil
}
