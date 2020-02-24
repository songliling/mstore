package main

import (
	cdc "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"time"
)

type SellOrder struct {
	OrderId            string        `json:"orderId"`
	Address            string        `json:"address"`
	Price              sdk.Int       `json:"price"`
	Rate               sdk.Dec       `json:"rate"`
	Amount             sdk.Coins     `json:"amount"`
	SellSize           sdk.Int       `json:"sellSize"`
	UnUseSize          sdk.Int       `json:"unUseSize"`
	Status             int           `json:"status"`
	CreateTime         time.Time     `json:"createTime"`
	CancelTimeDuration time.Duration `json:"cancelTimeDuration"`
	MarketAddress      string        `json:"marketAddress"`
	MinBuySize         sdk.Int       `json:"minBuySize"`  // config
	MinDuration        time.Duration `json:"minDuration"` // config
	MaxDuration        time.Duration `json:"maxDuration"` // config
	Reserve1           string        `json:"reserve1"`
}

func NewSellOrder(size int64) []byte {
	codec := cdc.New()
	order := SellOrder{
		OrderId:    ed25519.GenPrivKey().PubKey().Address().String(),
		Address:    "lambda1t9hr73cf77am25mmuy267sp90suyrxl7h08tkp",
		Price:      sdk.NewInt(1),
		Rate:       sdk.NewDec(1),
		Amount:     sdk.Coins{},
		SellSize:   sdk.NewInt(size),
		UnUseSize:  sdk.NewInt(size),
		Status:     0,
		CreateTime: time.Now(),
	}
	return codec.MustMarshalBinaryBare(order)
}
