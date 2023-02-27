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
	exitCodeGetter func(err *Error) int
}

func New(name string, app lifecycle.Lifecycle, options ...Option) Service {
	lss := &linuxSystemdService{
		app:            app,
		systemdNotify:  isLinuxSystemdService,
		notifySignals:  DefaultNotifySignals,
		signalCallback: DefaultSignalCallback,
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

func (lss *linuxSystemdService) setExitCodeGetter(exitCodeGetter func(err *Error) int) {
	lss.exitCodeGetter = exitCodeGetter
}

func (lss *linuxSystemdService) Run() int {
	go lss.app.Run()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, lss.notifySignals...)

	startedWaiter := lss.app.StartedWaiter()
	closedWaiter := make(lifecycle.ChanFuture[error], 1)
	for {
		select {
		case sig := <-sigChan:
			if lss.signalCallback(sig) {
				if lss.systemdNotify {
					SystemdNotify("STOPPING=1")
				}
				lss.app.Close(nil)
			}
		case err := <-startedWaiter:
			if err != nil {
				return lss.exitCodeGetter(&Error{Type: StartError, Err: err})
			}
			lss.app.AddClosedFuture(closedWaiter)
			if lss.systemdNotify {
				SystemdNotify("READY=1")
			}
		case err := <-closedWaiter:
			if err != nil {
				return lss.exitCodeGetter(&Error{Type: ExitError, Err: err})
			}
			return 0
		}
	}
}
