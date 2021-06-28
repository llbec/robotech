package compliquidation

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/compound"
	"github.com/llbec/robotech/armory/ethchain"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/llbec/robotech/logistics/db"
	"github.com/llbec/robotech/squadrons/compoundsquadron/ctokenevent"
	"github.com/llbec/robotech/utils"
)

var (
	Zero           = big.NewInt(0)
	CTokenMantissa = 1e18
	Mantissa       = big.NewFloat(CTokenMantissa)
	maxRange       = big.NewInt(1000)
	defaulthealth  = "1.1"
)

type debtorHealth struct {
	Debtor common.Address
	Health *big.Float
}
type HealthList []debtorHealth

func (p HealthList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p HealthList) Len() int           { return len(p) }
func (p HealthList) Less(i, j int) bool { return p[i].Health.Cmp(p[j].Health) < 0 }

type CompLiquidation struct {
	Chain *ethchain.EthChain
	//CTokens   map[string]*compound.CToken
	CStroller *compound.CompTroller
	FlashLoan *flashloan.CompoundLiquidation
	Oracle    *compound.PriceOracle

	startheight int64
	healthList  HealthList
	//markets map[common.Address]*big.Int //market last update debtor height
	//datamtx       sync.Mutex
	liquidationDB db.DB
	debtorDB      db.DB
	author        *bind.TransactOpts
	threshold     *big.Float
}

func NewReadOnce(rpc, troller, flash string, tsd *big.Float) *CompLiquidation {
	rpcclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Panic(err)
	}
	tc, err := compound.NewCompTroller(common.HexToAddress(troller), rpcclient)
	if err != nil {
		log.Panic(err)
	}
	flc, err := flashloan.NewCompoundLiquidation(common.HexToAddress(flash), rpcclient)
	if err != nil {
		log.Panic(err)
	}
	addr, err := tc.Oracle(nil)
	if err != nil {
		log.Panic(err)
	}
	oc, err := compound.NewPriceOracle(addr, rpcclient)
	if err != nil {
		log.Panic(err)
	}

	return &CompLiquidation{
		Chain:     &ethchain.EthChain{ReadClient: rpcclient},
		CStroller: tc,
		FlashLoan: flc,
		Oracle:    oc,
		threshold: tsd,
	}
}

func NewLiquidation(rpc, ws, troller, flash string, tsd *big.Float) *CompLiquidation {
	rpcclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Panic(err)
	}
	wsclient, err := ethclient.Dial(ws)
	if err != nil {
		log.Panic(err)
	}
	tc, err := compound.NewCompTroller(common.HexToAddress(troller), rpcclient)
	if err != nil {
		log.Panic(err)
	}
	flc, err := flashloan.NewCompoundLiquidation(common.HexToAddress(flash), rpcclient)
	if err != nil {
		log.Panic(err)
	}
	addr, err := tc.Oracle(nil)
	if err != nil {
		log.Panic(err)
	}
	oc, err := compound.NewPriceOracle(addr, rpcclient)
	if err != nil {
		log.Panic(err)
	}

	return &CompLiquidation{
		Chain:     &ethchain.EthChain{ReadClient: rpcclient, SubsClient: wsclient},
		CStroller: tc,
		FlashLoan: flc,
		Oracle:    oc,
		threshold: tsd,
	}
}

func (l *CompLiquidation) Set(ldb, ddb db.DB) {
	l.liquidationDB = ldb
	l.debtorDB = ddb
	h, err := ldb.Get([]byte("height"))
	if err != nil {
		l.startheight = 836494
	}
	height := int64(utils.BytesToUInt64(h))
	if height < 836494 {
		l.startheight = 836494
	} else {
		l.startheight = height
	}
	fmt.Printf("Start height: %v\n", l.startheight)
}

func (l *CompLiquidation) CheckAccount(address common.Address) int {
	e, liquidity, shortfall, err := l.CStroller.GetAccountLiquidity(nil, address)
	if err != nil || e.Cmp(Zero) != 0 {
		log.Printf("GetAccountLiquidity(%v): %v,%v,%v,%v\n", address, e, liquidity, shortfall, err)
	}
	if liquidity.Cmp(Zero) == 0 && shortfall.Cmp(Zero) > 0 {
		return 1
	}
	if liquidity.Cmp(Zero) == 0 && shortfall.Cmp(Zero) == 0 {
		return -1
	}
	return 0
}

func (l *CompLiquidation) GetAccountSnapshot(account, ctoken common.Address) (*Asset, *Asset, error) {
	c, err := compound.NewCToken(ctoken, l.Chain.ReadClient)
	if err != nil {
		return nil, nil, fmt.Errorf("NewCToken(%v):%v", ctoken, err)
	}
	eNum, balance, borrow, index, err := c.GetAccountSnapshot(nil, account)
	if err != nil || eNum.Cmp(Zero) != 0 {
		return nil, nil, fmt.Errorf("GetAccountSnapshot(%v):%v", ctoken, err)
	}
	fBalance := new(big.Float).SetInt(balance)
	fBorrow := new(big.Float).SetInt(borrow)
	fIndex := new(big.Float).SetInt(index)

	price, err := l.Oracle.GetUnderlyingPrice(nil, ctoken)
	if err != nil {
		return nil, nil, fmt.Errorf("Oracle(%v):%v", ctoken, err)
	}
	fPrice := new(big.Float).SetInt(price)

	collateral := big.NewFloat(1).Mul(fBalance.Quo(fBalance, Mantissa), fIndex.Quo(fIndex, Mantissa))
	cv := big.NewFloat(1).Mul(collateral, fPrice.Quo(fPrice, Mantissa))
	dv := big.NewFloat(1).Mul(fBorrow.Quo(fBorrow, Mantissa), fPrice)

	symb, _ := c.Name(nil)
	//fmt.Printf("%v:\n\tcollateral(%v)\n\tborrow(%v)\n\tcollateral value(%v)\n\tdebt value(%v)\n", symb, collateral, fBorrow, cv, dv)
	col := &Asset{Symbol: symb, Address: ctoken, Balance: collateral, Value: cv}
	bow := &Asset{Symbol: symb, Address: ctoken, Balance: fBorrow, Value: dv}

	return col, bow, nil
}

func (l *CompLiquidation) LoadDBDebtors() {
	iter := l.debtorDB.NewIterator()
	for {
		addr := common.BytesToAddress(iter.Key())
		health, ok := big.NewFloat(1).SetString(string(iter.Value()))
		if !ok {
			health, ok = big.NewFloat(1).SetString(defaulthealth)
			if !ok {
				panic(fmt.Errorf("%v to float filed", defaulthealth))
			}
		}
		l.healthList = append(l.healthList, debtorHealth{Debtor: addr, Health: health})
		if !iter.Next() {
			break
		}
	}
}

func (l *CompLiquidation) GetAccount(address common.Address) (*Account, error) {
	markets, err := l.CStroller.GetAssetsIn(nil, address)
	if err != nil {
		return nil, fmt.Errorf("GetAssetsIn: %v", err)
	}
	account := &Account{Collateral: big.NewFloat(0), Debt: big.NewFloat(0)}
	account.Address = address.String()

	for _, mkt := range markets {
		c, b, err := l.GetAccountSnapshot(address, mkt)
		if err != nil {
			return nil, fmt.Errorf("GetAccountSnapshot: %v", err)
		}
		if c.Value.Cmp(big.NewFloat(0)) > 0 {
			account.Collaterals = append(account.Collaterals, c)
			account.Collateral.Add(account.Collateral, c.Value)
		}

		if b.Value.Cmp(big.NewFloat(0)) > 0 {
			account.Debts = append(account.Debts, b)
			account.Debt.Add(account.Debt, b.Value)
		}
		//fmt.Println(c, b)
	}

	if account.Debt.Cmp(big.NewFloat(0)) == 0 {
		account.Health = big.NewFloat(math.MaxFloat64)
	} else if account.Collateral.Cmp(big.NewFloat(0)) == 0 {
		account.Health = big.NewFloat(0)
	} else {
		account.Health = big.NewFloat(1).Quo(account.Collateral, account.Debt)
	}

	return account, nil
}

func (l *CompLiquidation) LiquidationDebtor(address common.Address) error {
	act, err := l.GetAccount(address)
	if err != nil {
		return fmt.Errorf("read %v error: %v", address, err)
	}
	if act.Debt.Cmp(l.threshold) < 0 {
		log.Printf("Skip %v, debts(%v) is too small", address, act.Debt.String())
		return nil
	}
	if act.Collateral.Cmp(l.threshold) < 0 {
		log.Printf("Skip %v, Collateral(%v) is too small", address, act.Collateral.String())
		return nil
	}
	maxdebt := &Asset{Value: big.NewFloat(0)}
	for _, dt := range act.Debts {
		if dt.Value.Cmp(maxdebt.Value) > 0 {
			maxdebt = dt
		}
	}
	maxcollateral := &Asset{Value: big.NewFloat(0)}
	for _, cl := range act.Collaterals {
		if cl.Value.Cmp(maxcollateral.Value) > 0 {
			maxcollateral = cl
		}
	}

	halfdebt := big.NewFloat(1).Mul(maxdebt.Value, big.NewFloat(0.5))
	tocover := big.NewInt(0)
	if halfdebt.Cmp(maxcollateral.Value) > 0 {
		f := big.NewFloat(1).Mul(maxdebt.Balance, big.NewFloat(1).Quo(maxcollateral.Value, maxdebt.Value))
		f.Mul(f, Mantissa)
		f.Int(tocover)
	} else {
		f := big.NewFloat(1).Mul(halfdebt, Mantissa)
		f.Int(tocover)
	}

	log.Printf("Liquidation:\n\tBorrow:%v\n\tCollateral:%s-%v\n\tDebts:%s-%v\n",
		address,
		maxcollateral.Symbol, maxcollateral.Value.String(),
		maxdebt.Symbol, maxdebt.Value.String())

	param := flashloan.CompflashLiquidationAdapterLiquidationParams{
		CollateralCtoken: maxcollateral.Address,
		BorrowedCtoken:   maxdebt.Address,
		User:             address,
		DebtToCover:      tocover,
		UseEthPath:       false,
	}
	tx, err := l.FlashLoan.Execute(l.author, param)
	if err != nil {
		return fmt.Errorf("FlashLoan(%v) %v error: %v", address, param, err)
	}
	fmt.Println(tx)
	return nil
}

func (l *CompLiquidation) loadDebtors(ch chan uint64, markets []common.Address) {
	//get currunt height
	height, err := l.Chain.ReadClient.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("BlockNumber ", err)
	}

	start := big.NewInt(l.startheight)
	current := big.NewInt(int64(height))
	end := big.NewInt(1).Add(start, maxRange)

	for {
		batch := l.debtorDB.NewBatch()
		if end.Cmp(current) > 0 {
			end = current
		}
		query := ethereum.FilterQuery{
			Addresses: markets,
			FromBlock: start,
			ToBlock:   end,
		}
		log.Printf("%v@%v\n", start, end)
		logs, err := l.Chain.ReadClient.FilterLogs(context.Background(), query)
		if err != nil {
			log.Println("FilterLogs ", err)
			continue
		}
		for _, vLog := range logs {
			data, err := ctokenevent.HandleEvent(vLog, ctokenevent.BorrowEvent)
			if err == nil && data != nil {
				usr := data.(ctokenevent.BorrowEventData).Borrower
				//fmt.Printf("add usr:%v\n", usr)
				err = batch.Put(usr.Bytes(), []byte(defaulthealth))
				if err != nil {
					log.Println("batch ", err)
					err = batch.Write()
					if err != nil {
						ch <- uint64(l.startheight)
					}
					ch <- vLog.BlockNumber - 1
					return
				}
			}
		}
		err = batch.Write()
		if err != nil {
			ch <- start.Uint64()
		}
		if end.Cmp(current) >= 0 {
			break
		}
		start = end
		end = big.NewInt(1).Add(start, maxRange)
	}
	ch <- end.Uint64()
}

func (l *CompLiquidation) watchDebtor(markets []common.Address, start uint64, res chan uint64) {
	log.Printf("WatchDebt thread started!\n")
	currentblock := start
	query := ethereum.FilterQuery{
		Addresses: markets,
		FromBlock: big.NewInt(int64(currentblock)),
	}
	logs := make(chan types.Log)
	sub, err := l.Chain.SubsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("SubscribeFilterLogs: %v", err)
		res <- currentblock
		return
	}
	for {
		select {
		case err := <-sub.Err():
			log.Printf("sub.Err(%v): ", err)
			res <- currentblock
			return
		case vLog := <-logs:
			data, err := ctokenevent.HandleEvent(vLog, ctokenevent.BorrowEvent)
			if err == nil && data != nil {
				//xcliqudation.handleBorrowEvent(data.(poolevent.BorrowEventData))
				usr := data.(ctokenevent.BorrowEventData).Borrower
				err = l.debtorDB.Put(usr.Bytes(), []byte(defaulthealth))
				if err != nil {
					log.Printf("sub.Err(%v): ", err)
					res <- currentblock
					return
				}
				currentblock = vLog.BlockNumber
				log.Printf("%v: %v borrow %v(%v)@%v\n", vLog.BlockNumber, usr, data.(ctokenevent.BorrowEventData).BorrowAmount, data.(ctokenevent.BorrowEventData).TotalBorrows, vLog.BlockNumber)
			}
		}
	}
}

func (l *CompLiquidation) updateDebtorsThread(quit chan int) {
	markets, err := l.CStroller.GetAllMarkets(nil)
	if err != nil {
		log.Panic(err)
	}

	load := make(chan uint64, 1)
	restart := make(chan uint64, 1)

	go l.loadDebtors(load, markets)
	lastheight := <-load
	l.liquidationDB.Put([]byte("height"), utils.UInt64ToBytes(lastheight))
	restart <- lastheight
	go l.liquidationThread(quit)
loop:
	for {
		select {
		case lastheight = <-restart:
			go l.watchDebtor(markets, lastheight, restart)
		case <-quit:
			break loop
		}
	}
}

func (l *CompLiquidation) liquidationThread(quit chan int) {
	log.Printf("CompLiquidation scan thread started!\n")
	times := 0
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Duration(60) * time.Second)
			times += 1
			l.LoadDBDebtors()
			sort.Sort(l.healthList)
			for i := 0; i < len(l.healthList); i++ {
				r := l.CheckAccount(l.healthList[i].Debtor)
				if r == 1 {
					// should be liquidate
					l.LiquidationDebtor(l.healthList[i].Debtor)
				} else if r < 0 {
					//should be delete this debtor
					l.debtorDB.Delete(l.healthList[i].Debtor.Bytes())
				}
			}
			if times/1000 == 1 {
				times = 1
			}
		}
	}
}

func (l *CompLiquidation) Run() {
	if !l.isDBReady() {
		return
	}
	quit := make(chan int, 3)
	go l.updateDebtorsThread(quit)
	<-quit
}

func (l *CompLiquidation) isDBReady() bool {
	if l.debtorDB == nil || l.liquidationDB == nil {
		log.Printf("DB is null! quit!\n")
		return false
	}
	return true
}

func (l *CompLiquidation) Once() {
	if !l.isDBReady() {
		return
	}

	markets, err := l.CStroller.GetAllMarkets(nil)
	if err != nil {
		log.Panic(err)
	}

	load := make(chan uint64, 1)
	go l.loadDebtors(load, markets)
	lastheight := <-load
	l.liquidationDB.Put([]byte("height"), utils.UInt64ToBytes(lastheight))

	l.LoadDBDebtors()
	sort.Sort(l.healthList)
	for i := 0; i < len(l.healthList); i++ {
		r := l.CheckAccount(l.healthList[i].Debtor)
		if r == 1 {
			// should be liquidate
			l.LiquidationDebtor(l.healthList[i].Debtor)
		} else if r < 0 {
			//should be delete this debtor
			l.debtorDB.Delete(l.healthList[i].Debtor.Bytes())
		}
	}
}

func (l *CompLiquidation) SetAuther(secretkey string) {
	secret, err := crypto.HexToECDSA(secretkey)
	if err != nil {
		log.Fatal("InitEnv-HexToECDSA ", err)
	}
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := l.Chain.ReadClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("InitEnv-PendingNonceAt ", err)
	}
	chainid, err := l.Chain.ReadClient.ChainID(context.Background())
	if err != nil {
		log.Fatal("InitEnv-ChainID ", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(secret, chainid)
	if err != nil {
		log.Fatal("InitEnv-NewKeyedTransactorWithChainID ", err)
	}
	gasPrice, err := l.Chain.ReadClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("InitEnv-SuggestGasPrice ", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4400000) // in units
	auth.GasPrice = gasPrice
	l.author = auth
}

func (l *CompLiquidation) Author() string {
	return l.author.From.String()
}
