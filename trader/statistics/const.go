package statistics

import (
	"sync"
)

const (
	MarketBinance  = "bian"
	MarketBitfinex = "bfx"
	MarketHuobi    = "huobi"
	ActionBuy      = "buy"
	ActionSell     = "sell"
	CoinEOS        = "EOS"
	CoinETH        = "ETH"
	CoinIOST       = "IOST"
	CoinBTC        = "BTC"
	CoinUSDT       = "USDT"
)

func GetAllMarket() []string {
	return []string{
		MarketBinance,
		MarketBitfinex,
		MarketHuobi,
	}
}

var Running = false
var RunMutex = &sync.Mutex{}

func GetRunning() bool {
	RunMutex.Lock()
	defer RunMutex.Unlock()
	return Running
}

func Pause() {
	RunMutex.Lock()
	defer RunMutex.Unlock()
	Running = false
}

func Resume() {
	RunMutex.Lock()
	defer RunMutex.Unlock()
	Running = true
}

func GetOrderTableByMarket(market string) string {
	if market == MarketBitfinex {
		return "bitfinex_order"
	} else if market == MarketBinance {
		return "binance_order"
	} else if market == MarketHuobi {
		return "huobi_order"
	}
	return ""
}
