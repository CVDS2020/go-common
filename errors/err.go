package errors

type Err struct {
	Err error
}

func (e Err) Error() string {
	return e.Err.Error()
}

func (e Err) Unwrap() error {
	return e.Err
}
