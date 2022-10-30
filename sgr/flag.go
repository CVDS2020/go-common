package sgr

type Flag [2]uint64

var zeroFlag Flag

func Flags(flags ...Flag) (fs Flag) {
	for _, flag := range flags {
		fs = fs.Add(flag)
	}
	return
}

func (f Flag) Add(flag Flag) Flag {
	return Flag{f[0] | flag[0], f[1] | flag[1]}
}

func (f Flag) Remove(flag Flag) Flag {
	return Flag{f[0] & ^flag[0], f[1] & ^flag[1]}
}

func (f Flag) Has(flag Flag) bool {
	return f[0]&flag[0] != 0 || f[1]&flag[1] != 0
}

func (f Flag) HasMulti(flags ...Flag) bool {
	for _, flag := range flags {
		f[0] &= flag[0]
		f[1] &= flag[1]
	}
	return f != zeroFlag
}

func (f Flag) Color() Flag {
	f[1] &= 0xffff000000000000
	f[0] &= 0x000000000000ffff | 1<<17 | 1<<19 | 1<<21
	return f
}

func (f Flag) Meta() *SGRMeta {
	return flagMetaMap[f]
}

func (f Flag) byte(i int) byte {
	if i < 8 {
		return byte((f[1] >> (i * 8)) & 0xff)
	} else if i -= 8; i < 8 {
		return byte((f[0] >> (i * 8)) & 0xff)
	}
	panic("index out of range")
}

var (
	Reset = Flag{0, 1 << 0}

	Bold                 = Flag{0, 1 << 8}
	Faint                = Flag{0, 1 << 9}
	Italic               = Flag{0, 1 << 10}
	Fraktur              = Flag{0, 1 << 11}
	Reversed             = Flag{0, 1 << 12}
	NormalIntensity      = Flag{0, 1 << 13}
	NotItalicBlackLetter = Flag{0, 1 << 14}
	NotReversed          = Flag{0, 1 << 15}

	Underline        = Flag{0, 1 << 16}
	DoublyUnderlined = Flag{0, 1 << 17}
	SlowBlink        = Flag{0, 1 << 18}
	RapidBlink       = Flag{0, 1 << 19}
	Overlined        = Flag{0, 1 << 20}
	NotUnderlined    = Flag{0, 1 << 21}
	NotBlinking      = Flag{0, 1 << 22}
	NotOverlined     = Flag{0, 1 << 23}

	Conceal                    = Flag{0, 1 << 24}
	CrossedOut                 = Flag{0, 1 << 25}
	DisableProportionalSpacing = Flag{0, 1 << 26}
	Reveal                     = Flag{0, 1 << 27}
	NotCrossedOut              = Flag{0, 1 << 28}
	ProportionalSpacing        = Flag{0, 1 << 29}

	Framed                  = Flag{0, 1 << 32}
	Encircled               = Flag{0, 1 << 33}
	Superscript             = Flag{0, 1 << 34}
	Subscript               = Flag{0, 1 << 35}
	NotFramedEncircled      = Flag{0, 1 << 36}
	NotSuperscriptSubscript = Flag{0, 1 << 37}

	IdeogramUnderline       = Flag{0, 1 << 40}
	IdeogramDoubleUnderline = Flag{0, 1 << 41}
	IdeogramOverline        = Flag{0, 1 << 42}
	IdeogramDoubleOverline  = Flag{0, 1 << 43}
	IdeogramStressMarking   = Flag{0, 1 << 44}
	NoIdeogramAttributes    = Flag{0, 1 << 45}

	FgBlack   = Flag{0, 1 << 48}
	FgRed     = Flag{0, 1 << 49}
	FgGreen   = Flag{0, 1 << 50}
	FgYellow  = Flag{0, 1 << 51}
	FgBlue    = Flag{0, 1 << 52}
	FgMagenta = Flag{0, 1 << 53}
	FgCyan    = Flag{0, 1 << 54}
	FgWhite   = Flag{0, 1 << 55}

	BgBlack   = Flag{0, 1 << 56}
	BgRed     = Flag{0, 1 << 57}
	BgGreen   = Flag{0, 1 << 58}
	BgYellow  = Flag{0, 1 << 59}
	BgBlue    = Flag{0, 1 << 60}
	BgMagenta = Flag{0, 1 << 61}
	BgCyan    = Flag{0, 1 << 62}
	BgWhite   = Flag{0, 1 << 63}

	FgBrightBlack   = Flag{1 << 0, 0}
	FgBrightRed     = Flag{1 << 1, 0}
	FgBrightGreen   = Flag{1 << 2, 0}
	FgBrightYellow  = Flag{1 << 3, 0}
	FgBrightBlue    = Flag{1 << 4, 0}
	FgBrightMagenta = Flag{1 << 5, 0}
	FgBrightCyan    = Flag{1 << 6, 0}
	FgBrightWhite   = Flag{1 << 7, 0}

	BgBrightBlack   = Flag{1 << 8, 0}
	BgBrightRed     = Flag{1 << 9, 0}
	BgBrightGreen   = Flag{1 << 10, 0}
	BgBrightYellow  = Flag{1 << 11, 0}
	BgBrightBlue    = Flag{1 << 12, 0}
	BgBrightMagenta = Flag{1 << 13, 0}
	BgBrightCyan    = Flag{1 << 14, 0}
	BgBrightWhite   = Flag{1 << 15, 0}

	CustomFgColor         = Flag{1 << 16, 0}
	DefaultFgColor        = Flag{1 << 17, 0}
	CustomBgColor         = Flag{1 << 18, 0}
	DefaultBgColor        = Flag{1 << 19, 0}
	CustomUnderlineColor  = Flag{1 << 20, 0}
	DefaultUnderlineColor = Flag{1 << 21, 0}

	FgColor256        = Flag{1 << 24, 0}
	FgColorRGB        = Flag{1 << 25, 0}
	BgColor256        = Flag{1 << 26, 0}
	BgColorRGB        = Flag{1 << 27, 0}
	UnderlineColor256 = Flag{1 << 28, 0}
	UnderlineColorRGB = Flag{1 << 29, 0}

	PrimaryFont      = Flag{1 << 32, 0}
	AlternativeFont1 = Flag{1 << 33, 0}
	AlternativeFont2 = Flag{1 << 34, 0}
	AlternativeFont3 = Flag{1 << 35, 0}
	AlternativeFont4 = Flag{1 << 36, 0}
	AlternativeFont5 = Flag{1 << 37, 0}
	AlternativeFont6 = Flag{1 << 38, 0}
	AlternativeFont7 = Flag{1 << 39, 0}

	AlternativeFont8 = Flag{1 << 40, 0}
	AlternativeFont9 = Flag{1 << 41, 0}
)

func WrapFlag(s string, flag Flag) string {
	meta := flag.Meta()
	if meta != nil {
		return WrapCode(s, meta.Code)
	}
	return s
}
