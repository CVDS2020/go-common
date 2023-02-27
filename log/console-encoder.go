package log

import (
	"fmt"
	"gitee.com/sy_183/common/log/internal/bufferpool"
	"io"
	"sync"
)

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type ConsoleEncoderConfig struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	DisableLevel      bool   `json:"disable-level" yaml:"disable-level"`
	DisableTime       bool   `json:"disable-time" yaml:"disable-time"`
	DisableName       bool   `json:"disable-name" yaml:"disable-name"`
	DisableCaller     bool   `json:"disable-caller" yaml:"disable-caller"`
	DisableFunction   bool   `json:"disable-function" yaml:"disable-function"`
	DisableStacktrace bool   `json:"disable-stacktrace" yaml:"disable-stacktrace"`
	SkipLineEnding    bool   `json:"skip-line-ending" yaml:"skip-line-ending"`
	LineEnding        string `json:"line-ending" yaml:"line-ending"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    LevelEncoder    `json:"level-encoder" yaml:"level-encoder"`
	EncodeTime     TimeEncoder     `json:"time-encoder" yaml:"time-encoder"`
	EncodeDuration DurationEncoder `json:"duration-encoder" yaml:"duration-encoder"`
	EncodeCaller   CallerEncoder   `json:"caller-encoder" yaml:"caller-encoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName NameEncoder `json:"name-encoder" yaml:"name-encoder"`
	// Configure the encoder for interface{} type objects.
	// If not provided, objects are encoded using json.Encoder
	NewReflectedEncoder func(io.Writer) ReflectedEncoder `json:"-" yaml:"-"`
	// Configures the field separator used by the console encoder. Defaults
	// to tab.
	ConsoleSeparator string `json:"console-separator" yaml:"console-separator"`
}

type consoleEncoder struct {
	*ConsoleEncoderConfig
	*jsonEncoder
}

// NewConsoleEncoder creates an encoder whose output is designed for human -
// rather than machine - consumption. It serializes the core log entry data
// (message, level, timestamp, etc.) in a plain-text format and leaves the
// structured context as JSON.
//
// Note that although the console encoder doesn't use the keys specified in the
// encoder configuration, it will omit any element whose key is set to the empty
// string.
func NewConsoleEncoder(cfg ConsoleEncoderConfig) Encoder {
	if cfg.SkipLineEnding {
		cfg.LineEnding = ""
	} else if cfg.LineEnding == "" {
		cfg.LineEnding = DefaultLineEnding
	}

	if !cfg.DisableTime && cfg.EncodeTime == nil {
		cfg.EncodeTime = RFC3339NanoTimeEncoder
	}

	if !cfg.DisableLevel && cfg.EncodeLevel == nil {
		cfg.EncodeLevel = CapitalLevelEncoder
	}

	if !cfg.DisableName && cfg.EncodeName == nil {
		cfg.EncodeName = FullNameEncoder
	}

	if !cfg.DisableCaller && cfg.EncodeCaller == nil {
		cfg.EncodeCaller = FullCallerEncoder
	}

	// If no EncoderConfig.NewReflectedEncoder is provided by the user, then use default
	if cfg.NewReflectedEncoder == nil {
		cfg.NewReflectedEncoder = defaultReflectedEncoder
	}

	if cfg.ConsoleSeparator == "" {
		// Use a default delimiter of '\t' for backwards compatibility
		cfg.ConsoleSeparator = "\t"
	}

	return consoleEncoder{
		ConsoleEncoderConfig: &cfg,
		jsonEncoder: newJSONEncoder(JsonEncoderConfig{
			EncodeTime:     cfg.EncodeTime,
			EncodeDuration: cfg.EncodeDuration,
		}, true),
	}
}

func (c consoleEncoder) Clone() Encoder {
	return consoleEncoder{
		ConsoleEncoderConfig: c.ConsoleEncoderConfig,
		jsonEncoder:          c.jsonEncoder.Clone().(*jsonEncoder),
	}
}

func (c consoleEncoder) EncodeEntry(ent Entry, fields []Field) (*bufferpool.Buffer, error) {
	line := bufferpool.Get()

	// We don't want the entry's metadata to be quoted and escaped (if it's
	// encoded as strings), which means that we can't use the JSON encoder. The
	// simplest option is to use the memory encoder and fmt.Fprint.
	//
	// If this ever becomes a performance bottleneck, we can implement
	// ArrayEncoder for our plain-text format.
	arr := getSliceEncoder()
	if !c.DisableTime && c.EncodeTime != nil {
		c.EncodeTime(ent.Time, arr)
	}
	if !c.DisableLevel && c.EncodeLevel != nil {
		c.EncodeLevel(ent.Level, arr)
	}
	if ent.LoggerName != "" && !c.DisableName {
		nameEncoder := c.EncodeName

		if nameEncoder == nil {
			// Fall back to FullNameEncoder for backward compatibility.
			nameEncoder = FullNameEncoder
		}

		nameEncoder(ent.LoggerName, arr)
	}
	if ent.Caller.Defined {
		if !c.DisableCaller && c.EncodeCaller != nil {
			c.EncodeCaller(ent.Caller, arr)
		}
		if !c.DisableFunction {
			arr.AppendString(ent.Caller.Function)
		}
	}
	for i := range arr.elems {
		if i > 0 {
			line.AppendString(c.ConsoleSeparator)
		}
		fmt.Fprint(line, arr.elems[i])
	}
	putSliceEncoder(arr)

	// Add the message itself.
	if ent.Message != "" {
		c.addSeparatorIfNecessary(line)
		line.AppendString(ent.Message)
	}

	// Add any structured context.
	c.writeContext(line, fields)

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && !c.DisableStacktrace {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	line.AppendString(c.LineEnding)
	return line, nil
}

func (c consoleEncoder) writeContext(line *bufferpool.Buffer, extra []Field) {
	context := c.jsonEncoder.Clone().(*jsonEncoder)
	defer func() {
		// putJSONEncoder assumes the buffer is still used, but we write out the buffer so
		// we can free it.
		context.buf.Free()
		putJSONEncoder(context)
	}()

	addFields(context, extra)
	context.closeOpenNamespaces()
	if context.buf.Len() == 0 {
		return
	}

	c.addSeparatorIfNecessary(line)
	line.AppendByte('{')
	line.Write(context.buf.Bytes())
	line.AppendByte('}')
}

func (c consoleEncoder) addSeparatorIfNecessary(line *bufferpool.Buffer) {
	if line.Len() > 0 {
		line.AppendString(c.ConsoleSeparator)
	}
}
