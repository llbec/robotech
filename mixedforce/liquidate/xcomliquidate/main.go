package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/armory/contract"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/llbec/robotech/logistics/daemon"
	"github.com/llbec/robotech/utils"
)

var (
	fInit      bool
	fList      bool
	fRun       bool
	fTerminate bool
	sAccount   string
)

func init() {
	flag.BoolVar(&fInit, "i", false, "Generate a simple config file. default false")
	flag.BoolVar(&fList, "l", false, "list all assets price. default false")
	flag.BoolVar(&fRun, "r", false, "Running fallback process. default false")
	flag.BoolVar(&fTerminate, "t", false, "Terminate fallback process. default false")
	flag.StringVar(&sAccount, "a", "", "read account reserves data")
}

func main() {
	flag.Parse()
	workDir, e := os.Getwd()
	if e != nil {
		log.Printf("os.Getwd failed: %v\n", e)
		return
	}
	if fInit {
		GeneratConfig(workDir)
		return
	}

	if fTerminate {
		InitEnv(workDir)
		d := daemon.NewDaemon(DmPort, nil)
		d.Stop()
		return
	}

	if fList {
		InitEnv(workDir)
		initReserves()
		showReserves()
		return
	}

	if sAccount != "" {
		InitEnv(workDir)
		initReserves()
		updatePrice()
		showAccount(sAccount)
		return
	}

	if fRun {
		InitEnv(workDir)
		initReserves()
		initAccounts()
		d := daemon.NewDaemon(DmPort, liquidateThread)
		d.Run(filepath.Join(workDir, "run.log"))
		return
	}

	//HELP:
	flag.Usage()
}

func showDebtors() {
	count := 0
	for account, amount := range debtors {
		count += 1
		fmt.Printf("****** %v: %v - %v\n", count, account, amount)
	}
}

func getAsset(asset string) (string, int) {
	erc20, _ := contract.ConnectByABI(asset, rpcURL, "erc20.abi")
	var out []interface{}
	erc20.ContractCaller.Contract.Call(nil, &out, "symbol")
	var out1 []interface{}
	erc20.ContractCaller.Contract.Call(nil, &out1, "decimals")
	d, _ := strconv.Atoi(fmt.Sprintf("%v", out1[0]))
	return fmt.Sprintf("%v", out[0]), d
}

func showReserves() {
	updatePrice()
	for asset, price := range Reserves {
		symbol, _ := getAsset(asset.String())
		fmt.Printf(
			"%v\t%7v: %v\n",
			asset,
			symbol,
			utils.BigToRecognizable(price, 8))
	}
}

func showAccount(account string) {
	usr := common.HexToAddress(account)
	aData, err := LendingPool.GetUserAccountData(nil, usr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf(
		"%s status:\n\tHealth: %.6f\n\tCollateral(%.6f)-Debt(%.6f)\n\tAvailable borrow: %.6f\n",
		account,
		utils.BigToRecognizable(aData.HealthFactor, 18),
		utils.BigToRecognizable(aData.TotalCollateralETH, 8),
		utils.BigToRecognizable(aData.TotalDebtETH, 8),
		utils.BigToRecognizable(aData.AvailableBorrowsETH, 8),
	)
	for asset, price := range Reserves {
		symbol, decimals := getAsset(asset.String())
		accountReserve, err := ProtocaldataProvider.GetUserReserveData(nil, asset, usr)
		if err != nil {
			fmt.Printf("GetUserReserveData(%v,%v) %v\n", asset, usr, err)
			continue
		}
		if accountReserve.CurrentATokenBalance.Cmp(big.NewInt(0)) == 0 &&
			accountReserve.CurrentStableDebt.Cmp(big.NewInt(0)) == 0 &&
			accountReserve.CurrentVariableDebt.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		balance := big.NewInt(1).Mul(accountReserve.CurrentATokenBalance, price)
		sDebt := big.NewInt(1).Mul(accountReserve.CurrentStableDebt, price)
		vDebt := big.NewInt(1).Mul(accountReserve.CurrentVariableDebt, price)
		fmt.Printf(
			"%9s:Deposit: %10.6f(%f)\tDebt: Stable%10.6f(%f)\tVariab%10f(%f)\n",
			symbol,
			utils.BigToRecognizable(accountReserve.CurrentATokenBalance, decimals),
			utils.BigToRecognizable(balance, 8+decimals),
			utils.BigToRecognizable(accountReserve.CurrentStableDebt, decimals),
			utils.BigToRecognizable(sDebt, 8+decimals),
			utils.BigToRecognizable(accountReserve.CurrentVariableDebt, decimals),
			utils.BigToRecognizable(vDebt, 8+decimals),
		)
	}
}

func liquidation() {
	for account := range debtors {
		data, err := LendingPool.GetUserAccountData(
			nil,
			common.HexToAddress(account))
		if err != nil {
			log.Println("liquidation GetUserAccountData: ", err)
		} else {
			if data.HealthFactor.Cmp(healthThreshold) <= 0 {
				if data.TotalCollateralETH.Cmp(liquidateLevel1) >= 0 {
					//liquidate immediately
					go liquidationAccount(common.HexToAddress(account))
				} else {
					log.Printf("pass: %v total collateral %v,(line is %v)",
						account,
						data.TotalCollateralETH.String(),
						liquidateLevel1,
					)
				}
			}
			if data.TotalDebtETH.Cmp(common.Big0) == 0 {
				delete(debtors, account)
				addr := common.HexToAddress(account)
				usrDb.Delete(addr[:])
			}
		}
	}
}

func initReserves() {
	Reserves = make(map[common.Address]*big.Int)
	assets, err := LendingPool.GetReservesList(nil)
	if err != nil {
		log.Panicf("GetReservesList: %v", err)
	}
	for _, asset := range assets {
		Reserves[asset] = common.Big0
	}
}

func updatePrice() {
	assets, err := LendingPool.GetReservesList(nil)
	if err != nil {
		log.Printf("GetReservesList: %v", err)
	} else {
		for _, asset := range assets {
			if _, ok := Reserves[asset]; !ok {
				Reserves[asset] = common.Big0
			}
		}
	}
	for asset, price := range Reserves {
		newPrice, err := AaveOracle.GetAssetPrice(nil, asset)
		if err != nil {
			log.Printf("GetAssetPrice(%v) error: %v\n", asset, err)
			continue
		}
		if newPrice.Cmp(price) != 0 {
			Reserves[asset] = newPrice
		}
	}
}

func getUsrLiquidationData(usr common.Address) (common.Address, common.Address, *big.Int) {
	var collateralAsset common.Address
	//collateralBalance := big.NewInt(0)
	collateralWorth := big.NewInt(0)
	var debtAsset common.Address
	debtBalance := big.NewInt(0)
	debtWorth := big.NewInt(0)
	for reserve, price := range Reserves {
		accountReserve, err := ProtocaldataProvider.GetUserReserveData(nil, reserve, usr)
		if err != nil {
			log.Printf("GetUserReserveData(%v,%v) %v\n", reserve, usr, err)
			continue
		}

		if price.Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		w := big.NewInt(1).Mul(accountReserve.CurrentATokenBalance, price)
		if w.Cmp(collateralWorth) > 0 {
			//collateralBalance = accountReserve.CurrentATokenBalance
			collateralAsset = reserve
			collateralWorth = w
		}
		debt := big.NewInt(1).Add(accountReserve.CurrentStableDebt, accountReserve.CurrentVariableDebt)
		dw := big.NewInt(1).Mul(debt, price)
		if dw.Cmp(debtWorth) > 0 {
			debtAsset = reserve
			debtBalance = debt
			debtWorth = dw
		}
	}
	//log.Printf("debt: %v(%v)\n", debtAsset, debtBalance)
	if debtWorth.Div(debtWorth, big.NewInt(2)).Cmp(collateralWorth) > 0 {
		price := Reserves[debtAsset]
		return collateralAsset, debtAsset, debtBalance.Div(debtBalance.Div(collateralWorth, price), big.NewInt(2))
	}
	return collateralAsset, debtAsset, debtBalance.Div(debtBalance, big.NewInt(2))
}

func liquidationAccount(usr common.Address) {
	collateralAsset, debtAsset, debttocover := getUsrLiquidationData(usr)
	fmt.Println("liquidation execute", collateralAsset, debtAsset, usr, debttocover)
	txOps, err := GetAuther(sKey)
	if err != nil {
		log.Printf("liquidation - GetAuther: %v\n", err)
		return
	}
	tx, err := FlashLiquidationAdp.Execute(
		txOps,
		flashloan.FlashLiquidationAdapterLiquidationParams{
			CollateralAsset: collateralAsset,
			BorrowedAsset:   debtAsset,
			User:            usr,
			DebtToCover:     debttocover,
			UseEthPath:      false,
		},
	)
	if err != nil {
		log.Printf("liquidationCall failed: %v %v@%v %v err: %v\n", collateralAsset, debttocover, debtAsset, usr, err)
	} else {
		log.Printf("liquidate(%v): %v@(%v - %v)\n", tx.Hash(), usr, collateralAsset, debttocover)
	}
}

func liquidateThread(chSig, chExit chan int) {
	log.Printf("liquidate start ...\n")
	curHeight = recoveryDB()
	if curHeight < StartHeight {
		curHeight = StartHeight
	}
	//get currunt height
	height, err := Client.BlockNumber(context.Background())
	if err != nil {
		panic(fmt.Errorf("BlockNumber: %v", err))
	}
	log.Printf("sync block from %v, end %v\n", curHeight, height)
	if curHeight < height {
		syncBlock(int64(curHeight), int64(height))
		curHeight = height
	}
	log.Printf("sync block finished with %v debtors\n", len(debtors))
	tcLiq := time.NewTicker(time.Duration(LiquidatePeroid) * time.Second)
	defer tcLiq.Stop()
	tcBlk := time.NewTicker(time.Duration(BlockPeroid) * time.Second)
	defer tcBlk.Stop()
	for {
		select {
		//new block, handle lending pool event
		case <-tcBlk.C:
			tcBlk.Stop()
			height, err = Client.BlockNumber(context.Background())
			if err == nil && curHeight < height {
				syncBlock(int64(curHeight), int64(height))
				curHeight = height
			} else {
				log.Printf("BlockNumber %v (%v): %v", curHeight, height, err)
			}
			tcBlk.Reset(time.Duration(BlockPeroid) * time.Second)
		//timer, check accounts health factor
		case <-tcLiq.C:
			tcLiq.Stop()
			//debtors := accountMana.GetDebtors()
			updatePrice()
			liquidation()
			tcLiq.Reset(time.Duration(LiquidatePeroid) * time.Second)
		//quit signal
		case sig := <-chSig:
			log.Printf("Thread rcv signal: %v\n", sig)
			switch sig {
			case 1:
				showDebtors()
			default:
				goto EXIT
			}
		}
	}
EXIT:
	chExit <- 1
}
