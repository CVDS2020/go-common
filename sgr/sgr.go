package sgr

import (
	"gitee.com/sy_183/common/slice"
	"strconv"
)

var u8ts [256]string

func initU8ts() {
	for i := 0; i < 256; i++ {
		u8ts[i] = strconv.Itoa(i)
	}
}

func init() {
	initU8ts()
}

// SGR Select Graphic Rendition
type SGR struct {
	Flag       Flag
	FgColor256 Color256
	FgColorRGB ColorRGB
	BgColor256 Color256
	BgColorRGB ColorRGB
	UlColor256 Color256
	UlColorRGB ColorRGB
}

func NewSGR(flags ...Flag) *SGR {
	return &SGR{Flag: Flags(flags...)}
}

func (s *SGR) Options(opts ...Option) *SGR {
	for _, opt := range opts {
		opt.apply(s)
	}
	return s
}

func (s *SGR) CachedCodes() []string {
	var cs []string
	if c := cacheTable1[s.Flag.byte(1)]; c != "" {
		cs = append(cs, c)
	}
	if c := cacheTable2[s.Flag.byte(2)]; c != "" {
		cs = append(cs, c)
	}
	if c := cacheTable3[s.Flag.byte(3)&0b00111111]; c != "" {
		cs = append(cs, c)
	}
	if c := cacheTable4[s.Flag.byte(4)&0b00111111]; c != "" {
		cs = append(cs, c)
	}
	if c := cacheTable5[s.Flag.byte(5)&0b00111111]; c != "" {
		cs = append(cs, c)
	}
	return cs
}

func (s *SGR) FgColorCodes() []string {
	var cs []string
	if s.Flag.Has(DefaultFgColor) {
		cs = append(cs, DefaultFgColorCode)
	} else if s.Flag.Has(FgColorRGB) {
		cs = append(cs, CustomFgColorCode, "2", u8ts[s.FgColorRGB.R], u8ts[s.FgColorRGB.G], u8ts[s.FgColorRGB.B])
	} else if s.Flag.Has(FgColor256) {
		cs = append(cs, CustomFgColorCode, "5", u8ts[s.FgColor256])
	} else if c := fgNormalColorTable[s.Flag.byte(6)]; c != "" {
		cs = append(cs, c)
	} else if c := fgBrightColorTable[s.Flag.byte(8)]; c != "" {
		cs = append(cs, c)
	}
	return cs
}

func (s *SGR) BgColorCodes() []string {
	var cs []string
	if s.Flag.Has(DefaultBgColor) {
		cs = append(cs, DefaultBgColorCode)
	} else if s.Flag.Has(BgColorRGB) {
		cs = append(cs, CustomBgColorCode, "2", u8ts[s.BgColorRGB.R], u8ts[s.BgColorRGB.G], u8ts[s.BgColorRGB.B])
	} else if s.Flag.Has(BgColor256) {
		cs = append(cs, CustomBgColorCode, "5", u8ts[s.BgColor256])
	} else if c := bgNormalColorTable[s.Flag.byte(7)]; c != "" {
		cs = append(cs, c)
	} else if c := bgBrightColorTable[s.Flag.byte(9)]; c != "" {
		cs = append(cs, c)
	}
	return cs
}

func (s *SGR) UlColorCodes() []string {
	var cs []string
	if s.Flag.Has(DefaultUnderlineColor) {
		cs = append(cs, DefaultUnderlineColorCode)
	} else if s.Flag.Has(UnderlineColorRGB) {
		cs = append(cs, CustomUnderlineColorCode, "2", u8ts[s.UlColorRGB.R], u8ts[s.UlColorRGB.G], u8ts[s.UlColorRGB.B])
	} else if s.Flag.Has(UnderlineColor256) {
		cs = append(cs, CustomUnderlineColorCode, "5", u8ts[s.UlColor256])
	}
	return cs
}

func (s *SGR) Codes() []string {
	if s.Flag.Has(Reset) {
		return []string{ResetCode}
	}
	return slice.JoinNew(s.CachedCodes(), s.FgColorCodes(), s.BgColorCodes(), s.UlColorCodes())
}

func (s *SGR) CSI() string {
	return MakeCSI(s.Codes()...)
}

func (s *SGR) Wrap(str string) string {
	return WrapCodes(str, s.Codes()...)
}

func WrapSGR(s string, sgr *SGR) string {
	return WrapCodes(s, sgr.Codes()...)
}
