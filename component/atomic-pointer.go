package component

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type AtomicPointer[E any] struct {
	Init        func() *E
	initializer sync.Once
	elem        unsafe.Pointer
}

func (p *AtomicPointer[E]) store(e *E) {
	atomic.StorePointer(&p.elem, unsafe.Pointer(e))
}

func (p *AtomicPointer[E]) load() *E {
	return (*E)(atomic.LoadPointer(&p.elem))
}

func (p *AtomicPointer[E]) Get() *E {
	if ep := p.load(); ep != nil {
		return ep
	}
	p.initializer.Do(func() {
		p.store(p.Init())
	})
	return p.load()
}

func (p *AtomicPointer[E]) Set(e *E) {
	p.store(e)
}
