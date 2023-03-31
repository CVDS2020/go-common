package container

import _ "unsafe"

//go:linkname goPanicIndex runtime.goPanicIndex
func goPanicIndex(x int, y int)

//go:linkname goPanicSliceAcap runtime.goPanicSliceAcap
func goPanicSliceAcap(x int, y int)

//go:linkname goPanicSliceB runtime.goPanicSliceB
func goPanicSliceB(x int, y int)

type Queue[E any] struct {
	raw  []E
	head int
	tail int
	len  int
	cap  int
}

func NewQueue[E any](cap int) *Queue[E] {
	return &Queue[E]{
		raw: make([]E, cap, cap),
		cap: cap,
	}
}

func (q *Queue[E]) Len() int {
	return q.len
}

func (q *Queue[E]) Cap() int {
	return q.cap
}

func (q *Queue[E]) sliceTo(end int) (prefix, suffix []E) {
	if end < 0 {
		end = q.len
	} else if end > q.cap {
		goPanicSliceAcap(end, q.cap)
	}
	if end += q.head; end > q.cap {
		end -= q.cap
		return q.raw[q.head:], q.raw[:end]
	} else {
		return q.raw[q.head:end], nil
	}
}

func (q *Queue[E]) Slice(start, end int) (prefix, suffix []E) {
	if start < 0 {
		return q.sliceTo(end)
	}
	if end < 0 {
		end = q.len
	} else if end > q.cap {
		goPanicSliceAcap(end, q.cap)
	}
	if start > end {
		goPanicSliceB(start, end)
	}
	if start += q.head; start >= q.cap {
		start -= q.cap
		end += q.head - q.cap
		return nil, q.raw[start:end]
	} else if end += q.head; end > q.cap {
		end -= q.cap
		return q.raw[start:], q.raw[:end]
	} else {
		return q.raw[start:end], nil
	}
}

func (q *Queue[E]) tailGrow() {
	if q.tail++; q.tail == q.cap {
		q.tail = 0
	}
	q.len++
}

func (q *Queue[E]) headGrow() {
	if q.head--; q.head < 0 {
		q.head = q.cap - 1
	}
	q.len++
}

func (q *Queue[E]) Push(e E) bool {
	if q.len == q.cap {
		return false
	}
	q.raw[q.tail] = e
	q.tailGrow()
	return true
}

func (q *Queue[E]) Use() *E {
	if q.len == q.cap {
		return nil
	}
	q.tailGrow()
	return &q.raw[q.tail]
}

func (q *Queue[E]) PushHead(e E) bool {
	if q.len == q.cap {
		return false
	}
	q.headGrow()
	q.raw[q.head] = e
	return true
}

func (q *Queue[E]) UseHead() *E {
	if q.len == q.cap {
		return nil
	}
	q.headGrow()
	return &q.raw[q.head]
}

func (q *Queue[E]) Pop() (e E, ok bool) {
	if q.len == 0 {
		return
	}
	e, ok = q.raw[q.head], true
	if q.head++; q.head == q.cap {
		q.head = 0
	}
	q.len--
	return
}

func (q *Queue[E]) PopTail() (e E, ok bool) {
	if q.len == 0 {
		return
	}
	if q.tail--; q.tail < 0 {
		q.tail = q.cap - 1
	}
	e, ok = q.raw[q.tail], true
	q.len--
	return
}

func (q *Queue[E]) PopTo(es []E) (n int) {
	end := len(es)
	if end > q.len {
		end = q.len
	}
	if end += q.head; end > q.cap {
		end -= q.cap
		n = copy(es, q.raw[q.head:])
		n += copy(es[n:], q.raw[:end])
	} else {
		n = copy(es, q.raw[q.head:end])
	}
	q.len -= n
	q.head = end
	return
}

func (q *Queue[E]) PopAll() []E {
	es := make([]E, q.len)
	q.PopTo(es)
	return es
}

func (q *Queue[E]) Get(i int) E {
	return *q.Pointer(i)
}

func (q *Queue[E]) HeadPointer() *E {
	if q.len > 0 {
		return &q.raw[q.head]
	}
	return nil
}

func (q *Queue[E]) Head() (e E, exist bool) {
	if q.len > 0 {
		return q.raw[q.head], true
	}
	return
}

func (q *Queue[E]) tailIndex() int {
	tail := q.tail - 1
	if tail < 0 {
		tail = q.cap - 1
	}
	return tail
}

func (q *Queue[E]) TailPointer() *E {
	if q.len > 0 {
		return &q.raw[q.tailIndex()]
	}
	return nil
}

func (q *Queue[E]) Tail() (e E, exist bool) {
	if q.len > 0 {
		return q.raw[q.tailIndex()], true
	}
	return
}

func (q *Queue[E]) checkIndex(i int) {
	if i >= q.len || i < 0 {
		goPanicIndex(i, q.len)
	}
}

func (q *Queue[E]) rawIndex(i int) int {
	if i += q.head; i >= q.cap {
		return i - q.cap
	}
	return i
}

func (q *Queue[E]) Pointer(i int) *E {
	q.checkIndex(i)
	return &q.raw[q.rawIndex(i)]
}

func (q *Queue[E]) Set(i int, e E) {
	*q.Pointer(i) = e
}

func (q *Queue[E]) Swap(i, j int) {
	q.checkIndex(i)
	q.checkIndex(j)
	i, j = q.rawIndex(i), q.rawIndex(j)
	q.raw[i], q.raw[j] = q.raw[j], q.raw[i]
}

func (q *Queue[E]) Clear() {
	q.head = q.tail
	q.len = 0
}
