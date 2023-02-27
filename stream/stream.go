package stream

type Stream[E any] interface {
	Filter()

	Foreach(func(e E))

	Range(func(e E) bool)
}
