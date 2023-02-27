package lifecycle

type Starter interface {
	DoStart(Lifecycle) (runFn RunFunc, closeFn CloseFunc, err error)
}

type StarterFunc func(Lifecycle) (runFn RunFunc, closeFn CloseFunc, err error)

type starterFunc StarterFunc

func FuncStarter(starterFn StarterFunc) Starter {
	if starterFn == nil {
		starterFn = func(Lifecycle) (runFn RunFunc, closeFn CloseFunc, err error) {
			return nil, nil, nil
		}
	}
	return starterFunc(starterFn)
}

func (f starterFunc) DoStart(lifecycle Lifecycle) (runFn RunFunc, closeFn CloseFunc, err error) {
	return f(lifecycle)
}

type starterRunner struct {
	runnerSetter interface{ setRunner(runner Runner) }
	starter      Starter
	runFn        RunFunc
	closeFn      CloseFunc
}

func newStarterRunner(runnerSetter interface{ setRunner(runner Runner) }, starter Starter) *starterRunner {
	return &starterRunner{
		runnerSetter: runnerSetter,
		starter:      starter,
	}
}

func (r *starterRunner) Runner() Runner {
	return FuncRunner(r.start, nil, nil)
}

func (r *starterRunner) RunningRunner() Runner {
	return FuncRunner(nil, r.run, r.closeFn)
}

func (r *starterRunner) start(lifecycle Lifecycle) error {
	runFn, closeFn, err := r.starter.DoStart(lifecycle)
	if err != nil {
		return err
	}
	r.runFn, r.closeFn = runFn, closeFn
	r.runnerSetter.setRunner(r.RunningRunner())
	return nil
}

func (r *starterRunner) run(lifecycle Lifecycle) (err error) {
	if r.runFn != nil {
		err = r.runFn(lifecycle)
	}
	r.runnerSetter.setRunner(r.Runner())
	return err
}
