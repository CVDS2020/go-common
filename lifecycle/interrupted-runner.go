package lifecycle

import "gitee.com/sy_183/common/utils"

type InterruptedRunner interface {
	DoStart(lifecycle Lifecycle, interrupter chan struct{}) error

	DoRun(lifecycle Lifecycle, interrupter chan struct{}) error
}

type (
	InterruptedStartFunc = func(lifecycle Lifecycle, interrupter chan struct{}) error
	InterruptedRunFunc   = func(lifecycle Lifecycle, interrupter chan struct{}) error
)

func InterrupterHoldRun(_ Lifecycle, interrupter chan struct{}) error {
	<-interrupter
	return nil
}

type interruptedRunnerFunc struct {
	startFn InterruptedStartFunc
	runFn   InterruptedRunFunc
}

func FuncInterruptedRunner(startFn InterruptedStartFunc, runFn InterruptedRunFunc) InterruptedRunner {
	if startFn == nil {
		startFn = func(lifecycle Lifecycle, interrupter chan struct{}) error { return nil }
	}
	if runFn == nil {
		runFn = func(lifecycle Lifecycle, interrupter chan struct{}) error { return nil }
	}
	return interruptedRunnerFunc{startFn: startFn, runFn: runFn}
}

func (f interruptedRunnerFunc) DoStart(lifecycle Lifecycle, interrupter chan struct{}) error {
	return f.startFn(lifecycle, interrupter)
}

func (f interruptedRunnerFunc) DoRun(lifecycle Lifecycle, interrupter chan struct{}) error {
	return f.runFn(lifecycle, interrupter)
}

type canInterrupted interface {
	setRunner(runner Runner)
	ToClosing()
}

type interruptedRunner struct {
	canInterrupted canInterrupted
	interrupter    chan struct{}
	runner         InterruptedRunner
}

func newInterrupterRunner(canInterrupted canInterrupted, runner InterruptedRunner) *interruptedRunner {
	return &interruptedRunner{
		canInterrupted: canInterrupted,
		interrupter:    make(chan struct{}, 1),
		runner:         runner,
	}
}

func (r *interruptedRunner) DoStart(lifecycle Lifecycle) error {
	err := r.runner.DoStart(lifecycle, r.interrupter)
	if err != nil {
		utils.ChanTryPop(r.interrupter)
	}
	return err
}

func (r *interruptedRunner) DoRun(lifecycle Lifecycle) error {
	defer utils.ChanTryPop(r.interrupter)
	return r.runner.DoRun(lifecycle, r.interrupter)
}

func (r *interruptedRunner) DoClose(Lifecycle) error {
	r.canInterrupted.ToClosing()
	r.interrupter <- struct{}{}
	return nil
}
