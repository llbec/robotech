package syncer

import (
	"context"
	"time"

	"ether-indexer/internal/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FetchTransactions(
	ctx context.Context,
	client *ethclient.Client,
	projectID string,
	hashes []common.Hash,
) ([]model.Transaction, []model.TransactionLog, error) {

	var txs []model.Transaction
	var logs []model.TransactionLog

	for _, h := range hashes {

		tx, _, _ := client.TransactionByHash(ctx, h)
		rcpt, _ := client.TransactionReceipt(ctx, h)
		blk, _ := client.BlockByHash(ctx, rcpt.BlockHash)

		//msg, _ := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
		from, _ := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)

		txs = append(txs, model.Transaction{
			ProjectID:   projectID,
			TxHash:      h.Hex(),
			BlockNumber: rcpt.BlockNumber.Uint64(),
			BlockTime:   time.Unix(int64(blk.Time()), 0),
			TxIndex:     rcpt.TransactionIndex,
			FromAddress: from.Hex(),
			ToAddress:   addr(tx.To()),
			Value:       tx.Value().String(),
			Gas:         tx.Gas(),
			GasPrice:    gas(tx),
			Nonce:       tx.Nonce(),
			InputData:   tx.Data(),
			Status:      int(rcpt.Status),
			CreatedAt:   time.Now(),
		})

		for _, l := range rcpt.Logs {
			logs = append(logs, model.TransactionLog{
				ProjectID: projectID,
				TxHash:    h.Hex(),
				LogIndex:  l.Index,
				Address:   l.Address.Hex(),
				Topics:    topics(l.Topics),
				Data:      l.Data,
			})
		}
	}

	return txs, logs, nil
}

func addr(a *common.Address) string {
	if a == nil {
		return ""
	}
	return a.Hex()
}

func gas(tx *types.Transaction) string {
	if tx.GasPrice() != nil {
		return tx.GasPrice().String()
	}
	if tx.GasFeeCap() != nil {
		return tx.GasFeeCap().String()
	}
	return "0"
}

func topics(ts []common.Hash) string {
	s := "["
	for i, t := range ts {
		if i > 0 {
			s += ","
		}
		s += `"` + t.Hex() + `"`
	}
	return s + "]"
}
