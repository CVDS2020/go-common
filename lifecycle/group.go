package lifecycle

import (
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/lock"
	"sync"
	"sync/atomic"
)

const GroupFieldName = "$group"

type GroupLifecycleHolder struct {
	Lifecycle
	name  string
	group *Group

	removed atomic.Bool

	closeAllOnStartError atomic.Bool
	closeAllOnExit       atomic.Bool
	closeAllOnExitError  atomic.Bool
}

func (h *GroupLifecycleHolder) Name() string {
	return h.name
}

func (h *GroupLifecycleHolder) Group() *Group {
	return h.group
}

func (h *GroupLifecycleHolder) SetCloseAllOnStartError(enable bool) *GroupLifecycleHolder {
	h.closeAllOnStartError.Store(enable)
	return h
}

func (h *GroupLifecycleHolder) SetCloseAllOnExit(enable bool) *GroupLifecycleHolder {
	h.closeAllOnExit.Store(enable)
	return h
}

func (h *GroupLifecycleHolder) SetCloseAllOnExitError(enable bool) *GroupLifecycleHolder {
	h.closeAllOnExitError.Store(enable)
	return h
}

type (
	groupLifecycleContext = childLifecycleContext[*GroupLifecycleHolder]
	groupLifecycleChannel = childLifecycleChannel[*GroupLifecycleHolder]
)

type Group struct {
	Lifecycle
	lifecycle *DefaultLifecycle

	children     map[string]*GroupLifecycleHolder
	childrenLock sync.Mutex
	loaded       bool

	runningChannel *groupLifecycleChannel
	closedChannel  *groupLifecycleChannel
}

func NewGroup() *Group {
	g := &Group{
		children:       make(map[string]*GroupLifecycleHolder),
		runningChannel: newChildLifecycleChannel[*GroupLifecycleHolder](),
		closedChannel:  newChildLifecycleChannel[*GroupLifecycleHolder](),
	}
	g.lifecycle = NewWithInterruptedStart(g.start)
	g.Lifecycle = g.lifecycle
	return g
}

func (g *Group) Add(name string, lifecycle Lifecycle) (*GroupLifecycleHolder, error) {
	return lock.RLockGetDouble(g.lifecycle, func() (*GroupLifecycleHolder, error) {
		var loaded bool
		child, err := lock.LockGetDouble(&g.childrenLock, func() (*GroupLifecycleHolder, error) {
			if _, has := g.children[name]; has {
				return nil, errors.New("生命周期组件已经存在")
			}
			loaded = g.loaded
			child := &GroupLifecycleHolder{
				Lifecycle: lifecycle,
				name:      name,
				group:     g,
			}
			child.SetCloseAllOnStartError(true)
			child.SetCloseAllOnExit(true)
			child.SetCloseAllOnExitError(true)
			child.SetField(GroupFieldName, g)
			g.children[name] = child
			return child, nil
		})
		if err != nil {
			return nil, err
		}
		if loaded && !g.lifecycle.Closing() {
			child.AddStartedFuture(groupLifecycleContext{
				Lifecycle: child,
				channel:   g.runningChannel,
			})
			child.Background()
		}
		return child, nil
	})
}

func (g *Group) Remove(name string) *GroupLifecycleHolder {
	return lock.RLockGet(g.lifecycle, func() *GroupLifecycleHolder {
		child := lock.LockGet(&g.childrenLock, func() *GroupLifecycleHolder {
			child := g.children[name]
			if child == nil {
				return nil
			}
			child.removed.Store(true)
			child.RemoveField(GroupFieldName)
			delete(g.children, name)
			return child
		})
		if child == nil {
			return nil
		}
		if !g.lifecycle.Closed() {
			child.Close(nil)
		}
		return child
	})
}

func (g *Group) MustAdd(name string, lifecycle Lifecycle) *GroupLifecycleHolder {
	return assert.Must(g.Add(name, lifecycle))
}

func (g *Group) getChildren(setLoaded bool) []*GroupLifecycleHolder {
	return lock.LockGet(&g.childrenLock, func() (children []*GroupLifecycleHolder) {
		for _, child := range g.children {
			children = append(children, child)
		}
		if setLoaded {
			g.loaded = true
		}
		return
	})
}

func (g *Group) shutdownChildren(exclude map[*GroupLifecycleHolder]struct{}) {
	// 此函数必需在状态为Closing的情况下执行
	excluded := func(child *GroupLifecycleHolder) bool {
		if exclude != nil {
			_, has := exclude[child]
			return has
		}
		return false
	}
	children := g.getChildren(false)
	futures := make([]ChanFuture[error], 0)
	for _, child := range children {
		if !excluded(child) && !child.removed.Load() {
			future := make(ChanFuture[error], 1)
			child.Close(future)
			futures = append(futures, future)
		}
	}
	for _, future := range futures {
		<-future
	}
}

func (g *Group) handleRunningSignal(handleContext func(ctx groupLifecycleContext)) (closeAll bool) {
	contexts := g.runningChannel.Pop()
	var closedSet map[*GroupLifecycleHolder]struct{}
	for _, ctx := range contexts {
		child := ctx.Lifecycle
		if handleContext != nil {
			handleContext(ctx)
		}
		if child.removed.Load() {
			continue
		}
		if ctx.err != nil {
			if child.closeAllOnStartError.Load() {
				closeAll = true
			}
			if closedSet == nil {
				closedSet = make(map[*GroupLifecycleHolder]struct{})
			}
			closedSet[child] = struct{}{}
		} else {
			child.AddClosedFuture(groupLifecycleContext{
				Lifecycle: child,
				channel:   g.closedChannel,
			})
		}
	}
	if closeAll {
		lock.LockDo(g.lifecycle, func() { g.lifecycle.ToClosing() })
		g.shutdownChildren(closedSet)
		return true
	}
	return false
}

func (g *Group) handleClosedSignal(handleContext func(ctx groupLifecycleContext)) (closeAll bool) {
	contexts := g.closedChannel.Pop()
	var closedSet map[*GroupLifecycleHolder]struct{}
	for _, ctx := range contexts {
		child := ctx.Lifecycle
		if handleContext != nil {
			handleContext(ctx)
		}
		if child.removed.Load() {
			continue
		}
		if ctx.err != nil {
			if child.closeAllOnExit.Load() || child.closeAllOnExitError.Load() {
				closeAll = true
			}
		} else {
			if child.closeAllOnExit.Load() {
				closeAll = true
			}
		}
		if closedSet == nil {
			closedSet = make(map[*GroupLifecycleHolder]struct{})
		}
		closedSet[child] = struct{}{}
	}
	if closeAll {
		lock.LockDo(g.lifecycle, func() { g.lifecycle.ToClosing() })
		g.shutdownChildren(closedSet)
		return true
	}
	return false
}

func (g *Group) start(_ Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error) {
	defer func() {
		if err != nil {
			g.reset()
		}
	}()
	children := g.getChildren(true)
	for _, child := range children {
		child.AddStartedFuture(groupLifecycleContext{
			Lifecycle: child,
			channel:   g.runningChannel,
		})
		child.Background()
	}
	childSet := make(map[*GroupLifecycleHolder]struct{})
	for _, child := range children {
		childSet[child] = struct{}{}
	}
	startedSet := make(map[*GroupLifecycleHolder]struct{})
	childStarted := func(child *GroupLifecycleHolder) {
		if _, in := childSet[child]; in {
			startedSet[child] = struct{}{}
		}
	}
	allStarted := func() bool { return len(startedSet) == len(childSet) }

	for {
		select {
		case <-g.runningChannel.Signal():
			if g.handleRunningSignal(func(ctx groupLifecycleContext) {
				if child := ctx.Lifecycle; ctx.err == nil || !child.closeAllOnStartError.Load() || child.removed.Load() {
					childStarted(child)
				}
			}) {
				return nil, NewInterruptedError("生命周期组", "启动")
			}
			if allStarted() {
				return g.run, nil
			}
		case <-g.closedChannel.Signal():
			if g.handleClosedSignal(nil) {
				return nil, NewInterruptedError("生命周期组", "启动")
			}
		case <-interrupter:
			g.shutdownChildren(nil)
			return nil, NewInterruptedError("生命周期组", "启动")
		}
	}
}

func (g *Group) run(_ Lifecycle, interrupter chan struct{}) error {
	defer g.reset()
	for {
		select {
		case <-g.runningChannel.Signal():
			if g.handleRunningSignal(nil) {
				return nil
			}
		case <-g.closedChannel.Signal():
			if g.handleClosedSignal(nil) {
				return nil
			}
		case <-interrupter:
			g.shutdownChildren(nil)
			return nil
		}
	}
}

func (g *Group) reset() {
	lock.LockDo(&g.childrenLock, func() { g.loaded = false })
	g.runningChannel = newChildLifecycleChannel[*GroupLifecycleHolder]()
	g.closedChannel = newChildLifecycleChannel[*GroupLifecycleHolder]()
}
