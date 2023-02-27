package pool

type Option interface {
	Apply(target any) any
}

type optionFunc func(target any) any
