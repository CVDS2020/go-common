package sgr

import "strings"

const (
	SGR_CSIHead = "\x1b["
	SGR_CSIEnd  = "m"
)

const (
	ResetSGR_CSI = SGR_CSIHead + ResetCode + SGR_CSIEnd

	FgBlackCSI        = SGR_CSIHead + FgBlackCode + SGR_CSIEnd
	FgRedCSI          = SGR_CSIHead + FgRedCode + SGR_CSIEnd
	FgGreenCSI        = SGR_CSIHead + FgGreenCode + SGR_CSIEnd
	FgYellowCSI       = SGR_CSIHead + FgYellowCode + SGR_CSIEnd
	FgBlueCSI         = SGR_CSIHead + FgBlueCode + SGR_CSIEnd
	FgMagentaCSI      = SGR_CSIHead + FgMagentaCode + SGR_CSIEnd
	FgCyanCSI         = SGR_CSIHead + FgCyanCode + SGR_CSIEnd
	FgWhiteCSI        = SGR_CSIHead + FgWhiteCode + SGR_CSIEnd
	CustomFgColorCSI  = SGR_CSIHead + CustomFgColorCode + SGR_CSIEnd
	DefaultFgColorCSI = SGR_CSIHead + DefaultFgColorCode + SGR_CSIEnd

	BgBlackCSI        = SGR_CSIHead + BgBlackCode + SGR_CSIEnd
	BgRedCSI          = SGR_CSIHead + BgRedCode + SGR_CSIEnd
	BgGreenCSI        = SGR_CSIHead + BgGreenCode + SGR_CSIEnd
	BgYellowCSI       = SGR_CSIHead + BgYellowCode + SGR_CSIEnd
	BgBlueCSI         = SGR_CSIHead + BgBlueCode + SGR_CSIEnd
	BgMagentaCSI      = SGR_CSIHead + BgMagentaCode + SGR_CSIEnd
	BgCyanCSI         = SGR_CSIHead + BgCyanCode + SGR_CSIEnd
	BgWhiteCSI        = SGR_CSIHead + BgWhiteCode + SGR_CSIEnd
	CustomBgColorCSI  = SGR_CSIHead + CustomBgColorCode + SGR_CSIEnd
	DefaultBgColorCSI = SGR_CSIHead + DefaultBgColorCode + SGR_CSIEnd

	FgBrightBlackCSI   = SGR_CSIHead + FgBrightBlackCode + SGR_CSIEnd
	FgBrightRedCSI     = SGR_CSIHead + FgBrightRedCode + SGR_CSIEnd
	FgBrightGreenCSI   = SGR_CSIHead + FgBrightGreenCode + SGR_CSIEnd
	FgBrightYellowCSI  = SGR_CSIHead + FgBrightYellowCode + SGR_CSIEnd
	FgBrightBlueCSI    = SGR_CSIHead + FgBrightBlueCode + SGR_CSIEnd
	FgBrightMagentaCSI = SGR_CSIHead + FgBrightMagentaCode + SGR_CSIEnd
	FgBrightCyanCSI    = SGR_CSIHead + FgBrightCyanCode + SGR_CSIEnd
	FgBrightWhiteCSI   = SGR_CSIHead + FgBrightWhiteCode + SGR_CSIEnd

	BgBrightBlackCSI   = SGR_CSIHead + BgBrightBlackCode + SGR_CSIEnd
	BgBrightRedCSI     = SGR_CSIHead + BgBrightRedCode + SGR_CSIEnd
	BgBrightGreenCSI   = SGR_CSIHead + BgBrightGreenCode + SGR_CSIEnd
	BgBrightYellowCSI  = SGR_CSIHead + BgBrightYellowCode + SGR_CSIEnd
	BgBrightBlueCSI    = SGR_CSIHead + BgBrightBlueCode + SGR_CSIEnd
	BgBrightMagentaCSI = SGR_CSIHead + BgBrightMagentaCode + SGR_CSIEnd
	BgBrightCyanCSI    = SGR_CSIHead + BgBrightCyanCode + SGR_CSIEnd
	BgBrightWhiteCSI   = SGR_CSIHead + BgBrightWhiteCode + SGR_CSIEnd
)

func WrapCSI(s string, csi string) string {
	return csi + s + ResetSGR_CSI
}

func MakeCSI(codes ...string) string {
	switch len(codes) {
	case 0:
		return ""
	case 1:
		return SGR_CSIHead + codes[0] + SGR_CSIEnd
	}
	n := len(SGR_CSIHead) + len(SGR_CSIEnd) + len(codes) - 1
	for i := 0; i < len(codes); i++ {
		n += len(codes[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(SGR_CSIHead)
	b.WriteString(codes[0])
	for _, c := range codes[1:] {
		b.WriteByte(';')
		b.WriteString(c)
	}
	b.WriteString(SGR_CSIEnd)
	return b.String()
}
