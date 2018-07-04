package models

import (
	"github.com/gocraft/dbr"
	"github.com/echoturing/auto_trader/trader/statistics"
)

type DBOrder struct {
	ID               int64   `json:"id" db:"id"`
	Market           string  `json:"market"`
	Coin             string  `json:"coin" db:"coin"`
	MarketCoin       string  `json:"market_coin" db:"market_coin"`
	OrderID          string  `json:"order_id" db:"order_id"`    //订单id
	UniqueID         string  `json:"unique_id" db:"unique_id" ` //
	Price            float64 `json:"price" db:"price"`          //实际下单价格
	Qty              float64 `json:"qty" db:"qty"`              //下单数量
	Action           string  `json:"action" db:"action"`
	ExpectProfit     float64 `json:"expect_profit" db:"expect_profit"`
	ActualProfit     float64 `json:"actual_profit" db:"actual_profit"`
	ExpectProfitRate float64 `json:"expect_profit_rate" db:"expect_profit_rate"` //期望利润率
	ActualPrice      float64 `json:"actual_price" db:"actual_price"`             //实际成交价格,后续订单生成再更新
	ExpectPrice      float64 `json:"expect_price" db:"expect_price"`             //不加价的价格
	Skid             float64 `json:"skid" db:"skid"`                             //滑点
	CreatedTs        int64   `json:"created_ts" db:"created_ts"`
	LastModify       int64   `json:"last_modify" db:"last_modify"`
	Conn             *dbr.Connection
}

func NewDBOrder(market string, conn *dbr.Connection) *DBOrder {
	return &DBOrder{Market: market, Conn: conn}
}
func (b *DBOrder) Save() error {
	table := statistics.GetOrderTableByMarket(b.Market)
	session := b.Conn.NewSession(nil)
	res, err := session.InsertInto(table).
		Columns("order_id", "unique_id", "coin", "market_coin", "price", "qty", "action", "expect_profit", "expect_profit_rate", "expect_price", "skid", "created_ts", "last_modify").
		Record(b).
		Exec()
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}
