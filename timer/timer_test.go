package timer

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.AfterFunc(time.Second, func() {
		println("timeout")
	})

	time.Sleep(time.Second * 2)

	if stopped := timer.Stop(); stopped {
		println(stopped)
	}
}
