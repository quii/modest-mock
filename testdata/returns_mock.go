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
	call := struct {}{}
	r.Calls.Generate = append(r.Calls.Generate, call)
	return r.Returns.Generate[len(r.Calls.Generate)].Number
}
