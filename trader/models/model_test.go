package models

import (
	"strconv"
	"testing"
	"time"

	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/echoturing/auto_trader/trader/utils"

	"github.com/echoturing/auto_trader/trader/db"
	"github.com/labstack/gommon/log"
)

func TestSave(t *testing.T) {
	dbOrder := &DBOrder{
		Market:           statistics.MarketBinance,
		Coin:             statistics.CoinEOS,
		MarketCoin:       statistics.CoinBTC,
		OrderID:          "xxxxx",
		UniqueID:         strconv.FormatInt(time.Now().UnixNano(), 10),
		Price:            1,
		Qty:              2,
		Action:           statistics.ActionBuy,
		ExpectPrice:      0.9,
		ExpectProfitRate: 0.001,
		ExpectProfit:     0.1,
		Skid:             0.1,
		CreatedTs:        utils.GetCurrentTsMs(),
		LastModify:       utils.GetCurrentTsMs(),
		Conn:             db.Conn,
	}
	err := dbOrder.Save()
	log.Infof("%#v", err)
}
