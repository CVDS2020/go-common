package pool

type DataPool interface {
	Alloc(len uint) *Data

	AllocCap(len, cap uint) *Data
}
