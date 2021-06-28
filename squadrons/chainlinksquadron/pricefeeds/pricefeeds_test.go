package pricefeeds_test

import (
	"fmt"
	"testing"

	"github.com/llbec/robotech/squadrons/chainlinksquadron/pricefeeds"
)

func TestNewFeeds(t *testing.T) {
	pfs := pricefeeds.NewChainLinkFeeds("./feeds")
	fmt.Print(pfs.String())
}

func TestGetContract(t *testing.T) {
	pfs := pricefeeds.NewChainLinkFeeds("./feeds")
	a, d, e := pfs.GetContract("heco", false, "FIL", "USD")
	if e != nil {
		t.Fatal(e)
	}
	if a != "0x4d8869eCF1F8C78C0bd2439c4BcAE50AC8420bC4" {
		t.Fatalf("contract address got %v, expect 0x4d8869eCF1F8C78C0bd2439c4BcAE50AC8420bC4", a)
	}
	if d != 8 {
		t.Fatalf("decimal got %v, expect 8", d)
	}
}
