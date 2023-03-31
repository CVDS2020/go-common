package errors

import (
	"fmt"
	"strings"
)

func AddParentArgument(argument, parent string) string {
	if argument == "" {
		return parent
	} else {
		return parent + "." + argument
	}
}

func ReplaceParentArgument(argument, parent string) string {
	if argument == "" {
		return parent
	} else if i := strings.IndexByte(argument, '.'); i >= 0 {
		return parent + argument[i:]
	} else {
		return parent + "." + argument
	}
}

type InvalidArgument struct {
	Argument string
	Err      error
}

func NewInvalidArgument(argument string, err error) *InvalidArgument {
	return &InvalidArgument{Argument: argument, Err: err}
}

func (e *InvalidArgument) Error() string {
	if e == nil {
		return "<nil>"
	} else if e.Argument == "" {
		return "参数解析错误：" + e.Err.Error()
	} else {
		return fmt.Sprintf("参数解析错误(%s)：%s", e.Argument, e.Err.Error())
	}
}

func (e *InvalidArgument) AddParentArgument(parent string) {
	e.Argument = AddParentArgument(e.Argument, parent)
}

func (e *InvalidArgument) ReplaceParentArgument(parent string) {
	e.Argument = ReplaceParentArgument(e.Argument, parent)
}
