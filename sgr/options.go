package sgr

type Option interface {
	apply(sgr *SGR)
}

type optionFunc func(sgr *SGR)

func (f optionFunc) apply(sgr *SGR) {
	f(sgr)
}

func FgColor256Option(color256 Color256) Option {
	return optionFunc(func(sgr *SGR) {
		sgr.Flag = sgr.Flag.Add(FgColor256)
		sgr.FgColor256 = color256
	})
}

func FgColorRGBOption(rgb ColorRGB) Option {
	return optionFunc(func(sgr *SGR) {
		sgr.Flag = sgr.Flag.Add(FgColorRGB)
		sgr.FgColorRGB = rgb
	})
}

func FgRGBOption(r, g, b uint8) Option {
	return FgColorRGBOption(RGB(r, g, b))
}

func FgRGB24Option(rgb uint32) Option {
	return FgColorRGBOption(RGB24(rgb))
}

func BgColor256Option(color256 Color256) Option {
	return optionFunc(func(sgr *SGR) {
		sgr.Flag = sgr.Flag.Add(BgColor256)
		sgr.BgColor256 = color256
	})
}

func BgColorRGBOption(rgb ColorRGB) Option {
	return optionFunc(func(sgr *SGR) {
		sgr.Flag = sgr.Flag.Add(BgColorRGB)
		sgr.BgColorRGB = rgb
	})
}

func BgRGBOption(r, g, b uint8) Option {
	return BgColorRGBOption(RGB(r, g, b))
}

func BgRGB24Option(rgb uint32) Option {
	return BgColorRGBOption(RGB24(rgb))
}

func UlColor256Option(color256 Color256) Option {
	return optionFunc(func(sgr *SGR) {
		sgr.Flag = sgr.Flag.Add(UnderlineColor256)
		sgr.UlColor256 = color256
	})
}

func UlColorRGBOption(rgb ColorRGB) Option {
	return optionFunc(func(sgr *SGR) {
		sgr.Flag = sgr.Flag.Add(UnderlineColorRGB)
		sgr.UlColorRGB = rgb
	})
}

func UlRGBOption(r, g, b uint8) Option {
	return UlColorRGBOption(RGB(r, g, b))
}

func UlRGB24Option(rgb uint32) Option {
	return UlColorRGBOption(RGB24(rgb))
}
