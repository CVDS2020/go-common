package pool

import (
	"fmt"
	"gitee.com/sy_183/common/sgr"
	"sync/atomic"
)

type SwapPool[O any] struct {
	// queue capacity
	qc int
	// queue pool
	qp *SyncPool[*[]O]
	// consumer queue
	cq atomic.Pointer[[]O]
	// producer queue
	pq atomic.Pointer[[]O]
	// wait queue
	wq atomic.Pointer[[]O]
	// object creator
	creator func(p Pool[O]) O
}

func NewSwapPool[O any](cap int, creator func(p Pool[O]) O) *SwapPool[O] {
	p := &SwapPool[O]{
		qc: cap,
		qp: NewSyncPool(func(p *SyncPool[*[]O]) *[]O {
			qp := new([]O)
			*qp = make([]O, 0, cap)
			return qp
		}),
		creator: creator,
	}
	p.cq.Store(p.qp.Get())
	p.pq.Store(p.qp.Get())
	return p
}

func (s *SwapPool[O]) Get() (o O) {
	cqp := s.cq.Swap(nil)
	fmt.Println(sgr.WrapColor("get consumer queue", sgr.FgBlue))
	if cqp == nil {
		panic(fmt.Errorf("swap pool does not allow concurrent get object"))
	}
	pq := *cqp
	l := len(pq)
	if l == 0 {
		fmt.Println(sgr.WrapColor("get from creator start", sgr.FgRed))
		o = s.creator(s)
		fmt.Println(sgr.WrapColor("get from creator end", sgr.FgRed))
	} else {
		o = pq[l-1]
		*cqp = pq[:l-1]
		fmt.Println(sgr.WrapColor("get from consumer queue", sgr.FgGreen))
	}
	if !s.cq.CompareAndSwap(nil, cqp) {
		fmt.Println(sgr.WrapColor("consumer queue swapped, put queue to wait queue", sgr.FgCyan))
		if wqp := s.wq.Swap(cqp); wqp != nil {
			fmt.Println(sgr.WrapColor("wait queue has last queue, put last queue to queue pool", sgr.FgMagenta))
			*wqp = (*wqp)[:0]
			s.qp.Put(wqp)
		}
	} else {
		fmt.Println(sgr.WrapColor("revert consumer queue", sgr.FgBlue))
	}
	return
}

func (s *SwapPool[O]) Put(val O) {
	pqp := s.pq.Load()
	fmt.Println(sgr.WrapColor("get producer queue", sgr.FgBrightBlue))
	if pqp == nil {
		fmt.Println(sgr.WrapColor("producer queue is nil, get it from wait queue", sgr.FgBrightGreen))
		pqp = s.wq.Swap(nil)
		if pqp == nil {
			fmt.Println(sgr.WrapColor("wait queue is nil, get producer queue from queue pool", sgr.FgBrightMagenta))
			pqp = s.qp.Get()
		}
		s.pq.Store(pqp)
	}
	fmt.Println(sgr.WrapColor("put object to producer queue", sgr.FgBrightGreen))
	cq := append(*pqp, val)
	*pqp = cq
	if len(cq) >= cap(cq) {
		fmt.Println(sgr.WrapColor("producer queue full, put it to consumer queue", sgr.FgBrightBlue))
		cqp := s.cq.Swap(pqp)
		if cqp != nil {
			fmt.Println(sgr.WrapColor("old consumer queue not used, clear and put it to producer queue", sgr.FgBrightBlue))
			*cqp = (*cqp)[:0]
		} else {
			fmt.Println(sgr.WrapColor("old consumer used, set producer queue nil", sgr.FgBrightCyan))
		}
		s.pq.Store(cqp)
	}
}
