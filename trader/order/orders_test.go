package order

import (
	"log"
	"testing"

	"github.com/echoturing/auto_trader/trader/statistics"
)

func TestBinSendOrder(t *testing.T) {
	order := &Order{
		Market:     statistics.MarketBinance,
		Coin:       statistics.CoinIOST,
		MarketCoin: statistics.CoinBTC,
		Action:     statistics.ActionBuy,
		Price:      0.00000017,
		Qty:        1,
	}
	o, err := BinSendOrder(order)
	log.Printf("%#v", o)
	log.Printf("%s", err.Error())
}
