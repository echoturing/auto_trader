package utils

import (
	"strings"
	"math"
	"time"
)

func BfxSymbol(coin, marketCoin string) string {
	return coin + marketCoin
}

func BinSymbol(coin, marketCoin string) string {
	return coin + marketCoin
}


func HuoBiSymbol(coin, marketCoin string) string {
	return strings.ToLower(coin + marketCoin)
}

// 保留精度
func Round(x float64, n int) float64 {
	x = x * math.Pow10(n)
	t := math.Trunc(x)
	//if math.Abs(x-t) >= 0.5 {
	//	return (t + math.Copysign(1, x) ) / math.Pow10(n)
	//}
	return t / math.Pow10(n)
}

func GetCurrentTsMs() int64 {
	return time.Now().UnixNano() / 1000000
}
