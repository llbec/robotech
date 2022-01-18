package main

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/llbec/robotech/squadrons/aavesquadron/lendingpoolevent"
)

func registerEvent() {
	lendingpoolevent.RegisterHandle(lendingpoolevent.DepositEvent, depositLog)
	lendingpoolevent.RegisterHandle(lendingpoolevent.RepayEvent, repayLog)
}

func repayLog(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.RepayEventData)
	if !ok {
		log.Printf("Parse RepayEventData failed!")
		return
	}
	withdraw(data.Amount)
}

func depositLog(height uint64, logIndex uint, eventData interface{}) {
	data, ok := (eventData).(lendingpoolevent.RepayEventData)
	if !ok {
		log.Printf("Parse RepayEventData failed!")
		return
	}
	withdraw(data.Amount)
}

func withdraw(amount *big.Int) {
	balance, err := usrBalance()
	if err != nil {
		log.Printf("get balance failed: %v", err.Error())
		return
	}
	if amount.Cmp(balance) > 0 {
		amount = balance
	}
	if amount.Cmp(common.Big0) > 1 {
		auth, err := rpcClient.GetAuther(skey, gasLimit)
		if err != nil {
			log.Printf("getAuther failed: %v", err.Error())
			return
		}
		tx, err := lendingPool.Withdraw(auth, common.HexToAddress(asset), amount, common.HexToAddress(usrAddress))
		if err != nil {
			log.Printf("withdraw %v failed: %v\n", amount, err)
		} else {
			log.Printf("withdraw %v @ %v\n", amount, tx.Hash())
		}
	}
}

func usrBalance() (*big.Int, error) {
	return atoken.BalanceOf(nil, common.HexToAddress(usrAddress))
}
