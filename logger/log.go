package logger

import (
	"io"
	"log"
	"os"
)

func Init(filename string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func Info(v ...any)  { log.Println("[INFO]", v) }
func Error(v ...any) { log.Println("[ERROR]", v) }
