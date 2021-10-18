package main

import (
	"io/ioutil"
)

func ReadABI(abiFile string) error {
	f, err := ioutil.ReadFile(abiFile)
	if err != nil {
		return err
	}
	contractABI = string(f)
	return nil
}
