package sgr

import (
	"strconv"
)

type ColorRGB struct {
	R uint8
	G uint8
	B uint8
}

func RGB(r, g, b uint8) ColorRGB {
	return ColorRGB{R: r, G: g, B: b}
}

func RGB24(rgb uint32) ColorRGB {
	return ColorRGB{R: byte((rgb >> 16) & 0xff), G: byte((rgb >> 8) & 0xff), B: byte(rgb & 0xff)}
}

type Color256 uint8

func WrapColor(s string, color Flag) string {
	return WrapFlag(s, color.Color())
}

func WrapColor256(s string, color256 Color256) string {
	return WrapCodes(s, CustomFgColorCode, "5", u8ts[color256])
}

func WrapColorRGB(s string, rgb ColorRGB) string {
	return WrapCodes(s, CustomFgColorCode, "2", u8ts[rgb.R], u8ts[rgb.G], u8ts[rgb.B])
}

func WrapRGB(s string, r, g, b uint8) string {
	return WrapColorRGB(s, RGB(r, g, b))
}

func WrapRGB24(s string, rgb uint32) string {
	return WrapColorRGB(s, RGB24(rgb))
}

var (
	fgNormalColorTable [256]string
	fgBrightColorTable [256]string
	bgNormalColorTable [256]string
	bgBrightColorTable [256]string
)

func initFgNormalColorTable() {
	for i := 0; i < 8; i++ {
		s := strconv.Itoa(i + 30)
		for j := 1 << i; j < 1<<(i+1); j++ {
			fgNormalColorTable[j] = s
		}
	}
}

func initFgBrightColorTable() {
	for i := 0; i < 8; i++ {
		s := strconv.Itoa(i + 90)
		for j := 1 << i; j < 1<<(i+1); j++ {
			fgBrightColorTable[j] = s
		}
	}
}

func initBgNormalColorTable() {
	for i := 0; i < 8; i++ {
		s := strconv.Itoa(i + 40)
		for j := 1 << i; j < 1<<(i+1); j++ {
			bgNormalColorTable[j] = s
		}
	}
}

func initBgBrightColorTable() {
	for i := 0; i < 8; i++ {
		s := strconv.Itoa(i + 100)
		for j := 1 << i; j < 1<<(i+1); j++ {
			bgBrightColorTable[j] = s
		}
	}
}

func init() {
	initFgNormalColorTable()
	initFgBrightColorTable()
	initBgNormalColorTable()
	initBgBrightColorTable()
}
