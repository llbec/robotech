package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/llbec/robotech/logistics/daemon"
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

	if fRun {
		loadConfig(workDir)
		prepare()
		d := daemon.NewDaemon(6960, thread)
		d.Run(filepath.Join(workDir, "run.log"))
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
	registerEvent()
	err := envInit()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func startInfo() uint64 {
	height, err := rpcClient.GetHTTPClient().BlockNumber(context.Background())
	if err != nil {
		panic(fmt.Errorf("BlockNumber: %v", err.Error()))
	}
	chainID, err := rpcClient.GetHTTPClient().ChainID(context.Background())
	if err != nil {
		panic(fmt.Errorf("ChainID: %v", err.Error()))
	}
	fmt.Printf("Thread(%v) start at Chain(%v) from height %v\n", os.Getpid(), chainID.Uint64(), height)
	return height
}

func thread(chSig, chExit chan int) {
	curHeight := startInfo()
	tcBlk := time.NewTicker(time.Duration(blockPeroid) * time.Second)
	defer tcBlk.Stop()
	for {
		select {
		//new block, handle lending pool event
		case <-tcBlk.C:
			tcBlk.Stop()
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
}
