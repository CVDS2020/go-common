package lifecycle

import (
	"strings"
)

const (
	StateClosed = State(iota)
	StateStarting
	StateRunning
	StateClosing
)

type StateInfo struct {
	State     State
	UpperName string
	LowerName string
}

var StateInfos = []StateInfo{
	{State: StateClosed, UpperName: "CLOSED"},
	{State: StateStarting, UpperName: "STARTING"},
	{State: StateRunning, UpperName: "RUNNING"},
	{State: StateClosing, UpperName: "CLOSING"},
}

func init() {
	for i := range StateInfos {
		StateInfos[i].LowerName = strings.ToLower(StateInfos[i].UpperName)
	}
}

type State int

func (s State) Closed() bool {
	return s == StateClosed
}

func (s State) Starting() bool {
	return s == StateStarting
}

func (s State) Running() bool {
	return s == StateRunning
}

func (s State) Closing() bool {
	return s == StateClosing
}

func (s State) String() string {
	for _, info := range StateInfos {
		if s == info.State {
			return info.UpperName
		}
	}
	return "UNKNOWN"
}

func (s State) check() {
	if s > StateClosing {
		panic(NewUnknownStateError("", s))
	}
}

func (s *State) ToClosed() {
	s.check()
	*s = StateClosed
}

func (s *State) ToStarting() {
	s.check()
	if s.Running() {
		panic(NewStateNotAllowSwitchError("", StateRunning.String(), StateStarting.String()))
	}
	*s = StateStarting
}

func (s *State) ToRunning() {
	s.check()
	*s = StateRunning
}

func (s *State) ToClosing() {
	s.check()
	if s.Closed() {
		panic(NewStateNotAllowSwitchError("", StateClosed.String(), StateClosing.String()))
	}
	*s = StateClosing
}
