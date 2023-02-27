package lifecycle

import "gitee.com/sy_183/common/utils"

type InterruptedStarter interface {
	DoStart(lifecycle Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error)
}

type InterruptedStarterFunc = func(lifecycle Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error)

type interruptedStarterFunc InterruptedStarterFunc

func FuncInterruptedStarter(starterFn InterruptedStarterFunc) InterruptedStarter {
	if starterFn == nil {
		starterFn = func(lifecycle Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error) {
			return nil, nil
		}
	}
	return interruptedStarterFunc(starterFn)
}

func (f interruptedStarterFunc) DoStart(lifecycle Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error) {
	return f(lifecycle, interrupter)
}

type interruptedStarter struct {
	canInterrupted canInterrupted
	interrupter    chan struct{}
	starter        InterruptedStarter
	runFn          InterruptedRunFunc
}

func newInterruptedStarter(canInterrupted canInterrupted, starter InterruptedStarter) *interruptedStarter {
	return &interruptedStarter{
		canInterrupted: canInterrupted,
		interrupter:    make(chan struct{}, 1),
		starter:        starter,
	}
}

func (s *interruptedStarter) Runner() Runner {
	return FuncRunner(s.start, nil, s.close)
}

func (s *interruptedStarter) RunningRunner() Runner {
	return FuncRunner(nil, s.run, s.close)
}

func (s *interruptedStarter) start(lifecycle Lifecycle) (err error) {
	runFn, err := s.starter.DoStart(lifecycle, s.interrupter)
	if err != nil {
		utils.ChanTryPop(s.interrupter)
		return err
	}
	s.runFn = runFn
	s.canInterrupted.setRunner(s.RunningRunner())
	return nil
}

func (s *interruptedStarter) run(lifecycle Lifecycle) (err error) {
	defer utils.ChanTryPop(s.interrupter)
	defer s.canInterrupted.setRunner(s.Runner())
	if s.runFn != nil {
		return s.runFn(lifecycle, s.interrupter)
	}
	return nil
}

func (s *interruptedStarter) close(Lifecycle) error {
	s.canInterrupted.ToClosing()
	s.interrupter <- struct{}{}
	return nil
}
