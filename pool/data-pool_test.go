package pool

import (
	"testing"
)

func TestBufPool(t *testing.T) {
	p := NewDataPool(1024)
	for i := 0; i < 200000000; i++ {
		d := p.Alloc(1024)
		d.Release()
	}
}

func TestZapPool(t *testing.T) {
	p := NewBufferPool(1024)
	for i := 0; i < 200000000; i++ {
		d := p.Get()
		d.Free()
	}
}
