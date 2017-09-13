package main

func chain(funcs ...func()) func() {
	return func() {
		for i := range funcs {
			funcs[i]()
		}
	}
}
