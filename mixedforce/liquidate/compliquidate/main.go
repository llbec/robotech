package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/llbec/robotech/logistics/db"
	"github.com/llbec/robotech/squadrons/compoundsquadron/compliquidation"
)

var (
	fHelp      bool
	fInit      bool
	fPath      string
	fThreshold float64
	fManual    bool
)

func init() {
	flag.BoolVar(&fHelp, "h", false, "help")
	flag.BoolVar(&fInit, "i", false, "create config file with example values")
	flag.StringVar(&fPath, "p", ".", "specify config file path. defult is \".\"")
	flag.Float64Var(&fThreshold, "t", 10, "liquidation threshold. default is 10.0")
	flag.BoolVar(&fManual, "m", false, "scan once or loop. default is loop(false)")

	file := "./" + "liquidation" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[Compound liquidation]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func main() {
	flag.Parse()
	if fHelp {
		flag.Usage()
		return
	}
	if fInit {
		GeneratConfig(fPath)
		return
	}

	var secretkey string
	fmt.Println("Enter liquidator's private key")
	fmt.Scanf("%s", &secretkey)

	LoadConfig(fPath)

	ld := db.NewDB(db.LevelDBBackend, proDB, filepath.Join(dbPath, "program"))
	dd := db.NewDB(db.LevelDBBackend, debtorDB, filepath.Join(dbPath, "debtor"))
	defer func() {
		ld.Close()
		dd.Close()
	}()

	var cl *compliquidation.CompLiquidation
	if fManual {
		cl = compliquidation.NewReadOnce(rpc, compTroller, flashContact, big.NewFloat(fThreshold))
	} else {
		cl = compliquidation.NewLiquidation(rpc, wsHost, compTroller, flashContact, big.NewFloat(fThreshold))
	}
	cl.Set(ld, dd)
	cl.SetAuther(secretkey)

	fmt.Printf("Robot start for %v\n", cl.Author())

	if fManual {
		cl.Once()
	} else {
		cl.Run()
	}
}
