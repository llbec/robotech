package main

import "flag"

var (
	fRun       bool
	fTerminate bool
	fInit      bool
)

func initCmd() {
	flag.BoolVar(&fRun, "r", false, "running daemon")
	flag.BoolVar(&fTerminate, "t", false, "Terminate fallback process. default false")
	flag.BoolVar(&fInit, "i", false, "Generat config file. default false")
}

func parseCmd() {
	flag.Parse()
}

func help() {
	//HELP:
	flag.Usage()
}
