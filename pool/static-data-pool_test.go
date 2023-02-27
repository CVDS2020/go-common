package pool

import (
	"testing"
)

func TestBufPool(t *testing.T) {
	p := NewStaticDataPool(1024, ProvideSlicePool[*Data])
	for i := 0; i < 200000000; i++ {
		d := p.Alloc(1024)
		d.Release()
	}
}
