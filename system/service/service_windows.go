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
	exitCodeGetter func(err *Error) int
}

func New(name string, app lifecycle.Lifecycle, options ...Option) Service {
	ws := &windowsService{
		name:           name,
		app:            app,
		notifySignals:  DefaultNotifySignals,
		signalCallback: DefaultSignalCallback,
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

func (ws *windowsService) setExitCodeGetter(exitCodeGetter func(err *Error) int) {
	ws.exitCodeGetter = exitCodeGetter
}

func (ws *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	go ws.app.Run()
	s <- svc.Status{State: svc.StartPending}

	startedWaiter := ws.app.StartedWaiter()
	closedWaiter := make(lifecycle.ChanFuture[error], 1)

	for {
		select {
		case err := <-startedWaiter:
			if err != nil {
				return true, uint32(ws.exitCodeGetter(&Error{Type: StartError, Err: err}))
			}
			ws.app.AddClosedFuture(closedWaiter)
			s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
		case err := <-closedWaiter:
			if err != nil {
				return true, uint32(ws.exitCodeGetter(&Error{Type: ExitError, Err: err}))
			}
			return false, 0
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending}
				ws.app.Close(nil)
			}
		}
	}
}

func (ws *windowsService) Run() int {
	if isWindowsService {
		if err := svc.Run(ws.name, ws); err != nil {
			return ws.exitCodeGetter(&Error{Type: ServiceError, Err: err})
		}
		return ws.exitCode
	}
	go ws.app.Run()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, ws.notifySignals...)

	startedWaiter := ws.app.StartedWaiter()
	closedWaiter := make(lifecycle.ChanFuture[error], 1)
	for {
		select {
		case sig := <-sigChan:
			if ws.signalCallback(sig) {
				ws.app.Close(nil)
			}
		case err := <-startedWaiter:
			if err != nil {
				return ws.exitCodeGetter(&Error{Type: StartError, Err: err})
			}
			ws.app.AddClosedFuture(closedWaiter)
		case err := <-closedWaiter:
			if err != nil {
				return ws.exitCodeGetter(&Error{Type: ExitError, Err: err})
			}
			return 0
		}
	}
}
