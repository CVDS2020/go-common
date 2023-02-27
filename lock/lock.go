package lock

import "sync"

type RLocker interface {
	RLock()

	RUnlock()
}

func LockDo(locker sync.Locker, f func()) {
	locker.Lock()
	defer locker.Unlock()
	f()
}

func LockGet[V any](locker sync.Locker, f func() V) V {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

func LockGetDouble[V1, V2 any](locker sync.Locker, f func() (V1, V2)) (V1, V2) {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

func LockGetTriple[V1, V2, V3 any](locker sync.Locker, f func() (V1, V2, V3)) (V1, V2, V3) {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

func RLockDo(locker RLocker, f func()) {
	locker.RLock()
	defer locker.RUnlock()
	f()
}

func RLockGet[V any](locker RLocker, f func() V) V {
	locker.RLock()
	defer locker.RUnlock()
	return f()
}

func RLockGetDouble[V1, V2 any](locker RLocker, f func() (V1, V2)) (V1, V2) {
	locker.RLock()
	defer locker.RUnlock()
	return f()
}

func RLockGetTriple[V1, V2, V3 any](locker RLocker, f func() (V1, V2, V3)) (V1, V2, V3) {
	locker.RLock()
	defer locker.RUnlock()
	return f()
}
