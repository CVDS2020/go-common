package lifecycle

import (
	"math/rand"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	onStarting := func(lifecycle Lifecycle) { t.Logf("%s starting", lifecycle.Field("name")) }
	onStarted := func(lifecycle Lifecycle, err error) {
		if err != nil {
			t.Logf("%s started with error(%s)", lifecycle.Field("name"), err.Error())
		} else {
			t.Logf("%s started", lifecycle.Field("name"))
		}
	}
	onClose := func(lifecycle Lifecycle, err error) {
		if err != nil {
			t.Logf("%s close with error(%s)", lifecycle.Field("name"), err.Error())
		} else {
			t.Logf("%s close", lifecycle.Field("name"))
		}
	}
	onClosed := func(lifecycle Lifecycle, err error) {
		if err != nil {
			t.Logf("%s exit with error(%s)", lifecycle.Field("name"), err.Error())
		} else {
			t.Logf("%s exit", lifecycle.Field("name"))
		}
	}
	starter := func(lifecycle Lifecycle, interrupter chan struct{}) (runFn InterruptedRunFunc, err error) {
		select {
		case <-time.After(time.Millisecond*time.Duration(rand.Intn(1000)) + time.Second):
		case <-interrupter:
			time.Sleep(time.Millisecond*time.Duration(rand.Intn(300)) + time.Millisecond*300)
			return nil, NewInterruptedError("", "启动")
		}
		return func(Lifecycle, chan struct{}) error {
			select {
			case <-time.After(time.Millisecond*time.Duration(rand.Intn(5000)) + time.Second*5):
			case <-interrupter:
				time.Sleep(time.Millisecond*time.Duration(rand.Intn(500)) + time.Millisecond*500)
				return nil
			}
			return nil
		}, nil
	}
	g := NewGroup().
		MustAdd("test1", NewWithInterruptedStart(starter).OnStarting(onStarting).OnStarted(onStarted).OnClose(onClose).OnClosed(onClosed).SetField("name", "test1")).SetCloseAllOnExit(false).Group().
		MustAdd("test2", NewWithInterruptedStart(starter).OnStarting(onStarting).OnStarted(onStarted).OnClose(onClose).OnClosed(onClosed).SetField("name", "test2")).SetCloseAllOnExit(false).Group().
		MustAdd("test3", NewWithInterruptedStart(starter).OnStarting(onStarting).OnStarted(onStarted).OnClose(onClose).OnClosed(onClosed).SetField("name", "test3")).Group()
	g.Run()
}
