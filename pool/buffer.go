package pool

import (
	"fmt"
	"sync/atomic"
)

type Buffer struct {
	raw        []byte
	dataPool   []Data
	allocated  uint
	reversed   uint
	ref        atomic.Int64
	onReleased func(*Buffer)
}

func NewBuffer(size uint, reversed uint) *Buffer {
	return &Buffer{
		raw:      make([]byte, size),
		reversed: reversed,
	}
}

func (b *Buffer) SetOnReleased(onReleased func(*Buffer)) *Buffer {
	b.onReleased = onReleased
	return b
}

func (b *Buffer) Size() uint {
	return uint(len(b.raw))
}

func (b *Buffer) Remain() uint {
	return uint(len(b.raw)) - b.allocated
}

func (b *Buffer) Get() []byte {
	if b.Remain() < b.reversed {
		return nil
	}
	return b.raw[b.allocated:]
}

func (b *Buffer) Alloc(size uint) *Data {
	end := b.allocated + size
	if end > uint(len(b.raw)) {
		return nil
	}
	l := len(b.dataPool)
	if l < cap(b.dataPool) {
		b.dataPool = b.dataPool[:l+1]
	} else {
		b.dataPool = append(b.dataPool, Data{})
	}
	d := &b.dataPool[l]
	d.Data = b.raw[b.allocated:end]
	d.Reference = b
	b.AddRef()
	b.allocated += size
	return d
}

func (b *Buffer) Release() {
	c := b.ref.Add(-1)
	if c == 0 {
		b.allocated = 0
		b.dataPool = b.dataPool[:0]
		if onReleased := b.onReleased; onReleased != nil {
			onReleased(b)
		}
	} else if c < 0 {
		panic(fmt.Errorf("重复释放缓冲区(%p), 缓冲区引用计数[%d -> %d]", b, c+1, c))
	}
}

func (b *Buffer) AddRef() {
	c := b.ref.Add(1)
	if c <= 0 {
		panic(fmt.Errorf("无效的缓冲区(%p)引用计数[%d -> %d]", b, c-1, c))
	}
}

func (b *Buffer) Use() *Buffer {
	b.AddRef()
	return b
}
