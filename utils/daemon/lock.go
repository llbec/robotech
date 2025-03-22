package daemon

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Lock(pid int) error {
	dw, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(dw, ".lock"), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fd, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if string(fd) != "" {
		return fmt.Errorf("process(%s) existed", string(fd))
	}

	_, err = io.WriteString(file, fmt.Sprintf("%d", pid))
	if err != nil {
		return err
	}
	return nil
}

func Unlock() {
	dw, _ := os.Getwd()
	os.Remove(filepath.Join(dw, ".lock"))
}
