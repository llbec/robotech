package syncer

import (
	"context"
	"database/sql"
	"log"
	"time"

	"ether-indexer/internal/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Scheduler struct {
	client *ethclient.Client
	txRepo *storage.TxRepository
	cpRepo *storage.CheckpointRepo
}

func NewScheduler(c *ethclient.Client, tx *storage.TxRepository, cp *storage.CheckpointRepo) *Scheduler {
	return &Scheduler{c, tx, cp}
}

func (s *Scheduler) Run(ctx context.Context, projectID string) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		if err := s.tick(ctx, projectID); err != nil {
			log.Println("sync error:", err)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context, projectID string) error {

	cp, err := s.cpRepo.Load(projectID)
	if err != nil {
		return err
	}

	head, _ := s.client.HeaderByNumber(ctx, nil)
	latest := head.Number.Uint64()

	from := cp.CurrentBlock
	to := from + cp.Step
	to = min(to, latest)

	if from > to {
		return nil
	}

	var addrs []common.Address
	for _, a := range cp.Addresses() {
		addrs = append(addrs, common.HexToAddress(a))
	}

	hashes, err := ScanLogs(ctx, s.client, from, to, addrs, cp.Topics())
	if err != nil {
		return err
	}

	txs, logs, err := FetchTransactions(ctx, s.client, projectID, hashes)
	if err != nil {
		return err
	}

	return s.txRepo.SaveBatch(
		from,
		txs,
		logs,
		func(tx *sql.Tx) error {
			return s.cpRepo.UpdateTx(tx, projectID, to+1)
		},
	)
}
