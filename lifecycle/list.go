package lifecycle

import (
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/lock"
	"sync"
	"sync/atomic"
)

const ListFieldName = "$list"

type ListLifecycleHolder struct {
	Lifecycle
	index int
	list  *List

	closeAllOnStartError  atomic.Bool
	closeAllOnExit        atomic.Bool
	closeAllOnExitError   atomic.Bool
	stopStartOnStartError atomic.Bool
	closeBackOnExit       atomic.Bool
	closeBackOnExitError  atomic.Bool
}

func (h *ListLifecycleHolder) List() *List {
	return h.list
}

func (h *ListLifecycleHolder) SetCloseAllOnStartError(enable bool) *ListLifecycleHolder {
	h.closeAllOnStartError.Store(enable)
	return h
}

func (h *ListLifecycleHolder) SetCloseAllOnExit(enable bool) *ListLifecycleHolder {
	h.closeAllOnExit.Store(enable)
	return h
}

func (h *ListLifecycleHolder) SetCloseAllOnExitError(enable bool) *ListLifecycleHolder {
	h.closeAllOnExitError.Store(enable)
	return h
}

func (h *ListLifecycleHolder) SetStopStartOnStartError(enable bool) *ListLifecycleHolder {
	h.stopStartOnStartError.Store(enable)
	return h
}

func (h *ListLifecycleHolder) SetCloseBackOnExit(enable bool) *ListLifecycleHolder {
	h.closeBackOnExit.Store(enable)
	return h
}

func (h *ListLifecycleHolder) SetCloseBackOnExitError(enable bool) *ListLifecycleHolder {
	h.closeBackOnExitError.Store(enable)
	return h
}

type (
	listLifecycleContext = childLifecycleContext[*ListLifecycleHolder]
	listLifecycleChannel = childLifecycleChannel[*ListLifecycleHolder]
)

type List struct {
	Lifecycle
	lifecycle *DefaultLifecycle

	children     []*ListLifecycleHolder
	childrenLock sync.Mutex
	started      int

	closedChannel *listLifecycleChannel
}

func NewList() *List {
	l := &List{
		closedChannel: newChildLifecycleChannel[*ListLifecycleHolder](),
	}
	l.lifecycle = NewWithInterruptedStart(l.start)
	l.Lifecycle = l.lifecycle
	return l
}

func (l *List) Append(lifecycle Lifecycle) (*ListLifecycleHolder, error) {
	return lock.RLockGetDouble(l.lifecycle, func() (*ListLifecycleHolder, error) {
		if !l.lifecycle.Closed() {
			return nil, NewStateNotClosedError("")
		}
		child := &ListLifecycleHolder{
			Lifecycle: lifecycle,
			index:     len(l.children),
			list:      l,
		}
		child.SetCloseAllOnStartError(true)
		child.SetCloseAllOnExit(true)
		child.SetCloseAllOnExitError(true)
		child.SetStopStartOnStartError(true)
		child.SetCloseBackOnExit(true)
		child.SetCloseBackOnExitError(true)
		child.SetField(ListFieldName, l)
		l.children = append(l.children, child)
		return child, nil
	})
}

func (l *List) MustAppend(lifecycle Lifecycle) *ListLifecycleHolder {
	return assert.Must(l.Append(lifecycle))
}

func (l *List) shutdownChildren(children []*ListLifecycleHolder, exclude map[*ListLifecycleHolder]struct{}) {
	excluded := func(child *ListLifecycleHolder) bool {
		if exclude != nil {
			_, has := exclude[child]
			return has
		}
		return false
	}
	for i := len(children) - 1; i >= 0; i-- {
		if !excluded(children[i]) {
			future := make(ChanFuture[error], 1)
			children[i].Close(future)
			<-future
		}
	}
}

func (l *List) handleClosedSignal() (closeAll bool) {
	closeBackIndex := -1
	contexts := l.closedChannel.Pop()
	var closedSet map[*ListLifecycleHolder]struct{}
	for _, ctx := range contexts {
		child := ctx.Lifecycle
		if ctx.err != nil {
			if child.closeAllOnExit.Load() || child.closeAllOnExitError.Load() {
				closeAll = true
			} else if child.closeBackOnExit.Load() || child.closeBackOnExitError.Load() {
				if child.index < 0 || child.index < closeBackIndex {
					closeBackIndex = child.index
				}
			}
		} else {
			if child.closeAllOnExit.Load() {
				closeAll = true
			} else if child.closeBackOnExit.Load() {
				if child.index < 0 || child.index < closeBackIndex {
					closeBackIndex = child.index
				}
			}
		}
		if closedSet == nil {
			closedSet = make(map[*ListLifecycleHolder]struct{})
		}
		closedSet[child] = struct{}{}
	}
	if closeBackIndex >= 0 {
		l.shutdownChildren(l.children[closeBackIndex+1:l.started], closedSet)
	} else if closeAll {
		l.shutdownChildren(l.children[:l.started], closedSet)
	}
	return
}

func (l *List) start(_ Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error) {
next:
	for i, child := range l.children {
		future := make(ChanFuture[error], 1)
		child.AddStartedFuture(future)
		child.Background()
		l.started = i + 1
		for {
			select {
			case err := <-future:
				if err != nil {
					if child.closeAllOnStartError.Load() {
						l.shutdownChildren(l.children[:l.started], nil)
						return nil, NewInterruptedError("生命周期列表", "启动")
					} else if child.stopStartOnStartError.Load() {
						return l.run, nil
					}
				} else {
					child.AddClosedFuture(listLifecycleContext{
						Lifecycle: child,
						channel:   l.closedChannel,
					})
				}
				continue next
			case <-l.closedChannel.Signal():
				if l.handleClosedSignal() {
					return nil, NewInterruptedError("生命周期列表", "启动")
				}
			case <-interrupter:
				l.shutdownChildren(l.children[:l.started], nil)
				return nil, NewInterruptedError("生命周期列表", "启动")
			}
		}
	}
	return l.run, nil
}

func (l *List) run(_ Lifecycle, interrupter chan struct{}) error {
	for {
		select {
		case <-l.closedChannel.Signal():
			if l.handleClosedSignal() {
				return nil
			}
		case <-interrupter:
			l.shutdownChildren(l.children, nil)
			return nil
		}
	}
}
