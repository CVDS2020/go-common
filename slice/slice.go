package slice

import (
	"gitee.com/sy_183/common/generic"
	"gitee.com/sy_183/common/uns/math"
	"unsafe"
)

func _make[E any](ptr unsafe.Pointer, len, cap int) []E {
	return StructToSlice[E](Struct{Ptr: ptr, Len: len, Cap: cap})
}

func Make[E any](ptr unsafe.Pointer, len, cap int) []E {
	es := generic.Size[E]()
	mem, overflow := math.MulUintptr(es, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(es, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return _make[E](ptr, len, cap)
}

func SliceConvert[EI, EO any](s []EI, len, cap int) []EO {
	ss := StructOfPtr[EI](&s)
	return Make[EO](ss.Ptr, len, cap)
}

func Join[E any](ss ...[]E) []E {
	switch len(ss) {
	case 0:
		return nil
	case 1:
		return append([]E(nil), ss[0]...)
	case 2:
		return append(ss[0], ss[1]...)
	default:
		var n int
		for _, s := range ss {
			n += len(s)
		}
		s := make([]E, n)
		sp := copy(s, ss[0])
		for _, es := range ss[1:] {
			sp += copy(s[sp:], es)
		}
		return s
	}
}

func min[N generic.Ordered](a, b N) N {
	if a < b {
		return a
	} else {
		return b
	}
}

func Assign[E any](s []E, _len, _cap int) []E {
	if cap(s) >= _cap {
		return s[:_len:_cap]
	}
	n := make([]E, _len, _cap)
	copy(n, s)
	return n
}

func AssignLen[E any](s []E, _len int) []E {
	if cap(s) >= _len {
		return s[:_len]
	}
	n := make([]E, _len)
	copy(n, s)
	return n
}

func Map[E any](s []E, mapper func(E) E) {

}
