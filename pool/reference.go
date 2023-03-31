package pool

import "sync/atomic"

type Reference interface {
	Release() bool

	AddRef()
}

type AtomicRef atomic.Int64

func (r *AtomicRef) Release() bool {
	if ref := (*atomic.Int64)(r).Add(-1); ref == 0 {
		return true
	} else if ref < 0 {
		panic(NewInvalidReferenceError(ref, ref+1))
	}
	return false
}

func (r *AtomicRef) AddRef() {
	if ref := (*atomic.Int64)(r).Add(1); ref <= 0 {
		panic(NewInvalidReferenceError(ref, ref-1))
	}
}

type Relations struct {
	AtomicRef
	relations []Reference
}

func (r *Relations) AddRelation(relation Reference) {
	r.relations = append(r.relations, relation)
}

func (r *Relations) Clear() {
	for _, relation := range r.relations {
		relation.Release()
	}
	r.relations = r.relations[:0]
}

func (r *Relations) Release() bool {
	if r.AtomicRef.Release() {
		r.Clear()
		return true
	}
	return false
}
