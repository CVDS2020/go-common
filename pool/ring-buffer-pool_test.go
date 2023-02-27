package pool

import (
	"fmt"
	"gitee.com/sy_183/common/unit"
	"testing"
)

func TestRingBufferPool(t *testing.T) {
	p := NewRingBufferPool(2, unit.MeBiByte, 2048)
	for i := 0; i < 10; i++ {
		var blocks []*Data
		for j := 0; ; j++ {
			buf := p.Get()
			if buf == nil {
				fmt.Println("alloc full")
				break
			}
			size := len(buf)
			if size > 1400 {
				size = 1400
			}
			data := p.Alloc(uint(size))
			fmt.Printf("alloc data(seq:%d,size:%d)\n", j, size)
			blocks = append(blocks, data)
		}
		for _, block := range blocks {
			block.Release()
		}
		blocks = blocks[:0]
	}
}
