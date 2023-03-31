package io

import (
	"encoding/binary"
	"gitee.com/sy_183/common/uns"
	"io"
)

type Writer struct {
	Buf []byte
	off int
}

func (w *Writer) Len() int {
	return w.off
}

func (w *Writer) Remain() int {
	return len(w.Buf) - w.off
}

func (w *Writer) Bytes() []byte {
	return w.Buf[:w.off]
}

func (w *Writer) ResetBuf(buf []byte) {
	w.Buf = buf
	w.off = 0
}

func (w *Writer) Reset() {
	w.off = 0
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.WriteBytes(p), nil
}

func (w *Writer) WriteBytes(p []byte) (n int) {
	n = copy(w.Buf[w.off:], p)
	w.off += n
	return
}

func (w *Writer) WriteBytesAnd(p []byte) *Writer {
	w.WriteBytes(p)
	return w
}

func (w *Writer) WriteString(s string) (n int) {
	n = copy(w.Buf[w.off:], s)
	w.off += n
	return
}

func (w *Writer) WriteStringAnd(s string) *Writer {
	w.WriteString(s)
	return w
}

func (w *Writer) WriteByte(b byte) *Writer {
	w.Buf[w.off] = b
	w.off++
	return w
}

func (w *Writer) WriteUint16(u uint16) *Writer {
	binary.BigEndian.PutUint16(w.Buf[w.off:], u)
	w.off += 2
	return w
}

func (w *Writer) WriteUint32(u uint32) *Writer {
	binary.BigEndian.PutUint32(w.Buf[w.off:], u)
	w.off += 4
	return w
}

func (w *Writer) WriteUint64(u uint64) *Writer {
	binary.BigEndian.PutUint64(w.Buf[w.off:], u)
	w.off += 8
	return w
}

func Write(w io.Writer, p []byte, np *int64) error {
	if len(p) > 0 {
		n, err := w.Write(p)
		*np += int64(n)
		return err
	}
	return nil
}

func WritePanic(w io.Writer, p []byte, np *int64) {
	if err := Write(w, p, np); err != nil {
		panic(err)
	}
}

func WriteString(w io.Writer, s string, np *int64) error {
	return Write(w, uns.StringToBytes(s), np)
}

func WriteStringPanic(w io.Writer, s string, np *int64) {
	if err := WriteString(w, s, np); err != nil {
		panic(err)
	}
}

func WriteTo(wt io.WriterTo, w io.Writer, np *int64) error {
	n, err := wt.WriteTo(w)
	*np += n
	return err
}

func WriteToPanic(wt io.WriterTo, w io.Writer, np *int64) {
	if err := WriteTo(wt, w, np); err != nil {
		panic(err)
	}
}

func WriteAndReset(w *Writer, iow io.Writer, np *int64) error {
	if w.Len() > 0 {
		if err := Write(iow, w.Bytes(), np); err != nil {
			return err
		}
		w.Reset()
	}
	return nil
}

func WriteAndResetPanic(w *Writer, iow io.Writer, np *int64) {
	if err := WriteAndReset(w, iow, np); err != nil {
		panic(err)
	}
}

func HandleRecovery(e any) error {
	if e != nil {
		if err, is := e.(error); is {
			return err
		}
		panic(e)
	}
	return nil
}
