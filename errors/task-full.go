package errors

import "fmt"

type TaskFull struct {
	Action  string
	MaxTask int
}

func (e *TaskFull) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.MaxTask != 0 {
		return fmt.Sprintf("%s任务已满，最大任务数为%d", e.Action, e.MaxTask)
	} else {
		return fmt.Sprintf("%s任务已满", e.Action)
	}
}
