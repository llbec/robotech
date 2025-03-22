package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"time"
)

func cleanupDBDir(dir, name string) {
	err := os.RemoveAll(filepath.Join(dir, name) + ".db")
	if err != nil {
		panic(err)
	}
}

func BenchmarkGoLevelDBRandomReadsWrites(b *testing.B) {
	name := fmt.Sprintf("test_%x", time.Now().Format("2006-01-02 15:04:05"))
	db, err := NewLevelDB(name, "")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		db.Close()
		cleanupDBDir("", name)
	}()

	benchmarkRandomReadsWrites(b, db)
}

func BenchmarkGoLevelDBRandomTraverse(b *testing.B) {
	name := fmt.Sprintf("test_%x", time.Now().Format("2006-01-02 15:04:05"))
	db, err := NewLevelDB(name, "")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		db.Close()
		cleanupDBDir("", name)
	}()

	benchmarkTraverse(b, db)
}
