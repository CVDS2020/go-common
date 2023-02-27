package pool

import "sync"

type StackPool[O any] struct {
	stack []O
	mu    sync.Mutex
}

func NewStackPool[O any](new func(p *StackPool[O]) O, size uint) *StackPool[O] {
	p := &StackPool[O]{stack: make([]O, size)}
	for i := range p.stack {
		p.stack[i] = new(p)
	}
	return p
}

func StackPoolProvider[O any](size uint) PoolProvider[O] {
	return func(new func(p Pool[O]) O) Pool[O] {
		return NewStackPool(func(p *StackPool[O]) O {
			return new(p)
		}, size)
	}
}

func (p *StackPool[O]) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.stack)
}

func (p *StackPool[O]) Cap() int {
	return cap(p.stack)
}

func (p *StackPool[O]) Get() (o O) {
	p.mu.Lock()
	defer p.mu.Unlock()
	l := len(p.stack)
	if l == 0 {
		return
	}
	o = p.stack[l-1]
	p.stack = p.stack[:l-1]
	return
}

func (p *StackPool[O]) Put(obj O) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.stack) == cap(p.stack) {
		return
	}
	p.stack = append(p.stack, obj)
}
