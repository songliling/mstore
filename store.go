package mstore

import (
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

func InitStore() stypes.CommitID {
	var err error
	db, err = dbm.NewGoLevelDB(dbName, dbHome)
	if err != nil {
		panic(err)
	}

	cms = store.NewCommitMultiStore(db)
	return cms.LastCommitID()
}

func CloseStore() stypes.CommitID {
	status := cms.Commit()
	db.Close()
	cms = nil
	return status
}

func CreateNewCommitKV(key stypes.StoreKey) {
	cms.MountStoreWithDB(key, stypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
}

func GetCommitKV(key stypes.StoreKey) stypes.CommitKVStore {
	return cms.GetCommitKVStore(key)
}
