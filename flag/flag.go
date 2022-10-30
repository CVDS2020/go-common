package flag

func TestFlag(flag int, mask int) bool {
	return flag&mask != 0
}

func TestFlags(flag int, masks ...int) bool {
	for _, mask := range masks {
		flag &= mask
	}
	return flag != 0
}

func MaskFlag(flag int, mask int) int {
	return flag | mask
}

func MasksFlag(flag int, masks ...int) int {
	for _, mask := range masks {
		flag |= mask
	}
	return flag
}

func UnmaskFlag(flag int, mask int) int {
	return flag & ^mask
}

func UnmasksFlag(flag int, masks ...int) int {
	for _, mask := range masks {
		flag &= ^mask
	}
	return flag
}

func SwapFlagMask(flag int, unmask, mask int) int {
	flag &= ^unmask
	return flag | mask
}
