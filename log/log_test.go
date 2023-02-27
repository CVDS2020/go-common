package log

import (
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/sgr"
	"os"
	"sync/atomic"
	"testing"
)

func TestLog(t *testing.T) {
	//config := zap.Config{
	//	Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
	//	Development: true,
	//	Sampling: &zap.SamplingConfig{
	//		Initial:    100,
	//		Thereafter: 100,
	//	},
	//	DisableStacktrace: true,
	//	Encoding:          "console",
	//	EncoderConfig: zapcore.EncoderConfig{
	//		MessageKey:     "msg",
	//		LevelKey:       "level",
	//		TimeKey:        "time",
	//		NameKey:        "name",
	//		CallerKey:      "caller",
	//		FunctionKey:    "func",
	//		LineEnding:     zapcore.DefaultLineEnding,
	//		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	//		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999999"),
	//		EncodeDuration: zapcore.StringDurationEncoder,
	//		EncodeCaller:   zapcore.ShortCallerEncoder,
	//		EncodeName:     zapcore.FullNameEncoder,
	//	},
	//	OutputPaths:      []string{"stdout"},
	//	ErrorOutputPaths: []string{"stdout"},
	//}
	//logger := assert.Must(config.Build())
	//zap.NewExample()
	encoder := NewConsoleEncoder(ConsoleEncoderConfig{
		//DisableStacktrace: true,
		EncodeLevel:    CapitalColorLevelEncoder,
		EncodeTime:     TimeEncoderOfLayout("2006-01-02 15:04:05.999999999"),
		EncodeDuration: StringDurationEncoder,
		EncodeCaller:   ShortCallerEncoder,
		EncodeName:     FullNameEncoder,
	})
	logger := New(NewCore(encoder, os.Stdout, DebugLevel), AddCaller(), AddStacktrace(InfoLevel), Fields(String("module", "log")))
	logger.Debug(sgr.WrapRGB24("lili", 0x2DF3D8), String("hello", "hi"))
	logger.Info("hellofdfdddddddddddddddddddddd", String("hello", sgr.WrapColor("dds", sgr.FgBrightBlue)))
	logger.Warn("hello")
	logger.Error("hello")
}

func TestLog1(t *testing.T) {
	logger := assert.Must(Config{
		Level:             NewAtomicLevelAt(InfoLevel),
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoder: NewConsoleEncoder(ConsoleEncoderConfig{
			DisableCaller:     false,
			DisableStacktrace: false,
		}),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build())
	logger.Info("hello")
}

func TestCopy(t *testing.T) {
	buf1 := make([]byte, 1024)
	buf2 := make([]byte, 1024)
	for i := 0; i < 40000000; i++ {
		copy(buf2, buf1)
		copy(buf1, buf2)
	}
}

func TestAtomic(t *testing.T) {
	var i int64
	for j := 0; j < 200000000; j++ {
		atomic.AddInt64(&i, int64(j))
		atomic.AddInt64(&i, int64(j))
	}
}
