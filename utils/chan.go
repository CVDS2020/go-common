package utils

func ChanAsyncPush[O any](c chan O, o O) {
	select {
	case c <- o:
	default:
		go func(c chan O) {
			c <- o
		}(c)
	}
}

func ChanTryPush[O any](c chan O, o O) bool {
	select {
	case c <- o:
		return true
	default:
		return false
	}
}

func ChanTryPop[O any](c chan O) (o O, ok bool) {
	select {
	case o = <-c:
		ok = true
	default:
	}
	return
}
