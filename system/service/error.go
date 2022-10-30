package svc

type ErrorType int

const (
	// StartError is the type of error returned by the Start method of Lifecycle
	StartError = ErrorType(iota)

	// ExitError is the type of error returned by the Run method of Lifecycle
	ExitError

	// StopError is the type of error returned by the Close method of Lifecycle
	StopError

	// ServiceError is the type of system service error
	ServiceError
)

// Type method return the name of error type, If the error type is undefined, return
// UNKNOWN_ERROR
func (t ErrorType) Type() string {
	switch t {
	case StartError:
		return "START_ERROR"
	case ExitError:
		return "EXIT_ERROR"
	case StopError:
		return "STOP_ERROR"
	case ServiceError:
		return "SERVICE_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

type Error struct {
	Type ErrorType
	Err  error
}

func (e Error) Error() string {
	return e.Type.Type() + ": " + e.Err.Error()
}
