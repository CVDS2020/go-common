package uns

import "unsafe"

func Convert[SRC, DST any](src SRC) DST {
	return *(*DST)(ToPointer(&src))
}

func ConvertPointer[SRC, DST any](src *SRC) *DST {
	return (*DST)(unsafe.Pointer(src))
}

func FromPointer[SRC, DST any](src *SRC) DST {
	return *(*DST)(unsafe.Pointer(src))
}

func ToPointer[X any](x X) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&x))
}

func Copy(dst, src unsafe.Pointer, size int) {
	switch size {
	case 0:
		return
	case 1:
		*(*uint8)(dst) = *(*uint8)(src)
	case 2:
		*(*uint16)(dst) = *(*uint16)(src)
	case 4:
		*(*uint32)(dst) = *(*uint32)(src)
	case 8:
		*(*uint64)(dst) = *(*uint64)(src)
	case 16:
		*(*complex128)(dst) = *(*complex128)(src)
	case 24:
		*(*[24]byte)(dst) = *(*[24]byte)(src)
	case 32:
		*(*[32]byte)(dst) = *(*[32]byte)(src)
	default:
		copy(
			FromPointer[BytesStruct, []byte](&BytesStruct{Ptr: dst, Len: size, Cap: size}),
			FromPointer[BytesStruct, []byte](&BytesStruct{Ptr: src, Len: size, Cap: size}),
		)
	}
}
