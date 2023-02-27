package lifecycle

import (
	"gitee.com/sy_183/common/timer"
	"sync/atomic"
	"time"
)

const RetryableFieldName = "$retryable"

type Retryable[LIFECYCLE Lifecycle] struct {
	Lifecycle
	lifecycle LIFECYCLE

	interrupted bool

	lazyStart     atomic.Bool
	retryInterval atomic.Int64
}

func NewRetryable[LIFECYCLE Lifecycle](lifecycle LIFECYCLE) *Retryable[LIFECYCLE] {
	r := &Retryable[LIFECYCLE]{
		lifecycle: lifecycle,
	}
	lifecycle.SetField(RetryableFieldName, r)
	r.Lifecycle = NewWithInterruptedStart(r.start)
	return r
}

func (r *Retryable[LIFECYCLE]) Get() LIFECYCLE {
	return r.lifecycle
}

func (r *Retryable[LIFECYCLE]) SetLazyStart(lazyStart bool) *Retryable[LIFECYCLE] {
	r.lazyStart.Store(lazyStart)
	return r
}

func (r *Retryable[LIFECYCLE]) SetRetryInterval(interval time.Duration) *Retryable[LIFECYCLE] {
	r.retryInterval.Store(int64(interval))
	return r
}

func (r *Retryable[LIFECYCLE]) start(_ Lifecycle, interrupter chan struct{}) (InterruptedRunFunc, error) {
	var lazyStart bool

	run := func(_ Lifecycle, interrupter chan struct{}) error {
		var state State
		runningFuture := make(ChanFuture[error], 1)
		closedFuture := make(ChanFuture[error], 1)
		startRetryTimer := timer.NewTimer(make(chan struct{}, 1))
		defer func() {
			r.interrupted = false
			startRetryTimer.Stop()
		}()

		if lazyStart {
			state = StateClosed
			startRetryTimer.Trigger()
		} else {
			state = StateRunning
			r.lifecycle.AddClosedFuture(closedFuture)
		}

		for {
			select {
			case <-startRetryTimer.C:
				// 如果启动了定时器，说明生命周期组件一定为关闭状态，有中断信号则直接退出，所以此处一定没有标记中断，
				// 此时需要启动组件并添加启动完成的追踪器到组件
				r.lifecycle.AddStartedFuture(runningFuture)
				state = StateStarting
				r.lifecycle.Background()
			case err := <-runningFuture:
				// 生命周期组件启动完成，如果启动错误，在这种情况下如果标记了中断，则直接退出，否则启动定时器，定时器
				// 触发后执行重启。如果启动成功，不管有没有中断信号，向组件添加关闭的追踪器
				if err != nil {
					state.ToClosed()
					if r.interrupted {
						return nil
					}
					startRetryTimer.After(time.Duration(r.retryInterval.Load()))
					continue
				}
				state.ToRunning()
				r.lifecycle.AddClosedFuture(closedFuture)
			case <-closedFuture:
				// 生命周期组件退出，如果标记了中断，则直接退出，否则启动定时器，定时器触发后执行重启
				state.ToClosed()
				if r.interrupted {
					return nil
				}
				startRetryTimer.After(time.Duration(r.retryInterval.Load()))
			case <-interrupter:
				// 中断信号只会出现一次，如果此时生命周期组件为关闭状态，则直接退出，否则对组件执行关闭操作
				r.interrupted = true
				if state.Closed() {
					return nil
				}
				r.lifecycle.Close(nil)
			}
		}
	}

	if !r.lazyStart.Load() {
		r.lifecycle.Background()
		for {
			select {
			case err := <-r.lifecycle.StartedWaiter():
				if err != nil {
					r.interrupted = false
				}
				return run, err
			case <-interrupter:
				r.interrupted = true
				r.lifecycle.Close(nil)
			}
		}
	}
	lazyStart = true
	return run, nil
}
