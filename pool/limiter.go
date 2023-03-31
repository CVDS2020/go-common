package pool

import "sync/atomic"

type limiter struct {
	limit     int64
	allocated int64
}

func (l *limiter) setLimit(limit int64) {
	l.limit = limit
}

func (l *limiter) alloc() bool {
	if l.limit > 0 {
		if l.allocated >= l.limit {
			return false
		}
		if l.allocated++; l.allocated <= 0 {
			panic(NewInvalidReferenceError(l.allocated, l.allocated-1))
		}
	}
	return true
}

func (l *limiter) release() {
	if l.limit > 0 {
		if l.allocated--; l.allocated < 0 {
			panic(NewInvalidReferenceError(l.allocated, l.allocated+1))
		}
	}
}

type syncLimiter struct {
	limit     int64
	allocated atomic.Int64
}

func (l *syncLimiter) setLimit(limit int64) {
	l.limit = limit
}

func (l *syncLimiter) alloc() bool {
	if l.limit > 0 {
		if allocated := l.allocated.Add(1); allocated > l.limit {
			l.allocated.Add(-1)
			return false
		} else if allocated <= 0 {
			panic(NewInvalidReferenceError(allocated, allocated-1))
		}
	}
	return true
}

func (l *syncLimiter) release() {
	if l.limit > 0 {
		if allocated := l.allocated.Add(-1); allocated < 0 {
			panic(NewInvalidReferenceError(allocated, allocated+1))
		}
	}
}
