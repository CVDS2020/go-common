package log

import (
	"sync/atomic"
)

type LoggerProvider interface {
	Logger() *Logger

	SetLogger(logger *Logger)

	SwapLogger(logger *Logger) *Logger

	CompareAndSwapLogger(old *Logger, new *Logger) bool
}

type AtomicLogger struct {
	logger atomic.Pointer[Logger]
}

func NewAtomic(logger *Logger) *AtomicLogger {
	l := new(AtomicLogger)
	l.SetLogger(logger)
	return l
}

func (l *AtomicLogger) Logger() *Logger {
	return l.logger.Load()
}

func (l *AtomicLogger) SetLogger(logger *Logger) {
	l.logger.Store(logger)
}

func (l *AtomicLogger) SwapLogger(logger *Logger) *Logger {
	return l.logger.Swap(logger)
}

func (l *AtomicLogger) CompareAndSwapLogger(old *Logger, new *Logger) bool {
	return l.logger.CompareAndSwap(old, new)
}
