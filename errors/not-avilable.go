package errors

import "fmt"

type NotAvailable struct {
	Target string
}

func NewNotAvailable(target string) NotAvailable {
	return NotAvailable{Target: target}
}

func (e NotAvailable) Error() string {
	return fmt.Sprintf("%s不可用", e.Target)
}
