package order_maker

import (
	"testing"

	"github.com/echoturing/auto_trader/trader/config"
	"github.com/echoturing/auto_trader/trader/depth"
	"github.com/echoturing/auto_trader/trader/statistics"
)

func TestMakeOrder(t *testing.T) {
	coin := statistics.CoinEOS
	marketCoin := statistics.CoinBTC
	bfxBook := &depth.Book{
		Market: statistics.MarketBitfinex,
		Bids: depth.BidsDepth{
			[]float64{1.5, 0.1}, //有人高价买
			[]float64{0.9, 0.2},
		},
		Asks: depth.AsksDepth{
			[]float64{1.7, 1}, //更高价卖
			[]float64{1.9, 1},
		},
		Coin:       coin,
		MarketCoin: marketCoin,
	}
	binBook := &depth.Book{
		Market: statistics.MarketBinance,
		Bids: depth.BidsDepth{
			[]float64{1, 0.1}, //有人低价买
			[]float64{0.9, 0.2},
		},
		Asks: depth.AsksDepth{
			[]float64{1.2, 1}, //有人高价卖
			[]float64{1.3, 1},
		},
		Coin:       coin,
		MarketCoin: marketCoin,
	}
	minQty := 0.01
	bfxConfig := &config.BaseConfig{
		Market:     statistics.MarketBitfinex,
		Coin:       coin,
		MarketCoin: marketCoin,
		BuyProfit:  0.001,
		ChargeRate: 0.002,
		Precision:  5,
		MinQty:     minQty,
		MaxQty:     100,
	}
	binConfig := &config.BaseConfig{
		Market:     statistics.MarketBinance,
		Coin:       coin,
		MarketCoin: marketCoin,
		BuyProfit:  0.001,
		ChargeRate: 0.0005,
		Precision:  5,
		MinQty:     minQty,
		MaxQty:     100,
	}
	bfxOrder, binOrder, err := MakeBfxAndBinance(bfxBook, binBook, bfxConfig, binConfig)
	//这个市场应该是bfx卖,bin买
	//数量是0.1个
	//卖价格1.5,买价格1.2 所以毛利是  (1.5-1.2)*0.1=0.03 是在bfx卖出,
	//手续费是    (1.5*0.1)*0.002  +(1.2*0.1)*0.0005 =0.00036
	//利润是   毛利-手续费 = 0.03-0.00036=0.02964
	//利润率是  (毛利-手续费)/买入花费 = (0.03-0.00036)/(1.2*0.1)=0.247
	//最小下单数
	t.Log(err)
	t.Logf("%#v\n%#v\n", bfxOrder, binOrder)
}
