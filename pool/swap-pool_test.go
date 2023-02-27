package pool

import (
	"gitee.com/sy_183/common/utils"
	"math/rand"
	"testing"
	"time"
)

func pop[X any](s *[]X) (x X, ok bool) {
	if len(*s) == 0 {
		return
	}
	x = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return x, true
}

func TestSwapPool(t *testing.T) {
	var ips = make(chan *int, 1024)
	sp := NewSwapPool(100, func(p Pool[*int]) *int {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		ip := new(int)
		//fmt.Printf("make new int pointer %p\n", ip)
		return ip
	})

	go func() {
		for {
			ip := sp.Get()
			//fmt.Printf("get int pointer %p\n", ip)
			utils.ChanTryPush(ips, ip)
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
	}()

	go func() {
		for {
			select {
			case ip := <-ips:
				sp.Put(ip)
				//default:
			}
			//time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
	}()

	time.Sleep(time.Hour)
}
