package life_test

import (
	"fmt"
	"gitee.com/sy_183/common/lifecycle"
	"sync"
	"testing"
)

type TestServer struct {
	lifecycle.DestructibleLifecycle
	runner *lifecycle.DefaultDestructibleRunner

	closeChan chan struct{}
}

func NewTestServer() *TestServer {
	s := new(TestServer)
	s.closeChan = make(chan struct{})
	s.runner, s.DestructibleLifecycle = lifecycle.NewDestructible("test",
		lifecycle.Self(s),
		lifecycle.StartFn(s.start),
		lifecycle.RunFn(s.run),
		lifecycle.CloseFn(s.close),
		lifecycle.StartChecker(s.startChecker),
		lifecycle.CloseChecker(s.closeChecker),
		lifecycle.DestroyChecker(s.destroyChecker),
	)
	return s
}

func (t *TestServer) startChecker() error {
	fmt.Println("start checker...")
	return nil
}

func (t *TestServer) closeChecker() error {
	fmt.Println("close checker...")
	return nil
}

func (t *TestServer) destroyChecker() error {
	fmt.Println("destroy checker...")
	return nil
}

func (t *TestServer) Start() error {
	fmt.Println("Start...")
	return t.DestructibleLifecycle.Start()
}

func (t *TestServer) Close(future chan error) error {
	fmt.Println("Close...")
	return t.DestructibleLifecycle.Close(future)
}

func (t *TestServer) Destroy(future chan error) error {
	fmt.Println("Destroy...")
	return t.DestructibleLifecycle.Destroy(future)
}

func (t *TestServer) start() error {
	fmt.Println("start...")
	return nil
}

func (t *TestServer) run() error {
	fmt.Println("run...")
	<-t.closeChan
	return nil
}

func (t *TestServer) close() error {
	fmt.Println("close")
	t.closeChan <- struct{}{}
	return nil
}

func TestLife1(t *testing.T) {
	life, _ := lifecycle.NewDestructible("test")

	waiter := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			for i := 0; i < 1000000; i++ {
				life.Start()
				life.CloseWait()
				life.Start()
				life.Close(nil)
				life.Restart()
				life.CloseWait()
				life.Close(nil)
				life.Start()
				life.Restart()
				life.CloseWait()
				life.Restart()
				life.Start()
				life.CloseWait()
				life.Start()
				life.Restart()
				life.Close(nil)
				life.Restart()
				life.Start()
				life.Restart()
				life.CloseWait()
			}
			fmt.Println(life.DestroyWait())
			waiter.Done()
		}()
	}
	waiter.Wait()
}

func TestLife2(t *testing.T) {
	s := NewTestServer()

	s.Start()
	s.Restart()
	s.CloseWait()
	s.Restart()
	s.Close(nil)
	s.Restart()
	s.DestroyWait()
}
