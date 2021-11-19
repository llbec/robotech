package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/llbec/robotech/logistics/daemon"
	"github.com/llbec/robotech/squadrons/filecoinsquadron"
	"github.com/llbec/robotech/squadrons/filecoinsquadron/xlandteam"
)

var (
	fRun        bool
	fInit       bool
	fQuit       bool
	fHeight     int64
	fGet        int64
	fCheck      int64
	filecoinAPI *filecoinsquadron.FileCoinAPI
	xlandMap    map[string]*xlandteam.XlandTeam
	payeeMap    map[string]bool
	waitsecond  int
)

func init() {
	flag.BoolVar(&fRun, "r", false, "running daemon")
	flag.BoolVar(&fInit, "i", false, "create config file with example values")
	flag.BoolVar(&fQuit, "q", false, "quit")
	//flag.StringVar(&fPath, "p", ".", "specify config file path. defult is \".\"")
	flag.Int64Var(&fHeight, "n", math.MaxInt64, "start height. default is current block")
	flag.Int64Var(&fGet, "g", math.MaxInt64, "Get Block")
	flag.Int64Var(&fCheck, "c", math.MaxInt64, "Check height and send notifacation")

	/*file := "./" + "running" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[XlanSaveValue]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)*/
}

func main() {
	flag.Parse()

	wd, e := os.Getwd()
	if e != nil {
		log.Printf("os.Getwd failed: %v\n", e)
		return
	}

	if fInit {
		err := GeneratConfig(wd)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("create xland.yaml at %v\n", wd)
		}
		return
	}
	if fQuit {
		d := daemon.NewDaemon(1234, nil)
		d.Stop()
		return
	}
	if fGet != math.MaxInt64 {
		LoadConfig(wd)
		initEnv()
		getTipSet()
		return
	}
	if fCheck != math.MaxInt64 {
		LoadConfig(wd)
		initEnv()
		checkTipSet(fCheck)
		return
	}
	if fRun {
		LoadConfig(wd)
		if fHeight == math.MaxInt64 {
			fmt.Printf("Run failed: no start block number\n")
			goto HELP
		}
		d := daemon.NewDaemon(1234, radar)
		d.Run(filepath.Join(wd, "running.log"))
		return
	}

HELP:
	flag.Usage()
}

func initEnv() {
	rpc := filecoinsquadron.NewRpc(Server, "/rpc/v0", "http", Token)
	subs := filecoinsquadron.NewSubscribe(Server, "/rpc/v0", "ws", Token)
	filecoinAPI = filecoinsquadron.NewFileCoinAPI(rpc, subs)
	xlandMap = make(map[string]*xlandteam.XlandTeam)
	payeeMap = make(map[string]bool)
	for t, xl := range targets {
		xlandMap[t] = xlandteam.NewXlandTeam(xl)
		payeeMap[t] = true
	}
}

func radar(chSig, chExit chan int) {
	waitsecond = 60
	blocknotify := make(chan int64, 10)
	initEnv()

	go watchHeight(blocknotify)

	for {
		select {
		case current := <-blocknotify:
			err := tipSetRadar(fHeight, current)
			if err != nil {
				log.Println(err)
			} else {
				fHeight = current + 1
			}
		case <-chSig:
			goto EXIT
		}
	}

EXIT:
	chExit <- 1
}

func watchHeight(newHeight chan int64) {
	tryCount := 0
	for {
		if tryCount >= recLimit {
			log.Println("Quit connect failed ", tryCount)
			return
		}
		conn, err := filecoinAPI.ChainNotify()
		if err == nil {
			tryCount = 0
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("ReadMessage:", err)
					break
				}

				height, err := filecoinAPI.ReadHeightFromTipSet(message)
				if err != nil {
					log.Println("ReadHeightFromTipSet:", err)
					break
				}
				newHeight <- height
			}
		}
		log.Println("ChainNotify:", err)
		time.Sleep(time.Duration(waitsecond) * time.Second)
		tryCount += 1
		conn.Close()
	}
}

func getTipSet() {
	tipsetBytes, err := filecoinAPI.GetTipsetByHeight(fGet)
	if err != nil {
		fmt.Printf("getTipSet(%v): %v\n", fGet, err)
		return
	}

	//fmt.Print(filecoinAPI.TipsetString(tipsetBytes))

	tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
	if err != nil {
		fmt.Printf("getTipSet Height(%v): %v\n", fGet, err)
		return
	}

	//fmt.Print(filecoinAPI.BlockSting(tipset.Blocks[0]))

	count := 0
	targets := make(map[string]filecoinsquadron.MsgInfo)
	fmt.Printf("TipSet(%v):\n", fGet)
	for _, b := range tipset.Blocks {
		//fmt.Print(filecoinAPI.BlockSting(b))
		msgs, err := filecoinAPI.GetBlockMsgs(b)
		if err != nil {
			fmt.Printf("Block(%v): %v\n", fGet, err)
		}
		for _, m := range msgs {
			count += 1
			fmt.Printf(("\t%v: %v -> %v, %v\n"), m.Cid, m.From, m.To, m.Value)
			if m.To == Payee {
				targets[m.Cid] = m
			}
		}
	}
	fmt.Printf("msgs %d\n", count)
	fmt.Println(targets)
}

func tipSetRadar(start, current int64) error {
	for last := start; last <= current; last++ {
		targets := make(map[string]filecoinsquadron.MsgInfo)
		tipsetBytes, err := filecoinAPI.GetTipsetByHeight(last)
		if err != nil {
			return fmt.Errorf("GetTipsetByHeight Height(%v): %v", last, err)
		}

		tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
		if err != nil {
			return fmt.Errorf("ReadTipSet Height(%v): %v", last, err)
		}

		for _, b := range tipset.Blocks {
			msgs, err := filecoinAPI.PayeeRadarInBlock(payeeMap, b)
			if err != nil {
				return fmt.Errorf("PayeeRadarInBlock Height(%v-%s): %v", last, b, err)
			}
			for _, m := range msgs {
				targets[m.Cid] = m
			}
		}

		for _, tx := range targets {
			xland := xlandMap[tx.To]
			if xland == nil {
				log.Printf("Target(%v) no xland serve address", tx.To)
				continue
			}
			err := xland.XlandSaveValue(tipset.Timestamp, tx.Value.String())
			log.Printf("Block(%v), save value<%v-%v> @ %v\n",
				tipset.Height,
				tipset.Timestamp,
				tx.Value,
				tx.To,
			)
			if err != nil {
				log.Printf("save value failed(%v) @ %v\n", err, tx.To)
			}
		}
		log.Printf("Block(%v) match %v\n", tipset.Height, len(targets))
	}
	return nil
}

func checkTipSet(height int64) {
	targets := make(map[string]filecoinsquadron.MsgInfo)
	tipsetBytes, err := filecoinAPI.GetTipsetByHeight(height)
	if err != nil {
		fmt.Printf("checkTipSet Height(%v): %v", height, err)
		return
	}

	tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
	if err != nil {
		fmt.Printf("checkTipSet Height(%v): %v", height, err)
		return
	}

	for _, b := range tipset.Blocks {
		msgs, err := filecoinAPI.PayeeRadarInBlock(payeeMap, b)
		if err != nil {
			fmt.Printf("checkTipSet Height(%v-%s): %v", height, b, err)
			return
		}
		for _, m := range msgs {
			targets[m.Cid] = m
		}
	}

	for _, tx := range targets {
		xland := xlandMap[tx.To]
		if xland == nil {
			fmt.Printf("Target(%v) no xland serve address", tx.To)
			continue
		}
		err := xland.XlandSaveValue(tipset.Timestamp, tx.Value.String())
		fmt.Printf("Block(%v), save value<%v-%v> @ %v\n",
			tipset.Height,
			tipset.Timestamp,
			tx.Value,
			tx.To,
		)
		if err != nil {
			fmt.Printf("save value failed(%v) @ %v\n", err, tx.To)
		}
	}
	fmt.Printf("Block(%v) match %v\n", tipset.Height, len(targets))
}
