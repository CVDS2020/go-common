package log

import (
	"sync/atomic"
	"time"
)

const (
	_numLevels        = _maxLevel - _minLevel + 1
	_countersPerLevel = 4096
)

type counter struct {
	resetAt int64
	counter uint64
}

type counters [_numLevels][_countersPerLevel]counter

func newCounters() *counters {
	return &counters{}
}

func (cs *counters) get(lvl Level, key string) *counter {
	i := lvl - _minLevel
	j := fnv32a(key) % _countersPerLevel
	return &cs[i][j]
}

// fnv32a, adapted from "hash/fnv", but without a []byte(string) alloc
func fnv32a(s string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	hash := uint32(offset32)
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}

func (c *counter) IncCheckReset(t time.Time, tick time.Duration) uint64 {
	tn := t.UnixNano()
	resetAfter := atomic.LoadInt64(&c.resetAt)
	if resetAfter > tn {
		return atomic.AddUint64(&c.counter, 1)
	}

	atomic.StoreUint64(&c.counter, 1)

	newResetAfter := tn + tick.Nanoseconds()
	if !atomic.CompareAndSwapInt64(&c.resetAt, resetAfter, newResetAfter) {
		// We raced with another goroutine trying to reset, and it also reset
		// the counter to 1, so we need to reincrement the counter.
		return atomic.AddUint64(&c.counter, 1)
	}

	return 1
}

// SamplingDecision is a decision represented as a bit field made by sampler.
// More decisions may be added in the future.
type SamplingDecision uint32

const (
	// LogDropped indicates that the Sampler dropped a log entry.
	LogDropped SamplingDecision = 1 << iota
	// LogSampled indicates that the Sampler sampled a log entry.
	LogSampled
)

// SamplerOptionFunc wraps a func so it satisfies the SamplerOption interface.
type SamplerOptionFunc func(*sampler)

func (f SamplerOptionFunc) apply(s *sampler) {
	f(s)
}

// SamplerOption configures a Sampler.
type SamplerOption interface {
	apply(*sampler)
}

// nopSamplingHook is the default hook used by sampler.
func nopSamplingHook(Entry, SamplingDecision) {}

// SamplerHook registers a function  which will be called when Sampler makes a
// decision.
//
// This hook may be used to get visibility into the performance of the sampler.
// For example, use it to track metrics of dropped versus sampled logs.
//
//  var dropped atomic.Int64
//  SamplerHook(func(ent Entry, dec SamplingDecision) {
//    if dec&LogDropped > 0 {
//      dropped.Inc()
//    }
//  })
func SamplerHook(hook func(entry Entry, dec SamplingDecision)) SamplerOption {
	return SamplerOptionFunc(func(s *sampler) {
		s.hook = hook
	})
}

// NewSamplerWithOptions creates a Core that samples incoming entries, which
// caps the CPU and I/O load of logging while attempting to preserve a
// representative subset of your logs.
//
// Zap samples by logging the first N entries with a given level and message
// each tick. If more Entries with the same level and message are seen during
// the same interval, every Mth message is logged and the rest are dropped.
//
// For example,
//
//   core = NewSamplerWithOptions(core, time.Second, 10, 5)
//
// This will log the first 10 log entries with the same level and message
// in a one second interval as-is. Following that, it will allow through
// every 5th log entry with the same level and message in that interval.
//
// If thereafter is zero, the Core will drop all log entries after the first N
// in that interval.
//
// Sampler can be configured to report sampling decisions with the SamplerHook
// option.
//
// Keep in mind that Zap's sampling implementation is optimized for speed over
// absolute precision; under load, each tick may be slightly over- or
// under-sampled.
func NewSamplerWithOptions(core Core, tick time.Duration, first, thereafter int, opts ...SamplerOption) Core {
	s := &sampler{
		Core:       core,
		tick:       tick,
		counts:     newCounters(),
		first:      uint64(first),
		thereafter: uint64(thereafter),
		hook:       nopSamplingHook,
	}
	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

type sampler struct {
	Core

	counts            *counters
	tick              time.Duration
	first, thereafter uint64
	hook              func(Entry, SamplingDecision)
}

// NewSampler creates a Core that samples incoming entries, which
// caps the CPU and I/O load of logging while attempting to preserve a
// representative subset of your logs.
//
// Zap samples by logging the first N entries with a given level and message
// each tick. If more Entries with the same level and message are seen during
// the same interval, every Mth message is logged and the rest are dropped.
//
// Keep in mind that log's sampling implementation is optimized for speed over
// absolute precision; under load, each tick may be slightly over- or
// under-sampled.
//
// Deprecated: use NewSamplerWithOptions.
func NewSampler(core Core, tick time.Duration, first, thereafter int) Core {
	return NewSamplerWithOptions(core, tick, first, thereafter)
}

func (s *sampler) With(fields []Field) Core {
	return &sampler{
		Core:       s.Core.With(fields),
		tick:       s.tick,
		counts:     s.counts,
		first:      s.first,
		thereafter: s.thereafter,
		hook:       s.hook,
	}
}

func (s *sampler) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	if !s.Enabled(ent.Level) {
		return ce
	}

	if ent.Level >= _minLevel && ent.Level <= _maxLevel {
		counter := s.counts.get(ent.Level, ent.Message)
		n := counter.IncCheckReset(ent.Time, s.tick)
		if n > s.first && (s.thereafter == 0 || (n-s.first)%s.thereafter != 0) {
			s.hook(ent, LogDropped)
			return ce
		}
		s.hook(ent, LogSampled)
	}
	return s.Core.Check(ent, ce)
}
