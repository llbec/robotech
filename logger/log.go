package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Init(filename string) {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err) // 如果创建目录失败，则 panic
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func Info(v ...any)                 { log.Println("[INFO]", v) }
func Infof(format string, v ...any) { log.Printf("[INFO] "+format+"\n", v...) }
func Error(v ...any)                { log.Println("[ERROR]", v) }
func Fatal(v ...any)                { log.Fatal("[FATAL]", v) }
func Debug(v ...any)                { log.Println("[DEBUG]", v) }
