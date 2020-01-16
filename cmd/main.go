package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	stypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tjan147/mstore"
)

const (
	name     = "sample"
	scale    = 1000000
	interval = 1000
	cover    = 100

	padding = "vKXNi0SN9u9wt78TX6R7kVwg5phtPhNk77bXerwcR0eNBT+v5fNzjuK4CAyXASW4iC3h5a5kUTfI8LquZszYBas1xpxZJuvZtnhAqbhspe07EAipTOtil5VVPXYp18oxgC9G1ABttakubzRBaDzMFcMukumJo5po9IsAaEMuJCu5EzUaiAJQVROtb7HkmawnJnCtWlxqwJtjNxAH5/SBBzuPBXftY8OzacNLsFMsgUexnSnyYfPwVvd5S/SbdkuS0jLrVh93hmVkN4Gxn6aGouCFmtqMF7wk5l44wiL8yrPHYudsJUq0CYgQxO4lHuKP4NadvilDdoyiobKkFPTjGBrEM9SOCYgKmosU0tw1NJqUTAIwT6p0rIQzjDhCb6sLyfb7qHAE2VlEimQk2xsL7Zpg8o0YZLXiN+6f9l80J5wcA/pK6YrrCuNpkyX7+GRwyQpPMZDJw2RS9sagURR8ePoTtV4xp4QE6RZBQcbLzvu7G1jHiJx8vEQY3w9MNynUBNYrQf5BkSgul+3z2CxWzMScZnqNWAFfLSGeqWIrBPZB9xa+gCq5V5rISjXnK/sQKIjdp7A5ygFwCKALRiJPLdlcvvJyLTz/FzOiqJx1o9r0wNFRYiQ1A1OadT1lwjSB+R604cZqgLrzFt3kvt0TLCZZwyL8YGjVQJVPUrTTaAgel+n3ac7sztAaB5mgpUhtkhcawpAD95L9OtOneYSF+73CdJJc9K3mpadVKBxp9il4Jr2qWkeAJePHJ4yQnYfSo7X05J6mncImwzCVzNeJdARYMgu3pDmNxvwSYA5/v4LJjXvHwZvs7cQCK+KEIT6TMyBc16EAN5VfXajSqNECPVr4jp77OxmWMkxOFoRJVUqKPoJZyl4D8xqYRnUTVFbh4PcID1glYpQgIlHZQUOYehf0i8l9kQlY2mWibJYB+MCb2llz3eCj3jSI9pYpPuUKXx0txYcndSlrmZ1qyc1dIFYs/SAA8dhsdBG2Hn62P8ErKboYLNakp58v0X4h7DqMOPb7VoiVpNfbeCpu4BYmHvgnAwJ76/7p0zBa6MU4YSJ9rDuU8btc9Rc5IF+ICwzXc2H60diVf+4mBWifVqrn3/aIx5efP/dZiq8hXANRu7xVA8pTmv28TgVrbQ6LeGD+S7/ll0aGsLhDBskv4CaiUJF8rWMVzwnS1j3XWLXs/SfgdsY7HRKjVo+yI2X3ljAyHW9esIxmWP9d2X3js97409kDOZ+UWmt1B4cGn+vaUTip4TjqGn4W/rxtXHuPXpV6jxVdfz/Jfd/rS3RBDcnWoQIwBpNp/jK4vMjSH2dWNFlN19C4/i99RXXLO2eLXOzOhmxRq8h7MC+dfx/sYalblg=="
)

func genMockData(idx int) ([]byte, []byte) {
	key := strconv.Itoa(idx)
	val := fmt.Sprintf("%s,%s", idx, padding)
	return []byte(key), []byte(val)
}

func fillData(ckv stypes.CommitKVStore, size int, step int) {
	fmt.Printf("---> Filling database with %d records <---\n", size)

	for i := 0; i < size; i++ {
		k, v := genMockData(i)
		ckv.Set(k, v)

		if i > 0 && (i%step == 0) {
			id := ckv.Commit()
			fmt.Printf("[%7d/%d]: %s\n", i, size, id.String())
		}
	}
	ckv.Commit()
}

func verifyData(ckv stypes.CommitKVStore, size int, count int) {
	fmt.Printf("---> Verifying database with %d/%d samples <---\n", count, size)

	var done bool
	var idx int
	hitTable := make(map[int]bool)
	for i := 0; i < count; i++ {
		done = false
		idx = -1
		for !done {
			idx = rand.Intn(size)
			if _, hit := hitTable[idx]; !hit {
				hitTable[idx] = true
				done = true
			}
		}

		k, v := genMockData(idx)
		refV := ckv.Get(k)

		fmt.Printf("[%3d/%3d]: index:%7d, validation: %t\n", i+1, count, idx, bytes.Compare(v, refV) == 0)
	}
}

func fillRecordsIndexedByTime(ckv stypes.CommitKVStore, start time.Time, step time.Duration, count int) {
	fmt.Printf(
		"---> Fill %d time-indexed records from %s with %0.1fs step <---\n",
		count, start.Format("15:04:05"), step.Seconds())

	i := 1
	t := start
	for ; i <= count; i++ {
		t = t.Add(step)
		ckv.Set(sdk.PrefixEndBytes(sdk.FormatTimeBytes(t)), []byte(strconv.Itoa(i)))
		fmt.Printf("[%d/%d] records inserted, index %s\n", i+1, count, t.Format("15:04:05"))
	}
	ver := ckv.Commit()
	fmt.Printf("%d time-indexed records inserted, %s\n", i+1, ver.String())
}

func pickRecordsFilteredByTime(ckv stypes.CommitKVStore, start, end time.Time) {
	iter := ckv.Iterator(sdk.FormatTimeBytes(start), sdk.InclusiveEndBytes(sdk.FormatTimeBytes(end)))
	count := 0
	for ; iter.Valid(); iter.Next() {
		val, err := strconv.Atoi(string(iter.Value()))
		if err != nil {
			fmt.Printf("[%d]: err: %s\n", count+1, err.Error())
		} else {
			fmt.Printf("[%d]: %d\n", count+1, val)
		}
		count++
	}
	fmt.Printf("%d time-indexed records picked\n", count)
}

func demoCase1() {
	initVer := mstore.InitStore()
	fmt.Printf("- init store: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	fmt.Println("prepare kv")

	start := time.Now()
	fillData(ckv, scale, interval)
	elapse := time.Now().Sub(start)
	fmt.Printf("\n----- DONE -----\nfillData in %d scale costs %.2fs \n", scale, elapse.Seconds())
	fmt.Printf("%s CommitKVStore reversion: %s\n", key.Name(), ckv.LastCommitID().String())

	closeVer := mstore.CloseStore()
	fmt.Printf("close store: %s\n", closeVer.String())
}

func demoCase2() {
	start := time.Now()
	initVer := mstore.InitStore()
	fmt.Printf("- init store: reversion: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	elapse := time.Now().Sub(start)
	fmt.Printf("- prepare kv: cost: %.2fs, reversion: %s\n", elapse.Seconds(), ckv.LastCommitID().String())

	start = time.Now()
	verifyData(ckv, scale, cover)
	elapse = time.Now().Sub(start)
	fmt.Printf("\n----- DONE -----\nverifyData with %d/%d samples costs %.2fs \n", cover, scale, elapse.Seconds())
	fmt.Printf("%s CommitKVStore reversion: %s\n", key.Name(), ckv.LastCommitID().String())

	closeVer := mstore.CloseStore()
	fmt.Printf("close store: %s\n", closeVer.String())
}

func demoCase3() {
	initVer := mstore.InitStore()
	fmt.Printf("- init store: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	fmt.Println("prepare kv")

	start := time.Now()
	fillStart := time.Date(2019, 11, 30, 16, 20, 0, 0, time.UTC)
	fillStep, err := time.ParseDuration("1s")
	if err != nil {
		panic(err)
	}
	fillCount := 600
	fillRecordsIndexedByTime(ckv, fillStart, fillStep, fillCount)
	pickStart := time.Date(2019, 11, 30, 16, 25, 0, 0, time.UTC)
	pickEnd := time.Date(2019, 11, 30, 16, 26, 0, 0, time.UTC)
	pickRecordsFilteredByTime(ckv, pickStart, pickEnd)
	elapse := time.Now().Sub(start)
	fmt.Printf("\n----- DONE -----\nfill+pick costs %.2fs \n", elapse.Seconds())
	fmt.Printf("%s CommitKVStore reversion: %s\n", key.Name(), ckv.LastCommitID().String())

	closeVer := mstore.CloseStore()
	fmt.Printf("close store: %s\n", closeVer.String())
}

func demoCase4() {
	initVer := mstore.InitStore()
	fmt.Printf("- init store: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	fmt.Println("prepare kv")

	dKey := []byte("data")
	// 1st stage
	dVal := []byte(strconv.Itoa(111))
	ckv.Set(dKey, dVal)
	stage1Ver := mstore.GetStoreRecoverSpot()
	ckvID := ckv.LastCommitID()
	fmt.Printf("stage1Ver: %s, ckvVer: %s, dVal = %v\n", stage1Ver.String(), ckvID.String(), dVal)

	// 2nd stage
	dVal = []byte(strconv.Itoa(222))
	ckv.Set(dKey, dVal)
	stage2Ver := mstore.GetStoreRecoverSpot()
	ckvID = ckv.LastCommitID()
	fmt.Printf("stage2Ver: %s, ckvVer: %s, dVal = %v\n", stage2Ver.String(), ckvID.String(), dVal)

	// 3rd stage
	dVal = []byte(strconv.Itoa(333))
	ckv.Set(dKey, dVal)
	stage3Ver := mstore.GetStoreRecoverSpot()
	ckvID = ckv.LastCommitID()
	fmt.Printf("stage3Ver: %s, ckvVer: %s, dVal = %v\n", stage3Ver.String(), ckvID.String(), dVal)

	// latest stage
	lVal := ckv.Get(dKey)
	fmt.Printf("latest stage: lVal = %v\n", lVal)

	// restore the 1st stage
	if err := mstore.LoadStoreRecoverSpot(stage1Ver.Version); err != nil {
		panic(err)
	}
	ckv = mstore.GetCommitKV(key)
	lVal = ckv.Get(dKey)
	ckvID = ckv.LastCommitID()
	fmt.Printf("restore stage1Ver ckvID: %s, lVal = %v\n", ckvID.String(), lVal)

	// restore the 2nd stage
	if err := mstore.LoadStoreRecoverSpot(stage2Ver.Version); err != nil {
		panic(err)
	}
	ckv = mstore.GetCommitKV(key)
	lVal = ckv.Get(dKey)
	ckvID = ckv.LastCommitID()
	fmt.Printf("restore stage2Ver ckvID: %s, lVal = %v\n", ckvID.String(), lVal)

	closeVer := mstore.CloseStore()
	fmt.Printf("close store: %s\n", closeVer.String())
}

func demoCase5Step1() {
	initVer := mstore.InitStore()
	fmt.Printf("- Step1 - init store: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	dKey := []byte("data")

	// 1st stage
	dVal := []byte(strconv.Itoa(111))
	ckv.Set(dKey, dVal)
	stage1Ver := mstore.GetStoreRecoverSpot()
	fmt.Printf("stage1Ver: %s, dVal = %s\n", stage1Ver.String(), string(dVal))

	// 2nd stage
	dVal = []byte(strconv.Itoa(222))
	ckv.Set(dKey, dVal)
	stage2Ver := mstore.GetStoreRecoverSpot()
	fmt.Printf("stage2Ver: %s, dVal = %s\n", stage2Ver.String(), string(dVal))

	// 3rd stage
	dVal = []byte(strconv.Itoa(333))
	ckv.Set(dKey, dVal)
	stage3Ver := mstore.GetStoreRecoverSpot()
	fmt.Printf("stage3Ver: %s, dVal = %s\n", stage3Ver.String(), string(dVal))

	// restore the 1st stage
	if err := mstore.LoadStoreRecoverSpot(stage1Ver.Version); err != nil {
		panic(err)
	}
	ckv = mstore.GetCommitKV(key)
	lVal := ckv.Get(dKey)
	fmt.Printf("restore stage1Ver: lVal = %s\n", string(lVal))

	mstore.CloseStore()
}

func demoCase5Step2() {
	initVer := mstore.InitStore()
	fmt.Printf("- Step2 - init store: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	dKey := []byte("data")
	lVal := ckv.Get(dKey)
	fmt.Printf("reload lVal = %s\n", string(lVal))

	mstore.CloseStore()
}

func demoCase6Step1() {
	initVer := mstore.InitStore()
	fmt.Printf("- Step1 - init store: %s\n", initVer.String())

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	ckv := mstore.GetCommitKV(key)
	dKey := []byte("data")

	// 1st stage
	dVal := []byte(strconv.Itoa(111))
	ckv.Set(dKey, dVal)
	stage1Ver := mstore.GetStoreRecoverSpot()
	fmt.Printf("stage1Ver: %s, dVal = %s\n", stage1Ver.String(), string(dVal))

	// 2nd stage
	dVal = []byte(strconv.Itoa(222))
	ckv.Set(dKey, dVal)
	stage2Ver := mstore.GetStoreRecoverSpot()
	fmt.Printf("stage2Ver: %s, dVal = %s\n", stage2Ver.String(), string(dVal))

	// 3rd stage
	dVal = []byte(strconv.Itoa(333))
	ckv.Set(dKey, dVal)
	stage3Ver := mstore.GetStoreRecoverSpot()
	fmt.Printf("stage3Ver: %s, dVal = %s\n", stage3Ver.String(), string(dVal))

	// restore the 1st stage
	if err := mstore.LoadStoreRecoverSpotForOverwriting(stage1Ver.Version); err != nil {
		panic(err)
	}
	ckv = mstore.GetCommitKV(key)
	lVal := ckv.Get(dKey)
	fmt.Printf("restore stage1Ver lVal = %s\n", string(lVal))

	// change after restore & commit
	dVal = []byte(strconv.Itoa(444))
	ckv.Set(dKey, dVal)
	stageVer := mstore.GetStoreRecoverSpot()
	fmt.Printf("stageVer: %s\n", stageVer.String())

	mstore.CloseStore()
}

func demoCase6Step2() {
	demoCase5Step2()
}

func demoCase7() {
	mstore.InitStore()

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)

	// 1st stage, write directly
	mstore.GetCommitKV(key).Set([]byte("111"), []byte(strconv.Itoa(111)))

	// 2nd stage, cache changes but do not write
	mstore.GetCacheKV(key).Set([]byte("222"), []byte(strconv.Itoa(222)))

	// to valid if the cache is written
	mstore.GetStoreRecoverSpot()

	// 3rd stage, cache changes and write
	cache := mstore.GetCacheKV(key)
	cache.Set([]byte("333"), []byte(strconv.Itoa(333)))
	cache.Write()
	mstore.GetStoreRecoverSpot()

	// reopen store & validate result
	ckv := mstore.GetCommitKV(key)

	fmt.Printf("%s: %v\n", "111", ckv.Get([]byte("111")))
	fmt.Printf("%s: %v\n", "222", ckv.Get([]byte("222")))
	fmt.Printf("%s: %v\n", "333", ckv.Get([]byte("333")))

	mstore.CloseStore()
}

func demoCase8Step1() {
	mstore.InitStore()

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)

	// 1st stage, write directly
	mstore.GetCommitKV(key).Set([]byte("111"), []byte(strconv.Itoa(111)))
	fmt.Printf("1: %v\n", mstore.GetStoreRecoverSpot())

	mstore.GetCommitKV(key).Set([]byte("222"), []byte(strconv.Itoa(222)))
	fmt.Printf("2: %v\n", mstore.GetStoreRecoverSpot())

	mstore.CloseStore()
}

func demoCase8Step2() {
	mstore.InitStore()

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)

	mstore.LoadStoreRecoverSpotForOverwriting(1)
	fmt.Printf("reset 1: %v\n", mstore.GetStoreRecoverSpot())

	mstore.CloseStore()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	demoCase := flag.Int("case", 0, "select a demo case to run")
	demoCaseStep := flag.Int("step", 1, "select a specific step of the demo case to run")
	flag.Parse()

	switch *demoCase {
	case 1:
		demoCase1()
	case 2:
		demoCase2()
	case 3:
		demoCase3()
	case 4:
		demoCase4()
	case 5:
		switch *demoCaseStep {
		case 1:
			demoCase5Step1()
		case 2:
			demoCase5Step2()
		default:
			fmt.Println("error: invalid demo case step selection")
			fmt.Println("usage: cmd -case 5 -step (1|2)")
		}
	case 6:
		switch *demoCaseStep {
		case 1:
			demoCase6Step1()
		case 2:
			demoCase6Step2()
		default:
			fmt.Println("error: invalid demo case step selection")
			fmt.Println("usage: cmd -case 6 -step (1|2)")
		}
	case 7:
		demoCase7()
	case 8:
		switch *demoCaseStep {
		case 1:
			demoCase8Step1()
		case 2:
			demoCase8Step2()
		default:
			fmt.Println("no such step")
		}
	default:
		fmt.Println("error: invalid demo case selection")
	}
}
