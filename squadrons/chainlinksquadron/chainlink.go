package chainlinksquadron

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/chainlink"
	"github.com/llbec/robotech/squadrons/chainlinksquadron/pricefeeds"
)

type ChainNode struct {
	Rpc map[string]*ethclient.Client
}

func NewChainNode() *ChainNode {
	return &ChainNode{
		Rpc: make(map[string]*ethclient.Client),
	}
}

func (node *ChainNode) AddNode(net, url string) error {
	rpcclient, err := ethclient.Dial(url)
	if err != nil {
		return err
	}
	node.Rpc[net] = rpcclient
	return nil
}

type ChainLinkSquadron struct {
	rpcs   *ChainNode
	pFeeds *pricefeeds.ChainLinkFeeds
}

func NewChainLinkSquadron(rpc *ChainNode, path string) *ChainLinkSquadron {
	clSquadron := &ChainLinkSquadron{}

	clSquadron.rpcs = rpc
	clSquadron.pFeeds = pricefeeds.NewChainLinkFeeds(path)

	return clSquadron
}

func (clSquadron *ChainLinkSquadron) GetLastPrice(chain, asset, price string, testnet bool) (*big.Int, int, error) {
	addr, decimal, err := clSquadron.pFeeds.GetContract(chain, testnet, asset, price)
	if err != nil {
		return big.NewInt(0), 0, err
	}
	contract, err := chainlink.NewEACAggregatorProxy(common.HexToAddress(addr), clSquadron.rpcs.Rpc[chain])
	if err != nil {
		return big.NewInt(0), 0, fmt.Errorf("conect to contract failed %v", err)
	}
	wrt, err := contract.LatestAnswer(nil)
	if err != nil {
		return big.NewInt(0), 0, fmt.Errorf("LatestAnswer failed %v", err)
	}
	str, _ := contract.Description(nil)
	if str != fmt.Sprintf("%s / %s", asset, price) {
		return big.NewInt(0), 0, fmt.Errorf("contract expect%s / %s, get %s", asset, price, str)
	}
	return wrt, decimal, nil
}

func (clSquadron *ChainLinkSquadron) GetAssets(chain, price string, testnet bool) []string {
	return clSquadron.pFeeds.GetAssets(chain, price, testnet)
}
