package sgr

import (
	"testing"
)

func TestSGR(t *testing.T) {
	println(WrapSGR("hello", NewSGR(Bold).Options(UlRGB24Option(0x629755))))
	println(WrapRGB24("hi", 0x629755))
	println(WrapColor("lili", FgCyan))
	println(WrapColor256("pip", 3))
}
