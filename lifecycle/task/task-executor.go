package task

import (
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/lifecycle"
	"gitee.com/sy_183/common/utils"
	"sync"
	"sync/atomic"
)

var TaskExecutorClosedError = errors.New("任务执行器已关闭")

type TaskExecutor struct {
	lifecycle.Lifecycle
	maxTask int
	ch      atomic.Pointer[chan Task]
}

func NewTaskExecutor(maxTask int) *TaskExecutor {
	executor := &TaskExecutor{maxTask: maxTask}
	executor.Lifecycle = lifecycle.NewWithInterruptedRun(executor.start, executor.run)
	return executor
}

func (e *TaskExecutor) loadChannel() chan Task {
	if ch := e.ch.Load(); ch != nil {
		return *ch
	}
	return nil
}

func (e *TaskExecutor) start(_ lifecycle.Lifecycle, interrupter chan struct{}) error {
	ch := make(chan Task, e.maxTask)
	e.ch.Store(&ch)
	return nil
}

func (e *TaskExecutor) run(_ lifecycle.Lifecycle, interrupter chan struct{}) error {
	ch := e.loadChannel()
	defer func() {
		e.ch.Store(nil)
		close(ch)
		for task := range ch {
			task.Do(nil)
		}
	}()
	for {
		select {
		case task := <-ch:
			if task.Do(interrupter) {
				return nil
			}
		case <-interrupter:
			return nil
		}
	}
}

func (e *TaskExecutor) StartFunc() lifecycle.InterruptedStartFunc {
	return e.start
}

func (e *TaskExecutor) RunFunc() lifecycle.InterruptedRunFunc {
	return e.run
}

func (e *TaskExecutor) notPanicPush(ch chan Task, task Task) (err error) {
	defer func() {
		if _, is := recover().(error); is {
			err = TaskExecutorClosedError
		}
	}()
	ch <- task
	return nil
}

func (e *TaskExecutor) notPanicTryPush(ch chan Task, task Task) (ok bool, err error) {
	defer func() {
		if _, is := recover().(error); is {
			err = TaskExecutorClosedError
		}
	}()
	return utils.ChanTryPush(ch, task), nil
}

func (e *TaskExecutor) Async(task Task) error {
	ch := e.loadChannel()
	if ch == nil {
		return TaskExecutorClosedError
	}
	return e.notPanicPush(ch, task)
}

func (e *TaskExecutor) Sync(task Task) error {
	ch := e.loadChannel()
	if ch == nil {
		return TaskExecutorClosedError
	}
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	if err := e.notPanicPush(ch, Interrupted(func(interrupter chan struct{}) (interrupted bool) {
		defer waiter.Done()
		return task.Do(interrupter)
	})); err != nil {
		return err
	}
	waiter.Wait()
	return nil
}

func (e *TaskExecutor) Try(task Task) (ok bool, err error) {
	ch := e.loadChannel()
	if ch == nil {
		return false, TaskExecutorClosedError
	}
	return e.notPanicTryPush(ch, task)
}

func (e *TaskExecutor) Wait() error {
	return e.Sync(Nop())
}
