package mstore

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store"
	stypes "github.com/cosmos/cosmos-sdk/store/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	dbHome = "."
	dbName = "sample"
)

var (
	db  dbm.DB
	cms stypes.CommitMultiStore
)

func InitStore() {
	var err error
	db, err = dbm.NewGoLevelDB(dbName, dbHome)
	if err != nil {
		panic(err)
	}
	fmt.Println("backend db ready")

	cms = store.NewCommitMultiStore(db)
	fmt.Println("cms ready")
}

func CloseStore() (id stypes.CommitID) {
	id = cms.Commit()
	db.Close()
	cms = nil
	return
}

func CreateNewCommitKV(key stypes.StoreKey) {
	cms.MountStoreWithDB(key, stypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
}

func GetCommitKV(key stypes.StoreKey) stypes.CommitKVStore {
	return cms.GetCommitKVStore(key)
}
