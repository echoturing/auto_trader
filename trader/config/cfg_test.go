package config

import (
	"fmt"
	"testing"

	"github.com/echoturing/auto_trader/trader/statistics"
)

func TestInitFromBase(t *testing.T) {
	InitFromBase()
	market := statistics.MarketHuobi
	toMarket := statistics.MarketBinance
	coin := statistics.CoinEOS
	marketCoin := statistics.CoinBTC
	//cfgKey := GetConfigKey(coin, marketCoin, market, toMarket)
	c := &BaseConfig{
		Market:     market,
		ToMarket:   toMarket,
		Coin:       coin,
		MarketCoin: marketCoin,
		MinQty:     10,
	}
	UpdateConfig(statistics.CoinEOS, statistics.CoinBTC, statistics.MarketHuobi, statistics.MarketBinance, c)
	cfg := GetConfig()
	for k := range cfg {
		fmt.Printf("%#v\n", cfg[k])
	}
}
