package main

import "flag"

var (
	fRun       bool
	fTerminate bool
	fInit      bool
	fLoad      bool
	fDrop      bool
)

func initCmd() {
	flag.BoolVar(&fRun, "r", false, "running daemon")
	flag.BoolVar(&fTerminate, "t", false, "Terminate fallback process. default false")
	flag.BoolVar(&fInit, "i", false, "Generat config file. default false")
	flag.BoolVar(&fLoad, "l", false, "reload config, modify config file first")
}

func parseCmd() {
	flag.Parse()
}

func help() {
	//HELP:
	flag.Usage()
}
