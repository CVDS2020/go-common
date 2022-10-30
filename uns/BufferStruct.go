package uns

import (
	"bytes"
	"unsafe"
)

type Buffer struct {
	Buf      []byte
	Off      int
	LastRead int8
}

func BufferOf(buf *bytes.Buffer) *Buffer {
	return (*Buffer)(unsafe.Pointer(buf))
}
