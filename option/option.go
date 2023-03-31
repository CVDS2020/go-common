package option

// Option 接口用于配置一个对象的属性
type Option[T, O any] interface {
	Apply(target T) O
}

type (
	AnyOption             Option[any, any]
	CustomOption[T any]   Option[T, any]
	ProviderOption[O any] Option[any, O]
)

// OptionFunc 使用配置函数配置一个对象的属性，此函数类型实现了 Option 接口
type Func[T, O any] func(target T) O
type AnyFunc Func[any, any]

func (f Func[T, O]) Apply(target T) O {
	return f(target)
}

// Custom 与 Func 类似，但是没有返回值
type Custom[T any] func(target T)
type AnyCustom = Custom[any]

func (f Custom[T]) Apply(target T) any {
	f(target)
	return nil
}

type CustomFunc[T, O any] func(target T)

func (f CustomFunc[T, O]) Apply(target T) (o O) {
	f(target)
	return
}

// Provider 通过函数提供一个配置属性用于配置对象
type Provider[O any] func() O
type AnyProvider = Provider[any]

func (p Provider[O]) Apply(target any) O {
	return p()
}

type ProviderFunc[T, O any] func() O

func (p ProviderFunc[T, O]) Apply(target any) O {
	return p()
}
