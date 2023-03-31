package pool

import "gitee.com/sy_183/common/option"

func WithLimit(limit int64) option.AnyOption {
	type limitSetter interface {
		setLimit(limit int64)
	}
	return option.AnyCustom(func(target any) {
		if setter, is := target.(limitSetter); is {
			setter.setLimit(limit)
		}
	})
}
