package pool

import "fmt"

type InvalidReferenceError struct {
	CurRef  int64
	LastRef int64
}

func NewInvalidReferenceError(cur, last int64) InvalidReferenceError {
	return InvalidReferenceError{CurRef: cur, LastRef: last}
}

func (e InvalidReferenceError) Error() string {
	if e.CurRef > e.LastRef {
		return fmt.Sprintf("非法的引用计数增加[%d -> %d]", e.LastRef, e.CurRef)
	} else if e.LastRef < e.CurRef {
		return fmt.Sprintf("非法的引用计数减少[%d -> %d]", e.LastRef, e.CurRef)
	} else {
		return fmt.Sprintf("非法的引用计数(%d)", e.CurRef)
	}
}

type AllocError struct {
	Target string
}

func NewAllocError(target string) error {
	return AllocError{Target: target}
}

func (e AllocError) Error() string {
	return fmt.Sprintf("申请%s失败", e.Target)
}
