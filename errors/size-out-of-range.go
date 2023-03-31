package errors

import (
	"fmt"
	"math"
)

type SizeOutOfRange struct {
	Target   string
	MinSize  int64
	MaxSize  int64
	Size     int64
	RealSize bool
}

func NewSizeOutOfRange(target string, min, max, size int64, real bool) error {
	return &SizeOutOfRange{Target: target, MinSize: min, MaxSize: max, Size: size, RealSize: real}
}

func (e *SizeOutOfRange) Error() string {
	target := e.Target
	if target == "" {
		target = "数据"
	}
	var mode string
	if e.RealSize {
		mode = "实际"
	} else {
		mode = "已解析"
	}
	switch {
	case e.MinSize == math.MinInt64:
		return fmt.Sprintf("%s大小超过限制(,%d)，%s大小为(%d)", e.Target, e.MaxSize, mode, e.Size)
	case e.MaxSize == math.MaxInt64:
		return fmt.Sprintf("%s大小超过限制(%d,)，%s大小为(%d)", e.Target, e.MinSize, mode, e.Size)
	default:
		return fmt.Sprintf("%s大小超过限制(%d,%d)，%s大小为(%d)", e.Target, e.MinSize, e.MaxSize, mode, e.Size)
	}
}
