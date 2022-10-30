package errors

import (
	"fmt"
	"strings"
)

type Errors []error

func (es Errors) Append(err error) Errors {
	if err != nil {
		return append(es, err)
	}
	return es
}

func (es Errors) Error() string {
	if len(es) == 0 {
		return "<nil>"
	}
	ess := make([]string, len(es))
	for i, err := range es {
		s := fmt.Sprintf("%q", err)
		if i == 0 {
			s = "[" + s
		}
		if i == len(es)-1 {
			s += "]"
		}
		ess[i] = s
	}
	return strings.Join(ess, ", ")
}

func (es Errors) ToError() error {
	if len(es) == 0 {
		return nil
	}
	return es
}

func MakeErrors(es ...error) error {
	var err error
	for _, e := range es {
		err = Append(err, e)
	}
	return err
}

func Append(left error, right error) error {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	}

	if r, ok := right.(Errors); ok {
		if l, ok := left.(Errors); ok {
			return append(l, r...)
		}
		return append(Errors{left}, r...)
	} else if l, ok := left.(Errors); ok {
		return append(l, right)
	}
	return Errors{left, right}
}

type StringError string

func (s StringError) Error() string {
	return string(s)
}
