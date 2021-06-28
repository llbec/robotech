package pricefeeds

import "github.com/llbec/robotech/utils"

const (
	BSC  string = "https://docs.chain.link/docs/binance-smart-chain-addresses/"
	ETH  string = "https://docs.chain.link/docs/ethereum-addresses/"
	HECO string = "https://docs.chain.link/docs/huobi-eco-chain-price-feeds/"
)

type Feeds struct {
	token    string
	decimal  int
	contract string
}

type PriceFeeds struct {
	net   *utils.Chain
	feeds map[string][]*Feeds
}

type ChainLinkFeeds struct {
	nets []*PriceFeeds
}
