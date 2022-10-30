package def

func Condition[X any](condition bool, yes X, no X) X {
	if condition {
		return yes
	} else {
		return no
	}
}

func ConditionFunc[X any](condition bool, yes func() X, no func() X) X {
	if condition {
		return yes()
	} else {
		return no()
	}
}

func SetDefault[X comparable](v X, def X) X {
	var zero X
	if v == zero {
		return def
	}
	return v
}

func SetDefaultIf[X any](v X, def X, condition func(v X) bool) X {
	if condition(v) {
		return def
	}
	return v
}

func SetDefaultEqual[X comparable](v X, def X, ref X) X {
	if v == ref {
		return def
	}
	return v
}

func SetterDefault[X comparable](v X, setter func() X) X {
	var zero X
	if v == zero {
		return setter()
	}
	return v
}

func SetterDefaultIf[X any](v X, setter func() X, condition func(v X) bool) X {
	if condition(v) {
		return setter()
	}
	return v
}

func SetterDefaultEqual[X comparable](v X, setter func() X, ref X) X {
	if v == ref {
		return setter()
	}
	return v
}

func SetDefaultP[X comparable](vp *X, def X) {
	var zero X
	if *vp == zero {
		*vp = def
	}
}

func SetDefaultPIf[X any](vp *X, def X, condition func(v X) bool) {
	if condition(*vp) {
		*vp = def
	}
}

func SetDefaultPEqual[X comparable](vp *X, def X, ref X) {
	if *vp == ref {
		*vp = def
	}
}

func SetterDefaultP[X comparable](vp *X, setter func() X) {
	var zero X
	if *vp == zero {
		*vp = setter()
	}
}

func SetterDefaultPIf[X any](vp *X, setter func() X, condition func(v X) bool) {
	if condition(*vp) {
		*vp = setter()
	}
}

func SetterDefaultPIfEqual[X comparable](vp *X, setter func() X, ref X) {
	if *vp == ref {
		*vp = setter()
	}
}
