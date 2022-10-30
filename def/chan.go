package def

func MakeChan[E any](c chan E, size int) chan E {
	if c == nil {
		return make(chan E, size)
	}
	return c
}

func SetChan[E any](c chan E, def chan E) chan E {
	if c == nil {
		return def
	}
	return c
}

func SetterChan[E any](c chan E, setter func() chan E) chan E {
	if c == nil {
		return setter()
	}
	return c
}

func MakeChanP[E any](cp *chan E, size int) {
	if *cp == nil {
		*cp = make(chan E, size)
	}
}

func SetChanP[E any](cp *chan E, def chan E) {
	if *cp == nil {
		*cp = def
	}
}

func SetterChanP[E any](cp *chan E, setter func() chan E) {
	if *cp == nil {
		*cp = setter()
	}
}
