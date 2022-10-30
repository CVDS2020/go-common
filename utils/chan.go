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
