package main

import (
	"fmt"
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/lifecycle"
	"gitee.com/sy_183/common/log"
	svc "gitee.com/sy_183/common/system/service"
	"math/rand"
	"os"
	"time"
)

var logger = assert.Must(log.Config{
	Encoder: log.NewConsoleEncoder(log.ConsoleEncoderConfig{
		DisableCaller:     true,
		DisableFunction:   true,
		DisableStacktrace: true,
		EncodeLevel:       log.CapitalColorLevelEncoder,
		EncodeTime:        log.TimeEncoderOfLayout("2006-01-02 15:04:05.999999999"),
	}),
}.Build())

type TestServer struct {
	lifecycle.Lifecycle
	closeRequestChan chan struct{}
	name             string
}

func NewTestServer(name string) *TestServer {
	s := &TestServer{closeRequestChan: make(chan struct{}, 1), name: name}
	_, s.Lifecycle = lifecycle.New(fmt.Sprintf("test server %s", name), lifecycle.Core(s.start, s.run, s.close))
	return s
}

func (s *TestServer) start() error {
	logger.Info(fmt.Sprintf("test server %s starting...", s.name))
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+500))
	//if rand.Intn(2) == 1 {
	//	logger.Error(fmt.Sprintf("test server %s start failed", s.name))
	//	return io.EOF
	//}
	logger.Info(fmt.Sprintf("test server %s started", s.name))
	return nil
}

func (s *TestServer) run() error {
	defer logger.Info(fmt.Sprintf("test server %s closed", s.name))
	select {
	case <-time.After(time.Second * time.Duration(rand.Intn(5)+5)):
		return nil
	case <-s.closeRequestChan:
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+500))
		return nil
	}
}

func (s *TestServer) close() error {
	logger.Info(fmt.Sprintf("test server %s closing", s.name))
	s.closeRequestChan <- struct{}{}
	return nil
}

func main() {
	s := lifecycle.NewGroup("test group", []lifecycle.ChildLifecycle{
		{Lifecycle: NewTestServer("1")},
		{Lifecycle: NewTestServer("2")},
		{Lifecycle: NewTestServer("3")},
		{Lifecycle: NewTestServer("4")},
		{Lifecycle: NewTestServer("5")},
		{Lifecycle: NewTestServer("6")},
		{Lifecycle: NewTestServer("7")},
		{Lifecycle: NewTestServer("8")},
		{Lifecycle: NewTestServer("9")},
		{Lifecycle: NewTestServer("10")},
		{Lifecycle: NewTestServer("11")},
		{Lifecycle: NewTestServer("12")},
		{Lifecycle: NewTestServer("13")},
		{Lifecycle: NewTestServer("14")},
	}, lifecycle.PreStart(true))
	os.Exit(svc.New("test-group", s).Run())
}
