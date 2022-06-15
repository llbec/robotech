package main

import "flag"

var (
	fPower    bool
	fOracle   bool
	fShedule  bool
	fBomb     bool
	fDownload bool
	iLine     uint64
	iHeight   uint64
	iTest     uint64
)

func initCmd() {
	flag.BoolVar(&fPower, "p", false, "calc power/period")
	flag.BoolVar(&fOracle, "o", false, "calc block period")
	flag.BoolVar(&fShedule, "s", false, "an shedule")
	flag.BoolVar(&fBomb, "b", false, "bomb diff shedule")
	flag.BoolVar(&fDownload, "d", false, "download header datas")
	flag.Uint64Var(&iLine, "l", 0, "check delta time")
	flag.Uint64Var(&iHeight, "h", 0, "calc average at height")
	flag.Uint64Var(&iTest, "t", 0, "test difficulty increase")
}

func parseCmd() {
	flag.Parse()
}

func help() {
	//HELP:
	flag.Usage()
}
