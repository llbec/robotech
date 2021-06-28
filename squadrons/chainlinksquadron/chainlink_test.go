package chainlinksquadron_test

import (
	"fmt"
	"testing"

	"github.com/llbec/robotech/squadrons/chainlinksquadron"
)

func Test_getPrice(t *testing.T) {
	node := chainlinksquadron.NewChainNode()
	node.AddNode("heco", "https://http-mainnet-node.huobichain.com")

	clSquadron := chainlinksquadron.NewChainLinkSquadron(node, "./pricefeeds/feeds")
	price, deci, err := clSquadron.GetLastPrice("heco", "FIL", "USD", false)
	fmt.Println(price, deci, err)
}
