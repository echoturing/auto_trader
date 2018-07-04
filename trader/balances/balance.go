package balances

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/echoturing/auto_trader/trader/clients"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type BalanceDetail struct {
	Free   float64 `json:"free"`
	Frozen float64 `json:"frozen"`
	Mutex  *sync.Mutex
}

type Balance struct {
	Market string         `json:"market"`
	EOS    *BalanceDetail `json:"eos"`
	ETH    *BalanceDetail `json:"eth"`
	BTC    *BalanceDetail `json:"btc"`
	USDT   *BalanceDetail `json:"usdt"`
	BNB    *BalanceDetail `json:"bnb"`
	IOST   *BalanceDetail `json:"iost"`
	Error  string         `json:"error"`
}

func (b *Balance) GetDetailByCoin(coin string) *BalanceDetail {
	switch coin {
	case statistics.CoinBTC:
		return b.BTC
	case statistics.CoinEOS:
		return b.EOS
	case statistics.CoinETH:
		return b.ETH
	case statistics.CoinIOST:
		return b.IOST
	case statistics.CoinUSDT:
		return b.USDT
	}

	return nil
}
func (b *Balance) ToJson() string {
	j, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func (b *Balance) Add(coin string, amount float64) {
	if coin == statistics.CoinEOS {
		b.EOS.Mutex.Lock()
		defer b.EOS.Mutex.Unlock()
		b.EOS.Free += amount
	} else if coin == statistics.CoinETH {
		b.ETH.Mutex.Lock()
		defer b.ETH.Mutex.Unlock()
		b.ETH.Free += amount
	} else if coin == statistics.CoinBTC {
		b.BTC.Mutex.Lock()
		defer b.BTC.Mutex.Unlock()
		b.BTC.Free += amount
	} else if coin == statistics.CoinIOST {
		b.IOST.Mutex.Lock()
		defer b.IOST.Mutex.Unlock()
		b.IOST.Free += amount
	} else if coin == statistics.CoinUSDT {
		b.USDT.Mutex.Lock()
		defer b.USDT.Mutex.Unlock()
		b.USDT.Free += amount
	}
}

//大部分情况应该都是减余额吧,余额变负数就得跪
func (b *Balance) Minus(coin string, amount float64) error {
	if coin == statistics.CoinEOS {
		b.EOS.Mutex.Lock()
		defer b.EOS.Mutex.Unlock()
		if b.EOS.Free < amount {
			return errors.New(fmt.Sprintf("%s余额不足->需要%f,当前%f", coin, amount, b.EOS.Free))
		}
		b.EOS.Free -= amount
	} else if coin == statistics.CoinETH {
		b.ETH.Mutex.Lock()
		defer b.ETH.Mutex.Unlock()
		if b.ETH.Free < amount {
			return errors.New(fmt.Sprintf("%s余额不足->需要%f,当前%f", coin, amount, b.EOS.Free))
		}
		b.ETH.Free -= amount
	} else if coin == statistics.CoinBTC {
		b.BTC.Mutex.Lock()
		defer b.BTC.Mutex.Unlock()
		if b.BTC.Free < amount {
			return errors.New(fmt.Sprintf("%s余额不足->需要%f,当前%f", coin, amount, b.EOS.Free))
		}
		b.BTC.Free -= amount
	} else if coin == statistics.CoinIOST {
		b.IOST.Mutex.Lock()
		defer b.IOST.Mutex.Unlock()
		if b.IOST.Free < amount {
			return errors.New(fmt.Sprintf("%s余额不足->需要%f,当前%f", coin, amount, b.EOS.Free))
		}
		b.IOST.Free -= amount
	} else if coin == statistics.CoinUSDT {
		b.USDT.Mutex.Lock()
		defer b.USDT.Mutex.Unlock()
		if b.USDT.Free < amount {
			return errors.New(fmt.Sprintf("%s余额不足->需要%f,当前%f", coin, amount, b.EOS.Free))
		}
		b.USDT.Free -= amount
	}
	return nil
}

func NewBalance() *Balance {
	return &Balance{
		EOS: &BalanceDetail{
			Mutex: &sync.Mutex{},
		},
		ETH: &BalanceDetail{
			Mutex: &sync.Mutex{},
		},
		BTC: &BalanceDetail{
			Mutex: &sync.Mutex{},
		},
		BNB: &BalanceDetail{
			Mutex: &sync.Mutex{},
		},
		IOST: &BalanceDetail{
			Mutex: &sync.Mutex{},
		},
		USDT: &BalanceDetail{
			Mutex: &sync.Mutex{},
		},
	}
}

func GetBinBalanceFromApi() (*Balance, error) {
	binBalance := NewBalance()
	account, err := clients.BinClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		binBalance.Error = err.Error()
		return binBalance, errors.New(fmt.Sprintf("初始化binance余额失败:%s", err.Error()))
	}
	binBalance.Market = statistics.MarketBinance
	for _, b := range account.Balances {
		switch b.Asset {
		case "EOS":
			free, _ := strconv.ParseFloat(b.Free, 64)
			locked, _ := strconv.ParseFloat(b.Locked, 64)
			binBalance.EOS.Free = free
			binBalance.EOS.Frozen = locked
		case "ETH":
			free, _ := strconv.ParseFloat(b.Free, 64)
			locked, _ := strconv.ParseFloat(b.Locked, 64)
			binBalance.ETH.Free = free
			binBalance.ETH.Frozen = locked

		case "BTC":
			free, _ := strconv.ParseFloat(b.Free, 64)
			locked, _ := strconv.ParseFloat(b.Locked, 64)
			binBalance.BTC.Free = free
			binBalance.BTC.Frozen = locked
		case "BNB":
			free, _ := strconv.ParseFloat(b.Free, 64)
			locked, _ := strconv.ParseFloat(b.Locked, 64)
			binBalance.BNB.Free = free
			binBalance.BNB.Frozen = locked
		case "IOST":
			free, _ := strconv.ParseFloat(b.Free, 64)
			locked, _ := strconv.ParseFloat(b.Locked, 64)
			binBalance.IOST.Free = free
			binBalance.IOST.Frozen = locked

		case "USDT":
			free, _ := strconv.ParseFloat(b.Free, 64)
			locked, _ := strconv.ParseFloat(b.Locked, 64)
			binBalance.USDT.Free = free
			binBalance.USDT.Frozen = locked
		default:
			continue
		}

	}
	return binBalance, nil

}

func GetBfxBalanceFromApi() (*Balance, error) {
	bfxBalance := NewBalance()
	bfxBalance.Market = statistics.MarketBitfinex
	balances, err := clients.NewBfxClient().Balances.All()
	if err != nil {
		bfxBalance.Error = err.Error()
		return bfxBalance, errors.New(fmt.Sprintf("初始化bitfinex余额失败:%s", err.Error()))
	}
	for _, b := range balances {
		switch b.Currency {
		case "eth":
			available, _ := strconv.ParseFloat(b.Available, 64)
			amount, _ := strconv.ParseFloat(b.Amount, 64)
			bfxBalance.ETH.Free = available
			bfxBalance.ETH.Frozen = amount - available
		case "eos":
			available, _ := strconv.ParseFloat(b.Available, 64)
			amount, _ := strconv.ParseFloat(b.Amount, 64)
			bfxBalance.EOS.Free = available
			bfxBalance.EOS.Frozen = amount - available

		case "btc":
			available, _ := strconv.ParseFloat(b.Available, 64)
			amount, _ := strconv.ParseFloat(b.Amount, 64)
			bfxBalance.BTC.Free = available
			bfxBalance.BTC.Frozen = amount - available
		case "iost":
			available, _ := strconv.ParseFloat(b.Available, 64)
			amount, _ := strconv.ParseFloat(b.Amount, 64)
			bfxBalance.IOST.Free = available
			bfxBalance.IOST.Frozen = amount - available
		}
	}

	return bfxBalance, nil
}

func GetHuobiBalanceFromApi() (*Balance, error) {
	huobiBalance := NewBalance()
	huobiBalance.Market = statistics.MarketHuobi
	balances, err := clients.HuobiClient.GetAccountBalance()
	if err != nil {
		huobiBalance.Error = err.Error()
		return huobiBalance, errors.New(fmt.Sprintf("初始化火币余额失败:%s", err.Error()))
	}
	for _, b := range balances {
		if b.Currency == "eos" {
			count, _ := strconv.ParseFloat(b.Balance, 64)
			if b.Type == "trade" {
				huobiBalance.EOS.Free = count
			} else if b.Type == "frozen" {
				huobiBalance.EOS.Frozen = count
			}
		} else if b.Currency == "eth" {
			count, _ := strconv.ParseFloat(b.Balance, 64)
			if b.Type == "trade" {
				huobiBalance.ETH.Free = count
			} else if b.Type == "frozen" {
				huobiBalance.ETH.Frozen = count
			}

		} else if b.Currency == "btc" {
			count, _ := strconv.ParseFloat(b.Balance, 64)
			if b.Type == "trade" {
				huobiBalance.BTC.Free = count
			} else if b.Type == "frozen" {
				huobiBalance.BTC.Frozen = count
			}
		} else if b.Currency == "usdt" {
			count, _ := strconv.ParseFloat(b.Balance, 64)
			if b.Type == "trade" {
				huobiBalance.USDT.Free = count
			} else if b.Type == "frozen" {
				huobiBalance.USDT.Frozen = count
			}
		} else if b.Currency == "iost" {
			count, _ := strconv.ParseFloat(b.Balance, 64)
			if b.Type == "trade" {
				huobiBalance.IOST.Free = count
			} else if b.Type == "frozen" {
				huobiBalance.IOST.Frozen = count
			}
		}

	}
	return huobiBalance, nil

}

var BinBalance *Balance
var BfxBalance *Balance
var HuobiBalance *Balance

func GetBalanceByMarket(market string) *Balance {
	if market == statistics.MarketBinance {
		return BinBalance
	} else if market == statistics.MarketBitfinex {
		return BfxBalance
	} else if market == statistics.MarketHuobi {
		return HuobiBalance
	}
	return nil
}

func Init() {
	var err error
	BfxBalance, err = GetBfxBalanceFromApi()
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("初始化bitfinex余额:%s", BfxBalance.ToJson())
	BinBalance, err = GetBinBalanceFromApi()
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("初始化binance余额:%s", BinBalance.ToJson())
	HuobiBalance, err = GetHuobiBalanceFromApi()
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("初始化火币余额:%s", HuobiBalance.ToJson())
}
