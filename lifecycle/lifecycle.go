package lifecycle

type Lifecycle interface {
	OnceLifecycle

	AddRunningFuture(future chan error) chan error

	Run() error

	Start() error

	Restart() error
}

type Runner interface {
	OnceRunner

	DoStart() error
}

type defaultRunner struct {
	defaultOnceRunner
	startFn func() error
}

func (r *defaultRunner) setStartFn(starter func() error) {
	r.startFn = starter
}

func (r *defaultRunner) setDefault() {
	if r.startFn == nil {
		r.startFn = func() error { return nil }
	}
	r.defaultOnceRunner.setDefault()
}

func (r *defaultRunner) DoStart() error {
	return r.startFn()
}

type DefaultRunner struct {
	abstractLifecycle
	defaultRunner
}

func (r *DefaultRunner) setSelf(self any) {
	if this, ok := self.(Lifecycle); ok {
		r.this = &struct {
			Lifecycle
			abstractLifecyclePrivateI
		}{
			Lifecycle:                 this,
			abstractLifecyclePrivateI: r,
		}
	} else {
		r.this = r
	}
}

func (r *DefaultRunner) setDefault() {
	if r.this == nil {
		r.this = r
	}
	r.defaultRunner.setDefault()
}

type abstractLifecyclePrivateI interface {
	abstractOnceLifecyclePrivateI

	broadcastRunning(err error)

	addRunningFuture(future chan error, makeIfNil bool) chan error

	startCheck() error

	doStart() error

	start() (err error, started bool)

	restartCheck() error

	doRestart() error
}

type abstractLifecycleI interface {
	Lifecycle

	abstractLifecyclePrivateI
}

type abstractLifecycle struct {
	abstractOnceLifecycle
	runningFutures []chan error
	startChecker   func() error
}

func (l *abstractLifecycle) self() abstractLifecycleI {
	return l.this.(abstractLifecycleI)
}

func (l *abstractLifecycle) SetStartChecker(checker func() error) {
	l.Lock()
	l.startChecker = checker
	l.Unlock()
}

func (l *abstractLifecycle) broadcastRunning(err error) {
	for _, c := range l.runningFutures {
		c <- err
	}
	l.runningFutures = l.runningFutures[:0]
}

func (l *abstractLifecycle) addRunningFuture(future chan error, makeIfNil bool) chan error {
	if future == nil && makeIfNil {
		future = make(chan error, 1)
	}
	if future != nil {
		l.runningFutures = append(l.runningFutures, future)
		return future
	}
	return nil
}

func (l *abstractLifecycle) AddRunningFuture(future chan error) chan error {
	l.Lock()
	defer l.Unlock()
	if l.Running() {
		return l.self().completed(future, true)
	}
	return l.self().addRunningFuture(future, true)
}

func (l *abstractLifecycle) startCheck() error {
	if l.Restarting() {
		return StateRestartingError(l.name)
	}
	if l.startChecker != nil {
		return l.startChecker()
	}
	return nil
}

func (l *abstractLifecycle) doStart() error {
	if err := l.runner.(Runner).DoStart(); err != nil {
		l.err = err
		l.self().broadcastRunning(err)
		return err
	}
	l.err = nil
	l.ToRunning()
	l.self().broadcastRunning(nil)
	return nil
}

func (l *abstractLifecycle) start() (err error, started bool) {
	l.Lock()
	defer l.Unlock()
	if l.Running() {
		return nil, true
	}
	if err = l.self().startCheck(); err != nil {
		l.err = err
		return
	}
	l.err = l.self().doStart()
	return l.err, false
}

func (l *abstractLifecycle) Run() error {
	if err, started := l.self().start(); err != nil {
		return err
	} else if started {
		return l.self().Wait()
	}
	return l.self().doRun()
}

func (l *abstractLifecycle) Start() error {
	if err, started := l.self().start(); err != nil {
		return err
	} else if started {
		return nil
	}
	go l.self().doRun()
	return nil
}

func (l *abstractLifecycle) restartCheck() error {
	if l.Restarting() {
		return StateRestartingError(l.name)
	}
	return nil
}

func (l *abstractLifecycle) doRestart() error {
	// must be locked

	l.ToRestarting()
	// start restart, Start, Close and Restart are not allowed
	defer l.ToRestarted()
	// when restart end, Start, Close and Restart will be allowed

	if l.Running() {
		future := make(chan error, 1)
		if err := l.self().doClose(future); err != nil {
			// close error, current state is running
			return err
		}
		// close request success, current state is closing
		l.Unlock()
		if err := <-future; err != nil {
			// exit error, current state is closed
			l.Lock()
			return err
		}
		// close success, current state is closed
		l.Lock()
	}

	if err := l.self().doStart(); err != nil {
		// start error, current state is closed
		return err
	}
	// start success, current state is running
	go l.self().doRun()
	return nil
}

func (l *abstractLifecycle) Restart() error {
	l.Lock()
	defer l.Unlock()
	if err := l.self().restartCheck(); err != nil {
		return err
	}
	return l.self().doRestart()
}

func New(name string, options ...Option) (*DefaultRunner, Lifecycle) {
	r := &DefaultRunner{
		abstractLifecycle: abstractLifecycle{
			abstractOnceLifecycle: abstractOnceLifecycle{
				name:         name,
				customStates: make(map[string]*CustomState),
				State:        StateClosed,
			},
		},
	}
	r.runner = r
	for _, option := range options {
		option.apply(r)
	}
	r.setDefault()
	return r, &r.abstractLifecycle
}
