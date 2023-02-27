package lock

import "sync"

type LockWrapper[T any] struct {
	Elem T
	sync.Mutex
}
