package log

import (
	"fmt"
	"reflect"
)

// Encodes the given error into fields of an object. A field with the given
// name is added for the error message.
//
// If the error implements fmt.Formatter, a field with the name ${key}Verbose
// is also added with the full verbose error message.
//
// Finally, if the error implements errorGroup (from go.uber.org/multierr) or
// causer (from github.com/pkg/errors), a ${key}Causes field is added with an
// array of objects containing the errors this error was comprised of.
//
//  {
//    "error": err.Error(),
//    "errorVerbose": fmt.Sprintf("%+v", err),
//    "errorCauses": [
//      ...
//    ],
//  }
func encodeError(key string, err error, enc ObjectEncoder) (retErr error) {
	// Try to capture panics (from nil references or otherwise) when calling
	// the Error() method
	defer func() {
		if rerr := recover(); rerr != nil {
			// If it's a nil pointer, just say "<nil>". The likeliest causes are a
			// error that fails to guard against nil or a nil pointer for a
			// value receiver, and in either case, "<nil>" is a nice result.
			if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
				enc.AddString(key, "<nil>")
				return
			}

			retErr = fmt.Errorf("PANIC=%v", rerr)
		}
	}()

	switch e := err.(type) {
	case ArrayMarshaler:
		return enc.AddArray(key, e)
	case ObjectMarshaler:
		return enc.AddObject(key, e)
	case errorGroup:
		return enc.AddArray(key, errArray(e.Errors()))
	default:
		enc.AddString(key, err.Error())
	}
	return nil
}

type errorGroup interface {
	// Provides read-only access to the underlying list of errors, preferably
	// without causing any allocs.
	Errors() []error
}

// Note that errArray and errArrayElem are very similar to the version
// implemented in the top-level error.go file. We can't re-use this because
// that would require exporting errArray as part of the core API.

// Encodes a list of errors using the standard error encoding logic.
type errArray []error

func (errs errArray) MarshalLogArray(arr ArrayEncoder) error {
	for i := range errs {
		if errs[i] == nil {
			continue
		}

		switch e := errs[i].(type) {
		case ArrayMarshaler:
			arr.AppendArray(e)
		case ObjectMarshaler:
			arr.AppendObject(e)
		case errorGroup:
			arr.AppendArray(errArray(e.Errors()))
		default:
			arr.AppendString(e.Error())
		}
	}
	return nil
}

// Error is shorthand for the common idiom NamedError("error", err).
func Error(err error) Field {
	return NamedError("error", err)
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: key, Type: ErrorType, Interface: err}
}
