package clients

import (
	"github.com/adshao/go-binance"
	"github.com/echoturing/auto_trader/trader/config"
)

var BinClient *binance.Client

func InitBin() {
	BinClient = binance.NewClient(config.BinConf.Key, config.BinConf.Secret)
}
