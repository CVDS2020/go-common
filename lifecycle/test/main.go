package main

import (
	"fmt"
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/lifecycle"
	"gitee.com/sy_183/common/log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const DefaultTimeLayout = "2006-01-02 15:04:05.999999999"

var logger = assert.Must(log.Config{
	Level: log.NewAtomicLevelAt(log.DebugLevel),
	Encoder: log.NewConsoleEncoder(log.ConsoleEncoderConfig{
		DisableCaller:     true,
		DisableFunction:   true,
		DisableStacktrace: true,
		EncodeLevel:       log.CapitalColorLevelEncoder,
		EncodeTime:        log.TimeEncoderOfLayout(DefaultTimeLayout),
		EncodeDuration:    log.SecondsDurationEncoder,
	}),
}.Build())

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TestLifecycle struct {
	lifecycle.Lifecycle
	name string

	minInterruptedStartTime, maxInterruptedStartTime,
	minNotInterruptedStartTime, maxNotInterruptedStartTime,
	minStartCloseTime, maxStartCloseTime,
	minRunTime, maxRunTime,
	minRunCloseTime, maxRunCloseTime time.Duration

	startErrorRate float64
	exitErrorRate  float64
}

type TestLifecycleConfig struct {
	MinInterruptedStartTime, MaxInterruptedStartTime,
	MinNotInterruptedStartTime, MaxNotInterruptedStartTime,
	MinStartCloseTime, MaxStartCloseTime,
	MinRunTime, MaxRunTime,
	MinRunCloseTime, MaxRunCloseTime time.Duration

	StartErrorRate float64
	ExitErrorRate  float64
}

func NewTestLifecycle(name string, config TestLifecycleConfig) *TestLifecycle {
	l := &TestLifecycle{name: name}
	l.Lifecycle = lifecycle.NewWithInterruptedStart(l.start)
	l.Lifecycle.SetOnStarting(l.onStarting).
		SetOnStarted(l.onStarted).
		SetOnClose(l.onClose).
		SetOnClosed(l.onClosed)
	l.minInterruptedStartTime = config.MinInterruptedStartTime
	l.maxInterruptedStartTime = config.MaxInterruptedStartTime
	l.minNotInterruptedStartTime = config.MinNotInterruptedStartTime
	l.maxNotInterruptedStartTime = config.MaxNotInterruptedStartTime
	l.minStartCloseTime = config.MinStartCloseTime
	l.maxStartCloseTime = config.MaxStartCloseTime
	l.minRunTime = config.MinRunTime
	l.maxRunTime = config.MaxRunTime
	l.minRunCloseTime = config.MinRunCloseTime
	l.maxRunCloseTime = config.MaxRunCloseTime
	l.startErrorRate = config.StartErrorRate
	l.exitErrorRate = config.ExitErrorRate
	return l
}

func (l *TestLifecycle) randDuration(min, max time.Duration) time.Duration {
	return min + time.Duration(rand.Int63n(int64(max-min)))
}

func (l *TestLifecycle) randBool(rate float64) bool {
	return rand.Float64() < rate
}

func (l *TestLifecycle) start(_ lifecycle.Lifecycle, interrupter chan struct{}) (_ lifecycle.InterruptedRunFunc, err error) {
	select {
	case <-time.After(l.randDuration(l.minInterruptedStartTime, l.maxInterruptedStartTime)):
	case <-interrupter:
		time.Sleep(l.randDuration(l.minStartCloseTime, l.maxStartCloseTime))
		return nil, errors.New("生命周期组件启动被中断")
	}
	logger.Warn("生命周期组件启动已不可被中断", log.String("组件名称", l.name))
	time.Sleep(l.randDuration(l.minNotInterruptedStartTime, l.maxNotInterruptedStartTime))
	if l.randBool(l.startErrorRate) {
		err = errors.New("启动错误")
	}
	return func(lifecycle.Lifecycle, chan struct{}) (err error) {
		select {
		case <-time.After(l.randDuration(l.minRunTime, l.maxRunTime)):
		case <-interrupter:
			time.Sleep(l.randDuration(l.minRunCloseTime, l.maxRunCloseTime))
			return errors.New("生命周期组件运行被中断")
		}
		if l.randBool(l.exitErrorRate) {
			return errors.New("退出错误")
		}
		return nil
	}, err
}

func (l *TestLifecycle) onStarting(lifecycle.Lifecycle) {
	logger.Info("组件正在启动...", log.String("组件名称", l.name))
}

func (l *TestLifecycle) onStarted(_ lifecycle.Lifecycle, err error) {
	if err != nil {
		logger.Error("组件启动失败", log.String("组件名称", l.name), log.NamedError("错误原因", err))
	} else {
		logger.Info("组件启动成功", log.String("组件名称", l.name))
	}
}

func (l *TestLifecycle) onClose(_ lifecycle.Lifecycle, err error) {
	if err != nil {
		logger.Error("组件关闭失败", log.String("组件名称", l.name), log.NamedError("错误原因", err))
	} else {
		logger.Info("组件正在关闭", log.String("组件名称", l.name))
	}
}

func (l *TestLifecycle) onClosed(_ lifecycle.Lifecycle, err error) {
	if err != nil {
		logger.Error("组件退出错误", log.String("组件名称", l.name), log.NamedError("错误原因", err))
	} else {
		logger.Info("组件成功退出", log.String("组件名称", l.name))
	}
}

func testGroup() {
	group := lifecycle.NewGroup()
	config := TestLifecycleConfig{
		MinInterruptedStartTime:    time.Second,
		MaxInterruptedStartTime:    time.Second * 2,
		MinNotInterruptedStartTime: time.Second,
		MaxNotInterruptedStartTime: time.Second * 2,
		MinStartCloseTime:          time.Millisecond * 500,
		MaxStartCloseTime:          time.Millisecond * 800,
		MinRunTime:                 time.Second * 10,
		MaxRunTime:                 time.Second * 15,
		MinRunCloseTime:            time.Millisecond * 800,
		MaxRunCloseTime:            time.Millisecond * 1200,
		StartErrorRate:             0.1,
		ExitErrorRate:              0.1,
	}

	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("test%d", i)
		group.MustAdd(name, lifecycle.NewRetryable(NewTestLifecycle(name, config)).SetLazyStart(true).SetRetryInterval(time.Second))
	}

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill)
	go func() {
		<-exitChan
		group.Close(nil)
	}()
	go func() {
		for i := 20; i < 30; i++ {
			name := fmt.Sprintf("test%d", i)
			group.MustAdd(name, lifecycle.NewRetryable(NewTestLifecycle(name, config)).SetLazyStart(true).SetRetryInterval(time.Second))
			logger.Info("组件添加成功", log.String("组件名称", name))
			time.Sleep(time.Millisecond * 500)
		}
	}()
	group.Run()
}

func testList() {
	list := lifecycle.NewList()

	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("test%d", i)
		list.Append(NewTestLifecycle(name, TestLifecycleConfig{
			MinInterruptedStartTime:    time.Millisecond * 500,
			MaxInterruptedStartTime:    time.Millisecond * 800,
			MinNotInterruptedStartTime: time.Millisecond * 500,
			MaxNotInterruptedStartTime: time.Millisecond * 800,
			MinStartCloseTime:          time.Millisecond * 200,
			MaxStartCloseTime:          time.Millisecond * 500,
			MinRunTime:                 time.Second * 10,
			MaxRunTime:                 time.Second * 15,
			MinRunCloseTime:            time.Millisecond * 500,
			MaxRunCloseTime:            time.Millisecond * 800,
		}))
	}

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill)
	go func() {
		<-exitChan
		list.Close(nil)
	}()
	list.Run()
}

func testRetryable() {
	retryable := lifecycle.NewRetryable(NewTestLifecycle("test", TestLifecycleConfig{
		MinInterruptedStartTime:    time.Second,
		MaxInterruptedStartTime:    time.Second * 2,
		MinNotInterruptedStartTime: time.Second,
		MaxNotInterruptedStartTime: time.Second * 2,
		MinStartCloseTime:          time.Millisecond * 500,
		MaxStartCloseTime:          time.Millisecond * 800,
		MinRunTime:                 time.Second * 2,
		MaxRunTime:                 time.Second * 3,
		MinRunCloseTime:            time.Millisecond * 800,
		MaxRunCloseTime:            time.Millisecond * 1200,
		StartErrorRate:             0.5,
		ExitErrorRate:              0.5,
	})).SetLazyStart(true).SetRetryInterval(time.Second)
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill)
	go func() {
		<-exitChan
		retryable.Close(nil)
	}()
	retryable.Run()
}

func main() {
	testGroup()
	//testList()
	//testRetryable()
}
