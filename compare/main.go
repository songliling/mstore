package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	cstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tm-db"
	"math/rand"
	"time"
)

const (
	GoLevelDB = "golevel"
	IavlDB    = "iavl"
	Round     = 10
	Version   = "0.13.0"
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

func CreateLevelDB(size int64, prefix string) *db.GoLevelDB {
	newDB, err := db.NewGoLevelDB(fmt.Sprintf("leveldb%s_%d", prefix, size), "")
	if err != nil {
		panic(err)
	}
	return newDB
}

func CreateIavlDB(size int64, prefix string) (sdk.KVStore, sdk.CommitMultiStore) {
	levelDB, err := db.NewGoLevelDB(fmt.Sprintf("iavl%s_%s_%d", prefix, Version, size), "")
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

func leveldbMockData(size int64, long bool) {
	if !long {
		newDB := CreateLevelDB(size, "")
		defer newDB.Close()
		for i := int64(0); i < size; i++ {
			key := Int64ToBytes(i)
			value := ed25519.GenPrivKey().PubKey().Address().Bytes()
			newDB.Set(key, value)
		}
	} else {
		newDB := CreateLevelDB(size, "long")
		defer newDB.Close()
		for i := int64(0); i < size; i++ {
			key := append(Int64ToBytes(i), []byte("Store")...)
			value := NewSellOrder(size)
			newDB.Set(key, value)
		}
	}

}

func iavlMockData(size int64, long bool) {
	if !long {
		store, cms := CreateIavlDB(size, "")
		for i := int64(0); i < size; i++ {
			key := Int64ToBytes(i)
			value := ed25519.GenPrivKey().PubKey().Address().Bytes()
			store.Set(key, value)
			if i%100 == 0 {
				fmt.Printf("commit index %d\n", i)
				cms.Commit()
			}
		}
	} else {
		store, cms := CreateIavlDB(size, "long")
		for i := int64(0); i < size; i++ {
			key := append(Int64ToBytes(i), []byte("Store")...)
			value := NewSellOrder(size)
			store.Set(key, value)
		}
		cms.Commit()
	}
}

func generateDBMockData(size int64, backendType string, long bool) {
	switch backendType {
	case GoLevelDB:
		leveldbMockData(size, long)
	case IavlDB:
		iavlMockData(size, long)
	default:
		panic("invalid db type")
	}
}

func levelDB_GetKey_Time(size int64, long bool) {
	if !long {
		levelDB := CreateLevelDB(size, "")
		defer levelDB.Close()
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			levelDB.Get(Int64ToBytes(int64(keyNum)))
		}
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo level db size %d get %d key average time is %s\n", size, Round, ret.String())
	} else {
		levelDB := CreateLevelDB(size, "long")
		defer levelDB.Close()
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			levelDB.Get(append(Int64ToBytes(int64(keyNum)), []byte("Store")...))
		}
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo level db size %d get %d key average time is %s\n", size, Round, ret.String())
	}
}

func levelDB_SetKey_Time(size int64, long bool) {
	if !long {
		levelDB := CreateLevelDB(size, "")
		defer levelDB.Close()
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			levelDB.Set(Int64ToBytes(int64(keyNum)), ed25519.GenPrivKey().PubKey().Bytes())
		}
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo level db size %d set %d key average time is %s\n", size, Round, ret.String())
	} else {
		levelDB := CreateLevelDB(size, "long")
		defer levelDB.Close()
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			levelDB.Set(append(Int64ToBytes(int64(keyNum)), []byte("Store")...),
				NewSellOrder(size))
		}
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo level db size %d set %d key average time is %s\n", size, Round, ret.String())
	}
}

func iavlDB_GetKey_Time(size int64, long bool) {
	if !long {
		iavlDB, _ := CreateIavlDB(size, "")
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			iavlDB.Get(Int64ToBytes(int64(keyNum)))
		}
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo iavl db size %d get %d key average time is %s\n", size, Round, ret.String())
	} else {
		iavlDB, _ := CreateIavlDB(size, "long")
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			iavlDB.Get(append(Int64ToBytes(int64(keyNum)), []byte("Store")...))
		}
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo iavl db size %d get %d key average time is %s\n", size, Round, ret.String())
	}
}

func iavlDB_SetKey_Time(size int64, long bool) {
	if !long {
		iavlDB, cms := CreateIavlDB(size, "")
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			iavlDB.Set(Int64ToBytes(int64(keyNum)), ed25519.GenPrivKey().PubKey().Bytes())
		}
		cms.Commit()
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo iavl db size %d set %d key average time is %s\n", size, Round, ret.String())
	} else {
		iavlDB, cms := CreateIavlDB(size, "long")
		start := time.Now()
		for i := 0; i < Round; i++ {
			keyNum := random(size, 0)
			iavlDB.Set(append(Int64ToBytes(int64(keyNum)), []byte("Store")...),
				NewSellOrder(size))
		}
		cms.Commit()
		end := time.Since(start)
		ret := end / time.Duration(Round)
		fmt.Printf("demo iavl db size %d set %d key average time is %s\n", size, Round, ret.String())
	}
}

func calTime(size int64, backendType string, long bool) {
	switch backendType {
	case GoLevelDB:
		levelDB_GetKey_Time(size, long)
	case IavlDB:
		iavlDB_GetKey_Time(size, long)
	default:
		panic("invalid db type")
	}
}

func setKey(size int64, backendType string, long bool) {
	switch backendType {
	case GoLevelDB:
		levelDB_SetKey_Time(size, long)
	case IavlDB:
		iavlDB_SetKey_Time(size, long)
	default:
		panic("invalid db type")
	}
}

func demoCase1(cal bool, set bool, long bool) {
	size := int64(1e5)
	if cal {
		calTime(size, GoLevelDB, long)
		return
	}
	if set {
		setKey(size, GoLevelDB, long)
		return
	}
	generateDBMockData(size, GoLevelDB, long)
}

func demoCase2(cal bool, set bool, long bool) {
	size := int64(1e6)
	if cal {
		calTime(size, GoLevelDB, long)
		return
	}
	if set {
		setKey(size, GoLevelDB, long)
		return
	}
	generateDBMockData(size, GoLevelDB, long)
}

func demoCase3(cal bool, set bool, long bool) {
	size := int64(1e5)
	if cal {
		calTime(size, IavlDB, long)
		return
	}
	if set {
		setKey(size, IavlDB, long)
		return
	}
	generateDBMockData(1e5, IavlDB, long)
}

func demoCase4(cal bool, set bool, long bool) {
	size := int64(1e6)
	if cal {
		calTime(size, IavlDB, long)
		return
	}
	if set {
		setKey(size, IavlDB, long)
		return
	}
	generateDBMockData(1e6, IavlDB, long)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	demoCase := flag.Int("case", 0, "select a demo case run")
	time := flag.Bool("time", false, fmt.Sprintf("get %d key Average time", Round))
	set := flag.Bool("set", false, fmt.Sprintf("set %d key Average time", Round))
	long := flag.Bool("long", false, fmt.Sprintf("set long key"))
	flag.Parse()

	switch *demoCase {
	case 1:
		demoCase1(*time, *set, *long)
	case 2:
		demoCase2(*time, *set, *long)
	case 3:
		demoCase3(*time, *set, *long)
	case 4:
		demoCase4(*time, *set, *long)
	default:
		panic("error case number")
	}
}
