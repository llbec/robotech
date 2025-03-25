package storage

import (
	"log"
	"math/big"
	"path/filepath"
	"robotech/applications/marginv1/config"
	"robotech/applications/marginv1/types"
	"robotech/utils"
	"robotech/utils/db"

	"github.com/ethereum/go-ethereum/common"
)

var (
	chainDb         db.DB
	poolsDb         db.DB
	positionsDbList map[string]db.DB
	dir             string
	height          uint64
)

func init() {
	dir = filepath.Join(config.GetWorkDir(), ".storage")

	if chainDb == nil {
		chainDb = db.NewDB(
			db.LevelDBBackend,
			"chain.db",
			filepath.Join(dir, "chain"),
		)
	}

	if poolsDb == nil {
		poolsDb = db.NewDB(
			db.LevelDBBackend,
			"pools.db",
			filepath.Join(dir, "pools"),
		)
	}

	if positionsDbList == nil {
		positionsDbList = make(map[string]db.DB)

		iter := poolsDb.NewIterator()
		for {
			memeToken := common.BytesToAddress(iter.Key())
			dbFile := memeToken.String() + ".db"
			//debtors[addr.String()] = big.NewInt(int64(utils.BytesToUInt64(iter.Value())))
			positionsDbList[memeToken.String()] = db.NewDB(
				db.LevelDBBackend,
				dbFile,
				filepath.Join(dir, "positions"),
			)
			if !iter.Next() {
				break
			}
		}
	}

	ok, err := chainDb.Has([]byte("last_block_number"))
	if err != nil {
		panic(err)
	}
	if ok {
		h, err := chainDb.Get([]byte("last_block_number"))
		if err != nil {
			panic(err)
		}
		height = utils.BytesToUInt64(h)
	} else {
		height = config.GetBlockNumber()
	}
	log.Printf("storage initialized")
}

func Close() {
	chainDb.Close()
	poolsDb.Close()
	for _, db := range positionsDbList {
		db.Close()
	}
}

func SetLastBlockNumber(h uint64) error {
	err := chainDb.Put([]byte("last_block_number"), utils.UInt64ToBytes(h))
	if err != nil {
		return err
	}
	height = h
	return nil
}

func GetLastBlockNumber() uint64 {
	return height
}

func GetPools() []common.Address {
	pools := make([]common.Address, 0)
	iter := poolsDb.NewIterator()
	for {
		memeToken := common.BytesToAddress(iter.Key())
		pools = append(pools, memeToken)
		if !iter.Next() {
			break
		}
	}
	return pools
}

func HasPool(pool common.Address) (bool, error) {
	return poolsDb.Has(pool.Bytes())
}

func AddPool(pool common.Address) error {
	_, ok := positionsDbList[pool.String()]
	if ok {
		return nil
	}
	err := poolsDb.Put(pool.Bytes(), utils.UInt64ToBytes(0))
	if err != nil {
		return err
	}
	dbFile := pool.String() + ".db"
	positionsDbList[pool.String()] = db.NewDB(
		db.LevelDBBackend,
		dbFile,
		filepath.Join(dir, "positions"),
	)
	log.Printf("add pool %s", pool.String())
	return nil
}

func AddPosition(memeToken, account common.Address, position types.Position) error {
	err := AddPool(memeToken)
	if err != nil {
		log.Printf("add pool %s error %s", memeToken.String(), err.Error())
	}
	log.Printf("add position at pool(%s) %s-%v", memeToken.String(), account.String(), position.PositionId)
	return positionsDbList[memeToken.String()].Put(account.Bytes(), position.PositionId.Bytes())
}

func RemovePosition(memeToken, account common.Address) error {
	_, ok := positionsDbList[memeToken.String()]
	if !ok {
		return nil
	}
	log.Printf("remove position %s %s", memeToken.String(), account.String())
	return positionsDbList[memeToken.String()].Delete(account.Bytes())
}

func GetPositions(memeToken common.Address) []types.Position {
	clearZero := false
	positionDb, ok := positionsDbList[memeToken.String()]
	if !ok {
		return []types.Position{}
	}
	iter := positionDb.NewIterator()
	positions := make([]types.Position, 0)
	for {
		//debtors[addr.String()] = big.NewInt(int64(utils.BytesToUInt64(iter.Value())))
		if common.BytesToAddress(iter.Key()) == common.HexToAddress("0x0000000000000000000000000000000000000000") {
			log.Printf("clear zero %v %v", common.BytesToAddress(iter.Key()), big.NewInt(0).SetBytes(iter.Value()))
			clearZero = true
			if !iter.Next() {
				break
			}
			continue
		}
		positions = append(positions, types.Position{
			MemeToken:  memeToken,
			Account:    common.BytesToAddress(iter.Key()),
			PositionId: big.NewInt(0).SetBytes(iter.Value()),
		})
		if !iter.Next() {
			break
		}
	}
	if clearZero {
		positionDb.Delete(common.HexToAddress("0x0000000000000000000000000000000000000000").Bytes())
	}
	//log.Printf("get positions %s %d", memeToken.String(), len(positions))
	iter.Release()
	//log.Printf("get positions %v", positions)
	return positions
}
