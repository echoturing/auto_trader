package order_maker

import (
	"fmt"
	"math"

	"github.com/echoturing/auto_trader/trader/utils"
	"github.com/echoturing/auto_trader/trader/config"
	"github.com/echoturing/auto_trader/trader/depth"
	"github.com/echoturing/auto_trader/trader/order"
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

//根据book产生订单,只有遇到不能处理的错误,err才不是nil.双order是nil只是说明没行情
func MakeBfxAndBinance(
	m1Book, m2Book *depth.Book, //行情
	m1Config, m2Config *config.BaseConfig,
) (*order.Order, *order.Order, error) {
	//防傻逼
	var m1Order, m2Order *order.Order
	var err error

	if m1Book.Market != m1Config.Market {
		err = errors.New(fmt.Sprintf("参数错误~1"))
		return nil, nil, err
	}
	if m2Book.Market != m2Config.Market {
		err = errors.New(fmt.Sprintf("参数错误~2 %#v	%#v", m2Book, m2Config))
		return nil, nil, err
	}
	m1Order = &order.Order{
		Market:     m1Config.Market,
		Coin:       m1Book.Coin,
		MarketCoin: m1Book.MarketCoin,
		Skid:       m1Config.Skid,
	}
	m2Order = &order.Order{
		Market:     m2Config.Market,
		Coin:       m2Book.Coin,
		MarketCoin: m2Book.MarketCoin,
		Skid:       m2Config.Skid,
	}
	m1Sell1Price := m1Book.Asks[0][0]
	m1Sell1Qty := m1Book.Asks[0][1]
	m2Buy1Price := m2Book.Bids[0][0]
	m2Buy1Qty := m2Book.Bids[0][1]

	m2Order, m1Order = OrderBuilder(m2Buy1Price, m2Buy1Qty, m2Config, m2Order,
		m1Sell1Price, m1Sell1Qty, m1Config, m1Order,
	)
	if m2Order != nil && m1Order != nil {
		return m1Order, m2Order, nil
	}
	m1Order = &order.Order{
		Market:     m1Config.Market,
		Coin:       m1Book.Coin,
		MarketCoin: m1Book.MarketCoin,
		Skid:       m1Config.Skid,
	}
	m2Order = &order.Order{
		Market:     m2Config.Market,
		Coin:       m2Book.Coin,
		MarketCoin: m2Book.MarketCoin,
		Skid:       m2Config.Skid,
	}

	m1Buy1Price := m1Book.Bids[0][0]
	m1Buy1Qty := m1Book.Bids[0][1]
	m2Sell1Price := m2Book.Asks[0][0]
	m2Sell1Qty := m1Book.Asks[0][1]
	m1Order, m2Order = OrderBuilder(m1Buy1Price, m1Buy1Qty, m1Config, m1Order,
		m2Sell1Price, m2Sell1Qty, m2Config, m2Order,
	)
	if m2Order != nil && m1Order != nil {
		return m1Order, m2Order, nil
	}
	return nil, nil, nil
}

func OrderBuilder(
	m1Buy1Price, m1Qty float64, m1Config *config.BaseConfig, order1 *order.Order,
	m2Sell1Price, m2Qty float64, m2Config *config.BaseConfig, order2 *order.Order,
) (*order.Order, *order.Order) {
	precision := m1Config.Precision
	minQty := m1Config.MinQty
	maxQty := m1Config.MaxQty

	if m2Sell1Price < m1Buy1Price {
		marketLimit := math.Min(m1Qty, m2Qty)
		afterPrecision := utils.Round(marketLimit, int(precision))
		actualAmount := math.Min(afterPrecision, maxQty)
		if actualAmount < minQty || marketLimit < minQty {
			log.Infof("%s_%s:可下单数小于系统或市场限制:actualAmount:%f    marketLimit:%f ", order1.Coin, order1.MarketCoin, actualAmount, marketLimit)
			return nil, nil
		}
		m2BuyCost := m2Sell1Price * actualAmount //market 2 卖1价格很低,可以在m2买入
		m1SellEarn := m1Buy1Price * actualAmount //market 1 买1价格很高,可以在m1卖出

		serviceCharge := m2BuyCost*m2Config.ChargeRate + m1SellEarn*m1Config.ChargeRate
		expectTotalEarn := m1SellEarn - m2BuyCost     //期望赚(未扣除手续费)
		expectEarn := expectTotalEarn - serviceCharge //期望赚(扣除)
		expectProfitRate := expectEarn / m2BuyCost    //期望利润率(扣除)

		order2.Action = statistics.ActionBuy
		order2.Price = m2Sell1Price * (1 + m2Config.Skid) //买入的时候加价
		order2.ExpectPrice = m2Sell1Price                 //期望卖出价格
		order2.ExpectProfitRate = expectProfitRate        //期望利润率(扣除手续费)
		order2.ExpectProfit = expectEarn                  //期望利润(扣除手续费)

		order1.Action = statistics.ActionSell
		order1.Price = m1Buy1Price * (1 - m1Config.Skid) //卖出的时候降价
		order1.ExpectPrice = m1Buy1Price
		order1.ExpectProfitRate = expectProfitRate
		order1.ExpectProfit = expectEarn

		order1.Qty, order2.Qty = actualAmount, actualAmount
		var err error
		order1.Price, err = PriceFilter(order1.Price, order1.Coin, order1.MarketCoin, order1.Action)
		if err != nil {
			log.Errorf("PriceFilter err:%s", err.Error())
			return nil, nil
		}

		order2.Price, err = PriceFilter(order2.Price, order2.Coin, order2.MarketCoin, order2.Action)
		if err != nil {
			log.Errorf("PriceFilter err:", err.Error())
			return nil, nil
		}
		if order1.Action == statistics.ActionBuy {
			if expectProfitRate >= m1Config.BuyProfit {
				return order1, order2
			} else if expectProfitRate > 0 {
				log.Infof("发现行情,但是利润率不能触发下单:%#v	%#v", order1, order2)
				return nil, nil
			} else {
				return nil, nil
			}
		} else {
			if expectProfitRate >= m2Config.BuyProfit {
				return order1, order2
			} else if expectProfitRate > 0 {
				log.Infof("发现行情,但是利润率不能触发下单:%#v	%#v", order1, order2)
				return nil, nil
			} else {
				return nil, nil
			}
		}

	} else {
		return nil, nil
	}
}
