package def

func MakeMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return make(map[K]V)
	}
	return m
}

func MakeSizeMap[K comparable, V any](m map[K]V, size int) map[K]V {
	if m == nil {
		return make(map[K]V, size)
	}
	return m
}

func SetMap[K comparable, V any](m map[K]V, def map[K]V) map[K]V {
	if m == nil {
		return def
	}
	return m
}

func SetterMap[K comparable, V any](m map[K]V, setter func() map[K]V) map[K]V {
	if m == nil {
		return setter()
	}
	return m
}

func MakeMapP[K comparable, V any](mp *map[K]V) {
	if *mp == nil {
		*mp = make(map[K]V)
	}
}

func MakeSizeMapP[K comparable, V any](mp *map[K]V, size int) {
	if *mp == nil {
		*mp = make(map[K]V, size)
	}
}

func SetMapP[K comparable, V any](mp *map[K]V, def map[K]V) {
	if *mp == nil {
		*mp = def
	}
}

func SetterMapP[K comparable, V any](mp *map[K]V, setter func() map[K]V) {
	if *mp == nil {
		*mp = setter()
	}
}
