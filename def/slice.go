package def

func MakeSlice[E any](s []E, size int) []E {
	if s == nil {
		return make([]E, size)
	}
	return s
}

func SetSlice[E any](s []E, def []E) []E {
	if s == nil {
		return def
	}
	return s
}

func SetterSlice[E any](s []E, setter func() []E) []E {
	if s == nil {
		return setter()
	}
	return s
}

func MakeSliceP[E any](sp *[]E, size int) {
	if *sp == nil {
		*sp = make([]E, size)
	}
}

func SetSliceP[E any](sp *[]E, def []E) {
	if *sp == nil {
		*sp = def
	}
}

func SetterSliceP[E any](sp *[]E, setter func() []E) {
	if *sp == nil {
		*sp = setter()
	}
}
