package ding_service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/echoturing/auto_trader/trader/config"
	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/labstack/gommon/log"
)

var HttpClient = &http.Client{Timeout: time.Second * 10}

func Alert(content string, url string) {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"isAtAll": true,
	}
	msgString, _ := json.Marshal(msg)
	resp, err := HttpClient.Post(url, "application/json", strings.NewReader(string(msgString)))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Infof("ding:%s", body)
}

func TargetAlert(content, coin, marketCoin string) {
	var url string
	if coin == statistics.CoinEOS && marketCoin == statistics.CoinBTC {
		url = config.GlobalKeyConfigs.DingConfig.EOSBTC
	} else if coin == statistics.CoinEOS && marketCoin == statistics.CoinETH {
		url = config.GlobalKeyConfigs.DingConfig.EOSETH
	} else if coin == statistics.CoinETH && marketCoin == statistics.CoinBTC {
		url = config.GlobalKeyConfigs.DingConfig.ETHBTC
	} else if coin == statistics.CoinIOST && marketCoin == statistics.CoinBTC {
		url = config.GlobalKeyConfigs.DingConfig.IOSTBTC
	} else if coin == statistics.CoinIOST && marketCoin == statistics.CoinETH {
		url = config.GlobalKeyConfigs.DingConfig.IOSTETH
	} else {
		url = config.GlobalKeyConfigs.DingConfig.Alert
	}
	Alert(content, url)
}
