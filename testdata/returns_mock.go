package testdata

type ReturnsMock struct {
	Returns struct {
		Generate struct {
			Number int
		}
	}
}

func (r *ReturnsMock) Generate() (number int) {
	return r.Returns.Generate.Number
}
