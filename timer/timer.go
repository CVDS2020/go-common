package timer

import (
	"time"
)

type Timer struct {
	timer  *time.Timer
	closed bool
	C      chan struct{}
}

func NewTimer(c chan struct{}) *Timer {
	if c == nil {
		c = make(chan struct{}, 1)
	}
	return &Timer{C: c}
}

func (t *Timer) retryCallback() {
	t.C <- struct{}{}
}

func (t *Timer) Trigger() *Timer {
	t.C <- struct{}{}
	return t
}

func (t *Timer) After(duration time.Duration) *Timer {
	if t.timer == nil {
		t.timer = time.AfterFunc(duration, t.retryCallback)
	} else {
		t.timer.Reset(duration)
	}
	return t
}

func (t *Timer) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	}
}
