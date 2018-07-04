package balances

import (
	"testing"
	"time"
	"sync"
	"github.com/labstack/gommon/log"
)

func TestB(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			b, err := GetBfxBalanceFromApi()
			if err != nil {
				log.Error(err.Error())
			}
			log.Printf("%s\n", b.ToJson())
			time.Sleep(time.Second * 10)
		}
	}()
	wg.Wait()
}
