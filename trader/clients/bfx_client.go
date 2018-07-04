package clients

import (
	"github.com/echoturing/bitfinex-api-go/v1"
	"github.com/echoturing/auto_trader/trader/config"
)

func NewBfxClient() *bitfinex.Client {
	return bitfinex.NewClient().Auth(config.BfxConf.Key, config.BfxConf.Secret)
}
