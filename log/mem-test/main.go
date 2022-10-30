package main

import (
	"fmt"
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/log"
	"time"
)

const DefaultTimeLayout = "2006-01-02 15:04:05.999999999"

func main() {
	logger := assert.Must(log.Config{
		Level: log.NewAtomicLevelAt(log.DebugLevel),
		Encoder: log.NewConsoleEncoder(log.ConsoleEncoderConfig{
			DisableCaller:     true,
			DisableFunction:   true,
			DisableStacktrace: true,
			EncodeLevel:       log.CapitalColorLevelEncoder,
			EncodeTime:        log.TimeEncoderOfLayout(DefaultTimeLayout),
			EncodeDuration:    log.SecondsDurationEncoder,
		}),
		OutputPaths: []string{"stdout"},
	}.Build())
	for i := 0; ; i++ {
		logger.Debug(fmt.Sprintf("log out line %d", i))
		time.Sleep(0)
	}
}
