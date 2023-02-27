package pool

import "sync/atomic"

type RingBufferPool struct {
	buffers   []*Buffer
	allocated []atomic.Bool
	index     int
}

func NewRingBufferPool(count uint, size uint, reversed uint) *RingBufferPool {
	if count == 0 {
		count = 1
	}
	p := &RingBufferPool{
		buffers:   make([]*Buffer, count),
		allocated: make([]atomic.Bool, count),
	}
	for i := range p.buffers {
		p.buffers[i] = NewBuffer(size, reversed).SetOnReleased(func(*Buffer) {
			p.allocated[i].Store(false)
		})
	}
	p.allocated[0].Store(true)
	p.buffers[0].AddRef()
	return p
}

func (p *RingBufferPool) Get() []byte {
	if buf := p.buffers[p.index].Get(); buf != nil {
		return buf
	}
	p.buffers[p.index].Release()
	index := p.index + 1
	if index == len(p.buffers) {
		index = 0
	}
	if !p.allocated[index].CompareAndSwap(false, true) {
		return nil
	}
	p.index = index
	return p.buffers[index].Use().Get()
}

func (p *RingBufferPool) Alloc(size uint) *Data {
	return p.buffers[p.index].Alloc(size)
}
