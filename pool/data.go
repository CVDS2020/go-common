package pool

import (
	"bytes"
	"encoding/binary"
	"gitee.com/sy_183/common/uns"
	"sync/atomic"
)

type Data struct {
	raw  []byte
	Data []byte
	ref  int64
	pool *DataPool
}

func NewData(data []byte) *Data {
	return &Data{
		raw:  data,
		Data: data,
		ref:  1,
	}
}

func (d *Data) Len() uint {
	return uint(len(d.Data))
}

func (d *Data) Cap() uint {
	return uint(cap(d.Data))
}

func (d *Data) Get(i uint) byte {
	return d.Data[i]
}

func (d *Data) Cut(start uint, end uint) *Data {
	d.Data = d.Data[start:end]
	return d
}

func (d *Data) CutFrom(start uint) *Data {
	d.Data = d.Data[start:]
	return d
}

func (d *Data) CutTo(end uint) *Data {
	d.Data = d.Data[:end]
	return d
}

func (d *Data) CutCap(start, end, cap uint) *Data {
	d.Data = d.Data[start:end:cap]
	return d
}

func (d *Data) CutCapTo(end, cap uint) *Data {
	d.Data = d.Data[:end:cap]
	return d
}

func (d *Data) PutUint16(off uint, v uint16) *Data {
	binary.BigEndian.PutUint16(d.Data[off:], v)
	return d
}

func (d *Data) PutUint32(off uint, v uint32) *Data {
	binary.BigEndian.PutUint32(d.Data[off:], v)
	return d
}

func (d *Data) PutUint64(off uint, v uint64) *Data {
	binary.BigEndian.PutUint64(d.Data[off:], v)
	return d
}

func (d *Data) PutLittleUint16(off uint, v uint16) *Data {
	binary.LittleEndian.PutUint16(d.Data[off:], v)
	return d
}

func (d *Data) PutLittleUint32(off uint, v uint32) *Data {
	binary.LittleEndian.PutUint32(d.Data[off:], v)
	return d
}

func (d *Data) PutLittleUint64(off uint, v uint64) *Data {
	binary.LittleEndian.PutUint64(d.Data[off:], v)
	return d
}

func (d *Data) Uint16(off uint) uint16 {
	return binary.BigEndian.Uint16(d.Data[off:])
}

func (d *Data) Uint32(off uint) uint32 {
	return binary.BigEndian.Uint32(d.Data[off:])
}

func (d *Data) Uint64(off uint) uint64 {
	return binary.BigEndian.Uint64(d.Data[off:])
}

func (d *Data) LittleUint16(off uint) uint16 {
	return binary.LittleEndian.Uint16(d.Data[off:])
}

func (d *Data) LittleUint32(off uint) uint32 {
	return binary.LittleEndian.Uint32(d.Data[off:])
}

func (d *Data) LittleUint64(off uint) uint64 {
	return binary.LittleEndian.Uint64(d.Data[off:])
}

func (d *Data) CopyFrom(off uint, b []byte) int {
	return copy(d.Data[off:], b)
}

func (d *Data) CopyStringFrom(off uint, s string) int {
	return copy(d.Data[off:], s)
}

func (d *Data) Equal(o []byte) bool {
	return bytes.Equal(d.Data, o)
}

func (d *Data) Compare(o []byte) int {
	return bytes.Compare(d.Data, o)
}

func (d *Data) Count(sep []byte) int {
	return bytes.Count(d.Data, sep)
}

func (d *Data) Contains(sub []byte) bool {
	return bytes.Contains(d.Data, sub)
}

func (d *Data) ContainsAny(chars string) bool {
	return bytes.ContainsAny(d.Data, chars)
}

func (d *Data) ContainsRune(r rune) bool {
	return bytes.ContainsRune(d.Data, r)
}

func (d *Data) IndexByte(c byte) int {
	return bytes.IndexByte(d.Data, c)
}

func (d *Data) LastIndex(sep []byte) int {
	return bytes.LastIndex(d.Data, sep)
}

func (d *Data) LastIndexByte(c byte) int {
	return bytes.LastIndexByte(d.Data, c)
}

func (d *Data) IndexRune(r rune) int {
	return bytes.IndexRune(d.Data, r)
}

func (d *Data) IndexAny(chars string) int {
	return bytes.IndexAny(d.Data, chars)
}

func (d *Data) LastIndexAny(chars string) int {
	return bytes.LastIndexAny(d.Data, chars)
}

func (d *Data) HasPrefix(prefix []byte) bool {
	return bytes.HasPrefix(d.Data, prefix)
}

func (d *Data) HasSuffix(suffix []byte) bool {
	return bytes.HasSuffix(d.Data, suffix)
}

func (d *Data) String() string {
	return uns.BytesToString(d.Data)
}

func (d *Data) Release() {
	if c := atomic.AddInt64(&d.ref, -1); c == 0 {
		if d.pool != nil {
			d.pool.put(d)
		}
	} else if c < 0 {
		panic("data repeat release")
	}
}

func (d *Data) AddRef() {
	if atomic.AddInt64(&d.ref, 1) <= 0 {
		panic("negative data reference count")
	}
}

func (d *Data) Use() *Data {
	d.AddRef()
	return d
}
