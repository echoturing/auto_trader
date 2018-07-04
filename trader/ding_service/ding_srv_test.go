package ding_service

import (
	"testing"

	"github.com/echoturing/auto_trader/trader/statistics"
)

func TestAlert(t *testing.T) {
	//TargetAlert("测试消息1", "", "")
	TargetAlert("测试消息EOS_BTC", statistics.CoinEOS, statistics.CoinBTC)
}
