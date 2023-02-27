package lifecycle

import (
	"gitee.com/sy_183/common/lock"
	"sync"
	"sync/atomic"
)

type Future[V any] interface {
	Complete(value V)
}

type WaitableFuture[V any] interface {
	Future[V]

	Wait() V
}

type NopFuture[V any] struct{}

func (f NopFuture[V]) Complete(value V) {}

type ChanFuture[V any] chan V

func (c ChanFuture[V]) Complete(value V) { c <- value }
func (c ChanFuture[V]) Wait() V          { return <-c }

type WaiterFuture[V any] struct {
	waiter sync.WaitGroup
	value  V
	once   atomic.Bool
}

func NewWaiterFuture[V any]() *WaiterFuture[V] {
	future := new(WaiterFuture[V])
	future.waiter.Add(1)
	return future
}

func (w *WaiterFuture[V]) Complete(value V) {
	if w.once.CompareAndSwap(false, true) {
		w.value = value
		w.waiter.Done()
	}
}

func (w *WaiterFuture[V]) Wait() V {
	w.waiter.Wait()
	return w.value
}

type CallbackFuture[V any] struct {
	callback atomic.Pointer[func(value V)]
	value    atomic.Pointer[V]
	once     atomic.Bool
	done     atomic.Bool
}

func NewCallbackFuture[V any](callback func(value V)) *CallbackFuture[V] {
	future := new(CallbackFuture[V])
	if callback != nil {
		future.callback.Store(&callback)
	}
	return future
}

func (f CallbackFuture[V]) Complete(value V) {
	if f.once.CompareAndSwap(false, true) {
		f.value.Store(&value)
		if callback := f.callback.Load(); callback != nil && f.done.CompareAndSwap(false, true) {
			(*callback)(value)
		}
	}
}

func (f CallbackFuture[V]) SetCallback(callback func(value V)) {
	f.callback.Store(&callback)
	if value := f.value.Load(); value != nil && f.done.CompareAndSwap(false, true) {
		callback(*value)
	}
}

type Futures[V any] []Future[V]

func (fs Futures[V]) Complete(value V) {
	for _, future := range fs {
		future.Complete(value)
	}
}

type SyncFutures[V any] struct {
	futures Futures[V]
	sync.Mutex
}

func NewSyncFutures[V any]() *SyncFutures[V] {
	return new(SyncFutures[V])
}

func (f *SyncFutures[V]) Append(future Future[V]) {
	if future == nil {
		return
	}
	lock.LockDo(f, func() { f.futures = append(f.futures, future) })
}

func (f *SyncFutures[V]) LoadAndReset() Futures[V] {
	return lock.LockGet(f, func() (futures Futures[V]) {
		futures = append(futures, f.futures...)
		f.futures = f.futures[:0]
		return
	})
}
