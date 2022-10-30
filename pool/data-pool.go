package pool

import (
	"sync"
)

type DataPool struct {
	pool sync.Pool
	size uint
}

func NewDataPool(size uint) *DataPool {
	p := &DataPool{
		pool: sync.Pool{},
		size: size,
	}
	p.pool.New = p.new
	return p
}

func (p *DataPool) new() any {
	return &Data{
		raw:  make([]byte, p.size),
		pool: p,
	}
}

func (p *DataPool) put(data *Data) {
	p.pool.Put(data)
}

func (p *DataPool) Size() uint {
	return p.size
}

func (p *DataPool) Alloc(len uint) (d *Data) {
	return p.AllocCap(len, len)
}

func (p *DataPool) AllocCap(len, cap uint) (d *Data) {
	if cap > p.size {
		return NewData(make([]byte, len, cap))
	}
	d = p.pool.Get().(*Data).Use()
	d.Data = d.raw[:len:cap]
	return d
}
