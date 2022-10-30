package lifecycle

import (
	"gitee.com/sy_183/common/timer"
	"time"
)

type RetryableLifecycle struct {
	Lifecycle
	runner *DefaultRunner

	child            Lifecycle
	childCloseFuture chan error
	retryInterval    time.Duration
	retryTimer       *timer.Timer
	startRetry       bool

	closeRequestChan  chan struct{}
	closeResponseChan chan error

	startErrorCallback func(err error)
	exitErrorCallback  func(err error)
}

type RetryableLifecycleConfig struct {
	RetryInterval      time.Duration
	StartErrorCallback func(err error)
	ExitErrorCallback  func(err error)
}

func NewRetryable(name string, child Lifecycle, config *RetryableLifecycleConfig) *RetryableLifecycle {
	if config == nil {
		config = new(RetryableLifecycleConfig)
	}
	l := &RetryableLifecycle{
		child:              child,
		childCloseFuture:   make(chan error, 1),
		retryInterval:      config.RetryInterval,
		closeRequestChan:   make(chan struct{}, 1),
		closeResponseChan:  make(chan error, 1),
		startErrorCallback: config.StartErrorCallback,
		exitErrorCallback:  config.ExitErrorCallback,
	}
	l.runner, l.Lifecycle = New(name, StartFn(l.start), RunFn(l.run), CloseFn(l.close))
	return l
}

func (l *RetryableLifecycle) startChild() error {
	if !l.startRetry {
		// first start must be added closed future
		l.child.AddClosedFuture(l.childCloseFuture)
	} else {
		// last start error, closed future exist
	}
	if err := l.child.Start(); err != nil {
		if l.startErrorCallback != nil {
			l.startErrorCallback(err)
		}
		// start error, delay retry start
		l.startRetry = true
		return err
	}
	// start success
	l.startRetry = false
	return nil
}

func (l *RetryableLifecycle) start() error {
	l.retryTimer = timer.NewTimer(nil)
	l.startChild()
	return nil
}

func (l *RetryableLifecycle) run() error {
	closed := false
	for true {
		select {
		case err := <-l.childCloseFuture:
			// child closed
			closed = true
			if err != nil && l.exitErrorCallback != nil {
				// child exit error
				l.exitErrorCallback(err)
			}
			l.retryTimer.After(l.retryInterval)
		case <-l.retryTimer.C:
			// child retry start
			if err := l.startChild(); err == nil {
				// child start success
				closed = false
			} else {
				// child start failed
			}
		case <-l.closeRequestChan:
			goto doClose
		}
	}

doClose:
	if !closed {
		l.closeResponseChan <- l.child.Close(nil)
		// if not closed, child must be added closed future
		if err := <-l.childCloseFuture; err != nil && l.exitErrorCallback != nil {
			l.exitErrorCallback(err)
		}
	} else {
		// child has closed
		l.closeResponseChan <- nil
	}
	return nil
}

func (l *RetryableLifecycle) close() error {
	defer func() {
		l.retryTimer = nil
	}()
	l.closeRequestChan <- struct{}{}
	return <-l.closeResponseChan
}
