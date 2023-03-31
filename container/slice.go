package container

import "gitee.com/sy_183/common/slice"

type Slice[E any] interface {
	Len() int

	Cap() int

	Get(i int) E

	Pointer(i int) *E

	Cut(from, to int) Slice[E]

	CutFrom(from int) Slice[E]

	CutTo(to int) Slice[E]

	CutCap(from, to int, cap int) Slice[E]

	CutCapTo(to int, cap int) Slice[E]

	Append(es ...E) Slice[E]

	Delete(i int) Slice[E]

	DeleteIf(delete func(i int) bool) Slice[E]

	Range(f func(i int) bool)

	Set(i int, e E)

	Swap(i, j int)

	Clear() Slice[E]
}

type GoSlice[E any] []E

func (s GoSlice[E]) Len() int {
	return len(s)
}

func (s GoSlice[E]) Cap() int {
	return cap(s)
}

func (s GoSlice[E]) Get(i int) E {
	return s[i]
}

func (s GoSlice[E]) Pointer(i int) *E {
	return &s[i]
}

func (s GoSlice[E]) Cut(from, to int) Slice[E] {
	return s[from:to]
}

func (s GoSlice[E]) CutFrom(from int) Slice[E] {
	return s[from:]
}

func (s GoSlice[E]) CutTo(to int) Slice[E] {
	return s[:to]
}

func (s GoSlice[E]) CutCap(from, to int, cap int) Slice[E] {
	return s[from:to:cap]
}

func (s GoSlice[E]) CutCapTo(to int, cap int) Slice[E] {
	return s[:to:cap]
}

func (s GoSlice[E]) Append(es ...E) Slice[E] {
	return append(s, es...)
}

func (s GoSlice[E]) Delete(i int) Slice[E] {
	return append(s[:i], s[i+1:]...)
}

func (s GoSlice[E]) DeleteIf(delete func(i int) bool) Slice[E] {
	return GoSlice[E](slice.SliceDelete(s, delete))
}

func (s GoSlice[E]) Range(f func(i int) bool) {
	for i := range s {
		if !f(i) {
			break
		}
	}
}

func (s GoSlice[E]) Set(i int, e E) {
	s[i] = e
}

func (s GoSlice[E]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GoSlice[E]) Clear() Slice[E] {
	return s[:0]
}
