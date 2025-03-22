package db

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	// degradationWarnInterval specifies how often warning should be printed if the
	// leveldb database cannot keep up with requested writes.
	degradationWarnInterval = time.Minute

	// minCache is the minimum amount of memory in megabytes to allocate to leveldb
	// read and write caching, split half and half.
	minCache = 16

	// minHandles is the minimum number of files handles to allocate to the open
	// database files.
	minHandles = 16

	// metricsGatheringInterval specifies the interval to retrieve leveldb database
	// compaction, io and pause stats to report to the user.
	metricsGatheringInterval = 3 * time.Second

	// LevelDBBackend leveldb
	LevelDBBackend = "leveldb"
)

func init() {
	dbCreator := func(name string, dir string) (DB, error) {
		return NewLevelDB(name, dir)
	}
	registerDBCreator(LevelDBBackend, dbCreator, false)
}

// LevelDB is a persistent key-value store. Apart from basic data storage
// functionality it also supports batch writes and iterating over the keyspace in
// binary-alphabetical order.
type LevelDB struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance

	quitLock sync.Mutex // Mutex protecting the quit channel access
}

var _ DB = (*LevelDB)(nil)

// New returns a wrapped LevelDB object.
/*func New(name string, dir string, cache int, handles int) (*LevelDB, error) {
	// Ensure we have some minimal caching and file guarantees
	if cache < minCache {
		cache = minCache
	}
	if handles < minHandles {
		handles = minHandles
	}

	return NewLevelDBWithOpts(name, dir, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
}*/

// NewLevelDB returns a wrapped LevelDB object.
func NewLevelDB(name string, dir string) (*LevelDB, error) {
	return NewLevelDBWithOpts(name, dir, nil)
}

// NewLevelDBWithOpts returns a wrapped LevelDB object.
func NewLevelDBWithOpts(name string, dir string, o *opt.Options) (*LevelDB, error) {
	dbPath := filepath.Join(dir, name+".db")
	db, err := leveldb.OpenFile(dbPath, o)
	if err != nil {
		return nil, err
	}
	database := &LevelDB{
		db: db,
		fn: dbPath,
	}
	return database, nil
}

// Close stops the metrics collection, flushes any pending data to disk and closes
// all io accesses to the underlying key-value store.
func (db *LevelDB) Close() error {
	db.quitLock.Lock()
	defer db.quitLock.Unlock()

	return db.db.Close()
}

// Has retrieves if a key is present in the key-value store.
func (db *LevelDB) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get retrieves the given key if it's present in the key-value store.
func (db *LevelDB) Get(key []byte) ([]byte, error) {
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Put inserts the given value into the key-value store.
func (db *LevelDB) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// Delete removes the key from the key-value store.
func (db *LevelDB) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *LevelDB) NewBatch() Batch {
	return &leveldbBatch{
		db: db.db,
		b:  new(leveldb.Batch),
	}
}

// NewIterator creates a binary-alphabetical iterator over the entire keyspace
// contained within the leveldb database.
func (db *LevelDB) NewIterator() Iterator {
	return db.db.NewIterator(new(util.Range), nil)
}

// NewIteratorWithStart creates a binary-alphabetical iterator over a subset of
// database content starting at a particular initial key (or after, if it does
// not exist).
func (db *LevelDB) NewIteratorWithStart(start []byte) Iterator {
	return db.db.NewIterator(&util.Range{Start: start}, nil)
}

// NewIteratorWithPrefix creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix.
func (db *LevelDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	return db.db.NewIterator(util.BytesPrefix(prefix), nil)
}

// Stat returns a particular internal stat of the database.
func (db *LevelDB) Stat(property string) (string, error) {
	return db.db.GetProperty(property)
}

// Compact flattens the underlying data store for the given key range. In essence,
// deleted and overwritten versions are discarded, and the data is rearranged to
// reduce the cost of operations needed to access them.
//
// A nil start is treated as a key before all keys in the data store; a nil limit
// is treated as a key after all keys in the data store. If both is nil then it
// will compact entire data store.
func (db *LevelDB) Compact(start []byte, limit []byte) error {
	return db.db.CompactRange(util.Range{Start: start, Limit: limit})
}

// Path returns the path to the database directory.
func (db *LevelDB) Path() string {
	return db.fn
}

/**********************************************************************
                                  batch
**********************************************************************/

// leveldbBatch is a write-only leveldb batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type leveldbBatch struct {
	db   *leveldb.DB
	b    *leveldb.Batch
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *leveldbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *leveldbBatch) Delete(key []byte) error {
	b.b.Delete(key)
	b.size++
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *leveldbBatch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to disk.
func (b *leveldbBatch) Write() error {
	return b.db.Write(b.b, nil)
}

// Reset resets the batch for reuse.
func (b *leveldbBatch) Reset() {
	b.b.Reset()
	b.size = 0
}

// Replay replays the batch contents.
func (b *leveldbBatch) Replay(w KeyValueWriter) error {
	return b.b.Replay(&leveldbReplayer{writer: w})
}

/**********************************************************************
                                replayer
**********************************************************************/

// leveldbReplayer is a small wrapper to implement the correct replay methods.
type leveldbReplayer struct {
	writer  KeyValueWriter
	failure error
}

// Put inserts the given value into the key-value data store.
func (r *leveldbReplayer) Put(key, value []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Put(key, value)
}

// Delete removes the key from the key-value data store.
func (r *leveldbReplayer) Delete(key []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Delete(key)
}
