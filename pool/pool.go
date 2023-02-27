package pool

type Pool[O any] interface {
	Get() O

	Put(obj O)
}

type PoolProvider[O any] func(new func(p Pool[O]) O) Pool[O]
