package lifecycle

import (
	"fmt"
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/lock"
	"sync"
	"sync/atomic"
)

type (
	OnStartingFunc = func(Lifecycle)
	OnStartedFunc  = func(Lifecycle, error)
	OnCloseFunc    = func(Lifecycle, error)
	OnClosedFunc   = func(Lifecycle, error)
)

type Lifecycle interface {
	Run() error

	Background() error

	Start() error

	StartedWaiter() <-chan error

	Close(future Future[error]) error

	ClosedWaiter() <-chan error

	Shutdown() error

	AddStartedFuture(future Future[error]) Future[error]

	AddClosedFuture(future Future[error]) Future[error]

	OnStarting(OnStartingFunc) Lifecycle

	OnStarted(OnStartedFunc) Lifecycle

	OnClose(OnCloseFunc) Lifecycle

	OnClosed(OnClosedFunc) Lifecycle

	SetOnStarting(OnStartingFunc) Lifecycle

	SetOnStarted(OnStartedFunc) Lifecycle

	SetOnClose(OnCloseFunc) Lifecycle

	SetOnClosed(OnClosedFunc) Lifecycle

	Field(name string) any

	Fields() map[string]any

	SetField(name string, value any) Lifecycle

	SetDefaultField(name string, defaultValue any) (value any, exist bool)

	DeleteField(name string) Lifecycle

	RemoveField(name string) any

	RangeField(f func(name string, value any) bool) Lifecycle

	Error() error
}

func AddRunningFuture[FUTURE Future[error]](lifecycle Lifecycle, future FUTURE) FUTURE {
	lifecycle.AddStartedFuture(future)
	return future
}

func AddClosedFuture[FUTURE Future[error]](lifecycle Lifecycle, future FUTURE) FUTURE {
	lifecycle.AddClosedFuture(future)
	return future
}

type defaultLifecycleClass interface {
	// 在前台运行组件
	doRun() error

	// 在确定组件可以启动后，执行启动流程并在前台运行组件
	run() error

	// 检查组件是否可以启动，返回错误则不能启动
	startCheck() error

	// 检查组件是否已经启动，是否可以启动，如果返回的runnable为false，则不能启动，如果err
	// 不为nil，则是检查组件是否可以启动返回的错误
	preStart() (runnable bool, err error)

	// 在确定组件可以启动后，执行启动流程，如果返回的runnable为false，则不能运行，如果err
	// 不为nil，则是执行启动流程返回的错误
	doStart() (runnable bool, err error)

	// 检查组件是否可以启动并执行启动流程，如果返回的runnable为false，则不能运行，如果err
	// 不为nil，则是检查组件是否可以启动返回的错误，或是执行启动流程返回的错误
	start() (runnable bool, err error)

	// 检查组件是否可以关闭，返回错误则不能关闭
	closeCheck() error

	// 检查组件是否可以关闭并执行关闭流程，返回错误则说明关闭失败，组件的状态未改变
	doClose(future Future[error]) error

	doOnStarting()

	doOnStarted(err error)

	doOnClose(err error)

	doOnClosed(err error)

	lock.RLocker

	fmt.Stringer
}

type defaultLifecycle interface {
	Lifecycle

	defaultLifecycleClass
}

type DefaultLifecycle struct {
	this   defaultLifecycle
	runner Runner

	runningFutures SyncFutures[error]
	closedFutures  SyncFutures[error]

	err atomic.Pointer[error]

	onStarting     []OnStartingFunc
	onStartingLock sync.Mutex

	onStarted     []OnStartedFunc
	onStartedLock sync.Mutex

	onClose     []OnCloseFunc
	onCloseLock sync.Mutex

	onClosed     []OnClosedFunc
	onClosedLock sync.Mutex

	fields sync.Map

	State
	sync.RWMutex
}

func New(options ...Option) *DefaultLifecycle {
	lifecycle := new(DefaultLifecycle)
	lifecycle.init(lifecycle, options...)
	return lifecycle
}

func NewWithRunner(runner Runner, options ...Option) *DefaultLifecycle {
	return New(append(options, WithRunner(runner))...)
}

func NewWithRun(startFn StartFunc, runFn RunFunc, closeFn CloseFunc, options ...Option) *DefaultLifecycle {
	return NewWithRunner(FuncRunner(startFn, runFn, closeFn), options...)
}

func NewWithInterruptedRunner(runner InterruptedRunner, options ...Option) *DefaultLifecycle {
	return New(append(options, WithInterruptedRunner(runner))...)
}

func NewWithInterruptedRun(startFn InterruptedStartFunc, runFn InterruptedRunFunc, options ...Option) *DefaultLifecycle {
	return NewWithInterruptedRunner(FuncInterruptedRunner(startFn, runFn), options...)
}

func NewWithStarter(starter Starter, options ...Option) *DefaultLifecycle {
	return New(append(options, WithStarter(starter))...)
}

func NewWithStart(starterFn StarterFunc, options ...Option) *DefaultLifecycle {
	return NewWithStarter(FuncStarter(starterFn), options...)
}

func NewWithInterruptedStarter(starter InterruptedStarter, options ...Option) *DefaultLifecycle {
	return New(append(options, WithInterruptedStarter(starter))...)
}

func NewWithInterruptedStart(starterFn InterruptedStarterFunc, options ...Option) *DefaultLifecycle {
	return NewWithInterruptedStarter(FuncInterruptedStarter(starterFn), options...)
}

func (l *DefaultLifecycle) init(self defaultLifecycle, options ...Option) {
	l.this = self
	for _, option := range options {
		option.Apply(self)
	}
	if l.runner == nil {
		l.runner = FuncRunner(nil, nil, nil)
	}
}

func (l *DefaultLifecycle) self() defaultLifecycle {
	return l.this
}

func (l *DefaultLifecycle) setSelf(self any) {
	if this, is := self.(defaultLifecycle); is {
		l.this = this
	}
}

func (l *DefaultLifecycle) setRunner(runner Runner) {
	l.runner = runner
}

func (l *DefaultLifecycle) doRun() error {
	err := l.runner.DoRun(l.self())
	l.err.Store(&err)
	l.doOnClosed(err)
	lock.LockGet[Futures[error]](l, func() Futures[error] {
		l.ToClosed()
		return l.closedFutures.LoadAndReset()
	}).Complete(err)
	return err
}

func (l *DefaultLifecycle) run() error {
	if runnable, err := l.self().doStart(); !runnable {
		return err
	}
	return l.self().doRun()
}

func (l *DefaultLifecycle) Run() error {
	if runnable, err := l.self().start(); !runnable {
		return err
	}
	return l.self().doRun()
}

func (l *DefaultLifecycle) Background() error {
	if runnable, err := l.self().preStart(); !runnable {
		return err
	}
	go l.self().run()
	return nil
}

func (l *DefaultLifecycle) startCheck() error {
	return nil
}

func (l *DefaultLifecycle) preStart() (runnable bool, err error) {
	runnable, err = lock.LockGetDouble(l, func() (runnable bool, err error) {
		if !l.Closed() {
			return
		}
		err = l.self().startCheck()
		l.err.Store(&err)
		if err != nil {
			return
		}
		runnable = true
		l.ToStarting()
		return
	})
	if runnable {
		l.doOnStarting()
	}
	return
}

func (l *DefaultLifecycle) doStart() (runnable bool, err error) {
	err = l.runner.DoStart(l.self())
	if err != nil {
		l.err.Store(&err)
		l.doOnStarted(err)
		runningFutures, closedFutures := lock.LockGetDouble(l, func() (Futures[error], Futures[error]) {
			l.ToClosed()
			return l.runningFutures.LoadAndReset(), l.closedFutures.LoadAndReset()
		})
		runningFutures.Complete(err)
		closedFutures.Complete(err)
		return
	}
	l.doOnStarted(nil)
	lock.LockGet(l, func() Futures[error] {
		if !l.Closing() {
			l.ToRunning()
		}
		runnable = true
		return l.runningFutures.LoadAndReset()
	}).Complete(nil)
	return
}

func (l *DefaultLifecycle) start() (runnable bool, err error) {
	runnable, err = l.self().preStart()
	if !runnable {
		return
	}
	return l.self().doStart()
}

func (l *DefaultLifecycle) Start() error {
	if runnable, err := l.self().start(); !runnable {
		return err
	}
	go l.self().doRun()
	return nil
}

func (l *DefaultLifecycle) StartedWaiter() <-chan error {
	return AddRunningFuture(l.self(), make(ChanFuture[error], 1))
}

func (l *DefaultLifecycle) closeCheck() error {
	return nil
}

func (l *DefaultLifecycle) doClose(future Future[error]) error {
	if err := l.self().closeCheck(); err != nil {
		return err
	}
	if err := l.runner.DoClose(l.self()); err != nil {
		return err
	}
	l.ToClosing()
	l.closedFutures.Append(future)
	return nil
}

func (l *DefaultLifecycle) Close(future Future[error]) error {
	var closed bool
	var closing bool
	if err := lock.LockGet[error](l, func() error {
		if closed = l.Closed(); !closed {
			if closing = l.Closing(); !closing {
				return l.self().doClose(future)
			}
			l.closedFutures.Append(future)
			return nil
		}
		return l.self().Error()
	}); !closed {
		if !closing {
			l.doOnClose(err)
		}
		return err
	} else {
		if future != nil {
			future.Complete(err)
		}
		return nil
	}
}

func (l *DefaultLifecycle) ClosedWaiter() <-chan error {
	return AddClosedFuture(l.self(), make(ChanFuture[error], 1))
}

func (l *DefaultLifecycle) Shutdown() error {
	if err := l.self().Close(nil); err != nil {
		return err
	}
	<-l.self().ClosedWaiter()
	return nil
}

func (l *DefaultLifecycle) AddStartedFuture(future Future[error]) Future[error] {
	if future == nil {
		return nil
	}
	if lock.RLockGet(l, func() (completed bool) {
		if completed = l.Running() || l.Closing(); !completed {
			l.runningFutures.Append(future)
		}
		return
	}) {
		future.Complete(l.self().Error())
	}
	return future
}

func (l *DefaultLifecycle) AddClosedFuture(future Future[error]) Future[error] {
	if future == nil {
		return nil
	}
	if lock.RLockGet(l, func() (completed bool) {
		if completed = l.Closed(); !completed {
			l.closedFutures.Append(future)
		}
		return
	}) {
		future.Complete(l.self().Error())
	}
	return future
}

func (l *DefaultLifecycle) OnStarting(onStarting OnStartingFunc) Lifecycle {
	if onStarting != nil {
		lock.LockDo(&l.onStartingLock, func() { l.onStarting = append(l.onStarting, onStarting) })
	}
	return l.self()
}

func (l *DefaultLifecycle) OnStarted(onStarted OnStartedFunc) Lifecycle {
	if onStarted != nil {
		lock.LockDo(&l.onStartedLock, func() { l.onStarted = append(l.onStarted, onStarted) })
	}
	return l.self()
}

func (l *DefaultLifecycle) OnClose(onClose OnCloseFunc) Lifecycle {
	if onClose != nil {
		lock.LockDo(&l.onCloseLock, func() { l.onClose = append(l.onClose, onClose) })
	}
	return l.self()
}

func (l *DefaultLifecycle) OnClosed(onClosed OnClosedFunc) Lifecycle {
	if onClosed != nil {
		lock.LockDo(&l.onClosedLock, func() { l.onClosed = append(l.onClosed, onClosed) })
	}
	return l.self()
}

func (l *DefaultLifecycle) SetOnStarting(onStarting OnStartingFunc) Lifecycle {
	lock.LockDo(&l.onStartingLock, func() {
		if onStarting == nil {
			l.onStarting = l.onStarting[:0]
		} else {
			l.onStarting = append(l.onStarting, onStarting)
		}
	})
	return l.self()
}

func (l *DefaultLifecycle) SetOnStarted(onStarted OnStartedFunc) Lifecycle {
	lock.LockDo(&l.onStartedLock, func() {
		if onStarted == nil {
			l.onStarted = l.onStarted[:0]
		} else {
			l.onStarted = append(l.onStarted, onStarted)
		}
	})
	return l.self()
}

func (l *DefaultLifecycle) SetOnClose(onClose OnCloseFunc) Lifecycle {
	lock.LockDo(&l.onCloseLock, func() {
		if onClose == nil {
			l.onClose = l.onClose[:0]
		} else {
			l.onClose = append(l.onClose, onClose)
		}
	})
	return l.self()
}

func (l *DefaultLifecycle) SetOnClosed(onClosed OnClosedFunc) Lifecycle {
	lock.LockDo(&l.onClosedLock, func() {
		if onClosed == nil {
			l.onClosed = l.onClosed[:0]
		} else {
			l.onClosed = append(l.onClosed, onClosed)
		}
	})
	return l.self()
}

func (l *DefaultLifecycle) doOnStarting() {
	lock.LockDo(&l.onStartingLock, func() {
		for _, callback := range l.onStarting {
			callback(l.self())
		}
	})
}

func (l *DefaultLifecycle) doOnStarted(err error) {
	lock.LockDo(&l.onStartedLock, func() {
		for _, callback := range l.onStarted {
			callback(l.self(), err)
		}
	})
}

func (l *DefaultLifecycle) doOnClose(err error) {
	lock.LockDo(&l.onCloseLock, func() {
		for _, callback := range l.onClose {
			callback(l.self(), err)
		}
	})
}

func (l *DefaultLifecycle) doOnClosed(err error) {
	lock.LockDo(&l.onClosedLock, func() {
		for _, callback := range l.onClosed {
			callback(l.self(), err)
		}
	})
}

func (l *DefaultLifecycle) Field(name string) (value any) {
	value, _ = l.fields.Load(name)
	return
}

func (l *DefaultLifecycle) Fields() map[string]any {
	fields := make(map[string]any)
	l.self().RangeField(func(name string, value any) bool {
		fields[name] = value
		return true
	})
	return fields
}

func (l *DefaultLifecycle) SetField(name string, value any) Lifecycle {
	if value == nil {
		panic(errors.New("属性值不允许为空"))
	}
	l.fields.Store(name, value)
	return l.self()
}

func (l *DefaultLifecycle) SetDefaultField(name string, defaultValue any) (value any, exist bool) {
	return l.fields.LoadOrStore(name, defaultValue)
}

func (l *DefaultLifecycle) DeleteField(name string) Lifecycle {
	l.fields.Delete(name)
	return l.self()
}

func (l *DefaultLifecycle) RemoveField(name string) (value any) {
	value, _ = l.fields.LoadAndDelete(name)
	return
}

func (l *DefaultLifecycle) RangeField(f func(name string, value any) bool) Lifecycle {
	l.fields.Range(func(key, value any) bool { return f(key.(string), value) })
	return l.self()
}

func (l *DefaultLifecycle) Error() error {
	if err := l.err.Load(); err != nil {
		return *err
	}
	return nil
}

func (l *DefaultLifecycle) String() string {
	return fmt.Sprintf("生命周期组件(%p)[%s]", l, lock.RLockGet(l, func() State { return l.State }))
}
