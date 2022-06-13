package main

import (
	"fmt"
	"os"

	"github.com/llbec/robotech/utils"
)

func main() {
	if err := utils.LocalEnv(); err != nil {
		fmt.Printf("env file error: %v\n", err)
		os.Exit(-1)
	}
}
