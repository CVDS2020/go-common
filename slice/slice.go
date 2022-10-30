package slice

import (
	"gitee.com/sy_183/common/generic"
)

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
