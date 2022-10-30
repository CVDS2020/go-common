package uns

type BufferWriter struct {
	buf []byte
	st  *BytesStruct
}

func NewBufferWriter(buf []byte) *BufferWriter {
	w := &BufferWriter{
		buf: buf,
	}
	w.st = BytesStructOfPtr(&w.buf)
	return w
}

func (w *BufferWriter) Write(data []byte) {
	if len(data)+len(w.buf) > cap(w.buf) {
		panic("length out of range")
	}
	ol := len(w.buf)
	w.st.Len += len(data)
	copy(w.buf[ol:], data)
}

func (w *BufferWriter) WriteString(s string) {
	if len(s)+len(w.buf) > cap(w.buf) {
		panic("length out of range")
	}
	ol := len(w.buf)
	w.st.Len += len(s)
	copy(w.buf[ol:], s)
}

func (w *BufferWriter) WriteByte(b byte) {
	if len(w.buf) >= cap(w.buf) {
		panic("length out of range")
	}
	w.st.Len += 1
	w.buf[w.st.Len-1] = b
}
