package order_maker

import (
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/pkg/errors"
	"fmt"
	"github.com/echoturing/auto_trader/trader/utils"
)

type Filter struct {
	MinPrice      float64
	MaxPrice      float64
	PriceTickSize float64
	MinQty        float64
	MaxQty        float64
	QtyStepSize   float64
	MinNotional   float64
	Precision     int
	ValidFactor   float64
}

func GetFilter(coin, marketCoin string) *Filter {
	var filter *Filter
	switch {
	case coin == statistics.CoinEOS && marketCoin == statistics.CoinETH:
		filter = &Filter{
			MinPrice:      0.00000100,
			MaxPrice:      100000.00000000,
			PriceTickSize: 0.00000100,
			MinQty:        0.01000000,
			MaxQty:        90000000.00000000,
			QtyStepSize:   0.01000000,
			MinNotional:   0.01000000,
			Precision:     6,
			ValidFactor:   1000000,
		}
	case coin == statistics.CoinEOS && marketCoin == statistics.CoinBTC:
		filter = &Filter{
			MinPrice:      0.0000001,
			MaxPrice:      100000.00000000,
			PriceTickSize: 0.00000010,
			MinQty:        0.01000000,
			MaxQty:        90000000.00000000,
			QtyStepSize:   0.01000000,
			MinNotional:   0.00100000,
			Precision:     7,
			ValidFactor:   10000000,
		}
	case coin == statistics.CoinETH && marketCoin == statistics.CoinBTC:
		filter = &Filter{
			MinPrice:      0.00000100,
			MaxPrice:      100000.00000000,
			PriceTickSize: 0.00000100,
			MinQty:        0.00100000,
			MaxQty:        100000.00000000,
			QtyStepSize:   0.00100000,
			MinNotional:   0.00100000,
			Precision:     6,
			ValidFactor:   1000000,
		}
	case coin == statistics.CoinIOST && marketCoin == statistics.CoinBTC:
		filter = &Filter{
			MinPrice:      0.00000001,
			MaxPrice:      100000.00000000,
			PriceTickSize: 0.00000001,
			MinQty:        1,
			MaxQty:        90000000.00000000,
			QtyStepSize:   1,
			MinNotional:   0.00100000, //  0.00000617 *100
			Precision:     8,
			ValidFactor:   100000000,
		}
	case coin == statistics.CoinIOST && marketCoin == statistics.CoinETH:
		filter = &Filter{
			MinPrice:      0.00000001,
			MaxPrice:      100000.00000000,
			PriceTickSize: 0.00000001,
			MinQty:        1,
			MaxQty:        90000000.00000000,
			QtyStepSize:   1,
			MinNotional:   0.01000000,  //0.00007628 *100
			Precision:     8,
			ValidFactor:   100000000,
		}
	}
	return filter
}

func (f *Filter) Filter(price float64, action string) (float64, error) {
	if price < f.MinPrice {
		return -1, errors.New(fmt.Sprintf("price<minPrice:%f    %f", price, f.MinPrice))
	}
	if price > f.MaxPrice {
		return -1, errors.New(fmt.Sprintf("price>maxPrice:%f    %f", price, f.MaxPrice))
	}
	//	买单需要加价,算公式的时候
	priceInt64 := int64(price * f.ValidFactor)
	modPrice := (priceInt64 - int64(f.MinPrice*f.ValidFactor)) % int64(f.PriceTickSize*f.ValidFactor)
	if action == statistics.ActionBuy {
		//lot_size 应该不会出事,qty和precision限制了
		//quantity >= minQty
		//quantity <= maxQty
		//(quantity-minQty) % stepSize == 0

		//price >= minPrice
		//price <= maxPrice
		//(price-minPrice) % tickSize == 0
		afterCalc := utils.Round(float64(priceInt64-modPrice)/f.ValidFactor+f.PriceTickSize, f.Precision)
		return afterCalc, nil
	} else { //  卖单需要减价
		afterCalc := utils.Round(float64(priceInt64-modPrice)/f.ValidFactor-f.PriceTickSize, f.Precision)
		return afterCalc, nil
	}
}

func PriceFilter(price float64, coin, marketCoin string, action string) (float64, error) {
	return GetFilter(coin, marketCoin).Filter(price, action)
}
