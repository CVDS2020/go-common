package log

import (
	"bytes"
	"errors"
	"fmt"
	"gitee.com/sy_183/common/uns"
	"gopkg.in/yaml.v3"
	"reflect"
	"sync/atomic"
)

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

// ParseLevel parses a level based on the lower-case or all-caps ASCII
// representation of the log level. If the provided ASCII representation is
// invalid an error is returned.
//
// This is particularly useful when dealing with text input to configure log
// levels.
func ParseLevel(text string) (Level, error) {
	var level Level
	err := level.UnmarshalText([]byte(text))
	return level, err
}

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DPanicLevel:
		return "dpanic"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// CapitalString returns an all-caps ASCII representation of the log level.
func (l Level) CapitalString() string {
	// Printing levels in all-caps is common enough that we should export this
	// functionality.
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DPanicLevel:
		return "DPANIC"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// MarshalText marshals the Level to text. Note that the text representation
// drops the -Level suffix (see example).
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText unmarshals text to a level. Like MarshalText, UnmarshalText
// expects the text representation of a Level to drop the -Level suffix (see
// example).
//
// In particular, this makes it easy to configure logging levels using YAML,
// TOML, or JSON files.
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "dpanic", "DPANIC":
		*l = DPanicLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

// Set sets the level for the flag.Value interface.
func (l *Level) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

// Get gets the level for the flag.Getter interface.
func (l *Level) Get() interface{} {
	return *l
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// LevelEnabler decides whether a given logging level is enabled when logging a
// message.
//
// Enablers are intended to be used to implement deterministic filters;
// concerns like sampling are better implemented as a Core.
//
// Each concrete Level value implements a static LevelEnabler which returns
// true for itself and all higher logging levels. For example WarnLevel.Enabled()
// will return true for WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, and
// FatalLevel, but return false for InfoLevel and DebugLevel.
type LevelEnabler interface {
	Enabled(Level) bool
}

// LevelEnablerFunc is a convenient way to implement LevelEnabler with
// an anonymous function.
//
// It's particularly useful when splitting log output between different
// outputs (e.g., standard error and standard out). For sample code, see the
// package-level AdvancedConfiguration example.
type LevelEnablerFunc func(Level) bool

// Enabled calls the wrapped function.
func (f LevelEnablerFunc) Enabled(lvl Level) bool { return f(lvl) }

// An AtomicLevel is an atomically changeable, dynamic logging level. It lets
// you safely change the log level of a tree of loggers (the root logger and
// any children created by adding context) at runtime.
//
// The AtomicLevel itself is an http.Handler that serves a JSON endpoint to
// alter its level.
//
// AtomicLevels must be created with the NewAtomicLevel constructor to allocate
// their internal atomic pointer.
type AtomicLevel struct {
	l *int32
}

// NewAtomicLevel creates an AtomicLevel with InfoLevel and above logging
// enabled.
func NewAtomicLevel() AtomicLevel {
	l := int32(InfoLevel)
	return AtomicLevel{
		l: &l,
	}
}

// NewAtomicLevelAt is a convenience function that creates an AtomicLevel
// and then calls SetLevel with the given level.
func NewAtomicLevelAt(l Level) AtomicLevel {
	a := NewAtomicLevel()
	a.SetLevel(l)
	return a
}

// ParseAtomicLevel parses an AtomicLevel based on a lowercase or all-caps ASCII
// representation of the log level. If the provided ASCII representation is
// invalid an error is returned.
//
// This is particularly useful when dealing with text input to configure log
// levels.
func ParseAtomicLevel(text string) (AtomicLevel, error) {
	a := NewAtomicLevel()
	l, err := ParseLevel(text)
	if err != nil {
		return a, err
	}

	a.SetLevel(l)
	return a, nil
}

// Enabled implements the LevelEnabler interface, which allows the
// AtomicLevel to be used in place of traditional static levels.
func (lvl AtomicLevel) Enabled(l Level) bool {
	return lvl.Level().Enabled(l)
}

// Level returns the minimum enabled log level.
func (lvl AtomicLevel) Level() Level {
	return Level(int8(atomic.LoadInt32(lvl.l)))
}

// SetLevel alters the logging level.
func (lvl AtomicLevel) SetLevel(l Level) {
	atomic.StoreInt32(lvl.l, int32(l))
}

// String returns the string representation of the underlying Level.
func (lvl AtomicLevel) String() string {
	return lvl.Level().String()
}

// UnmarshalText unmarshals the text to an AtomicLevel. It uses the same text
// representations as the static Levels ("debug", "info", "warn",
// "error", "dpanic", "panic", and "fatal").
func (lvl *AtomicLevel) UnmarshalText(text []byte) error {
	zero := int32(0)
	if lvl.l == nil {
		lvl.l = &zero
	}

	var l Level
	if err := l.UnmarshalText(text); err != nil {
		return err
	}

	lvl.SetLevel(l)
	return nil
}

// MarshalText marshals the AtomicLevel to a byte slice. It uses the same
// text representation as the static Levels ("debug", "info", "warn",
// "error", "dpanic", "panic", and "fatal").
func (lvl AtomicLevel) MarshalText() (text []byte, err error) {
	return lvl.Level().MarshalText()
}

func (lvl *AtomicLevel) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		if err := lvl.UnmarshalText(uns.StringToBytes(value.Value)); err != nil {
			return yamlTypeError(value, reflect.TypeOf(lvl).Elem(), err)
		}
		return nil
	}
	return yamlTypeError(value, reflect.TypeOf(lvl).Elem(), nil)
}
