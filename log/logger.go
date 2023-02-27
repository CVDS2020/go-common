package log

import (
	"errors"
	"fmt"
	"gitee.com/sy_183/common/log/internal/bufferpool"
	"io/ioutil"
	"os"
	"strings"
)

// A Logger provides fast, leveled, structured logging. All methods are safe
// for concurrent use.
//
// The Logger is designed for contexts in which every microsecond and every
// allocation matters, so its API intentionally favors performance and type
// safety over brevity. For most applications, the SugaredLogger strikes a
// better balance between performance and ergonomics.
type Logger struct {
	core Core

	development bool
	addCaller   bool
	onFatal     CheckWriteAction // default is WriteThenFatal

	name        string
	errorOutput WriteSyncer

	addStack LevelEnabler

	callerSkip int

	clock Clock
}

// New constructs a new Logger from the provided Core and Options. If
// the passed Core is nil, it falls back to using a no-op
// implementation.
//
// This is the most flexible way to construct a Logger, but also the most
// verbose. For typical use cases, the highly-opinionated presets
// (NewProduction, NewDevelopment, and NewExample) or the Config struct are
// more convenient.
//
// For sample code, see the package-level AdvancedConfiguration example.
func New(core Core, options ...Option) *Logger {
	if core == nil {
		return NewNop()
	}
	log := &Logger{
		core:        core,
		errorOutput: Lock(os.Stderr),
		addStack:    FatalLevel + 1,
		clock:       DefaultClock,
	}
	return log.WithOptions(options...)
}

// NewNop returns a no-op Logger. It never writes out logs or internal errors,
// and it never runs user-defined hooks.
//
// Using WithOptions to replace the Core or error output of a no-op Logger can
// re-enable logging.
func NewNop() *Logger {
	return &Logger{
		core:        NewNopCore(),
		errorOutput: AddSync(ioutil.Discard),
		addStack:    FatalLevel + 1,
		clock:       DefaultClock,
	}
}

// NewProduction builds a sensible production Logger that writes InfoLevel and
// above logs to standard error as JSON.
//
// It's a shortcut for NewProductionConfig().Build(...Option).
func NewProduction(options ...Option) (*Logger, error) {
	return NewProductionConfig().Build(options...)
}

// NewDevelopment builds a development Logger that writes DebugLevel and above
// logs to standard error in a human-friendly format.
//
// It's a shortcut for NewDevelopmentConfig().Build(...Option).
func NewDevelopment(options ...Option) (*Logger, error) {
	return NewDevelopmentConfig().Build(options...)
}

// NewExample builds a Logger that's designed for use in log's testable
// examples. It writes DebugLevel and above logs to standard out as JSON, but
// omits the timestamp and calling function to keep example output
// short and deterministic.
func NewExample(options ...Option) *Logger {
	encoderCfg := JsonEncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    LowercaseLevelEncoder,
		EncodeTime:     ISO8601TimeEncoder,
		EncodeDuration: StringDurationEncoder,
	}
	core := NewCore(NewJSONEncoder(encoderCfg), os.Stdout, DebugLevel)
	return New(core).WithOptions(options...)
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func (log *Logger) Sugar() *SugaredLogger {
	core := log.clone()
	core.callerSkip += 2
	return &SugaredLogger{core}
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (log *Logger) Named(name string) *Logger {
	l := log.clone()
	l.name = name
	return l
}

func (log *Logger) SubNamed(name string) *Logger {
	if name == "" {
		return log
	}
	l := log.clone()
	if log.name == "" {
		l.name = name
	} else {
		l.name = strings.Join([]string{l.name, name}, ".")
	}
	return l
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (log *Logger) WithOptions(opts ...Option) *Logger {
	c := log.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (log *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {
		return log
	}
	l := log.clone()
	l.core = l.core.With(fields)
	return l
}

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
func (log *Logger) Check(lvl Level, msg string) *CheckedEntry {
	return log.check(lvl, msg)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Debug(msg string, fields ...Field) {
	if ce := log.check(DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Info(msg string, fields ...Field) {
	if ce := log.check(InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Warn(msg string, fields ...Field) {
	if ce := log.check(WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Error(msg string, fields ...Field) {
	if ce := log.check(ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// ErrorWith logs a message and error at ErrorLevel, return err. The message
// includes any fields passed at the log site, as well as any fields accumulated
// on the logger.
func (log *Logger) ErrorWith(msg string, err error, fields ...Field) error {
	if ce := log.check(ErrorLevel, msg); ce != nil {
		ce.Write(append([]Field{Error(err)}, fields...)...)
	}
	return err
}

func (log *Logger) ErrorMsg(err string, fields ...Field) error {
	if ce := log.check(ErrorLevel, err); ce != nil {
		ce.Write(fields...)
	}
	return errors.New(err)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (log *Logger) DPanic(msg string, fields ...Field) {
	if ce := log.check(DPanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (log *Logger) Panic(msg string, fields ...Field) {
	if ce := log.check(PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (log *Logger) Fatal(msg string, fields ...Field) {
	if ce := log.check(FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (log *Logger) Sync() error {
	return log.core.Sync()
}

// Core returns the Logger's underlying Core.
func (log *Logger) Core() Core {
	return log.core
}

func (log *Logger) clone() *Logger {
	clone := *log
	return &clone
}

func (log *Logger) check(lvl Level, msg string) *CheckedEntry {
	// Logger.check must always be called directly by a method in the
	// Logger interface (e.g., Check, Info, Fatal).
	// This skips Logger.check and the Info/Fatal/Check/etc. method that
	// called it.
	const callerSkipOffset = 2

	// Check the level first to reduce the cost of disabled log calls.
	// Since Panic and higher may exit, we skip the optimization for those levels.
	if lvl < DPanicLevel && !log.core.Enabled(lvl) {
		return nil
	}

	// Create basic checked entry thru the core; this will be non-nil if the
	// log message will actually be written somewhere.
	ent := Entry{
		LoggerName: log.name,
		Time:       log.clock.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce := log.core.Check(ent, nil)
	willWrite := ce != nil

	// Set up any required terminal behavior.
	switch ent.Level {
	case PanicLevel:
		ce = ce.Should(ent, WriteThenPanic)
	case FatalLevel:
		onFatal := log.onFatal
		// Noop is the default value for CheckWriteAction, and it leads to
		// continued execution after a Fatal which is unexpected.
		if onFatal == WriteThenNoop {
			onFatal = WriteThenFatal
		}
		ce = ce.Should(ent, onFatal)
	case DPanicLevel:
		if log.development {
			ce = ce.Should(ent, WriteThenPanic)
		}
	}

	// Only do further annotation if we're going to write this message; checked
	// entries that exist only for terminal behavior don't benefit from
	// annotation.
	if !willWrite {
		return ce
	}

	// Thread the error output through to the CheckedEntry.
	ce.ErrorOutput = log.errorOutput

	addStack := log.addStack.Enabled(ce.Level)
	if !log.addCaller && !addStack {
		return ce
	}

	// Adding the caller or stack trace requires capturing the callers of
	// this function. We'll share information between these two.
	stackDepth := stacktraceFirst
	if addStack {
		stackDepth = stacktraceFull
	}
	stack := captureStacktrace(log.callerSkip+callerSkipOffset, stackDepth)
	defer stack.Free()

	if stack.Count() == 0 {
		if log.addCaller {
			fmt.Fprintf(log.errorOutput, "%v Logger.check error: failed to get caller\n", ent.Time.UTC())
			log.errorOutput.Sync()
		}
		return ce
	}

	frame, more := stack.Next()

	if log.addCaller {
		ce.Caller = EntryCaller{
			Defined:  frame.PC != 0,
			PC:       frame.PC,
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		}
	}

	if addStack {
		buffer := bufferpool.Get()
		defer buffer.Free()

		stackfmt := newStackFormatter(buffer)

		// We've already extracted the first frame, so format that
		// separately and defer to stackfmt for the rest.
		stackfmt.FormatFrame(frame)
		if more {
			stackfmt.FormatStack(stack)
		}
		ce.Stack = buffer.String()
	}

	return ce
}
