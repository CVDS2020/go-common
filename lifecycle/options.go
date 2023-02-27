package lifecycle

type Option interface {
	Apply(lifecycle Lifecycle)
}

type optionFunc func(lifecycle Lifecycle)

func (f optionFunc) Apply(lifecycle Lifecycle) {
	f(lifecycle)
}

func WithSelf(self any) Option {
	return optionFunc(func(lifecycle Lifecycle) {
		if setter, is := lifecycle.(interface{ setSelf(self any) }); is {
			setter.setSelf(self)
		}
	})
}

func WithRunner(runner Runner) Option {
	return optionFunc(func(lifecycle Lifecycle) {
		if setter, is := lifecycle.(interface{ setRunner(runner Runner) }); is {
			setter.setRunner(runner)
		}
	})
}

func WithInterruptedRunner(runner InterruptedRunner) Option {
	return optionFunc(func(lifecycle Lifecycle) {
		if canInterrupted, is := lifecycle.(canInterrupted); is {
			canInterrupted.setRunner(newInterrupterRunner(canInterrupted, runner))
		}
	})
}

func WithStarter(starter Starter) Option {
	return optionFunc(func(lifecycle Lifecycle) {
		if setter, is := lifecycle.(interface{ setRunner(runner Runner) }); is {
			setter.setRunner(newStarterRunner(setter, starter).Runner())
		}
	})
}

func WithInterruptedStarter(starter InterruptedStarter) Option {
	return optionFunc(func(lifecycle Lifecycle) {
		if canInterrupted, is := lifecycle.(canInterrupted); is {
			canInterrupted.setRunner(newInterruptedStarter(canInterrupted, starter).Runner())
		}
	})
}
