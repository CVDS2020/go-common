package lifecycle

// An GroupOption configures a Logger.
type GroupOption interface {
	apply(group *Group)
}

// groupOptionFunc wraps a func so it satisfies the Option interface.
type groupOptionFunc func(group *Group)

func (f groupOptionFunc) apply(group *Group) {
	f(group)
}

func PreStart(b bool) GroupOption {
	return groupOptionFunc(func(group *Group) {
		group.preStart = true
	})
}
