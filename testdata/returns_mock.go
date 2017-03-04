package testdata

type ReturnsMock struct {
	Calls struct {
		Generate []struct {
		}
	}

	Returns struct {
		Generate []struct {
			Number int
		}
	}
}

func (r *ReturnsMock) Generate() (Number int) {
	return r.Returns.Generate[len(r.Calls.Generate)].Number
}
