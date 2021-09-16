package aavesquadron_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
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

var (
	list = []string{"0x94aE37C1043a7Dba3cC56A3CE95485C8fad9924B",
		"0x7C48C159F57EBA48FA42D8a4944569e5D6fEd0EC",
		"0x7F9B46940ced99C919ff416bc70e7D98D55160F2",
		"0x6011fa84D30430353eDcDD1A8c24806e66B418AA",
		"0xf1AfcF0cC6521D196aBC295215a51Ed90C47e384",
		"0x475882F54ff00bb75955BA2107FA4C31A8477bFe",
		"0x085Ac76060de8d3bF1C6A20522c2c1770CEB6180",
		"0xE6bfB34C43a50e4EdFdc11167D6b00836e8c3e0e",
		"0x4426280d8139053938b21440Fdfa69443b71650D",
		"0xA0126462a2031EdA82BbcfC6B63dc5eCE02048b9",
		"0x1ae3211e249fEBFE34359c885e8c4127b86b2324",
		"0x9118ad4d5444c710Ac74Fb7FcD3e967CD7fEe9e9",
		"0xc160dA89253C167e7fe85b43B2f316929B58fa76",
		"0x802C4bcc581Ac02f80C6901ac994b7236332F6EF"}
)

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
	for i := 0; i < len(list); i++ {
		data, err := squadron.LendingPool.GetUserAccountData(
			nil,
			common.HexToAddress(list[i]))
		if err != nil {
			t.Fatal(err)
		}
		//collateral := big.NewInt(1).Div(data.TotalCollateralETH, big.NewInt(100000000))
		//debt := big.NewInt(1).Div(data.TotalDebtETH, big.NewInt(100000000))
		fmt.Printf("%s: collateral(%10v), debt(%v)\n", list[i], data.TotalCollateralETH, data.TotalDebtETH)
	}
}

func Test_readLog(t *testing.T) {
	squadron := aavesquadron.NewAaveSquadron()
	err := squadron.SetRpcClient("https://http-mainnet-node.huobichain.com")
	if err != nil {
		t.Fatal(err)
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress("0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")},
		FromBlock: big.NewInt(int64(8084871)),
		ToBlock:   big.NewInt(int64(8084871)),
	}
	logs, err := squadron.RpcClient.FilterLogs(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	for _, vLog := range logs {
		fmt.Printf("height(%v)-index(%v)\n", vLog.BlockNumber, vLog.Index)
	}
}
