package laf

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewLAFAgent(cfgFile string) *LafAgent {
	cfg := &LafAgentConfig{}
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Printf("read config file %s failed: %v", cfgFile, err)
		return nil
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Printf("unmarshal config file %s failed: %v", cfgFile, err)
		return nil
	}
	client, err := ethclient.Dial(cfg.RpcUrl)
	if err != nil {
		log.Printf("dial rpc url %s failed: %v", cfg.RpcUrl, err)
		return nil
	}
	return &LafAgent{
		client:           client,
		lafContract:      common.HexToAddress(cfg.LafContract),
		stakingContract:  common.HexToAddress(cfg.StakingContract),
		referralContract: common.HexToAddress(cfg.ReferralContract),
		usdtContract:     common.HexToAddress(cfg.USDTContract),
		swapContract:     common.HexToAddress(cfg.SwapContract),
	}
}

// Filter transactions containing logs
func (agent *LafAgent) FilterLogs(fromBlock, toBlock uint64) (txs []common.Hash, err error) {
	if fromBlock > toBlock {
		err = fmt.Errorf("fromBlock %d is greater than toBlock %d", fromBlock, toBlock)
		return
	}
	if agent.client == nil {
		err = fmt.Errorf("client is nil")
		return
	}
	if toBlock-fromBlock > 9 {
		err = fmt.Errorf("block range %d is too large", toBlock-fromBlock)
		return
	}
	lafABI, err := abi.JSON(strings.NewReader(LAFABI))
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{agent.lafContract},
		Topics:    [][]common.Hash{{lafABI.Events["Transfer"].ID}},
	}

	logs, err := agent.client.FilterLogs(context.Background(), query)
	if err != nil {
		err = fmt.Errorf("failed to filter logs[%v-%v]: %v", fromBlock, toBlock, err)
		return
	}
	txRecords := make(map[common.Hash]bool)
	for _, v := range logs {
		if txRecords[v.TxHash] {
			continue
		}
		txRecords[v.TxHash] = true
		txs = append(txs, v.TxHash)
	}
	return
}
