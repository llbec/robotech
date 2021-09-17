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
	filecoinAPI *filecoinsquadron.FileCoinAPI
	xland       *xlandteam.XlandTeam
	waitsecond  int
)

func init() {
	flag.BoolVar(&fRun, "r", false, "running daemon")
	flag.BoolVar(&fInit, "i", false, "create config file with example values")
	flag.BoolVar(&fQuit, "q", false, "quit")
	//flag.StringVar(&fPath, "p", ".", "specify config file path. defult is \".\"")
	flag.Int64Var(&fHeight, "n", math.MaxInt64, "start height. default is current block")
	flag.Int64Var(&fGet, "g", math.MaxInt64, "Get Block")

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
		getTipSet()
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

func radar(chSig, chExit chan int) {
	waitsecond = 60
	rpc := filecoinsquadron.NewRpc(Server, "/rpc/v0", "http", Token)
	subs := filecoinsquadron.NewSubscribe(Server, "/rpc/v0", "ws", Token)
	filecoinAPI = filecoinsquadron.NewFileCoinAPI(rpc, subs)
	xland = xlandteam.NewXlandTeam(XlandServer)

	blocknotify := make(chan int64, 10)

	go watchHeight(blocknotify)

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

EXIT:
	chExit <- 1
}

func watchHeight(newHeight chan int64) {
	for {
		conn, err := filecoinAPI.ChainNotify()
		if err != nil {
			log.Println("ChainNotify:", err)
			time.Sleep(time.Duration(waitsecond) * time.Second)
			continue
		}
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage:", err)
				break
			}

			height, err := filecoinAPI.ReadHeightFromTipSet(message)
			if err != nil {
				log.Println("ReadMessage:", err)
			}
			newHeight <- height
		}
	}
}

func getTipSet() {
	rpc := filecoinsquadron.NewRpc(Server, "/rpc/v0", "http", Token)
	filecoinAPI = filecoinsquadron.NewFileCoinAPI(rpc, nil)
	tipsetBytes, err := filecoinAPI.GetTipsetByHeight(fGet)
	if err != nil {
		fmt.Printf("Height(%v): %v\n", fGet, err)
		return
	}

	//fmt.Print(filecoinAPI.TipsetString(tipsetBytes))

	tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
	if err != nil {
		fmt.Printf("Height(%v): %v\n", fGet, err)
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
			return fmt.Errorf("Height(%v): %v", last, err)
		}

		tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
		if err != nil {
			return fmt.Errorf("Height(%v): %v", last, err)
		}

		for _, b := range tipset.Blocks {
			msgs, err := filecoinAPI.PayeeRadarInBlock(Payee, b)
			if err != nil {
				return fmt.Errorf("Height(%v): %v", last, err)
			}
			for _, m := range msgs {
				targets[m.Cid] = m
			}
		}

		for _, tx := range targets {
			err := xland.XlandSaveValue(tipset.Timestamp, tx.Value.String())
			log.Printf("Block(%v), save value<%v-%v>\n", tipset.Height, tipset.Timestamp, tx.Value)
			if err != nil {
				log.Printf("save value failed(%v)\n", err)
			}
		}
		log.Printf("Block(%v) match %v\n", tipset.Height, len(targets))
	}
	return nil
}
