package fallbackteam_test

import (
	"fmt"
	"testing"

	"github.com/llbec/robotech/squadrons/balasquadron/fallbackteam"
)

func Test_Getprice(t *testing.T) {
	fb := fallbackteam.NewFallback("https://http-mainnet.hoosmartchain.com", "0x2C2882158DC6774BF0cF2168065494A70135D40d")
	assets := []string{"BCH", "BTC"}
	for _, a := range assets {
		p, e := fb.GetPrice(a)
		fmt.Println(p, e)
	}
}
