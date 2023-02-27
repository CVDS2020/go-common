package chanUtils

func AsyncPush[O any](c chan O, o O) {
	select {
	case c <- o:
	default:
		go func(c chan O) {
			c <- o
		}(c)
	}
}

func TryPush[O any](c chan O, o O) bool {
	select {
	case c <- o:
		return true
	default:
		return false
	}
}

func TryPop[O any](c chan O) (o O, ok bool) {
	select {
	case o = <-c:
		ok = true
	default:
	}
	return
}
