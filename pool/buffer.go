package pool

type Buffer struct {
	raw        []byte
	dataPool   []Data
	allocated  uint
	reversed   uint
	ref        AtomicRef
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

func (b *Buffer) Used() []byte {
	return b.raw[:b.allocated]
}

func (b *Buffer) Unused() []byte {
	return b.raw[b.allocated:]
}

func (b *Buffer) Get() []byte {
	if b.Remain() < b.reversed {
		return nil
	}
	return b.Unused()
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

func (b *Buffer) Release() bool {
	if b.ref.Release() {
		b.allocated = 0
		b.dataPool = b.dataPool[:0]
		if onReleased := b.onReleased; onReleased != nil {
			onReleased(b)
		}
		return true
	}
	return false
}

func (b *Buffer) AddRef() {
	b.ref.AddRef()
}

func (b *Buffer) Use() *Buffer {
	b.AddRef()
	return b
}
