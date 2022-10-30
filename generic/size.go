package generic

import "unsafe"

func Size[T any]() uintptr {
	var o T
	return unsafe.Sizeof(o)
}
