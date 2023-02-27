package pool

type DefaultBufferPool struct {
	pool   Pool[*Buffer]
	buffer *Buffer
}

func NewDefaultBufferPool(size uint, reversed uint, poolProvider PoolProvider[*Buffer]) *DefaultBufferPool {
	return &DefaultBufferPool{pool: poolProvider(func(p Pool[*Buffer]) *Buffer {
		return NewBuffer(size, reversed).SetOnReleased(func(buffer *Buffer) {
			p.Put(buffer)
		})
	})}
}

func (p *DefaultBufferPool) Get() []byte {
	if p.buffer == nil {
		p.buffer = p.pool.Get().Use()
	}
	if buf := p.buffer.Get(); buf != nil {
		return buf
	}
	p.buffer.Release()
	p.buffer = p.pool.Get().Use()
	return p.buffer.Get()
}

func (p *DefaultBufferPool) Alloc(size uint) *Data {
	if p.buffer == nil {
		p.buffer = p.pool.Get().Use()
	}
	return p.buffer.Alloc(size)
}
