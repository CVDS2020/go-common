package io

import (
	"encoding/binary"
	"fmt"
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/uns"
	"io"
	"math"
	"unsafe"
)

type Data interface {
	Bytes() []byte
}

const (
	opTypeBytes = iota
	opTypeString
	opTypeEmbedded
	opTypeWriter

	bytesSize  = int(unsafe.Sizeof([]byte{}))
	stringSize = int(unsafe.Sizeof(""))
	byteSize   = 1
	uint16Size = 2
	uint32Size = 4
	uint64Size = 8
	writerSize = int(unsafe.Sizeof(io.WriterTo(nil)))

	BytesOpSize  = 1 + bytesSize
	StringOpSize = 1 + stringSize
)

type OpWriter struct {
	opBuf  []byte
	size   int
	typ    uint8
	typOff int
	embLen uint16
}

func NewOpWriter(buf []byte) *OpWriter {
	return &OpWriter{opBuf: buf[:0]}
}

func (w *OpWriter) SetBuf(buf []byte) *OpWriter {
	w.opBuf = append(buf[:0], w.opBuf...)
	return w
}

func (w *OpWriter) AppendBytes(bs []byte) *OpWriter {
	w.typ, w.typOff = opTypeBytes, len(w.opBuf)
	w.opBuf = append(append(w.opBuf, opTypeBytes),
		uns.MakeBytesUnchecked(uns.ToPointer(&bs), bytesSize, bytesSize)...)
	w.size += len(bs)
	return w
}

func (w *OpWriter) AppendString(s string) *OpWriter {
	w.typ, w.typOff = opTypeString, len(w.opBuf)
	w.opBuf = append(append(w.opBuf, opTypeString),
		uns.MakeBytesUnchecked(uns.ToPointer(&s), stringSize, stringSize)...)
	w.size += len(s)
	return w
}

func (w *OpWriter) AppendEmbeddedBytes(bs []byte) *OpWriter {
	if len(bs) > math.MaxUint16 {
		panic(errors.NewSizeOutOfRange("嵌入型数据", 0, math.MaxUint16, int64(len(bs)), true))
	}
	bsl := uint16(len(bs))
	if w.typ != opTypeEmbedded || w.embLen+bsl > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), bsl
		w.opBuf = append(binary.LittleEndian.AppendUint16(append(w.opBuf, opTypeEmbedded), bsl), bs...)
	} else {
		w.opBuf = append(w.opBuf, bs...)
		w.embLen += bsl
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += int(bsl)
	return w
}

func (w *OpWriter) AppendEmbeddedString(s string) *OpWriter {
	if len(s) > math.MaxUint16 {
		panic(errors.NewSizeOutOfRange("嵌入型字符串", 0, math.MaxUint16, int64(len(s)), true))
	}
	sl := uint16(len(s))
	if w.typ != opTypeEmbedded || w.embLen+sl > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), sl
		w.opBuf = append(binary.LittleEndian.AppendUint16(append(w.opBuf, opTypeEmbedded), sl), s...)
	} else {
		w.opBuf = append(w.opBuf, s...)
		w.embLen += sl
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += int(sl)
	return w
}

func (w *OpWriter) AppendByte(b byte) *OpWriter {
	if w.typ != opTypeEmbedded || w.embLen+byteSize > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), byteSize
		w.opBuf = append(binary.LittleEndian.AppendUint16(append(w.opBuf, opTypeEmbedded), byteSize), b)
	} else {
		w.opBuf = append(w.opBuf, b)
		w.embLen += byteSize
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += byteSize
	return w
}

func (w *OpWriter) AppendUint16(u uint16) *OpWriter {
	if w.typ != opTypeEmbedded || w.embLen+uint16Size > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), uint16Size
		w.opBuf = binary.LittleEndian.AppendUint16(binary.LittleEndian.AppendUint16(
			append(w.opBuf, opTypeEmbedded), uint16Size), u)
	} else {
		w.opBuf = binary.LittleEndian.AppendUint16(w.opBuf, u)
		w.embLen += uint16Size
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += uint16Size
	return w
}

func (w *OpWriter) AppendUint32(u uint32) *OpWriter {
	if w.typ != opTypeEmbedded || w.embLen+uint32Size > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), uint32Size
		w.opBuf = binary.LittleEndian.AppendUint32(binary.LittleEndian.AppendUint16(
			append(w.opBuf, opTypeEmbedded), uint32Size), u)
	} else {
		w.opBuf = binary.LittleEndian.AppendUint32(w.opBuf, u)
		w.embLen += uint32Size
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += uint32Size
	return w
}

func (w *OpWriter) AppendUint64(u uint64) *OpWriter {
	if w.typ != opTypeEmbedded || w.embLen+uint64Size > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), uint64Size
		w.opBuf = binary.LittleEndian.AppendUint64(binary.LittleEndian.AppendUint16(
			append(w.opBuf, opTypeEmbedded), uint64Size), u)
	} else {
		w.opBuf = binary.LittleEndian.AppendUint64(w.opBuf, u)
		w.embLen += uint64Size
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += uint64Size
	return w
}

func OpWriterAppendEmbedded[T any](w *OpWriter, emb T) *OpWriter {
	ss := unsafe.Sizeof(emb)
	if ss > math.MaxUint16 {
		panic(errors.NewSizeOutOfRange("嵌入型数据结构", 0, math.MaxUint16, int64(ss), true))
	}
	s := uint16(ss)
	if w.typ != opTypeEmbedded || w.embLen+s > math.MaxUint16 {
		w.typ, w.typOff, w.embLen = opTypeEmbedded, len(w.opBuf), s
		w.opBuf = append(binary.LittleEndian.AppendUint16(append(w.opBuf, opTypeEmbedded), s),
			uns.MakeBytesUnchecked(uns.ToPointer(&emb), int(s), int(s))...)
	} else {
		w.opBuf = append(w.opBuf, uns.MakeBytesUnchecked(uns.ToPointer(&emb), int(s), int(s))...)
		w.embLen += s
		*uns.ConvertPointer[byte, uint16](&w.opBuf[w.typOff+1]) = w.embLen
	}
	w.size += int(s)
	return w
}

func (w *OpWriter) AppendWriter(writer io.WriterTo, size int) *OpWriter {
	if size < 0 {
		panic(errors.NewSizeOutOfRange("Writer", 0, math.MinInt64, int64(size), true))
	}
	w.typ, w.typOff = opTypeWriter, len(w.opBuf)
	w.opBuf = append(binary.LittleEndian.AppendUint64(append(w.opBuf, opTypeWriter), uint64(size)),
		uns.MakeBytesUnchecked(uns.ToPointer(&writer), writerSize, writerSize)...)
	w.size += size
	return w
}

func (w *OpWriter) Size() int {
	return w.size
}

func (w *OpWriter) Bytes() []byte {
	writer := Writer{Buf: make([]byte, w.Size())}
	n, _ := w.WriteTo(&writer)
	return writer.Buf[:n]
}

func (w *OpWriter) Read(p []byte) (n int, err error) {
	writer := Writer{Buf: p}
	_n, _ := w.WriteTo(&writer)
	return int(_n), io.EOF
}

func (w *OpWriter) WriteTo(_w io.Writer) (n int64, err error) {
	defer func() { err = HandleRecovery(recover()) }()
	p := w.opBuf
	for len(p) > 0 {
		switch p[0] {
		case opTypeBytes:
			WritePanic(_w, uns.FromPointer[byte, []byte](&p[1]), &n)
			p = p[1+bytesSize:]
		case opTypeString:
			WritePanic(_w, uns.StringToBytes(uns.FromPointer[byte, string](&p[1])), &n)
			p = p[1+stringSize:]
		case opTypeEmbedded:
			s := int(uns.FromPointer[byte, uint16](&p[1]))
			WritePanic(_w, uns.MakeBytesUnchecked(uns.ToPointer(&p[1+uint16Size]), s, s), &n)
			p = p[1+uint16Size+s:]
		case opTypeWriter:
			WriteToPanic(uns.FromPointer[byte, io.WriterTo](&p[1+uint64Size]), _w, &n)
			p = p[1+uint64Size+writerSize:]
		default:
			panic(fmt.Errorf("非法的操作类型(%d)", p[0]))
		}
	}
	return
}

func (w *OpWriter) Reset() {
	w.opBuf = w.opBuf[:0]
	w.typ = opTypeBytes
	w.typOff = 0
	w.embLen = 0
}

func (w *OpWriter) Clear() {
	w.Reset()
}
