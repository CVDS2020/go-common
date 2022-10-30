package lifecycle

import (
	"gitee.com/sy_183/common/def"
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/timer"
	"time"
)

const (
	contextStateStart = iota
	contextStateStarted
	contextStateStartFailed
	contextStateClosed
)

type ChildLifecycle struct {
	Lifecycle
	RetryInterval      time.Duration
	OnStarted          func(l Lifecycle)
	OnClosed           func(l Lifecycle)
	StartErrorCallback func(err error)
	CloseErrorCallback func(err error)
	ExitErrorCallback  func(err error)
}

type context struct {
	ChildLifecycle
	index      int
	retryTimer *timer.Timer
}

type contextState struct {
	*context
	state int
}

func newContext(index int, lifecycle ChildLifecycle) *context {
	ctx := &context{
		ChildLifecycle: lifecycle,
		index:          index,
		retryTimer:     timer.NewTimer(nil),
	}
	def.SetDefaultP(&ctx.RetryInterval, time.Second)
	return ctx
}

func (c *context) delayStart(stateChan chan contextState) {
	c.retryTimer.After(c.RetryInterval)
	go func() {
		<-c.retryTimer.C
		stateChan <- contextState{
			context: c,
			state:   contextStateStart,
		}
	}()
}

func (c *context) start(stateChan chan contextState) {
	go func() {
		if err := c.Start(); err != nil {
			if c.StartErrorCallback != nil {
				c.StartErrorCallback(err)
			}
			stateChan <- contextState{
				context: c,
				state:   contextStateStartFailed,
			}
			return
		}
		if c.OnStarted != nil {
			c.OnStarted(c.Lifecycle)
		}
		stateChan <- contextState{
			context: c,
			state:   contextStateStarted,
		}
		future := c.AddClosedFuture(nil)
		go func() {
			if err := <-future; err != nil {
				if c.ExitErrorCallback != nil {
					c.ExitErrorCallback(err)
				}
			}
			if c.OnClosed != nil {
				c.OnClosed(c.Lifecycle)
			}
			stateChan <- contextState{
				context: c,
				state:   contextStateClosed,
			}
		}()
	}()
}

func (c *context) stop(future chan error) {
	if err := c.Close(future); err != nil {
		if c.CloseErrorCallback != nil {
			c.CloseErrorCallback(err)
		}
	}
}

type Group struct {
	Lifecycle
	contexts         []*context
	stateChan        chan contextState
	preStart         bool
	closeRequestChan chan struct{}
}

func NewGroup(name string, children []ChildLifecycle, options ...GroupOption) *Group {
	g := &Group{
		contexts:         make([]*context, 0, len(children)),
		stateChan:        make(chan contextState, len(children)),
		closeRequestChan: make(chan struct{}, 1),
	}
	for _, option := range options {
		option.apply(g)
	}
	for i, child := range children {
		g.contexts = append(g.contexts, newContext(i, child))
	}
	_, g.Lifecycle = New(name, Core(g.start, g.run, g.close))
	return g
}

func (g *Group) start() error {
	if g.preStart {
		var doClose bool
		// running lifecycles index
		running := make(map[int]struct{})
		// starting lifecycles index
		starting := make(map[int]struct{})
		// start all lifecycle
		for _, ctx := range g.contexts {
			starting[ctx.index] = struct{}{}
			ctx.start(g.stateChan)
		}
		var es error
		for {
			select {
			case ctx := <-g.stateChan:
				switch ctx.state {
				case contextStateStarted:
					delete(starting, ctx.index)
					running[ctx.index] = struct{}{}
					if doClose {
						ctx.stop(nil)
					} else if len(starting) == 0 {
						for i, ctx := range g.contexts {
							if _, in := running[i]; !in {
								g.stateChan <- contextState{
									context: ctx,
									state:   contextStateClosed,
								}
							}
						}
						return nil
					}
				case contextStateStartFailed:
					es = errors.Append(es, ctx.Error())
					delete(starting, ctx.index)
					if !doClose {
						doClose = true
						// stop all running lifecycle
						for i, ctx := range g.contexts {
							if _, in := running[i]; in {
								ctx.stop(nil)
							}
						}
					}
					if len(running) == 0 && len(starting) == 0 {
						// all lifecycle has been stopped
						return es
					}
				case contextStateClosed:
					delete(running, ctx.index)
					if doClose {
						// has lifecycle stopped
						if len(running) == 0 && len(starting) == 0 {
							// all lifecycle has been stopped
							return es
						}
					}
				}
			}
		}
	}
	return nil
}

func (g *Group) run() error {
	var doClose bool

	// running lifecycles index
	running := make(map[int]struct{})
	// starting lifecycles index
	starting := make(map[int]struct{})

	// start all lifecycle
	if g.preStart {
		for _, ctx := range g.contexts {
			running[ctx.index] = struct{}{}
		}
	} else {
		for _, ctx := range g.contexts {
			starting[ctx.index] = struct{}{}
			ctx.start(g.stateChan)
		}
	}

	// lifecycle monitor
	for {
		select {
		case ctx := <-g.stateChan:
			switch ctx.state {
			case contextStateStart:
				if !doClose {
					// do start
					starting[ctx.index] = struct{}{}
					ctx.start(g.stateChan)
				}
			case contextStateStarted:
				delete(starting, ctx.index)
				running[ctx.index] = struct{}{}
				if doClose {
					ctx.stop(nil)
				}
			case contextStateStartFailed:
				delete(starting, ctx.index)
				if doClose {
					if len(running) == 0 && len(starting) == 0 {
						// all lifecycle has been stopped
						return nil
					}
					continue
				} else {
					ctx.delayStart(g.stateChan)
				}
			case contextStateClosed:
				// has lifecycle stopped
				delete(running, ctx.index)
				if doClose {
					if len(running) == 0 && len(starting) == 0 {
						// all lifecycle has been stopped
						return nil
					}
					continue
				}
				// wait a time, restart lifecycle
				ctx.delayStart(g.stateChan)
			}
		case <-g.closeRequestChan:
			doClose = true
			// stop all running lifecycle
			for i, ctx := range g.contexts {
				if _, in := running[i]; in {
					ctx.stop(nil)
				}
			}
			if len(running) == 0 && len(starting) == 0 {
				// all lifecycle has been stopped
				return nil
			}
		}
	}
}

func (g *Group) close() error {
	g.closeRequestChan <- struct{}{}
	return nil
}
