package main

import "flag"

var (
	fRun       bool
	fTerminate bool
)

func initCmd() {
	flag.BoolVar(&fRun, "r", false, "running daemon")
	flag.BoolVar(&fTerminate, "t", false, "Terminate fallback process. default false")
}

func parseCmd() {
	flag.Parse()
}

func help() {
	//HELP:
	flag.Usage()
}
