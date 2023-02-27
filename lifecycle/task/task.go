package task

type Task interface {
	Do(interrupter chan struct{}) (interrupted bool)
}

type nopTask struct{}

func Nop() Task {
	return nopTask{}
}

func (t nopTask) Do(interrupter chan struct{}) (interrupted bool) { return false }

type interruptedTaskFunc func(interrupter chan struct{}) (interrupted bool)

func Interrupted(fn func(interrupter chan struct{}) (interrupted bool)) Task {
	if fn == nil {
		return Nop()
	}
	return interruptedTaskFunc(fn)
}

func (f interruptedTaskFunc) Do(interrupter chan struct{}) (interrupted bool) { return f(interrupter) }

type uninterruptedTaskFunc func()

func Func(fn func()) Task {
	if fn == nil {
		return Nop()
	}
	return uninterruptedTaskFunc(fn)
}

func (f uninterruptedTaskFunc) Do(interrupter chan struct{}) (interrupted bool) {
	f()
	return false
}
