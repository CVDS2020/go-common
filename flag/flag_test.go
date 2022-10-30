package flag

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestFlagMask(t *testing.T) {
	start := time.Now()
	for i := 0; i < math.MaxInt; i++ {
		if TestFlag(i, 0x1000000000000000) {
			return
		}
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		//_ = i | 0x100011101
		if i%100000000 == 0 {
			d := time.Now().Sub(start)
			speed := float64(i) / d.Seconds()
			fmt.Printf("%10f\n", speed)
		}
	}
}
