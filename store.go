package mstore

//
//import (
//	"fmt"
//
//	"github.com/cosmos/cosmos-sdk/store"
//	stypes "github.com/cosmos/cosmos-sdk/store/types"
//	dbm "github.com/tendermint/tm-db"
//)
//
//const (
//	dbHome = "."
//	dbName = "sample"
//)
//
//var (
//	db   dbm.DB
//	fcms stypes.ForkableCommitMultiStore
//)
//
//func InitStore() stypes.CommitID {
//	var err error
//	db, err = dbm.NewGoLevelDB(dbName, dbHome)
//	if err != nil {
//		panic(err)
//	}
//
//	fcms = store.NewForkableCommitMultiStore(db)
//	fcms.SetPruning(stypes.PruneSyncable)
//	return fcms.LastCommitID()
//}
//
//func CloseStore() stypes.CommitID {
//	status := fcms.LastCommitID()
//	db.Close()
//	return status
//}
//
//func CreateNewCommitKV(key stypes.StoreKey) {
//	fcms.MountStoreWithDB(key, stypes.StoreTypeIAVL, db)
//	fcms.LoadLatestVersion()
//}
//
//func GetCacheKV(key stypes.StoreKey) stypes.CacheKVStore {
//	wrapper := fcms.GetKVStore(key).CacheWrap()
//	cacheKV, ok := wrapper.(stypes.CacheKVStore)
//	if !ok {
//		panic(fmt.Errorf("Unsupported StoreType\n"))
//	}
//	return cacheKV
//}
//
//func GetCommitKV(key stypes.StoreKey) stypes.CommitKVStore {
//	return fcms.GetCommitKVStore(key)
//}
//
//func GetStoreRecoverSpot() stypes.CommitID {
//	return fcms.Commit()
//}
//
//func LoadStoreRecoverSpot(rev int64) error {
//	return fcms.LoadVersion(rev)
//}
//
//func LoadStoreRecoverSpotForOverwriting(rev int64) error {
//	return fcms.LoadVersionForOverwriting(rev)
//}
