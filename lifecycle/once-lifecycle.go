package lifecycle

import (
	"gitee.com/sy_183/common/utils"
	"sync"
)

type OnceLifecycle interface {
	AddClosedFuture(future chan error) chan error

	SetClosedFuture(future chan error) chan error

	AddCustomStateFuture(name string, future chan bool) chan bool

	Close(future chan error) error

	Wait() error

	CloseWait() (closeErr error, exitErr error)

	Error() error
}

type OnceRunner interface {
	DoRun() error

	DoClose() error
}

type defaultOnceRunner struct {
	runFn   func() error
	closeFn func() error
}

func (r *defaultOnceRunner) setRunFn(runFn func() error) {
	r.runFn = runFn
}

func (r *defaultOnceRunner) setCloseFn(closeFn func() error) {
	r.closeFn = closeFn
}

func (r *defaultOnceRunner) setDefault() {
	if r.runFn == nil && r.closeFn == nil {
		closeChan := make(chan struct{}, 1)
		r.runFn = func() error {
			<-closeChan
			return nil
		}
		r.closeFn = func() error {
			closeChan <- struct{}{}
			return nil
		}
	} else if r.runFn == nil {
		r.runFn = func() error { return nil }
	} else if r.closeFn == nil {
		r.runFn = func() error { return nil }
	}
}

func (r *defaultOnceRunner) DoRun() error {
	return r.runFn()
}

func (r *defaultOnceRunner) DoClose() error {
	return r.closeFn()
}

type DefaultOnceRunner struct {
	abstractOnceLifecycle
	defaultOnceRunner
}

func (r *DefaultOnceRunner) setSelf(self any) {
	if this, ok := self.(OnceLifecycle); ok {
		r.this = &struct {
			OnceLifecycle
			abstractOnceLifecyclePrivateI
		}{
			OnceLifecycle:                 this,
			abstractOnceLifecyclePrivateI: r,
		}
	} else {
		r.this = r
	}
}

func (r *DefaultOnceRunner) setDefault() {
	if r.this == nil {
		r.this = r
	}
	r.defaultOnceRunner.setDefault()
}

type abstractOnceLifecyclePrivateI interface {
	Name() string

	completed(future chan error, makeIfNil bool) chan error

	notifyClosed(err error)

	addClosedFuture(future chan error, makeIfNil bool) chan error

	setClosedFuture(future chan error) chan error

	AddCustomState(name string)

	triggerCustomStates()

	TriggerCustomState(name string)

	doRun() error

	closeCheck() error

	doClose(future chan error) error
}

type abstractOnceLifecycleI interface {
	OnceLifecycle

	abstractOnceLifecyclePrivateI
}

type CustomState struct {
	Name    string
	state   bool
	Futures []chan bool
}

type abstractOnceLifecycle struct {
	this   abstractOnceLifecycleI
	runner OnceRunner
	name   string

	closedFutures []chan error
	closeChecker  func() error

	customStates map[string]*CustomState

	err error
	State
	sync.Mutex
}

func (l *abstractOnceLifecycle) self() abstractOnceLifecycleI {
	return l.this
}

func (l *abstractOnceLifecycle) Name() string {
	return l.name
}

func (l *abstractOnceLifecycle) ClosingLock() bool {
	return l.Closing()
}

func (l *abstractOnceLifecycle) SetCloseChecker(checker func() error) {
	l.Lock()
	l.closeChecker = checker
	l.Unlock()
}

func (l *abstractOnceLifecycle) completed(future chan error, makeIfNil bool) chan error {
	if future == nil && makeIfNil {
		future = make(chan error, 1)
	}
	if future != nil {
		select {
		case future <- l.err:
		default:
			go func(err error) {
				future <- err
			}(l.err)
		}
	}
	return future
}

func (l *abstractOnceLifecycle) notifyClosed(err error) {
	for _, c := range l.closedFutures {
		utils.ChanAsyncPush(c, err)
	}
	l.closedFutures = l.closedFutures[:0]
}

func (l *abstractOnceLifecycle) addClosedFuture(future chan error, makeIfNil bool) chan error {
	if future == nil && makeIfNil {
		future = make(chan error, 1)
	}
	if future != nil {
		l.closedFutures = append(l.closedFutures, future)
		return future
	}
	return nil
}

func (l *abstractOnceLifecycle) setClosedFuture(future chan error) chan error {
	if future == nil {
		l.closedFutures = l.closedFutures[:0]
		return nil
	}
	l.closedFutures = append(l.closedFutures[:0], future)
	return future
}

func (l *abstractOnceLifecycle) AddCustomState(name string) {
	l.Lock()
	defer l.Unlock()
	l.customStates[name] = &CustomState{Name: name}
}

func (l *abstractOnceLifecycle) triggerCustomStates() {
	for _, state := range l.customStates {
		if !state.state {
			for _, future := range state.Futures {
				utils.ChanAsyncPush(future, false)
			}
			state.Futures = state.Futures[:0]
		}
	}
}

func (l *abstractOnceLifecycle) TriggerCustomState(name string) {
	l.Lock()
	defer l.Unlock()
	if l.State.Closed() {
		return
	}
	if customState := l.customStates[name]; customState != nil {
		if !customState.state {
			customState.state = true
			for _, future := range customState.Futures {
				utils.ChanAsyncPush(future, true)
			}
			customState.Futures = customState.Futures[:0]
		}
	}
}

func (l *abstractOnceLifecycle) doRun() error {
	err := l.runner.DoRun()
	l.Lock()
	l.err = err
	l.ToClosed()
	l.self().notifyClosed(err)
	l.self().triggerCustomStates()
	l.Unlock()
	return err
}

func (l *abstractOnceLifecycle) closeCheck() error {
	if l.Restarting() {
		return StateRestartingError(l.name)
	}
	if l.closeChecker != nil {
		return l.closeChecker()
	}
	return nil
}

func (l *abstractOnceLifecycle) doClose(future chan error) error {
	if l.Closing() {
		l.self().addClosedFuture(future, false)
		return nil
	}
	if err := l.runner.DoClose(); err != nil {
		return err
	}
	l.ToClosing()
	l.self().addClosedFuture(future, false)
	return nil
}

func (l *abstractOnceLifecycle) AddClosedFuture(future chan error) chan error {
	l.Lock()
	defer l.Unlock()
	if l.Closed() {
		return l.self().completed(future, true)
	}
	return l.self().addClosedFuture(future, true)
}

func (l *abstractOnceLifecycle) SetClosedFuture(future chan error) chan error {
	l.Lock()
	defer l.Unlock()
	if l.Closed() {
		return l.self().completed(future, false)
	}
	return l.self().setClosedFuture(future)
}

func (l *abstractOnceLifecycle) AddCustomStateFuture(name string, future chan bool) chan bool {
	l.Lock()
	defer l.Unlock()
	if customState := l.customStates[name]; customState != nil {
		if future == nil {
			future = make(chan bool, 1)
		}
		if customState.state || l.State.Closed() {
			utils.ChanAsyncPush(future, customState.state)
		} else {
			customState.Futures = append(customState.Futures, future)
		}
		return future
	}
	return nil
}

func (l *abstractOnceLifecycle) Close(future chan error) error {
	l.Lock()
	defer l.Unlock()
	if l.Closed() {
		l.self().completed(future, false)
		return nil
	}
	if err := l.self().closeCheck(); err != nil {
		return err
	}
	return l.self().doClose(future)
}

func (l *abstractOnceLifecycle) Wait() error {
	return <-l.self().AddClosedFuture(nil)
}

func (l *abstractOnceLifecycle) CloseWait() (closeErr error, exitErr error) {
	future := make(chan error, 1)
	if err := l.self().Close(future); err != nil {
		return err, nil
	}
	return nil, <-future
}

func (l *abstractOnceLifecycle) Error() error {
	l.Lock()
	defer l.Unlock()
	return l.err
}

func NewOnce(name string, options ...Option) (*DefaultOnceRunner, OnceLifecycle) {
	r := &DefaultOnceRunner{
		abstractOnceLifecycle: abstractOnceLifecycle{
			name:         name,
			customStates: make(map[string]*CustomState),
			State:        StateRunning,
		},
	}
	r.runner = r
	for _, option := range options {
		option.apply(r)
	}
	r.setDefault()
	go r.self().doRun()
	return r, &r.abstractOnceLifecycle
}
