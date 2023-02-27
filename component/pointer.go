package component

import "sync"

type Pointer[E any] struct {
	Init        func() *E
	initializer sync.Once
	elem        *E
}

func NewPointer[E any](init func() *E) *Pointer[E] {
	return &Pointer[E]{Init: init}
}

func (c *Pointer[E]) Get() *E {
	if c.elem != nil {
		return c.elem
	}
	c.initializer.Do(func() {
		c.elem = c.Init()
	})
	return c.elem
}
