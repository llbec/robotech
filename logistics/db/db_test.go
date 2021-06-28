package db

import (
	"bytes"
	"encoding/binary"
	"testing"

	"math/rand"
	"time"
)

func benchmarkTraverse(b *testing.B, db DB) {
	b.StopTimer()
	// create db data
	const numItems = int64(1000000)
	internal := map[int64]int64{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < int(numItems); i++ {
		internal[int64(i)] = int64(r.Int())
		db.Put(int642Bytes(int64(i)), int642Bytes(internal[int64(i)]))
	}

	// start to test traversal
	b.StartTimer()

	iter := db.NewIterator()
	defer iter.Release()
	count := 0
	for iter.Next() {
		count++
		idx := bytes2Int64(iter.Key())
		val := bytes2Int64(iter.Value())
		if internal[idx] != val {
			b.Errorf("Expected %v for %v, got %v", internal[idx], idx, val)
		}
	}
	if int64(count) != numItems {
		b.Errorf("Expected length %v, got %v", numItems, count)
	}
}

func benchmarkRandomReadsWrites(b *testing.B, db DB) {
	b.StopTimer()

	// create dummy data
	const numItems = int64(1000000)
	internal := map[int64]int64{}
	for i := 0; i < int(numItems); i++ {
		internal[int64(i)] = int64(0)
	}

	// fmt.Println("ok, starting")
	b.StartTimer()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < b.N; i++ {
		// Write something
		{
			idx := int64(r.Int()) % numItems
			internal[idx]++
			val := internal[idx]
			idxBytes := int642Bytes(int64(idx))
			valBytes := int642Bytes(int64(val))
			//fmt.Printf("Set %X -> %X\n", idxBytes, valBytes)
			db.Put(idxBytes, valBytes)
		}

		// Read something
		{
			idx := int64(r.Int()) % numItems
			valExp := internal[idx]
			idxBytes := int642Bytes(int64(idx))
			valBytes, _ := db.Get(idxBytes)
			//fmt.Printf("Get %X -> %X\n", idxBytes, valBytes)
			if valExp == 0 {
				if !bytes.Equal(valBytes, nil) {
					b.Errorf("Expected %v for %v, got %X", nil, idx, valBytes)
					break
				}
			} else {
				if len(valBytes) != 8 {
					b.Errorf("Expected length 8 for %v, got %X", idx, valBytes)
					break
				}
				valGot := bytes2Int64(valBytes)
				if valExp != valGot {
					b.Errorf("Expected %v for %v, got %v", valExp, idx, valGot)
					break
				}
			}
		}

	}
}

func int642Bytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func bytes2Int64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
