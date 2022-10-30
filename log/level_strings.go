package log

import (
	"gitee.com/sy_183/common/sgr"
)

var (
	_levelToColor = map[Level]string{
		DebugLevel:  sgr.FgMagentaCSI,
		InfoLevel:   sgr.FgBlueCSI,
		WarnLevel:   sgr.FgYellowCSI,
		ErrorLevel:  sgr.FgRedCSI,
		DPanicLevel: sgr.FgRedCSI,
		PanicLevel:  sgr.FgRedCSI,
		FatalLevel:  sgr.FgRedCSI,
	}
	_unknownLevelColor = sgr.DefaultFgColorCSI

	_levelToLowercaseColorString = make(map[Level]string, len(_levelToColor))
	_levelToCapitalColorString   = make(map[Level]string, len(_levelToColor))
)

func init() {
	for level, _color := range _levelToColor {
		_levelToLowercaseColorString[level] = sgr.WrapCSI(level.String(), _color)
		_levelToCapitalColorString[level] = sgr.WrapCSI(level.CapitalString(), _color)
	}
}
