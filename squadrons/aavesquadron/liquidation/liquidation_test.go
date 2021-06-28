package liquidation_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
)

func Test_Price(t *testing.T) {
	url := "https://http-mainnet-node.huobichain.com"
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(fmt.Errorf("ethclient.Dial error %v, %v", err, url))
	}

	oracle, err := aave.NewPriceOracle(common.HexToAddress("0x1be62C4A97a45B04628E2B4f38F2eC71cC709e24"), client)
	if err != nil {
		t.Fatal(fmt.Errorf("NewPriceOracle: %v", err))
	}

	token := common.HexToAddress("0x5545153CCFcA01fbd7Dd11C0b23ba694D9509A6F")
	ap, err := oracle.GetAssetPrice(nil, token)
	if err != nil {
		t.Fatal(fmt.Errorf("GetAssetPrice: %v", err))
	}
	fmt.Printf("%v(%v)\n", token, ap)
}
