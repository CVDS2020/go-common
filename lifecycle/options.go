package lifecycle

// An Option configures a Logger.
type Option interface {
	apply(runner OnceRunner)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(runner OnceRunner)

func (f optionFunc) apply(runner OnceRunner) {
	f(runner)
}

func Self(self any) Option {
	type selfSetter interface {
		setSelf(self any)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(selfSetter); ok {
			setter.setSelf(self)
		}
	})
}

func StartFn(startFn func() error) Option {
	type startFnSetter interface {
		setStartFn(checker func() error)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(startFnSetter); ok {
			setter.setStartFn(startFn)
		}
	})
}

func RunFn(runFn func() error) Option {
	type runFnSetter interface {
		setRunFn(checker func() error)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(runFnSetter); ok {
			setter.setRunFn(runFn)
		}
	})
}

func CloseFn(closeFn func() error) Option {
	type closeFnSetter interface {
		setCloseFn(checker func() error)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(closeFnSetter); ok {
			setter.setCloseFn(closeFn)
		}
	})
}

func Core(startFn func() error, runFn func() error, closeFn func() error) Option {
	return optionFunc(func(runner OnceRunner) {
		StartFn(startFn).apply(runner)
		RunFn(runFn).apply(runner)
		CloseFn(closeFn).apply(runner)
	})
}

func Context(runFn func(interrupter chan struct{}) error) Option {
	interrupter := make(chan struct{})
	return optionFunc(func(runner OnceRunner) {
		RunFn(func() error {
			return runFn(interrupter)
		}).apply(runner)
		CloseFn(func() error {
			interrupter <- struct{}{}
			return nil
		}).apply(runner)
	})
}

func StartChecker(checker func() error) Option {
	type startCheckerSetter interface {
		SetStartChecker(checker func() error)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(startCheckerSetter); ok {
			setter.SetStartChecker(checker)
		}
	})
}

func CloseChecker(checker func() error) Option {
	type closeCheckerSetter interface {
		SetCloseChecker(checker func() error)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(closeCheckerSetter); ok {
			setter.SetCloseChecker(checker)
		}
	})
}

func DestroyChecker(checker func() error) Option {
	type destroyCheckerSetter interface {
		SetDestroyChecker(checker func() error)
	}
	return optionFunc(func(runner OnceRunner) {
		if setter, ok := runner.(destroyCheckerSetter); ok {
			setter.SetDestroyChecker(checker)
		}
	})
}
