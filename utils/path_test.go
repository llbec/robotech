package utils

import (
	"fmt"
	"testing"
)

func Test_HomePath(t *testing.T) {
	home, err := Home()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(home)
}

func Test_FileList(t *testing.T) {
	files, err := FileList("db", "go")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(files)
}
