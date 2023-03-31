package pool

import "gitee.com/sy_183/common/option"

type DefaultBufferPool struct {
	pool   Pool[*Buffer]
	buffer *Buffer
}

func NewDefaultBufferPool(size uint, reversed uint, poolProvider PoolProvider[*Buffer], poolOptions ...option.AnyOption) *DefaultBufferPool {
	return &DefaultBufferPool{pool: poolProvider(func(p Pool[*Buffer]) *Buffer {
		return NewBuffer(size, reversed).SetOnReleased(func(buffer *Buffer) {
			p.Put(buffer)
		})
	}, poolOptions...)}
}

func (p *DefaultBufferPool) getBuffer() *Buffer {
	if p.buffer == nil {
		buffer := p.pool.Get()
		if buffer == nil {
			return nil
		}
		p.buffer = buffer.Use()
	}
	return p.buffer
}

func (p *DefaultBufferPool) Get() []byte {
	if buffer := p.getBuffer(); buffer != nil {
		if buf := buffer.Get(); buf != nil {
			return buf
		}
		p.buffer.Release()
		p.buffer = nil
		if buffer := p.getBuffer(); buffer != nil {
			return buffer.Get()
		}
	}
	return nil
}

func (p *DefaultBufferPool) Alloc(size uint) *Data {
	if p.buffer == nil {
		p.buffer = p.pool.Get().Use()
	}
	return p.buffer.Alloc(size)
}
