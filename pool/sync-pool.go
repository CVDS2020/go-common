package pool

import "sync"

type SyncPool[O any] struct {
	p     sync.Pool
	newFn func(p *SyncPool[O]) O
}

func NewSyncPool[O any](new func(p *SyncPool[O]) O) *SyncPool[O] {
	p := &SyncPool[O]{newFn: new}
	p.p.New = p.new
	return p
}

func ProvideSyncPool[O any](new func(p Pool[O]) O) Pool[O] {
	return NewSyncPool(func(p *SyncPool[O]) O {
		return new(p)
	})
}

func (p *SyncPool[O]) new() any {
	return p.newFn(p)
}

func (p *SyncPool[O]) Get() O {
	return p.p.Get().(O)
}

func (p *SyncPool[O]) Put(obj O) {
	p.p.Put(obj)
}
