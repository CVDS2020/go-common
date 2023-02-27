package container

import "fmt"

type LinkedMapEntry[K comparable, V any] struct {
	key   K
	value V
	next  *LinkedMapEntry[K, V]
	prev  *LinkedMapEntry[K, V]
}

func NewLinkedMapEntry[K comparable, V any](key K, value V) *LinkedMapEntry[K, V] {
	return &LinkedMapEntry[K, V]{key: key, value: value}
}

func (e *LinkedMapEntry[K, V]) Key() K {
	return e.key
}

func (e *LinkedMapEntry[K, V]) Value() V {
	return e.value
}

func (e *LinkedMapEntry[K, V]) SetValue(v V) {
	e.value = v
}

func (e *LinkedMapEntry[K, V]) Next() *LinkedMapEntry[K, V] {
	if e.next == nil || e.next.next == nil {
		return nil
	}
	return e.next
}

func (e *LinkedMapEntry[K, V]) Prev() *LinkedMapEntry[K, V] {
	if e.prev == nil || e.prev.prev == nil {
		return nil
	}
	return e.prev
}

func (e *LinkedMapEntry[K, V]) Removed() bool {
	return e.next == nil && e.prev == nil
}

func (e *LinkedMapEntry[K, V]) String() string {
	return fmt.Sprintf("{key: %v, value: %v}", e.key, e.value)
}

type LinkedMap[K comparable, V any] struct {
	head *LinkedMapEntry[K, V]
	tail *LinkedMapEntry[K, V]
	m    map[K]*LinkedMapEntry[K, V]
}

func NewLinkedMap[K comparable, V any](size int) *LinkedMap[K, V] {
	m := &LinkedMap[K, V]{
		head: new(LinkedMapEntry[K, V]),
		tail: new(LinkedMapEntry[K, V]),
		m:    make(map[K]*LinkedMapEntry[K, V], size),
	}
	m.head.next = m.tail
	m.tail.prev = m.head
	return m
}

func (m *LinkedMap[K, V]) addLink(e *LinkedMapEntry[K, V], p *LinkedMapEntry[K, V]) {
	e.prev = p
	e.next = p.next
	p.next.prev = e
	p.next = e
}

func (m *LinkedMap[K, V]) removeLink(e *LinkedMapEntry[K, V]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next, e.prev = nil, nil
}

func (m *LinkedMap[K, V]) put(e *LinkedMapEntry[K, V]) {
	m.m[e.key] = e
	m.addLink(e, m.tail.prev)
}

func (m *LinkedMap[K, V]) remove(e *LinkedMapEntry[K, V]) {
	delete(m.m, e.key)
	m.removeLink(e)
}

func (m *LinkedMap[K, V]) Len() int {
	return len(m.m)
}

func (m *LinkedMap[K, V]) Size() uint {
	return uint(len(m.m))
}

func (m *LinkedMap[K, V]) Empty() bool {
	return len(m.m) == 0
}

func (m *LinkedMap[K, V]) Has(key K) bool {
	_, has := m.m[key]
	return has
}

func (m *LinkedMap[K, V]) Get(key K) (value V, exist bool) {
	if e, exist := m.m[key]; exist {
		return e.value, exist
	}
	return value, false
}

func (m *LinkedMap[K, V]) First() (value V, exist bool) {
	return m.head.next.value, len(m.m) != 0
}

func (m *LinkedMap[K, V]) Last() (value V, exist bool) {
	return m.tail.prev.value, len(m.m) != 0
}

func (m *LinkedMap[K, V]) GetEntry(key K) *LinkedMapEntry[K, V] {
	return m.m[key]
}

func (m *LinkedMap[K, V]) FirstEntry() *LinkedMapEntry[K, V] {
	if len(m.m) == 0 {
		return nil
	}
	return m.head.next
}

func (m *LinkedMap[K, V]) LastEntry() *LinkedMapEntry[K, V] {
	if len(m.m) == 0 {
		return nil
	}
	return m.tail.prev
}

func (m *LinkedMap[K, V]) Put(key K, value V) (entry, old *LinkedMapEntry[K, V]) {
	if e, exist := m.m[key]; exist {
		old = e
		m.remove(e)
	}
	entry = NewLinkedMapEntry[K, V](key, value)
	m.put(entry)
	return
}

func (m *LinkedMap[K, V]) PutIfAbsent(key K, value V) (entry *LinkedMapEntry[K, V], exist bool) {
	return m.PutEntryIfAbsent(NewLinkedMapEntry[K, V](key, value))
}

func (m *LinkedMap[K, V]) PutEntry(e *LinkedMapEntry[K, V]) *LinkedMapEntry[K, V] {
	if !e.Removed() {
		panic("linked map entry not removed")
	}
	if old, exist := m.m[e.key]; exist {
		m.remove(old)
		m.put(e)
		return old
	}
	m.put(e)
	return nil
}

func (m *LinkedMap[K, V]) PutEntryIfAbsent(e *LinkedMapEntry[K, V]) (entry *LinkedMapEntry[K, V], exist bool) {
	if !e.Removed() {
		panic("linked map entry not removed")
	}
	if entry, exist = m.m[e.key]; exist {
		return entry, true
	}
	m.put(e)
	return e, false
}

func (m *LinkedMap[K, V]) Remove(key K) *LinkedMapEntry[K, V] {
	if e, exist := m.m[key]; exist {
		m.remove(e)
		return e
	}
	return nil
}

func (m *LinkedMap[K, V]) RemoveFirst() *LinkedMapEntry[K, V] {
	if len(m.m) == 0 {
		return nil
	}
	e := m.head.next
	m.remove(e)
	return e
}

func (m *LinkedMap[K, V]) RemoveLast() *LinkedMapEntry[K, V] {
	if len(m.m) == 0 {
		return nil
	}
	e := m.tail.prev
	m.remove(e)
	return e
}

func (m *LinkedMap[K, V]) RemoveEntry(e *LinkedMapEntry[K, V]) bool {
	if e.Removed() {
		return false
	}
	m.remove(e)
	return true
}

func (m *LinkedMap[K, V]) Clear() {
	for _, e := range m.m {
		e.next, e.prev = nil, nil
	}
	for key := range m.m {
		delete(m.m, key)
	}
	m.head.next = m.tail
	m.tail.prev = m.head
}
