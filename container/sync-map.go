package container

import "sync"

type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func (s *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	var v any
	v, ok = s.m.Load(key)
	if ok {
		value = v.(V)
	}
	return
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}

func (s *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var v any
	v, loaded = s.m.LoadOrStore(key, value)
	if loaded {
		actual = v.(V)
	}
	return
}

func (s *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var v any
	v, loaded = s.m.LoadAndDelete(key)
	if loaded {
		value = v.(V)
	}
	return
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}

func (s *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	s.m.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}

func (s *SyncMap[K, V]) Map() map[K]V {
	m := make(map[K]V)
	s.Range(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

func (s *SyncMap[K, V]) Keys() (keys []K) {
	s.Range(func(key K, value V) bool {
		keys = append(keys, key)
		return true
	})
	return
}

func (s *SyncMap[K, V]) Values() (values []V) {
	s.Range(func(key K, value V) bool {
		values = append(values, value)
		return true
	})
	return
}
