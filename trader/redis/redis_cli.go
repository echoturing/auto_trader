package redis

import (
	"github.com/go-redis/redis"
	"time"
)

var RedisCli *redis.Client

func init() {
	opt := &redis.Options{
		Addr:     "localhost:6379",
		Password: "qRo$t,88=*z3)+2b484MZ9[Y$9{FC3a&",
		DB:       1,
	}
	RedisCli = redis.NewClient(opt)
}

func GetBalanceAlertKey(coin, market string) string {
	return "balance_alert_" + coin + "_" + market
}

func BalanceNeedAlert(coin, market string) bool {
	key := GetBalanceAlertKey(coin, market)
	v := RedisCli.Get(key).Val()
	if v != "" {
		return true
	}
	return false
}

func AddBalanceNotEnough(coin, market string) error {
	key := GetBalanceAlertKey(coin, market)
	err := RedisCli.Set(key, "not need", time.Second*30).Err()
	if err != nil {
		return err
	}
	return nil
}
