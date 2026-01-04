package syncer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func ScanLogs(
	ctx context.Context,
	client ethereum.LogFilterer,
	from, to uint64,
	addresses []common.Address,
	topics [][]common.Hash,
) ([]common.Hash, error) {

	logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: addresses,
		Topics:    topics,
	})
	if err != nil {
		return nil, err
	}

	txSet := make(map[common.Hash]struct{})
	for _, l := range logs {
		txSet[l.TxHash] = struct{}{}
	}

	var hashes []common.Hash
	for h := range txSet {
		hashes = append(hashes, h)
	}
	return hashes, nil
}
