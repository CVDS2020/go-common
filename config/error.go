package config

type ErrorType int

const (
	// ParseError occurs when an error occurs in the configuration parsing
	ParseError = ErrorType(iota)

	// CheckError occurs when an error occurs in the configuration checking
	CheckError

	// ReloadParseError occurs when an error occurs in the configuration
	// reloaded and parsing
	ReloadParseError

	// ReloadCheckError occurs when an error occurs in the configuration
	// reloaded and checking
	ReloadCheckError
)

// Type method return the name of error type, If the error type is undefined, return
// UNKNOWN_ERROR
func (t ErrorType) Type() string {
	switch t {
	case ParseError:
		return "PARSE_ERROR"
	case CheckError:
		return "CHECK_ERROR"
	case ReloadParseError:
		return "RELOAD_PARSE_ERROR"
	case ReloadCheckError:
		return "RELOAD_CHECK_ERROR"
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
