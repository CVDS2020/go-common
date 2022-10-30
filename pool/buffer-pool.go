package pool

import "sync"

// A BufferPool is a type-safe wrapper around a sync.Pool.
type BufferPool struct {
	p *sync.Pool
}

// NewBufferPool constructs a new Pool.
func NewBufferPool(size uint) *BufferPool {
	return &BufferPool{p: &sync.Pool{
		New: func() interface{} {
			return &Buffer{bs: make([]byte, 0, size)}
		},
	}}
}

// Get retrieves a Buffer from the pool, creating one if necessary.
func (p *BufferPool) Get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p *BufferPool) put(buf *Buffer) {
	p.p.Put(buf)
}
