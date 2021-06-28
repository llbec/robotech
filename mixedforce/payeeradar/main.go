package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/llbec/robotech/squadrons/filecoinsquadron"
	"github.com/llbec/robotech/squadrons/filecoinsquadron/xlandteam"
)

var (
	fHelp       bool
	fInit       bool
	fPath       string
	fHeight     int64
	filecoinAPI *filecoinsquadron.FileCoinAPI
	xland       *xlandteam.XlandTeam
)

func init() {
	flag.BoolVar(&fHelp, "h", false, "help")
	flag.BoolVar(&fInit, "i", false, "create config file with example values")
	flag.StringVar(&fPath, "p", ".", "specify config file path. defult is \".\"")
	flag.Int64Var(&fHeight, "s", math.MaxInt64, "start height. default is current block")

	file := "./" + "running" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[XlanSaveValue]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func main() {
	flag.Parse()
	if fHelp {
		flag.Usage()
		return
	}
	if fInit {
		err := GeneratConfig(fPath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("create xland.yaml at %v\n", fPath)
		}
		return
	}

	// init
	LoadConfig(fPath)
	rpc := filecoinsquadron.NewRpc(Server, "/rpc/v0", "http", Token)
	subs := filecoinsquadron.NewSubscribe(Server, "/rpc/v0", "ws", Token)
	filecoinAPI = filecoinsquadron.NewFileCoinAPI(rpc, subs)
	xland = xlandteam.NewXlandTeam(XlandServer)

	chainnotify := make(chan int64, 10)
	go tipSetRadar(fHeight, chainnotify)
	watchHeight(chainnotify)
}

func watchHeight(newHeight chan int64) {
	for {
		conn, err := filecoinAPI.ChainNotify()
		if err != nil {
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

func tipSetRadar(start int64, chainnotify chan int64) {
	last := start
	for {
		current := <-chainnotify
		if last > current {
			last = current
		}
		for ; last < current; last++ {
			targets := make(map[string]filecoinsquadron.MsgInfo)
			tipsetBytes, err := filecoinAPI.GetTipsetByHeight(last)
			if err != nil {
				log.Println("GetTipsetByHeight:", err)
				continue
			}

			tipset, err := filecoinAPI.ReadTipSet(tipsetBytes)
			if err != nil {
				log.Println("ReadTipSet:", err)
				continue
			}

			for _, b := range tipset.Blocks {
				msgs, err := filecoinAPI.PayeeRadarInBlock(Payee, b)
				if err != nil {
					log.Println("PayeeRadarInBlock:", err)
					continue
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
	}
}
