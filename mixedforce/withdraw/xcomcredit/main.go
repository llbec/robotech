package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/llbec/robotech/logistics/daemon"
	"github.com/llbec/robotech/utils"
)

func init() {
	initCmd()
}

func main() {
	parseCmd()

	workDir, e := os.Getwd()
	if e != nil {
		log.Printf("os.Getwd failed: %v\n", e)
		return
	}

	if fInit {
		generatConfig(workDir)
		return
	}

	if fRun {
		loadConfig(workDir)
		prepare()
		inputKey()
		d := daemon.NewDaemon(6960, thread)
		d.Run(filepath.Join(workDir, "run.log"))
		return
	}

	if fLoad {
		loadConfig(workDir)
		inputKey()
		d := daemon.NewDaemon(6960, nil)
		d.Input(1)
		return
	}

	if fTerminate {
		loadConfig(workDir)
		d := daemon.NewDaemon(6960, nil)
		d.Stop()
		return
	}
	help()
}

func prepare() {
	if fDrop {
		registerEvent()
	}
	err := envInit()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

/*func checkThread(chSig, chExit chan int) {
	curHeight, balance := startInfo()
	tcBlk := time.NewTicker(time.Duration(blockPeroid) * time.Second)
	defer tcBlk.Stop()
	if balance.Cmp(big.NewInt(0)) <= 0 {
		goto EXIT
	}
	for {
		select {
		//timeout,check pool balance
		case <-tcBlk.C:
			tcBlk.Stop()
			amount, err := poolBalance()
			if err == nil {
				log.Printf("%v balance is %.6f", usrAddress, utils.BigToRecognizable(amount, 18))
				if balance.Cmp(big.NewInt(0)) <= 0 {
					goto EXIT
				}
			}
			height, err := rpcClient.GetHTTPClient().BlockNumber(context.Background())
			if err == nil && curHeight < height {
				handleBlock(int64(curHeight), int64(height))
				curHeight = height
			} else {
				log.Printf("BlockNumber %v (%v): %v", curHeight, height, err)
			}
			tcBlk.Reset(time.Duration(blockPeroid) * time.Second)
		//quit signal
		case sig := <-chSig:
			switch sig {
			default:
				goto EXIT
			}
		}
	}
EXIT:
	chExit <- 1
}*/

func thread(chSig, chExit chan int) {
	currentHeight, _ = startInfo()
	tcBlk := time.NewTicker(time.Duration(blockPeroid) * time.Second)
	defer tcBlk.Stop()
	for {
		select {
		//new block, handle lending pool event
		case <-tcBlk.C:
			tcBlk.Stop()
			if fDrop {
				if !checkBlock() {
					goto EXIT
				}
			}
			if fRun {
				if !getProcess() {
					goto EXIT
				}
			}
			tcBlk.Reset(time.Duration(blockPeroid) * time.Second)
		//signal handle
		case sig := <-chSig:
			switch sig {
			case 1: //set internal
				tcBlk.Stop()
				tcBlk.Reset(time.Duration(blockPeroid) * time.Second)
			default:
				goto EXIT
			}
		}
	}
EXIT:
	chExit <- 1
}

func checkBlock() bool {
	balance, err := usrBalance()
	if err == nil {
		log.Printf("%v balance is %.6f", usrAddress, utils.BigToRecognizable(balance, 18))
		if balance.Cmp(big.NewInt(0)) <= 0 {
			return false
		}
	}
	height, err := rpcClient.GetHTTPClient().BlockNumber(context.Background())
	if err == nil && currentHeight < height {
		handleBlock(int64(currentHeight), int64(height))
		currentHeight = height
	} else {
		log.Printf("BlockNumber %v (%v): %v", currentHeight, height, err)
	}
	return true
}

func getProcess() bool {
	height, _ := rpcClient.GetHTTPClient().BlockNumber(context.Background())
	amount, err := poolBalance()
	if err != nil {
		log.Printf("Block@ %v read pool balance failed: %v", height, err)
	} else {
		if amount.Cmp(big.NewInt(1e18)) > 0 {
			balance, err := usrBalance()
			if err == nil {
				log.Printf("%v balance is %.6f", usrAddress, utils.BigToRecognizable(balance, 18))
				if balance.Cmp(big.NewInt(0)) <= 0 {
					return false
				}
				simpleWithdraw(amount, balance)
				balance, err = usrBalance()
				if err == nil && balance.Cmp(big.NewInt(1e17)) < 0 {
					return false
				}
			}
		} else {
			log.Printf("@%v pool have no balance %.6f\n", height, utils.BigToRecognizable(amount, 18))
		}
	}
	return true
}

func startInfo() (uint64, *big.Int) {
	height, err := rpcClient.GetHTTPClient().BlockNumber(context.Background())
	if err != nil {
		panic(fmt.Errorf("BlockNumber: %v", err.Error()))
	}
	chainID, err := rpcClient.GetHTTPClient().ChainID(context.Background())
	if err != nil {
		panic(fmt.Errorf("ChainID: %v", err.Error()))
	}
	balance, err := usrBalance()
	if err != nil {
		panic(fmt.Errorf("usrBalance: %v", err.Error()))
	}
	fmt.Printf("Thread(%v) start at Chain(%v) from height %v\n%v balance is %.6f\n",
		os.Getpid(),
		chainID.Uint64(),
		height,
		usrAddress,
		utils.BigToRecognizable(balance, 18),
	)
	return height, balance
}

func inputKey() {
	var skip string
	fmt.Printf("Input y to set key\n")
	fmt.Scan(&skip)
	if skip != "y" {
		return
	}
	fmt.Printf("Enter key...\n")
	fmt.Scan(&skey)
}
