package main

import (
	"context"
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/lifecycle"
	"gitee.com/sy_183/common/log"
	svc "gitee.com/sy_183/common/system/service"
	"os"
	"time"
)

type app struct {
	lifecycle.Lifecycle
	runner     *lifecycle.DefaultLifecycle
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     *log.Logger
}

func newApp() *app {
	a := new(app)
	a.ctx, a.cancelFunc = context.WithTimeout(context.Background(), time.Hour)
	a.runner = lifecycle.NewWithRun(a.start, a.run, a.close)
	a.Lifecycle = a.runner
	a.logger = assert.Must(log.Config{
		Level: log.NewAtomicLevelAt(log.InfoLevel),
		Encoder: log.NewConsoleEncoder(log.ConsoleEncoderConfig{
			DisableTime:       true,
			DisableName:       true,
			DisableCaller:     true,
			DisableStacktrace: true,
			DisableFunction:   true,
			EncodeLevel:       log.CapitalColorLevelEncoder,
		}),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build())
	return a
}

func (a *app) start(lifecycle.Lifecycle) error {
	a.logger.Info("app is starting...")
	time.Sleep(time.Second)
	a.logger.Info("app is running...")
	return nil
}

func (a *app) run(lifecycle.Lifecycle) error {
	<-a.ctx.Done()
	time.Sleep(time.Second)
	a.logger.Info("app is stopped")
	return nil
}

func (a *app) close(lifecycle.Lifecycle) error {
	a.logger.Info("app is stopping...")
	a.cancelFunc()
	return nil
}

func main() {
	os.Exit(svc.New("app", newApp()).Run())
}
