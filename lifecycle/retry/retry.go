package retry

import (
	"errors"
	"gitee.com/sy_183/common/timer"
	"math"
	"time"
)

var InterruptedError = errors.New("重试被中断")

type RetryContext struct {
	Do func() error

	Interval time.Duration

	Delay func(ctx *RetryContext) time.Duration

	Interrupter chan struct{}

	MaxRetry int

	Retrievable func(ctx *RetryContext) bool

	Ignorable func(ctx *RetryContext) bool

	Error error

	Retry int
}

type Retry RetryContext

func MakeRetry(retry Retry) *RetryContext {
	c := new(RetryContext)
	*c = RetryContext(retry)
	return c
}

func (c *RetryContext) do(retryTimer *timer.Timer) bool {
	maxRetry := c.MaxRetry
	if maxRetry < 0 {
		maxRetry = math.MaxInt
	}
	// do func
	if err := c.Do(); err != nil {
		// check error ignorable
		if ignorable := c.Ignorable; ignorable != nil {
			if ignorable(c) {
				return false
			}
		}
		c.Error = err
		if err == InterruptedError {
			return false
		}

		// check retrievable
		if c.Retry >= maxRetry {
			return false
		}
		if retrievable := c.Retrievable; retrievable != nil {
			if !retrievable(c) {
				return false
			}
		}
		c.Retry++

		// get retry delay
		delay := c.Interval
		if c.Delay != nil {
			delay = c.Delay(c)
		}
		if delay <= 0 {
			delay = time.Second
		}

		retryTimer.After(delay)
		return true
	}
	return false
}

func (c *RetryContext) Todo() error {
	retryTimer := timer.NewTimer(make(chan struct{}))
	defer retryTimer.Stop()

	if retry := c.do(retryTimer); c.Error == nil {
		return nil
	} else if !retry {
		return c.Error
	}

	for {
		select {
		case <-retryTimer.C:
			if retry := c.do(retryTimer); c.Error == nil {
				return nil
			} else if !retry {
				return c.Error
			}
		case <-c.Interrupter:
			return InterruptedError
		}
	}
}
