package sgr

import "strings"

const (
	ResetCode      = "0"
	BoldCode       = "1"
	FaintCode      = "2"
	ItalicCode     = "3"
	UnderlineCode  = "4"
	SlowBlinkCode  = "5"
	RapidBlinkCode = "6"
	ReversedCode   = "7"
	ConcealCode    = "8"
	CrossedOutCode = "9"

	PrimaryFontCode      = "10"
	AlternativeFont1Code = "11"
	AlternativeFont2Code = "12"
	AlternativeFont3Code = "13"
	AlternativeFont4Code = "14"
	AlternativeFont5Code = "15"
	AlternativeFont6Code = "16"
	AlternativeFont7Code = "17"
	AlternativeFont8Code = "18"
	AlternativeFont9Code = "19"

	FrakturCode              = "20"
	DoublyUnderlinedCode     = "21"
	NormalIntensityCode      = "22"
	NotItalicBlackLetterCode = "23"
	NotUnderlinedCode        = "24"
	NotBlinkingCode          = "25"
	ProportionalSpacingCode  = "26"
	NotReversedCode          = "27"
	RevealCode               = "28"
	NotCrossedOutCode        = "29"

	FgBlackCode        = "30"
	FgRedCode          = "31"
	FgGreenCode        = "32"
	FgYellowCode       = "33"
	FgBlueCode         = "34"
	FgMagentaCode      = "35"
	FgCyanCode         = "36"
	FgWhiteCode        = "37"
	CustomFgColorCode  = "38"
	DefaultFgColorCode = "39"

	BgBlackCode        = "40"
	BgRedCode          = "41"
	BgGreenCode        = "42"
	BgYellowCode       = "43"
	BgBlueCode         = "44"
	BgMagentaCode      = "45"
	BgCyanCode         = "46"
	BgWhiteCode        = "47"
	CustomBgColorCode  = "48"
	DefaultBgColorCode = "49"

	DisableProportionalSpacingCode = "50"
	FramedCode                     = "51"
	EncircledCode                  = "52"
	OverlinedCode                  = "53"
	NotFramedEncircledCode         = "54"
	NotOverlinedCode               = "55"
	CustomUnderlineColorCode       = "58"
	DefaultUnderlineColorCode      = "59"

	IdeogramUnderlineCode       = "60"
	IdeogramDoubleUnderlineCode = "61"
	IdeogramOverlineCode        = "62"
	IdeogramDoubleOverlineCode  = "63"
	IdeogramStressMarkingCode   = "64"
	NoIdeogramAttributesCode    = "65"

	SuperscriptCode             = "73"
	SubscriptCode               = "74"
	NotSuperscriptSubscriptCode = "75"

	FgBrightBlackCode   = "90"
	FgBrightRedCode     = "91"
	FgBrightGreenCode   = "92"
	FgBrightYellowCode  = "93"
	FgBrightBlueCode    = "94"
	FgBrightMagentaCode = "95"
	FgBrightCyanCode    = "96"
	FgBrightWhiteCode   = "97"

	BgBrightBlackCode   = "100"
	BgBrightRedCode     = "101"
	BgBrightGreenCode   = "102"
	BgBrightYellowCode  = "103"
	BgBrightBlueCode    = "104"
	BgBrightMagentaCode = "105"
	BgBrightCyanCode    = "106"
	BgBrightWhiteCode   = "107"
)

func WrapCode(s string, c string) string {
	return SGR_CSIHead + c + SGR_CSIEnd + s + ResetSGR_CSI
}

func WrapCodes(s string, cs ...string) string {
	switch len(cs) {
	case 0:
		return ""
	case 1:
		return WrapCode(s, cs[0])
	}
	n := len(SGR_CSIHead) + len(SGR_CSIEnd) + len(cs) + len(s) + len(ResetSGR_CSI) - 1
	for i := 0; i < len(cs); i++ {
		n += len(cs[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(SGR_CSIHead)
	b.WriteString(cs[0])
	for _, c := range cs[1:] {
		b.WriteByte(';')
		b.WriteString(c)
	}
	b.WriteString(SGR_CSIEnd)
	b.WriteString(s)
	b.WriteString(ResetSGR_CSI)
	return b.String()
}
