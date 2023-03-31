package slice

import _ "unsafe"

//go:linkname GoPanicIndex runtime.goPanicIndex
func GoPanicIndex(x int, y int)

//go:linkname GoPanicSliceAlen runtime.goPanicSliceAlen
func GoPanicSliceAlen(x int, y int)

//go:linkname GoPanicSliceB runtime.goPanicSliceB
func GoPanicSliceB(x int, y int)

type Chunks[E any] [][]E

func Join[E any](ss [][]E) []E {
	switch len(ss) {
	case 0:
		return nil
	case 1:
		return ss[0]
	case 2:
		rs := ss[0]
		return append(rs[:len(rs):len(rs)], ss[1]...)
	default:
		var n int
		for _, s := range ss {
			n += len(s)
		}
		rs := make([]E, n)
		sp := copy(rs, ss[0])
		for _, es := range ss[1:] {
			sp += copy(rs[sp:], es)
		}
		return rs
	}
}

func Join1N[E any](s []E, ss [][]E) []E {
	switch len(ss) {
	case 0:
		return s
	case 1:
		return append(s[:len(s):len(s)], ss[0]...)
	default:
		var n = len(s)
		for _, s := range ss {
			n += len(s)
		}
		rs := make([]E, n)
		sp := copy(rs, s)
		for _, es := range ss {
			sp += copy(rs[sp:], es)
		}
		return rs
	}
}

func JoinN1[E any](ss [][]E, s []E) []E {
	switch len(ss) {
	case 0:
		return s
	case 1:
		rs := ss[0]
		return append(rs[:len(rs):len(rs)], s...)
	default:
		var n = len(s)
		for _, s := range ss {
			n += len(s)
		}
		rs := make([]E, n)
		sp := copy(rs, ss[0])
		for _, es := range ss[1:] {
			sp += copy(rs[sp:], es)
		}
		copy(rs[sp:], s)
		return rs
	}
}

func Join11[E any](s1, s2 []E) []E {
	return append(s1[:len(s1):len(s1)], s2...)
}

func Join1N1[E any](s1 []E, ss [][]E, s2 []E) []E {
	switch len(ss) {
	case 0:
		return append(s1[:len(s1):len(s1)], s2...)
	default:
		var n = len(s1) + len(s2)
		for _, s := range ss {
			n += len(s)
		}
		rs := make([]E, n)
		sp := copy(rs, s1)
		for _, es := range ss {
			sp += copy(rs[sp:], es)
		}
		copy(rs[sp:], s2)
		return rs
	}
}

func Join111[E any](s1, s2, s3 []E) []E {
	rs := make([]E, len(s1)+len(s2)+len(s3))
	sp := copy(rs, s1)
	sp += copy(rs[sp:], s2)
	copy(rs[sp:], s3)
	return rs
}

func Locate[E any](cs [][]E, index int) (int, int) {
	i := index
	for j, chunk := range cs {
		if i < len(chunk) {
			return j, i
		}
		i -= len(chunk)
	}
	GoPanicIndex(index, index-i)
	return 0, 0
}

func LocateStart[E any](cs [][]E, start int) (int, int) {
	if len(cs) == 0 {
		return 0, 0
	}
	i := start
	for j, chunk := range cs {
		if i < len(chunk) {
			return j, i
		}
		i -= len(chunk)
	}
	if i > 0 {
		GoPanicSliceAlen(start, start-i)
	}
	return len(cs) - 1, len(cs[len(cs)-1])
}

func LocateEnd[E any](cs [][]E, end int) (int, int) {
	if len(cs) == 0 {
		return 0, 0
	}
	i := end
	for j, chunk := range cs {
		if i <= len(chunk) {
			return j, i
		}
		i -= len(chunk)
	}
	GoPanicSliceAlen(end, end-i)
	return 0, 0
}

func TrimChunks[E any](chunks [][]E) [][]E {
	return TrimChunksEnd(TrimChunksStart(chunks))
}

func TrimChunksStart[E any](chunks [][]E) [][]E {
	for len(chunks) > 0 {
		if len(chunks[0]) == 0 {
			chunks = chunks[1:]
		} else {
			break
		}
	}
	return chunks
}

func TrimChunksEnd[E any](chunks [][]E) [][]E {
	for i := len(chunks) - 1; i >= 0; i-- {
		if len(chunks[i]) == 0 {
			chunks = chunks[:i]
		} else {
			break
		}
	}
	return chunks
}

func (cs Chunks[E]) Len() int {
	var l int
	for _, c := range cs {
		l += len(c)
	}
	return l
}

func (cs Chunks[E]) Get(i int) E {
	x, y := Locate(cs, i)
	return cs[x][y]
}

func (cs Chunks[E]) Pointer(i int) *E {
	x, y := Locate(cs, i)
	return &cs[x][y]
}

func (cs Chunks[E]) First() (e E, exist bool) {
	cs = cs.TrimStart()
	if len(cs) == 0 {
		return
	}
	return cs[0][0], true
}

func (cs Chunks[E]) FirstPointer() *E {
	cs = cs.TrimStart()
	if len(cs) == 0 {
		return nil
	}
	return &cs[0][0]
}

func (cs Chunks[E]) Last() (e E, exist bool) {
	cs = cs.TrimEnd()
	if len(cs) == 0 {
		return
	}
	c := cs[len(cs)-1]
	return c[len(c)-1], true
}

func (cs Chunks[E]) LastPointer() *E {
	cs = cs.TrimEnd()
	if len(cs) == 0 {
		return nil
	}
	c := cs[len(cs)-1]
	return &c[len(c)-1]
}

func (cs Chunks[E]) TrimStart() Chunks[E] {
	return TrimChunksStart(cs)
}

func (cs Chunks[E]) TrimEnd() Chunks[E] {
	return TrimChunksEnd(cs)
}

func (cs Chunks[E]) Trim() Chunks[E] {
	return TrimChunks(cs)
}

func (cs Chunks[E]) sliceFrom(sx, sy int) []E {
	if len(cs) == 0 {
		return nil
	}
	return Join1N(cs[sx][sy:], cs[sx+1:])
}

func (cs Chunks[E]) sliceTo(ex, ey int) []E {
	if len(cs) == 0 {
		return nil
	}
	return JoinN1(cs[:ex], cs[ex][:ey])
}

func (cs Chunks[E]) slice(sx, sy, ex, ey int) []E {
	if len(cs) == 0 {
		return nil
	}
	switch {
	case sx > ex:
		return nil
	case sx == ex:
		return cs[sx][sy:ey]
	case sx == ex-1:
		return Join11(cs[sx][sy:], cs[ex][:ey])
	default:
		return Join1N1(cs[sx][sy:], cs[sx+1:ex], cs[ex][:ey])
	}
}

func (cs Chunks[E]) Slice(start, end int) []E {
	if start >= 0 && end >= 0 {
		if start > end {
			GoPanicSliceB(start, end)
		}
		cs = cs.TrimStart().TrimEnd()
		sx, sy := LocateStart(cs, start)
		ex, ey := LocateEnd(cs, end)
		return cs.slice(sx, sy, ex, ey)
	}
	cs = cs.TrimStart().TrimEnd()
	switch {
	case start < 0 && end < 0:
		return Join(cs)
	case start < 0:
		return cs.sliceTo(LocateEnd(cs, end))
	default:
		return cs.sliceFrom(LocateStart(cs, start))
	}
}

func (cs Chunks[E]) cutFrom(sx, sy int) Chunks[E] {
	if len(cs) == 0 {
		return nil
	}
	if sy == 0 {
		return cs[sx:]
	}
	return Join11(Chunks[E]{cs[sx][sy:]}, cs[sx+1:])
}

func (cs Chunks[E]) cutTo(ex, ey int) Chunks[E] {
	if len(cs) == 0 {
		return nil
	}
	if ey == len(cs[ex]) {
		return cs[:ex+1]
	}
	return Join11(cs[:ex], Chunks[E]{cs[ex][:ey]})
}

func (cs Chunks[E]) cut(sx, sy, ex, ey int) Chunks[E] {
	if len(cs) == 0 {
		return nil
	}
	switch {
	case sx > ex:
		return nil
	case sx == ex:
		return Chunks[E]{cs[sx][sy:ey]}
	case sx == ex-1:
		return Chunks[E]{cs[sx][sy:], cs[ex][:ey]}
	default:
		if sy == 0 && ey == len(cs[ex]) {
			return cs[sx : ex+1]
		}
		return Join111(Chunks[E]{cs[sx][sy:]}, cs[sx+1:ex], Chunks[E]{cs[ex][:ey]})
	}
}

func (cs Chunks[E]) Cut(start, end int) Chunks[E] {
	if start >= 0 && end >= 0 {
		if start > end {
			GoPanicSliceB(start, end)
		}
		cs = cs.TrimStart().TrimEnd()
		sx, sy := LocateStart(cs, start)
		ex, ey := LocateEnd(cs, end)
		return cs.cut(sx, sy, ex, ey)
	}
	cs = cs.TrimStart().TrimEnd()
	switch {
	case start < 0 && end < 0:
		return cs
	case start < 0:
		return cs.cutTo(LocateEnd(cs, end))
	default:
		return cs.cutFrom(LocateStart(cs, start))
	}
}

func (cs Chunks[E]) Set(i int, e E) {
	*cs.Pointer(i) = e
}

func (cs Chunks[E]) GetAndSet(i int, e E) (old E) {
	ptr := cs.Pointer(i)
	old = *ptr
	*ptr = e
	return
}

func (cs Chunks[E]) Swap(i, j int) {
	x1, y1 := Locate(cs, i)
	x2, y2 := Locate(cs, j)
	cs[x1][y1], cs[x2][y2] = cs[x2][y2], cs[x1][y1]
}
