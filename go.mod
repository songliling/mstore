module github.com/tjan147/mstore

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.38.1
	github.com/tendermint/iavl v0.13.0
	github.com/tendermint/tendermint v0.31.5
	github.com/tendermint/tm-db v0.5.1
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/LambdaIM/cosmos-sdk v0.35.11
	github.com/tendermint/iavl => github.com/LambdaIM/iavl v0.13.2-dev1
	github.com/tendermint/tendermint => github.com/LambdaIM/tendermint v0.31.5-fix1
)
