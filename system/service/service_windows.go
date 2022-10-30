package svc

import (
	"gitee.com/sy_183/common/lifecycle"
	"golang.org/x/sys/windows/svc"
	"os"
	"os/signal"
)

var isWindowsService = false

func init() {
	var err error
	isWindowsService, err = svc.IsWindowsService()
	if err != nil {
		panic(err)
	}
}

type windowsService struct {
	name     string
	app      lifecycle.Lifecycle
	exitCode int

	notifySignals  []os.Signal
	signalCallback func(sig os.Signal) (exit bool)

	errorCallback  func(err *Error)
	exitCodeGetter func(err *Error) int
	onStarted      func(s Service, l lifecycle.Lifecycle)
	onClosed       func(s Service, l lifecycle.Lifecycle)
}

func New(name string, app lifecycle.Lifecycle, options ...Option) Service {
	ws := &windowsService{
		name:           name,
		app:            app,
		notifySignals:  DefaultNotifySignals,
		signalCallback: DefaultSignalCallback,
		errorCallback:  DefaultErrorCallback,
		exitCodeGetter: DefaultExitCodeGetter,
	}
	for _, option := range options {
		option.apply(ws)
	}
	return ws
}

func (ws *windowsService) setSignalNotify(callback func(sig os.Signal) (exit bool), sig ...os.Signal) {
	ws.signalCallback = callback
	ws.notifySignals = sig
}

func (ws *windowsService) setErrorCallback(callback func(err *Error)) {
	ws.errorCallback = callback
}

func (ws *windowsService) setExitCodeGetter(exitCodeGetter func(err *Error) int) {
	ws.exitCodeGetter = exitCodeGetter
}

func (ws *windowsService) setOnStarted(callback func(s Service, l lifecycle.Lifecycle)) {
	ws.onStarted = callback
}

func (ws *windowsService) setOnClosed(callback func(s Service, l lifecycle.Lifecycle)) {
	ws.onClosed = callback
}

func (ws *windowsService) errorExit(typ ErrorType, err error) int {
	e := &Error{Type: typ, Err: err}
	ws.errorCallback(e)
	ws.exitCode = ws.exitCodeGetter(e)
	return ws.exitCode
}

func (ws *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}

	if err := ws.app.Start(); err != nil {
		return true, uint32(ws.errorExit(StartError, err))
	}
	future := ws.app.AddClosedFuture(make(chan error, 1))
	if ws.onStarted != nil {
		ws.onStarted(ws, ws.app)
	}
	defer func() {
		if ws.onClosed != nil {
			ws.onClosed(ws, ws.app)
		}
	}()
	s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	for {
		select {
		case err := <-future:
			if err != nil {
				return true, uint32(ws.errorExit(ExitError, err))
			}
			return false, 0
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending}
				if err := ws.app.Close(nil); err != nil {
					ws.errorCallback(&Error{Type: StopError, Err: err})
				}
				continue
			}
		}
	}
}

func (ws *windowsService) Run() int {
	if isWindowsService {
		if err := svc.Run(ws.name, ws); err != nil {
			return ws.errorExit(ServiceError, err)
		}
		return ws.exitCode
	}

	if err := ws.app.Start(); err != nil {
		return ws.errorExit(StartError, err)
	}
	closedFuture := ws.app.AddClosedFuture(make(chan error, 1))
	if ws.onStarted != nil {
		ws.onStarted(ws, ws.app)
	}
	defer func() {
		if ws.onClosed != nil {
			ws.onClosed(ws, ws.app)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, ws.notifySignals...)

	for {
		select {
		case sig := <-sigChan:
			if ws.signalCallback(sig) {
				if err := ws.app.Close(nil); err != nil {
					ws.errorCallback(&Error{Type: StopError, Err: err})
				}
			}
			continue

		case err := <-closedFuture:
			if err != nil {
				return ws.errorExit(ExitError, err)
			}
			return 0
		}
	}
}
