module github.com/tjan147/mstore

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.38.1
	github.com/tendermint/tendermint v0.33.1
	github.com/tendermint/tm-db v0.4.0
)

replace github.com/cosmos/cosmos-sdk v0.38.1 => ../../go-project/src/github.com/cosmos-sdk

replace github.com/tendermint/iavl v0.13.0 => ../../go-project/src/github.com/iavl
