package container

import (
	"fmt"
)

type LinkedListEntry[E any] struct {
	elem E
	next *LinkedListEntry[E]
	prev *LinkedListEntry[E]
}

func NewLinkedListEntry[E any](e E) *LinkedListEntry[E] {
	return &LinkedListEntry[E]{elem: e}
}

func (e *LinkedListEntry[E]) Value() E {
	return e.elem
}

func (e *LinkedListEntry[E]) Next() *LinkedListEntry[E] {
	if e.next == nil || e.next.next == nil {
		return nil
	}
	return e.next
}

func (e *LinkedListEntry[E]) Prev() *LinkedListEntry[E] {
	if e.prev == nil || e.prev.prev == nil {
		return nil
	}
	return e.prev
}

func (e *LinkedListEntry[E]) Removed() bool {
	return e.next == nil && e.prev == nil
}

func (e *LinkedListEntry[E]) String() string {
	return fmt.Sprint(e.elem)
}

type LinkedList[E any] struct {
	len  int
	head *LinkedListEntry[E]
	tail *LinkedListEntry[E]
}

func NewLinkedList[E any]() *LinkedList[E] {
	l := new(LinkedList[E])
	l.head = new(LinkedListEntry[E])
	l.tail = new(LinkedListEntry[E])
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

func (l *LinkedList[E]) addLink(e, p *LinkedListEntry[E]) {
	e.prev = p
	e.next = p.next
	p.next.prev = e
	p.next = e
}

func (l *LinkedList[E]) removeLink(e *LinkedListEntry[E]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next, e.prev = nil, nil
}

func (l *LinkedList[E]) Len() int {
	return l.len
}

func (l *LinkedList[E]) Head() (e E, exist bool) {
	return l.head.next.elem, l.len != 0
}

func (l *LinkedList[E]) Tail() (e E, exist bool) {
	return l.tail.prev.elem, l.len != 0
}

func (l *LinkedList[E]) HeadEntry() *LinkedListEntry[E] {
	if l.len == 0 {
		return nil
	}
	return l.head.next
}

func (l *LinkedList[E]) TailEntry() *LinkedListEntry[E] {
	if l.len == 0 {
		return nil
	}
	return l.tail.prev
}

func (l *LinkedList[E]) insert(e, p *LinkedListEntry[E]) {
	l.addLink(e, p)
	l.len++
}

func (l *LinkedList[E]) AddTail(e E) (entry *LinkedListEntry[E]) {
	entry = NewLinkedListEntry(e)
	l.insert(entry, l.tail.prev)
	return entry
}

func (l *LinkedList[E]) AddEntryTail(entry *LinkedListEntry[E]) {
	if !entry.Removed() {
		panic("linked list entry not removed")
	}
	l.insert(entry, l.tail.prev)
}

func (l *LinkedList[E]) AddHead(e E) (entry *LinkedListEntry[E]) {
	entry = NewLinkedListEntry(e)
	l.insert(entry, l.head)
	return entry
}

func (l *LinkedList[E]) AddEntryHead(entry *LinkedListEntry[E]) {
	if !entry.Removed() {
		panic("linked list entry not removed")
	}
	l.insert(entry, l.head)
}

func (l *LinkedList[E]) InsertBack(e E, ref *LinkedListEntry[E]) (entry *LinkedListEntry[E]) {
	entry = NewLinkedListEntry(e)
	l.insert(entry, ref)
	return entry
}

func (l *LinkedList[E]) InsertEntryBack(entry, ref *LinkedListEntry[E]) {
	if !entry.Removed() {
		panic("linked list entry not removed")
	}
	if ref.Removed() {
		panic("reference linked list entry removed")
	}
	l.insert(entry, ref)
}

func (l *LinkedList[E]) InsertFront(e E, ref *LinkedListEntry[E]) (entry *LinkedListEntry[E]) {
	entry = NewLinkedListEntry(e)
	l.insert(entry, ref.prev)
	return entry
}

func (l *LinkedList[E]) InsertEntryFront(entry, ref *LinkedListEntry[E]) {
	if !entry.Removed() {
		panic("linked list entry not removed")
	}
	if ref.Removed() {
		panic("reference linked list entry removed")
	}
	l.insert(entry, ref.prev)
}

func (l *LinkedList[E]) RemoveTail() (entry *LinkedListEntry[E]) {
	if l.len == 0 {
		return nil
	}
	entry = l.tail.prev
	l.removeLink(entry)
	return
}

func (l *LinkedList[E]) RemoveHead() (entry *LinkedListEntry[E]) {
	if l.len == 0 {
		return nil
	}
	entry = l.head.next
	l.removeLink(entry)
	return
}

func (l *LinkedList[E]) Clear() {
	for entry := l.head.next; entry != l.tail; {
		next := entry.next
		entry.prev, entry.next = nil, nil
		entry = next
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	l.len = 0
}
