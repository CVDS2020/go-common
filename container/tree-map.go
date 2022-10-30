package container

import (
	"fmt"
	"gitee.com/sy_183/common/generic"
)

const (
	treeMapColorNone = iota
	treeMapColorRed
	treeMapColorBlack
)

type TreeMapEntry[K, V any] struct {
	key    K
	value  V
	left   *TreeMapEntry[K, V]
	right  *TreeMapEntry[K, V]
	parent *TreeMapEntry[K, V]
	color  int
}

func NewTreeMapEntry[K, V any](key K, value V) *TreeMapEntry[K, V] {
	return &TreeMapEntry[K, V]{key: key, value: value}
}

func (e *TreeMapEntry[K, V]) Key() K {
	return e.key
}

func (e *TreeMapEntry[K, V]) Value() V {
	return e.value
}

func (e *TreeMapEntry[K, V]) next() *TreeMapEntry[K, V] {
	if e.right != nil {
		p := e.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		p := e.parent
		ch := e
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}
}

func (e *TreeMapEntry[K, V]) prev() *TreeMapEntry[K, V] {
	if e.left != nil {
		p := e.left
		for p.right != nil {
			p = p.right
		}
		return p
	} else {
		p := e.parent
		ch := e
		for p != nil && ch == p.left {
			ch = p
			p = p.parent
		}
		return p
	}
}

func (e *TreeMapEntry[K, V]) Next() *TreeMapEntry[K, V] {
	if e.Removed() {
		return nil
	}
	return e.next()
}

func (e *TreeMapEntry[K, V]) Prev() *TreeMapEntry[K, V] {
	if e.Removed() {
		return nil
	}
	return e.prev()
}

func (e *TreeMapEntry[K, V]) Removed() bool {
	return e.color == treeMapColorNone
}

func (e *TreeMapEntry[K, V]) String() string {
	return fmt.Sprintf("{key: %v, value: %v}", e.key, e.value)
}

type Comparator[V any] func(o1, o2 V) int

type TreeMap[K, V any] struct {
	root *TreeMapEntry[K, V]
	head *TreeMapEntry[K, V]
	tail *TreeMapEntry[K, V]
	cpr  Comparator[K]
	size uint
}

func NewTreeMap[K, V any](cpr Comparator[K]) *TreeMap[K, V] {
	return &TreeMap[K, V]{cpr: cpr}
}

func (m *TreeMap[K, V]) Size() uint {
	return m.size
}

func (m *TreeMap[K, V]) Empty() bool {
	return m.size == 0
}

func (m *TreeMap[K, V]) Has(key K) bool {
	return m.getEntry(key) != nil
}

func (m *TreeMap[K, V]) Get(key K) (value V, exist bool) {
	if e := m.getEntry(key); e != nil {
		return e.value, true
	}
	return value, false
}

func (m *TreeMap[K, V]) GetFirst() (value V, exist bool) {
	if m.head == nil {
		return value, false
	}
	return m.head.value, true
}

func (m *TreeMap[K, V]) GetLast() (value V, exist bool) {
	if m.tail == nil {
		return value, false
	}
	return m.tail.value, true
}

func (m *TreeMap[K, V]) GetEntry(key K) *TreeMapEntry[K, V] {
	return m.getEntry(key)
}

func (m *TreeMap[K, V]) GetFirstEntry() *TreeMapEntry[K, V] {
	return m.head
}

func (m *TreeMap[K, V]) GetLastEntry() *TreeMapEntry[K, V] {
	return m.tail
}

func (m *TreeMap[K, V]) Put(key K, value V) (entry *TreeMapEntry[K, V], old *TreeMapEntry[K, V]) {
	entry = NewTreeMapEntry[K, V](key, value)
	old = m.addEntry(entry, true)
	return
}

func (m *TreeMap[K, V]) PutIfAbsent(key K, value V) (entry *TreeMapEntry[K, V], exist bool) {
	return m.PutEntryIfAbsent(NewTreeMapEntry[K, V](key, value))
}

func (m *TreeMap[K, V]) PutEntry(e *TreeMapEntry[K, V]) *TreeMapEntry[K, V] {
	if !e.Removed() {
		panic("tree map entry not removed")
	}
	return m.addEntry(e, true)
}

func (m *TreeMap[K, V]) PutEntryIfAbsent(e *TreeMapEntry[K, V]) (entry *TreeMapEntry[K, V], exist bool) {
	if !e.Removed() {
		panic("tree map entry not removed")
	}
	if entry = m.addEntry(e, false); entry != nil {
		return entry, true
	}
	return e, false
}

func (m *TreeMap[K, V]) Remove(key K) *TreeMapEntry[K, V] {
	if p := m.getEntry(key); p != nil {
		m.deleteEntry(p)
		return p
	}
	return nil
}

func (m *TreeMap[K, V]) RemoveFirst() *TreeMapEntry[K, V] {
	if m.head != nil {
		removed := m.head
		m.deleteEntry(m.head)
		return removed
	}
	return nil
}

func (m *TreeMap[K, V]) RemoveLast() *TreeMapEntry[K, V] {
	if m.tail != nil {
		removed := m.tail
		m.deleteEntry(m.tail)
		return removed
	}
	return nil
}

func (m *TreeMap[K, V]) RemoveEntry(e *TreeMapEntry[K, V]) bool {
	if e.Removed() {
		return false
	}
	m.deleteEntry(e)
	return true
}

func (m *TreeMap[K, V]) Clear() {
	var next *TreeMapEntry[K, V]
	for e := m.head; e != nil; e = next {
		next = e.next()
		m.deleteEntry(e)
	}
}

func (m *TreeMap[K, V]) redirectRelation(p, l, r, o, n *TreeMapEntry[K, V]) {
	if p != nil {
		if p.left == o {
			p.left = n
		} else if p.right == o {
			p.right = n
		} else {
			panic("parent not parent of old")
		}
	}
	if l != nil {
		l.parent = n
	}
	if r != nil {
		r.parent = n
	}
}

func (m *TreeMap[K, V]) replaceRelation(old, e *TreeMapEntry[K, V]) {
	e.parent = old.parent
	e.left, e.right = old.left, old.right
	e.color = old.color
	m.redirectRelation(old.parent, old.left, old.right, old, e)
	if m.root == old {
		m.root = e
	}
	if m.head == old {
		m.head = e
	}
	if m.tail == old {
		m.tail = e
	}
	old.parent = nil
	old.left, old.right = nil, nil
	old.color = treeMapColorNone
}

func (m *TreeMap[K, V]) swapRelation(e1, e2 *TreeMapEntry[K, V]) {
	p1, l1, r1 := e1.parent, e1.left, e1.right
	p2, l2, r2 := e2.parent, e2.left, e2.right
	m.redirectRelation(p1, l1, r1, e1, e2)
	m.redirectRelation(p2, l2, r2, e2, e1)
	e1.parent, e2.parent = e2.parent, e1.parent
	e1.left, e2.left = e2.left, e1.left
	e1.right, e2.right = e2.right, e1.right
	e1.color, e2.color = e2.color, e1.color
	if m.root == e1 {
		m.root = e2
	} else if m.root == e2 {
		m.root = e1
	}
	if m.head == e1 {
		m.head = e2
	} else if m.head == e2 {
		m.head = e1
	}
	if m.tail == e1 {
		m.tail = e2
	} else if m.tail == e2 {
		m.tail = e1
	}
}

func (m *TreeMap[K, V]) getEntry(key K) *TreeMapEntry[K, V] {
	t := m.root
	if t == nil {
		return nil
	}
	if m.cpr != nil {
		for t != nil {
			cmp := m.cpr(key, t.key)
			if cmp < 0 {
				t = t.left
			} else if cmp > 0 {
				t = t.right
			} else {
				return t
			}
		}
	} else {
		com := (any(key)).(Comparable)
		for t != nil {
			cmp := com.Compare((any(t.key)).(Comparable))
			if cmp < 0 {
				t = t.left
			} else if cmp > 0 {
				t = t.right
			} else {
				return t
			}
		}
	}
	return nil
}

func (m *TreeMap[K, V]) addEntry(e *TreeMapEntry[K, V], replaceOld bool) *TreeMapEntry[K, V] {
	t := m.root
	e.color = treeMapColorBlack
	if t == nil {
		m.addEntryToEmpty(e)
		return nil
	}
	var parent *TreeMapEntry[K, V]
	var cmp int
	if m.cpr != nil {
		for t != nil {
			parent = t
			cmp = m.cpr(e.key, t.key)
			if cmp < 0 {
				t = t.left
			} else if cmp > 0 {
				t = t.right
			} else {
				if replaceOld {
					m.replaceRelation(t, e)
				}
				return t
			}
		}
	} else {
		com := (any(e.key)).(Comparable)
		for t != nil {
			parent = t
			cmp = com.Compare((any(t.key)).(Comparable))
			if cmp < 0 {
				t = t.left
			} else if cmp > 0 {
				t = t.right
			} else {
				if replaceOld {
					m.replaceRelation(t, e)
				}
				return t
			}
		}
	}
	m.addEntryToParent(e, parent, cmp < 0)
	return nil
}

func (m *TreeMap[K, V]) addEntryToEmpty(e *TreeMapEntry[K, V]) {
	m.root = e
	m.head, m.tail = e, e
	m.size = 1
}

func (m *TreeMap[K, V]) addEntryToParent(e, p *TreeMapEntry[K, V], addToLeft bool) {
	e.parent = p
	if addToLeft {
		p.left = e
		if p == m.head {
			m.head = e
		}
	} else {
		p.right = e
		if p == m.tail {
			m.tail = e
		}
	}
	m.fixAfterInsertion(e)
	m.size++
}

func (m *TreeMap[K, V]) deleteEntry(p *TreeMapEntry[K, V]) {
	m.size--

	// If strictly internal, copy successor's element to p and then make p
	// point to successor.
	if p.left != nil && p.right != nil {
		m.swapRelation(p, p.next())
	} // p has 2 children

	if p == m.head {
		m.head = p.next()
	}
	if p == m.tail {
		m.tail = p.prev()
	}

	// Start fixup at replacement node, if it exists.
	var replacement *TreeMapEntry[K, V]
	if p.left != nil {
		replacement = p.left
	} else {
		replacement = p.right
	}

	if replacement != nil {
		// Link replacement to parent
		replacement.parent = p.parent
		if p.parent == nil {
			m.root = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		// Null out links so they are OK to use by fixAfterDeletion.
		p.left, p.right, p.parent = nil, nil, nil

		// Fix replacement
		if p.color == treeMapColorBlack {
			m.fixAfterDeletion(replacement)
		}
	} else if p.parent == nil { // return if we are the only node.
		m.root = nil
	} else { //  No children. Use self as phantom replacement and unlink.
		if p.color == treeMapColorBlack {
			m.fixAfterDeletion(p)
		}

		if p.parent != nil {
			if p == p.parent.left {
				p.parent.left = nil
			} else if p == p.parent.right {
				p.parent.right = nil
			}
			p.parent = nil
		}
	}

	p.color = treeMapColorNone
}

func (*TreeMap[K, V]) colorOf(p *TreeMapEntry[K, V]) int {
	if p == nil {
		return treeMapColorBlack
	}
	return p.color
}

func (*TreeMap[K, V]) parentOf(p *TreeMapEntry[K, V]) *TreeMapEntry[K, V] {
	if p == nil {
		return nil
	}
	return p.parent
}

func (*TreeMap[K, V]) setColor(p *TreeMapEntry[K, V], color int) {
	if p != nil {
		p.color = color
	}
}

func (*TreeMap[K, V]) leftOf(p *TreeMapEntry[K, V]) *TreeMapEntry[K, V] {
	if p == nil {
		return nil
	}
	return p.left
}

func (*TreeMap[K, V]) rightOf(p *TreeMapEntry[K, V]) *TreeMapEntry[K, V] {
	if p == nil {
		return nil
	}
	return p.right
}

func (m *TreeMap[K, V]) rotateLeft(p *TreeMapEntry[K, V]) {
	if p != nil {
		r := p.right
		p.right = r.left
		if r.left != nil {
			r.left.parent = p
		}
		r.parent = p.parent
		if p.parent == nil {
			m.root = r
		} else if p.parent.left == p {
			p.parent.left = r
		} else {
			p.parent.right = r
		}
		r.left = p
		p.parent = r
	}
}

func (m *TreeMap[K, V]) rotateRight(p *TreeMapEntry[K, V]) {
	if p != nil {
		l := p.left
		p.left = l.right
		if l.right != nil {
			l.right.parent = p
		}
		l.parent = p.parent
		if p.parent == nil {
			m.root = l
		} else if p.parent.right == p {
			p.parent.right = l
		} else {
			p.parent.left = l
		}
		l.right = p
		p.parent = l
	}
}

func (m *TreeMap[K, V]) fixAfterInsertion(x *TreeMapEntry[K, V]) {
	x.color = treeMapColorRed

	for x != nil && x != m.root && x.parent.color == treeMapColorRed {
		if m.parentOf(x) == m.leftOf(m.parentOf(m.parentOf(x))) {
			y := m.rightOf(m.parentOf(m.parentOf(x)))
			if m.colorOf(y) == treeMapColorRed {
				m.setColor(m.parentOf(x), treeMapColorBlack)
				m.setColor(y, treeMapColorBlack)
				m.setColor(m.parentOf(m.parentOf(x)), treeMapColorRed)
				x = m.parentOf(m.parentOf(x))
			} else {
				if x == m.rightOf(m.parentOf(x)) {
					x = m.parentOf(x)
					m.rotateLeft(x)
				}
				m.setColor(m.parentOf(x), treeMapColorBlack)
				m.setColor(m.parentOf(m.parentOf(x)), treeMapColorRed)
				m.rotateRight(m.parentOf(m.parentOf(x)))
			}
		} else {
			y := m.leftOf(m.parentOf(m.parentOf(x)))
			if m.colorOf(y) == treeMapColorRed {
				m.setColor(m.parentOf(x), treeMapColorBlack)
				m.setColor(y, treeMapColorBlack)
				m.setColor(m.parentOf(m.parentOf(x)), treeMapColorRed)
				x = m.parentOf(m.parentOf(x))
			} else {
				if x == m.leftOf(m.parentOf(x)) {
					x = m.parentOf(x)
					m.rotateRight(x)
				}
				m.setColor(m.parentOf(x), treeMapColorBlack)
				m.setColor(m.parentOf(m.parentOf(x)), treeMapColorRed)
				m.rotateLeft(m.parentOf(m.parentOf(x)))
			}
		}
	}
	m.root.color = treeMapColorBlack
}

func (m *TreeMap[K, V]) fixAfterDeletion(x *TreeMapEntry[K, V]) {
	for x != m.root && m.colorOf(x) == treeMapColorBlack {
		if x == m.leftOf(m.parentOf(x)) {
			sib := m.rightOf(m.parentOf(x))

			if m.colorOf(sib) == treeMapColorRed {
				m.setColor(sib, treeMapColorBlack)
				m.setColor(m.parentOf(x), treeMapColorRed)
				m.rotateLeft(m.parentOf(x))
				sib = m.rightOf(m.parentOf(x))
			}

			if m.colorOf(m.leftOf(sib)) == treeMapColorBlack &&
				m.colorOf(m.rightOf(sib)) == treeMapColorBlack {
				m.setColor(sib, treeMapColorRed)
				x = m.parentOf(x)
			} else {
				if m.colorOf(m.rightOf(sib)) == treeMapColorBlack {
					m.setColor(m.leftOf(sib), treeMapColorBlack)
					m.setColor(sib, treeMapColorRed)
					m.rotateRight(sib)
					sib = m.rightOf(m.parentOf(x))
				}
				m.setColor(sib, m.colorOf(m.parentOf(x)))
				m.setColor(m.parentOf(x), treeMapColorBlack)
				m.setColor(m.rightOf(sib), treeMapColorBlack)
				m.rotateLeft(m.parentOf(x))
				x = m.root
			}
		} else { // symmetric
			sib := m.leftOf(m.parentOf(x))

			if m.colorOf(sib) == treeMapColorRed {
				m.setColor(sib, treeMapColorBlack)
				m.setColor(m.parentOf(x), treeMapColorRed)
				m.rotateRight(m.parentOf(x))
				sib = m.leftOf(m.parentOf(x))
			}

			if m.colorOf(m.rightOf(sib)) == treeMapColorBlack &&
				m.colorOf(m.leftOf(sib)) == treeMapColorBlack {
				m.setColor(sib, treeMapColorRed)
				x = m.parentOf(x)
			} else {
				if m.colorOf(m.leftOf(sib)) == treeMapColorBlack {
					m.setColor(m.rightOf(sib), treeMapColorBlack)
					m.setColor(sib, treeMapColorRed)
					m.rotateLeft(sib)
					sib = m.leftOf(m.parentOf(x))
				}
				m.setColor(sib, m.colorOf(m.parentOf(x)))
				m.setColor(m.parentOf(x), treeMapColorBlack)
				m.setColor(m.leftOf(sib), treeMapColorBlack)
				m.rotateRight(m.parentOf(x))
				x = m.root
			}
		}
	}

	m.setColor(x, treeMapColorBlack)
}

type OrderedKeyTreeMap[K generic.Ordered, V any] struct {
	TreeMap[K, V]
}

func NewOrderedKeyTreeMap[K generic.Ordered, V any]() *OrderedKeyTreeMap[K, V] {
	return &OrderedKeyTreeMap[K, V]{TreeMap: TreeMap[K, V]{cpr: OrderedCompare[K]}}
}

type ComparableKeyTreeMap[K Comparable, V any] struct {
	TreeMap[K, V]
}

func NewComparableKeyTreeMap[K Comparable, V any]() *ComparableKeyTreeMap[K, V] {
	return &ComparableKeyTreeMap[K, V]{}
}
