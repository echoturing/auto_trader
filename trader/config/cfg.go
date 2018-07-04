package config

import (
	"encoding/json"
	"sync"

	"github.com/echoturing/auto_trader/trader/redis"

	"path/filepath"

	"os"

	"io/ioutil"

	"github.com/echoturing/auto_trader/trader/statistics"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Key    string
	Secret string
}
type DingConfig struct {
	EOSETH  string `yaml:"eos_eth"`
	EOSBTC  string `yaml:"eos_btc"`
	ETHBTC  string `yaml:"eth_btc"`
	IOSTBTC string `yaml:"iost_btc"`
	IOSTETH string `yaml:"iost_eth"`
	Alert   string `yaml:"alert"`
}

type KeyConfigs struct {
	Huobi      *Config     `yaml:"huobi"`
	Binance    *Config     `yaml:"binance"`
	Bitfinex   *Config     `yaml:"bfx"`
	DingConfig *DingConfig `yaml:"dingTalk"`
}

var BfxConf *Config
var BinConf *Config
var HuobiConf *Config
var GlobalKeyConfigs *KeyConfigs

func LoadConfigFromFile(filePath string) (*KeyConfigs, error) {
	file, err := filepath.Abs(filePath)
	if err != nil {
		log.Warn("load config file err:", err.Error())
		return nil, err
	}
	f, err := os.Open(file)
	if err != nil {
		log.Warn("open config file err:", err.Error())
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Warn("read config file err:", err.Error())
		return nil, err
	}
	var keyConfig KeyConfigs
	err = yaml.Unmarshal(b, &keyConfig)
	if err != nil {
		log.Warn("unmarshal config file err:", err.Error())
		return nil, err
	}
	GlobalKeyConfigs = &keyConfig
	return &keyConfig, nil
}
func InitConfig(filePath string) {
	keyConfig, err := LoadConfigFromFile(filePath)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	BfxConf = keyConfig.Bitfinex
	BinConf = keyConfig.Binance
	HuobiConf = keyConfig.Huobi
}

type BaseConfig struct {
	Market     string  `json:"market"`
	ToMarket   string  `json:"to_market"`
	Coin       string  `json:"coin"`
	MarketCoin string  `json:"market_coin"`
	BuyProfit  float64 `json:"buy_profit"`
	ChargeRate float64 `json:"charge_rate"`
	Precision  int64   `json:"precision"`
	Skid       float64 `json:"skid"`
	MinQty     float64 `json:"min_qty"`
	MaxQty     float64 `json:"max_qty"`
}

const baseConfigMapKey = "NEW_BASE_CONFIG"

var baseConfigMap = make(map[string]*BaseConfig)
var m = &sync.Mutex{}

func GetConfig() map[string]*BaseConfig {
	return baseConfigMap
}

func GetConfigKey(coin, marketCoin, market, toMarket string) string {
	return market + "_" + toMarket + "_" + coin + marketCoin
}

func UpdateConfig(coin, marketCoin, market, toMarket string, cfg *BaseConfig) {
	m.Lock()
	defer m.Unlock()
	key := GetConfigKey(coin, marketCoin, market, toMarket)
	markets := statistics.GetAllMarket()
	for _, m := range markets {
		if m != market { //同symbol的数据需要同步qty和precision
			for _, toM := range markets {
				if m != toM {
					tempKey := GetConfigKey(coin, marketCoin, m, toM)
					baseConfigMap[tempKey].MinQty = cfg.MinQty
					baseConfigMap[tempKey].MaxQty = cfg.MaxQty
					baseConfigMap[tempKey].Precision = cfg.Precision

					tempKey2 := GetConfigKey(coin, marketCoin, toM, m)
					baseConfigMap[tempKey2].MinQty = cfg.MinQty
					baseConfigMap[tempKey2].MaxQty = cfg.MaxQty
					baseConfigMap[tempKey2].Precision = cfg.Precision
				}

			}

		}
	}
	baseConfigMap[key] = cfg
	cfgString, err := json.Marshal(baseConfigMap)
	if err != nil {
		log.Fatal(err)
	}
	redis.RedisCli.Set(baseConfigMapKey, cfgString, 0)
}

func getCfgFromRedis() string {
	configMapString := redis.RedisCli.Get(baseConfigMapKey).Val()
	return configMapString
}

func InitFromBase() {
	coins := []string{statistics.CoinEOS, statistics.CoinETH, statistics.CoinIOST}
	marketCoins := []string{statistics.CoinETH, statistics.CoinBTC}
	markets := []string{statistics.MarketBitfinex, statistics.MarketBinance, statistics.MarketHuobi}
	for _, coin := range coins {
		for _, marketCoin := range marketCoins {
			if coin != marketCoin {
				for _, market := range markets {
					for _, toMarket := range markets {
						if market != toMarket {
							key := GetConfigKey(coin, marketCoin, market, toMarket)
							value := redis.RedisCli.Get(key).Val()
							cfg := &BaseConfig{}
							if value != "" {
								err := json.Unmarshal([]byte(value), cfg)
								if err != nil {
									log.Fatal(err)
								}
								baseConfigMap[key] = cfg
							} else {
								baseConfigMap[key] = cfg
								cfg.Market = market
								cfg.Coin = coin
								cfg.MarketCoin = marketCoin
								cfg.ToMarket = toMarket
							}
						}
					}
				}
			}
		}
	}
}

func init() {
	configMapString := getCfgFromRedis()
	if configMapString != "" {
		err := json.Unmarshal([]byte(configMapString), &baseConfigMap)
		if err != nil {
			log.Fatal(err)
		}

		for key := range baseConfigMap {
			log.Infof("%#v", *baseConfigMap[key])
		}
		return
	}
	InitFromBase()
}
