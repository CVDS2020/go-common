package lifecycle

import (
	"gitee.com/sy_183/common/container"
	"gitee.com/sy_183/common/lock"
	"gitee.com/sy_183/common/utils"
	"sync"
)

type childLifecycleContext[L any] struct {
	Lifecycle L
	channel   *childLifecycleChannel[L]
	err       error
}

func (c childLifecycleContext[L]) Clone(err error) (nc childLifecycleContext[L]) {
	nc = c
	nc.err = err
	return
}

func (c childLifecycleContext[L]) Complete(err error) {
	c.channel.Push(c.Clone(err))
}

type childLifecycleChannel[L any] struct {
	signal   chan struct{}
	contexts *container.LinkedList[childLifecycleContext[L]]
	mu       sync.Mutex
}

func newChildLifecycleChannel[L any]() *childLifecycleChannel[L] {
	return &childLifecycleChannel[L]{
		signal:   make(chan struct{}, 1),
		contexts: container.NewLinkedList[childLifecycleContext[L]](),
	}
}

func (c *childLifecycleChannel[L]) Signal() <-chan struct{} {
	return c.signal
}

func (c *childLifecycleChannel[L]) Push(ctx childLifecycleContext[L]) {
	lock.LockDo(&c.mu, func() {
		c.contexts.AddTail(ctx)
	})
	utils.ChanTryPush(c.signal, struct{}{})
}

func (c *childLifecycleChannel[L]) Pop() []childLifecycleContext[L] {
	return lock.LockGet(&c.mu, func() []childLifecycleContext[L] {
		contexts := make([]childLifecycleContext[L], 0, c.contexts.Len())
		for entry := c.contexts.HeadEntry(); entry != nil; entry = entry.Next() {
			contexts = append(contexts, entry.Value())
		}
		c.contexts.Clear()
		return contexts
	})
}
