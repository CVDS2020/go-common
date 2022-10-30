package lifecycle

type DestructibleLifecycle interface {
	Lifecycle

	Destroy(future chan error) error

	DestroyWait() (closeErr error, exitErr error)
}

type DefaultDestructibleRunner struct {
	abstractDestructibleLifecycle
	defaultRunner
}

func (r *DefaultDestructibleRunner) setDefault() {
	if r.this == nil {
		r.this = r
	}
	r.defaultRunner.setDefault()
}

func (r *DefaultDestructibleRunner) setSelf(self any) {
	if this, ok := self.(DestructibleLifecycle); ok {
		r.this = &struct {
			DestructibleLifecycle
			abstractDestructibleLifecyclePrivateI
		}{DestructibleLifecycle: this, abstractDestructibleLifecyclePrivateI: r}
	} else {
		r.this = r
	}
}

type abstractDestructibleLifecyclePrivateI interface {
	abstractLifecyclePrivateI

	destroyCheck() error

	doDestroy(future chan error) error
}

type abstractDestructibleLifecycleI interface {
	DestructibleLifecycle

	abstractDestructibleLifecyclePrivateI
}

type abstractDestructibleLifecycle struct {
	abstractLifecycle
	destroyChecker func() error
}

func (l *abstractDestructibleLifecycle) self() abstractDestructibleLifecycleI {
	return l.this.(abstractDestructibleLifecycleI)
}

func (l *abstractDestructibleLifecycle) SetDestroyChecker(checker func() error) {
	l.Lock()
	l.destroyChecker = checker
	l.Unlock()
}

func (l *abstractDestructibleLifecycle) startCheck() error {
	if l.Destroyed() {
		return StateDestroyedError(l.name)
	}
	return l.abstractLifecycle.startCheck()
}

func (l *abstractDestructibleLifecycle) doRun() error {
	err := l.runner.DoRun()
	l.Lock()
	if l.Destroying() {
		l.ToDestroyed()
	} else {
		l.ToClosed()
	}
	l.self().notifyClosed(err)
	l.Unlock()
	return err
}

func (l *abstractDestructibleLifecycle) restartCheck() error {
	if l.Destroying() {
		return StateDestroyingError(l.name)
	} else if l.Destroyed() {
		return StateDestroyedError(l.name)
	}
	return l.abstractLifecycle.restartCheck()
}

func (l *abstractDestructibleLifecycle) destroyCheck() error {
	if err := l.self().closeCheck(); err != nil {
		return err
	}
	if l.destroyChecker != nil {
		return l.destroyChecker()
	}
	return nil
}

func (l *abstractDestructibleLifecycle) doDestroy(future chan error) error {
	if l.Closing() {
		l.ToDestroying()
		l.self().addClosedFuture(future, false)
		return nil
	}
	if err := l.runner.DoClose(); err != nil {
		return err
	}
	l.ToDestroying()
	l.self().addClosedFuture(future, false)
	return nil
}

func (l *abstractDestructibleLifecycle) Destroy(future chan error) error {
	l.Lock()
	defer l.Unlock()
	if err := l.self().destroyCheck(); err != nil {
		return err
	}
	return l.self().doDestroy(future)
}

func (l *abstractDestructibleLifecycle) DestroyWait() (closeErr error, exitErr error) {
	future := make(chan error, 1)
	if err := l.self().Destroy(future); err != nil {
		return err, nil
	}
	return nil, <-future
}

func NewDestructible(name string, options ...Option) (*DefaultDestructibleRunner, DestructibleLifecycle) {
	r := &DefaultDestructibleRunner{
		abstractDestructibleLifecycle: abstractDestructibleLifecycle{
			abstractLifecycle: abstractLifecycle{
				abstractOnceLifecycle: abstractOnceLifecycle{
					name:         name,
					customStates: make(map[string]*CustomState),
					State:        StateClosed,
				},
			},
		},
	}
	r.runner = r
	for _, option := range options {
		option.apply(r)
	}
	r.setDefault()
	return r, &r.abstractDestructibleLifecycle
}
