package lifecycle

import (
	"fmt"
)

type UnknownStateError struct {
	Target string
	State  State
}

func NewUnknownStateError(target string, state State) *UnknownStateError {
	return &UnknownStateError{Target: target, State: state}
}

func (e *UnknownStateError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return fmt.Sprintf("未知的生命周期组件状态(%s)", e.State)
	} else {
		return fmt.Sprintf("未知的%s状态(%s)", e.Target, e.State)
	}
}

type StateDependencyError struct {
	Target string
	State  string
	Depend string
}

func NewStateDependencyError(target string, state string, depend string) *StateDependencyError {
	return &StateDependencyError{Target: target, State: state, Depend: depend}
}

func (e *StateDependencyError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return fmt.Sprintf("生命周期组件状态(%s)依赖(%s)", e.State, e.Depend)
	} else {
		return fmt.Sprintf("%s状态(%s)依赖(%s)", e.Target, e.State, e.Depend)
	}
}

type StateConflictError struct {
	Target   string
	State    string
	Conflict string
}

func NewStateConflictError(target string, state string, conflict string) *StateConflictError {
	return &StateConflictError{Target: target, State: state, Conflict: conflict}
}

func (e *StateConflictError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return fmt.Sprintf("生命周期组件状态(%s)与(%s)冲突", e.State, e.Conflict)
	} else {
		return fmt.Sprintf("%s状态(%s)与(%s)冲突", e.Target, e.State, e.Conflict)
	}
}

type StateNotAllowSwitchError struct {
	Target string
	Form   string
	To     string
}

func NewStateNotAllowSwitchError(target string, from string, to string) *StateNotAllowSwitchError {
	return &StateNotAllowSwitchError{Target: target, Form: from, To: to}
}

func (e *StateNotAllowSwitchError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return fmt.Sprintf("生命周期组件不可以从(%s)转变到(%s)", e.Form, e.To)
	} else {
		return fmt.Sprintf("%s不可以从(%s)转变到(%s)", e.Target, e.Form, e.To)
	}
}

type StateMustSwitchError struct {
	Target string
	Form   string
	To     string
}

func NewStateMustSwitchError(target string, from string, to string) *StateMustSwitchError {
	return &StateMustSwitchError{Target: target, Form: from, To: to}
}

func (e *StateMustSwitchError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return fmt.Sprintf("生命周期组件必须从(%s)转变到(%s)", e.Form, e.To)
	} else {
		return fmt.Sprintf("%s必须从(%s)转变到(%s)", e.Target, e.Form, e.To)
	}
}

type InterruptedError struct {
	Target string
	Action string
}

func NewInterruptedError(target string, action string) *InterruptedError {
	return &InterruptedError{Target: target, Action: action}
}

func (e *InterruptedError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		if e.Action == "" {
			return "生命周期组件已被中断"
		} else {
			return fmt.Sprintf("生命周期组件%s已被中断", e.Action)
		}
	} else {
		return fmt.Sprintf("%s%s已被中断", e.Target, e.Action)
	}
}

type StateStartingError struct {
	Target string
}

func NewStateStartingError(target string) *StateStartingError {
	return &StateStartingError{Target: target}
}

func (e *StateStartingError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return "生命周期组件正在启动"
	} else {
		return fmt.Sprintf("%s正在启动", e.Target)
	}
}

type StateRunningError struct {
	Target string
}

func NewStateRunningError(target string) *StateRunningError {
	return &StateRunningError{Target: target}
}

func (e *StateRunningError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return "生命周期组件正在运行"
	} else {
		return fmt.Sprintf("%s正在运行", e.Target)
	}
}

type StateNotRunningError struct {
	Target string
}

func NewStateNotRunningError(target string) *StateNotRunningError {
	return &StateNotRunningError{Target: target}
}

func (e *StateNotRunningError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return "生命周期组件未运行"
	} else {
		return fmt.Sprintf("%s未运行", e.Target)
	}
}

type StateClosingError struct {
	Target string
}

func NewStateClosingError(target string) *StateClosingError {
	return &StateClosingError{Target: target}
}

func (e *StateClosingError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return "生命周期组件正在关闭"
	} else {
		return fmt.Sprintf("%s正在关闭", e.Target)
	}
}

type StateClosedError struct {
	Target string
}

func NewStateClosedError(target string) *StateClosedError {
	return &StateClosedError{Target: target}
}

func (e *StateClosedError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return "生命周期组件已经关闭"
	} else {
		return fmt.Sprintf("%s已经关闭", e.Target)
	}
}

type StateNotClosedError struct {
	Target string
}

func NewStateNotClosedError(target string) *StateNotClosedError {
	return &StateNotClosedError{Target: target}
}

func (e *StateNotClosedError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Target == "" {
		return "生命周期组件未关闭"
	} else {
		return fmt.Sprintf("%s未关闭", e.Target)
	}
}
