package order_maker

import (
	"testing"
	"auto_trader/trader/statistics"
)

func TestPriceFilter(t *testing.T) {
	var err error
	coin := statistics.CoinEOS
	marketCoin := statistics.CoinETH
	actionBuy := statistics.ActionBuy
	actionSell := statistics.ActionSell

	validPriceList := []float64{
		0.001,
		0.002,
		0.000001,
		0.000201,
	}
	for _, p := range validPriceList {
		_, err = PriceFilter(p, coin, marketCoin, actionBuy)
		if err != nil {
			t.Error(err.Error())
		}
	}

	for _, p := range validPriceList {
		_, err = PriceFilter(p, coin, marketCoin, actionSell)
		if err != nil {
			t.Error(err.Error())
		}
	}
	invalidPriceList := []float64{
		0.00000099, 100001,
	}
	for _, p := range invalidPriceList {
		_, err = PriceFilter(p, coin, marketCoin, actionBuy)
		if err == nil {
			t.Error("minPrice,maxPrice not work")
		}
		t.Log(err)
	}
	for _, p := range invalidPriceList {
		_, err = PriceFilter(p, coin, marketCoin, actionSell)
		if err == nil {
			t.Error("minPrice,maxPrice not work")
		}
		t.Log(err)
	}

	needCorrectPrice := 0.0002113
	correctPrice, err := PriceFilter(needCorrectPrice, coin, marketCoin, actionBuy)
	if err != nil {
		t.Error(err.Error())
	}
	if correctPrice != 0.000212 { //不知道浮点数是会多还是少..所以只能在最小值有个浮动
		t.Error(correctPrice)
	}

	needCorrectPrice = 0.123321123456
	correctPrice, err = PriceFilter(needCorrectPrice, coin, marketCoin, actionBuy)
	if err != nil {
		t.Error(err.Error())
	}
	//log.Info(correctPrice)
	if correctPrice != 0.123322 {
		t.Error(correctPrice)
	}

	needCorrectPrice = 0.123321123456
	correctPrice, err = PriceFilter(needCorrectPrice, coin, marketCoin, actionSell)
	if err != nil {
		t.Error(err.Error())
	}
	//log.Info(correctPrice)
	if correctPrice != 0.123320 {
		t.Error(correctPrice)
	}
}
