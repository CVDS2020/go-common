package lifecycle

import "fmt"

const (
	// StateClosed indicates instance has been closed, conflicts with
	// StateRunning and StateClosing, cannot change to StateClosing
	StateClosed = 1 << iota

	// StateRunning indicates instance is running, conflicts with
	// StateClosed
	StateRunning

	// StateClosing indicates instance is running, conflicts with
	// StateClosed, must coexist with StateRunning
	StateClosing

	// StateRestarting indicates instance is restart, can coexist with
	// any other state
	StateRestarting

	StateDestroying

	StateDestroyed
)

const stateMask = StateClosed | StateRunning | StateClosing | StateRestarting | StateDestroying | StateDestroyed

type State int

func (s State) Is(state State) bool {
	return s&state != 0
}

func (s State) Closed() bool {
	return s.Is(StateClosed)
}

func (s State) Running() bool {
	return s.Is(StateRunning)
}

func (s State) Closing() bool {
	return s.Is(StateClosing)
}

func (s State) Restarting() bool {
	return s.Is(StateRestarting)
}

func (s State) Destroying() bool {
	return s.Is(StateDestroying)
}

func (s State) Destroyed() bool {
	return s.Is(StateDestroyed)
}

func (s State) check() {
	if s&stateMask == 0 {
		panic(fmt.Errorf("unknown lifecycle state"))
	}
	if s.Destroying() && !s.Closing() {
		// destroying must with closing
		panic(fmt.Errorf("lifecycle state destroying must with closing"))
	}
	if s.Destroyed() && !s.Closed() {
		panic(fmt.Errorf("lifecycle state destroyed must with closed"))
	}
	if s.Closing() && !s.Running() {
		// closing must with running
		panic(fmt.Errorf("lifecycle state closing must with running"))
	}
	if s.Closed() == s.Running() {
		// closed conflicts with running
		panic(fmt.Errorf("lifecycle state closed conflicts with running"))
	}
	if s.Restarting() && s.Destroying() {
		// restarting conflicts with destroying
		panic(fmt.Errorf("lifecycle state restarting conflicts with destroying"))
	}
	if s.Restarting() && s.Destroyed() {
		// restarting conflicts with destroyed
		panic(fmt.Errorf("lifecycle state restarting conflicts with destroyed"))
	}
}

func (s State) notAllowSwitch(old, new string) {
	panic(fmt.Errorf("lifecycle state cannot be changed from %s to %s", old, new))
}

func (s State) mustSwitch(old, new string) {
	panic(fmt.Errorf("lifecycle state must be changed from %s to %s", old, new))
}

func (s *State) Set(state State) {
	switch state {
	case StateClosed:
		s.ToClosed()
	case StateRunning:
		s.ToRunning()
	case StateClosing:
		s.ToClosing()
	case StateRestarting:
		s.ToRestarting()
	case StateDestroying:
		s.ToDestroying()
	case StateDestroyed:
		s.ToDestroyed()
	}
}

func (s *State) ToClosed() {
	s.check()
	if s.Destroyed() {
		s.notAllowSwitch("destroyed", "closed")
	}
	*s &= ^(StateRunning | StateClosing)
	*s |= StateClosed
}

func (s *State) ToRunning() {
	s.check()
	if s.Destroyed() {
		s.notAllowSwitch("destroyed", "running")
	}
	*s &= ^StateClosed
	*s |= StateRunning
}

func (s *State) ToClosing() {
	s.check()
	if s.Destroyed() {
		s.notAllowSwitch("destroyed", "closing")
	} else if s.Closed() {
		s.notAllowSwitch("closed", "closing")
	}
	*s |= StateClosing
}

func (s *State) ToRestarting() {
	s.check()
	if s.Destroyed() {
		s.notAllowSwitch("destroyed", "restarting")
	}
	*s |= StateRestarting
}

func (s *State) ToRestarted() {
	s.check()
	if s.Destroyed() {
		s.notAllowSwitch("destroyed", "restarted")
	} else if s.Closing() {
		// closing, restarted is impossible
		s.notAllowSwitch("closing", "restarted")
	} else if !s.Restarting() {
		// if switch state to 'restarted', must be restarting
		s.mustSwitch("restarting", "restarted")
	}
	*s &= ^StateRestarting
}

func (s *State) ToDestroying() {
	s.check()
	if s.Destroyed() {
		// destroyed, destroying is impossible
		s.notAllowSwitch("destroyed", "destroying")
	} else if s.Closed() {
		// closed, destroying is impossible
		s.notAllowSwitch("closed", "destroying")
	}
	// destroying must with closing
	if !s.Closing() {
		s.ToClosing()
	}
	*s |= StateDestroying
}

func (s *State) ToDestroyed() {
	s.check()
	if !s.Closed() {
		s.ToClosed()
	}
	*s &= ^StateDestroying
	*s |= StateDestroyed
}
