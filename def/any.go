package def

import "reflect"

func SetAny[O any](o O, def O) O {
	if reflect.ValueOf(o).IsNil() {
		return def
	}
	return o
}

func SetterAny[O any](o O, setter func() O) O {
	if reflect.ValueOf(o).IsNil() {
		return setter()
	}
	return o
}

func SetAnyP[O any](op *O, def O) {
	if reflect.ValueOf(*op).IsNil() {
		*op = def
	}
}

func SetterAnyP[O any](op *O, setter func() O) {
	if reflect.ValueOf(*op).IsNil() {
		*op = setter()
	}
}
