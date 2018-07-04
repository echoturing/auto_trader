package depth

import (
	"encoding/json"

	"github.com/echoturing/auto_trader/trader/redis"
	"github.com/labstack/gommon/log"
)

type BidsDepth [][]float64

func (b BidsDepth) Len() int {
	return len(b)
}

func (b BidsDepth) Less(i, j int) bool {
	return b[i][0] > b[j][0]
}

func (b BidsDepth) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type AsksDepth [][]float64

func (a AsksDepth) Len() int {
	return len(a)
}

func (a AsksDepth) Less(i, j int) bool {
	return a[i][0] < a[j][0]
}

func (a AsksDepth) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Book struct {
	Ts         int64     `json:"ts"`
	Symbol     string    `json:"symbol"`
	Coin       string    `json:"coin"`
	MarketCoin string    `json:"market_coin"`
	Bids       BidsDepth `json:"bids"`
	Asks       AsksDepth `json:"asks"`
	Market     string    `json:"market"`
}

//从redis获取不同市场的depth
// market	huobi  bian  bfx  ...
// coin  UPPER CASE     EOS .
// marketCoin UPPER CASE  	ETH .
func GetDepth(market1, market2, coin, marketCoin string) (*Book, *Book) {
	m1Key := market1 + "_" + coin + "_" + marketCoin
	m2Key := market2 + "_" + coin + "_" + marketCoin
	res, err := redis.RedisCli.MGet(m1Key, m2Key).Result()
	if err != nil {
		log.Fatal(err)
	}
	m1DepthString := []byte(res[0].(string))
	m2DepthString := []byte(res[1].(string))
	m1Dep := &Book{Coin: coin, MarketCoin: marketCoin}
	m2Dep := &Book{Coin: coin, MarketCoin: marketCoin}
	err = json.Unmarshal(m1DepthString, m1Dep)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(m2DepthString, m2Dep)
	if err != nil {
		log.Fatal(err)
	}
	m1Dep.Coin = coin
	m2Dep.Coin = coin
	m1Dep.MarketCoin = marketCoin
	m2Dep.MarketCoin = marketCoin
	m1Dep.Market = market1
	m2Dep.Market = market2
	return m1Dep, m2Dep
}
