package generic

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Integer interface {
	Unsigned | Signed
}

type Float interface {
	~float32 | ~float64
}

type Complex interface {
	~complex64 | ~complex128
}

type OrderedNumber interface {
	Integer | Float
}

type Number interface {
	OrderedNumber | Complex
}

type Ordered interface {
	OrderedNumber | ~string
}
