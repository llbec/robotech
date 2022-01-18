package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/squadrons/aavesquadron/lendingpoolevent"
)

func handleBlock(from, to int64) {
	start := from
	end := start + maxEnd
	for {
		if end > to {
			end = to
		}
		log.Printf("block from %v to %v\n", start, end)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(contract)},
			FromBlock: big.NewInt(int64(start)),
			ToBlock:   big.NewInt(int64(end)),
		}
		logs, err := rpcClient.GetHTTPClient().FilterLogs(context.Background(), query)
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
}
