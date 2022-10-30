package svc

import (
	"gitee.com/sy_183/common/lifecycle"
	"os"
	"os/signal"
)

var isLinuxSystemdService = false

func init() {
	if _, err := os.Stat("/run/systemd/system"); err != nil {
		isLinuxSystemdService = true
	}
}

type linuxSystemdService struct {
	app      lifecycle.Lifecycle
	exitCode int

	systemdNotify bool

	notifySignals  []os.Signal
	signalCallback func(sig os.Signal) (exit bool)

	errorCallback  func(err *Error)
	exitCodeGetter func(err *Error) int
	onStarted      func(s Service, l lifecycle.Lifecycle)
	onClosed       func(s Service, l lifecycle.Lifecycle)
}

func New(name string, app lifecycle.Lifecycle, options ...Option) Service {
	lss := &linuxSystemdService{
		app:            app,
		systemdNotify:  true,
		notifySignals:  DefaultNotifySignals,
		signalCallback: DefaultSignalCallback,
		errorCallback:  DefaultErrorCallback,
		exitCodeGetter: DefaultExitCodeGetter,
	}
	for _, option := range options {
		option.apply(lss)
	}
	return lss
}

func (lss *linuxSystemdService) setSignalNotify(callback func(sig os.Signal) (exit bool), sig ...os.Signal) {
	lss.signalCallback = callback
	lss.notifySignals = sig
}

func (lss *linuxSystemdService) setErrorCallback(callback func(err *Error)) {
	lss.errorCallback = callback
}

func (lss *linuxSystemdService) setExitCodeGetter(exitCodeGetter func(err *Error) int) {
	lss.exitCodeGetter = exitCodeGetter
}

func (lss *linuxSystemdService) setOnStarted(callback func(s Service, l lifecycle.Lifecycle)) {
	lss.onStarted = callback
}

func (lss *linuxSystemdService) setOnClosed(callback func(s Service, l lifecycle.Lifecycle)) {
	lss.onClosed = callback
}

func (lss *linuxSystemdService) errorExit(typ ErrorType, err error) int {
	e := &Error{Type: typ, Err: err}
	lss.errorCallback(e)
	lss.exitCode = lss.exitCodeGetter(e)
	return lss.exitCode
}

func (lss *linuxSystemdService) Run() int {
	if err := lss.app.Start(); err != nil {
		return lss.errorExit(StartError, err)
	}
	closedFuture := lss.app.AddClosedFuture(make(chan error, 1))
	if lss.onStarted != nil {
		lss.onStarted(lss, lss.app)
	}
	defer func() {
		if lss.onClosed != nil {
			lss.onClosed(lss, lss.app)
		}
	}()
	if lss.systemdNotify {
		SystemdNotify("READY=1")
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, lss.notifySignals...)

	for {
		select {
		case sig := <-sigChan:
			if lss.signalCallback(sig) {
				if lss.systemdNotify {
					SystemdNotify("STOPPING=1")
				}
				if err := lss.app.Close(nil); err != nil {
					lss.errorCallback(&Error{Type: StopError, Err: err})
				}
			}
			continue
		case err := <-closedFuture:
			if err != nil {
				return lss.errorExit(ExitError, err)
			}
			return 0
		}
	}
}
