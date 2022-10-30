package lifecycle

import (
	"errors"
	"fmt"
)

var InvalidStateErrorTypeId = errors.New("invalid state error type id")

type StateError struct {
	Type   uint64 `json:"code"`
	Runner string `json:"runner"`
}

const (
	StateClosedErrorType = iota
	StateRunningErrorType
	StateClosingErrorType
	StateRestartingErrorType
	StateDestroyingErrorType
	StateDestroyedErrorType
)

var Descriptions = map[uint64]string{
	StateClosedErrorType:     "runner has been closed",
	StateRunningErrorType:    "runner is running",
	StateClosingErrorType:    "runner is closing",
	StateRestartingErrorType: "runner is restarting",
	StateDestroyingErrorType: "runner is destroying",
	StateDestroyedErrorType:  "runner has been destroyed",
}

var ErrorFormats = map[uint64]string{
	StateClosedErrorType:     "runner %s has been closed",
	StateRunningErrorType:    "runner %s is running",
	StateClosingErrorType:    "runner %s is closing",
	StateRestartingErrorType: "runner %s is restarting",
	StateDestroyingErrorType: "runner %s is destroying",
	StateDestroyedErrorType:  "runner %s has been destroyed",
}

func NewStateError(e StateError, runner string) StateError {
	e.Runner = runner
	return e
}

func StateClosedError(runner string) StateError {
	return StateError{Type: StateClosedErrorType, Runner: runner}
}

func StateRunningError(runner string) StateError {
	return StateError{Type: StateRunningErrorType, Runner: runner}
}

func StateClosingError(runner string) StateError {
	return StateError{Type: StateClosingErrorType, Runner: runner}
}

func StateRestartingError(runner string) StateError {
	return StateError{Type: StateRestartingErrorType, Runner: runner}
}

func StateDestroyingError(runner string) StateError {
	return StateError{Type: StateDestroyingErrorType, Runner: runner}
}

func StateDestroyedError(runner string) StateError {
	return StateError{Type: StateDestroyedErrorType, Runner: runner}
}

func IsStateError(err error) bool {
	_, is := err.(*StateError)
	return is
}

func (e StateError) Description() string {
	if des, has := Descriptions[e.Type]; has {
		return des
	}
	panic("invalid state error type id")
}

func (e StateError) Error() string {
	if format, has := ErrorFormats[e.Type]; has {
		return fmt.Sprintf(format, e.Runner)
	}
	panic(InvalidStateErrorTypeId)
}
