package log

// ObjectMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
//
// Note: ObjectMarshaler is only used when log.Object is used or when
// passed directly to log.Any. It is not used when reflection-based
// encoding is used.
type ObjectMarshaler interface {
	MarshalLogObject(ObjectEncoder) error
}

// ObjectMarshalerFunc is a type adapter that turns a function into an
// ObjectMarshaler.
type ObjectMarshalerFunc func(ObjectEncoder) error

// MarshalLogObject calls the underlying function.
func (f ObjectMarshalerFunc) MarshalLogObject(enc ObjectEncoder) error {
	return f(enc)
}

// ArrayMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
//
// Note: ArrayMarshaler is only used when log.Array is used or when
// passed directly to log.Any. It is not used when reflection-based
// encoding is used.
type ArrayMarshaler interface {
	MarshalLogArray(ArrayEncoder) error
}

// ArrayMarshalerFunc is a type adapter that turns a function into an
// ArrayMarshaler.
type ArrayMarshalerFunc func(ArrayEncoder) error

// MarshalLogArray calls the underlying function.
func (f ArrayMarshalerFunc) MarshalLogArray(enc ArrayEncoder) error {
	return f(enc)
}
