package aavesquadron_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/squadrons/aavesquadron"
)

func Test_LendingPool(t *testing.T) {
	squadron := aavesquadron.NewAaveSquadron()
	err := squadron.SetRpcClient("https://http-mainnet-node.huobichain.com")
	if err != nil {
		t.Fatal(err)
	}
	err = squadron.SetLendingPool("0x1BeB0e1d334a5289b235a4bdF8CA54146627A11a")
	if err != nil {
		t.Fatal(err)
	}
	data, err := squadron.LendingPool.GetUserAccountData(
		nil,
		common.HexToAddress("0x6623741e2900C03A1eA0c5Ef7687f44e30cAE819"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data.TotalCollateralETH, data.TotalDebtETH)
	xFilCfg, err := squadron.LendingPool.GetReserveData(
		nil,
		common.HexToAddress("0xd6cBE226aEaA37fD022039419ef525253e2d819e"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(xFilCfg.ATokenAddress)
	xxfil, err := squadron.GetAToken(xFilCfg.ATokenAddress.Hex())
	if err != nil {
		t.Fatal(err)
	}
	balance, err := xxfil.BalanceOf(nil, common.HexToAddress("0x6623741e2900C03A1eA0c5Ef7687f44e30cAE819"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("BalanceOf is %v\n", balance)
}

func Test_LendingPool_User(t *testing.T) {
	squadron := aavesquadron.NewAaveSquadron()
	err := squadron.SetRpcClient("https://http-mainnet-node.huobichain.com")
	if err != nil {
		t.Fatal(err)
	}
	err = squadron.SetLendingPool("0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")
	if err != nil {
		t.Fatal(err)
	}
	data, err := squadron.LendingPool.GetUserAccountData(
		nil,
		common.HexToAddress("0x7c0Fc56fE45c98cCF000DD548F457D2fCA9ECC6e"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data.TotalCollateralETH, data.TotalDebtETH)
}
