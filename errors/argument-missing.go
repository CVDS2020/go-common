package errors

import (
	"strings"
)

type ArgumentMissing struct {
	Arguments []string
}

func NewArgumentMissing(arguments ...string) ArgumentMissing {
	return ArgumentMissing{Arguments: arguments}
}

func (e ArgumentMissing) Error() string {
	return "缺少必要的参数(" + strings.Join(e.Arguments, ",") + ")"
}

func (e ArgumentMissing) AddParentArgument(parent string) {
	for i, argument := range e.Arguments {
		e.Arguments[i] = AddParentArgument(argument, parent)
	}
}

func (e ArgumentMissing) ReplaceParentArgument(parent string) {
	for i, argument := range e.Arguments {
		e.Arguments[i] = ReplaceParentArgument(argument, parent)
	}
}

type ArgumentMissingOne struct {
	Arguments []string
}

func NewArgumentMissingOne(arguments ...string) ArgumentMissingOne {
	return ArgumentMissingOne{Arguments: arguments}
}

func (e ArgumentMissingOne) Error() string {
	return "缺少必要的参数之一(" + strings.Join(e.Arguments, ",") + ")"
}

func (e ArgumentMissingOne) AddParentArgument(parent string) {
	for i, argument := range e.Arguments {
		e.Arguments[i] = AddParentArgument(argument, parent)
	}
}

func (e ArgumentMissingOne) ReplaceParentArgument(parent string) {
	for i, argument := range e.Arguments {
		e.Arguments[i] = ReplaceParentArgument(argument, parent)
	}
}
