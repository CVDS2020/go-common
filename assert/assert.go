package assert

import (
	"fmt"
	"reflect"
)

func Assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func NotNil[V any](obj V, name string) V {
	if any(obj) == nil || reflect.ValueOf(obj).IsNil() {
		panic(fmt.Errorf("%s must be not nil", name))
	}
	return obj
}

func PtrNotNil[V any](ptr *V, name string) *V {
	if ptr == nil {
		panic(fmt.Errorf("%s must be not nil", name))
	}
	return ptr
}

func SliceNotNil[E any](slice []E, name string) []E {
	if slice == nil {
		panic(fmt.Errorf("%s must be not nil", name))
	}
	return slice
}

func MapNotNil[K comparable, V any](m map[K]V, name string) map[K]V {
	if m == nil {
		panic(fmt.Errorf("%s must be not nil", name))
	}
	return m
}

func ChanNotNil[E any](c chan E, name string) chan E {
	if c == nil {
		panic(fmt.Errorf("%s must be not nil", name))
	}
	return c
}

func IsNil[V any](obj V, name string) V {
	if any(obj) != nil && !reflect.ValueOf(obj).IsNil() {
		panic(fmt.Errorf("%s must nil", name))
	}
	return obj
}

func PtrIsNil[V any](ptr *V, name string) *V {
	if ptr != nil {
		panic(fmt.Errorf("%s must be not nil", name))
	}
	return ptr
}

func NotEmpty[V any](obj V, name string) V {
	Assert(reflect.ValueOf(obj).Len() != 0, fmt.Sprintf("%s must be not empty", name))
	return obj
}

func MustEmpty[V any](obj V, name string) V {
	Assert(reflect.ValueOf(obj).Len() == 0, fmt.Sprintf("%s must be not empty", name))
	return obj
}

func NotZero[V any](obj V, name string) V {
	Assert(!reflect.ValueOf(obj).IsZero(), fmt.Sprintf("%s must be not zero", name))
	return obj
}

func IsZero[V any](obj V, name string) V {
	Assert(reflect.ValueOf(obj).IsZero(), fmt.Sprintf("%s must be not zero", name))
	return obj
}
