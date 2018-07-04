package order

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/echoturing/auto_trader/trader/utils"
	"github.com/echoturing/auto_trader/trader/clients"
	"github.com/adshao/go-binance"
	"github.com/echoturing/bitfinex-api-go/v1"
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/labstack/gommon/log"
	"github.com/leek-box/sheep/proto"
)

type Result struct {
	OrderID string `json:"order_id"`
	Detail  interface{}
}
type Order struct {
	Market           string  `json:"market"`             //市场
	MarketCoin       string  `json:"market_coin"`        //计价币种
	Coin             string  `json:"coin"`               //交易币种
	Action           string  `json:"action"`             //buy or sell
	Price            float64 `json:"price"`              //实际下单价格
	Qty              float64 `json:"qty"`                //实际下单数量
	ExpectProfit     float64 `json:"expect_profit"`      //期望利润
	ExpectProfitRate float64 `json:"expect_profit_rate"` //期望利润率
	ExpectPrice      float64 `json:"expect_price"`       //理论价格
	Skid             float64 `json:"skid"`               //加价百分比(滑点)
	Result           *Result `json:"result"`             //订单结果
}

//发送订单
func (o *Order) Send(wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
		if r := recover(); r != nil {
			var err error
			switch r := r.(type) {
			case error:
				err = r
			default:
				err = fmt.Errorf("%v", r)
			}
			log.Error(err.Error())
		}
	}()
	var err error
	ts1 := time.Now().UnixNano()
	log.Infof("start send order:%#v", o)
	if o.Market == statistics.MarketBitfinex {
		_, err = BfxSendOrder(o)
	} else if o.Market == statistics.MarketBinance {
		_, err = BinSendOrder(o)
	} else if o.Market == statistics.MarketHuobi {
		_, err = HuoBiSendOrder(o)
	}
	ts2 := time.Now().UnixNano()
	log.Infof("%s send order cost %s", o.Market, ts2-ts1)
	if err != nil {
		return err
	}
	return nil
}

func BfxSendOrder(order *Order) (*Order, error) {
	symbol := utils.BfxSymbol(order.Coin, order.MarketCoin)
	var qty = order.Qty
	if order.Action == statistics.ActionSell {
		// api卖是负数
		qty = -1 * order.Qty
	}
	o, err := clients.NewBfxClient().Orders.Create(symbol, qty, order.Price, bitfinex.OrderTypeExchangeLimit) //限价买
	log.Printf("BfxSendOrder:%#v    %#v", o, err)
	if err != nil {
		if strings.Contains(err.Error(), "once") {
			o, err = clients.NewBfxClient().Orders.Create(symbol, qty, order.Price, bitfinex.OrderTypeExchangeLimit) //限价买
			log.Printf("BfxSendOrder retry:%#v    %#v", o, err)
			if err != nil {
				return nil, err
			}
			order.Result = &Result{
				OrderID: strconv.FormatInt(o.ID, 10),
				Detail:  o,
			}
			return order, nil
		}
	}
	order.Result = &Result{
		OrderID: strconv.FormatInt(o.ID, 10),
		Detail:  o,
	}
	return order, nil
}

func BinSendOrder(order *Order) (*Order, error) {
	symbol := utils.BinSymbol(order.Coin, order.MarketCoin)
	var side binance.SideType
	if order.Action == statistics.ActionBuy {
		side = binance.SideTypeBuy
	} else {
		side = binance.SideTypeSell
	}
	qtyStr := strconv.FormatFloat(order.Qty, 'f', -1, 64)
	priceStr := strconv.FormatFloat(order.Price, 'f', -1, 64)
	o, err := clients.BinClient.NewCreateOrderService().Symbol(symbol).Side(side).
		Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceGTC).
		Quantity(qtyStr).Price(priceStr).Do(context.Background())
	log.Printf("BinSendOrder:%#v    %#v", o, err)
	if err != nil {
		return nil, err
	}
	order.Result = &Result{
		OrderID: strconv.FormatInt(o.OrderID, 10),
		Detail:  o,
	}
	return order, nil
}

func HuoBiSendOrder(order *Order) (*Order, error) {
	var typ string
	if order.Action == statistics.ActionBuy {
		typ = proto.OrderPlaceTypeBuyLimit
	} else if order.Action == statistics.ActionSell {
		typ = proto.OrderPlaceTypeSellLimit
	}
	param := &proto.OrderPlaceParams{
		Price:           order.Price,
		Amount:          order.Qty,
		Type:            typ,
		BaseCurrencyID:  order.Coin,
		QuoteCurrencyID: order.MarketCoin,
	}
	r, err := clients.HuobiClient.OrderPlace(param)
	log.Printf("HuobiSendOrder:%#v    %#v", r, err)
	if err != nil {
		return nil, err
	}
	order.Result = &Result{
		OrderID: r.OrderID,
		Detail:  r,
	}
	return order, nil
}
