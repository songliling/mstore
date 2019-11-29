package main

import (
	"fmt"
	"strconv"
	"time"

	stypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/tjan147/mstore"
)

const (
	name     = "sample"
	scale    = 1000000
	interval = 1000

	padding = "vKXNi0SN9u9wt78TX6R7kVwg5phtPhNk77bXerwcR0eNBT+v5fNzjuK4CAyXASW4iC3h5a5kUTfI8LquZszYBas1xpxZJuvZtnhAqbhspe07EAipTOtil5VVPXYp18oxgC9G1ABttakubzRBaDzMFcMukumJo5po9IsAaEMuJCu5EzUaiAJQVROtb7HkmawnJnCtWlxqwJtjNxAH5/SBBzuPBXftY8OzacNLsFMsgUexnSnyYfPwVvd5S/SbdkuS0jLrVh93hmVkN4Gxn6aGouCFmtqMF7wk5l44wiL8yrPHYudsJUq0CYgQxO4lHuKP4NadvilDdoyiobKkFPTjGBrEM9SOCYgKmosU0tw1NJqUTAIwT6p0rIQzjDhCb6sLyfb7qHAE2VlEimQk2xsL7Zpg8o0YZLXiN+6f9l80J5wcA/pK6YrrCuNpkyX7+GRwyQpPMZDJw2RS9sagURR8ePoTtV4xp4QE6RZBQcbLzvu7G1jHiJx8vEQY3w9MNynUBNYrQf5BkSgul+3z2CxWzMScZnqNWAFfLSGeqWIrBPZB9xa+gCq5V5rISjXnK/sQKIjdp7A5ygFwCKALRiJPLdlcvvJyLTz/FzOiqJx1o9r0wNFRYiQ1A1OadT1lwjSB+R604cZqgLrzFt3kvt0TLCZZwyL8YGjVQJVPUrTTaAgel+n3ac7sztAaB5mgpUhtkhcawpAD95L9OtOneYSF+73CdJJc9K3mpadVKBxp9il4Jr2qWkeAJePHJ4yQnYfSo7X05J6mncImwzCVzNeJdARYMgu3pDmNxvwSYA5/v4LJjXvHwZvs7cQCK+KEIT6TMyBc16EAN5VfXajSqNECPVr4jp77OxmWMkxOFoRJVUqKPoJZyl4D8xqYRnUTVFbh4PcID1glYpQgIlHZQUOYehf0i8l9kQlY2mWibJYB+MCb2llz3eCj3jSI9pYpPuUKXx0txYcndSlrmZ1qyc1dIFYs/SAA8dhsdBG2Hn62P8ErKboYLNakp58v0X4h7DqMOPb7VoiVpNfbeCpu4BYmHvgnAwJ76/7p0zBa6MU4YSJ9rDuU8btc9Rc5IF+ICwzXc2H60diVf+4mBWifVqrn3/aIx5efP/dZiq8hXANRu7xVA8pTmv28TgVrbQ6LeGD+S7/ll0aGsLhDBskv4CaiUJF8rWMVzwnS1j3XWLXs/SfgdsY7HRKjVo+yI2X3ljAyHW9esIxmWP9d2X3js97409kDOZ+UWmt1B4cGn+vaUTip4TjqGn4W/rxtXHuPXpV6jxVdfz/Jfd/rS3RBDcnWoQIwBpNp/jK4vMjSH2dWNFlN19C4/i99RXXLO2eLXOzOhmxRq8h7MC+dfx/sYalblg=="
)

func genKV(idx int) ([]byte, []byte) {
	key := strconv.Itoa(idx)
	val := fmt.Sprintf("%s,%s", idx, padding)
	return []byte(key), []byte(val)
}

func fillMockData(ckv stypes.CommitKVStore, size int, step int) {
	fmt.Printf("---> Filling database with %d records <---\n", size)

	for i := 0; i < size; i++ {
		k, v := genKV(i)
		ckv.Set(k, v)

		if i > 0 && (i%step == 0) {
			id := ckv.Commit()
			fmt.Printf("[%7d/%d]: %s\n", i, size, id.String())
		}
	}
}

func main() {
	mstore.InitStore()
	fmt.Println("init store")

	key := stypes.NewKVStoreKey(name)
	mstore.CreateNewCommitKV(key)
	fmt.Println("new kv")
	ckv := mstore.GetCommitKV(key)

	start := time.Now()
	fillMockData(ckv, scale, interval)
	elapse := time.Now().Sub(start)
	fmt.Printf("\n----- DONE -----\nfillMockData costs %.2fs \n", elapse.Seconds())

	closeVer := mstore.CloseStore()
	fmt.Printf("close store: %s\n", closeVer.String())
}
