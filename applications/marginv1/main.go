package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"
	"robotech/applications/marginv1/config"
	"robotech/applications/marginv1/eventhandle"
	"robotech/applications/marginv1/liquidation"
	"robotech/applications/marginv1/storage"
	"robotech/utils/daemon"
	"robotech/utils/ethUtils"
	"time"
)

var (
	fList bool
	fRun  bool
	fStop bool
)

func init() {
	flag.BoolVar(&fList, "l", false, "show all Positions. default false")
	flag.BoolVar(&fRun, "r", false, "Running fallback process. default false")
	flag.BoolVar(&fStop, "s", false, "Stop fallback process. default false")
}

func main() {
	flag.Parse()

	port := config.GetDaemonPort()

	if fStop {
		d := daemon.NewDaemon(port, nil)
		d.Stop()
		return
	}

	if fRun {
		d := daemon.NewDaemon(port, mainThread)
		d.Run(filepath.Join(config.GetWorkDir(), "run.log"))
		return
	}

	if fList {
		d := daemon.NewDaemon(port, nil)
		d.Input(1)
		return
	}

	flag.Usage()
}

func mainThread(chSig, chExit chan int) {
	log.Printf("Thread is running...")
	tcBlk := time.NewTicker(time.Duration(config.GetReadBlockPeriod()) * time.Second)
	defer tcBlk.Stop()

	lastHeight := storage.GetLastBlockNumber()

	for {
		select {
		// block ticker
		case <-tcBlk.C:
			tcBlk.Stop()
			height, err := ethUtils.GetClient().BlockNumber(context.Background())
			if err == nil && lastHeight < height {
				h, poolList := eventhandle.SyncBlock(int64(lastHeight), int64(height))
				lastHeight = uint64(h)
				if len(poolList) > 0 {
					log.Printf("BlockNumber %v (%v) to: %v", lastHeight, height, poolList)
					go liquidation.Liquidation(poolList)
				}
			} else {
				log.Printf("BlockNumber %v (%v): %v", lastHeight, height, err)
			}
			tcBlk.Reset(time.Duration(config.GetReadBlockPeriod()) * time.Second)
		//daemon signal
		case sig := <-chSig:
			log.Printf("Thread rcv signal: %v\n", sig)
			switch sig {
			default:
				goto EXIT
			}
		}
	}
EXIT:
	storage.Close()
	chExit <- 1
}
