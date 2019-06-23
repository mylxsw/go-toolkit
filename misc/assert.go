package misc

func AssertError(err error) {
	if err != nil {
		panic(err)
	}
}
