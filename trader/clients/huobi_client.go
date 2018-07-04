package clients

import (
	"github.com/leek-box/sheep/huobi"
	"github.com/echoturing/auto_trader/trader/config"
)

var HuobiClient *huobi.Huobi

func InitHuobi() {
	var err error
	HuobiClient, err = huobi.NewHuobi(config.HuobiConf.Key, config.HuobiConf.Secret)
	if err != nil {
		panic(err)
	}
}
