package uns

import "unsafe"

type StringStruct struct {
	Ptr unsafe.Pointer
	Len int
}

func StringStructOf(s string) StringStruct {
	return *(*StringStruct)(unsafe.Pointer(&s))
}

func StringStructOfPtr(sp *string) *StringStruct {
	return (*StringStruct)(unsafe.Pointer(sp))
}

func (ss *StringStruct) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&SliceStruct{Ptr: ss.Ptr, Len: ss.Len, Cap: ss.Len}))
}

func (ss *StringStruct) ToString() string {
	return *(*string)(unsafe.Pointer(ss))
}

func (ss *StringStruct) ToBytesStruct() *BytesStruct {
	return &BytesStruct{Ptr: ss.Ptr, Len: ss.Len, Cap: ss.Len}
}

func MakeString(ptr unsafe.Pointer, len int) string {
	ss := StringStruct{Ptr: ptr, Len: len}
	return ss.ToString()
}

func StringToBytes(s string) []byte {
	ss := (*StringStruct)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&SliceStruct{Ptr: ss.Ptr, Len: ss.Len, Cap: ss.Len}))
}
