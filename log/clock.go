package log

import "time"

// DefaultClock is the default clock used by Zap in operations that require
// time. This clock uses the system clock for all operations.
var DefaultClock = systemClock{}

// Clock is a source of time for logged entries.
type Clock interface {
	// Now returns the current local time.
	Now() time.Time

	// NewTicker returns *time.Ticker that holds a channel
	// that delivers "ticks" of a clock.
	NewTicker(time.Duration) *time.Ticker
}

// systemClock implements default Clock that uses system time.
type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func (systemClock) NewTicker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}
