package pool

import (
	"gitee.com/sy_183/common/option"
	"sync"
)

type SyncPool[O any] struct {
	p     sync.Pool
	newFn func(p *SyncPool[O]) O
	syncLimiter
}

func NewSyncPool[O any](new func(p *SyncPool[O]) O, options ...option.AnyOption) *SyncPool[O] {
	p := &SyncPool[O]{newFn: new}
	p.p.New = p.new
	for _, opt := range options {
		opt.Apply(p)
	}
	return p
}

func ProvideSyncPool[O any](new func(p Pool[O]) O, options ...option.AnyOption) Pool[O] {
	return NewSyncPool(func(p *SyncPool[O]) O {
		return new(p)
	}, options...)
}

func (p *SyncPool[O]) new() any {
	return p.newFn(p)
}

func (p *SyncPool[O]) Get() (o O) {
	if p.alloc() {
		return p.p.Get().(O)
	}
	return
}

func (p *SyncPool[O]) Put(obj O) {
	p.p.Put(obj)
	p.release()
}
