package pool

import "sync"

type SlicePool[O any] struct {
	slice []O
	newFn func(p *SlicePool[O]) O
	mu    sync.Mutex
}

func NewSlicePool[O any](new func(p *SlicePool[O]) O) *SlicePool[O] {
	return &SlicePool[O]{newFn: new}
}

func ProvideSlicePool[O any](new func(p Pool[O]) O) Pool[O] {
	return NewSlicePool(func(p *SlicePool[O]) O {
		return new(p)
	})
}

func (p *SlicePool[O]) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.slice)
}

func (p *SlicePool[O]) Cap() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return cap(p.slice)
}

func (p *SlicePool[O]) Get() (o O) {
	p.mu.Lock()
	l := len(p.slice)
	if l == 0 {
		p.mu.Unlock()
		return p.newFn(p)
	}
	o = p.slice[l-1]
	p.slice = p.slice[:l-1]
	p.mu.Unlock()
	return
}

func (p *SlicePool[O]) Put(obj O) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.slice = append(p.slice, obj)
}
