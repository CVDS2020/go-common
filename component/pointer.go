package component

import "sync"

type Pointer[E any] struct {
	Init        func() *E
	initializer sync.Once
	elem        *E
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
