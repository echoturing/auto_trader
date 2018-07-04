package strategies

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/echoturing/auto_trader/trader/balances"
	"github.com/echoturing/auto_trader/trader/config"
	"github.com/echoturing/auto_trader/trader/db"
	"github.com/echoturing/auto_trader/trader/depth"
	"github.com/echoturing/auto_trader/trader/ding_service"
	"github.com/echoturing/auto_trader/trader/models"
	"github.com/echoturing/auto_trader/trader/order"
	"github.com/echoturing/auto_trader/trader/order_maker"
	"github.com/echoturing/auto_trader/trader/redis"
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/echoturing/auto_trader/trader/utils"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

func RefreshBalance(sendOrder <-chan struct{}) {
	timer := time.After(15 * time.Second)
	for {
		select {
		case <-sendOrder:
			timer = time.After(15 * time.Second)
		case <-timer:
			balances.Init()
			timer = time.After(15 * time.Second)
		}
	}
}

func Run(coin, marketCoin string, market1, market2 string, sendOrder chan struct{}) {
	errCount := time.Second
	for {
		if !statistics.GetRunning() {
			log.Warnf("暂停中...")
			time.Sleep(5 * time.Second)
			continue
		}
		//正常跑策略
		strategyConfig := config.GetConfig()
		m1Config := strategyConfig[config.GetConfigKey(coin, marketCoin, market1, market2)]
		m2Config := strategyConfig[config.GetConfigKey(coin, marketCoin, market2, market1)]
		err, success := Move(coin, marketCoin, m1Config, m2Config, market1, market2)
		if err != nil {
			//下单失败至少休息1分钟
			log.Error(err.Error())
			ding_service.TargetAlert(err.Error(), coin, marketCoin)
			time.Sleep(60 * errCount)
			errCount += time.Second
		} else if success {
			//下单成功后,
			sendOrder <- struct{}{}
			errCount = time.Second
			time.Sleep(10 * time.Second) //下单成功后休息10秒
		}
	}

}

func GetOrders(market1, market2, coin, marketCoin string, market1Conf, market2Conf *config.BaseConfig) (*order.Order, *order.Order) {
	m1Book, m2Book := depth.GetDepth(market1, market2, coin, marketCoin)
	err := checkTs(m2Book.Ts, m1Book.Ts)
	if err != nil {
		log.Infof("%s_%s:%s    %s:%d    %s%d", coin, marketCoin, err.Error(), market1, m1Book.Ts, market2, m2Book.Ts)
		return nil, nil
	}
	m1Order, m2Order, err := order_maker.MakeBfxAndBinance(m1Book, m2Book, market1Conf, market2Conf)
	if err != nil {
		log.Warnf("订单:%s", err.Error())
	}
	return m1Order, m2Order
}

func Move(coin, marketCoin string, m1Config, m2Config *config.BaseConfig, market1, market2 string) (error, bool) {
	var err error
	var m1Order, m2Order *order.Order
	for {
		if !statistics.GetRunning() {
			return nil, false
		}
		m1Order, m2Order = GetOrders(
			market1, market2,
			coin, marketCoin,
			m1Config, m2Config,
		)
		if m1Order == nil && m2Order == nil {
			time.Sleep(time.Millisecond * 100) // 从redis爬数据还是要休息一下
			continue
		} else if m1Order != nil && m2Order != nil {
			break
		} else {
			return errors.Errorf("程序出错了:%#v	%#v", m1Order, m2Order), false
		}
	}
	m1Balance := balances.GetBalanceByMarket(market1)
	m2Balance := balances.GetBalanceByMarket(market2)

	err = CheckAndMinusBalance(m1Balance, m1Order)
	//err了的话,就不会minus
	if err != nil {
		e := fmt.Sprintf("%s余额不足:%s	%#v		%#v", market1, err.Error(), *m1Order, *m2Order)
		log.Warn(e)
		ding_service.TargetAlert(e, "", "")
		time.Sleep(time.Millisecond * 500)
		return nil, false
	}
	err = CheckAndMinusBalance(m2Balance, m2Order)
	if err != nil {
		AddBalance(m1Balance, m1Order) //因为之前把m1的余额给减了,这里要给别人加回来
		e := fmt.Sprintf("%s余额不足:%s	%#v		%#v", market2, err.Error(), *m1Order, *m2Order)
		log.Warn(e)
		ding_service.TargetAlert(e, "", "")
		time.Sleep(time.Millisecond * 500)
		return nil, false
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	log.Infof("开始发送订单:%#v    %#v", m1Order, m2Order)
	go m1Order.Send(wg)
	go m2Order.Send(wg)
	wg.Wait()
	//下单完毕
	if m1Order.Result != nil && m2Order.Result != nil {
		//入库
		uniqueID := strconv.FormatInt(time.Now().UnixNano(), 10)
		now := utils.GetCurrentTsMs()
		m1DBModel := models.NewDBOrder(market1, db.Conn)
		m1DBModel.Coin = coin
		m1DBModel.MarketCoin = marketCoin
		m1DBModel.OrderID = m1Order.Result.OrderID
		m1DBModel.UniqueID = uniqueID
		m1DBModel.Price = m1Order.Price
		m1DBModel.Qty = m1Order.Qty
		m1DBModel.Action = m1Order.Action
		m1DBModel.ExpectProfit = m1Order.ExpectProfit
		m1DBModel.ExpectProfitRate = m1Order.ExpectProfitRate
		m1DBModel.ExpectPrice = m1Order.ExpectPrice
		m1DBModel.Skid = m1Order.Skid
		m1DBModel.CreatedTs = now
		m1DBModel.LastModify = now
		err := m1DBModel.Save()
		if err != nil {
			log.Error(err.Error())
		}
		m2DBModel := models.NewDBOrder(market2, db.Conn)
		m2DBModel.Coin = coin
		m2DBModel.MarketCoin = marketCoin
		m2DBModel.OrderID = m2Order.Result.OrderID
		m2DBModel.UniqueID = uniqueID
		m2DBModel.Price = m2Order.Price
		m2DBModel.Qty = m2Order.Qty
		m2DBModel.Action = m2Order.Action
		m2DBModel.ExpectProfit = m2Order.ExpectProfit
		m2DBModel.ExpectProfitRate = m2Order.ExpectProfitRate
		m2DBModel.ExpectPrice = m2Order.ExpectPrice
		m2DBModel.Skid = m2Order.Skid
		m2DBModel.CreatedTs = now
		m2DBModel.LastModify = now
		if err != nil {
			log.Error(err.Error())
		}
		m2DBModel.Save()
		ding_service.TargetAlert(fmt.Sprintf("下单成功\n%#v\n%#v", m1DBModel, m2DBModel), coin, marketCoin)
		return nil, true
	}
	return errors.New(fmt.Sprintf("send orders failed:%s->%#v	%s->%#v", market1, m1Order, market2, m2Order)), false
}

func checkTs(ts1, ts2 int64) error {
	if ts1 == 0 || ts2 == 0 {
		return errors.New("depth in redis is nil")
	}
	diff := math.Abs(float64(ts1 - ts2))
	if diff > 1000 { //1s内的行情
		return errors.New(fmt.Sprintf("时间间隔过大:%f", diff))
	}
	now := time.Now().UnixNano() / 1000000
	if now-ts1 > 1000 || now-ts2 > 1000 {
		log.Warnf("行情很久没更新了:%d		%d		%d", now, ts1, ts2)
		return errors.New("行情过去时间太久")
	}
	return nil
}

func CheckAndMinusBalance(balance *balances.Balance, order *order.Order) error {
	//防傻逼
	if balance.Market != order.Market {
		return errors.New("傻逼了吧..")
	}
	var balanceDetail *balances.BalanceDetail
	if order.Action == statistics.ActionBuy {
		balanceDetail = balance.GetDetailByCoin(order.MarketCoin) //买入花费marketCoin
		amount := order.Qty * order.Price
		if balanceDetail.Free < amount {
			if redis.BalanceNeedAlert(order.MarketCoin, order.Market) {
				redis.AddBalanceNotEnough(order.MarketCoin, order.Market)
			}
			return errors.New(fmt.Sprintf("买入%s<%f>,花费%s<%f>,余额不足:<%#v>", order.Coin, order.Qty, order.MarketCoin, amount, *balanceDetail))
		}
		balance.Minus(order.MarketCoin, amount)
	} else {
		balanceDetail = balance.GetDetailByCoin(order.Coin) //卖出花费coin
		if balanceDetail.Free < order.Qty {
			if redis.BalanceNeedAlert(order.Coin, order.Market) {
				redis.AddBalanceNotEnough(order.Coin, order.Market)
			}
			redis.AddBalanceNotEnough(order.Coin, order.Market)
			return errors.New(fmt.Sprintf("卖出%s<%f>,余额不足:<%#v>", order.Coin, order.Qty, *balanceDetail))
		}
		balance.Minus(order.Coin, order.Qty)
	}
	return nil
}

func AddBalance(balance *balances.Balance, order *order.Order) {
	if order.Action == statistics.ActionBuy {
		amount := order.Qty * order.Price
		balance.Add(order.MarketCoin, amount)
	} else {
		amount := order.Qty
		balance.Add(order.Coin, amount)
	}
}
