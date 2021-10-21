package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/logistics/db"
	"github.com/llbec/robotech/squadrons/aavesquadron/accounts"
	"github.com/llbec/robotech/squadrons/aavesquadron/lendingpoolevent"
	"github.com/llbec/robotech/utils"
)

const (
	actDeposit  = 1
	actWithdraw = 2
	actBorrow   = 3
	actRepay    = 4
	maxEnd      = 1000
)

var (
	usrDb       db.DB
	accountMana = accounts.NewAccountMana()
)

func initAccounts() {
	usrDb = db.NewDB(
		db.LevelDBBackend,
		"usrs.db",
		filepath.Join(workDir, "database"),
	)
	lendingpoolevent.RegisterHandle(lendingpoolevent.DepositEvent, logDeposit)
	lendingpoolevent.RegisterHandle(lendingpoolevent.BorrowEvent, logBorrow)
	lendingpoolevent.RegisterHandle(lendingpoolevent.WithdrawEvent, logWithdraw)
	lendingpoolevent.RegisterHandle(lendingpoolevent.RepayEvent, logRepay)
}

//************
//db functions
//************
func keyEncode(height uint64, index uint) []byte {
	bts := utils.UInt64ToBytes(height)
	return append(bts, utils.UInt64ToBytes(uint64(index))...)
}

func keyDecode(bts []byte) (height uint64, index uint, err error) {
	if len(bts) != 16 {
		return 0, 0, fmt.Errorf("invalid length %v", len(bts))
	}
	return utils.BytesToUInt64(bts[:8]), uint(utils.BytesToUInt64(bts[8:])), nil
}

func valueEncode(
	account common.Address,
	asset common.Address,
	action uint,
	amount *big.Int,
) []byte {
	v := append(account[:], asset[:]...)
	v = append(v, utils.UInt64ToBytes(uint64(action))...)
	v = append(v, amount.Bytes()...)
	return v
}

func valueDecode(bts []byte) (
	account common.Address,
	asset common.Address,
	action uint,
	amount *big.Int,
	err error,
) {
	if len(bts) <= 48 {
		err = fmt.Errorf("invalid length %v", len(bts))
		return
	}
	account = common.BytesToAddress(bts[:20])
	asset = common.BytesToAddress(bts[20:40])
	action = uint(utils.BytesToUInt64(bts[40:48]))
	amount = big.NewInt(0).SetBytes(bts[48:])
	return
}

func recoveryDB() uint64 {
	iter := usrDb.NewIterator()
	for {
		height, _, err := keyDecode(iter.Key())
		if err != nil {
			panic(err)
		}

		account, asset, action, amount, err := valueDecode(iter.Value())
		if err != nil {
			panic(err)
		}
		switch action {
		case actDeposit:
			accountMana.Deposit(account.String(), asset.String(), amount)
		case actWithdraw:
			accountMana.Withdraw(account.String(), asset.String(), amount)
		case actBorrow:
			accountMana.Borrow(account.String(), asset.String(), amount)
		case actRepay:
			accountMana.Repay(account.String(), asset.String(), amount)
		default:
			panic(fmt.Errorf("unknow action %d", action))
		}

		if !iter.Next() {
			return height
		}
	}
}

func syncBlock(from, to int64) {
	start := from
	end := start + maxEnd
	for {
		if end > to {
			end = to
		}
		log.Printf("sync block from %v to %v\n", start, end)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(LendingPoolAddress)},
			FromBlock: big.NewInt(int64(start)),
			ToBlock:   big.NewInt(int64(end)),
		}
		logs, err := Client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("FilterLogs(%v-%v): %v\n", start, end, err)
		}
		for _, vLog := range logs {
			err = lendingpoolevent.LendingPoolEventHandle(vLog)
			if err != nil {
				log.Printf("HandleLog(%v-%v): %v\n", vLog.BlockNumber, vLog.Index, err)
			}
		}
		if end == to {
			break
		}
		start = end + 1
		end = start + maxEnd
	}
}

func logDeposit(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.DepositEventData)
	if !ok {
		log.Printf("Parse DepositEventData failed!")
		return
	}
	key := keyEncode(height, logIndex)
	v, err := usrDb.Get(key)
	if v == nil {
		log.Println("DB Get:", err)
		err = usrDb.Put(
			key,
			valueEncode(
				data.OnBehalfOf,
				data.Reserve,
				actDeposit,
				data.Amount,
			),
		)
		if err == nil {
			accountMana.Deposit(
				data.OnBehalfOf.String(),
				data.Reserve.String(),
				data.Amount,
			)
		} else {
			log.Printf("Put deposit log(%v-%v):%v", height, logIndex, err)
		}
	} else {
		log.Printf("repeat deposit log(%v-%v)", height, logIndex)
	}
}
func logWithdraw(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.WithdrawEventData)
	if !ok {
		log.Printf("Parse WithdrawEventData failed!")
		return
	}
	key := keyEncode(height, logIndex)
	v, err := usrDb.Get(key)
	if v == nil {
		log.Println("DB Get:", err)
		err = usrDb.Put(
			key,
			valueEncode(
				data.User,
				data.Reserve,
				actWithdraw,
				data.Amount,
			),
		)
		if err == nil {
			accountMana.Withdraw(
				data.User.String(),
				data.Reserve.String(),
				data.Amount,
			)
		} else {
			log.Printf("Put withdraw log(%v-%v):%v", height, logIndex, err)
		}
	} else {
		log.Printf("repeat withdraw log(%v-%v)", height, logIndex)
	}
}
func logBorrow(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.BorrowEventData)
	if !ok {
		log.Printf("Parse BorrowEventData failed!")
		return
	}
	key := keyEncode(height, logIndex)
	v, err := usrDb.Get(key)
	if v == nil {
		log.Println("DB Get:", err)
		err = usrDb.Put(
			key,
			valueEncode(
				data.User,
				data.Reserve,
				actBorrow,
				data.Amount,
			),
		)
		if err == nil {
			accountMana.Borrow(
				data.User.String(),
				data.Reserve.String(),
				data.Amount,
			)
		} else {
			log.Printf("Put borrow log(%v-%v):%v", height, logIndex, err)
		}
	} else {
		log.Printf("repeat borrow log(%v-%v)", height, logIndex)
	}
}
func logRepay(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.RepayEventData)
	if !ok {
		log.Printf("Parse RepayEventData failed!")
		return
	}
	key := keyEncode(height, logIndex)
	v, err := usrDb.Get(key)
	if v == nil {
		log.Println("DB Get:", err)
		err = usrDb.Put(
			key,
			valueEncode(
				data.User,
				data.Reserve,
				actRepay,
				data.Amount,
			),
		)
		if err == nil {
			accountMana.Repay(
				data.User.String(),
				data.Reserve.String(),
				data.Amount,
			)
		} else {
			log.Printf("Put repay log(%v-%v):%v", height, logIndex, err)
		}
	} else {
		log.Printf("repeat repay log(%v-%v)", height, logIndex)
	}
}
