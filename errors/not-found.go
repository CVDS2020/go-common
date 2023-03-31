package errors

import "fmt"

type NotFound struct {
	Target string
}

func NewNotFound(target string) NotFound {
	return NotFound{Target: target}
}

func (e NotFound) Error() string {
	return fmt.Sprintf("%s未找到", e.Target)
}
