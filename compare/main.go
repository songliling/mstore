package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	cstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/db"
	"math/rand"
	"time"
)

const (
	GoLevelDB = "golevel"
	IavlDB    = "iavl"
	Round     = 10
)

func random(max, min int64) int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Int63n(max-min) + min
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func CreateLevelDB(size int64) *db.GoLevelDB {
	newDB, err := db.NewGoLevelDB(fmt.Sprintf("leveldb_%d", size), "")
	if err != nil {
		panic(err)
	}
	return newDB
}

func CreateIavlDB(size int64) (sdk.KVStore, sdk.CommitMultiStore) {
	levelDB, err := db.NewGoLevelDB(fmt.Sprintf("iavl_%d", size), "")
	if err != nil {
		panic(err)
	}
	cms := cstore.NewCommitMultiStore(levelDB)
	cms.SetPruning(cstore.PruneSyncable)

	storeKey := sdk.NewKVStoreKey(fmt.Sprintf("iavl_%d", size))
	cms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, levelDB)
	if err = cms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	return cms.GetKVStore(storeKey), cms
}

func leveldbMockData(size int64) {
	newDB := CreateLevelDB(size)
	defer newDB.Close()
	for i := int64(0); i < size; i++ {
		key := Int64ToBytes(i)
		value := ed25519.GenPrivKey().PubKey().Address().Bytes()
		newDB.Set(key, value)
	}
}

func iavlMockData(size int64) {
	store, cms := CreateIavlDB(size)
	for i := int64(0); i < size; i++ {
		key := Int64ToBytes(i)
		value := ed25519.GenPrivKey().PubKey().Address().Bytes()
		store.Set(key, value)
	}
	cms.Commit()
}

func generateDBMockData(size int64, backendType string) {
	switch backendType {
	case GoLevelDB:
		leveldbMockData(size)
	case IavlDB:
		iavlMockData(size)
	default:
		panic("invalid db type")
	}
}

func levelDB_GetKey_Time(size int64) {
	levelDB := CreateLevelDB(size)
	defer levelDB.Close()
	start := time.Now()
	for i := 0; i < Round; i++ {
		keyNum := random(size, 0)
		levelDB.Get(Int64ToBytes(int64(keyNum)))
	}
	end := time.Since(start)
	ret := end / time.Duration(Round)
	fmt.Printf("demo level db size %d get %d key average time is %s", size, Round, ret)
}

func iavlDB_GetKey_Time(size int64) {
	iavlDB, _ := CreateIavlDB(size)
	start := time.Now()
	for i := 0; i < Round; i++ {
		keyNum := random(size, 0)
		iavlDB.Get(Int64ToBytes(int64(keyNum)))
	}
	end := time.Since(start)
	ret := end / time.Duration(Round)
	fmt.Printf("demo iavl db size %d get %d key average time is %s", size, Round, ret)
}

func calTime(size int64, backendType string) {
	switch backendType {
	case GoLevelDB:
		levelDB_GetKey_Time(size)
	case IavlDB:
		iavlDB_GetKey_Time(size)
	default:
		panic("invalid db type")
	}
}

func demoCase1(flag bool) {
	size := int64(1e5)
	if flag {
		calTime(size, GoLevelDB)
		return
	}
	generateDBMockData(size, GoLevelDB)
}

func demoCase2(flag bool) {
	size := int64(1e6)
	if flag {
		calTime(size, GoLevelDB)
		return
	}
	generateDBMockData(size, GoLevelDB)
}

func demoCase3(flag bool) {
	size := int64(1e5)
	if flag {
		calTime(size, IavlDB)
		return
	}
	generateDBMockData(1e5, IavlDB)
}

func demoCase4(flag bool) {
	size := int64(1e6)
	if flag {
		calTime(size, IavlDB)
		return
	}
	generateDBMockData(1e6, IavlDB)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	demoCase := flag.Int("case", 0, "select a demo case run")
	time := flag.Bool("time", false, "get 100 key Average time")
	flag.Parse()

	switch *demoCase {
	case 1:
		demoCase1(*time)
	case 2:
		demoCase2(*time)
	case 3:
		demoCase3(*time)
	case 4:
		demoCase4(*time)
	default:
		panic("error case number")
	}
}
