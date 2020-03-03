package main

import (
	"encoding/hex"
	"fmt"
	cstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tm-db"
	tdb "github.com/tendermint/tm-db"
	"math/rand"
	mrand "math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

/**
1. 生成一个db，两个db实例，cache size不同
*/
const (
	KeyPrefix     = "Store"
	eachStepScale = 10000
	valueLen      = 24
)

func CreateIavlDB(name string, cacheSize int) (sdk.KVStore, sdk.CommitMultiStore, *tdb.GoLevelDB) {

	levelDB, err := db.NewGoLevelDB(name, "")

	if err != nil {
		panic(err)
	}
	cms := cstore.NewCommitMultiStore(levelDB, cacheSize)
	cms.SetPruning(cstore.PruneSyncable)

	storeKey := sdk.NewKVStoreKey(fmt.Sprintf("iavl_%d", cacheSize))
	cms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, levelDB)
	if err = cms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	return cms.GetKVStore(storeKey), cms, levelDB
}

func cleanupDBDir(name, dir string) {
	if err := os.RemoveAll(filepath.Join(dir, name) + ".db"); err != nil {
		panic(err)
	}
}

type dbTestFunc = func(stores []sdk.KVStore, cmss []sdk.CommitMultiStore, steps int) string

func generateData(stores []sdk.KVStore, cmss []sdk.CommitMultiStore, steps int) string {
	start := time.Now()
	num := len(stores)
	dbKeyNum := eachStepScale / num
	for j := 0; j < num; j++ {
		for i := j * dbKeyNum; i < (j+1)*dbKeyNum; i++ {
			key := []byte(KeyPrefix + strconv.Itoa(i))
			val := make([]byte, valueLen)
			if _, err := rand.Read(val); err != nil {
				panic(err)
			}
			stores[j].Set(key, []byte(hex.EncodeToString(val)))
			if i%100 == 0 {
				cmss[j].Commit()
			}
		}
	}
	return fmt.Sprintf("cost %s", time.Since(start))
}

func testDB(name string, cacheSize int, stepQuence []string, suite map[string]dbTestFunc, dbNum int) {
	stores := make([]sdk.KVStore, dbNum)
	cmss := make([]sdk.CommitMultiStore, dbNum)
	dbs := make([]*tdb.GoLevelDB, dbNum)
	for i := 0; i < dbNum; i++ {
		stores[i], cmss[i], dbs[i] = CreateIavlDB(fmt.Sprintf("%s_%d", name, i), cacheSize)
		defer cleanupDBDir(fmt.Sprintf("%s_%d", name, i), "")
	}

	for _, mKey := range stepQuence {
		fmt.Printf("%s %s: %s\n", name, mKey, suite[mKey](stores, cmss, 1000))
	}

	for _, db := range dbs {
		db.Close()
	}
}

func testSet(stores []sdk.KVStore, cmss []sdk.CommitMultiStore, steps int) (out string) {
	start := time.Now()
	num := len(stores)
	numKey := steps / num
	for i := 0; i < num; i++ {
		for step := 0; step < numKey; step++ {
			key := []byte(KeyPrefix + strconv.Itoa(mrand.Intn(eachStepScale/num)+i*(eachStepScale/num)))
			value := make([]byte, valueLen)
			rand.Read(value)
			stores[i].Set(key, value)
		}
		cmss[i].Commit()
	}
	out += fmt.Sprintf("cost %s", time.Since(start))
	return
}

func testGet(stores []sdk.KVStore, cmss []sdk.CommitMultiStore, steps int) (out string) {
	// update set
	start := time.Now()
	num := len(stores)
	for i := 0; i < num; i++ {
		for step := 0; step < steps/num; step++ {
			key := []byte(KeyPrefix + strconv.Itoa(mrand.Intn(eachStepScale/num)+i*(eachStepScale/num)))
			stores[i].Get(key)
		}
	}
	out += fmt.Sprintf("cost %s", time.Since(start))
	return
}

// find a way to warp this & make it work
// TODO: fix this
func reopen(name, dir string, dtype tdb.BackendType, scale int, db tdb.DB) string {
	start := time.Now()
	db.Close()

	db = tdb.NewDB(name, dtype, dir)
	return fmt.Sprintf("reopen, %dms", time.Since(start).Milliseconds())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// diff cache size instance
	stepQ := []string{"gen", "set", "get"}
	cacheSuite := map[string]dbTestFunc{
		"gen": generateData,
		"set": testSet,
		"get": testGet,
	}
	runtime.GC()
	size := 10000
	testDB("cache_1W", size, stepQ, cacheSuite, 1)

	runtime.GC()
	size = 100000
	testDB("cache_10W", size, stepQ, cacheSuite, 1)

	fmt.Println("---------------------------------------------")
	// key Dispersed diff db

	stepQ = []string{"gen", "set", "get"}
	cacheSuite = map[string]dbTestFunc{
		"gen": generateData,
		"set": testSet,
		"get": testGet,
	}
	size = 100000
	testDB("cache_10W_4db", size, stepQ, cacheSuite, 4)
	testDB("cache_10W_8db", size, stepQ, cacheSuite, 8)
}
