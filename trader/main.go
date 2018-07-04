package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/echoturing/auto_trader/trader/balances"
	"github.com/echoturing/auto_trader/trader/ding_service"
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/echoturing/auto_trader/trader/strategies"
	"github.com/echoturing/auto_trader/trader/web_service"

	"path/filepath"

	"github.com/echoturing/auto_trader/trader/clients"
	"github.com/echoturing/auto_trader/trader/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func httpServer(wg *sync.WaitGroup) {
	e := echo.New()
	e.Use(
		middleware.CORS(),    // 支持跨域
		middleware.Logger(),  // 标准log
		middleware.Recover(), //恢复panic
	)
	e.Logger.SetLevel(log.INFO)
	e.Static("/new/static", "static")
	webService := web_service.WebService{}
	e.GET("/new/config", webService.GetConfig)
	e.POST("/new/config", webService.UpdateConfig)
	e.POST("/new/pause", webService.Pause)
	e.POST("/new/resume", webService.Resume)
	e.GET("/new/status", webService.Status)
	e.GET("/new/balance", webService.Balance)
	go func() {
		if err := e.Start(":8099"); err != nil {
			e.Logger.Info(fmt.Sprintf("shutting down the server:%s", err.Error()))
			wg.Done()
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func run(coin, marketCoin string, market1, market2 string, sendOrder chan struct{}) {
	strategies.Run(coin, marketCoin, market1, market2, sendOrder)
}

func main() {
	cfgFilePath, err := filepath.Abs("trader/etc/config.yaml")
	if err != nil {
		panic(err)
	}
	config.InitConfig(cfgFilePath)
	clients.InitBin()
	clients.InitHuobi()
	balances.Init()
	var wg sync.WaitGroup
	wg.Add(1)
	go httpServer(&wg)
	coin1 := flag.String("coin1", "", "coin1")
	marketCoin1 := flag.String("market_coin1", "", "marketCoin1")
	coin2 := flag.String("coin2", "", "coin2")
	marketCoin2 := flag.String("market_coin2", "", "marketCoin2")

	coin3 := flag.String("coin3", "", "coin3")
	marketCoin3 := flag.String("market_coin3", "", "marketCoin3")
	//m1 := flag.String("m1", "bian", "market1")
	//m2 := flag.String("m2", "bfx", "market2")

	flag.Parse()
	*coin1 = strings.ToUpper(*coin1)
	*marketCoin1 = strings.ToUpper(*marketCoin1)
	*coin2 = strings.ToUpper(*coin2)
	*marketCoin2 = strings.ToUpper(*marketCoin2)
	*coin3 = strings.ToUpper(*coin3)
	*marketCoin3 = strings.ToUpper(*marketCoin3)

	sendOrder := make(chan struct{})
	go strategies.RefreshBalance(sendOrder)

	go func() {
		run(statistics.CoinEOS, statistics.CoinBTC, statistics.MarketHuobi, statistics.MarketBinance, sendOrder)
	}()
	go func() {
		run(statistics.CoinEOS, statistics.CoinETH, statistics.MarketHuobi, statistics.MarketBinance, sendOrder)
	}()

	go func() {
		run(statistics.CoinEOS, statistics.CoinBTC, statistics.MarketHuobi, statistics.MarketBitfinex, sendOrder)
	}()

	go func() {
		run(statistics.CoinEOS, statistics.CoinETH, statistics.MarketHuobi, statistics.MarketBitfinex, sendOrder)
	}()

	go func() {
		run(statistics.CoinEOS, statistics.CoinBTC, statistics.MarketBinance, statistics.MarketBitfinex, sendOrder)
	}()

	go func() {
		run(statistics.CoinEOS, statistics.CoinETH, statistics.MarketBinance, statistics.MarketBitfinex, sendOrder)
	}()
	go func() {
		run(statistics.CoinIOST, statistics.CoinBTC, statistics.MarketHuobi, statistics.MarketBinance, sendOrder)
	}()
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		select {
		case <-ticker.C:
			if balances.BinBalance.BNB.Free < 2 {
				ding_service.TargetAlert("BNB余额不足,当前剩余小于2个", "", "")
			}
		}
	}()

	wg.Wait()
}
