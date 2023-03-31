package pool

import (
	"gitee.com/sy_183/common/option"
	"sync"
)

type SlicePool[O any] struct {
	slice []O
	newFn func(p *SlicePool[O]) O
	limiter
	mu sync.Mutex
}

func NewSlicePool[O any](new func(p *SlicePool[O]) O, options ...option.AnyOption) *SlicePool[O] {
	p := &SlicePool[O]{newFn: new}
	for _, opt := range options {
		opt.Apply(p)
	}
	return p
}

func ProvideSlicePool[O any](new func(p Pool[O]) O, options ...option.AnyOption) Pool[O] {
	return NewSlicePool(func(p *SlicePool[O]) O {
		return new(p)
	}, options...)
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

func (p *SlicePool[O]) getCached() (o O, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	ok = true
	if !p.alloc() {
		return
	}
	l := len(p.slice)
	if l == 0 {
		ok = false
		return
	}
	o = p.slice[l-1]
	p.slice = p.slice[:l-1]
	return
}

func (p *SlicePool[O]) Get() (o O) {
	var ok bool
	o, ok = p.getCached()
	if ok {
		return
	}
	return p.newFn(p)
}

func (p *SlicePool[O]) Put(obj O) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.slice = append(p.slice, obj)
	p.release()
}
