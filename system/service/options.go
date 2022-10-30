package svc

import (
	"gitee.com/sy_183/common/lifecycle"
	"gitee.com/sy_183/common/log"
	"os"
	"syscall"
)

// An Option configures a Service.
type Option interface {
	apply(service Service)
}

// optionFunc wraps a func, so it satisfies the Option interface.
type optionFunc func(service Service)

func (f optionFunc) apply(service Service) {
	f(service)
}

// ErrorCallback Option specifies the callback function for the Service
// when an error occurs
func ErrorCallback(callback func(err *Error)) Option {
	type errorCallbackSetter interface {
		setErrorCallback(callback func(err *Error))
	}
	return optionFunc(func(service Service) {
		if setter, is := service.(errorCallbackSetter); is {
			setter.setErrorCallback(callback)
		}
	})
}

func DefaultErrorCallback(err *Error) {}

func LogErrorCallback(log *log.Logger) func(err *Error) {
	return func(err *Error) {
		switch err.Type {
		case StartError:
			log.ErrorWith("service start error", err.Err)
		case ExitError:
			log.ErrorWith("service exit error", err.Err)
		case StopError:
			log.ErrorWith("service stop error", err.Err)
		case ServiceError:
			log.ErrorWith("service internal error", err.Err)
		}
	}
}

// ExitCodeGetter Option specifies the callback function for the Service
// to obtain the program exit code
func ExitCodeGetter(exitCodeGetter func(err *Error) int) Option {
	type exitCodeGetterSetter interface {
		setExitCodeGetter(exitCodeGetter func(err *Error) int)
	}
	return optionFunc(func(service Service) {
		if setter, is := service.(exitCodeGetterSetter); is {
			setter.setExitCodeGetter(exitCodeGetter)
		}
	})
}

func DefaultExitCodeGetter(err *Error) int {
	if err == nil {
		return 0
	}
	return 1
}

// SignalNotify Option specifies the callback function when capturing the
// signal to be notified, if callback return true, If the callback function
// returns true, the program will start to exit
func SignalNotify(callback func(sig os.Signal) (exit bool), sig ...os.Signal) Option {
	type signalNotifySetter interface {
		setSignalNotify(callback func(sig os.Signal) (exit bool), sig ...os.Signal)
	}
	return optionFunc(func(service Service) {
		if setter, is := service.(signalNotifySetter); is {
			setter.setSignalNotify(callback, sig...)
		}
	})
}

var DefaultNotifySignals = []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP}

func DefaultSignalCallback(sig os.Signal) bool {
	switch sig {
	case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP:
		return true
	}
	return false
}

func OnStarted(callback func(s Service, l lifecycle.Lifecycle)) Option {
	type onStartedSetter interface {
		setOnStarted(callback func(s Service, l lifecycle.Lifecycle))
	}
	return optionFunc(func(service Service) {
		if setter, is := service.(onStartedSetter); is {
			setter.setOnStarted(callback)
		}
	})
}

func OnClosed(callback func(s Service, l lifecycle.Lifecycle)) Option {
	type onClosedSetter interface {
		setOnClosed(callback func(s Service, l lifecycle.Lifecycle))
	}
	return optionFunc(func(service Service) {
		if setter, is := service.(onClosedSetter); is {
			setter.setOnClosed(callback)
		}
	})
}
