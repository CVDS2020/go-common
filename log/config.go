package log

import (
	"sort"
	"time"
)

// SamplingConfig sets a sampling strategy for the logger. Sampling caps the
// global CPU and I/O load that logging puts on your process while attempting
// to preserve a representative subset of your logs.
//
// If specified, the Sampler will invoke the Hook after each decision.
//
// Values configured here are per-second. See NewSamplerWithOptions for
// details.
type SamplingConfig struct {
	Initial    int                           `json:"initial" yaml:"initial"`
	Thereafter int                           `json:"thereafter" yaml:"thereafter"`
	Hook       func(Entry, SamplingDecision) `json:"-" yaml:"-"`
}

// Config offers a declarative way to construct a logger. It doesn't do
// anything that can't be done with New, Options, and the various
// WriteSyncer and Core wrappers, but it's a simpler way to
// toggle common options.
//
// Note that Config intentionally supports only the most common options. More
// unusual logging setups (logging to network connections or message queues,
// splitting output between multiple files, etc.) are possible, but require
// direct use of the core package. For sample code, see the package-level
// BasicConfiguration and AdvancedConfiguration examples.
//
// For an example showing runtime log level changes, see the documentation for
// AtomicLevel.
type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disable-caller" yaml:"disable-caller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disable-stacktrace" yaml:"disable-stacktrace"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoder sets the logger's encoder, See Encoder for details
	Encoder Encoder `json:"-" yaml:"-"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"output-paths" yaml:"output-paths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"error-output-paths" yaml:"error-output-paths"`
	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initial-fields" yaml:"initial-fields"`
}

// NewProductionConfig is a reasonable production logging configuration.
// Logging is enabled at InfoLevel and above.
//
// It uses a JSON encoder, writes to standard error, and enables sampling.
// Stacktraces are automatically included on logs of ErrorLevel and above.
func NewProductionConfig() Config {
	return Config{
		Level:       NewAtomicLevelAt(InfoLevel),
		Development: false,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoder:          NewJSONEncoder(NewProductionJsonEncoderConfig()),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// NewDevelopmentConfig is a reasonable development logging configuration.
// Logging is enabled at DebugLevel and above.
//
// It enables development mode (which makes DPanicLevel logs panic), uses a
// console encoder, writes to standard error, and disables sampling.
// Stacktraces are automatically included on logs of WarnLevel and above.
func NewDevelopmentConfig() Config {
	return Config{
		Level:            NewAtomicLevelAt(DebugLevel),
		Development:      true,
		Encoder:          NewJSONEncoder(NewDevelopmentJsonEncoderConfig()),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// Build constructs a logger from the Config and Options.
func (cfg Config) Build(opts ...Option) (*Logger, error) {
	if cfg.Level.l == nil {
		cfg.Level = NewAtomicLevel()
	}
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stdout"}
	}
	if len(cfg.ErrorOutputPaths) == 0 {
		cfg.ErrorOutputPaths = []string{"stderr"}
	}
	if cfg.Encoder == nil {
		cfg.Encoder = NewConsoleEncoder(ConsoleEncoderConfig{})
	}
	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	log := New(
		NewCore(cfg.Encoder, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

func (cfg Config) buildOptions(errSink WriteSyncer) []Option {
	opts := []Option{ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, AddCaller())
	}

	stackLevel := ErrorLevel
	if cfg.Development {
		stackLevel = WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, AddStacktrace(stackLevel))
	}

	if scfg := cfg.Sampling; scfg != nil {
		opts = append(opts, WrapCore(func(core Core) Core {
			var samplerOpts []SamplerOption
			if scfg.Hook != nil {
				samplerOpts = append(samplerOpts, SamplerHook(scfg.Hook))
			}
			return NewSamplerWithOptions(
				core,
				time.Second,
				cfg.Sampling.Initial,
				cfg.Sampling.Thereafter,
				samplerOpts...,
			)
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, Fields(fs...))
	}

	return opts
}

func (cfg Config) openSinks() (WriteSyncer, WriteSyncer, error) {
	sink, closeOut, err := Open(cfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}
	errSink, _, err := Open(cfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}
	return sink, errSink, nil
}
