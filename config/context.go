package config

import (
	"gitee.com/sy_183/common/container"
	"gitee.com/sy_183/common/lock"
	"gitee.com/sy_183/common/log"
	"sync"
	"sync/atomic"
)

// An Option configures a Context.
type Option[C any] interface {
	apply(ctx *Context[C])
}

// optionFunc wraps a func, so it satisfies the Option interface.
type optionFunc[C any] func(ctx *Context[C])

func (f optionFunc[C]) apply(ctx *Context[C]) {
	f(ctx)
}

func ErrorCallback[C any](callback func(err *Error)) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.errorCallback = callback
	})
}

func LoggerErrorCallback(logger *log.Logger) func(err *Error) {
	return func(err *Error) {
		switch err.Type {
		case ParseError:
			logger.Fatal("加载配置失败", log.Error(err))
		case CheckError:
			logger.Fatal("检查配置失败", log.Error(err))
		case ReloadParseError:
			logger.ErrorWith("重新加载配置失败", err)
		case ReloadCheckError:
			logger.ErrorWith("重新加载配置检查失败", err)
		}
	}
}

func AddBytes[C any](bs []byte, typ Type) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.AddBytes(bs, typ)
	})
}

func SetBytes[C any](bs []byte, typ Type) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.SetBytes(bs, typ)
	})
}

func AddFile[C any](path string, typ *Type) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.AddFile(path, typ)
	})
}

func SetFile[C any](path string, typ *Type) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.SetFile(path, typ)
	})
}

func AddFilePrefix[C any](prefix string, types ...Type) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.AddFilePrefix(prefix, types...)
	})
}

func SetFilePrefix[C any](prefix string, types ...Type) Option[C] {
	return optionFunc[C](func(ctx *Context[C]) {
		ctx.SetFilePrefix(prefix, types...)
	})
}

type (
	OnConfigReloaded[C any]    func(oc, nc *C)
	ConfigReloadChecker[C any] func(oc, nc *C) error
)

type Context[C any] struct {
	Parser

	config      atomic.Pointer[C]
	initializer sync.Once

	errorCallback func(err *Error)

	configReloadedCallbackId atomic.Uint64
	configReloadCheckerId    atomic.Uint64

	configReloadedCallbacks       *container.LinkedMap[uint64, OnConfigReloaded[C]]
	configReloadedCallbacksLocker sync.Mutex
	configReloadCheckers          *container.LinkedMap[uint64, ConfigReloadChecker[C]]
	configReloadCheckersLocker    sync.Mutex
}

func NewContext[C any](options ...Option[C]) *Context[C] {
	ctx := &Context[C]{
		configReloadedCallbacks: container.NewLinkedMap[uint64, OnConfigReloaded[C]](0),
		configReloadCheckers:    container.NewLinkedMap[uint64, ConfigReloadChecker[C]](0),
	}
	for _, option := range options {
		option.apply(ctx)
	}
	return ctx
}

func (c *Context[C]) RegisterConfigReloadedCallback(callback OnConfigReloaded[C]) uint64 {
	return lock.LockGet(&c.configReloadedCallbacksLocker, func() uint64 {
		id := c.configReloadedCallbackId.Add(1)
		c.configReloadedCallbacks.PutIfAbsent(id, callback)
		return id
	})
}

func (c *Context[C]) UnregisterConfigReloadedCallback(id uint64) {
	lock.LockDo(&c.configReloadedCallbacksLocker, func() {
		c.configReloadedCallbacks.Remove(id)
	})
}

func (c *Context[C]) RegisterConfigReloadChecker(checker ConfigReloadChecker[C]) uint64 {
	return lock.LockGet(&c.configReloadCheckersLocker, func() uint64 {
		id := c.configReloadCheckerId.Add(1)
		c.configReloadCheckers.PutIfAbsent(id, checker)
		return id
	})
}

func (c *Context[C]) UnregisterConfigReloadChecker(id uint64) {
	lock.LockDo(&c.configReloadCheckersLocker, func() {
		c.configReloadCheckers.Remove(id)
	})
}

func (c *Context[C]) configReloaded(oc, nc *C) {
	lock.LockDo(&c.configReloadedCallbacksLocker, func() {
		for entry := c.configReloadedCallbacks.FirstEntry(); entry != nil; entry = entry.Next() {
			entry.Value()(oc, nc)
		}
	})
}

func (c *Context[C]) configReloadCheck(oc, nc *C) error {
	return lock.LockGet(&c.configReloadCheckersLocker, func() error {
		for entry := c.configReloadCheckers.FirstEntry(); entry != nil; entry = entry.Next() {
			if err := entry.Value()(oc, nc); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *Context[C]) ConfigP() *C {
	if cfg := c.config.Load(); cfg != nil {
		return cfg
	}
	c.initializer.Do(c.initConfig)
	return c.config.Load()
}

func (c *Context[C]) Config() C {
	return *c.ConfigP()
}

func (c *Context[C]) initConfig() {
	nc := new(C)
	if err := c.Parser.Unmarshal(nc); err != nil {
		if c.errorCallback != nil {
			c.errorCallback(&Error{Type: ParseError, Err: err})
		}
		panic(err)
	}
	if err := c.configReloadCheck(nil, nc); err != nil {
		if c.errorCallback != nil {
			c.errorCallback(&Error{Type: CheckError, Err: err})
		}
		panic(err)
	}
	c.config.Store(nc)
	c.configReloaded(nil, nc)
}

func (c *Context[C]) ReloadConfig() {
	oc := c.ConfigP()
	nc := new(C)
	if err := c.Parser.Unmarshal(nc); err != nil {
		if c.errorCallback != nil {
			c.errorCallback(&Error{Type: ReloadParseError, Err: err})
		}
		return
	}
	if err := c.configReloadCheck(oc, nc); err != nil {
		if c.errorCallback != nil {
			c.errorCallback(&Error{Type: ReloadCheckError, Err: err})
		}
		return
	}
	c.config.Store(nc)
	c.configReloaded(oc, nc)
}
