package config

import (
	"errors"
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/def"
	"gitee.com/sy_183/common/log"
)

const DefaultTimeLayout = "2006-01-02 15:04:05.999999999"

var (
	LogConfigNotFoundError  = errors.New("log config not found")
	InvalidLogEncodingError = errors.New("invalid log encoding")
)

func newTrue() *bool {
	b := true
	return &b
}

type JsonEncoder struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey      string `json:"message-key" yaml:"message-key" default:"msg"`
	LevelKey        string `json:"level-key" yaml:"level-key" default:"level"`
	TimeKey         string `json:"time-key" yaml:"time-key" default:"time"`
	NameKey         string `json:"name-key" yaml:"name-key" default:"logger"`
	CallerKey       string `json:"caller-key" yaml:"caller-key" default:"caller"`
	FunctionKey     string `json:"function-key" yaml:"function-key" default:"func"`
	StacktraceKey   string `json:"stacktrace-key" yaml:"stacktrace-key"`
	SkipLineEndingP *bool  `json:"skip-line-ending" yaml:"skip-line-ending"`
	LineEnding      string `json:"line-ending" yaml:"line-ending"`
	EscapeESCP      *bool  `json:"escape-esc" yaml:"escape-esc"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    log.LevelEncoder    `json:"level-encoder" yaml:"level-encoder"`
	EncodeTime     log.TimeEncoder     `json:"time-encoder" yaml:"time-encoder"`
	EncodeDuration log.DurationEncoder `json:"duration-encoder" yaml:"duration-encoder"`
	EncodeCaller   log.CallerEncoder   `json:"caller-encoder" yaml:"caller-encoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName log.NameEncoder `json:"name-encoder" yaml:"name-encoder"`
}

func (e *JsonEncoder) PreModify() (nc any, modified bool) {
	e.EncodeLevel = log.CapitalLevelEncoder
	e.EncodeTime = log.TimeEncoderOfLayout(DefaultTimeLayout)
	e.EncodeDuration = log.StringDurationEncoder
	e.EncodeCaller = log.ShortCallerEncoder
	e.EncodeName = log.FullNameEncoder
	return e, true
}

func (e *JsonEncoder) SkipLineEnding() bool {
	if e.SkipLineEndingP == nil {
		return false
	}
	return *e.SkipLineEndingP
}

func (e *JsonEncoder) EscapeESC() bool {
	if e.EscapeESCP == nil {
		return false
	}
	return *e.EscapeESCP
}

type ConsoleEncoder struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	DisableLevelP      *bool  `json:"disable-level" yaml:"disable-level"`
	DisableTimeP       *bool  `json:"disable-time" yaml:"disable-time"`
	DisableNameP       *bool  `json:"disable-name" yaml:"disable-name"`
	DisableCallerP     *bool  `json:"disable-caller" yaml:"disable-caller" default:"true"`
	DisableFunctionP   *bool  `json:"disable-function" yaml:"disable-function" default:"true"`
	DisableStacktraceP *bool  `json:"disable-stacktrace" yaml:"disable-stacktrace" default:"true"`
	SkipLineEndingP    *bool  `json:"skip-line-ending" yaml:"skip-line-ending"`
	LineEnding         string `json:"line-ending" yaml:"line-ending"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    log.LevelEncoder    `json:"level-encoder" yaml:"level-encoder"`
	EncodeTime     log.TimeEncoder     `json:"time-encoder" yaml:"time-encoder"`
	EncodeDuration log.DurationEncoder `json:"duration-encoder" yaml:"duration-encoder"`
	EncodeCaller   log.CallerEncoder   `json:"caller-encoder" yaml:"caller-encoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName log.NameEncoder `json:"name-encoder" yaml:"name-encoder"`
	// Configures the field separator used by the console encoder. Defaults
	// to tab.
	ConsoleSeparator string `json:"console-separator" yaml:"console-separator"`
}

func (e *ConsoleEncoder) PreModify() (nc any, modified bool) {
	e.EncodeLevel = log.CapitalLevelEncoder
	e.EncodeTime = log.TimeEncoderOfLayout(DefaultTimeLayout)
	e.EncodeDuration = log.StringDurationEncoder
	e.EncodeCaller = log.ShortCallerEncoder
	e.EncodeName = log.FullNameEncoder
	return e, false
}

func (e *ConsoleEncoder) DisableLevel() bool {
	if e.DisableLevelP == nil {
		return false
	}
	return *e.DisableLevelP
}

func (e *ConsoleEncoder) DisableTime() bool {
	if e.DisableTimeP == nil {
		return false
	}
	return *e.DisableTimeP
}

func (e *ConsoleEncoder) DisableName() bool {
	if e.DisableNameP == nil {
		return false
	}
	return *e.DisableNameP
}

func (e *ConsoleEncoder) DisableCaller() bool {
	if e.DisableCallerP == nil {
		return false
	}
	return *e.DisableCallerP
}

func (e *ConsoleEncoder) DisableFunction() bool {
	if e.DisableFunctionP == nil {
		return false
	}
	return *e.DisableFunctionP
}

func (e *ConsoleEncoder) DisableStacktrace() bool {
	if e.DisableStacktraceP == nil {
		return false
	}
	return *e.DisableStacktraceP
}

func (e *ConsoleEncoder) SkipLineEnding() bool {
	if e.SkipLineEndingP == nil {
		return false
	}
	return *e.SkipLineEndingP
}

type LoggerConfig struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level log.AtomicLevel `json:"level" yaml:"level"`
	// DevelopmentP puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	DevelopmentP *bool `json:"development" yaml:"development"`
	// DisableCallerP stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCallerP *bool `json:"disable-caller" yaml:"disable-caller"`
	// DisableStacktraceP completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktraceP *bool `json:"disable-stacktrace" yaml:"disable-stacktrace" default:"true"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *log.SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding" default:"console"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"output-paths" yaml:"output-paths" default:"[stdout]"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"error-output-paths" yaml:"error-output-paths" default:"[stderr]"`
	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initial-fields" yaml:"initial-fields"`

	ConsoleEncoder ConsoleEncoder `yaml:"console-encoder" json:"console-encoder"`
	JsonEncoder    JsonEncoder    `yaml:"json-encoder" json:"json-encoder"`
}

func (c *LoggerConfig) PreModify() (nc any, modified bool) {
	c.Level = log.NewAtomicLevelAt(log.InfoLevel)
	return c, true
}

func (c *LoggerConfig) Development() bool {
	if c.DevelopmentP == nil {
		return false
	}
	return *c.DevelopmentP
}

func (c *LoggerConfig) DisableCaller() bool {
	if c.DisableCallerP == nil {
		return false
	}
	return *c.DisableCallerP
}

func (c *LoggerConfig) DisableStacktrace() bool {
	if c.DisableStacktraceP == nil {
		return false
	}
	return *c.DisableStacktraceP
}

func (c *LoggerConfig) Build() (*log.Logger, error) {
	var encoder log.Encoder
	switch c.Encoding {
	case "json":
		jsonConfig := c.JsonEncoder
		encoder = log.NewJSONEncoder(log.JsonEncoderConfig{
			MessageKey:     jsonConfig.MessageKey,
			LevelKey:       jsonConfig.LevelKey,
			TimeKey:        jsonConfig.TimeKey,
			NameKey:        jsonConfig.NameKey,
			CallerKey:      jsonConfig.CallerKey,
			FunctionKey:    jsonConfig.FunctionKey,
			StacktraceKey:  jsonConfig.StacktraceKey,
			SkipLineEnding: jsonConfig.SkipLineEnding(),
			LineEnding:     jsonConfig.LineEnding,
			EscapeESC:      jsonConfig.EscapeESC(),
			EncodeLevel:    jsonConfig.EncodeLevel,
			EncodeTime:     jsonConfig.EncodeTime,
			EncodeDuration: jsonConfig.EncodeDuration,
			EncodeCaller:   jsonConfig.EncodeCaller,
			EncodeName:     jsonConfig.EncodeName,
		})
	case "console":
		consoleConfig := c.ConsoleEncoder
		encoder = log.NewConsoleEncoder(log.ConsoleEncoderConfig{
			DisableLevel:      consoleConfig.DisableLevel(),
			DisableTime:       consoleConfig.DisableTime(),
			DisableName:       consoleConfig.DisableName(),
			DisableCaller:     consoleConfig.DisableCaller(),
			DisableFunction:   consoleConfig.DisableFunction(),
			DisableStacktrace: consoleConfig.DisableStacktrace(),
			SkipLineEnding:    consoleConfig.SkipLineEnding(),
			LineEnding:        consoleConfig.LineEnding,
			EncodeLevel:       consoleConfig.EncodeLevel,
			EncodeTime:        consoleConfig.EncodeTime,
			EncodeDuration:    consoleConfig.EncodeDuration,
			EncodeCaller:      consoleConfig.EncodeCaller,
			EncodeName:        consoleConfig.EncodeName,
			ConsoleSeparator:  consoleConfig.ConsoleSeparator,
		})
	default:
		panic("internal error: invalid log encoding")
	}

	logConfig := &log.Config{
		Level:             c.Level,
		Development:       c.Development(),
		DisableCaller:     c.DisableCaller(),
		DisableStacktrace: c.DisableStacktrace(),
		Sampling:          c.Sampling,
		Encoder:           encoder,
		OutputPaths:       c.OutputPaths,
		ErrorOutputPaths:  c.ErrorOutputPaths,
	}

	if len(c.InitialFields) > 0 {
		logConfig.InitialFields = make(map[string]interface{}, len(c.InitialFields))
		for key, value := range c.InitialFields {
			logConfig.InitialFields[key] = value
		}
	}

	return logConfig.Build()
}

func (c *LoggerConfig) MustBuild() *log.Logger {
	return assert.Must(c.Build())
}

type Config struct {
	LoggerConfig `yaml:",inline"`
	Configs      map[string]*LoggerConfig `yaml:",inline" json:"configs" default:"{}"`
	Modules      map[string]string        `yaml:"modules" json:"modules" default:"{}"`
}

func (c *Config) PostHandle() (nc any, modified bool, err error) {
	if c.Encoding != "json" && c.Encoding != "console" {
		return nil, false, InvalidLogEncodingError
	}

	for _, logConfig := range c.Configs {
		// post handle custom logger config
		def.SetDefaultP(&logConfig.Level, c.Level)
		def.SetDefaultP(&logConfig.DevelopmentP, c.DevelopmentP)
		def.SetDefaultP(&logConfig.DisableCallerP, c.DisableCallerP)
		def.SetDefaultP(&logConfig.DisableStacktraceP, c.DisableStacktraceP)
		def.SetDefaultP(&logConfig.Sampling, c.Sampling)
		def.SetDefaultP(&logConfig.Encoding, c.Encoding)
		if len(logConfig.OutputPaths) == 0 {
			logConfig.OutputPaths = c.OutputPaths
		}
		if len(logConfig.ErrorOutputPaths) == 0 {
			logConfig.ErrorOutputPaths = c.ErrorOutputPaths
		}
		if len(logConfig.InitialFields) == 0 && len(c.InitialFields) != 0 {
			if logConfig.InitialFields == nil {
				logConfig.InitialFields = make(map[string]interface{}, len(c.InitialFields))
			}
			for key, value := range c.InitialFields {
				logConfig.InitialFields[key] = value
			}
		}

		// post handle custom logger json encoder config
		jsonEncoder := &logConfig.JsonEncoder
		def.SetDefaultP(&jsonEncoder.MessageKey, c.JsonEncoder.MessageKey)
		def.SetDefaultP(&jsonEncoder.TimeKey, c.JsonEncoder.TimeKey)
		def.SetDefaultP(&jsonEncoder.NameKey, c.JsonEncoder.NameKey)
		def.SetDefaultP(&jsonEncoder.CallerKey, c.JsonEncoder.CallerKey)
		def.SetDefaultP(&jsonEncoder.FunctionKey, c.JsonEncoder.FunctionKey)
		def.SetDefaultP(&jsonEncoder.StacktraceKey, c.JsonEncoder.StacktraceKey)
		def.SetDefaultP(&jsonEncoder.LineEnding, c.JsonEncoder.LineEnding)
		def.SetDefaultP(&jsonEncoder.SkipLineEndingP, c.JsonEncoder.SkipLineEndingP)

		def.SetAnyP(&jsonEncoder.EncodeLevel, c.JsonEncoder.EncodeLevel)
		def.SetAnyP(&jsonEncoder.EncodeTime, c.JsonEncoder.EncodeTime)
		def.SetAnyP(&jsonEncoder.EncodeDuration, c.JsonEncoder.EncodeDuration)
		def.SetAnyP(&jsonEncoder.EncodeCaller, c.JsonEncoder.EncodeCaller)
		def.SetAnyP(&jsonEncoder.EncodeName, c.JsonEncoder.EncodeName)

		// post handle custom logger console encoder config
		consoleEncoder := &logConfig.ConsoleEncoder
		def.SetDefaultP(&consoleEncoder.DisableLevelP, c.ConsoleEncoder.DisableLevelP)
		def.SetDefaultP(&consoleEncoder.DisableTimeP, c.ConsoleEncoder.DisableTimeP)
		def.SetDefaultP(&consoleEncoder.DisableNameP, c.ConsoleEncoder.DisableNameP)
		def.SetDefaultP(&consoleEncoder.DisableCallerP, c.ConsoleEncoder.DisableCallerP)
		def.SetDefaultP(&consoleEncoder.DisableFunctionP, c.ConsoleEncoder.DisableFunctionP)
		def.SetDefaultP(&consoleEncoder.DisableStacktraceP, c.ConsoleEncoder.DisableStacktraceP)
		def.SetDefaultP(&consoleEncoder.SkipLineEndingP, c.ConsoleEncoder.SkipLineEndingP)

		def.SetAnyP(&consoleEncoder.EncodeLevel, c.ConsoleEncoder.EncodeLevel)
		def.SetAnyP(&consoleEncoder.EncodeTime, c.ConsoleEncoder.EncodeTime)
		def.SetAnyP(&consoleEncoder.EncodeDuration, c.ConsoleEncoder.EncodeDuration)
		def.SetAnyP(&consoleEncoder.EncodeCaller, c.ConsoleEncoder.EncodeCaller)
		def.SetAnyP(&consoleEncoder.EncodeName, c.ConsoleEncoder.EncodeName)
	}

	// check module logger config
	for _, conf := range c.Modules {
		if _, has := c.Configs[conf]; !has {
			if conf != "default" {
				return nil, false, InvalidLogEncodingError
			}
		}
	}
	return c, true, nil
}

type Module interface {
	Module() string
}

func (c *Config) Build(modules ...string) (*log.Logger, error) {
	var conf *LoggerConfig
	for _, module := range modules {
		if c, has := c.Configs[module]; has {
			conf = c
			break
		}
	}
	// logger config not found, use default
	if conf == nil {
		conf = &c.LoggerConfig
	}

	return conf.Build()
}

func (c *Config) MustBuild(modules ...string) *log.Logger {
	return assert.Must(c.Build(modules...))
}
