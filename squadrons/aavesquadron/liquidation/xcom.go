package liquidation

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"path/filepath"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/llbec/robotech/logistics/db"
	"github.com/llbec/robotech/squadrons/aavesquadron/lendingpoolevent"
	"github.com/llbec/robotech/utils"
)

var (
	startblock = big.NewInt(4000000)
	//startblock = big.NewInt(3600000)
	healthline = big.NewInt(1e18)
	maxRange   = big.NewInt(1000)
	//listdividing = big.NewInt(11e17)
)

type SAccount struct {
	count  uint
	health *big.Float
}

type accountdata struct {
	TotalCollateralETH          *big.Int
	TotalDebtETH                *big.Int
	AvailableBorrowsETH         *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}

func (xc *XCom) ReadAccount(user string) {
	account, err := Contracts.LendingPool.GetUserAccountData(nil, common.HexToAddress(user))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user, account)
}

func (xc *XCom) GetDebtors() (map[common.Address]uint, uint64) {
	debtors := make(map[common.Address]uint)

	//get currunt height
	height, err := Clients[HTTPRPC].BlockNumber(context.Background())
	if err != nil {
		log.Fatal("GetDebtors:BlockNumber ", err)
	}

	start := startblock
	current := big.NewInt(int64(height))
	end := big.NewInt(1).Add(start, maxRange)

	fmt.Printf("Read: %v-%v\n", start.Uint64(), current.Uint64())

	for {
		if end.Cmp(current) > 0 {
			end = current
		}
		fmt.Printf("\t%v-%v\n", start, end)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{Contracts.lendingpoolAddr},
			FromBlock: start,
			ToBlock:   end,
		}
		//fmt.Printf("%v@%v\n", query, start)
		logs, err := Clients[HTTPRPC].FilterLogs(context.Background(), query)
		if err != nil {
			log.Fatal("GetDebtors:FilterLogs ", err)
		}
		for _, vLog := range logs {
			data, err := lendingpoolevent.HandleLendingPoolEvent(vLog, lendingpoolevent.BorrowEvent)
			if err == nil && data != nil {
				usr := data.(lendingpoolevent.BorrowEventData).User
				fmt.Printf("add usr:%v\n", usr)
				debtors[usr] = debtors[usr] + 1
			}
		}
		if end.Cmp(current) >= 0 {
			break
		}
		start = end
		end = big.NewInt(1).Add(start, maxRange)
	}
	return debtors, height
}

func (xc *XCom) GetLiqudationList() map[string]uint64 {
	list := make(map[string]uint64)
	debtors, _ := xc.GetDebtors()
	for debtor := range debtors {
		account, err := Contracts.LendingPool.GetUserAccountData(nil, debtor)
		if err != nil {
			log.Fatal("GetLiqudationList:GetUserAccountData ", err)
		}
		fmt.Printf("%v@%v\n", debtor, big.NewInt(1).Div(account.HealthFactor, big.NewInt(1e15)).Uint64())
		if account.TotalDebtETH.Cmp(big.NewInt(1).Mul(healthline, big.NewInt(10))) < 0 {
			fmt.Printf("%v@%v skip debt %v\n", debtor, big.NewInt(1).Div(account.HealthFactor, big.NewInt(1e15)).Uint64(), account.TotalDebtETH)
			continue
		}
		if account.HealthFactor.Cmp(healthline) <= 0 {
			list[debtor.String()] = account.HealthFactor.Uint64()
		}
	}
	return list
}

func (xc *XCom) Run() {
	log.Printf("xcom running!\n")

	dbpath := filepath.Join(xcPath(), "database")
	xc.dbs = db.NewDB(db.LevelDBBackend, "uses.db", dbpath)
	xc.LoadDB()

	dbbatch := xc.dbs.NewBatch()
	debtors, start := xc.GetDebtors()
	for debtor, n := range debtors {
		account, err := Contracts.LendingPool.GetUserAccountData(nil, debtor)
		if err != nil {
			log.Fatal("Run:GetUserAccountData ", err)
		}
		fhealth := big.NewFloat(calcFloat(account.HealthFactor, healthline))
		xc.datas[debtor] = SAccount{count: n, health: fhealth}
		dbbatch.Put(debtor.Bytes(), []byte(fhealth.String()))
	}
	dbbatch.Put([]byte("height"), utils.UInt64ToBytes(start))
	dbbatch.Write()
	log.Printf("xcom init height %v with %v debtors\n", start, len(debtors))

	//var wg sync.WaitGroup
	go xc.WatchDebt(start)
	go xc.ScanDebt()
	<-xc.channelInt
	//wg.Wait()
}

func calcFloat(x *big.Int, y *big.Int) float64 {
	return float64(x.Uint64()) / float64(y.Uint64())
}

func (xc *XCom) LoadDB() {
	lastheight, err := xc.dbs.Get([]byte("height"))
	if err == nil {
		startn := big.NewInt(1).SetBytes(lastheight)
		if startn.Cmp(startblock) > 0 {
			startblock = big.NewInt(1).Sub(startn, big.NewInt(1))
		}
	}

	count := 0

	iter := xc.dbs.NewIterator()
	for {
		addr := iter.Key()
		if !reflect.DeepEqual(addr, []byte("height")) {
			fh := new(big.Float)
			fh.SetString(string(iter.Value()))
			xc.datas[common.BytesToAddress(iter.Key())] = SAccount{count: 1, health: fh}
			count += 1
		}
		if !iter.Next() {
			break
		}
	}

	log.Printf("xcom load database at block %v, usr %v\n", startblock.Uint64(), count)
}

func (xc *XCom) WatchDebt(start uint64) {
	log.Printf("WatchDebt thread started!\n")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{Contracts.lendingpoolAddr},
		FromBlock: big.NewInt(int64(start)),
	}

	logs := make(chan types.Log)
	sub, err := Clients[WEBSOCKETRPC].SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal("WatchDebt:SubscribeFilterLogs ", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal("WatchDebt:sub.Err ", err)
		case vLog := <-logs:
			data, err := lendingpoolevent.HandleLendingPoolEvent(vLog, lendingpoolevent.BorrowEvent)
			if err == nil && data != nil {
				//xcliqudation.handleBorrowEvent(data.(poolevent.BorrowEventData))
				usr := data.(lendingpoolevent.BorrowEventData).User

				xc.datamtx.Lock()
				info := xc.datas[usr]
				xc.datamtx.Unlock()

				info.count += 1

				account, err := Contracts.LendingPool.GetUserAccountData(nil, usr)
				if err != nil {
					log.Printf("WatchDebt:GetUserAccountData %v\n", err)
				}

				info.health = big.NewFloat(calcFloat(account.HealthFactor, healthline))

				xc.datamtx.Lock()
				xc.datas[usr] = info
				xc.datamtx.Unlock()

				xc.dbs.Put(usr.Bytes(), []byte(info.health.String()))
				xc.dbs.Put([]byte("height"), utils.UInt64ToBytes(vLog.BlockNumber))
				log.Printf("%v borrow %v(%v)@%v\n", usr, data.(lendingpoolevent.BorrowEventData).Reserve, data.(lendingpoolevent.BorrowEventData).Amount, vLog.BlockNumber)
			}
		}
	}
}

func (xc *XCom) ScanDebt() {
	log.Printf("ScanDebt thread started!\n")
	times := 0
	for {
		times += 1
		xc.liqudationlist = make(map[common.Address]accountdata)
		xc.datamtx.Lock()
		for debtor, info := range xc.datas {
			if info.health.Cmp(big.NewFloat(1.1)) <= 0 || times/1000 == 1 {
				account, err := Contracts.LendingPool.GetUserAccountData(nil, debtor)
				if err != nil {
					log.Printf("ScanDebt:GetUserAccountData %v\n", err)
				}
				xc.datas[debtor] = SAccount{count: info.count, health: big.NewFloat(calcFloat(account.HealthFactor, healthline))}
				//log.Printf("%v @ %v\n", debtor, xc.datas[debtor].health)
				if xc.datas[debtor].health.Cmp(big.NewFloat(1.0)) <= 0 {
					xc.liqudationlist[debtor] = account
					if account.TotalDebtETH.Cmp(LiquidationThreshold) > 0 {
						xc.LiquidationAccount(debtor)
						log.Printf("ID:%v-%v\n\tLiquidation %v-%v\n", times, len(xc.datas), debtor, account)
					} else {
						log.Printf("ID:%v-%v\n\tSkip(%v): %v-%v\n", times, len(xc.datas), LiquidationThreshold, debtor, account)
					}

				}
			}
		}
		if times/1000 == 1 {
			times = 1
		}
		xc.datamtx.Unlock()
		time.Sleep(time.Duration(60) * time.Second)
	}
}

func (xc *XCom) GetUsrLiquidationData(usr common.Address) (common.Address, common.Address, *big.Int) {
	var collateralAsset common.Address
	//collateralBalance := big.NewInt(0)
	collateralWorth := big.NewInt(0)
	var debtAsset common.Address
	debtBalance := big.NewInt(0)
	debtWorth := big.NewInt(0)
	for s, reserve := range Reserves {
		accountReserve, err := Contracts.ProtocaldataProvider.GetUserReserveData(nil, reserve, usr)
		if err != nil {
			log.Printf("GetUsrLiquidationData:GetUserReserveData(%v,%v) %v\n", s, usr, err)
			continue
		}
		price := xc.GetAssetPrice(reserve)
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
		price := xc.GetAssetPrice(debtAsset)
		return collateralAsset, debtAsset, debtBalance.Div(debtBalance.Div(collateralWorth, price), big.NewInt(2))
	}
	return collateralAsset, debtAsset, debtBalance.Div(debtBalance, big.NewInt(2))
}

func (xc *XCom) LiquidationAccount(usr common.Address) {
	collateralAsset, debtAsset, debttocover := xc.GetUsrLiquidationData(usr)
	fmt.Println(collateralAsset, debtAsset, usr, debttocover)
	tx, err := Contracts.FlashLiquidationAdp.Execute(
		GetAuther(),
		flashloan.FlashLiquidationAdapterLiquidationParams{
			CollateralAsset: collateralAsset,
			BorrowedAsset:   debtAsset,
			User:            usr,
			DebtToCover:     debttocover,
			UseEthPath:      false,
		},
	)
	if err != nil {
		log.Printf("LiquidationCall failed: %v %v@%v %v err: %v\n", collateralAsset, debttocover, debtAsset, usr, err)
	} else {
		log.Println(tx.Hash())
	}
}

func (xc *XCom) GetAssetPrice(asset common.Address) *big.Int {
	price, err := Contracts.AaveOracle.GetAssetPrice(nil, asset)
	if err != nil {
		log.Printf("GetWorth(%v) error: %v\n", asset, err)
		return big.NewInt(0)
	}
	return price
}
