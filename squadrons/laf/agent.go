package laf

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"robotech/armory/abilibs/erc20abi"
	"robotech/armory/abilibs/lafabi"
	"robotech/armory/abilibs/uniswapv2abi"
	"robotech/armory/txstore"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	LafABI      abi.ABI
	StakingABI  abi.ABI
	ReferralABI abi.ABI
	SwapABI     abi.ABI
	UsdtABI     abi.ABI
)

func init() {
	var err error
	LafABI, err = abi.JSON(strings.NewReader(lafabi.LAFABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse LAFABI: %v", err))
	}
	StakingABI, err = abi.JSON(strings.NewReader(lafabi.STAKINGABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse StakingABI: %v", err))
	}
	ReferralABI, err = abi.JSON(strings.NewReader(lafabi.REFERRALABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse ReferralABI: %v", err))
	}
	SwapABI, err = abi.JSON(strings.NewReader(uniswapv2abi.SWAPABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse SwapABI: %v", err))
	}
	UsdtABI, err = abi.JSON(strings.NewReader(erc20abi.ERC20ABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse UsdtABI: %v", err))
	}
}

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
		routeContract:    common.HexToAddress(cfg.RouteContract),
	}
}

// Filter transactions containing expected logs
func (agent *LafAgent) FilterTxs(fromBlock, toBlock uint64) (txs []common.Hash, err error) {
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

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{
			agent.lafContract,
			agent.stakingContract,
			agent.referralContract},
		Topics: [][]common.Hash{{
			LafABI.Events["Transfer"].ID,
			LafABI.Events["OwnershipTransferred"].ID,
			StakingABI.Events["OwnershipTransferred"].ID,
			ReferralABI.Events["SetOperator"].ID}},
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

// Parse transaction logs
func (agent *LafAgent) ParseTx(tx common.Hash) (
	txEvent txstore.TxEvent,
	err error) {
	if agent.client == nil {
		err = fmt.Errorf("client is nil")
		return
	}

	receipt, err := agent.client.TransactionReceipt(context.Background(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get transaction receipt: %v", err)
		return
	}
	block, err := agent.client.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		err = fmt.Errorf("failed to get block: %v", err)
		return
	}
	txEvent = txstore.TxEvent{
		BlockHeight: int64(receipt.BlockNumber.Int64()),
		BlockTime:   int64(block.Time()),
		TxIndex:     int64(receipt.TransactionIndex),
		TxHash:      receipt.TxHash.String(),
	}

	for _, log := range receipt.Logs {
		data := make(map[string]any)
		switch log.Address {
		case agent.lafContract:
			switch log.Topics[0].String() {
			case LafABI.Events["Transfer"].ID.String():
				err = LafABI.UnpackIntoMap(data, log.Topics[0].String(), log.Data)
				if err != nil {
					err = fmt.Errorf("failed to unpack log: %v", err)
					return
				}
			}
		}
	}
	return
}
