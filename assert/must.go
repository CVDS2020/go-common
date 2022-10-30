package assert

func MustSuccess(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[V any](v V, err error) V {
	MustSuccess(err)
	return v
}
