package lifecycle

type Runner interface {
	DoStart(Lifecycle) error

	DoRun(Lifecycle) error

	DoClose(Lifecycle) error
}

type (
	StartFunc = func(Lifecycle) error
	RunFunc   = func(Lifecycle) error
	CloseFunc = func(Lifecycle) error
)

type runnerFunc struct {
	startFn StartFunc
	runFn   RunFunc
	closeFn CloseFunc
}

func FuncRunner(startFn, runFn, closeFn func(Lifecycle) error) Runner {
	if startFn == nil {
		startFn = func(Lifecycle) error { return nil }
	}
	if runFn == nil {
		runFn = func(Lifecycle) error { return nil }
	}
	if closeFn == nil {
		closeFn = func(Lifecycle) error { return nil }
	}
	return runnerFunc{startFn: startFn, runFn: runFn, closeFn: closeFn}
}

func (f runnerFunc) DoStart(lifecycle Lifecycle) error {
	return f.startFn(lifecycle)
}

func (f runnerFunc) DoRun(lifecycle Lifecycle) error {
	return f.runFn(lifecycle)
}

func (f runnerFunc) DoClose(lifecycle Lifecycle) error {
	return f.closeFn(lifecycle)
}
