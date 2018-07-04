package web_service

import (
	"net/http"
	"sort"

	"github.com/echoturing/auto_trader/trader/balances"
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/echoturing/auto_trader/trader/config"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type WebService struct {
}

func (w *WebService) GetConfig(c echo.Context) error {
	cfg := config.GetConfig()
	cfgList := make([]*config.BaseConfig, 0)
	keyList := make([]string, len(cfg))
	index := 0
	for key := range cfg {
		keyList[index] = key
		index += 1
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		cfgList = append(cfgList, cfg[k])
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"data": cfgList,
	})
}

func (w *WebService) UpdateConfig(c echo.Context) error {
	cfg := &config.BaseConfig{}
	if err := c.Bind(cfg); err != nil {
		log.Error(err.Error())
		return err
	}
	config.UpdateConfig(cfg.Coin, cfg.MarketCoin, cfg.Market, cfg.ToMarket, cfg)
	return c.JSON(http.StatusOK, map[string]interface{}{"code": 200, "data": true})
}

func (w *WebService) Pause(c echo.Context) error {
	statistics.Pause()
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func (w *WebService) Resume(c echo.Context) error {
	statistics.Resume()
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func (w *WebService) Status(c echo.Context) error {
	status := statistics.GetRunning()
	return c.JSON(http.StatusOK, map[string]bool{"running": status})
}

func (w *WebService) Balance(c echo.Context) error {
	res := map[string]*balances.Balance{
		balances.BfxBalance.Market:   balances.BfxBalance,
		balances.BinBalance.Market:   balances.BinBalance,
		balances.HuobiBalance.Market: balances.HuobiBalance,
	}
	return c.JSON(http.StatusOK, res)
}
