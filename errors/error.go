package errors

import (
	_ "unsafe"
)

//go:linkname New errors.New
func New(text string) error

//go:linkname Unwrap errors.Unwrap
func Unwrap(err error) error

//go:linkname As errors.As
func As(err error, target any) bool

//go:linkname Is errors.Is
func Is(err, target error) bool
