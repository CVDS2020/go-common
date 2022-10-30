package config

import (
	"gitee.com/sy_183/common/log"
	"sync"
	"sync/atomic"
	"unsafe"
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
			logger.Fatal("load config error", log.Error(err))
		case CheckError:
			logger.Fatal("check config error", log.Error(err))
		case ReloadParseError:
			logger.ErrorWith("reload config error", err)
		case ReloadCheckError:
			logger.ErrorWith("reload check config error", err)
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

type Context[C any] struct {
	Parser

	config      unsafe.Pointer
	initializer sync.Once

	errorCallback func(err *Error)

	configReloadedCallbacks       []func(oc, nc *C)
	configReloadedCallbacksLocker sync.Mutex
	configReloadCheckers          []func(oc, nc *C) error
	configReloadCheckersLocker    sync.Mutex
}

func NewContext[C any](options ...Option[C]) *Context[C] {
	ctx := new(Context[C])
	for _, option := range options {
		option.apply(ctx)
	}
	return ctx
}

func (c *Context[C]) store(config *C) {
	atomic.StorePointer(&c.config, unsafe.Pointer(config))
}

func (c *Context[C]) load() *C {
	return (*C)(atomic.LoadPointer(&c.config))
}

func (c *Context[C]) RegisterConfigReloadedCallback(callback func(oc, nc *C)) {
	c.configReloadedCallbacksLocker.Lock()
	c.configReloadedCallbacks = append(c.configReloadedCallbacks, callback)
	c.configReloadedCallbacksLocker.Unlock()
}

func (c *Context[C]) RegisterConfigReloadChecker(checker func(oc, nc *C) error) {
	c.configReloadCheckersLocker.Lock()
	c.configReloadCheckers = append(c.configReloadCheckers, checker)
	c.configReloadCheckersLocker.Unlock()
}

func (c *Context[C]) configReloaded(oc, nc *C) {
	c.configReloadedCallbacksLocker.Lock()
	for _, callback := range c.configReloadedCallbacks {
		callback(oc, nc)
	}
	c.configReloadedCallbacksLocker.Unlock()
}

func (c *Context[C]) configReloadCheck(oc, nc *C) error {
	c.configReloadCheckersLocker.Lock()
	defer c.configReloadCheckersLocker.Unlock()
	for _, checker := range c.configReloadCheckers {
		if err := checker(oc, nc); err != nil {
			return err
		}
	}
	return nil
}

func (c *Context[C]) ConfigP() *C {
	if cfg := c.load(); cfg != nil {

		return cfg
	}
	c.initializer.Do(c.initConfig)
	return c.load()
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
	c.store(nc)
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
	c.store(nc)
	c.configReloaded(oc, nc)
}
