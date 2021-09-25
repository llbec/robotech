package main

import (
	"context"
	"log"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/logistics/db"
	"github.com/llbec/robotech/squadrons/aavesquadron/lendingpoolevent"
	"github.com/llbec/robotech/utils"
)

const (
	actDeposit  = 1
	actWithdraw = 2
	actBorrow   = 3
	actRepay    = 4
	maxEnd      = 1000
)

var (
	usrDb   db.DB
	logDb   db.DB
	debtors map[string]*big.Int
)

func initAccounts() {
	usrDb = db.NewDB(
		db.LevelDBBackend,
		"usrs.db",
		filepath.Join(workDir, "database"),
	)

	logDb = db.NewDB(
		db.LevelDBBackend,
		"log.db",
		filepath.Join(workDir, "database"),
	)

	debtors = make(map[string]*big.Int)

	lendingpoolevent.RegisterHandle(lendingpoolevent.BorrowEvent, logBorrow)
}

func recoveryDB() uint64 {
	iter := usrDb.NewIterator()
	for {
		addr := common.BytesToAddress(iter.Key())
		debtors[addr.String()] = big.NewInt(int64(utils.BytesToUInt64(iter.Value())))

		if !iter.Next() {
			break
		}
	}
	h, err := logDb.Get([]byte("height"))
	if err != nil {
		log.Printf("load height failed: %v", err)
	}
	return utils.BytesToUInt64(h)
}

func syncBlock(from, to int64) {
	log.Printf("block sync task from %v to %v, step %v", from, to, maxEnd)
	start := from
	end := start + maxEnd
	for {
		if end > to {
			end = to
		}
		log.Printf("sync block from %v to %v\n", start, end)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(LendingPoolAddress)},
			FromBlock: big.NewInt(int64(start)),
			ToBlock:   big.NewInt(int64(end)),
		}
		logs, err := Client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("FilterLogs(%v-%v): %v\n", start, end, err)
		}
		for _, vLog := range logs {
			err = lendingpoolevent.LendingPoolEventHandle(vLog)
			if err != nil {
				log.Printf("HandleLog(%v-%v): %v\n", vLog.BlockNumber, vLog.Index, err)
			}
		}
		if end == to {
			break
		}
		start = end + 1
		end = start + maxEnd
	}
	logDb.Put([]byte("height"), utils.UInt64ToBytes(uint64(end)))
}

func logBorrow(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.BorrowEventData)
	if !ok {
		log.Printf("Parse BorrowEventData failed!\n")
		return
	}
	_, ok = debtors[data.User.String()]
	if !ok {
		usrDb.Put(data.User[:], utils.UInt64ToBytes(data.Amount.Uint64()))
	}
	act, err := LendingPool.GetUserAccountData(nil, data.User)
	if err != nil {
		log.Printf("logBorrow: GetUserAccountData error: %v\n", err)
		return
	}
	debtors[data.User.String()] = act.TotalDebtETH
}
