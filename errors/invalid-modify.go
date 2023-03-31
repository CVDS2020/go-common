package errors

import "fmt"

type InvalidModify struct {
	Target string
}

func (e *InvalidModify) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s不可以被修改", e.Target)
}
