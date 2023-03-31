package uns

import "unsafe"

type BytesStruct SliceStruct

func BytesStructOf(bs []byte) BytesStruct {
	return *(*BytesStruct)(unsafe.Pointer(&bs))
}

func BytesStructOfPtr(bsp *[]byte) *BytesStruct {
	return (*BytesStruct)(unsafe.Pointer(bsp))
}

func BytesStructOfString(s string) *BytesStruct {
	ss := (*StringStruct)(unsafe.Pointer(&s))
	return &BytesStruct{Ptr: ss.Ptr, Len: ss.Len, Cap: ss.Len}
}

func (bss *BytesStruct) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(bss))
}

func (bss *BytesStruct) ToString() string {
	return *(*string)(unsafe.Pointer(bss))
}

func (bss *BytesStruct) ToStringStruct() StringStruct {
	return StringStruct{Ptr: bss.Ptr, Len: bss.Len}
}

func MakeBytes(ptr unsafe.Pointer, len, cap int) []byte {
	return MakeSlice[byte](ptr, len, cap)
}

func MakeBytesUnchecked(ptr unsafe.Pointer, len, cap int) (bs []byte) {
	*ConvertPointer[[]byte, SliceStruct](&bs) = SliceStruct{Ptr: ptr, Len: len, Cap: cap}
	return
}

func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
