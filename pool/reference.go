package pool

type Reference interface {
	Release()

	AddRef()
}
