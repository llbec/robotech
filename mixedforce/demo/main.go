package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/llbec/robotech/logistics/daemon"
)

var (
	fRun  bool
	fStop bool
)

func init() {
	flag.BoolVar(&fRun, "r", false, "running")
	flag.BoolVar(&fStop, "s", false, "stop")
}

func main() {
	flag.Parse()

	if fRun {
		d := daemon.NewDaemon(1234, proc)
		d.Run("demo.log")
		return
	}

	if fStop {
		d := daemon.NewDaemon(1234, nil)
		d.Stop()
		return
	}

	flag.Usage()
}

func proc(chSig, chExit chan int) {
	tc := time.NewTicker(time.Duration(3) * time.Second)
	defer tc.Stop()

	count := 0

WORKING:
	for {
		select {
		case <-tc.C:
			log.Printf("Process %d: ticked %d\n", os.Getpid(), count)
			count++
		case <-chSig:
			break WORKING
		}
	}
	chExit <- 1
}
