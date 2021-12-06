package liquidation

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/armory/aave"
)

func Test_GetDebtors(t *testing.T) {
	//xc, _ := NewXCom("wss://ws-testnet.hecochain.com", "0x074Dd89fBEb91C6D3C248390E091cbd02CdB1B9c")
	//xc, _ := NewXCom("https://http-mainnet-node.huobichain.com", "wss://heco.apron.network/ws/v1/afdff8d504f14ee8a2643112ec9578aa", "0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")
	xc, _ := NewXCom()
	list, _ := xc.GetDebtors()
	if len(list) == 0 {
		t.Errorf("Length invalid: %v", len(list))
	}
	fmt.Println(len(list))
	fmt.Println(list[common.HexToAddress("0x7c0Fc56fE45c98cCF000DD548F457D2fCA9ECC6e")])
	fmt.Println(list[common.HexToAddress("0x744808B56E9B470CE901ABaa9A6261E92928cE82")])
}

func Test_GetLiqudationList(t *testing.T) {
	//xc, _ := NewXCom("https://http-mainnet-node.huobichain.com", "wss://heco.apron.network/ws/v1/afdff8d504f14ee8a2643112ec9578aa", "0x074Dd89fBEb91C6D3C248390E091cbd02CdB1B9c")
	xc, _ := NewXCom()
	debtors := xc.GetLiqudationList()
	fmt.Println(len(debtors))
	fmt.Println(debtors)
	fmt.Println(debtors["0x520cCF78ed1810801296Fb7d06c20344A16Bf9AA"])
}

func Test_ReadAccount(t *testing.T) {
	//app, _ := NewXCom("https://http-testnet.hecochain.com", "0x074Dd89fBEb91C6D3C248390E091cbd02CdB1B9c")
	//xc, _ := NewXCom("https://http-mainnet-node.huobichain.com", "wss://heco.apron.network/ws/v1/afdff8d504f14ee8a2643112ec9578aa", "0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")
	xc, _ := NewXCom()
	//xc.ReadAccount("0x7e3e766dca5131DaeC1d62401FCeC03c05D5B1dc")
	//xc.ReadAccount("0x3a277438aeBd6745e770d729a37F60bD4925740a")

	//xc.ReadAccount("0x1DB58997BE23366B0e435cb253bf23154D7356CB")
	xc.ReadAccount("0x1ae3211e249fEBFE34359c885e8c4127b86b2324")
	//xc.ReadAccount("0x4426280d8139053938b21440Fdfa69443b71650D")
	//xc.ReadAccount("0xA87Dc7dc8890c519502415D5CDf9C7a3D32fAaE8")
	//xc.ReadAccount("0x475882F54ff00bb75955BA2107FA4C31A8477bFe")
}

func Test_WatchDebt(t *testing.T) {
	//xc, _ := NewXCom("https://http-mainnet-node.huobichain.com", "wss://heco.apron.network/ws/v1/afdff8d504f14ee8a2643112ec9578aa", "0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")
	xc, _ := NewXCom()
	xc.WatchDebt(startblock.Uint64())
}

func Test_GetPrice(t *testing.T) {
	xc, _ := NewXCom()
	fmt.Println(xc.GetAssetPrice(common.HexToAddress("0xae3a768f9aB104c69A7CD6041fE16fFa235d1810")))
}

func Test_AccountStatus(t *testing.T) {
	//xc, _ := NewXCom()
	collateralAsset := ""
	collateralBalance := big.NewInt(0)
	collateralWorth := big.NewInt(0)
	debtAsset := ""
	debtBalance := big.NewInt(0)
	debtWorth := big.NewInt(0)

	usr := common.HexToAddress("0x3d095f92152a37c6b0e84707d48080c807126a2c")
	initEnv()
	account, err := Contracts.LendingPool.GetUserAccountData(nil, usr)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf(
		"%v status:\n\tHealth: %v\n\tCollateral(%v)-Debt(%v)\n\tAvailable borrow: %v\n",
		usr,
		account.HealthFactor,
		account.TotalCollateralETH,
		account.TotalDebtETH,
		account.AvailableBorrowsETH,
	)

	for s, reserve := range Reserves {
		accountReserve, err := Contracts.ProtocaldataProvider.GetUserReserveData(nil, reserve, usr)
		if err != nil {
			t.Errorf("%v get reserve data error: %v", s, err)
		}
		price, err := Contracts.AaveOracle.GetAssetPrice(nil, reserve)
		if err != nil {
			t.Errorf("%v get price error: %v", s, err)
		}

		w := big.NewInt(1).Mul(accountReserve.CurrentATokenBalance, price)
		if w.Cmp(collateralWorth) > 0 {
			collateralBalance = accountReserve.CurrentATokenBalance
			collateralAsset = s
			collateralWorth = w
		}
		dsw := big.NewInt(1).Mul(accountReserve.CurrentStableDebt, price)
		if dsw.Cmp(debtWorth) > 0 {
			debtAsset = s
			debtBalance = accountReserve.CurrentStableDebt
			debtWorth = dsw
		}
		dvw := big.NewInt(1).Mul(accountReserve.CurrentVariableDebt, price)
		if dvw.Cmp(debtWorth) > 0 {
			debtAsset = s
			debtBalance = accountReserve.CurrentVariableDebt
			debtWorth = dvw
		}
		fmt.Println(s, accountReserve, price, w, dsw, dvw)
	}
	fmt.Printf("collateral:%v-%v\ndebt:%v-%v\n", collateralAsset, collateralBalance, debtAsset, debtBalance)
}

func Test_Liquidation(t *testing.T) {
	xc, _ := NewXCom()
	xc.ReadAccount("0x1ae3211e249fEBFE34359c885e8c4127b86b2324")
	xc.LiquidationAccount(common.HexToAddress("0x1ae3211e249fEBFE34359c885e8c4127b86b2324"))
	xc.ReadAccount("0x1ae3211e249fEBFE34359c885e8c4127b86b2324")
}

func Test_CollateralManager(t *testing.T) {
	xc, _ := NewXCom()
	debtor := common.HexToAddress("0x1ae3211e249fEBFE34359c885e8c4127b86b2324")
	collateralAsset, debtAsset, debtToCover := xc.GetUsrLiquidationData(debtor)
	fmt.Println(collateralAsset, debtAsset, debtToCover)

	collateralManager, err := aave.NewLendingPoolCollateralManager(common.HexToAddress("0x2a7330aCF4A8526924b82Ef8D4B4EC0070882fAC"), Clients[HTTPRPC])
	if err != nil {
		t.Error(err)
	}
	tx, err := collateralManager.LiquidationCall(GetAuther(), collateralAsset, debtAsset, debtor, debtToCover, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tx.Hash())
}
