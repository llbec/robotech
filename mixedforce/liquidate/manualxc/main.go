package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/squadrons/aavesquadron/liquidation"
)

func main() {
	xc, err := liquidation.NewXCom()
	if err != nil {
		fmt.Println(err)
		return
	}
	list := xc.GetLiqudationList()
	fmt.Println("number of liquidation ", len(list))

	for d := range list {
		fmt.Printf("liquidation: %v\n", d)
		xc.LiquidationAccount(common.HexToAddress(d))
	}
}
