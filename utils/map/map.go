package mapUtils

func Copy[K comparable, V any](m map[K]V) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		n[k] = v
	}
	return n
}

func Keys[K comparable, V any](m map[K]V) (keys []K) {
	keys = make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func Values[K comparable, V any](m map[K]V) (values []V) {
	values = make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return
}
