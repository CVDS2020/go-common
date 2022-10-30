package pool

type Buf Data

func (b *Buf) Write(p []byte) (n int, err error) {
	if len(b.Data)+len(p) <= cap(b.Data) {
		b.Data = append(b.Data, p...)
		return len(p), nil
	}
	b.Data = append(b.Data, p...)
	b.raw = b.Data
	b.pool = nil
	return len(p), nil
}
