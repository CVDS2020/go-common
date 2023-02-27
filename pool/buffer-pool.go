package pool

type BufferPool interface {
	Get() []byte

	Alloc(size uint) *Data
}
