package pool

import "gitee.com/sy_183/common/option"

type Pool[O any] interface {
	Get() O

	Put(obj O)
}

type PoolProvider[O any] func(new func(p Pool[O]) O, options ...option.AnyOption) Pool[O]
